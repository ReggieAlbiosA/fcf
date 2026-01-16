package platform

// LinuxDistro holds information about the Linux distribution
type LinuxDistro struct {
	ID      string // ubuntu, debian, fedora, arch, alpine, rhel, centos, opensuse, etc.
	Name    string // Pretty name (e.g., "Ubuntu 22.04 LTS")
	Version string // Version ID (e.g., "22.04")
}

// DetectLinuxDistro parses /etc/os-release to identify the Linux distribution
// Platform-specific implementation in distro_unix.go and distro_windows.go
func DetectLinuxDistro() LinuxDistro {
	return detectLinuxDistro()
}

// GetFdInstallHint returns the package manager command to install fd
// Platform-specific implementation in distro_unix.go and distro_windows.go
func GetFdInstallHint() string {
	return getFdInstallHint()
}

// detectLinuxDistro is the platform-specific implementation
// Implemented in distro_unix.go and distro_windows.go

