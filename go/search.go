package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// hasFd checks if fd is available in PATH
func hasFd() bool {
	_, err := exec.LookPath("fd")
	return err == nil
}

// searchWithFd uses fd for fast parallel search
func searchWithFd(pattern, searchPath string, opts *Options) ([]string, error) {
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

	cmd := exec.Command("fd", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var results []string
	scanner := bufio.NewScanner(stdout)
	count := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		count++
		results = append(results, line)

		// Display result in real-time (streaming)
		if opts.MaxDisplay == 0 || count <= opts.MaxDisplay {
			showResult(line, count)
		}
	}

	cmd.Wait()
	return results, nil
}

// searchWithWalk uses filepath.WalkDir as fallback
func searchWithWalk(pattern, searchPath string, opts *Options) ([]string, error) {
	var results []string
	count := 0

	err := filepath.WalkDir(searchPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
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
		results = append(results, path)

		// Display result in real-time (streaming)
		if opts.MaxDisplay == 0 || count <= opts.MaxDisplay {
			showResult(path, count)
		}

		return nil
	})

	return results, err
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
func search(pattern, searchPath string) ([]string, bool) {
	// Resolve search path
	absPath, err := filepath.Abs(searchPath)
	if err != nil {
		absPath = searchPath
	}

	usingFd := hasFd()
	showSearchInfo(absPath, pattern, usingFd)

	fmt.Printf("%s %s\n\n", colors.Bold("Results:"), colors.Dim("(streaming in real-time...)"))

	var results []string
	if usingFd {
		results, _ = searchWithFd(pattern, absPath, &opts)
	} else {
		results, _ = searchWithWalk(pattern, absPath, &opts)
	}

	return results, usingFd
}
