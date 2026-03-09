package binman

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Option configures a Manager.
type Option func(*Manager)

// WithSystemPaths sets system paths to search for binaries.
func WithSystemPaths(paths []string) Option {
	return func(m *Manager) {
		m.systemPaths = paths
	}
}

// Manager handles binary installation and resolution.
type Manager struct {
	binDir      string
	systemPaths []string
}

// NewManager creates a new binary manager.
func NewManager(binDir string, opts ...Option) *Manager {
	m := &Manager{binDir: binDir}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// BinDir returns the managed binary directory.
func (m *Manager) BinDir() string {
	return m.binDir
}

// EnsureDir creates the bin directory if it doesn't exist.
func (m *Manager) EnsureDir() error {
	return os.MkdirAll(m.binDir, 0750)
}

// ResolvePath returns the path to a binary, checking: env override > system paths > binDir.
func (m *Manager) ResolvePath(def BinaryDef) (string, error) {
	// Check environment override
	if def.EnvOverride != "" {
		if envPath := os.Getenv(def.EnvOverride); envPath != "" {
			if _, err := os.Stat(envPath); err == nil {
				return envPath, nil
			}
		}
	}

	// Check system paths
	for _, dir := range m.systemPaths {
		p := filepath.Join(dir, def.Name)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	// Check bin directory
	p := filepath.Join(m.binDir, def.Name)
	if _, err := os.Stat(p); err == nil {
		return p, nil
	}

	return "", fmt.Errorf("binary %s not found", def.Name)
}

// IsInstalled checks if a binary is available via any path.
func (m *Manager) IsInstalled(def BinaryDef) bool {
	_, err := m.ResolvePath(def)
	return err == nil
}

// IsPlatformSupported checks if the current platform is supported by the binary def.
func (m *Manager) IsPlatformSupported(def BinaryDef) bool {
	if def.Platforms == nil {
		return true
	}
	arches, ok := def.Platforms[runtime.GOOS]
	if !ok {
		return false
	}
	if arches == nil {
		return true
	}
	for _, a := range arches {
		if a == runtime.GOARCH {
			return true
		}
	}
	return false
}

// EnsureInstalled resolves an existing binary or downloads it.
func (m *Manager) EnsureInstalled(def BinaryDef, fn ProgressFunc) (string, error) {
	if path, err := m.ResolvePath(def); err == nil {
		return path, nil
	}
	if err := m.Download(def, def.PinnedVersion, fn); err != nil {
		return "", err
	}
	return m.ResolvePath(def)
}

// Download downloads, verifies, and installs a binary.
func (m *Manager) Download(def BinaryDef, version string, fn ProgressFunc) error {
	url := m.BuildURL(def, version)

	tmpPath, err := httpDownload(url, fn)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	// Verify checksum if configured
	if def.ChecksumURL != "" {
		checksumURL := m.buildChecksumURL(def, version)
		expectedFilename := filepath.Base(url)
		if err := fetchAndVerifyChecksum(checksumURL, tmpPath, expectedFilename); err != nil {
			return fmt.Errorf("checksum verification failed for %s: %w", def.Name, err)
		}
	}

	// Extract from archive if needed
	installPath := tmpPath
	if def.IsArchive() {
		var extractedPath string
		switch def.GetArchiveType() {
		case "zip":
			extractedPath, err = extractZip(tmpPath, def.Name)
		case "tar.xz":
			extractedPath, err = extractTarXz(tmpPath, def.Name)
		default:
			err = fmt.Errorf("unsupported archive type: %s", def.GetArchiveType())
		}
		if err != nil {
			return fmt.Errorf("failed to extract %s from archive: %w", def.Name, err)
		}
		defer os.Remove(extractedPath)
		installPath = extractedPath
	}

	// Install to binDir
	if err := m.EnsureDir(); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	destPath := filepath.Join(m.binDir, def.Name)

	data, err := os.ReadFile(installPath)
	if err != nil {
		return fmt.Errorf("failed to read downloaded binary: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0755); err != nil {
		return fmt.Errorf("failed to install %s: %w", def.Name, err)
	}

	return nil
}

// Remove removes a binary from the managed bin directory.
func (m *Manager) Remove(def BinaryDef) error {
	p := filepath.Join(m.binDir, def.Name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil // already gone
	}
	return os.Remove(p)
}

// BuildURL expands a BinaryDef's URLPattern for a given version.
func (m *Manager) BuildURL(def BinaryDef, version string) string {
	return buildURL(def, version)
}

func (m *Manager) buildChecksumURL(def BinaryDef, version string) string {
	// Use the same expansion logic on ChecksumURL
	tmpDef := def
	tmpDef.URLPattern = def.ChecksumURL
	return buildURL(tmpDef, version)
}

// buildURL expands placeholders in a URL pattern.
func buildURL(def BinaryDef, version string) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	r := strings.NewReplacer(
		"{version}", version,
		"{os}", osName,
		"{arch}", arch,
		"{ext}", ext,
	)
	url := r.Replace(def.URLPattern)

	// Apply custom arch mappings
	platform := osName + "/" + arch
	for placeholder, mapping := range def.ArchMappings {
		if replacement, ok := mapping[platform]; ok {
			url = strings.ReplaceAll(url, "{"+placeholder+"}", replacement)
		}
	}

	return url
}
