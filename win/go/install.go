package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

// runInstall is the main entry point for the install subcommand
func runInstall() {
	initColors()

	fmt.Println(colors.Bold(colors.Cyan("╔════════════════════════════════════════╗")))
	fmt.Println(colors.Bold(colors.Cyan("║")) + "   " + colors.Bold("fcf") + " - Installation                  " + colors.Bold(colors.Cyan("║")))
	fmt.Println(colors.Bold(colors.Cyan("╚════════════════════════════════════════╝")))
	fmt.Println()

	// Check for elevated privileges
	if !isElevated() {
		fmt.Println(colors.Red("Error: Installation requires elevated privileges."))
		fmt.Println()
		if runtime.GOOS == "windows" {
			fmt.Println("Please run this command as Administrator:")
			fmt.Println(colors.Cyan("  Right-click PowerShell -> Run as Administrator"))
		} else {
			fmt.Println("Please run with sudo:")
			fmt.Println(colors.Cyan("  sudo ./fcf install"))
		}
		os.Exit(1)
	}

	// Get install path
	installPath := getInstallPath()
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

	// Show success message and shell integration instructions
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

// showInstallSuccess displays success message and shell integration instructions
func showInstallSuccess() {
	fmt.Println()
	fmt.Println(colors.Bold(colors.Green("Installation successful!")))
	fmt.Println()
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(colors.Bold("Shell Integration (required for navigation):"))
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()
	fmt.Println("Add the following to your shell configuration file:")
	fmt.Println()

	fmt.Println(getShellIntegration())

	fmt.Println()
	fmt.Println(colors.Dim("After adding, restart your shell or run:"))
	if runtime.GOOS == "windows" {
		fmt.Println(colors.Cyan("  . $PROFILE"))
	} else {
		fmt.Println(colors.Cyan("  source ~/.bashrc  # or ~/.zshrc"))
	}
	fmt.Println()

	// Show fd installation hint if not installed
	if !hasFd() {
		fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		fmt.Println(colors.Yellow("Optional: Install 'fd' for faster searching:"))
		fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		fmt.Println()
		fmt.Println(colors.Cyan("  " + getFdInstallHint()))
		fmt.Println()
	}
}

// getShellIntegration returns the shell wrapper function for the current OS
func getShellIntegration() string {
	if runtime.GOOS == "windows" {
		return getPowerShellIntegration()
	}
	return getBashZshIntegration()
}

// getBashZshIntegration returns the Bash/Zsh wrapper function
func getBashZshIntegration() string {
	return colors.Dim("# Add to ~/.bashrc or ~/.zshrc") + `
` + colors.Cyan(`fcf() {
    local nav_file="/tmp/fcf_nav_path"
    rm -f "$nav_file"
    command fcf "$@"
    if [[ -f "$nav_file" ]]; then
        local target
        target=$(cat "$nav_file")
        rm -f "$nav_file"
        if [[ -d "$target" ]]; then
            cd "$target" || return
        fi
    fi
}`)
}

// getPowerShellIntegration returns the PowerShell wrapper function
func getPowerShellIntegration() string {
	return colors.Dim("# Add to $PROFILE (run: notepad $PROFILE)") + `
` + colors.Cyan(`function fcf {
    $navFile = Join-Path $env:TEMP "fcf_nav_path"
    if (Test-Path $navFile) { Remove-Item $navFile -Force }
    & "C:\Program Files\fcf\fcf.exe" @args
    if (Test-Path $navFile) {
        $target = Get-Content $navFile -Raw
        Remove-Item $navFile -Force
        if (Test-Path $target -PathType Container) {
            Set-Location $target
        }
    }
}`)
}
