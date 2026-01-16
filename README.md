# FCF (Find File or Folder)

A fast, interactive file and folder finder with parallel search and real-time streaming results.

## What is FCF?

FCF is an interactive command-line tool for finding files and folders with advanced pattern matching, real-time streaming results, and easy navigation. It uses parallel search for speed and provides a seamless workflow for locating and navigating to files.

Think of it as `find` with a friendly interface - interactive prompts, color-coded output, and instant navigation.

## Key Features

- **Parallel Search** - Uses `fd` for blazing fast parallel file searching
- **Real-Time Streaming** - Results appear line-by-line as they're found
- **Interactive Mode** - Step-by-step guided search workflow
- **Navigation** - Jump directly to any result or path
- **Pattern Matching** - Glob patterns, partial names, extensions
- **Color-Coded Output** - Visual distinction for folders, files, executables, symlinks
- **Loop Workflow** - Search again without restarting

## Quick Start

### Installation

FCF is a compiled Go binary with cross-platform support. Download the binary for your platform and run `fcf install` to auto-detect your shell and configure everything.

#### Linux

**Option 1: System-wide installation (recommended)**

```bash
# Download binary for your architecture (amd64 is most common)
curl -sSL https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-linux-amd64 -o fcf
chmod +x fcf

# Install and configure (requires sudo)
sudo ./fcf install
```

**Option 2: User-only installation (no sudo required)**

```bash
# Download binary
curl -sSL https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-linux-amd64 -o fcf
chmod +x fcf

# Install for current user only
./fcf install --user
```

#### macOS

**Option 1: System-wide installation (recommended)**

```bash
# Download for your Mac architecture
# Apple Silicon (M1/M2/M3):
curl -sSL https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-darwin-arm64 -o fcf

# OR Intel Mac:
# curl -sSL https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-darwin-amd64 -o fcf

chmod +x fcf

# Install and configure (requires sudo)
sudo ./fcf install
```

**Option 2: User-only installation (no sudo required)**

```bash
# Download for your Mac architecture (see above)
chmod +x fcf

# Install for current user only
./fcf install --user
```

#### Windows (PowerShell)

Run PowerShell as Administrator:

```powershell
# Download binary
Invoke-WebRequest -Uri "https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-windows-amd64.exe" -OutFile fcf.exe

# Install and configure
.\fcf.exe install
```

The `fcf install` command automatically:
- Detects your shell (bash, zsh, fish on Unix; PowerShell on Windows)
- Adds the required shell wrapper function for navigation
- Configures PATH if needed
- Provides reload instructions

## Usage

### Interactive Mode

Simply run `fcf` without arguments:

```bash
fcf
```

This starts the interactive workflow:

```
╔════════════════════════════════════════╗
║   fcf - Find File or Folder            ║
╚════════════════════════════════════════╝

Step 1: Enter path to search
(Press Enter for current directory)
Path: ~/projects

Step 2: Enter file/folder name or pattern to find
Examples: *.log, config, .env, src, *.js
Pattern: *.ts

Results: (streaming in real-time...)
  [1] 📄 ./src/index.ts
  [2] 📄 ./src/utils/helpers.ts
  [3] 📄 ./src/components/App.ts

Found 3 match(es) in 0.05s

Step 3: Enter path to navigate to
Navigate to: 1

✓ Navigated to: ./src

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Options:
  [f] Find again (new search)
  [r] Repeat search (same path)
  [n] Exit
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Direct Mode

Skip the prompts by providing arguments:

```bash
# Find all .log files in current directory
fcf "*.log"

# Find in specific directory
fcf "*.js" ~/projects

# Case-insensitive search
fcf -i "*.PNG"

# Find only directories
fcf -t d src

# Find only files
fcf -t f config

# Include hidden files
fcf -H ".env*"

# Show file sizes
fcf --show-size "*.mp4"
```

## Options

| Option | Description |
|--------|-------------|
| `-h, --help` | Show help message |
| `-i, --ignore-case` | Case-insensitive pattern matching |
| `-t, --type TYPE` | Filter by type: `f` (file) or `d` (directory) |
| `-H, --hidden` | Include hidden files/folders |
| `--show-size` | Display file sizes |
| `--max-display NUM` | Maximum results to display |

## Interactive Workflow

### Step 1: Path Selection
Enter the directory to search in. Press Enter to use the current directory.

### Step 2: Pattern Input
Enter the file/folder name or pattern to find:
- Glob patterns: `*.txt`, `*.log`, `dist`
- Partial names: `config`, `test`, `README`
- Extensions: `.js`, `.py`, `.sh`

### Step 3: Navigation
Choose a result to navigate to:
- Enter a **number** (e.g., `3`) to navigate to result #3
- Enter a **full path** to navigate anywhere
- Press **Enter** to skip navigation

### Options Menu
After navigation, choose your next action:
- `f` - Find again (restart from Step 1 with new path)
- `r` - Repeat search (same path, new pattern)
- `n` - Exit

## Output Icons

| Icon | Type |
|------|------|
| 📁 | Directory |
| 📄 | Regular file |
| ⚡ | Executable |
| 🔗 | Symbolic link |

## Performance

FCF uses `fd` for parallel searching when available:

| Tool | Type | Speed |
|------|------|-------|
| `fd` | Parallel | 5-10x faster |
| `find` / `Get-ChildItem` | Sequential | Fallback |

The installer will offer to install `fd` automatically. You can also install it manually:

### Ubuntu/Linux

```bash
# Debian/Ubuntu
sudo apt install fd-find

# Fedora
sudo dnf install fd-find

# Arch Linux
sudo pacman -S fd

# macOS
brew install fd
```

### Windows

```powershell
# Using winget (recommended)
winget install sharkdp.fd

# Using Chocolatey
choco install fd

# Using Scoop
scoop install fd
```

## Examples

### Find configuration files
```bash
fcf "*.config.*"
fcf -i "config"
```

### Find all TypeScript files
```bash
fcf "*.ts" ~/projects
fcf "*.tsx" ./src
```

### Find directories named 'test'
```bash
fcf -t d test
fcf -t d "__tests__"
```

### Find large video files
```bash
fcf --show-size "*.mp4"
fcf --show-size "*.mkv" ~/Videos
```

### Find hidden files
```bash
fcf -H ".env*"
fcf -H ".git*"
```

## Installation Details

### What Gets Installed

FCF is a **compiled Go binary** with automatic shell integration configuration.

**Installed Binary:**
- Linux/macOS: `~/.local/bin/fcf` (user) or `/usr/local/bin/fcf` (system)
- Windows: `%USERPROFILE%\.local\bin\fcf.exe` (user) or `C:\Program Files\fcf\fcf.exe` (system)

**Shell Configuration:**
- Automatically detected: Bash, Zsh, Fish (Unix/Linux/macOS), PowerShell (Windows)
- Wrapper function added to config file for directory navigation
- Idempotent: Won't duplicate on reinstall

### Installation Scopes

| Scope | Linux/macOS | Windows | Requires Sudo |
|-------|-------------|---------|---------------|
| **System-wide** | `/usr/local/bin/fcf` | `C:\Program Files\fcf\fcf.exe` | Yes |
| **User-only** | `~/.local/bin/fcf` | `%USERPROFILE%\.local\bin\fcf.exe` | No |

### Command Options

```bash
fcf install              # System-wide (requires sudo on Unix)
fcf install --user      # User-only (no sudo required)
fcf install --no-shell  # Binary only, skip shell integration
fcf install --shell zsh # Force specific shell configuration
```

### Supported Shells

| Shell | Config File | Platform | Auto-Detected |
|-------|------------|----------|---------------|
| Bash | `~/.bashrc` | Linux/macOS | Yes |
| Bash | `~/.bash_profile` | macOS | Yes |
| Zsh | `~/.zshrc` | Linux/macOS | Yes |
| Fish | `~/.config/fish/config.fish` | Linux/macOS | Yes |
| PowerShell | `$PROFILE` | Windows | Yes |

## Updating

To update to the latest version, re-download the binary and run `fcf install` again:

```bash
# Linux/macOS
curl -sSL https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-linux-amd64 -o fcf
chmod +x fcf
sudo ./fcf install       # or ./fcf install --user
```

```powershell
# Windows
Invoke-WebRequest -Uri "https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-windows-amd64.exe" -OutFile fcf.exe
.\fcf.exe install
```

The installer detects your existing installation and updates the shell configuration if needed.

## Available Binaries

Pre-compiled binaries for all platforms are available on the [GitHub Releases](https://github.com/ReggieAlbiosA/fcf/releases) page.

**Linux:**
- `fcf-linux-amd64` - Intel/AMD 64-bit
- `fcf-linux-arm64` - ARM 64-bit

**macOS:**
- `fcf-darwin-amd64` - Intel
- `fcf-darwin-arm64` - Apple Silicon (M1/M2/M3)

**Windows:**
- `fcf-windows-amd64.exe` - 64-bit
- `fcf-windows-386.exe` - 32-bit

Download the appropriate binary for your platform and run `fcf install` to complete setup.

## Uninstallation

### Linux/macOS

```bash
# User installation
rm ~/.local/bin/fcf

# System-wide installation
sudo rm /usr/local/bin/fcf

# Optional: Remove shell integration from config files
# Edit ~/.bashrc, ~/.zshrc, or ~/.config/fish/config.fish
# and remove the fcf function block
```

### Windows (PowerShell)

```powershell
# User installation
Remove-Item "$env:USERPROFILE\.local\bin\fcf.exe" -Force

# System-wide installation (as Administrator)
Remove-Item "C:\Program Files\fcf" -Recurse -Force

# Optional: Remove shell integration from PowerShell profile
# Edit $PROFILE and remove the fcf function block
```

**Note:** Shell integration functions are marked with unique comments and can be safely removed manually. They won't interfere with FCF if you reinstall later.

## Troubleshooting

### General

**Shell integration not working?**

Reload your shell configuration after installation:

```bash
# Bash
source ~/.bashrc

# Zsh
source ~/.zshrc

# Fish
source ~/.config/fish/config.fish

# PowerShell
. $PROFILE
```

**Reinstall shell integration:**

```bash
# Manually reconfigure the installed shell
fcf install --shell bash   # or zsh, fish
```

### Linux/macOS

**Command not found after installation:**

- For system-wide: Restart terminal or source your shell config
- For user-only: Ensure `~/.local/bin` is in your PATH

**Slow search:**

Install `fd` for faster parallel searching:

```bash
# Debian/Ubuntu
sudo apt install fd-find

# Fedora
sudo dnf install fd-find

# Arch
sudo pacman -S fd

# macOS
brew install fd
```

### Windows

**PowerShell profile not found:**

Create it manually:

```powershell
New-Item -Path $PROFILE -ItemType File -Force
```

**Slow search:**

Install `fd` via package manager:

```powershell
winget install sharkdp.fd
```

## Support

For bug reports, feature requests, or questions:

- GitHub Issues: [ReggieAlbiosA/fcf/issues](https://github.com/ReggieAlbiosA/fcf/issues)
- GitHub Discussions: [ReggieAlbiosA/fcf/discussions](https://github.com/ReggieAlbiosA/fcf/discussions)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Project Structure

FCF uses a clean, modular Go package structure:

```
fcf/
├── cmd/
│   └── fcf/
│       └── main.go           # Entry point
│
├── internal/
│   ├── command/              # CLI commands, flags, help
│   │   ├── root.go
│   │   ├── help.go
│   │   └── version.go
│   │
│   ├── search/               # File/folder search logic
│   │   └── search.go
│   │
│   ├── ui/                   # Display and interactive mode
│   │   ├── display.go
│   │   └── interactive.go
│   │
│   ├── navigation/           # Directory navigation
│   │   └── navigate.go
│   │
│   ├── install/              # Installation command
│   │   ├── install.go
│   │   ├── install_unix.go   # Unix-specific (build tag: unix)
│   │   ├── install_windows.go # Windows-specific (build tag: windows)
│   │   └── shell/            # Shell integration
│   │       ├── shell.go
│   │       ├── shell_unix.go
│   │       └── shell_windows.go
│   │
│   └── platform/             # Platform-specific utilities
│       ├── exec_unix.go      # Executable detection (Unix)
│       ├── exec_windows.go   # Executable detection (Windows)
│       ├── distro_unix.go    # Linux distro detection
│       └── distro_windows.go # Windows stub
│
├── go.mod
├── go.sum
└── README.md
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/ReggieAlbiosA/fcf.git
cd fcf

# Build the binary
go build -o fcf ./cmd/fcf

# Run it
./fcf --help
```

### Development

- **Go version:** 1.21+
- **Dependencies:** Listed in `go.mod`
- **Platform-specific code:** Uses Go build tags (`//go:build unix`, `//go:build windows`)
- **Testing:** `go test ./...`

## License

MIT License - see LICENSE file for details

## Author

**Reggie Albios**
- GitHub: [@ReggieAlbiosA](https://github.com/ReggieAlbiosA)

## Changelog

### v3.2.0 (2026-01-16)
- **Major:** Refactored to modular Go package structure
- Organized code into clean packages: `cmd/`, `internal/command/`, `internal/ui/`, `internal/search/`, `internal/navigation/`, `internal/install/`, `internal/platform/`
- Improved code maintainability and testability
- Updated CI/CD workflow for new build structure
- Enhanced developer experience with clear project organization

### v3.1.0 (2026-01-16)
- **Major:** Automatic shell integration
- `fcf install` now auto-detects and configures shell (bash, zsh, fish, PowerShell)
- Added `--user`, `--shell`, and `--no-shell` installation flags
- Idempotent installation with marker-based detection
- No manual shell configuration required

### v3.0.0 (2026-01-14)
- **Major:** Unified cross-platform Go codebase
- Consolidated all code to single Go project (Linux, macOS, Windows)
- Platform-specific logic via Go build tags
- Multi-platform CI/CD builds (5 platforms: Linux amd64/arm64, macOS Intel/Apple Silicon, Windows)
- Linux distro detection via `/etc/os-release`
- Automatic `fdfind` → `fd` alias support for Debian/Ubuntu
- Simplified installation: Binary-only approach, no shell script installers needed
- macOS support for both Intel and Apple Silicon

### v2.0.0 (2024-12-30)
- **Windows:** Rewritten in Go (replaces PowerShell script)
- Fixed ANSI color rendering issues on Windows
- Resolved PowerShell execution policy problems
- GitHub Actions automated binary builds
- Improved Windows terminal compatibility

### v1.0.0 (2024-12-28)
- Initial release
- Interactive 3-step workflow
- Parallel search with fd
- Real-time streaming results
- Navigation to results
- Loop workflow (find again, repeat, exit)
- Color-coded output with icons
- Shell-based implementation (Bash for Linux, PowerShell for Windows)

---

**Happy searching!**
