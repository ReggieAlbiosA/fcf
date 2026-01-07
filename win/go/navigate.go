package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// getNavFilePath returns the path to the navigation temp file
func getNavFilePath() string {
	tempDir := os.TempDir()
	return filepath.Join(tempDir, "fcf_nav_path")
}

// writeNavPath writes the navigation path to temp file
func writeNavPath(targetPath string) error {
	navFile := getNavFilePath()
	return os.WriteFile(navFile, []byte(targetPath), 0644)
}

// cleanupNavFile removes the navigation temp file if it exists
func cleanupNavFile() {
	navFile := getNavFilePath()
	os.Remove(navFile)
}

// navigateToPath handles navigation to a selected path
func navigateToPath(targetPath string) bool {
	// Get file info
	info, err := os.Stat(targetPath)
	if err != nil {
		fmt.Printf("%s Directory '%s' does not exist\n", colors.Red("ERROR:"), targetPath)
		return false
	}

	// If it's a file, get the parent directory
	if !info.IsDir() {
		targetPath = filepath.Dir(targetPath)
	}

	// Verify it's a valid directory
	info, err = os.Stat(targetPath)
	if err != nil || !info.IsDir() {
		fmt.Printf("%s '%s' is not a valid directory\n", colors.Red("ERROR:"), targetPath)
		return false
	}

	// Get absolute path
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		absPath = targetPath
	}

	// Write path to temp file for shell integration
	if err := writeNavPath(absPath); err != nil {
		fmt.Printf("%s Could not save navigation path: %v\n", colors.Red("ERROR:"), err)
		return false
	}

	fmt.Printf("%s %s\n", colors.Green("✓ Will navigate to:"), colors.Cyan(absPath))
	fmt.Println()

	// Show directory contents
	fmt.Println(colors.Dim("Contents:"))
	showDirectoryContents(absPath)

	return true
}

// showDirectoryContents displays the contents of a directory
func showDirectoryContents(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("  %s\n", colors.Red("Could not read directory"))
		return
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := entry.Name()
		if entry.IsDir() {
			fmt.Printf("  %s\n", colors.Blue(name+"/"))
		} else {
			size := formatSize(info.Size())
			fmt.Printf("  %s %s\n", name, colors.Dim("("+size+")"))
		}
	}
	fmt.Println()
}
