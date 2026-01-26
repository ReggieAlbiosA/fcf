package install

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/ReggieAlbiosA/fcf/internal/install/shell"
	"github.com/ReggieAlbiosA/fcf/internal/platform"
	"github.com/ReggieAlbiosA/fcf/internal/search"
	"github.com/ReggieAlbiosA/fcf/internal/ui"
)

// runInstall is the main entry point for the install subcommand
func RunInstall() {
	ui.InitColors()

	// Parse install flags
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	shellOverride := fs.String("shell", "", "Override shell detection (bash, zsh, fish)")
	noShell := fs.Bool("no-shell", false, "Skip shell integration")
	shellOnly := fs.Bool("shell-only", false, "Only install shell integration (skip binary installation)")
	local := fs.Bool("local", false, "Install from current binary (for development testing)")
	force := fs.Bool("force", false, "Force install without prompts")
	fs.BoolVar(force, "f", false, "Force install without prompts (shorthand)")
	fs.Parse(os.Args[2:])

	fmt.Println(ui.Colors.Bold(ui.Colors.Cyan("╔════════════════════════════════════════╗")))
	fmt.Println(ui.Colors.Bold(ui.Colors.Cyan("║")) + "   " + ui.Colors.Bold("fcf") + " - Installation                  " + ui.Colors.Bold(ui.Colors.Cyan("║")))
	fmt.Println(ui.Colors.Bold(ui.Colors.Cyan("╚════════════════════════════════════════╝")))
	fmt.Println()

	// Shell-only mode: skip binary installation and privilege check
	if *shellOnly {
		if *noShell {
			fmt.Println(ui.Colors.Red("Error: --shell-only and --no-shell cannot be used together"))
			os.Exit(1)
		}
		fmt.Printf("%s %s\n", ui.Colors.Blue("Mode:"), ui.Colors.Cyan("Shell integration only"))
		fmt.Println()
		installShellIntegration(*shellOverride)
		showInstallSuccess()
		return
	}

	// Check for elevated privileges
	if !isElevated() {
		fmt.Println(ui.Colors.Red("Error: Installation requires elevated privileges."))
		fmt.Println()
		if runtime.GOOS == "windows" {
			fmt.Println("Please run this command as Administrator:")
			fmt.Println(ui.Colors.Cyan("  Right-click PowerShell -> Run as Administrator"))
		} else {
			fmt.Println("Please run with sudo:")
			fmt.Println(ui.Colors.Cyan("  sudo ./fcf install"))
		}
		os.Exit(1)
	}

	// Get install path
	installPath := getInstallPath()

	// Show install mode
	if *local {
		fmt.Printf("%s %s\n", ui.Colors.Blue("Mode:"), ui.Colors.Yellow("Local development"))
	}
	fmt.Printf("%s %s\n", ui.Colors.Blue("Install scope:"), ui.Colors.Cyan("System"))
	fmt.Printf("%s %s\n", ui.Colors.Blue("Install location:"), ui.Colors.Cyan(installPath))

	// Detect OS/distro
	fmt.Printf("%s %s\n", ui.Colors.Blue("Operating System:"), ui.Colors.Cyan(getOSInfo()))
	fmt.Println()

	// Check if already installed (unless --local or --force)
	if !*local && !*force {
		if _, err := os.Stat(installPath); err == nil {
			fmt.Println(ui.Colors.Yellow("fcf is already installed at this location."))
			fmt.Println()
			fmt.Print(ui.Colors.Bold("Do you want to overwrite? [y/N] "))

			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
				fmt.Println()
				fmt.Println(ui.Colors.Yellow("Installation cancelled."))
				fmt.Println(ui.Colors.Dim("Tip: Use --force to skip this prompt"))
				return
			}
			fmt.Println()
		}
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("%s %s\n", ui.Colors.Red("Error:"), "Could not determine executable path")
		os.Exit(1)
	}

	// Ensure install directory exists
	if err := ensureInstallDir(); err != nil {
		fmt.Printf("%s %s\n", ui.Colors.Red("Error:"), err.Error())
		os.Exit(1)
	}

	// Copy binary to install location
	fmt.Printf("%s", ui.Colors.Yellow("Installing fcf... "))
	if err := copyFile(execPath, installPath); err != nil {
		fmt.Println(ui.Colors.Red("FAILED"))
		fmt.Printf("%s %s\n", ui.Colors.Red("Error:"), err.Error())
		os.Exit(1)
	}

	// Make executable (Unix only, no-op on Windows)
	if err := makeExecutable(installPath); err != nil {
		fmt.Println(ui.Colors.Red("FAILED"))
		fmt.Printf("%s %s\n", ui.Colors.Red("Error:"), err.Error())
		os.Exit(1)
	}

	fmt.Println(ui.Colors.Green("OK"))

	// Platform-specific post-install (e.g., add to PATH on Windows)
	if err := postInstall(); err != nil {
		fmt.Printf("%s %s\n", ui.Colors.Yellow("Warning:"), err.Error())
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
		distro := platform.DetectLinuxDistro()
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
	homeDir, err := getRealUserHomeDir()
	if err != nil {
		fmt.Printf("%s %s\n", ui.Colors.Yellow("Warning:"), "Could not determine home directory for shell integration")
		return
	}

	fmt.Println()
	fmt.Println(ui.Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(ui.Colors.Bold("Configuring shell integration..."))
	fmt.Println(ui.Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()

	// Get shells to configure
	var shells []shell.ShellInfo
	if shellOverride != "" {
		// Override shell detection
		var shellType shell.ShellType
		switch shellOverride {
		case "bash":
			shellType = shell.ShellBash
		case "zsh":
			shellType = shell.ShellZsh
		case "fish":
			shellType = shell.ShellFish
		default:
			fmt.Printf("%s Unknown shell: %s\n", ui.Colors.Red("Error:"), shellOverride)
			return
		}
		configPath := shell.GetShellConfigPath(homeDir, shellType)
		if configPath != "" {
			shells = append(shells, shell.ShellInfo{
				Type:       shellType,
				Name:       shellType.String(),
				ConfigPath: configPath,
				Detected:   true,
			})
		}
	} else {
		// Auto-detect shells
		shells = shell.DetectShellsForInstallation(homeDir)
	}

	if len(shells) == 0 {
		fmt.Printf("%s No shells detected for configuration\n", ui.Colors.Yellow("Warning:"))
		return
	}

	// Configure each detected shell
	var successCount int
	for _, s := range shells {
		fmt.Printf("%s %s... ", ui.Colors.Yellow("Configuring"), s.Name)

		// Check if already installed
		if shell.HasExistingInstallation(s.ConfigPath) {
			fmt.Println(ui.Colors.Yellow("already configured"))
			continue
		}

		// Add shell integration
		if err := shell.AddShellIntegration(s.ConfigPath, s.Type); err != nil {
			fmt.Printf("%s %s\n", ui.Colors.Red("FAILED"), err.Error())
			continue
		}

		fmt.Println(ui.Colors.Green("OK"))
		successCount++
	}

	fmt.Println()

	if successCount > 0 {
		fmt.Println(ui.Colors.Bold(ui.Colors.Green("Shell integration installed successfully!")))
		fmt.Println()
		fmt.Println(ui.Colors.Bold("To enable the new function, reload your shell config:"))
		fmt.Println()
		for _, shellVar := range shells {
			if !shell.HasExistingInstallation(shellVar.ConfigPath) {
				continue
			}
			cmd := shell.GetShellReloadCommand(shellVar.Type)
			fmt.Printf("%s %s\n", ui.Colors.Cyan("  "+shellVar.Name+":"), cmd)
		}
		fmt.Println()
	}
}

// showInstallSuccess displays success message
func showInstallSuccess() {
	fmt.Println()
	fmt.Println(ui.Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(ui.Colors.Bold(ui.Colors.Green("Installation complete!")))
	fmt.Println(ui.Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()

	// Show fd installation hint if not installed
	if !search.HasFd() {
		fmt.Println(ui.Colors.Bold("Optional: Install 'fd' for faster searching:"))
		fmt.Println()
		fmt.Println(ui.Colors.Cyan("  " + platform.GetFdInstallHint()))
		fmt.Println()
	}

	fmt.Println(ui.Colors.Dim("For more information, run: fcf --help"))
	fmt.Println()
}
