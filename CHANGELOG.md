# Changelog

All notable changes to FCF (Find File or Folder) will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
