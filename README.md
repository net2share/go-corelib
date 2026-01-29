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

#### Full-Screen Interactive Components

Interactive menus, confirmations, and inputs using bubbletea.

```go
// Selection menu
selected, err := tui.RunMenu(tui.MenuConfig{
    Title:       "Choose an option",
    Description: "Use arrow keys to navigate",
    Options: []tui.MenuOption{
        {Label: "Option 1", Value: "opt1"},
        {Label: "Option 2", Value: "opt2"},
    },
})

// Simple menu with just labels
selected, err := tui.RunMenuSimple("Select Mode", []string{"Development", "Production"})

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

// Full-screen progress indicator
progress := tui.StartProgress(tui.ProgressConfig{
    Title:   "Installing",
    Message: "Please wait...",
})
// ... do work ...
progress.Done()
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
