package command

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ReggieAlbiosA/fcf/internal/navigation"
	"github.com/ReggieAlbiosA/fcf/internal/search"
	"github.com/ReggieAlbiosA/fcf/internal/ui"
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

	fmt.Printf("%s Enter path to search\n", ui.Colors.Bold("Step 1:"))
	fmt.Printf("%s\n", ui.Colors.Dim(fmt.Sprintf("(Press Enter for current directory: %s)", cwd)))
	fmt.Println()

	userPath := readLine(ui.Colors.Cyan("Path: "))

	if userPath == "" {
		fmt.Println(ui.Colors.Green("Using current directory"))
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
		fmt.Printf("%s Directory '%s' does not exist\n", ui.Colors.Red("ERROR:"), userPath)
		readLine("Press Enter to try again...")
		return ""
	}

	fmt.Println()
	return userPath
}

// getPattern prompts for and returns the search pattern (Step 2)
func getPattern() string {
	fmt.Printf("%s Enter file/folder name or pattern to find\n", ui.Colors.Bold("Step 2:"))
	fmt.Printf("%s\n", ui.Colors.Dim("Examples: *.log, config, .env, src, *.js"))
	fmt.Println()

	pattern := readLine(ui.Colors.Cyan("Pattern: "))

	if pattern == "" {
		fmt.Printf("%s Pattern cannot be empty\n", ui.Colors.Red("ERROR:"))
		readLine("Press Enter to try again...")
		return ""
	}

	fmt.Println()
	return pattern
}

// SelectResult prompts user to select a result for navigation (Step 3)
func SelectResult(results []string) string {
	fmt.Println()
	fmt.Printf("%s Enter path to navigate to\n", ui.Colors.Bold("Step 3:"))
	fmt.Printf("%s\n", ui.Colors.Dim("(Enter a number from results, full path, or press Enter to skip)"))
	fmt.Println()

	navInput := readLine(ui.Colors.Cyan("Navigate to: "))

	if navInput == "" {
		fmt.Println(ui.Colors.Dim("Skipped navigation"))
		return ""
	}

	// Check if input is a number (result index)
	if index, err := strconv.Atoi(navInput); err == nil {
		idx := index - 1 // Convert to 0-based index
		if idx >= 0 && idx < len(results) {
			return results[idx]
		}
		fmt.Printf("%s Invalid result number\n", ui.Colors.Red("ERROR:"))
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
	fmt.Println(ui.Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(ui.Colors.Bold("Options:"))
	fmt.Printf("  %s Find again (new search)\n", ui.Colors.Cyan("[f]"))
	fmt.Printf("  %s Repeat search (same path)\n", ui.Colors.Cyan("[r]"))
	fmt.Printf("  %s Exit\n", ui.Colors.Cyan("[n]"))
	fmt.Println(ui.Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()

	choice := readLine(ui.Colors.Cyan("Choose: "))

	switch strings.ToLower(choice) {
	case "f":
		return 1 // Go to Step 1
	case "r":
		return 2 // Go to Step 2
	default:
		return 0 // Exit
	}
}

// RunInteractiveMode runs the main interactive loop
func RunInteractiveMode() {
	currentStep := 1
	var searchPath, pattern string

	for {
		// Reset results
		var results []string

		// Show header
		ui.ShowHeader()

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
			ui.ShowHeader()
			pattern = getPattern()
			if pattern == "" {
				continue
			}
		}

		// Show header again before search
		ui.ShowHeader()

		// Execute search
		startTime := getTime()
		results, _ = search.Search(pattern, searchPath)
		elapsed := getTime() - startTime

		// Show summary
		ui.ShowSummary(len(results), elapsed)

		// Step 3: Navigate to path
		if len(results) > 0 {
			targetPath := SelectResult(results)
			if targetPath != "" {
				fmt.Println()
				navigation.NavigateToPath(targetPath)
			}
		}

		// Show options menu
		choice := showOptionsMenu()

		switch choice {
		case 0: // Exit
			fmt.Println(ui.Colors.Green("Goodbye!"))
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
