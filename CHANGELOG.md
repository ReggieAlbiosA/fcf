# Changelog

All notable changes to FCF (Find File or Folder) will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [3.0.0] - 2026-01-14

### Added
- **Unified cross-platform codebase** - Single Go implementation for Linux, macOS, and Windows
- **Self-installing binary** - New `fcf install` command for automated system-wide installation
- **Linux distro detection** - Automatic detection via /etc/os-release for distro-specific behavior
- **fdfind alias support** - Automatic support for Debian/Ubuntu's fdfind package naming
- **Platform-aware help text** - Installation instructions adapt to the current operating system
- **Multi-platform CI/CD** - GitHub Actions builds for 5 targets: linux/amd64, darwin/amd64, darwin/arm64, windows/amd64, windows/386

### Changed
- **Installation paths** - System-wide installation to /usr/local/bin (Unix) or C:\Program Files\fcf (Windows)
- **Build architecture** - Introduced platform-specific build tags for Unix/Windows code paths
- **Executable resolution** - Platform-aware handling for fd/fdfind executable detection
- **File permissions** - Automatic executable bit setting on Unix systems during installation

### Architecture
- Separated platform-specific logic with build tags (+build unix/windows)
- New files: distro_unix.go, distro_windows.go, executable_unix.go, executable_windows.go
- Modular installation system: install.go, install_unix.go, install_windows.go
- Updated display.go, search.go, main.go for cross-platform compatibility

### Build System
- Enhanced GitHub Actions workflow for multi-platform releases
- Automated builds for Linux (amd64), macOS (amd64, arm64), Windows (amd64, 386)
- Streamlined release artifact generation

---

## [2.0.2] – 2026-01-07

### Fixed
- **GitHub Actions Windows build paths**  
  Corrected workflow paths after Go source reorganization (win/go → go) to ensure Windows binaries are built and uploaded correctly.

- **Temporary navigation file conflicts**  
  Standardized navigation path to /tmp/fcf_nav_path, preventing permission and multi-user conflicts on Linux.

- **System installation reliability**  
  Installer now explicitly uses sudo when installing to system directories, avoiding silent permission failures.

### Changed
- **Go project structure**  
  Unified Go source location by moving Windows Go files from win/go/ to go/ for cleaner cross-platform maintenance.

- **Linux navigation handling**  
  Simplified navigation logic by removing root/sudo detection and related warning indicators.

- **Installer behavior (Linux)**  
  Installation now runs both user-level and system-level installs automatically, instead of branching based on root detection.

- **Shell integration responsibility**  
  Removed automatic shell wrapper injection from the installer; shell integration is now documented and user-managed.

- **Output cleanliness**  
  Removed sudo-lock indicators (🔒) and related legends for clearer search output.

### Removed
- Root/sudo detection logic from [fcf.sh](http://fcf.sh)
- Permission-check helpers and sudo navigation warnings
- Automatic shell function injection into .bashrc / .zshrc

### Internal
- Refactored path handling to always resolve absolute paths before navigation
- Reduced installer complexity by eliminating conditional execution paths
- Improved consistency between CI, installer, and runtime behavior

---

## [2.0.1] - 2024-12-30

### Fixed
- **fd auto-installation** - Now properly detects installation success using `fd --version` and checks exit code
- **False success messages** - Installer no longer reports "fd installed" when installation actually failed

### Added
- **PATH auto-refresh** - After fd installation, PATH is refreshed automatically (no manual PowerShell restart needed)
- **Version check** - Installer skips download if same version already installed ("fcf is already up to date")
- **CI push trigger** - GitHub Actions now triggers on push to any branch for testing

### Changed
- Simplified fd installation to use winget only (most common on Windows)
- Removed package manager selection prompt (was overly complex)

---

## [2.0.0] - 2024-12-30

### Added
- **Go rewrite for Windows** - Native binary replaces PowerShell script
- GitHub Actions workflow for automated builds
- Automatic release binary uploads to GitHub Releases
- Legacy upgrade support (seamless migration from v1.x)

### Changed
- Windows version now distributed as `fcf.exe` (Go binary)
- Installer downloads from GitHub Releases instead of raw source
- Profile function updated to call `.exe` instead of `.ps1`
- Search fallback uses Go's `filepath.WalkDir` instead of `Get-ChildItem`

### Removed
- `win/fcf.ps1` - Replaced by Go binary

### Fixed
- ANSI color issues on PowerShell 5.1 (native Go handles colors correctly)
- No more PowerShell execution policy issues

### Migration
Users with v1.x installed will be automatically upgraded:
- Old `fcf.ps1` file is deleted
- Profile function is updated to point to `fcf.exe`
- All settings and workflow remain the same

---

## [1.0.1] - 2024-12-30

### Fixed
- Resolve PSStyle property error on PowerShell 5.1
- Use `$([char]27)` for ANSI escape codes (PS 5.1 compatible)

---

## [1.0.0] - 2024-12-30

### Added
- Initial release of FCF (Find File or Folder)
- Interactive 3-step workflow (Path → Pattern → Navigate)
- Real-time streaming search results
- `fd` integration for fast parallel search
- Fallback to native search when `fd` not available
- File type icons and color coding
- Cross-platform support (Ubuntu/Linux and Windows)
- Installation scripts for both platforms
- Shell integration for directory navigation

### Platforms
- **Ubuntu/Linux**: `fcf.sh` (Bash script)
- **Windows**: `fcf.ps1` (PowerShell script)
