#!/bin/bash

# fcf - Find File or Folder Command
# Interactive file/folder finder with parallel search and real-time streaming

VERSION="1.0.0"

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color
BOLD='\033[1m'
DIM='\033[2m'

# Default settings
IGNORE_CASE=false
TYPE_FILTER=""
SHOW_SIZE=false
MAX_DISPLAY=0  # 0 = unlimited
HIDDEN_FILES=false

# Store results for navigation
declare -a RESULTS=()

# Function to show help
show_help() {
    cat << EOF
${BOLD}fcf - Find File or Folder${NC} v${VERSION}

${BOLD}USAGE:${NC}
    fcf [OPTIONS] [PATTERN] [PATH]
    fcf                          # Interactive mode

${BOLD}DESCRIPTION:${NC}
    Interactive tool to find files and folders with pattern matching
    and real-time streaming results. Uses parallel search for speed.

${BOLD}OPTIONS:${NC}
    ${CYAN}-h, --help${NC}              Show this help message
    ${CYAN}-i, --ignore-case${NC}       Case-insensitive pattern matching
    ${CYAN}-t, --type TYPE${NC}         Filter by type: ${YELLOW}f${NC}(file) or ${YELLOW}d${NC}(directory)
    ${CYAN}-H, --hidden${NC}            Include hidden files/folders
    ${CYAN}--show-size${NC}             Display file sizes
    ${CYAN}--max-display NUM${NC}       Maximum results to display (default: unlimited)

${BOLD}EXAMPLES:${NC}
    ${GREEN}# Interactive mode${NC}
    fcf

    ${GREEN}# Find all .log files${NC}
    fcf "*.log"

    ${GREEN}# Find in specific directory${NC}
    fcf "*.js" ~/projects

    ${GREEN}# Case-insensitive search for PNG files${NC}
    fcf -i "*.PNG"

    ${GREEN}# Find only directories named 'src'${NC}
    fcf -t d src

    ${GREEN}# Find with file sizes${NC}
    fcf --show-size "*.mp4"

${BOLD}INTERACTIVE WORKFLOW:${NC}
    Step 1: Enter path to search
    Step 2: Enter pattern to find
    Step 3: Navigate to a result path

${BOLD}NAVIGATION OPTIONS:${NC}
    After navigation, choose:
    ${CYAN}f${NC} - Find again (restart from Step 1)
    ${CYAN}r${NC} - Repeat search (go to Step 2, same path)
    ${CYAN}n${NC} - Exit

${BOLD}PERFORMANCE:${NC}
    - Uses 'fd' for fast parallel searching (if installed)
    - Falls back to 'find' if fd is not available
    - Install fd: ${CYAN}sudo apt install fd-find${NC} (Debian/Ubuntu)
                  ${CYAN}brew install fd${NC} (macOS)

${BOLD}EXIT CODES:${NC}
    0 - Success
    1 - No matches found
    2 - User cancelled

EOF
}

# Function to check if fd is available
has_fd() {
    command -v fd &> /dev/null || command -v fdfind &> /dev/null
}

# Get the fd command name (fd or fdfind)
get_fd_cmd() {
    if command -v fd &> /dev/null; then
        echo "fd"
    elif command -v fdfind &> /dev/null; then
        echo "fdfind"
    fi
}

# Function to format file size
format_size() {
    local bytes=$1
    if (( bytes >= 1073741824 )); then
        echo "$(awk "BEGIN {printf \"%.1f\", $bytes/1073741824}")G"
    elif (( bytes >= 1048576 )); then
        echo "$(awk "BEGIN {printf \"%.1f\", $bytes/1048576}")M"
    elif (( bytes >= 1024 )); then
        echo "$(awk "BEGIN {printf \"%.1f\", $bytes/1024}")K"
    else
        echo "${bytes}B"
    fi
}

# Function to get file info
get_file_info() {
    local file=$1
    local size=""

    if [[ "$SHOW_SIZE" == true ]] && [[ -f "$file" ]]; then
        local bytes=$(stat -c%s "$file" 2>/dev/null || stat -f%z "$file" 2>/dev/null)
        if [[ -n "$bytes" ]]; then
            size=" ${DIM}($(format_size $bytes))${NC}"
        fi
    fi

    echo "$size"
}

# Function to display a found item with styling
display_result() {
    local file=$1
    local count=$2
    local info=$(get_file_info "$file")

    if [[ -d "$file" ]]; then
        # Directory - blue with folder icon
        echo -e "${CYAN}  [$count]${NC} ${BLUE}📁 $file/${NC}$info"
    elif [[ -x "$file" ]]; then
        # Executable - green
        echo -e "${CYAN}  [$count]${NC} ${GREEN}⚡ $file${NC}$info"
    elif [[ -L "$file" ]]; then
        # Symlink - magenta
        echo -e "${CYAN}  [$count]${NC} ${MAGENTA}🔗 $file${NC}$info"
    else
        # Regular file
        echo -e "${CYAN}  [$count]${NC} 📄 $file$info"
    fi
}

# Build fd command for parallel search
build_fd_command() {
    local pattern=$1
    local search_path=$2
    local fd_cmd=$(get_fd_cmd)

    # Base options
    local opts="--color never"

    # Include hidden files if requested
    if [[ "$HIDDEN_FILES" == true ]]; then
        opts+=" --hidden"
    fi

    # Don't respect .gitignore for complete results
    opts+=" --no-ignore"

    # Type filter
    if [[ -n "$TYPE_FILTER" ]]; then
        opts+=" -t $TYPE_FILTER"
    fi

    # Case sensitivity
    if [[ "$IGNORE_CASE" == true ]]; then
        opts+=" -i"
    else
        opts+=" -s"
    fi

    # Build command with glob pattern
    echo "$fd_cmd $opts -g \"$pattern\" \"$search_path\""
}

# Build find command as fallback
build_find_command() {
    local pattern=$1
    local search_path=$2
    local find_cmd="find \"$search_path\""

    # Type filter
    if [[ -n "$TYPE_FILTER" ]]; then
        find_cmd+=" -type $TYPE_FILTER"
    fi

    # Pattern matching
    if [[ "$IGNORE_CASE" == true ]]; then
        find_cmd+=" -iname \"$pattern\""
    else
        find_cmd+=" -name \"$pattern\""
    fi

    # Exclude hidden if not requested
    if [[ "$HIDDEN_FILES" == false ]]; then
        find_cmd+=" -not -path '*/.*'"
    fi

    echo "$find_cmd"
}

# Function to show header
show_header() {
    clear
    echo -e "${BOLD}${CYAN}╔════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${CYAN}║${NC}   ${BOLD}fcf${NC} - Find File or Folder          ${BOLD}${CYAN}║${NC}"
    echo -e "${BOLD}${CYAN}╚════════════════════════════════════════╝${NC}"
    echo ""
}

# Function to navigate to path
navigate_to_path() {
    local target_path=$1

    # If it's a file, get the directory
    if [[ -f "$target_path" ]]; then
        target_path=$(dirname "$target_path")
    fi

    # Check if path exists
    if [[ ! -d "$target_path" ]]; then
        echo -e "${RED}ERROR:${NC} Directory '$target_path' does not exist"
        return 1
    fi

    # Change to the directory
    cd "$target_path" 2>/dev/null
    if [[ $? -eq 0 ]]; then
        echo -e "${GREEN}✓ Navigated to:${NC} ${CYAN}$target_path${NC}"
        echo ""
        echo -e "${DIM}Contents:${NC}"
        ls -la --color=auto 2>/dev/null || ls -la
        return 0
    else
        echo -e "${RED}ERROR:${NC} Could not navigate to '$target_path'"
        return 1
    fi
}

# Function to show options menu
show_options_menu() {
    echo ""
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BOLD}Options:${NC}"
    echo -e "  ${CYAN}[f]${NC} Find again (new search)"
    echo -e "  ${CYAN}[r]${NC} Repeat search (same path)"
    echo -e "  ${CYAN}[n]${NC} Exit"
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    read -p "$(echo -e ${CYAN}Choose:${NC} )" choice

    case $choice in
        f|F)
            return 1  # Go to Step 1
            ;;
        r|R)
            return 2  # Go to Step 2
            ;;
        n|N|*)
            return 0  # Exit
            ;;
    esac
}

# Parse command line arguments
PATTERN=""
SEARCH_PATH=""
CLI_MODE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -i|--ignore-case)
            IGNORE_CASE=true
            shift
            ;;
        -t|--type)
            TYPE_FILTER="$2"
            if [[ "$TYPE_FILTER" != "f" && "$TYPE_FILTER" != "d" ]]; then
                echo -e "${RED}ERROR:${NC} Invalid type '$TYPE_FILTER'. Use 'f' for files or 'd' for directories"
                exit 1
            fi
            shift 2
            ;;
        -H|--hidden)
            HIDDEN_FILES=true
            shift
            ;;
        --show-size)
            SHOW_SIZE=true
            shift
            ;;
        --max-display)
            MAX_DISPLAY="$2"
            shift 2
            ;;
        -*)
            echo -e "${RED}ERROR:${NC} Unknown option: $1"
            echo "Use 'fcf --help' for usage information"
            exit 1
            ;;
        *)
            if [[ -z "$PATTERN" ]]; then
                PATTERN="$1"
                CLI_MODE=true
            else
                SEARCH_PATH="$1"
            fi
            shift
            ;;
    esac
done

# Main loop
CURRENT_STEP=1

while true; do
    # Reset results array
    RESULTS=()

    # Show header
    show_header

    # Step 1: Get search path
    if [[ $CURRENT_STEP -eq 1 ]]; then
        if [[ -z "$SEARCH_PATH" ]] || [[ "$CLI_MODE" == false ]]; then
            echo -e "${BOLD}Step 1:${NC} Enter path to search"
            echo -e "${DIM}(Press Enter for current directory: $PWD)${NC}"
            echo ""
            read -p "$(echo -e ${CYAN}Path:${NC} )" user_path

            if [[ -z "$user_path" ]]; then
                SEARCH_PATH="."
                echo -e "${GREEN}Using current directory${NC}"
            else
                # Expand tilde to home directory
                SEARCH_PATH="${user_path/#\~/$HOME}"

                if [[ ! -d "$SEARCH_PATH" ]]; then
                    echo -e "${RED}ERROR:${NC} Directory '$SEARCH_PATH' does not exist"
                    read -p "Press Enter to try again..."
                    continue
                fi
            fi
            echo ""
        fi
        CURRENT_STEP=2
    fi

    # Step 2: Get pattern
    if [[ $CURRENT_STEP -eq 2 ]]; then
        if [[ -z "$PATTERN" ]] || [[ "$CLI_MODE" == false ]]; then
            echo -e "${BOLD}Step 2:${NC} Enter file/folder name or pattern to find"
            echo -e "${DIM}Examples: *.log, config, .env, src, *.js${NC}"
            echo ""
            read -p "$(echo -e ${CYAN}Pattern:${NC} )" PATTERN

            if [[ -z "$PATTERN" ]]; then
                echo -e "${RED}ERROR:${NC} Pattern cannot be empty"
                read -p "Press Enter to try again..."
                continue
            fi
            echo ""
        fi
    fi

    # Show search info
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Searching in:${NC} ${CYAN}$SEARCH_PATH${NC}"
    echo -e "${BLUE}Pattern:${NC} ${YELLOW}$PATTERN${NC}"

    # Show search method
    if has_fd; then
        echo -e "${BLUE}Method:${NC} ${GREEN}fd (parallel search)${NC}"
    else
        echo -e "${BLUE}Method:${NC} ${YELLOW}find (sequential - install 'fd' for faster search)${NC}"
    fi
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    # Build search command
    if has_fd; then
        search_cmd=$(build_fd_command "$PATTERN" "$SEARCH_PATH")
    else
        search_cmd=$(build_find_command "$PATTERN" "$SEARCH_PATH")
    fi

    # Execute search with real-time streaming
    echo -e "${BOLD}Results:${NC} ${DIM}(streaming in real-time...)${NC}"
    echo ""

    count=0
    start_time=$(date +%s.%N)

    while IFS= read -r file; do
        # Skip empty lines
        [[ -z "$file" ]] && continue

        ((count++))

        # Store result for navigation
        RESULTS+=("$file")

        # Check max display limit
        if [[ $MAX_DISPLAY -gt 0 ]] && [[ $count -gt $MAX_DISPLAY ]]; then
            continue  # Still count but don't display
        fi

        # Display result immediately (streaming)
        display_result "$file" "$count"

    done < <(eval "$search_cmd" 2>/dev/null)

    end_time=$(date +%s.%N)
    elapsed=$(echo "$end_time - $start_time" | bc 2>/dev/null || echo "0")

    # Summary
    echo ""
    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

    if [[ $count -eq 0 ]]; then
        echo -e "${YELLOW}No matches found${NC} for pattern: ${CYAN}$PATTERN${NC}"
        echo ""
        echo -e "${DIM}Tips:${NC}"
        echo -e "  - Try a different pattern"
        echo -e "  - Use ${CYAN}-i${NC} for case-insensitive search"
        echo -e "  - Use ${CYAN}-H${NC} to include hidden files"

        # Show options even when no results
        show_options_menu
        result=$?

        if [[ $result -eq 0 ]]; then
            exit 0
        elif [[ $result -eq 1 ]]; then
            CURRENT_STEP=1
            SEARCH_PATH=""
            PATTERN=""
            CLI_MODE=false
            continue
        elif [[ $result -eq 2 ]]; then
            CURRENT_STEP=2
            PATTERN=""
            CLI_MODE=false
            continue
        fi
    else
        echo -e "${GREEN}${BOLD}Found $count match(es)${NC} in ${CYAN}${elapsed}s${NC}"

        if [[ $MAX_DISPLAY -gt 0 ]] && [[ $count -gt $MAX_DISPLAY ]]; then
            echo -e "${YELLOW}(Displayed first $MAX_DISPLAY of $count)${NC}"
        fi
    fi

    echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

    # Step 3: Navigate to path
    echo ""
    echo -e "${BOLD}Step 3:${NC} Enter path to navigate to"
    echo -e "${DIM}(Enter a number from results, full path, or press Enter to skip)${NC}"
    echo ""
    read -p "$(echo -e ${CYAN}Navigate to:${NC} )" nav_input

    if [[ -n "$nav_input" ]]; then
        # Check if input is a number (result index)
        if [[ "$nav_input" =~ ^[0-9]+$ ]]; then
            index=$((nav_input - 1))
            if [[ $index -ge 0 ]] && [[ $index -lt ${#RESULTS[@]} ]]; then
                target_path="${RESULTS[$index]}"
            else
                echo -e "${RED}ERROR:${NC} Invalid result number"
                target_path=""
            fi
        else
            # Treat as path (expand tilde)
            target_path="${nav_input/#\~/$HOME}"
        fi

        if [[ -n "$target_path" ]]; then
            echo ""
            navigate_to_path "$target_path"
        fi
    else
        echo -e "${DIM}Skipped navigation${NC}"
    fi

    # Show options menu after Step 3
    show_options_menu
    result=$?

    if [[ $result -eq 0 ]]; then
        echo -e "${GREEN}Goodbye!${NC}"
        exit 0
    elif [[ $result -eq 1 ]]; then
        # Go to Step 1
        CURRENT_STEP=1
        SEARCH_PATH=""
        PATTERN=""
        CLI_MODE=false
    elif [[ $result -eq 2 ]]; then
        # Go to Step 2
        CURRENT_STEP=2
        PATTERN=""
        CLI_MODE=false
    fi
done

exit 0
