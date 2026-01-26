package navigation

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/ReggieAlbiosA/fcf/internal/ui"
)

// sudoUserInfo holds information about the original user when running under sudo
var sudoUserInfo struct {
	detected bool
	username string
	uid      int
}

// getNavFilePath returns the path to the navigation temp file
// On Unix systems, includes UID to prevent multi-user permission conflicts
// On Windows, %TEMP% is already user-specific so no UID needed
// When running under sudo, uses SUDO_USER's UID so navigation works with user's shell
func getNavFilePath() string {
	tempDir := os.TempDir()

	if runtime.GOOS == "windows" {
		// Windows: %TEMP% is already user-specific
		return filepath.Join(tempDir, "fcf_nav_path")
	}

	// Unix (Linux/macOS): Check if running under sudo
	uid := os.Getuid()

	// If running as root, check for SUDO_USER
	if uid == 0 {
		if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
			// Check SUDO_COMMAND to distinguish "sudo fcf" from "sudo su -> fcf"
			sudoCmd := os.Getenv("SUDO_COMMAND")
			isSudoSu := false
			if sudoCmd != "" {
				// If SUDO_COMMAND contains su or a shell, we're in "sudo su" mode
				base := filepath.Base(sudoCmd)
				for _, sh := range []string{"su", "bash", "zsh", "fish", "sh"} {
					if len(base) >= len(sh) && base[:len(sh)] == sh {
						isSudoSu = true
						break
					}
				}
			}

			// Only use SUDO_USER's UID if this is direct "sudo fcf", not "sudo su -> fcf"
			if !isSudoSu {
				if u, err := user.Lookup(sudoUser); err == nil {
					if sudoUID, err := strconv.Atoi(u.Uid); err == nil {
						// Store for later message display
						sudoUserInfo.detected = true
						sudoUserInfo.username = sudoUser
						sudoUserInfo.uid = sudoUID
						uid = sudoUID
					}
				}
			}
		}
	}

	return filepath.Join(tempDir, fmt.Sprintf("fcf_nav_path_%d", uid))
}

// ShowSudoNavigationNote displays a note if running under sudo
// Call this before navigation to inform the user
func ShowSudoNavigationNote() {
	if sudoUserInfo.detected {
		fmt.Printf("%s Running under sudo - navigation will apply to %s's shell\n",
			ui.Colors.Yellow("Note:"),
			ui.Colors.Cyan(sudoUserInfo.username))
		fmt.Println()
	}
}

// writeNavPath writes the navigation path to temp file
func writeNavPath(targetPath string) error {
	navFile := getNavFilePath()
	return os.WriteFile(navFile, []byte(targetPath), 0644)
}

// cleanupNavFile removes the navigation temp file if it exists
func CleanupNavFile() {
	navFile := getNavFilePath()
	os.Remove(navFile)
}

// navigateToPath handles navigation to a selected path
func NavigateToPath(targetPath string) bool {
	// Get file info
	info, err := os.Stat(targetPath)
	if err != nil {
		fmt.Printf("%s Directory '%s' does not exist\n", ui.Colors.Red("ERROR:"), targetPath)
		return false
	}

	// If it's a file, get the parent directory
	if !info.IsDir() {
		targetPath = filepath.Dir(targetPath)
	}

	// Verify it's a valid directory
	info, err = os.Stat(targetPath)
	if err != nil || !info.IsDir() {
		fmt.Printf("%s '%s' is not a valid directory\n", ui.Colors.Red("ERROR:"), targetPath)
		return false
	}

	// Get absolute path
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		absPath = targetPath
	}

	// Write path to temp file for shell integration
	if err := writeNavPath(absPath); err != nil {
		fmt.Printf("%s Could not save navigation path: %v\n", ui.Colors.Red("ERROR:"), err)
		return false
	}

	fmt.Printf("%s %s\n", ui.Colors.Green("✓ Will navigate to:"), ui.Colors.Cyan(absPath))
	fmt.Println()

	// Show directory contents
	fmt.Println(ui.Colors.Dim("Contents:"))
	showDirectoryContents(absPath)

	return true
}

// showDirectoryContents displays the contents of a directory
func showDirectoryContents(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("  %s\n", ui.Colors.Red("Could not read directory"))
		return
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := entry.Name()
		if entry.IsDir() {
			fmt.Printf("  %s\n", ui.Colors.Blue(name+"/"))
		} else {
			size := ui.FormatSize(info.Size())
			fmt.Printf("  %s %s\n", name, ui.Colors.Dim("("+size+")"))
		}
	}
	fmt.Println()
}
