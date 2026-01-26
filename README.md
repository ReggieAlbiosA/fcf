# FCF (Find File or Folder)

A fast, interactive file and folder finder with parallel search and real-time streaming results.

## What is FCF?

FCF is an interactive command-line tool for finding files and folders with advanced pattern matching, real-time streaming results, and easy navigation. It uses parallel search for speed and provides a seamless workflow for locating and navigating to files.

Think of it as `find` with a friendly interface - interactive prompts, color-coded output, and instant navigation.

## Key Features

- **Parallel Search** - Uses `fd` for blazing fast parallel file searching
- **Real-Time Streaming** - Results appear line-by-line as they're found
- **Stoppable Search** - Press `s` to stop search immediately and use current results
- **Interactive Mode** - Step-by-step guided search workflow
- **Navigation** - Jump directly to any result or path
- **Pattern Matching** - Glob patterns, partial names, extensions
- **Color-Coded Output** - Visual distinction for folders, files, executables, symlinks
- **Loop Workflow** - Search again without restarting
- **Self-Update** - Built-in `fcf update` command to get latest version

## Quick Start

### Installation

One-liner installation - downloads, installs, and configures everything automatically.

#### Linux

```bash
# amd64 (most common)
curl -fsSL https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-linux-amd64 -o /tmp/fcf && chmod +x /tmp/fcf && sudo /tmp/fcf install && rm /tmp/fcf

# arm64
curl -fsSL https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-linux-arm64 -o /tmp/fcf && chmod +x /tmp/fcf && sudo /tmp/fcf install && rm /tmp/fcf
```

#### macOS

```bash
# Apple Silicon (M1/M2/M3)
curl -fsSL https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-darwin-arm64 -o /tmp/fcf && chmod +x /tmp/fcf && sudo /tmp/fcf install && rm /tmp/fcf

# Intel Mac
curl -fsSL https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-darwin-amd64 -o /tmp/fcf && chmod +x /tmp/fcf && sudo /tmp/fcf install && rm /tmp/fcf
```

#### Windows (PowerShell as Administrator)

```powershell
irm https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-windows-amd64.exe -OutFile $env:TEMP\fcf.exe; & $env:TEMP\fcf.exe install; Remove-Item $env:TEMP\fcf.exe
```

**After installation, reload your shell or open a new terminal to start using `fcf`.**

The installer automatically:
- Copies the binary to system path (`/usr/local/bin` or `C:\Program Files\fcf`)
- Detects your shell (Bash, Zsh, Fish, or PowerShell)
- Adds the shell wrapper function for directory navigation
- Works immediately after shell reload

#### Package Managers (Alternative)

**Homebrew (macOS/Linux):**
```bash
brew tap ReggieAlbiosA/tap
brew install fcf
fcf install --shell-only  # Add shell integration
```

**Scoop (Windows):**
```powershell
scoop bucket add fcf https://github.com/ReggieAlbiosA/scoop-bucket
scoop install fcf
fcf install --shell-only  # Add shell integration
```

### Managing FCF

```bash
# Update to latest version
sudo fcf update

# Uninstall fcf
sudo fcf uninstall
```

On Windows, run these commands in PowerShell as Administrator (without `sudo`).

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
- Linux/macOS: `/usr/local/bin/fcf`
- Windows: `C:\Program Files\fcf\fcf.exe`

**Shell Configuration:**
- Automatically detected: Bash, Zsh, Fish (Unix/Linux/macOS), PowerShell (Windows)
- Wrapper function added to config file for directory navigation
- Idempotent: Won't duplicate on reinstall

### Command Options

```bash
fcf install               # Full installation (requires sudo on Unix)
fcf install --no-shell    # Binary only, skip shell integration
fcf install --shell zsh   # Force specific shell configuration
fcf install --shell-only  # Only install shell integration (skip binary)
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

Re-run the same one-liner installation command to update to the latest version:

```bash
# Linux (amd64)
curl -fsSL https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-linux-amd64 -o /tmp/fcf && chmod +x /tmp/fcf && sudo /tmp/fcf install && rm /tmp/fcf
```

```powershell
# Windows (PowerShell as Administrator)
irm https://github.com/ReggieAlbiosA/fcf/releases/latest/download/fcf-windows-amd64.exe -OutFile $env:TEMP\fcf.exe; & $env:TEMP\fcf.exe install; Remove-Item $env:TEMP\fcf.exe
```

The installer detects your existing installation and updates automatically.

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
# Remove binary
sudo rm /usr/local/bin/fcf

# Optional: Remove shell integration from config files
# Edit ~/.bashrc, ~/.zshrc, or ~/.config/fish/config.fish
# and remove the fcf function block
```

### Windows (PowerShell as Administrator)

```powershell
# Remove binary
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

### v3.3.1 (2026-01-26)
- **Fix:** Shell integration now installs to correct user's config when using `sudo`
- Previously, running `sudo ./fcf install` wrote shell wrapper to `/root/.bashrc` instead of the user's config
- Now correctly detects `SUDO_USER` environment variable to find the real user's home directory
- Navigation feature now works properly after one-liner installation

### v3.3.0 (2026-01-23)
- **Major:** Package management and search control
- Added `fcf update` command for self-updating from GitHub releases
- Added `fcf uninstall` command for clean removal (binary + shell integration)
- Added press `s` to stop search immediately during execution
- Cross-platform keyboard input support for stop functionality
- One-liner installation commands in README
- Updated help command with new subcommands

### v3.2.0 (2026-01-16)
- **Major:** Refactored to modular Go package structure
- Organized code into clean packages: `cmd/`, `internal/command/`, `internal/ui/`, `internal/search/`, `internal/navigation/`, `internal/install/`, `internal/platform/`
- Improved code maintainability and testability
- Updated CI/CD workflow for new build structure
- Enhanced developer experience with clear project organization

### v3.1.0 (2026-01-16)
- **Major:** Automatic shell integration
- `fcf install` now auto-detects and configures shell (bash, zsh, fish, PowerShell)
- Added `--shell`, `--no-shell`, and `--shell-only` installation flags
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
