//go:build unix

package platform

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)


// detectLinuxDistro parses /etc/os-release to identify the Linux distribution
func detectLinuxDistro() LinuxDistro {
	// Only applicable on Linux
	if runtime.GOOS != "linux" {
		return LinuxDistro{ID: "macos", Name: "macOS"}
	}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		return LinuxDistro{ID: "unknown", Name: "Unknown Linux"}
	}
	defer file.Close()

	var distro LinuxDistro
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "ID=") {
			distro.ID = strings.Trim(strings.TrimPrefix(line, "ID="), `"`)
		} else if strings.HasPrefix(line, "NAME=") {
			distro.Name = strings.Trim(strings.TrimPrefix(line, "NAME="), `"`)
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			distro.Version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
		}
	}

	if distro.ID == "" {
		distro.ID = "unknown"
	}
	if distro.Name == "" {
		distro.Name = "Linux"
	}

	return distro
}

// getFdInstallHint returns the package manager command to install fd
func getFdInstallHint() string {
	if runtime.GOOS == "darwin" {
		return "brew install fd"
	}

	distro := detectLinuxDistro()

	switch distro.ID {
	case "ubuntu", "debian", "linuxmint", "pop":
		return "sudo apt install fd-find"
	case "fedora":
		return "sudo dnf install fd-find"
	case "rhel", "centos", "rocky", "almalinux":
		return "sudo yum install fd-find"
	case "arch", "manjaro", "endeavouros":
		return "sudo pacman -S fd"
	case "opensuse", "opensuse-leap", "opensuse-tumbleweed":
		return "sudo zypper install fd"
	case "alpine":
		return "sudo apk add fd"
	case "void":
		return "sudo xbps-install fd"
	case "gentoo":
		return "sudo emerge sys-apps/fd"
	case "nixos":
		return "nix-env -iA nixpkgs.fd"
	default:
		return "Install 'fd' using your package manager"
	}
}
