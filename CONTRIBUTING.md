# Contributing to FCF

Thank you for your interest in contributing to FCF (Find File or Folder)! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs **actual behavior**
- **Environment details** (OS, version, shell)
- **Screenshots or error messages** (if applicable)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub Discussions. Please use the [Feature Idea template](.github/DISCUSSION_TEMPLATE/feature-idea.md) when suggesting new features.

### Pull Requests

Pull requests are welcome! Please follow these guidelines:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make your changes**
4. **Test your changes** on your platform
5. **Commit your changes** (`git commit -m 'Add amazing feature'`)
6. **Push to the branch** (`git push origin feature/amazing-feature`)
7. **Open a Pull Request**

## Development Setup

### Project Structure

```
fcf/
├── ubuntu/           # Linux/Unix implementation (Bash)
│   ├── fcf.sh        # Main script
│   └── install.sh    # Installer
├── win/              # Windows implementation
│   ├── go/           # Go source code
│   └── install.ps1   # Installer
├── .github/          # GitHub templates and workflows
└── wiki/             # Wiki documentation
```

### Linux/Ubuntu Development

The Linux version is a Bash script. To develop:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/ReggieAlbiosA/fcf.git
   cd fcf
   ```

2. **Make the script executable:**
   ```bash
   chmod +x ubuntu/fcf.sh
   ```

3. **Test locally:**
   ```bash
   ./ubuntu/fcf.sh
   # or create a symlink
   ln -s $(pwd)/ubuntu/fcf.sh ~/.local/bin/fcf
   ```

4. **Test the installer:**
   ```bash
   ./ubuntu/install.sh
   ```

### Windows Development

The Windows version is written in Go. To develop:

1. **Install Go** (1.19 or later):
   ```powershell
   winget install GoLang.Go
   ```

2. **Clone the repository:**
   ```powershell
   git clone https://github.com/ReggieAlbiosA/fcf.git
   cd fcf
   ```

3. **Navigate to Go source:**
   ```powershell
   cd win/go
   ```

4. **Build the binary:**
   ```powershell
   go build -o fcf.exe
   ```

5. **Test locally:**
   ```powershell
   .\fcf.exe
   ```

6. **Run tests (if available):**
   ```powershell
   go test ./...
   ```

## Coding Standards

### Bash Script (Linux)

- Use **4 spaces** for indentation (no tabs)
- Follow **shellcheck** guidelines when possible
- Add comments for complex logic
- Use descriptive variable names
- Keep functions focused and small
- Test on multiple shells (bash, zsh) if possible

### Go Code (Windows)

- Follow **gofmt** formatting (run `go fmt ./...`)
- Follow **golint** guidelines
- Add comments for exported functions
- Keep functions focused and testable
- Write unit tests for new features

### General Guidelines

- **Keep it simple** - FCF is meant to be straightforward
- **Maintain cross-platform consistency** - Features should work similarly on both platforms
- **Update documentation** - If you add features, update README.md and wiki
- **Update CHANGELOG.md** - Document your changes

## Testing

### Before Submitting

- Test your changes on your platform
- Test both interactive and direct modes
- Test edge cases (empty results, special characters, etc.)
- Verify the installer still works
- Check that existing functionality still works

### Test Scenarios

- Search for files with different patterns
- Test navigation functionality
- Test with `fd` installed and without
- Test with various file types and extensions
- Test with hidden files
- Test in different directory structures

## Commit Messages

Write clear, descriptive commit messages:

```
feat: Add fuzzy search support
fix: Resolve navigation path issue on Windows
docs: Update installation instructions
refactor: Simplify search logic
test: Add tests for pattern matching
```

Use conventional commit prefixes:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

## Pull Request Process

1. **Update CHANGELOG.md** - Add your changes under "Unreleased" or appropriate version
2. **Update documentation** - README.md, wiki, or both if needed
3. **Ensure tests pass** - If applicable
4. **Request review** - Tag maintainers if needed
5. **Respond to feedback** - Be open to suggestions and improvements

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] No new warnings or errors
- [ ] Tested on target platform(s)

## Documentation

When adding features, please update:

- **README.md** - Main documentation
- **CHANGELOG.md** - Version history
- **Wiki** (if applicable) - Detailed guides
- **Code comments** - Inline documentation

## Questions?

- Open a [GitHub Discussion](https://github.com/ReggieAlbiosA/fcf/discussions)
- Check existing [Issues](https://github.com/ReggieAlbiosA/fcf/issues)
- Review the [README](README.md) and [Wiki](https://github.com/ReggieAlbiosA/fcf/wiki)

## License

By contributing, you agree that your contributions will be licensed under the same [MIT License](LICENSE) that covers the project.

---

**Thank you for contributing to FCF! 🎉**



