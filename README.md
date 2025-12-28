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

Install with a single command:

```bash
curl -sSL https://raw.githubusercontent.com/ReggieAlbiosA/fcf/refs/heads/main/install.sh | bash
```

The installer automatically installs to **both locations**:
- **User:** `~/.local/bin/fcf` - Available for your user account
- **System:** `/usr/local/bin/fcf` - Available for all users (requires sudo)

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
| `find` | Sequential | Fallback |

The installer will offer to install `fd` automatically. You can also install it manually:

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

### User Installation
- **Location:** `~/.local/bin/fcf`
- **Available to:** Current user only
- **Benefit:** Works without sudo, safe and isolated

### System-Wide Installation
- **Location:** `/usr/local/bin/fcf`
- **Requires:** sudo (installer will prompt)
- **Available to:** All users on the system

## Updating

Re-run the installation command to update:

```bash
curl -sSL https://raw.githubusercontent.com/ReggieAlbiosA/fcf/refs/heads/main/install.sh | bash
```

The installer will detect existing installation and upgrade automatically.

## Manual Installation

```bash
# Download
curl -sSL https://raw.githubusercontent.com/ReggieAlbiosA/fcf/refs/heads/main/fcf.sh -o fcf

# Make executable
chmod +x fcf

# Move to PATH
mv fcf ~/.local/bin/  # or /usr/local/bin with sudo

# Add to PATH if needed
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

## Uninstallation

```bash
# User installation
rm ~/.local/bin/fcf

# System-wide installation
sudo rm /usr/local/bin/fcf

# Remove logs
rm -rf ~/.fcf
```

## Troubleshooting

### Command not found
```bash
source ~/.bashrc  # Reload shell config
# or restart terminal
```

### Permission denied
```bash
chmod +x ~/.local/bin/fcf
```

### Slow search
Install `fd` for parallel searching:
```bash
sudo apt install fd-find  # Debian/Ubuntu
```

## Installation Logs

All installations are logged to: `~/.fcf/install.log`

```bash
cat ~/.fcf/install.log  # View installation history
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Author

**Reggie Albios**
- GitHub: [@ReggieAlbiosA](https://github.com/ReggieAlbiosA)

## Changelog

### v1.0.0 (2024-12-28)
- Initial release
- Interactive 3-step workflow
- Parallel search with fd
- Real-time streaming results
- Navigation to results
- Loop workflow (find again, repeat, exit)
- Color-coded output with icons
- APT-style installer

---

**Happy searching!**
