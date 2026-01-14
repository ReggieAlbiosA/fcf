//go:build unix

package main

import "os"

// isExecutable checks if a file has execute permission (Unix)
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	// Check if any execute bit is set (owner, group, or other)
	return info.Mode().Perm()&0111 != 0
}
