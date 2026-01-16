package platform

// IsExecutable checks if a file has execute permission
// Platform-specific implementation in exec_unix.go and exec_windows.go
func IsExecutable(path string) bool {
	return isExecutable(path)
}

// isExecutable is the platform-specific implementation
// Implemented in exec_unix.go and exec_windows.go

