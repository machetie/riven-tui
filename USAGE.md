# Riven TUI Usage Guide

This guide provides comprehensive instructions for using the Riven TUI application.

## Quick Start

### 1. Installation
```bash
# Download and build
git clone <repository-url>
cd riven-tui
go build -o riven-tui cmd/riven-tui/main.go

# Or install directly
go install github.com/your-org/riven-tui/cmd/riven-tui@latest
```

### 2. Configuration
Create a configuration file at one of these locations:
- `~/.config/riven-tui/config.yaml` (recommended)
- `~/.riven-tui.yaml`
- `./config.yaml`

Minimal configuration:
```yaml
api:
  endpoint: "http://your-riven-instance:8080"
  token: "your-bearer-token"  # Optional
```

### 3. First Run
```bash
riven-tui
```

## Screen Overview

### Dashboard (Press 'd')
The dashboard provides an overview of your Riven system:
- **System Statistics**: CPU, memory, disk usage
- **Service Status**: Status of Riven services
- **Recent Activity**: Latest system events
- **Quick Actions**: Common operations

**Navigation:**
- `r` - Refresh data
- `↑/↓` - Navigate between sections
- `Enter` - Expand/collapse sections

### Media Browser (Press 'm')
Browse and manage your media library:
- **Item List**: Paginated list of media items
- **Search**: Find specific items
- **Filters**: Filter by state, type, etc.
- **Actions**: Perform operations on items

**Navigation:**
- `/` - Open search
- `f` - Toggle filters (Failed → Completed → Downloaded → All)
- `s` - Cycle sort order (Date ↓ → Date ↑ → Title ↑ → Title ↓)
- `c` - Clear search and filters
- `n/p` - Next/previous page
- `a` - Show actions menu
- `Enter` - View item details
- `↑/↓` - Navigate items

**Search Tips:**
- Search by title, year, or ID
- Use partial matches
- Search is case-insensitive

**Available Actions:**
- Retry failed items
- Reset item state
- Pause/unpause items
- Remove items from library

### Item Details (Automatic when selecting item)
Detailed view of a specific media item with three tabs:

**Details Tab (1):**
- Basic information (title, year, type)
- External IDs (TMDB, TVDB, IMDB)
- Overview and metadata

**Streams Tab (2):**
- Available torrent streams
- Quality, size, and cache status
- Stream rankings and metadata

**Actions Tab (3):**
- Item-specific actions
- Retry processing
- Reset state
- Manage streams

**Navigation:**
- `1-3` - Switch tabs directly
- `Tab` - Next tab
- `Shift+Tab` - Previous tab
- `Esc` - Return to media browser
- `↑/↓` - Navigate streams (in Streams tab)

### Settings (Press 's')
View and manage Riven configuration:
- **Settings Overview**: Summary of configuration
- **Categories**: Organized settings groups
- **Validation**: Check configuration validity

**Note**: Settings editing requires the Riven web interface. The TUI provides read-only access.

**Navigation:**
- `r` - Refresh settings
- `↑/↓` - Navigate settings

### Logs (Press 'l')
Monitor system logs and events:
- **Log Entries**: Recent system logs
- **Filtering**: Filter by log level
- **Real-time**: Live log updates

**Navigation:**
- `r` - Refresh logs
- `f` - Filter by level
- `c` - Clear displayed logs
- `↑/↓` - Navigate log entries

### Help (Press '?')
Interactive help system:
- **Keyboard Shortcuts**: All available keybindings
- **Context Help**: Help for current screen
- **Tips**: Usage tips and tricks

## Advanced Features

### Real-time Updates
The TUI automatically refreshes data at configurable intervals:
- Dashboard: Every 5 seconds (default)
- Media items: Every 30 seconds
- Logs: Real-time streaming

Configure refresh intervals in your config file:
```yaml
ui:
  refresh_interval: "10s"  # Global refresh interval
```

### Themes
Choose from multiple themes:
```yaml
ui:
  theme: "dark"  # "default", "dark", "light"
```

**Theme Characteristics:**
- **Default**: Balanced colors for most terminals
- **Dark**: Optimized for dark backgrounds
- **Light**: Optimized for light backgrounds

### Search and Filtering
Advanced search capabilities:
- **Text Search**: Search titles, descriptions
- **State Filtering**: Filter by item state
- **Type Filtering**: Movies vs TV shows
- **Combined Filters**: Use multiple filters together

### Keyboard Shortcuts

#### Global Shortcuts
- `q` - Quit application
- `?` - Toggle help
- `r` - Refresh current view
- `Esc` - Go back/cancel

#### Screen Navigation
- `d` - Dashboard
- `m` - Media browser
- `s` - Settings
- `l` - Logs

#### Movement
- `↑/↓` or `k/j` - Navigate up/down
- `←/→` or `h/l` - Navigate left/right
- `Enter` - Select/confirm
- `Space` - Toggle selection

#### Media Browser Specific
- `/` - Search
- `f` - Filter
- `s` - Sort
- `c` - Clear
- `n/p` - Page navigation
- `a` - Actions menu

## Configuration Guide

### API Configuration
```yaml
api:
  endpoint: "http://localhost:8080"
  token: "your-bearer-token"
  timeout: "30s"
  debug: false
```

### UI Customization
```yaml
ui:
  theme: "default"
  page_size: 50
  refresh_interval: "5s"
  mouse_support: false
```

### Performance Tuning
```yaml
performance:
  max_concurrent_requests: 10
  cache:
    enabled: true
    ttl:
      items: "30s"
      stats: "10s"
```

## Troubleshooting

### Connection Issues
1. **Check Riven Status**: Ensure Riven is running
   ```bash
   curl http://your-riven-instance:8080/health
   ```

2. **Verify Bearer Token**: Test authentication
   ```bash
   curl -H "Authorization: Bearer your-token" http://your-riven-instance:8080/api/v1/stats
   ```

3. **Check Configuration**: Validate config file
   ```bash
   riven-tui --validate-config
   ```

### Performance Issues
1. **Reduce Refresh Interval**: Increase refresh intervals
2. **Decrease Page Size**: Show fewer items per page
3. **Disable Animations**: Turn off UI animations
4. **Check Network**: Verify connection to Riven

### Display Issues
1. **Terminal Compatibility**: Ensure 256-color support
2. **Terminal Size**: Minimum 80x24 recommended
3. **Font**: Use monospace font
4. **Theme**: Try different themes

### Common Error Messages

**"Connection refused"**
- Riven is not running
- Wrong endpoint URL
- Network connectivity issues

**"Unauthorized"**
- Invalid or missing Bearer token
- Bearer token not configured in Riven

**"Timeout"**
- Slow network connection
- Riven instance overloaded
- Increase timeout in configuration

## Tips and Best Practices

### Efficient Navigation
- Use keyboard shortcuts instead of mouse
- Learn the most common shortcuts (`d`, `m`, `s`, `l`)
- Use search (`/`) to quickly find items

### Performance Optimization
- Adjust page size based on your library size
- Use filters to reduce data loading
- Configure appropriate refresh intervals

### Monitoring
- Keep the dashboard open for system overview
- Use real-time logs for troubleshooting
- Monitor system resources during heavy operations

### Configuration Management
- Keep configuration in version control
- Use environment variables for sensitive data
- Test configuration changes in development

## Integration Examples

### Environment Variables
```bash
export RIVEN_API_ENDPOINT="http://localhost:8080"
export RIVEN_API_TOKEN="your-bearer-token"
export RIVEN_API_TIMEOUT="60s"
riven-tui
```

### Docker Usage
```bash
docker run -it \
  -v ~/.config/riven-tui:/root/.config/riven-tui \
  -e RIVEN_API_ENDPOINT="http://host.docker.internal:8080" \
  riven-tui
```

### Systemd Service
```ini
[Unit]
Description=Riven TUI
After=network.target

[Service]
Type=simple
User=riven
ExecStart=/usr/local/bin/riven-tui
Restart=always
Environment=RIVEN_API_ENDPOINT=http://localhost:8080

[Install]
WantedBy=multi-user.target
```

## Getting Help

- **Built-in Help**: Press `?` in the application
- **Configuration**: See `examples/config-example.yaml`
- **Issues**: Report bugs on GitHub
- **Documentation**: Check README.md for latest updates
