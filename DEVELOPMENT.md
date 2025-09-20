# Development Guide

This guide covers development setup, architecture, and contribution guidelines for Riven TUI.

## Development Setup

### Prerequisites
- Go 1.21 or later
- Git
- A running Riven instance for testing
- Terminal with 256-color support

### Environment Setup
```bash
# Clone the repository
git clone <repository-url>
cd riven-tui

# Install dependencies
go mod tidy

# Build the application
go build -o riven-tui cmd/riven-tui/main.go

# Run tests
go test ./...
```

### Development Configuration
Create a `config.dev.yaml` for development:
```yaml
api:
  endpoint: "http://localhost:8080"
  api_key: "your-dev-api-key"
  debug: true

ui:
  refresh_interval: "2s"
  theme: "default"

logging:
  level: "debug"
  file: "riven-tui-dev.log"
```

## Architecture Overview

### Project Structure
```
riven-tui/
├── cmd/riven-tui/          # Application entry point
│   └── main.go             # Main function and CLI setup
├── pkg/
│   ├── api/                # API client and HTTP handling
│   │   ├── client.go       # Main API client
│   │   └── endpoints.go    # API endpoint definitions
│   ├── config/             # Configuration management
│   │   ├── config.go       # Configuration struct and loading
│   │   └── validation.go   # Configuration validation
│   ├── models/             # Data models and types
│   │   ├── common.go       # Common types and utilities
│   │   ├── items.go        # Media item models
│   │   ├── settings.go     # Settings models
│   │   └── scraping.go     # Scraping-related models
│   └── tui/                # Terminal UI components
│       ├── app.go          # Main application controller
│       ├── dashboard.go    # Dashboard screen
│       ├── items.go        # Media browser screen
│       ├── item_detail.go  # Item detail view
│       ├── settings.go     # Settings screen
│       ├── logs.go         # Logs viewer
│       ├── components.go   # Reusable UI components
│       └── styles.go       # Themes and styling
├── config.yaml             # Default configuration
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── README.md               # User documentation
└── DEVELOPMENT.md          # This file
```

### Key Design Patterns

#### Model-View-Update (MVU)
The application follows the Bubble Tea MVU pattern:
- **Model**: Application state and data
- **View**: Rendering logic that converts state to UI
- **Update**: Event handling that modifies state

#### Component Architecture
- Reusable UI components in `components.go`
- Theme system for consistent styling
- Screen-based navigation with shared state

#### API Client Design
- Context-based request handling
- Automatic retry with exponential backoff
- Comprehensive error handling
- Type-safe request/response handling

## Development Workflow

### Adding New Features

1. **Plan the Feature**
   - Define the user story
   - Identify required API endpoints
   - Design the UI flow

2. **Implement Models**
   - Add data models in `pkg/models/`
   - Update API client if needed
   - Add configuration options

3. **Create UI Components**
   - Implement screen logic in `pkg/tui/`
   - Use existing components where possible
   - Follow theme system for styling

4. **Add Navigation**
   - Update main app controller
   - Add keyboard shortcuts
   - Update help system

5. **Test and Document**
   - Add unit tests
   - Update documentation
   - Test with real Riven instance

### Code Style Guidelines

#### Go Code Style
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and small

#### TUI Code Style
- Use theme system for all styling
- Handle window resize events
- Provide keyboard shortcuts
- Include loading states

#### Error Handling
- Always handle errors gracefully
- Provide user-friendly error messages
- Include retry mechanisms where appropriate
- Log errors for debugging

### Testing

#### Unit Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/api/
```

#### Integration Tests
```bash
# Test against real Riven instance
RIVEN_TEST_ENDPOINT=http://localhost:8080 go test ./pkg/api/ -tags=integration
```

#### Manual Testing
- Test with different terminal sizes
- Test with different themes
- Test error conditions
- Test with large datasets

### Debugging

#### Debug Mode
```bash
# Enable debug logging
RIVEN_TUI_DEBUG=1 ./riven-tui

# Use debug configuration
./riven-tui -config config.dev.yaml
```

#### Common Debug Techniques
- Add debug prints to Update functions
- Use Go debugger (delve) for complex issues
- Check API responses with debug logging
- Test UI components in isolation

## API Integration

### Adding New Endpoints

1. **Define Models**
   ```go
   // In pkg/models/
   type NewFeatureResponse struct {
       Data    interface{} `json:"data"`
       Message string      `json:"message"`
   }
   ```

2. **Add Client Method**
   ```go
   // In pkg/api/client.go
   func (c *Client) GetNewFeature(ctx context.Context) (*NewFeatureResponse, error) {
       var response NewFeatureResponse
       err := c.makeRequest(ctx, "GET", "/api/new-feature", nil, &response)
       return &response, err
   }
   ```

3. **Update TUI**
   ```go
   // In appropriate TUI file
   type newFeatureMsg struct {
       data *NewFeatureResponse
       err  error
   }
   
   func (m *Model) fetchNewFeature() tea.Cmd {
       return tea.Cmd(func() tea.Msg {
           data, err := m.client.GetNewFeature(m.ctx)
           return newFeatureMsg{data: data, err: err}
       })
   }
   ```

### Error Handling Patterns
- Use context for request cancellation
- Implement exponential backoff for retries
- Provide meaningful error messages to users
- Log detailed errors for debugging

## UI Development

### Theme System
```go
// Use theme for consistent styling
style := m.theme.CardStyle()
content := style.Render("Card content")

// State-specific styling
stateStyle := m.theme.StateStyle("Completed")
```

### Component Development
```go
// Create reusable components
loading := NewLoadingComponent("Loading data...", theme)
progress := NewProgressComponent("Processing...", theme)
```

### Responsive Design
```go
// Handle window resize
func (m *Model) SetSize(width, height int) {
    m.width = width
    m.height = height
    // Update child components
}
```

## Performance Considerations

### Memory Management
- Limit data retention in models
- Use pagination for large datasets
- Clean up resources in Update functions

### Network Efficiency
- Cache API responses when appropriate
- Use context for request cancellation
- Implement request debouncing

### UI Performance
- Minimize re-renders
- Use efficient string building
- Optimize table rendering for large datasets

## Contributing

### Pull Request Process
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests and documentation
5. Run quality checks
6. Submit pull request

### Code Review Guidelines
- Focus on code clarity and maintainability
- Ensure proper error handling
- Verify UI responsiveness
- Check for memory leaks
- Validate API integration

### Quality Checks
```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Vet code
go vet ./...

# Run tests
go test ./...

# Check for security issues
gosec ./...
```

## Release Process

### Version Management
- Use semantic versioning (MAJOR.MINOR.PATCH)
- Tag releases in Git
- Update version in main.go

### Build Process
```bash
# Build for multiple platforms
make build-all

# Create release packages
make package
```

### Documentation Updates
- Update README.md
- Update CHANGELOG.md
- Update API documentation
- Update configuration examples

## Troubleshooting

### Common Issues
- **Build failures**: Check Go version and dependencies
- **API connection issues**: Verify Riven instance and credentials
- **UI rendering issues**: Check terminal compatibility
- **Performance issues**: Profile with Go tools

### Debug Resources
- Go documentation: https://golang.org/doc/
- Bubble Tea examples: https://github.com/charmbracelet/bubbletea/tree/master/examples
- Lip Gloss styling: https://github.com/charmbracelet/lipgloss
