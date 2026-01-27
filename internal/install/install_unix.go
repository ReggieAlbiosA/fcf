//go:build unix

package install

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// isElevated checks if the process is running as root
func isElevated() bool {
	return os.Geteuid() == 0
}

// getInstallPath returns the system-wide installation path
func getInstallPath() string {
	return "/usr/local/bin/fcf"
}

// ensureInstallDir ensures the installation directory exists
func ensureInstallDir() error {
	dir := filepath.Dir(getInstallPath())
	return os.MkdirAll(dir, 0755)
}

// makeExecutable sets the executable permission on the installed binary
func makeExecutable(path string) error {
	return os.Chmod(path, 0755)
}

// postInstall performs any post-installation tasks (Unix: none needed)
func postInstall() error {
	// On Unix, /usr/local/bin is typically already in PATH
	// No additional setup required
	return nil
}

// postUninstall performs any post-uninstallation cleanup (Unix: none needed)
func postUninstall() {
	// On Unix, no additional cleanup required
	// /usr/local/bin remains in PATH (used by other tools)
}

// isInSudoSuMode detects if we're running in "sudo su" mode vs direct "sudo command"
// In "sudo su" mode, we want to use root's home, not SUDO_USER's home
func isInSudoSuMode() bool {
	sudoCmd := os.Getenv("SUDO_COMMAND")
	if sudoCmd == "" {
		return false
	}

	// Get the base command name
	base := filepath.Base(sudoCmd)

	// Check if SUDO_COMMAND is a shell or su command
	// This indicates "sudo su" or "sudo bash" etc., not direct "sudo fcf"
	shellCommands := []string{"su", "bash", "zsh", "fish", "sh", "-bash", "-zsh", "-fish", "-sh"}
	for _, sh := range shellCommands {
		if base == sh || strings.HasPrefix(base, sh+" ") {
			return true
		}
	}

	return false
}

// getRootHomeDirIfNeeded returns root's home directory when running
// with sudo (not sudo su) and the real user isn't root.
// Returns "" if root's shell integration doesn't need updating.
func getRootHomeDirIfNeeded(realUserHome string) string {
	if !isElevated() || isInSudoSuMode() {
		return ""
	}
	u, err := user.Lookup("root")
	if err != nil {
		return ""
	}
	if u.HomeDir == realUserHome {
		return "" // already targeting root's home
	}
	return u.HomeDir
}

// getRealUserHomeDir returns the home directory of the actual user,
// even when running under sudo. This is essential for shell integration
// to be written to the correct user's config files.
//
// When running "sudo fcf install": returns original user's home (e.g., /home/reggie)
// When running "sudo su" then "fcf install": returns root's home (/root)
func getRealUserHomeDir() (string, error) {
	// Check if running under sudo
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" && !isInSudoSuMode() {
		// Direct "sudo fcf install" - use the original user's home
		u, err := user.Lookup(sudoUser)
		if err != nil {
			// Fall back to os.UserHomeDir if lookup fails
			return os.UserHomeDir()
		}
		return u.HomeDir, nil
	}

	// Not running under sudo, or in "sudo su" mode - use standard method
	// In "sudo su" mode, os.UserHomeDir() correctly returns /root
	return os.UserHomeDir()
}
