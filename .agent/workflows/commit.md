---
description: Stage, commit with conventional commits, and push changes
---

# Conventional Commit and Push Workflow

This workflow helps you stage changes, create a message following the [Conventional Commits](https://www.conventionalcommits.org/) specification, and push to the remote repository.

## Automated Helper

You can run the interactive helper script:
// turbo
`./.scripts/commit.sh`

## Manual Steps

1. **Review Changes**
   List all modified and untracked files.
   // turbo
   `git status`

2. **Stage Changes**
   Ask the user which files to stage. You can stage all with `git add .` or specific files.
   // turbo
   `git add .` (Defaulting to all changes for simplicity, but confirm with user if they want subset)

3. **Generate Conventional Commit Message**
   Identify the type of change (feat, fix, refactor, chore, docs, etc.), an optional scope, and a concise description.

   Example format: `<type>(<scope>): <description>`

4. **Commit**
   // turbo
   `git commit -m "<type>(<scope>): <description>"`

5. **Push**
   // turbo
   `git push`

---

_Note: If you have many changes, I can help you group them into logical commits._
