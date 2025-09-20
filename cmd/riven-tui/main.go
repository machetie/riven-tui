package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"riven-tui/pkg/config"
	"riven-tui/pkg/tui"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	version    = flag.Bool("version", false, "Show version information")
	help       = flag.Bool("help", false, "Show help information")
)

const (
	appName = "Riven TUI"
)

var (
	// appVersion can be overridden at build time with -ldflags "-X main.appVersion=v1.2.3"
	appVersion = "0.2.0"
)

func main() {
	flag.Parse()

	if *version {
		// appVersion may already include 'v' prefix when set via ldflags
		if strings.HasPrefix(appVersion, "v") {
			fmt.Printf("%s %s\n", appName, appVersion)
		} else {
			fmt.Printf("%s v%s\n", appName, appVersion)
		}
		os.Exit(0)
	}

	if *help {
		showHelp()
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate that we can connect to the API
	if err := validateConnection(cfg); err != nil {
		log.Fatalf("Failed to connect to Riven API: %v", err)
	}

	// Create and run the TUI application
	app := tui.NewApp(cfg)

	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

func showHelp() {
	fmt.Printf(`%s v%s

A Terminal User Interface (TUI) client for Riven media management system.

USAGE:
    riven-tui [OPTIONS]

OPTIONS:
    -config <path>    Path to configuration file
    -version          Show version information
    -help             Show this help message

CONFIGURATION:
    The application looks for configuration in the following order:
    1. File specified by -config flag
    2. ~/.config/riven-tui/config.yaml
    3. ~/.riven-tui.yaml
    4. ./config.yaml

    You can also use environment variables:
    - RIVEN_API_ENDPOINT: Riven API endpoint URL
    - RIVEN_API_TOKEN: Bearer token for authentication
    - RIVEN_API_TIMEOUT: Request timeout (default: 30s)

EXAMPLE CONFIG FILE:
    api:
      endpoint: "http://localhost:8080"
      token: "your-bearer-token-here"
      timeout: 30s
    ui:
      refresh_interval: 5s
      theme: "default"
      page_size: 50

NAVIGATION:
    Arrow Keys / hjkl: Navigate menus and lists
    Enter:             Select item or confirm action
    Tab:               Switch between panels
    Esc:               Go back or cancel
    q:                 Quit application
    r:                 Refresh current view
    ?:                 Show help

    d:                 Dashboard
    m:                 Media items
    s:                 Settings
    l:                 Logs

For more information, visit: https://github.com/rivenmedia/riven
`, appName, appVersion)
}

func validateConnection(cfg *config.Config) error {
	// Create a temporary client to test the connection
	client := tui.NewTestClient(cfg)

	// Try to get the root endpoint to validate connection
	if err := client.TestConnection(); err != nil {
		return fmt.Errorf("cannot connect to Riven API at %s: %w", cfg.API.Endpoint, err)
	}

	return nil
}
