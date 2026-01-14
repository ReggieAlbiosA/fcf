//go:build windows

package main

// LinuxDistro holds information about the Linux distribution (stub for Windows)
type LinuxDistro struct {
	ID      string
	Name    string
	Version string
}

// detectLinuxDistro returns empty struct on Windows (not applicable)
func detectLinuxDistro() LinuxDistro {
	return LinuxDistro{ID: "windows", Name: "Windows"}
}

// getFdInstallHint returns Windows-specific fd installation instructions
func getFdInstallHint() string {
	return "winget install sharkdp.fd\n  or: choco install fd\n  or: scoop install fd"
}
