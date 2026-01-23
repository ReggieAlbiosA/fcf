package install

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/ReggieAlbiosA/fcf/internal/ui"
)

const (
	githubRepo    = "ReggieAlbiosA/fcf"
	githubAPIURL  = "https://api.github.com/repos/" + githubRepo + "/releases/latest"
	downloadURL   = "https://github.com/" + githubRepo + "/releases/download"
)

// GitHubRelease represents a GitHub release response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// RunUpdate is the main entry point for the update subcommand
func RunUpdate(currentVersion string) {
	ui.InitColors()

	fmt.Println(ui.Colors.Bold(ui.Colors.Cyan("╔════════════════════════════════════════╗")))
	fmt.Println(ui.Colors.Bold(ui.Colors.Cyan("║")) + "   " + ui.Colors.Bold("fcf") + " - Update                        " + ui.Colors.Bold(ui.Colors.Cyan("║")))
	fmt.Println(ui.Colors.Bold(ui.Colors.Cyan("╚════════════════════════════════════════╝")))
	fmt.Println()
	fmt.Printf("%s %s\n", ui.Colors.Blue("Current version:"), ui.Colors.Cyan("v"+currentVersion))

	// Check for latest version
	fmt.Printf("%s", ui.Colors.Yellow("Checking for updates... "))
	latestVersion, err := getLatestVersion()
	if err != nil {
		fmt.Println(ui.Colors.Red("FAILED"))
		fmt.Printf("%s %s\n", ui.Colors.Red("Error:"), err.Error())
		os.Exit(1)
	}
	fmt.Println(ui.Colors.Green("OK"))

	fmt.Printf("%s %s\n", ui.Colors.Blue("Latest version:"), ui.Colors.Cyan("v"+latestVersion))
	fmt.Println()

	// Compare versions
	if !isNewerVersion(latestVersion, currentVersion) {
		fmt.Println(ui.Colors.Green("You are already running the latest version!"))
		return
	}

	fmt.Printf("%s v%s -> v%s\n", ui.Colors.Yellow("Update available:"), currentVersion, latestVersion)
	fmt.Println()

	// Check for elevated privileges
	if !isElevated() {
		fmt.Println(ui.Colors.Red("Error: Update requires elevated privileges."))
		fmt.Println()
		if runtime.GOOS == "windows" {
			fmt.Println("Please run this command as Administrator:")
			fmt.Println(ui.Colors.Cyan("  Right-click PowerShell -> Run as Administrator"))
		} else {
			fmt.Println("Please run with sudo:")
			fmt.Println(ui.Colors.Cyan("  sudo fcf update"))
		}
		os.Exit(1)
	}

	// Download and install update
	if err := downloadAndInstall(latestVersion); err != nil {
		fmt.Printf("%s %s\n", ui.Colors.Red("Error:"), err.Error())
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(ui.Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(ui.Colors.Bold(ui.Colors.Green("Update complete!")))
	fmt.Println(ui.Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()
	fmt.Printf("fcf has been updated to %s\n", ui.Colors.Cyan("v"+latestVersion))
	fmt.Println()
}

// getLatestVersion fetches the latest version from GitHub releases
func getLatestVersion() (string, error) {
	resp, err := http.Get(githubAPIURL)
	if err != nil {
		return "", fmt.Errorf("could not fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("could not parse release info: %w", err)
	}

	// Remove 'v' prefix if present
	version := strings.TrimPrefix(release.TagName, "v")
	return version, nil
}

// isNewerVersion compares two version strings (simple comparison)
func isNewerVersion(latest, current string) bool {
	// Simple string comparison works for semver-like versions
	// For more robust comparison, use a semver library
	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	for i := 0; i < len(latestParts) && i < len(currentParts); i++ {
		if latestParts[i] > currentParts[i] {
			return true
		}
		if latestParts[i] < currentParts[i] {
			return false
		}
	}

	return len(latestParts) > len(currentParts)
}

// downloadAndInstall downloads the latest binary and replaces the current one
func downloadAndInstall(version string) error {
	binaryName := getBinaryName()
	url := fmt.Sprintf("%s/v%s/%s", downloadURL, version, binaryName)

	fmt.Printf("%s %s\n", ui.Colors.Yellow("Downloading:"), ui.Colors.Dim(url))

	// Download to temp file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "fcf-update-*")
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Copy downloaded content
	fmt.Printf("%s", ui.Colors.Yellow("Downloading binary... "))
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return fmt.Errorf("could not save update: %w", err)
	}
	tmpFile.Close()
	fmt.Println(ui.Colors.Green("OK"))

	// Make executable (Unix only)
	if err := makeExecutable(tmpPath); err != nil {
		return fmt.Errorf("could not make executable: %w", err)
	}

	// Replace current binary
	fmt.Printf("%s", ui.Colors.Yellow("Installing update... "))
	installPath := getInstallPath()

	// Remove old binary first (Windows requires this)
	os.Remove(installPath)

	if err := copyFile(tmpPath, installPath); err != nil {
		return fmt.Errorf("could not install update: %w", err)
	}

	if err := makeExecutable(installPath); err != nil {
		return fmt.Errorf("could not set permissions: %w", err)
	}

	fmt.Println(ui.Colors.Green("OK"))
	return nil
}

// getBinaryName returns the appropriate binary name for the current platform
func getBinaryName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch os {
	case "windows":
		return fmt.Sprintf("fcf-windows-%s.exe", arch)
	case "darwin":
		return fmt.Sprintf("fcf-darwin-%s", arch)
	default:
		return fmt.Sprintf("fcf-linux-%s", arch)
	}
}
