package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// runInstall is the main entry point for the install subcommand
func runInstall() {
	initColors()

	// Parse install flags
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	userScope := fs.Bool("user", false, "Install for current user only (no elevated privileges required)")
	shellOverride := fs.String("shell", "", "Override shell detection (bash, zsh, fish)")
	noShell := fs.Bool("no-shell", false, "Skip shell integration")
	fs.Parse(os.Args[2:])

	fmt.Println(colors.Bold(colors.Cyan("╔════════════════════════════════════════╗")))
	fmt.Println(colors.Bold(colors.Cyan("║")) + "   " + colors.Bold("fcf") + " - Installation                  " + colors.Bold(colors.Cyan("║")))
	fmt.Println(colors.Bold(colors.Cyan("╚════════════════════════════════════════╝")))
	fmt.Println()

	// Check for elevated privileges (skip for user-scope install)
	if !*userScope && !isElevated() {
		fmt.Println(colors.Red("Error: System-wide installation requires elevated privileges."))
		fmt.Println()
		if runtime.GOOS == "windows" {
			fmt.Println("Please run this command as Administrator:")
			fmt.Println(colors.Cyan("  Right-click PowerShell -> Run as Administrator"))
		} else {
			fmt.Println("Please run with sudo:")
			fmt.Println(colors.Cyan("  sudo ./fcf install"))
		}
		fmt.Println()
		fmt.Println("Alternatively, install for current user only:")
		fmt.Println(colors.Cyan("  ./fcf install --user"))
		os.Exit(1)
	}

	// Get install path
	var installPath string
	if *userScope && runtime.GOOS != "windows" {
		// User-scope installation on Unix
		if err := ensureUserBinDirectory(); err != nil {
			fmt.Printf("%s %s\n", colors.Red("Error:"), err.Error())
			os.Exit(1)
		}
		homeDir, _ := os.UserHomeDir()
		installPath = filepath.Join(homeDir, ".local", "bin", "fcf")
	} else {
		installPath = getInstallPath()
	}

	fmt.Printf("%s %s\n", colors.Blue("Install scope:"), colors.Cyan(map[bool]string{true: "User", false: "System"}[*userScope]))
	fmt.Printf("%s %s\n", colors.Blue("Install location:"), colors.Cyan(installPath))

	// Detect OS/distro
	fmt.Printf("%s %s\n", colors.Blue("Operating System:"), colors.Cyan(getOSInfo()))
	fmt.Println()

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("%s %s\n", colors.Red("Error:"), "Could not determine executable path")
		os.Exit(1)
	}

	// Ensure install directory exists
	if err := ensureInstallDir(); err != nil {
		fmt.Printf("%s %s\n", colors.Red("Error:"), err.Error())
		os.Exit(1)
	}

	// Copy binary to install location
	fmt.Printf("%s", colors.Yellow("Installing fcf... "))
	if err := copyFile(execPath, installPath); err != nil {
		fmt.Println(colors.Red("FAILED"))
		fmt.Printf("%s %s\n", colors.Red("Error:"), err.Error())
		os.Exit(1)
	}

	// Make executable (Unix only, no-op on Windows)
	if err := makeExecutable(installPath); err != nil {
		fmt.Println(colors.Red("FAILED"))
		fmt.Printf("%s %s\n", colors.Red("Error:"), err.Error())
		os.Exit(1)
	}

	fmt.Println(colors.Green("OK"))

	// Platform-specific post-install (e.g., add to PATH on Windows)
	if err := postInstall(); err != nil {
		fmt.Printf("%s %s\n", colors.Yellow("Warning:"), err.Error())
	}

	// Add user bin to PATH if needed (user-scope on Unix)
	if *userScope && runtime.GOOS != "windows" && !isUserBinInPath() {
		homeDir, _ := os.UserHomeDir()
		if err := addUserBinToPath(homeDir); err != nil {
			fmt.Printf("%s %s\n", colors.Yellow("Warning:"), err.Error())
		}
	}

	// Shell integration
	if !*noShell {
		installShellIntegration(*shellOverride)
	}

	// Show success message
	showInstallSuccess()
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("could not copy file: %w", err)
	}

	return nil
}

// getOSInfo returns a human-readable OS description
func getOSInfo() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows"
	case "darwin":
		return "macOS"
	case "linux":
		distro := detectLinuxDistro()
		if distro.Version != "" {
			return fmt.Sprintf("%s %s", distro.Name, distro.Version)
		}
		return distro.Name
	default:
		return runtime.GOOS
	}
}

// installShellIntegration detects shells and installs wrapper functions
func installShellIntegration(shellOverride string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%s %s\n", colors.Yellow("Warning:"), "Could not determine home directory for shell integration")
		return
	}

	fmt.Println()
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(colors.Bold("Configuring shell integration..."))
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()

	// Get shells to configure
	var shells []ShellInfo
	if shellOverride != "" {
		// Override shell detection
		var shellType ShellType
		switch shellOverride {
		case "bash":
			shellType = ShellBash
		case "zsh":
			shellType = ShellZsh
		case "fish":
			shellType = ShellFish
		default:
			fmt.Printf("%s Unknown shell: %s\n", colors.Red("Error:"), shellOverride)
			return
		}
		configPath := getShellConfigPath(homeDir, shellType)
		if configPath != "" {
			shells = append(shells, ShellInfo{
				Type:       shellType,
				Name:       shellType.String(),
				ConfigPath: configPath,
				Detected:   true,
			})
		}
	} else {
		// Auto-detect shells
		shells = detectShellsForInstallation(homeDir)
	}

	if len(shells) == 0 {
		fmt.Printf("%s No shells detected for configuration\n", colors.Yellow("Warning:"))
		return
	}

	// Configure each detected shell
	var successCount int
	for _, shell := range shells {
		fmt.Printf("%s %s... ", colors.Yellow("Configuring"), shell.Name)

		// Check if already installed
		if hasExistingInstallation(shell.ConfigPath) {
			fmt.Println(colors.Yellow("already configured"))
			continue
		}

		// Add shell integration
		if err := addShellIntegration(shell.ConfigPath, shell.Type); err != nil {
			fmt.Printf("%s %s\n", colors.Red("FAILED"), err.Error())
			continue
		}

		fmt.Println(colors.Green("OK"))
		successCount++
	}

	fmt.Println()

	if successCount > 0 {
		fmt.Println(colors.Bold(colors.Green("Shell integration installed successfully!")))
		fmt.Println()
		fmt.Println(colors.Bold("To enable the new function, reload your shell config:"))
		fmt.Println()
		for _, shell := range shells {
			if !hasExistingInstallation(shell.ConfigPath) {
				continue
			}
			cmd := getShellReloadCommand(shell.Type)
			fmt.Printf("%s %s\n", colors.Cyan("  "+shell.Name+":"), cmd)
		}
		fmt.Println()
	}
}

// showInstallSuccess displays success message
func showInstallSuccess() {
	fmt.Println()
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(colors.Bold(colors.Green("Installation complete!")))
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()

	// Show fd installation hint if not installed
	if !hasFd() {
		fmt.Println(colors.Bold("Optional: Install 'fd' for faster searching:"))
		fmt.Println()
		fmt.Println(colors.Cyan("  " + getFdInstallHint()))
		fmt.Println()
	}

	fmt.Println(colors.Dim("For more information, run: fcf --help"))
	fmt.Println()
}
