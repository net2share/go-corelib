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

Terminal UI utilities using [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss).

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

// Format text
bold := tui.Bold("important")
code := tui.Code("go build")
highlight := tui.Highlight("key value")

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
tui.WaitForEnter()  // "Press Enter to continue..."
tui.ClearLine()     // Clear current terminal line

// Progress spinner
spinner := tui.NewSpinner("Loading")
spinner.Start()
// ... do work ...
spinner.Stop()
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

## Supported Distributions

- Fedora, RHEL, CentOS, Rocky, Alma, Oracle Linux (dnf/yum)
- Debian, Ubuntu, Linux Mint, Pop!\_OS (apt)
- Arch, Manjaro, EndeavourOS (pacman)
- openSUSE, SLES (zypper)
- Alpine (apk)

Distributions with `ID_LIKE` matching the above are also supported.
