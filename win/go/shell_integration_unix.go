//go:build unix

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getShellConfigPaths returns the config file paths for each shell type on Unix
func getShellConfigPaths(homeDir string) map[ShellType]string {
	return map[ShellType]string{
		ShellBash: filepath.Join(homeDir, ".bashrc"),
		ShellZsh:  filepath.Join(homeDir, ".zshrc"),
		ShellFish: filepath.Join(homeDir, ".config", "fish", "config.fish"),
	}
}

// detectShellsForInstallation detects which shells need configuration
func detectShellsForInstallation(homeDir string) []ShellInfo {
	var shells []ShellInfo

	// First try to detect from $SHELL env var (most reliable for default shell)
	if shellType, _ := detectShellFromEnv(); shellType != ShellUnknown {
		configPath := getShellConfigPath(homeDir, shellType)
		shells = append(shells, ShellInfo{
			Type:       shellType,
			Name:       shellType.String(),
			ConfigPath: configPath,
			Detected:   true,
		})
	} else {
		// Fallback: check for config files
		shells = detectShellsFromConfigFiles(homeDir)
	}

	return shells
}

// getShellConfigPath returns the primary config file path for a shell type
func getShellConfigPath(homeDir string, shellType ShellType) string {
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

// ensureUserBinDirectory creates ~/.local/bin if it doesn't exist
func ensureUserBinDirectory() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}

	binDir := filepath.Join(homeDir, ".local", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("could not create user bin directory: %w", err)
	}

	return nil
}

// isUserBinInPath checks if ~/.local/bin is in the user's PATH
func isUserBinInPath() bool {
	pathEnv := os.Getenv("PATH")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	userBin := filepath.Join(homeDir, ".local", "bin")
	for _, p := range strings.Split(pathEnv, ":") {
		if p == userBin {
			return true
		}
	}
	return false
}

// addUserBinToPath adds ~/.local/bin to PATH in shell config files
func addUserBinToPath(homeDir string) error {
	// Get the default shell
	shellType, _ := detectShellFromEnv()
	if shellType == ShellUnknown {
		// Try to detect from config files
		shells := detectShellsFromConfigFiles(homeDir)
		if len(shells) > 0 {
			shellType = shells[0].Type
		}
	}

	if shellType == ShellUnknown {
		return nil // Can't add PATH if shell unknown
	}

	configPath := getShellConfigPath(homeDir, shellType)
	if configPath == "" {
		return nil
	}

	// Check if PATH is already set
	content, err := os.ReadFile(configPath)
	if err == nil && strings.Contains(string(content), ".local/bin") {
		return nil // Already in PATH
	}

	// Append PATH configuration
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open config file: %w", err)
	}
	defer f.Close()

	pathExport := ""
	switch shellType {
	case ShellBash, ShellZsh:
		pathExport = `export PATH="$HOME/.local/bin:$PATH"`
	case ShellFish:
		pathExport = `set -gx PATH $HOME/.local/bin $PATH`
	default:
		return nil
	}

	// Add newline before PATH export for readability
	if _, err := f.WriteString("\n\n# Add ~/.local/bin to PATH\n"); err != nil {
		return fmt.Errorf("could not write to config file: %w", err)
	}

	if _, err := f.WriteString(pathExport + "\n"); err != nil {
		return fmt.Errorf("could not write PATH export: %w", err)
	}

	return nil
}

// getShellReloadCommand returns the command to reload shell config
func getShellReloadCommand(shellType ShellType) string {
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
