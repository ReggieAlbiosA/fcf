# Security Policy

## Supported Versions

We actively support the following versions of FCF with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 2.0.x   | :white_check_mark: |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please report it responsibly.

### How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **Email (Preferred)**: Send details to the maintainer at their GitHub profile
2. **GitHub Security Advisory**: Use GitHub's [private vulnerability reporting](https://github.com/ReggieAlbiosA/fcf/security/advisories/new) feature

### What to Include

When reporting a vulnerability, please include:

- **Description**: A clear description of the vulnerability
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Impact**: Potential impact and severity assessment
- **Suggested Fix**: If you have ideas on how to fix it (optional but appreciated)
- **Affected Versions**: Which versions are affected

### What to Expect

- **Acknowledgment**: You will receive an acknowledgment within 48 hours
- **Initial Assessment**: We will assess the vulnerability within 7 days
- **Updates**: We will keep you informed of our progress
- **Resolution**: We will work to resolve critical issues as quickly as possible

### Disclosure Policy

- We will coordinate with you before publicly disclosing the vulnerability
- We aim to release a fix before public disclosure
- You will be credited for the discovery (unless you prefer to remain anonymous)

## Security Best Practices

When using FCF, please follow these security best practices:

### Installation

- **Verify Installation Scripts**: Review installation scripts before running them
- **Use Official Sources**: Only download from the official GitHub repository
- **Check File Integrity**: Verify checksums when available

### Usage

- **Permissions**: Be cautious when running FCF with elevated permissions
- **Search Paths**: Be aware of what directories you're searching
- **Script Execution**: FCF does not execute files, but be careful with navigation to untrusted locations

### System Security

- **Keep Updated**: Regularly update FCF to the latest version
- **System Updates**: Keep your operating system and dependencies updated
- **Path Environment**: Ensure your PATH environment variable is secure

## Scope

### In Scope

- Remote code execution vulnerabilities
- Privilege escalation issues
- Information disclosure vulnerabilities
- Path traversal vulnerabilities
- Installation script security issues
- Authentication/authorization bypasses

### Out of Scope

- Issues requiring physical access to the system
- Issues requiring social engineering
- Denial of service (DoS) attacks that don't affect security
- Issues in third-party dependencies (please report to their maintainers)
- Issues in unsupported versions

## Security Updates

Security updates will be:

- Released as patch versions (e.g., 2.0.1 → 2.0.2)
- Documented in the [CHANGELOG.md](CHANGELOG.md)
- Tagged with security-related labels
- Announced via GitHub releases

## Contact

For security-related questions or concerns:

- **GitHub**: [@ReggieAlbiosA](https://github.com/ReggieAlbiosA)
- **Repository**: [fcf](https://github.com/ReggieAlbiosA/fcf)
- **Security Advisories**: [View Security Advisories](https://github.com/ReggieAlbiosA/fcf/security/advisories)

---

**Thank you for helping keep FCF and its users safe!**

