# Dependency Graph

This document outlines all dependencies for the FCF project across different platforms and components.

## Overview

FCF has different dependencies depending on the platform:
- **Linux/Ubuntu**: Bash script with optional external tools
- **Windows**: Go binary with Go module dependencies

---

## Linux/Ubuntu Dependencies

### Runtime Dependencies

#### Required
- **Bash** (4.0+) - Shell interpreter
- **find** - Standard Unix file search utility (fallback search method)
- **Standard Unix utilities**: `command`, `which`, `awk`, `cat`, `chmod`, `mkdir`, `rm`

#### Optional (Recommended)
- **fd** or **fdfind** - Fast parallel file search tool
  - Provides 5-10x faster search performance
  - Falls back to `find` if not available
  - Installation:
    - Debian/Ubuntu: `sudo apt install fd-find`
    - Fedora: `sudo dnf install fd-find`
    - Arch Linux: `sudo pacman -S fd`
    - macOS: `brew install fd`

### Installation Script Dependencies

The installer (`ubuntu/install.sh`) requires:

- **curl** or **wget** - For downloading the script
- **Package manager** (for optional fd installation):
  - `apt` (Debian/Ubuntu)
  - `dnf` (Fedora)
  - `yum` (RHEL/CentOS)
  - `pacman` (Arch Linux)
  - `brew` (macOS)
  - `zypper` (openSUSE)

### Shell Integration

For navigation to work, users need to add a shell function to:
- `~/.bashrc` (Bash)
- `~/.zshrc` (Zsh)

---

## Windows Dependencies

### Build Dependencies

- **Go 1.21+** - Required to build the Windows binary
  - Download: https://go.dev/dl/
  - Or install via: `winget install GoLang.Go`

### Runtime Dependencies

The Windows version is a standalone binary (`fcf.exe`) with no external runtime dependencies. All functionality is self-contained.

### Go Module Dependencies

Located in `win/go/go.mod`:

#### Direct Dependencies
```
github.com/fatih/color v1.16.0
  └── Terminal color output support

github.com/mattn/go-isatty v0.0.20
  └── Terminal detection (TTY checks)
```

#### Indirect Dependencies
```
github.com/mattn/go-colorable v0.1.13
  └── Windows color support (dependency of fatih/color)

golang.org/x/sys v0.14.0
  └── System-level utilities (dependency of go-isatty)
```

### Installation Script Dependencies

The Windows installer (`win/install.ps1`) requires:

- **PowerShell 5.1+** - For running the installer
- **Internet connection** - To download the binary from GitHub Releases
- **Administrator privileges** - For system-wide installation (optional)

---

## Dependency Graph Visualization

### Linux/Ubuntu Runtime

```
fcf.sh
├── Required
│   ├── Bash (4.0+)
│   ├── find
│   └── Standard Unix utilities
│
└── Optional (Performance)
    └── fd / fdfind
        └── Provides parallel search (5-10x faster)
```

### Windows Runtime

```
fcf.exe (Go binary)
├── Go 1.21+ (build-time only)
│
└── Runtime (bundled in binary)
    ├── github.com/fatih/color
    │   └── github.com/mattn/go-colorable
    │
    └── github.com/mattn/go-isatty
        └── golang.org/x/sys
```

---

## Dependency Management

### Linux/Ubuntu

- **No package manager** - FCF is a standalone script
- **Optional fd installation** - Handled by installer or user
- **System dependencies** - Usually pre-installed on Linux systems

### Windows

- **Go modules** - Managed via `go.mod` and `go.sum`
- **Vendor directory** - Not used (relies on Go module cache)
- **Binary distribution** - Pre-built binaries available on GitHub Releases

---

## Version Compatibility

### Linux/Ubuntu

| Component | Minimum Version | Notes |
|-----------|----------------|-------|
| Bash | 4.0 | Most modern systems have 4.4+ |
| find | Any | Standard on all Unix-like systems |
| fd | Any | Optional, latest recommended |

### Windows

| Component | Minimum Version | Notes |
|-----------|----------------|-------|
| Go | 1.21 | Required for building |
| PowerShell | 5.1 | Required for installer |
| Windows | 10+ | Tested on Windows 10/11 |

---

## Security Considerations

### External Tools

- **fd/fdfind**: Third-party tool, maintained by [sharkdp](https://github.com/sharkdp/fd)
- **Go modules**: All dependencies are from trusted sources
  - `github.com/fatih/color` - Well-maintained, widely used
  - `github.com/mattn/go-isatty` - Standard library alternative
  - `golang.org/x/sys` - Official Go extended standard library

### Dependency Updates

- **Go modules**: Update with `go get -u ./...` and `go mod tidy`
- **fd**: Users should update via their package manager
- **Bash/find**: System-level updates via OS package manager

---

## Troubleshooting Dependencies

### Linux/Ubuntu

**Issue**: Slow search performance
- **Solution**: Install `fd` for parallel search

**Issue**: Installer fails to download
- **Solution**: Ensure `curl` or `wget` is installed

**Issue**: Navigation doesn't work
- **Solution**: Add shell function to `~/.bashrc` or `~/.zshrc`

### Windows

**Issue**: Can't build from source
- **Solution**: Install Go 1.21+ and ensure it's in PATH

**Issue**: Binary won't run
- **Solution**: Download pre-built binary from GitHub Releases

**Issue**: Colors don't work
- **Solution**: Use Windows Terminal or PowerShell 7+ for better color support

---

## License Compatibility

All dependencies are compatible with FCF's MIT License:

- **fd**: Apache-2.0 or MIT
- **Go modules**: MIT or compatible licenses
- **Standard Unix tools**: Various (GPL, BSD, etc.) - system-level, not bundled

---

## Dependency Audit

To audit dependencies:

### Linux/Ubuntu
- Check installed versions: `fd --version`, `bash --version`
- Review script for external command usage

### Windows
```powershell
cd win/go
go list -m all          # List all dependencies
go mod verify           # Verify module checksums
go list -json -m all    # Detailed dependency info
```

---

_Last updated: 2026-01-07_

