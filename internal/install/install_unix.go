//go:build unix

package install

import (
	"os"
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
