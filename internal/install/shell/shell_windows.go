//go:build windows

package shell

import (
	"os"
	"path/filepath"
)

// getShellConfigPaths returns the config file paths for PowerShell on Windows
func getShellConfigPaths(homeDir string) map[ShellType]string {
	// Get PowerShell profile path from environment or use default
	profilePath := getPowerShellProfilePath()

	return map[ShellType]string{
		ShellPowerShell: profilePath,
	}
}

// getPowerShellProfilePath returns the path to the PowerShell profile
func getPowerShellProfilePath() string {
	// Try to use PROFILE environment variable set by PowerShell
	if profile := os.Getenv("PROFILE"); profile != "" {
		return profile
	}

	// Fallback to standard location for CurrentUser scope
	userProfile := os.Getenv("USERPROFILE")
	if userProfile != "" {
		// Try PowerShell 5.1+ location first
		ps7Profile := filepath.Join(userProfile, "Documents", "PowerShell", "profile.ps1")
		if fileExists(ps7Profile) {
			return ps7Profile
		}

		// Try Windows PowerShell location (v5.1 and earlier)
		ps5Profile := filepath.Join(userProfile, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
		return ps5Profile
	}

	return ""
}

// DetectShellsForInstallation returns PowerShell as the only shell option on Windows
func DetectShellsForInstallation(homeDir string) []ShellInfo {
	profilePath := getPowerShellProfilePath()
	if profilePath == "" {
		return []ShellInfo{}
	}

	return []ShellInfo{
		{
			Type:       ShellPowerShell,
			Name:       "PowerShell",
			ConfigPath: profilePath,
			Detected:   true,
		},
	}
}

// GetShellConfigPath returns the PowerShell profile path on Windows
func GetShellConfigPath(homeDir string, shellType ShellType) string {
	if shellType == ShellPowerShell {
		return getPowerShellProfilePath()
	}
	return ""
}

// GetShellReloadCommand returns the command to reload PowerShell profile
func GetShellReloadCommand(shellType ShellType) string {
	if shellType == ShellPowerShell {
		return ". $PROFILE"
	}
	return "restart PowerShell"
}
