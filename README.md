# go-corelib

Shared Go library providing common utilities for CLI applications.

## Packages

### tui

Terminal UI components for building interactive CLI applications.

```go
import "github.com/net2share/go-corelib/tui"

// Display functions
tui.PrintTitle("My Application")
tui.PrintSuccess("Operation completed")
tui.PrintError("Something went wrong")
tui.PrintWarning("Warning message")
tui.PrintInfo("Information")
tui.PrintStep(1, 5, "Processing...")

// Prompts
name := tui.Prompt("Enter name")
value := tui.PromptWithDefault("Enter value", "default")
num := tui.PromptInt("Enter number", 10, 1, 100)
choice := tui.PromptChoice("Select mode", []string{"ssh", "socks"}, "ssh")
confirmed := tui.Confirm("Continue?", true)

// Menu
options := []tui.MenuOption{
    {Key: "1", Label: "Option 1"},
    {Key: "2", Label: "Option 2"},
}
tui.ShowMenu(options)

// Box display
tui.PrintBox("Title", []string{"Line 1", "Line 2"})

// Terminal control
tui.ClearScreen()
tui.ClearLine()
tui.WaitForEnter()

// Progress bar
tui.PrintProgress(50, 100)
```

### osdetect

OS detection and package management utilities.

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

// Get system info
arch := osdetect.GetArch()        // "amd64", "arm64", etc.
iface, _ := osdetect.GetDefaultInterface()
port := osdetect.DetectSSHPort()  // "22"
```

## Supported Distributions

- Fedora, RHEL, CentOS, Rocky, Alma, Oracle Linux (dnf/yum)
- Debian, Ubuntu, Linux Mint, Pop!\_OS (apt)
- Arch, Manjaro, EndeavourOS (pacman)
- openSUSE, SLES (zypper)
- Alpine (apk)

Distributions with `ID_LIKE` matching the above are also supported.
