# Contributing to Riven TUI

Thank you for your interest in contributing to Riven TUI! This document provides guidelines and information for contributors.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and constructive in all interactions.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- A running Riven instance for testing
- Basic familiarity with terminal applications and Go

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/riven-tui.git
   cd riven-tui
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Run Tests**
   ```bash
   go test ./...
   ```

4. **Build and Test**
   ```bash
   go build -o riven-tui cmd/riven-tui/main.go
   ./riven-tui --help
   ```

## Development Guidelines

### Code Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and small
- Use proper error handling

### Project Structure

```
riven-tui/
├── cmd/riven-tui/          # Main application entry point
├── pkg/
│   ├── api/                # Riven API client
│   ├── config/             # Configuration management
│   ├── models/             # Data models
│   └── tui/                # TUI components and logic
├── examples/               # Example configurations
├── .github/workflows/      # CI/CD workflows
└── docs/                   # Documentation
```

### Making Changes

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Your Changes**
   - Write clean, well-documented code
   - Add tests for new functionality
   - Update documentation as needed

3. **Test Your Changes**
   ```bash
   go test ./...
   go build -o riven-tui cmd/riven-tui/main.go
   ./riven-tui --version
   ```

4. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

### Commit Message Format

Use conventional commits format:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

Examples:
- `feat: add media filtering by genre`
- `fix: resolve pagination issue in item browser`
- `docs: update installation instructions`

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/api/
```

### Writing Tests

- Add unit tests for new functions
- Test error conditions and edge cases
- Use table-driven tests where appropriate
- Mock external dependencies (API calls, etc.)

## Submitting Changes

### Pull Request Process

1. **Update Documentation**
   - Update README.md if needed
   - Add/update code comments
   - Update USAGE.md for new features

2. **Create Pull Request**
   - Use a descriptive title
   - Provide detailed description of changes
   - Reference any related issues
   - Include screenshots for UI changes

3. **Review Process**
   - Address reviewer feedback
   - Keep PR focused and atomic
   - Ensure CI passes

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] Added tests for new functionality
- [ ] Manual testing completed

## Screenshots (if applicable)
Add screenshots for UI changes
```

## Release Process

Releases are automated via GitHub Actions:

1. **Version Bump**
   - Update version in `cmd/riven-tui/main.go`
   - Update CHANGELOG.md

2. **Create Release**
   ```bash
   git tag v0.x.x
   git push origin v0.x.x
   ```

3. **Automated Build**
   - GitHub Actions builds cross-platform binaries
   - Creates GitHub release with assets
   - Generates SHA256 checksums

## Getting Help

- **Issues**: Open an issue for bugs or feature requests
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: Check README.md and USAGE.md

## Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes for significant contributions
- Special thanks in documentation

Thank you for contributing to Riven TUI! 🎉
