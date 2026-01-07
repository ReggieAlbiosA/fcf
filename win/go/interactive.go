package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// readLine reads a line of input from stdin
func readLine(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// getSearchPath prompts for and returns the search path (Step 1)
func getSearchPath() string {
	cwd, _ := os.Getwd()

	fmt.Printf("%s Enter path to search\n", colors.Bold("Step 1:"))
	fmt.Printf("%s\n", colors.Dim(fmt.Sprintf("(Press Enter for current directory: %s)", cwd)))
	fmt.Println()

	userPath := readLine(colors.Cyan("Path: "))

	if userPath == "" {
		fmt.Println(colors.Green("Using current directory"))
		return "."
	}

	// Expand ~ to home directory
	if strings.HasPrefix(userPath, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			userPath = strings.Replace(userPath, "~", home, 1)
		}
	}

	// Expand environment variables
	userPath = os.ExpandEnv(userPath)

	// Validate path exists
	info, err := os.Stat(userPath)
	if err != nil || !info.IsDir() {
		fmt.Printf("%s Directory '%s' does not exist\n", colors.Red("ERROR:"), userPath)
		readLine("Press Enter to try again...")
		return ""
	}

	fmt.Println()
	return userPath
}

// getPattern prompts for and returns the search pattern (Step 2)
func getPattern() string {
	fmt.Printf("%s Enter file/folder name or pattern to find\n", colors.Bold("Step 2:"))
	fmt.Printf("%s\n", colors.Dim("Examples: *.log, config, .env, src, *.js"))
	fmt.Println()

	pattern := readLine(colors.Cyan("Pattern: "))

	if pattern == "" {
		fmt.Printf("%s Pattern cannot be empty\n", colors.Red("ERROR:"))
		readLine("Press Enter to try again...")
		return ""
	}

	fmt.Println()
	return pattern
}

// selectResult prompts user to select a result for navigation (Step 3)
func selectResult(results []string) string {
	fmt.Println()
	fmt.Printf("%s Enter path to navigate to\n", colors.Bold("Step 3:"))
	fmt.Printf("%s\n", colors.Dim("(Enter a number from results, full path, or press Enter to skip)"))
	fmt.Println()

	navInput := readLine(colors.Cyan("Navigate to: "))

	if navInput == "" {
		fmt.Println(colors.Dim("Skipped navigation"))
		return ""
	}

	// Check if input is a number (result index)
	if index, err := strconv.Atoi(navInput); err == nil {
		idx := index - 1 // Convert to 0-based index
		if idx >= 0 && idx < len(results) {
			return results[idx]
		}
		fmt.Printf("%s Invalid result number\n", colors.Red("ERROR:"))
		return ""
	}

	// Expand ~ and environment variables
	if strings.HasPrefix(navInput, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			navInput = strings.Replace(navInput, "~", home, 1)
		}
	}
	navInput = os.ExpandEnv(navInput)

	return navInput
}

// showOptionsMenu displays the options menu and returns user choice
func showOptionsMenu() int {
	fmt.Println()
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(colors.Bold("Options:"))
	fmt.Printf("  %s Find again (new search)\n", colors.Cyan("[f]"))
	fmt.Printf("  %s Repeat search (same path)\n", colors.Cyan("[r]"))
	fmt.Printf("  %s Exit\n", colors.Cyan("[n]"))
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()

	choice := readLine(colors.Cyan("Choose: "))

	switch strings.ToLower(choice) {
	case "f":
		return 1 // Go to Step 1
	case "r":
		return 2 // Go to Step 2
	default:
		return 0 // Exit
	}
}

// runInteractiveMode runs the main interactive loop
func runInteractiveMode() {
	currentStep := 1
	var searchPath, pattern string

	for {
		// Reset results
		var results []string

		// Show header
		showHeader()

		// Step 1: Get search path
		if currentStep == 1 {
			searchPath = getSearchPath()
			if searchPath == "" {
				continue
			}
			currentStep = 2
		}

		// Step 2: Get pattern
		if currentStep == 2 {
			showHeader()
			pattern = getPattern()
			if pattern == "" {
				continue
			}
		}

		// Show header again before search
		showHeader()

		// Execute search
		startTime := getTime()
		results, _ = search(pattern, searchPath)
		elapsed := getTime() - startTime

		// Show summary
		showSummary(len(results), elapsed)

		// Step 3: Navigate to path
		if len(results) > 0 {
			targetPath := selectResult(results)
			if targetPath != "" {
				fmt.Println()
				navigateToPath(targetPath)
			}
		}

		// Show options menu
		choice := showOptionsMenu()

		switch choice {
		case 0: // Exit
			fmt.Println(colors.Green("Goodbye!"))
			return
		case 1: // Go to Step 1
			currentStep = 1
			searchPath = ""
			pattern = ""
		case 2: // Go to Step 2
			currentStep = 2
			pattern = ""
		}
	}
}

// getTime returns current time in seconds
func getTime() float64 {
	return float64(timeNow().UnixNano()) / 1e9
}
