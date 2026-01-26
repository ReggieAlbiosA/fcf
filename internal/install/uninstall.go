package install

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ReggieAlbiosA/fcf/internal/install/shell"
	"github.com/ReggieAlbiosA/fcf/internal/ui"
)

// RunUninstall is the main entry point for the uninstall subcommand
func RunUninstall() {
	ui.InitColors()

	fmt.Println(ui.Colors.Bold(ui.Colors.Cyan("╔════════════════════════════════════════╗")))
	fmt.Println(ui.Colors.Bold(ui.Colors.Cyan("║")) + "   " + ui.Colors.Bold("fcf") + " - Uninstallation                " + ui.Colors.Bold(ui.Colors.Cyan("║")))
	fmt.Println(ui.Colors.Bold(ui.Colors.Cyan("╚════════════════════════════════════════╝")))
	fmt.Println()

	// Check for elevated privileges
	if !isElevated() {
		fmt.Println(ui.Colors.Red("Error: Uninstallation requires elevated privileges."))
		fmt.Println()
		if runtime.GOOS == "windows" {
			fmt.Println("Please run this command as Administrator:")
			fmt.Println(ui.Colors.Cyan("  Right-click PowerShell -> Run as Administrator"))
		} else {
			fmt.Println("Please run with sudo:")
			fmt.Println(ui.Colors.Cyan("  sudo fcf uninstall"))
		}
		os.Exit(1)
	}

	// Confirm uninstallation
	fmt.Println(ui.Colors.Yellow("This will remove:"))
	fmt.Printf("  - Binary at %s\n", ui.Colors.Cyan(getInstallPath()))
	fmt.Println("  - Shell integration from your shell config files")
	fmt.Println()
	fmt.Print(ui.Colors.Bold("Are you sure you want to uninstall fcf? [y/N] "))

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		fmt.Println()
		fmt.Println(ui.Colors.Yellow("Uninstallation cancelled."))
		return
	}

	fmt.Println()

	// Remove shell integration first (while binary still exists)
	removeShellIntegration()

	// Remove binary
	fmt.Printf("%s", ui.Colors.Yellow("Removing binary... "))
	installPath := getInstallPath()
	if err := os.Remove(installPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(ui.Colors.Yellow("not found (already removed)"))
		} else {
			fmt.Println(ui.Colors.Red("FAILED"))
			fmt.Printf("%s %s\n", ui.Colors.Red("Error:"), err.Error())
		}
	} else {
		fmt.Println(ui.Colors.Green("OK"))
	}

	// Platform-specific cleanup
	postUninstall()

	// Show success message
	fmt.Println()
	fmt.Println(ui.Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(ui.Colors.Bold(ui.Colors.Green("Uninstallation complete!")))
	fmt.Println(ui.Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()
	fmt.Println(ui.Colors.Dim("Thank you for using fcf. To reinstall, visit:"))
	fmt.Println(ui.Colors.Cyan("  https://github.com/ReggieAlbiosA/fcf"))
	fmt.Println()
}

// removeShellIntegration removes the shell wrapper function from config files
func removeShellIntegration() {
	homeDir, err := getRealUserHomeDir()
	if err != nil {
		fmt.Printf("%s %s\n", ui.Colors.Yellow("Warning:"), "Could not determine home directory")
		return
	}

	fmt.Println(ui.Colors.Bold("Removing shell integration..."))
	fmt.Println()

	// Detect shells and remove integration
	shells := shell.DetectShellsForInstallation(homeDir)
	if len(shells) == 0 {
		// Try to detect from config files
		shells = shell.DetectShellsFromConfigFiles(homeDir)
	}

	for _, s := range shells {
		fmt.Printf("%s %s... ", ui.Colors.Yellow("Checking"), s.Name)

		if !shell.HasExistingInstallation(s.ConfigPath) {
			fmt.Println(ui.Colors.Dim("not installed"))
			continue
		}

		if err := shell.RemoveShellIntegration(s.ConfigPath); err != nil {
			fmt.Printf("%s %s\n", ui.Colors.Red("FAILED"), err.Error())
			continue
		}

		fmt.Println(ui.Colors.Green("removed"))
	}

	fmt.Println()
}
