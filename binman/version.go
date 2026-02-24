package binman

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// VersionManifest tracks installed binary versions.
type VersionManifest struct {
	Versions  map[string]string `json:"versions"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// NewManifest creates an empty version manifest.
func NewManifest() *VersionManifest {
	return &VersionManifest{
		Versions: make(map[string]string),
	}
}

// LoadManifest loads a version manifest from path. Returns an empty manifest on ENOENT.
func LoadManifest(path string) (*VersionManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewManifest(), nil
		}
		return nil, err
	}

	var m VersionManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if m.Versions == nil {
		m.Versions = make(map[string]string)
	}
	return &m, nil
}

// Save writes the manifest to path.
func (m *VersionManifest) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}

	m.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetVersion returns the installed version for a binary name.
func (m *VersionManifest) GetVersion(name string) string {
	return m.Versions[name]
}

// SetVersion sets the installed version for a binary name.
func (m *VersionManifest) SetVersion(name, version string) {
	m.Versions[name] = version
}

// CompareVersions compares two version strings.
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2.
// Handles both semver (v1.23.0) and date-based (v2026.01.29) versions.
func CompareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	if v1 == "" && v2 == "" {
		return 0
	}
	if v1 == "" {
		return -1
	}
	if v2 == "" {
		return 1
	}

	if isDevVersion(v1) && !isDevVersion(v2) {
		return -1
	}
	if !isDevVersion(v1) && isDevVersion(v2) {
		return 1
	}
	if isDevVersion(v1) && isDevVersion(v2) {
		return 0
	}

	// Date-based versions (YYYY.MM.DD)
	datePattern := regexp.MustCompile(`^\d{4}\.\d{2}\.\d{2}$`)
	if datePattern.MatchString(v1) && datePattern.MatchString(v2) {
		return strings.Compare(v1, v2)
	}

	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 int
		if i < len(parts1) {
			p1 = parts1[i]
		}
		if i < len(parts2) {
			p2 = parts2[i]
		}
		if p1 < p2 {
			return -1
		}
		if p1 > p2 {
			return 1
		}
	}

	return 0
}

// IsNewer returns true if candidate is newer than current.
func IsNewer(current, candidate string) bool {
	return CompareVersions(current, candidate) < 0
}

func isDevVersion(v string) bool {
	switch v {
	case "dev", "unknown", "latest":
		return true
	}
	for _, c := range v {
		if c >= '0' && c <= '9' {
			return false
		}
	}
	return true
}

func parseVersion(v string) []int {
	re := regexp.MustCompile(`[^\d]+`)
	parts := re.Split(v, -1)

	var result []int
	for _, p := range parts {
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		result = append(result, n)
	}
	return result
}
