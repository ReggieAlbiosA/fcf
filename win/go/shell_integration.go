package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ShellType represents detected shell types
type ShellType int

const (
	ShellUnknown ShellType = iota
	ShellBash
	ShellZsh
	ShellFish
	ShellPowerShell
)

// ShellInfo contains detected shell information
type ShellInfo struct {
	Type       ShellType
	Name       string
	ConfigPath string
	Detected   bool
}

const (
	fcfMarkerStart = "# FCF shell integration - DO NOT EDIT (managed by fcf install)"
	fcfMarkerEnd   = "# END FCF shell integration"
)

// getShellName returns the display name for a shell type
func (s ShellType) String() string {
	switch s {
	case ShellBash:
		return "Bash"
	case ShellZsh:
		return "Zsh"
	case ShellFish:
		return "Fish"
	case ShellPowerShell:
		return "PowerShell"
	default:
		return "Unknown"
	}
}

// detectShellFromEnv detects the shell from $SHELL environment variable
func detectShellFromEnv() (ShellType, string) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return ShellUnknown, ""
	}

	shellName := filepath.Base(shell)
	switch shellName {
	case "bash":
		return ShellBash, shell
	case "zsh":
		return ShellZsh, shell
	case "fish":
		return ShellFish, shell
	default:
		return ShellUnknown, shell
	}
}

// detectShellsFromConfigFiles detects available shells by checking for config files
func detectShellsFromConfigFiles(homeDir string) []ShellInfo {
	var shells []ShellInfo

	// Check Bash config files (both .bashrc and .bash_profile)
	bashConfigs := []string{".bashrc", ".bash_profile"}
	for _, config := range bashConfigs {
		path := filepath.Join(homeDir, config)
		if fileExists(path) {
			shells = append(shells, ShellInfo{
				Type:       ShellBash,
				Name:       "Bash",
				ConfigPath: path,
				Detected:   true,
			})
			break // Only add once, prefer .bashrc
		}
	}

	// Check Zsh config
	zshPath := filepath.Join(homeDir, ".zshrc")
	if fileExists(zshPath) {
		shells = append(shells, ShellInfo{
			Type:       ShellZsh,
			Name:       "Zsh",
			ConfigPath: zshPath,
			Detected:   true,
		})
	}

	// Check Fish config
	fishPath := filepath.Join(homeDir, ".config", "fish", "config.fish")
	if fileExists(fishPath) {
		shells = append(shells, ShellInfo{
			Type:       ShellFish,
			Name:       "Fish",
			ConfigPath: fishPath,
			Detected:   true,
		})
	}

	return shells
}

// getShellFunction returns the shell wrapper function for the given shell type
func getShellFunction(shellType ShellType) string {
	switch shellType {
	case ShellBash, ShellZsh:
		return getBashZshFunction()
	case ShellFish:
		return getFishFunction()
	case ShellPowerShell:
		return getPowerShellFunction()
	default:
		return ""
	}
}

// getBashZshFunction returns the Bash/Zsh wrapper function
func getBashZshFunction() string {
	return fcfMarkerStart + `
fcf() {
    local nav_file="/tmp/fcf_nav_path"
    rm -f "$nav_file"
    command fcf "$@"
    if [[ -f "$nav_file" ]]; then
        local target
        target=$(cat "$nav_file")
        rm -f "$nav_file"
        if [[ -d "$target" ]]; then
            cd "$target" || return
        fi
    fi
}
` + fcfMarkerEnd
}

// getFishFunction returns the Fish shell wrapper function
func getFishFunction() string {
	return fcfMarkerStart + `
function fcf
    set nav_file /tmp/fcf_nav_path
    rm -f $nav_file
    command fcf $argv
    if test -f $nav_file
        set target (cat $nav_file)
        rm -f $nav_file
        if test -d $target
            cd $target
        end
    end
end
` + fcfMarkerEnd
}

// getPowerShellFunction returns the PowerShell wrapper function
func getPowerShellFunction() string {
	return `# FCF shell integration - DO NOT EDIT (managed by fcf install)
function fcf {
    $navFile = Join-Path $env:TEMP "fcf_nav_path"
    if (Test-Path $navFile) { Remove-Item $navFile -Force }
    & "C:\Program Files\fcf\fcf.exe" @args
    if (Test-Path $navFile) {
        $target = Get-Content $navFile -Raw
        Remove-Item $navFile -Force
        if (Test-Path $target -PathType Container) {
            Set-Location $target
        }
    }
}
# END FCF shell integration`
}

// hasExistingInstallation checks if shell integration is already installed
func hasExistingInstallation(configPath string) bool {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), fcfMarkerStart)
}

// addShellIntegration adds the shell wrapper function to the config file
func addShellIntegration(configPath string, shellType ShellType) error {
	// Create config directory if needed (especially for Fish)
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	// Read existing content
	content, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not read config file: %w", err)
	}

	contentStr := string(content)

	// Check for existing installation
	if strings.Contains(contentStr, fcfMarkerStart) {
		// Update existing installation
		return updateExistingInstallation(configPath, contentStr, shellType)
	}

	// Append new installation
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open config file: %w", err)
	}
	defer f.Close()

	// Add newlines for separation if file is not empty
	if len(content) > 0 {
		if !strings.HasSuffix(contentStr, "\n\n") {
			if _, err := f.WriteString("\n\n"); err != nil {
				return fmt.Errorf("could not write to config file: %w", err)
			}
		}
	}

	shellFunc := getShellFunction(shellType)
	if _, err := f.WriteString(shellFunc + "\n"); err != nil {
		return fmt.Errorf("could not write shell function: %w", err)
	}

	return nil
}

// updateExistingInstallation replaces existing shell integration with new one
func updateExistingInstallation(configPath string, content string, shellType ShellType) error {
	startIdx := strings.Index(content, fcfMarkerStart)
	endIdx := strings.Index(content, fcfMarkerEnd)

	if startIdx == -1 || endIdx == -1 {
		// Markers not properly formatted, skip update
		return nil
	}

	// Find the end of the END marker
	endMarkerEnd := endIdx + len(fcfMarkerEnd)

	// Build new content
	newFunc := getShellFunction(shellType)
	newContent := content[:startIdx] + newFunc + content[endMarkerEnd:]

	// Write back to file
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
