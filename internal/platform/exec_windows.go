//go:build windows

package platform

import (
	"path/filepath"
	"strings"
)

// isExecutable checks if a file has an executable extension (Windows)
func isExecutable(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	execExts := []string{".exe", ".bat", ".cmd", ".ps1", ".com"}
	for _, e := range execExts {
		if ext == e {
			return true
		}
	}
	return false
}
