package command

import (
	"fmt"
	"runtime"

	"github.com/ReggieAlbiosA/fcf/internal/platform"
	"github.com/ReggieAlbiosA/fcf/internal/ui"
)

// ShowHelp displays the help message
func ShowHelp() {
	fmt.Printf("%s v%s\n", ui.Colors.Bold("fcf - Find File or Folder"), Version)
	fmt.Println()
	fmt.Println(ui.Colors.Bold("USAGE:"))
	fmt.Println("    fcf [OPTIONS] [PATTERN] [PATH]")
	fmt.Println("    fcf                          # Interactive mode")
	fmt.Println("    fcf install                  # Install fcf system-wide")
	fmt.Println("    fcf update                   # Update to latest version")
	fmt.Println("    fcf uninstall                # Remove fcf from system")
	fmt.Println()
	fmt.Println(ui.Colors.Bold("DESCRIPTION:"))
	fmt.Println("    Interactive tool to find files and folders with pattern matching")
	fmt.Println("    and real-time streaming results. Uses parallel search for speed.")
	fmt.Println()
	fmt.Println(ui.Colors.Bold("OPTIONS:"))
	fmt.Printf("    %s               Show this help message\n", ui.Colors.Cyan("-h, --help"))
	fmt.Printf("    %s                  Case-insensitive pattern matching\n", ui.Colors.Cyan("-i"))
	fmt.Printf("    %s           Filter by type: %s(file) or %s(directory)\n",
		ui.Colors.Cyan("-t TYPE"), ui.Colors.Yellow("f"), ui.Colors.Yellow("d"))
	fmt.Printf("    %s           Display file sizes\n", ui.Colors.Cyan("--show-size"))
	fmt.Printf("    %s    Maximum results to display (default: unlimited)\n", ui.Colors.Cyan("--max-display NUM"))
	fmt.Printf("    %s                  Skip navigation after search\n", ui.Colors.Cyan("-S"))
	fmt.Println()
	fmt.Println(ui.Colors.Bold("COMMANDS:"))
	fmt.Printf("    %s              Install/update fcf system-wide (requires sudo/admin)\n", ui.Colors.Cyan("install"))
	fmt.Printf("    %s               Update fcf to the latest version from GitHub\n", ui.Colors.Cyan("update"))
	fmt.Printf("    %s            Remove fcf from system\n", ui.Colors.Cyan("uninstall"))
	fmt.Println()
	fmt.Println(ui.Colors.Bold("SHELL INTEGRATION (for navigation to work):"))
	showShellIntegrationHelp()
	fmt.Println()
	fmt.Println(ui.Colors.Bold("EXAMPLES:"))
	fmt.Printf("    %s\n", ui.Colors.Green("# Interactive mode"))
	fmt.Println("    fcf")
	fmt.Println()
	fmt.Printf("    %s\n", ui.Colors.Green("# Find all .log files"))
	fmt.Println("    fcf \"*.log\"")
	fmt.Println()
	fmt.Printf("    %s\n", ui.Colors.Green("# Find in specific directory"))
	showExamplePath()
	fmt.Println()
	fmt.Printf("    %s\n", ui.Colors.Green("# Case-insensitive search for PNG files"))
	fmt.Println("    fcf -i \"*.PNG\"")
	fmt.Println()
	fmt.Printf("    %s\n", ui.Colors.Green("# Find only directories named 'src'"))
	fmt.Println("    fcf -t d src")
	fmt.Println()
	fmt.Printf("    %s\n", ui.Colors.Green("# Find with file sizes"))
	fmt.Println("    fcf --show-size \"*.mp4\"")
	fmt.Println()
	fmt.Println(ui.Colors.Bold("INTERACTIVE WORKFLOW:"))
	fmt.Println("    Step 1: Enter path to search")
	fmt.Println("    Step 2: Enter pattern to find")
	fmt.Println("    Step 3: Navigate to a result path")
	fmt.Println()
	fmt.Println(ui.Colors.Bold("NAVIGATION OPTIONS:"))
	fmt.Println("    After navigation, choose:")
	fmt.Printf("    %s - Find again (restart from Step 1)\n", ui.Colors.Cyan("f"))
	fmt.Printf("    %s - Repeat search (go to Step 2, same path)\n", ui.Colors.Cyan("r"))
	fmt.Printf("    %s - Exit\n", ui.Colors.Cyan("n"))
	fmt.Println()
	fmt.Println(ui.Colors.Bold("PERFORMANCE:"))
	fmt.Println("    - Uses 'fd' for fast parallel searching (if installed)")
	fmt.Println("    - Falls back to Go's filepath.WalkDir if fd is not available")
	fmt.Printf("    - Install fd: %s\n", ui.Colors.Cyan(platform.GetFdInstallHint()))
	fmt.Println()
}

// showShellIntegrationHelp displays platform-specific shell integration instructions
func showShellIntegrationHelp() {
	if runtime.GOOS == "windows" {
		fmt.Printf("    Add this to your %s:\n", ui.Colors.Cyan("$PROFILE"))
		fmt.Println()
		fmt.Println(ui.Colors.Yellow("    function fcf {"))
		fmt.Println(ui.Colors.Yellow("        $navFile = Join-Path $env:TEMP \"fcf_nav_path\""))
		fmt.Println(ui.Colors.Yellow("        if (Test-Path $navFile) { Remove-Item $navFile -Force }"))
		fmt.Println(ui.Colors.Yellow("        & \"C:\\Program Files\\fcf\\fcf.exe\" @args"))
		fmt.Println(ui.Colors.Yellow("        if (Test-Path $navFile) {"))
		fmt.Println(ui.Colors.Yellow("            $target = Get-Content $navFile -Raw"))
		fmt.Println(ui.Colors.Yellow("            Remove-Item $navFile -Force"))
		fmt.Println(ui.Colors.Yellow("            if (Test-Path $target -PathType Container) {"))
		fmt.Println(ui.Colors.Yellow("                Set-Location $target"))
		fmt.Println(ui.Colors.Yellow("            }"))
		fmt.Println(ui.Colors.Yellow("        }"))
		fmt.Println(ui.Colors.Yellow("    }"))
	} else {
		fmt.Printf("    Add this to your %s or %s:\n", ui.Colors.Cyan("~/.bashrc"), ui.Colors.Cyan("~/.zshrc"))
		fmt.Println()
		fmt.Println(ui.Colors.Yellow("    fcf() {"))
		fmt.Println(ui.Colors.Yellow("        local nav_file=\"/tmp/fcf_nav_path_$(id -u)\""))
		fmt.Println(ui.Colors.Yellow("        rm -f \"$nav_file\""))
		fmt.Println(ui.Colors.Yellow("        command fcf \"$@\""))
		fmt.Println(ui.Colors.Yellow("        if [[ -f \"$nav_file\" ]]; then"))
		fmt.Println(ui.Colors.Yellow("            local target"))
		fmt.Println(ui.Colors.Yellow("            target=$(cat \"$nav_file\")"))
		fmt.Println(ui.Colors.Yellow("            rm -f \"$nav_file\""))
		fmt.Println(ui.Colors.Yellow("            if [[ -d \"$target\" ]]; then"))
		fmt.Println(ui.Colors.Yellow("                cd \"$target\" || return"))
		fmt.Println(ui.Colors.Yellow("            fi"))
		fmt.Println(ui.Colors.Yellow("        fi"))
		fmt.Println(ui.Colors.Yellow("    }"))
	}
}

// showExamplePath displays a platform-appropriate example path
func showExamplePath() {
	if runtime.GOOS == "windows" {
		fmt.Println("    fcf \"*.js\" C:\\Projects")
	} else {
		fmt.Println("    fcf \"*.js\" ~/projects")
	}
}
