// Package osdetect provides OS detection and package management utilities.
package osdetect

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ErrNotRoot is returned when root privileges are required but not present.
var ErrNotRoot = errors.New("this program must be run as root")

// OSInfo contains detected OS information.
type OSInfo struct {
	ID             string // e.g., "fedora", "ubuntu", "debian"
	IDLike         string // e.g., "rhel fedora", "debian"
	PrettyName     string // e.g., "Fedora Linux 39"
	VersionID      string // e.g., "39", "22.04"
	PackageManager string // "dnf", "apt", "yum", "pacman", "zypper", "apk"
	InstallCmd     string // Full install command, e.g., "dnf install -y"
}

// Detect reads /etc/os-release and determines package manager.
func Detect() (*OSInfo, error) {
	info := &OSInfo{}

	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			info.ID = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
		if strings.HasPrefix(line, "ID_LIKE=") {
			info.IDLike = strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), "\"")
		}
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			info.PrettyName = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			info.VersionID = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
		}
	}

	info.PackageManager, info.InstallCmd = detectPackageManager(info.ID, info.IDLike)

	return info, nil
}

// detectPackageManager determines the package manager based on OS ID.
func detectPackageManager(id, idLike string) (manager, installCmd string) {
	switch id {
	case "fedora", "rhel", "centos", "rocky", "alma", "ol":
		// Check if dnf is available, fallback to yum
		if _, err := exec.LookPath("dnf"); err == nil {
			return "dnf", "dnf install -y"
		}
		if _, err := exec.LookPath("yum"); err == nil {
			return "yum", "yum install -y"
		}
		return "dnf", "dnf install -y"

	case "debian", "ubuntu", "linuxmint", "pop":
		return "apt", "apt-get install -y"

	case "arch", "manjaro", "endeavouros":
		return "pacman", "pacman -S --noconfirm"

	case "opensuse", "opensuse-leap", "opensuse-tumbleweed", "sles":
		return "zypper", "zypper install -y"

	case "alpine":
		return "apk", "apk add"
	}

	// Try to detect from ID_LIKE
	if strings.Contains(idLike, "debian") || strings.Contains(idLike, "ubuntu") {
		return "apt", "apt-get install -y"
	}
	if strings.Contains(idLike, "fedora") || strings.Contains(idLike, "rhel") {
		if _, err := exec.LookPath("dnf"); err == nil {
			return "dnf", "dnf install -y"
		}
		return "yum", "yum install -y"
	}
	if strings.Contains(idLike, "arch") {
		return "pacman", "pacman -S --noconfirm"
	}

	// Fallback: try to detect by available commands
	managers := []struct {
		cmd        string
		name       string
		installCmd string
	}{
		{"apt", "apt", "apt-get install -y"},
		{"dnf", "dnf", "dnf install -y"},
		{"yum", "yum", "yum install -y"},
		{"pacman", "pacman", "pacman -S --noconfirm"},
		{"zypper", "zypper", "zypper install -y"},
		{"apk", "apk", "apk add"},
	}

	for _, m := range managers {
		if _, err := exec.LookPath(m.cmd); err == nil {
			return m.name, m.installCmd
		}
	}

	return "unknown", ""
}

// InstallPackage installs a package using the detected package manager.
func (o *OSInfo) InstallPackage(pkg string) error {
	if o.InstallCmd == "" {
		return fmt.Errorf("could not detect package manager for OS '%s'", o.ID)
	}

	// Update package cache for apt
	if o.PackageManager == "apt" {
		cmd := exec.Command("apt-get", "update", "-qq")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run() // Ignore errors from update
	}

	// Split the install command and append the package
	parts := strings.Fields(o.InstallCmd)
	parts = append(parts, pkg)

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsRoot checks if running as root (uid == 0).
func IsRoot() bool {
	return os.Geteuid() == 0
}

// RequireRoot returns ErrNotRoot if not running as root.
// Use this for consistent root check error messages across projects.
func RequireRoot() error {
	if !IsRoot() {
		return ErrNotRoot
	}
	return nil
}

// HasSystemd checks if systemctl is available.
func HasSystemd() bool {
	_, err := exec.LookPath("systemctl")
	return err == nil
}

// GetArch returns the system architecture in common naming (amd64, arm64, armv7, 386).
func GetArch() string {
	return runtime.GOARCH
}

// HasIPv6 checks if IPv6 is available on any interface.
func HasIPv6() bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() == nil && ipnet.IP.To16() != nil {
					return true
				}
			}
		}
	}

	return false
}

// GetDefaultInterface returns the default network interface name.
func GetDefaultInterface() (string, error) {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip header

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[1] == "00000000" {
			return fields[0], nil
		}
	}

	return "", scanner.Err()
}

// DetectSSHPort reads the SSH port from sshd_config.
func DetectSSHPort() string {
	file, err := os.Open("/etc/ssh/sshd_config")
	if err != nil {
		return "22"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.ToLower(fields[0]) == "port" {
			return fields[1]
		}
	}

	return "22"
}
