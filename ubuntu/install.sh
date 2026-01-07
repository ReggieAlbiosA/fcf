#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m' # No Color

# GitHub repository details
GITHUB_USER="ReggieAlbiosA"
GITHUB_REPO="fcf"
GITHUB_BRANCH="main"
SCRIPT_URL="https://raw.githubusercontent.com/${GITHUB_USER}/${GITHUB_REPO}/refs/heads/${GITHUB_BRANCH}/ubuntu/fcf.sh"

# Check if fd is installed (for fast parallel searching)
has_fd() {
    command -v fd &> /dev/null || command -v fdfind &> /dev/null
}

# Detect package manager
detect_package_manager() {
    if command -v apt &> /dev/null; then
        echo "apt"
    elif command -v dnf &> /dev/null; then
        echo "dnf"
    elif command -v yum &> /dev/null; then
        echo "yum"
    elif command -v pacman &> /dev/null; then
        echo "pacman"
    elif command -v brew &> /dev/null; then
        echo "brew"
    elif command -v zypper &> /dev/null; then
        echo "zypper"
    else
        echo "unknown"
    fi
}

# Install fd based on package manager
install_fd() {
    local pkg_manager=$(detect_package_manager)

    print_status "step" "Installing fd (fast file finder)..."
    log "Installing fd using $pkg_manager"

    case $pkg_manager in
        apt)
            # fd is named fd-find on Debian/Ubuntu
            if sudo apt install -y fd-find 2>/dev/null; then
                # Create symlink if needed (fd-find installs as 'fdfind')
                if command -v fdfind &> /dev/null && ! command -v fd &> /dev/null; then
                    sudo ln -sf $(which fdfind) /usr/local/bin/fd 2>/dev/null || true
                fi
                return 0
            fi
            ;;
        dnf|yum)
            sudo $pkg_manager install -y fd-find 2>/dev/null && return 0
            ;;
        pacman)
            sudo pacman -S --noconfirm fd 2>/dev/null && return 0
            ;;
        brew)
            brew install fd 2>/dev/null && return 0
            ;;
        zypper)
            sudo zypper install -y fd 2>/dev/null && return 0
            ;;
    esac

    return 1
}

# Log file location
LOG_DIR="$HOME/.fcf"
LOG_FILE="$LOG_DIR/install.log"

# Create log directory if it doesn't exist
mkdir -p "$LOG_DIR" 2>/dev/null || true

# Logging function
log() {
    local message="$1"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[$timestamp] $message" >> "$LOG_FILE"
}

# Status printing functions (apt-style)
print_status() {
    local status="$1"
    local message="$2"
    case $status in
        "ok")
            echo -e "${GREEN}[${BOLD}✓${NC}${GREEN}]${NC} $message"
            log "[OK] $message"
            ;;
        "info")
            echo -e "${CYAN}[${BOLD}*${NC}${CYAN}]${NC} $message"
            log "[INFO] $message"
            ;;
        "warn")
            echo -e "${YELLOW}[${BOLD}!${NC}${YELLOW}]${NC} $message"
            log "[WARN] $message"
            ;;
        "error")
            echo -e "${RED}[${BOLD}✗${NC}${RED}]${NC} $message"
            log "[ERROR] $message"
            ;;
        "step")
            echo -e "${BLUE}${BOLD}==>${NC} $message"
            log "[STEP] $message"
            ;;
    esac
}

# Progress indicator (apt-style)
progress() {
    local message="$1"
    echo -ne "${message}..."
    log "$message"
}

progress_done() {
    echo -e " ${GREEN}Done${NC}"
    log "Done"
}

print_step() {
    echo -e "${BOLD}${BLUE}▸ $1${NC}"
}

# Header
log "========================================="
log "FCF Installation Started"
log "========================================="

echo -e "${BOLD}${CYAN}╔════════════════════════════════════════╗${NC}"
echo -e "${BOLD}${CYAN}║${NC}   ${BOLD}FCF Installer${NC}                       ${BOLD}${CYAN}║${NC}"
echo -e "${BOLD}${CYAN}║${NC}   ${DIM}Find File or Folder${NC}                 ${BOLD}${CYAN}║${NC}"
echo -e "${BOLD}${CYAN}╚════════════════════════════════════════╝${NC}"
echo ""

# Function to install to user bin
install_user() {
    BIN_DIR="$HOME/.local/bin"
    local IS_UPDATE=false

    print_step "User Installation"
    echo ""

    # Check if already installed
    if [ -f "$BIN_DIR/fcf" ]; then
        IS_UPDATE=true
        print_status "info" "Package 'fcf' is already installed"
        print_status "info" "Preparing to upgrade fcf..."
        log "Action: UPDATE"
    else
        print_status "info" "Preparing to install fcf..."
        log "Action: FRESH INSTALL"
    fi

    log "Installation Location: $BIN_DIR"

    # Create directory
    progress "Creating directories"
    mkdir -p "$BIN_DIR"
    progress_done
    log "Created/Verified directory: $BIN_DIR"

    # Download
    print_status "step" "Fetching fcf from repository..."
    progress "  Downloading $SCRIPT_URL"
    log "Downloading from: $SCRIPT_URL"

    if command -v curl &> /dev/null; then
        if curl -sSL "$SCRIPT_URL" -o "$BIN_DIR/fcf" 2>/dev/null; then
            progress_done
            log "Download method: curl"
        else
            echo -e " ${RED}Failed${NC}"
            print_status "error" "Download failed"
            log "ERROR: Download failed (curl)"
            exit 1
        fi
    elif command -v wget &> /dev/null; then
        if wget -q "$SCRIPT_URL" -O "$BIN_DIR/fcf" 2>/dev/null; then
            progress_done
            log "Download method: wget"
        else
            echo -e " ${RED}Failed${NC}"
            print_status "error" "Download failed"
            log "ERROR: Download failed (wget)"
            exit 1
        fi
    else
        echo ""
        print_status "error" "Neither curl nor wget found. Please install one of them."
        log "ERROR: No download tool available (curl/wget)"
        exit 1
    fi

    # Set permissions
    progress "Setting up fcf"
    chmod +x "$BIN_DIR/fcf"
    progress_done
    log "Made executable: $BIN_DIR/fcf"

    # Configure PATH
    if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
        progress "Configuring environment (PATH)"

        if ! grep -q "export PATH=\"\$HOME/.local/bin:\$PATH\"" "$HOME/.bashrc" 2>/dev/null; then
            echo '' >> "$HOME/.bashrc"
            echo '# Added by fcf installer' >> "$HOME/.bashrc"
            echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
            progress_done
            log "Added PATH to ~/.bashrc"
        else
            progress_done
            log "PATH already in ~/.bashrc (skipped)"
        fi
    else
        log "$BIN_DIR already in PATH"
    fi

    # Add shell wrapper function for navigation
    add_shell_function() {
        local PROFILE_FILE="$1"

        if grep -q "function fcf\|fcf()" "$PROFILE_FILE" 2>/dev/null; then
            print_status "info" "fcf shell function already exists in $(basename $PROFILE_FILE)"
            return 0
        fi

        cat >> "$PROFILE_FILE" << 'WRAPPER'

# Added by fcf installer - enables directory navigation
fcf() {
    command fcf "$@"
    local nav_path="$HOME/.fcf_nav_path"
    if [[ -f "$nav_path" ]]; then
        cd "$(cat "$nav_path")"
        rm -f "$nav_path"
    fi
}
WRAPPER
        print_status "ok" "Added fcf shell function to $(basename $PROFILE_FILE)"
    }

    # Add to user's shell profile
    if [ -f "$HOME/.bashrc" ]; then
        add_shell_function "$HOME/.bashrc"
    fi
    if [ -f "$HOME/.zshrc" ]; then
        add_shell_function "$HOME/.zshrc"
    fi

    echo ""
    if [ "$IS_UPDATE" = true ]; then
        print_status "ok" "fcf upgraded successfully"
        log "User upgrade completed successfully"
    else
        print_status "ok" "fcf installed successfully"
        log "User installation completed successfully"
    fi

    print_status "info" "Location: $BIN_DIR/fcf"
    echo ""
    print_status "warn" "Please run: ${BOLD}source ~/.bashrc${NC}"
    print_status "warn" "Or restart your terminal to use: ${BOLD}${GREEN}fcf${NC}"
}

# Function to install system-wide
install_system() {
    BIN_DIR="/usr/local/bin"
    local IS_UPDATE=false

    print_step "System-Wide Installation"
    echo ""

    # Check if already installed
    if [ -f "$BIN_DIR/fcf" ]; then
        IS_UPDATE=true
        print_status "info" "Package 'fcf' is already installed (system-wide)"
        print_status "info" "Preparing to upgrade fcf..."
        log "Action: UPDATE (SYSTEM-WIDE)"
    else
        print_status "info" "Preparing to install fcf (system-wide)..."
        log "Action: FRESH INSTALL (SYSTEM-WIDE)"
    fi

    log "Installation Location: $BIN_DIR"

    # Download to temp
    TMP_FILE=$(mktemp)
    print_status "step" "Fetching fcf from repository..."
    progress "  Downloading $SCRIPT_URL"
    log "Downloading from: $SCRIPT_URL"
    log "Temporary file: $TMP_FILE"

    if command -v curl &> /dev/null; then
        if curl -sSL "$SCRIPT_URL" -o "$TMP_FILE" 2>/dev/null; then
            progress_done
            log "Download method: curl"
        else
            echo -e " ${RED}Failed${NC}"
            print_status "error" "Download failed"
            log "ERROR: Download failed (curl)"
            exit 1
        fi
    elif command -v wget &> /dev/null; then
        if wget -q "$SCRIPT_URL" -O "$TMP_FILE" 2>/dev/null; then
            progress_done
            log "Download method: wget"
        else
            echo -e " ${RED}Failed${NC}"
            print_status "error" "Download failed"
            log "ERROR: Download failed (wget)"
            exit 1
        fi
    else
        echo ""
        print_status "error" "Neither curl nor wget found. Please install one of them."
        log "ERROR: No download tool available (curl/wget)"
        exit 1
    fi

    # Install to system directory (already running as root)
    progress "Setting up fcf"
    mv "$TMP_FILE" "$BIN_DIR/fcf"
    chmod +x "$BIN_DIR/fcf"
    progress_done
    log "Moved to system bin"
    log "Made executable: $BIN_DIR/fcf"

    echo ""
    if [ "$IS_UPDATE" = true ]; then
        print_status "ok" "fcf upgraded successfully (system-wide)"
        log "System-wide upgrade completed successfully"
    else
        print_status "ok" "fcf installed successfully (system-wide)"
        log "System-wide installation completed successfully"
    fi

    print_status "info" "Location: $BIN_DIR/fcf"
    print_status "info" "Available to all users"
}

# Detect if running as root/sudo
if [[ $EUID -eq 0 ]]; then
    # Running as root - system-wide installation
    print_step "System-Wide Installation (running as root)"
    echo ""
    log "Installing to system directory (running as root)"
    install_system
else
    # Running as regular user - user installation only
    print_step "User Installation"
    echo ""
    log "Installing to user directory"
    install_user
fi

# Check for fd (optional fast search dependency)
echo ""
print_step "Checking Optional Dependencies"
echo ""

if has_fd; then
    print_status "ok" "fd is installed (fast parallel search enabled)"
    log "fd: already installed"
else
    print_status "warn" "fd not found (fcf will use slower 'find' command)"
    log "fd: not installed"

    # Ask user if they want to install fd
    echo ""
    echo -e "  ${CYAN}fd${NC} enables ${BOLD}5-10x faster${NC} parallel file searching."
    echo -ne "  Install fd now? [Y/n]: "

    # Handle piped input (curl | bash)
    if [ -t 0 ]; then
        read -r install_fd_choice </dev/tty
    else
        read -r install_fd_choice
    fi

    if [[ "$install_fd_choice" =~ ^[Yy]?$ ]]; then
        if install_fd; then
            # Check again after installation (handle fdfind -> fd symlink)
            if has_fd || command -v fdfind &> /dev/null; then
                print_status "ok" "fd installed successfully"
                log "fd: installed successfully"
            else
                print_status "warn" "fd installation may require restart"
                log "fd: installation completed, may need restart"
            fi
        else
            print_status "warn" "Could not install fd automatically"
            print_status "info" "Install manually: ${CYAN}sudo apt install fd-find${NC}"
            log "fd: auto-install failed"
        fi
    else
        print_status "info" "Skipped fd installation"
        print_status "info" "Install later: ${CYAN}sudo apt install fd-find${NC}"
        log "fd: user skipped installation"
    fi
fi

# Summary
echo ""
echo -e "${BOLD}${CYAN}╔════════════════════════════════════════╗${NC}"
echo -e "${BOLD}${CYAN}║${NC}   ${GREEN}${BOLD}Installation Complete!${NC}             ${BOLD}${CYAN}║${NC}"
echo -e "${BOLD}${CYAN}╚════════════════════════════════════════╝${NC}"
echo ""
print_status "info" "Usage: ${BOLD}${GREEN}fcf${NC} (interactive mode)"
print_status "info" "Usage: ${BOLD}${GREEN}fcf \"*.log\"${NC} (direct search)"
echo ""

# Show feature status
echo -e "${BOLD}Features:${NC}"
if has_fd || command -v fdfind &> /dev/null; then
    echo -e "  ${GREEN}✓${NC} Fast parallel search (fd)"
else
    echo -e "  ${YELLOW}○${NC} Fast parallel search (install fd for boost)"
fi
echo -e "  ${GREEN}✓${NC} Real-time streaming results"
echo -e "  ${GREEN}✓${NC} Interactive navigation"
echo -e "  ${GREEN}✓${NC} Pattern matching (glob, regex)"
echo ""

print_status "info" "Installation log: ${CYAN}$LOG_FILE${NC}"
echo ""

log "========================================="
log "FCF Installation Finished Successfully"
log "========================================="
