# go-corelib

Shared Go library providing common utilities for CLI applications.

## Packages

### osdetect

OS detection, package management, and system utilities.

```go
import "github.com/net2share/go-corelib/osdetect"

// Detect OS
info, err := osdetect.Detect()
fmt.Println(info.ID)             // "fedora"
fmt.Println(info.PrettyName)     // "Fedora Linux 39"
fmt.Println(info.PackageManager) // "dnf"

// Install packages
err = info.InstallPackage("nginx")

// System checks
if osdetect.IsRoot() { ... }
if osdetect.HasSystemd() { ... }
if osdetect.HasIPv6() { ... }

// Require root (returns error if not root)
if err := osdetect.RequireRoot(); err != nil {
    return err  // "this program must be run as root"
}

// Get system info
arch := osdetect.GetArch()        // "amd64", "arm64", etc.
iface, _ := osdetect.GetDefaultInterface()
port := osdetect.DetectSSHPort()  // "22"
```

### tui

Terminal UI utilities using [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) and [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea).

```go
import "github.com/net2share/go-corelib/tui"

// Print messages with consistent styling
tui.PrintSuccess("Operation completed!")
tui.PrintError("Something went wrong")
tui.PrintWarning("Proceed with caution")
tui.PrintInfo("Note: this is informational")
tui.PrintStatus("Processing...")
tui.PrintStep(1, 5, "Installing dependencies")

// Display a styled box
tui.PrintBox("Configuration", []string{
    "Domain: example.com",
    "Mode: production",
})
tui.PrintBoxSimple([]string{"Line 1", "Line 2"})

// Format text
kv := tui.KV("Domain", "example.com")  // styled key-value pair
header := tui.Header("Section Title")
code := tui.Code("go build")
value := tui.Value("some-value")
muted := tui.Muted("subdued text")

// Display app banner
tui.PrintBanner(tui.BannerConfig{
    AppName:   "My App",
    Version:   "1.0.0",
    BuildTime: "2024-01-01",
    ASCII:     asciiArt,  // optional
})

// Simple banner without ASCII art
tui.PrintSimpleBanner("My App", "1.0.0", "2024-01-01")

// Terminal helpers
tui.WaitForEnter()   // "Press Enter to continue..."
tui.ClearLine()      // Clear current terminal line
tui.ClearScreen()    // Clear entire terminal screen

// Inline progress bar
tui.PrintProgress(downloaded, total)
tui.PrintProgressWithLabel("Downloading", downloaded, total)

// Session management (alternate screen buffer)
tui.BeginSession()
defer tui.EndSession()
if tui.InSession() { ... }

// Set global app info (shown in fullscreen component footers)
tui.SetAppInfo("My App", "1.0.0", "2024-01-01")
```

#### Full-Screen Interactive Components

Interactive menus, confirmations, and inputs using bubbletea.

```go
// Selection menu
selected, err := tui.RunMenu(tui.MenuConfig{
    Header:      "Status line",       // optional header above title
    Title:       "Choose an option",
    Description: "Use arrow keys to navigate",
    Options: []tui.MenuOption{
        {Label: "Option 1", Value: "opt1"},
        {Label: "Option 2", Value: "opt2"},
        {Separator: true},             // non-selectable divider
        {Label: "Option 3", Value: "opt3"},
    },
})

// Simple menu with just labels
selected, err := tui.RunMenuSimple("Select Mode", []string{"Development", "Production"})

// RunSelect is an alias for RunMenu
selected, err := tui.RunSelect(tui.SelectConfig{ ... })

// Confirmation dialog
confirmed, err := tui.RunConfirm(tui.ConfirmConfig{
    Title:       "Confirm Action",
    Description: "This will delete all data",
    Default:     false,
})

// Text input
value, ok, err := tui.RunInput(tui.InputConfig{
    Title:       "Enter Name",
    Placeholder: "your-name",
    Password:    false,  // set true for password input
})

// Display message (waits for OK)
tui.ShowMessage(tui.AppMessage{
    Type:    "success",  // success, error, warning, info
    Message: "Operation completed!",
})

// Display a non-interactive list
tui.ShowList(tui.ListConfig{
    Title:     "Installed Packages",
    Items:     []string{"nginx", "curl", "jq"},
    EmptyText: "No packages installed",
})

// Display structured info with sections and columns
tui.ShowInfo(tui.InfoConfig{
    Title: "System Info",
    Sections: []tui.InfoSection{
        {
            Title: "Network",
            Rows: []tui.InfoRow{
                {Key: "IP", Value: "192.168.1.1"},
                {Columns: []string{"eth0", "UP", "1000Mbps"}},
            },
        },
    },
})

// Full-screen progress spinner
progress := tui.StartProgress(tui.ProgressConfig{
    Title:   "Installing",
    Message: "Please wait...",
})
progress.Update("Still working...")
progress.Done()

// Multi-line progress view with scrolling
pv := tui.NewProgressView("Installing")
pv.AddInfo("Downloading package...")
pv.AddStatus("Extracting files...")
pv.AddSuccess("Package installed")
pv.AddWarning("Optional dependency missing")
pv.AddError("Failed to configure")
pv.AddText("Plain text line")
pv.Done()       // show "press any key" prompt
pv.Dismiss()    // close immediately without prompt
```

#### Theme Colors

The tui package uses a consistent color theme:

| Color | Usage |
|-------|-------|
| `Theme.Primary` | Cyan - Primary actions, highlights |
| `Theme.Secondary` | Magenta - Secondary elements |
| `Theme.Success` | Green - Success messages |
| `Theme.Error` | Red - Error messages |
| `Theme.Warning` | Yellow - Warning messages |
| `Theme.Info` | Blue - Informational messages |
| `Theme.Muted` | Gray - Subdued text |

### binman

Binary management: download, install, version tracking, checksum verification, and self-update.

```go
import "github.com/net2share/go-corelib/binman"

// Define binaries
def := binman.BinaryDef{
    Name:          "my-tool",
    EnvOverride:   "MY_TOOL_PATH",           // env var to override path resolution
    URLPattern:    "https://github.com/org/repo/releases/download/{version}/my-tool-{os}-{arch}",
    PinnedVersion: "v1.2.0",
    ChecksumURL:   "https://github.com/org/repo/releases/download/{version}/checksums.sha256",
}

// For archived binaries (tar.xz)
archiveDef := binman.BinaryDef{
    Name:          "sslocal",
    URLPattern:    "https://example.com/releases/{version}/pkg-{version}.{platform}.tar.xz",
    PinnedVersion: "v1.21.2",
    Archive:       true,
    ChecksumURL:   "https://example.com/releases/{version}/pkg-{version}.{platform}.tar.xz.sha256",
    ArchMappings: map[string]binman.ArchMapping{
        "platform": {
            "linux/amd64":  "x86_64-unknown-linux-gnu",
            "linux/arm64":  "aarch64-unknown-linux-gnu",
            "darwin/arm64": "aarch64-apple-darwin",
        },
    },
}
```

#### Manager

Handles binary resolution, download, and installation.

```go
mgr := binman.NewManager("/usr/local/bin",
    binman.WithSystemPaths([]string{"/usr/bin", "/usr/sbin"}),
)

// Check and resolve
mgr.IsInstalled(def)                          // true if found anywhere
mgr.IsPlatformSupported(def)                  // check platform constraints
path, err := mgr.ResolvePath(def)             // env > system paths > binDir
url := mgr.BuildURL(def, "v1.2.0")           // expand URL pattern

// Download and install
err := mgr.Download(def, "v1.2.0", func(downloaded, total int64) {
    fmt.Printf("\r%d/%d bytes", downloaded, total)
})

// Resolve or download if missing
path, err := mgr.EnsureInstalled(def, progressFn)

// Remove from managed directory
mgr.Remove(def)
```

#### Version Manifest

Track installed binary versions with a JSON manifest.

```go
// Create or load
manifest := binman.NewManifest()
manifest, err := binman.LoadManifest("/path/to/versions.json")  // empty on ENOENT

// Track versions
manifest.SetVersion("my-tool", "v1.2.0")
ver := manifest.GetVersion("my-tool")  // "v1.2.0"
manifest.Save("/path/to/versions.json")

// Compare versions (semver and date-based)
binman.CompareVersions("v1.2.0", "v1.3.0")  // -1
binman.IsNewer("v1.2.0", "v1.3.0")          // true
```

#### GitHub Releases and Self-Update

```go
// Check for updates
release, err := binman.GetLatestRelease("org/repo")
fmt.Println(release.TagName)

// Self-update check
latest, available, err := binman.CheckSelfUpdate("org/repo", "v1.0.0")

// Perform self-update (atomic binary replacement)
err := binman.SelfUpdate(binman.SelfUpdateConfig{
    Repo:       "org/repo",
    URLPattern: "https://github.com/org/repo/releases/download/{version}/app-{os}-{arch}",
    StatusFn:   func(msg string) { fmt.Println(msg) },
}, "v1.1.0")
```

## Supported Distributions

- Fedora, RHEL, CentOS, Rocky, Alma, Oracle Linux (dnf/yum)
- Debian, Ubuntu, Linux Mint, Pop!\_OS (apt)
- Arch, Manjaro, EndeavourOS (pacman)
- openSUSE, SLES (zypper)
- Alpine (apk)

Distributions with `ID_LIKE` matching the above are also supported.
