//go:build unix

package install

import (
	"os"
	"os/user"
	"path/filepath"
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

// getRealUserHomeDir returns the home directory of the actual user,
// even when running under sudo. This is essential for shell integration
// to be written to the correct user's config files.
func getRealUserHomeDir() (string, error) {
	// Check if running under sudo
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		// Look up the actual user's home directory
		u, err := user.Lookup(sudoUser)
		if err != nil {
			// Fall back to os.UserHomeDir if lookup fails
			return os.UserHomeDir()
		}
		return u.HomeDir, nil
	}

	// Not running under sudo, use standard method
	return os.UserHomeDir()
}
