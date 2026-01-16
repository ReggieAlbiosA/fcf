//go:build windows

package platform

// detectLinuxDistro returns empty struct on Windows (not applicable)
func detectLinuxDistro() LinuxDistro {
	return LinuxDistro{ID: "windows", Name: "Windows"}
}

// getFdInstallHint returns Windows-specific fd installation instructions
func getFdInstallHint() string {
	return "winget install sharkdp.fd\n  or: choco install fd\n  or: scoop install fd"
}
