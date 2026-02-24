// Package binman provides binary management: download, install, update, and version tracking.
package binman

// ArchMapping maps platform identifiers (e.g. "linux/amd64") to replacement strings.
type ArchMapping map[string]string

// BinaryDef defines a downloadable binary.
type BinaryDef struct {
	Name          string
	EnvOverride   string                   // env var to override binary path
	URLPattern    string                   // supports {version}, {os}, {arch}, {ext}
	PinnedVersion string                   // version to install ("latest" = head of release)
	Archive       bool                     // true if download is tar.xz
	ChecksumURL   string                   // URL pattern for checksum file (empty = skip)
	Platforms     map[string][]string      // os -> []arch; nil = all
	SkipUpdate    bool                     // exclude from update checks
	ArchMappings  map[string]ArchMapping   // custom placeholder -> platform mapping
}

// ProgressFunc reports download progress.
type ProgressFunc func(downloaded, total int64)

// StatusFunc reports status messages.
type StatusFunc func(message string)
