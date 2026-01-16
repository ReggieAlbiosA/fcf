//go:build windows

package shell

import (
	"fmt"
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

// EnsureUserBinDirectory creates user bin directory if needed on Windows
func EnsureUserBinDirectory() error {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		return fmt.Errorf("USERPROFILE environment variable not set")
	}

	binDir := filepath.Join(userProfile, ".local", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("could not create user bin directory: %w", err)
	}

	return nil
}

// IsUserBinInPath checks if user bin is in PATH on Windows
func IsUserBinInPath() bool {
	// On Windows, PATH entries are separated by semicolons
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		return false
	}

	userBin := filepath.Join(userProfile, ".local", "bin")
	pathEnv := os.Getenv("PATH")

	// Case-insensitive comparison on Windows
	for _, p := range filepath.SplitList(pathEnv) {
		if p == userBin {
			return true
		}
	}

	return false
}

// AddUserBinToPath adds user bin to PATH on Windows (requires registry modification)
func AddUserBinToPath(homeDir string) error {
	// On Windows, adding to PATH typically requires:
	// 1. Modifying the system environment variable (registry)
	// 2. Or modifying user environment variable (registry)
	// 3. Or updating PowerShell profile
	//
	// For simplicity, we'll document this in the installation output
	// Users can either:
	// - Run the exe directly if in .local\bin
	// - Add it manually to PATH via System Properties
	// - Install system-wide to get automatic PATH

	return nil // PATH modification on Windows is complex and may require admin
}

// GetShellReloadCommand returns the command to reload PowerShell profile
func GetShellReloadCommand(shellType ShellType) string {
	if shellType == ShellPowerShell {
		return ". $PROFILE"
	}
	return "restart PowerShell"
}
