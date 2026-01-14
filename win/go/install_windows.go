//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	modadvapi32          = syscall.NewLazyDLL("advapi32.dll")
	procOpenProcessToken = modadvapi32.NewProc("OpenProcessToken")
	procGetTokenInfo     = modadvapi32.NewProc("GetTokenInformation")
)

const (
	tokenQuery          = 0x0008
	tokenInfoElevation  = 20
)

// isElevated checks if the process is running with Administrator privileges
func isElevated() bool {
	var token syscall.Token
	currentProcess, _ := syscall.GetCurrentProcess()

	err := syscall.OpenProcessToken(currentProcess, tokenQuery, &token)
	if err != nil {
		return false
	}
	defer token.Close()

	var elevation uint32
	var size uint32

	err = syscall.GetTokenInformation(token, tokenInfoElevation, (*byte)(unsafe.Pointer(&elevation)), uint32(unsafe.Sizeof(elevation)), &size)
	if err != nil {
		return false
	}

	return elevation != 0
}

// getInstallPath returns the system-wide installation path
func getInstallPath() string {
	return `C:\Program Files\fcf\fcf.exe`
}

// ensureInstallDir ensures the installation directory exists
func ensureInstallDir() error {
	dir := filepath.Dir(getInstallPath())
	return os.MkdirAll(dir, 0755)
}

// makeExecutable is a no-op on Windows (executability is determined by extension)
func makeExecutable(path string) error {
	return nil
}

// postInstall adds fcf to the system PATH if not already present
func postInstall() error {
	installDir := filepath.Dir(getInstallPath())

	// Check if already in PATH
	currentPath := os.Getenv("PATH")
	if strings.Contains(strings.ToLower(currentPath), strings.ToLower(installDir)) {
		return nil
	}

	// Add to system PATH using PowerShell
	script := fmt.Sprintf(`
		$installDir = '%s'
		$currentPath = [Environment]::GetEnvironmentVariable('Path', 'Machine')
		if ($currentPath -notlike "*$installDir*") {
			$newPath = $currentPath + ';' + $installDir
			[Environment]::SetEnvironmentVariable('Path', $newPath, 'Machine')
		}
	`, installDir)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not add to PATH: %w", err)
	}

	fmt.Printf("%s Added %s to system PATH\n", colors.Green("OK:"), colors.Cyan(installDir))
	fmt.Println(colors.Yellow("Note: You may need to restart your terminal for PATH changes to take effect."))

	return nil
}
