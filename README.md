# Riven TUI

[![Release](https://img.shields.io/github/v/release/machetie/riven-tui)](https://github.com/machetie/riven-tui/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/machetie/riven-tui)](https://golang.org/)
[![License](https://img.shields.io/github/license/machetie/riven-tui)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/machetie/riven-tui/release.yml)](https://github.com/machetie/riven-tui/actions)

A comprehensive Terminal User Interface (TUI) client for the Riven media management system, built with Go and the Bubble Tea framework.

## Features

### 🎯 Core Functionality
- **Dashboard**: Real-time system statistics and health monitoring
- **Media Browser**: Browse, search, and manage your media library with advanced filtering
- **Item Details**: Comprehensive item view with streams, metadata, and actions
- **Settings Management**: View and configure Riven settings
- **Logs Viewer**: Monitor system logs and events
- **Interactive Help**: Built-in help system with keyboard shortcuts

### 🎨 User Interface
- **Modern TUI Design**: Clean, responsive terminal interface
- **Multiple Themes**: Default, dark, and light themes with customizable styling
- **Responsive Layout**: Adapts to different terminal sizes
- **Rich Components**: Tables, progress bars, spinners, toast notifications, and confirmation dialogs
- **Keyboard Navigation**: Intuitive vim-like key bindings

### 🔧 Advanced Features
- **Real-time Updates**: Auto-refreshing data with configurable intervals
- **Search & Filtering**: Powerful search and filtering capabilities with state-based filters
- **Pagination**: Efficient handling of large datasets
- **Error Handling**: Comprehensive error handling with retry mechanisms and user feedback
- **Configuration Management**: Flexible YAML-based configuration system

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/machetie/riven-tui/releases):

```bash
# Linux (amd64)
curl -L -o riven-tui https://github.com/machetie/riven-tui/releases/download/v0.2.0/riven-tui-linux-amd64
chmod +x riven-tui

# macOS (Intel)
curl -L -o riven-tui https://github.com/machetie/riven-tui/releases/download/v0.2.0/riven-tui-darwin-amd64
chmod +x riven-tui

# macOS (Apple Silicon)
curl -L -o riven-tui https://github.com/machetie/riven-tui/releases/download/v0.2.0/riven-tui-darwin-arm64
chmod +x riven-tui

# Windows
curl -L -o riven-tui.exe https://github.com/machetie/riven-tui/releases/download/v0.2.0/riven-tui-windows-amd64.exe
```

### Prerequisites

- Access to a running Riven instance
- Bearer token for Riven API authentication

### Build from Source

```bash
git clone https://github.com/machetie/riven-tui.git
cd riven-tui
go build -o riven-tui cmd/riven-tui/main.go
```

### Install with Go

```bash
go install github.com/machetie/riven-tui/cmd/riven-tui@latest
```

## Configuration

> **⚠️ Breaking Change Notice**: Starting from version 0.2.0, the authentication method has changed from `api_key` to `token` using standard Bearer token authentication. Please update your configuration files to use `token:` instead of `api_key:`. The old `RIVEN_API_KEY` environment variable is still supported but deprecated - please use `RIVEN_API_TOKEN` instead.

Create a configuration file at `~/.config/riven-tui/config.yaml`:

```yaml
api:
  endpoint: "http://localhost:8080"
  token: "your-bearer-token-here"
  timeout: 30s

ui:
  refresh_interval: 5s
  theme: "default"
```

### Environment Variables

You can also configure using environment variables:

- `RIVEN_API_ENDPOINT`: Riven API endpoint URL
- `RIVEN_API_TOKEN`: Bearer token for authentication
- `RIVEN_API_TIMEOUT`: Request timeout (default: 30s)

## Usage

```bash
# Start the TUI
./riven-tui

# Specify custom config file
./riven-tui --config /path/to/config.yaml

# Use environment variables
RIVEN_API_ENDPOINT=http://localhost:8080 RIVEN_API_TOKEN=your-token ./riven-tui
```

### Navigation

- **Arrow Keys / hjkl**: Navigate menus and lists
- **Enter**: Select item or confirm action
- **Tab**: Switch between panels
- **Esc**: Go back or cancel
- **q**: Quit application
- **r**: Refresh current view
- **?**: Show help

### Screens

1. **Dashboard** (`d`): System overview and quick stats
2. **Media Items** (`m`): Browse and manage media library
3. **Settings** (`s`): Configure Riven settings
4. **Logs** (`l`): View system logs and events
5. **Help** (`?`): Show keyboard shortcuts

## Development

### Project Structure

```
riven-tui/
├── cmd/riven-tui/          # Main application entry point
├── pkg/
│   ├── api/                # API client and HTTP handling
│   ├── config/             # Configuration management
│   ├── models/             # Data models and structures
│   └── tui/                # TUI components and screens
├── go.mod
├── go.sum
└── README.md
```

### Running in Development

```bash
go run cmd/riven-tui/main.go
```

### Building

```bash
# Build for current platform
go build -o riven-tui cmd/riven-tui/main.go

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o riven-tui-linux-amd64 cmd/riven-tui/main.go
GOOS=windows GOARCH=amd64 go build -o riven-tui-windows-amd64.exe cmd/riven-tui/main.go
GOOS=darwin GOARCH=amd64 go build -o riven-tui-darwin-amd64 cmd/riven-tui/main.go
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Quick Start for Contributors

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Setup

```bash
git clone https://github.com/machetie/riven-tui.git
cd riven-tui
go mod download
go test ./...
```

## Release Process

Releases are automated via GitHub Actions when a new tag is pushed:

```bash
git tag v0.3.0
git push origin v0.3.0
```

This will automatically:
- Build cross-platform binaries
- Generate SHA256 checksums
- Create a GitHub release with all assets

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The TUI framework that powers this application
- [Riven](https://github.com/rivenmedia/riven) - The media management system this TUI interfaces with
- [Bubbles](https://github.com/charmbracelet/bubbles) - Pre-built TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling library for terminal applications

---

**Made with ❤️ for the Riven community**
