#!/bin/bash

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting Conventional Commit Workflow...${NC}"

# Check for changes
if [[ -z $(git status -s) ]]; then
    echo -e "${YELLOW}No changes detected.${NC}"
    exit 0
fi

# 1. Staging
echo -e "\n${BLUE}Current Status:${NC}"
git status -s

read -p "Stage all changes? (y/n): " STAGE_ALL
if [[ "$STAGE_ALL" == "y" ]]; then
    git add .
else
    echo "Please stage your files manually using 'git add' and run this script again."
    exit 0
fi

# 2. Select Type
TYPES=("feat: A new feature" "fix: A bug fix" "docs: Documentation only changes" "style: Changes that do not affect the meaning of the code" "refactor: A code change that neither fixes a bug nor adds a feature" "perf: A code change that improves performance" "test: Adding missing tests or correcting existing tests" "build: Changes that affect the build system or external dependencies" "ci: Changes to our CI configuration files and scripts" "chore: Other changes that don't modify src or test files" "revert: Reverts a previous commit")

CHOICE=$(printf "%s\n" "${TYPES[@]}" | fzf --prompt="Select commit type: " --height=~15)
TYPE=$(echo $CHOICE | cut -d':' -f1)

if [[ -z "$TYPE" ]]; then
    echo -e "${YELLOW}No type selected. Exiting.${NC}"
    exit 1
fi

# 3. Scope
read -p "Enter scope (optional, e.g., 'cli', 'core'): " SCOPE
if [[ -n "$SCOPE" ]]; then
    SCOPE="($SCOPE)"
fi

# 4. Description
read -p "Enter short description: " DESC
if [[ -z "$DESC" ]]; then
    echo -e "${YELLOW}Description is required.${NC}"
    exit 1
fi

# Commit
COMMIT_MSG="$TYPE$SCOPE: $DESC"
echo -e "\n${BLUE}Committing with message: ${GREEN}$COMMIT_MSG${NC}"
git commit -m "$COMMIT_MSG"

# 5. Push
read -p "Push to remote? (y/n): " PUSH
if [[ "$PUSH" == "y" ]]; then
    echo -e "${BLUE}Pushing...${NC}"
    git push
fi

echo -e "${GREEN}Done!${NC}"
