package search

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ReggieAlbiosA/fcf/internal/input"
	"github.com/ReggieAlbiosA/fcf/internal/ui"
)

// SearchResult contains the search results and metadata
type SearchResult struct {
	Results []string
	Stopped bool // true if search was stopped by user
}

// getFdCommand returns the fd command name if available
// Checks for "fd" first (standard), then "fdfind" (Debian/Ubuntu package name)
func getFdCommand() string {
	if _, err := exec.LookPath("fd"); err == nil {
		return "fd"
	}
	if _, err := exec.LookPath("fdfind"); err == nil {
		return "fdfind"
	}
	return ""
}

// hasFd checks if fd is available in PATH
func HasFd() bool {
	return getFdCommand() != ""
}

// SearchWithFd uses fd for fast parallel search
func SearchWithFd(pattern, searchPath string, opts *ui.Options, stopChan <-chan struct{}) (*SearchResult, error) {
	args := []string{"--color", "never", "--hidden", "--no-ignore"}

	// Type filter
	if opts.Type == "f" {
		args = append(args, "-t", "f")
	} else if opts.Type == "d" {
		args = append(args, "-t", "d")
	}

	// Case sensitivity
	if opts.IgnoreCase {
		args = append(args, "-i")
	} else {
		args = append(args, "-s")
	}

	// Glob pattern and path
	args = append(args, "-g", pattern, searchPath)

	cmd := exec.Command(getFdCommand(), args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	result := &SearchResult{
		Results: []string{},
		Stopped: false,
	}
	scanner := bufio.NewScanner(stdout)
	count := 0

	for scanner.Scan() {
		// Check for stop signal
		select {
		case <-stopChan:
			cmd.Process.Kill()
			result.Stopped = true
			return result, nil
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		count++
		result.Results = append(result.Results, line)

		// Display result in real-time (streaming)
		if opts.MaxDisplay == 0 || count <= opts.MaxDisplay {
			ui.ShowResult(line, count)
		}
	}

	cmd.Wait()
	return result, nil
}

// SearchWithWalk uses filepath.WalkDir as fallback
func SearchWithWalk(pattern, searchPath string, opts *ui.Options, stopChan <-chan struct{}) (*SearchResult, error) {
	result := &SearchResult{
		Results: []string{},
		Stopped: false,
	}
	count := 0

	err := filepath.WalkDir(searchPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}

		// Check for stop signal
		select {
		case <-stopChan:
			result.Stopped = true
			return filepath.SkipAll
		default:
		}

		// Skip the root directory itself
		if path == searchPath {
			return nil
		}

		// Type filter
		if opts.Type == "f" && d.IsDir() {
			return nil
		}
		if opts.Type == "d" && !d.IsDir() {
			return nil
		}

		// Pattern matching
		name := d.Name()
		matched := matchPattern(name, pattern, opts.IgnoreCase)
		if !matched {
			return nil
		}

		count++
		result.Results = append(result.Results, path)

		// Display result in real-time (streaming)
		if opts.MaxDisplay == 0 || count <= opts.MaxDisplay {
			ui.ShowResult(path, count)
		}

		return nil
	})

	return result, err
}

// matchPattern checks if name matches the glob pattern
func matchPattern(name, pattern string, ignoreCase bool) bool {
	if ignoreCase {
		name = strings.ToLower(name)
		pattern = strings.ToLower(pattern)
	}

	matched, err := filepath.Match(pattern, name)
	if err != nil {
		return false
	}
	return matched
}

// search performs the search using fd or fallback
func Search(pattern, searchPath string) ([]string, bool) {
	result, _ := SearchWithStop(pattern, searchPath)
	return result.Results, HasFd()
}

// SearchWithStop performs the search with ability to stop via 's' key
func SearchWithStop(pattern, searchPath string) (*SearchResult, bool) {
	// Resolve search path
	absPath, err := filepath.Abs(searchPath)
	if err != nil {
		absPath = searchPath
	}

	usingFd := HasFd()
	ui.ShowSearchInfo(absPath, pattern, usingFd)

	fmt.Printf("%s %s  %s\n\n",
		ui.Colors.Bold("Results:"),
		ui.Colors.Dim("(streaming in real-time...)"),
		ui.Colors.Yellow("[press 's' to stop]"))

	// Set up stop channel and key listener
	stopChan := make(chan struct{})
	keyChan := make(chan string, 10)
	stopListener := input.StartKeyListener(keyChan)
	defer stopListener()

	// Goroutine to handle 's' key press
	go func() {
		for key := range keyChan {
			if strings.ToLower(key) == "s" {
				close(stopChan)
				return
			}
		}
	}()

	var result *SearchResult
	if usingFd {
		result, _ = SearchWithFd(pattern, absPath, &ui.Opts, stopChan)
	} else {
		result, _ = SearchWithWalk(pattern, absPath, &ui.Opts, stopChan)
	}

	return result, usingFd
}
