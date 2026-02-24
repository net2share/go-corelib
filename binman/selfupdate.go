package binman

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// SelfUpdateConfig configures a self-update operation.
type SelfUpdateConfig struct {
	Repo       string     // "owner/name"
	URLPattern string     // with {version}, {os}, {arch}
	StatusFn   StatusFunc
}

// CheckSelfUpdate checks if a newer version is available.
func CheckSelfUpdate(repo, currentVersion string) (latestVersion string, available bool, err error) {
	release, err := GetLatestRelease(repo)
	if err != nil {
		return "", false, err
	}
	if IsNewer(currentVersion, release.TagName) {
		return release.TagName, true, nil
	}
	return release.TagName, false, nil
}

// SelfUpdate downloads and replaces the running binary.
func SelfUpdate(cfg SelfUpdateConfig, newVersion string) error {
	if cfg.StatusFn != nil {
		cfg.StatusFn(fmt.Sprintf("Downloading %s...", newVersion))
	}

	url := expandURL(cfg.URLPattern, newVersion)

	tmpFile, err := os.CreateTemp("", "selfupdate-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	tmpPath, err = httpDownload(url, nil)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	if err := os.Chmod(tmpPath, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Resolve the actual executable path (follow symlinks)
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find current executable: %w", err)
	}
	targetPath, err := filepath.EvalSymlinks(currentExe)
	if err != nil {
		targetPath = currentExe
	}

	if cfg.StatusFn != nil {
		cfg.StatusFn("Installing new version...")
	}

	// Remove old binary first (Linux keeps it in memory for current process)
	os.Remove(targetPath)

	// Move new binary into place
	if err := os.Rename(tmpPath, targetPath); err != nil {
		// Cross-device fallback: copy
		if err := copyFile(tmpPath, targetPath); err != nil {
			return fmt.Errorf("failed to install binary: %w", err)
		}
	}

	if cfg.StatusFn != nil {
		cfg.StatusFn(fmt.Sprintf("Updated to %s", newVersion))
	}

	return nil
}

// expandURL replaces {version}, {os}, {arch} in a URL pattern.
func expandURL(pattern, version string) string {
	r := strings.NewReplacer(
		"{version}", version,
		"{os}", runtime.GOOS,
		"{arch}", runtime.GOARCH,
	)
	return r.Replace(pattern)
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}
