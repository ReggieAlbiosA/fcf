//go:build unix

package shell

import (
	"path/filepath"
)

// getShellConfigPaths returns the config file paths for each shell type on Unix
func getShellConfigPaths(homeDir string) map[ShellType]string {
	return map[ShellType]string{
		ShellBash: filepath.Join(homeDir, ".bashrc"),
		ShellZsh:  filepath.Join(homeDir, ".zshrc"),
		ShellFish: filepath.Join(homeDir, ".config", "fish", "config.fish"),
	}
}

// DetectShellsForInstallation detects which shells need configuration
func DetectShellsForInstallation(homeDir string) []ShellInfo {
	var shells []ShellInfo

	// First try to detect from $SHELL env var (most reliable for default shell)
	if shellType, _ := detectShellFromEnv(); shellType != ShellUnknown {
		configPath := GetShellConfigPath(homeDir, shellType)
		shells = append(shells, ShellInfo{
			Type:       shellType,
			Name:       shellType.String(),
			ConfigPath: configPath,
			Detected:   true,
		})
	} else {
		// Fallback: check for config files
		shells = DetectShellsFromConfigFiles(homeDir)
	}

	return shells
}

// GetShellConfigPath returns the primary config file path for a shell type
func GetShellConfigPath(homeDir string, shellType ShellType) string {
	switch shellType {
	case ShellBash:
		// On macOS, prefer ~/.bash_profile for login shells
		// On Linux, prefer ~/.bashrc
		bashProfile := filepath.Join(homeDir, ".bash_profile")
		bashrc := filepath.Join(homeDir, ".bashrc")

		// If .bash_profile exists, use it (typical macOS convention)
		if fileExists(bashProfile) {
			return bashProfile
		}
		// Otherwise use .bashrc (typical Linux convention)
		return bashrc

	case ShellZsh:
		return filepath.Join(homeDir, ".zshrc")
	case ShellFish:
		return filepath.Join(homeDir, ".config", "fish", "config.fish")
	default:
		return ""
	}
}

// GetShellReloadCommand returns the command to reload shell config
func GetShellReloadCommand(shellType ShellType) string {
	switch shellType {
	case ShellBash:
		return "source ~/.bashrc"
	case ShellZsh:
		return "source ~/.zshrc"
	case ShellFish:
		return "source ~/.config/fish/config.fish"
	default:
		return "restart your terminal"
	}
}
