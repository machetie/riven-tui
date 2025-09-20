package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"riven-tui/pkg/api"
	"riven-tui/pkg/config"
)

// Screen represents different screens in the application
type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenItems
	ScreenItemDetail
	ScreenSettings
	ScreenLogs
	ScreenHelp
)

// App represents the main application state
type App struct {
	client        *api.Client
	config        *config.Config
	currentScreen Screen
	width         int
	height        int
	ctx           context.Context

	// Screen models
	dashboard  *DashboardModel
	items      *ItemsModel
	itemDetail *ItemDetailModel
	settings   *SettingsModel
	logs       *LogsModel
	help       *HelpModel

	// Navigation
	keys KeyMap

	// Theme and UI components
	theme        Theme
	loading      *LoadingComponent
	error        *ErrorComponent
	status       *StatusComponent
	confirmation *ConfirmationComponent
	toast        *ToastComponent
	helpOverlay  *HelpComponent

	// UI state
	showHelp    bool
	showConfirm bool
	lastError   string
}

// Common message types
type refreshMsg struct{}
type errorMsg struct {
	err error
}
type statusMsg struct {
	message    string
	statusType StatusType
}
type toastMsg struct {
	message    string
	statusType StatusType
	duration   time.Duration
}

// KeyMap defines the key bindings
type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Enter   key.Binding
	Back    key.Binding
	Quit    key.Binding
	Refresh key.Binding
	Help    key.Binding

	// Screen navigation
	Dashboard key.Binding
	Items     key.Binding
	Settings  key.Binding
	Logs      key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Dashboard: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "dashboard"),
		),
		Items: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "media items"),
		),
		Settings: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "settings"),
		),
		Logs: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "logs"),
		),
	}
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config) *App {
	client := api.NewClient(cfg)
	ctx := context.Background()

	app := &App{
		client:        client,
		config:        cfg,
		currentScreen: ScreenDashboard,
		ctx:           ctx,
		keys:          DefaultKeyMap(),
		theme:         GetTheme(cfg.UI.Theme),
	}

	// Initialize screen models
	app.dashboard = NewDashboardModel(client, ctx)
	app.items = NewItemsModel(client, ctx)
	app.settings = NewSettingsModel(client, ctx)
	app.logs = NewLogsModel(client, ctx)
	app.help = NewHelpModel(app.keys)

	return app
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.dashboard.Init(),
		tea.EnterAltScreen,
	)
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		// Update all screen models with new size
		a.dashboard.SetSize(msg.Width, msg.Height)
		a.items.SetSize(msg.Width, msg.Height)
		a.settings.SetSize(msg.Width, msg.Height)
		a.logs.SetSize(msg.Width, msg.Height)
		a.help.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		// Handle escape key for navigation
		if msg.String() == "esc" && a.currentScreen == ScreenItemDetail {
			a.currentScreen = ScreenItems
			return a, nil
		}

		// Global key bindings
		switch {
		case key.Matches(msg, a.keys.Quit):
			return a, tea.Quit
		case key.Matches(msg, a.keys.Dashboard):
			a.currentScreen = ScreenDashboard
			return a, a.dashboard.Init()
		case key.Matches(msg, a.keys.Items):
			a.currentScreen = ScreenItems
			return a, a.items.Init()
		case key.Matches(msg, a.keys.Settings):
			a.currentScreen = ScreenSettings
			return a, a.settings.Init()
		case key.Matches(msg, a.keys.Logs):
			a.currentScreen = ScreenLogs
			return a, a.logs.Init()
		case key.Matches(msg, a.keys.Help):
			a.currentScreen = ScreenHelp
			return a, nil
		}

	case showItemDetailMsg:
		// Navigate to item detail view
		a.itemDetail = NewItemDetailModel(a.client, a.ctx, msg.itemID)
		a.itemDetail.SetSize(a.width, a.height)
		a.currentScreen = ScreenItemDetail
		return a, a.itemDetail.Init()
	}

	// Update current screen
	switch a.currentScreen {
	case ScreenDashboard:
		a.dashboard, cmd = a.dashboard.Update(msg)
		cmds = append(cmds, cmd)
	case ScreenItems:
		a.items, cmd = a.items.Update(msg)
		cmds = append(cmds, cmd)
	case ScreenItemDetail:
		if a.itemDetail != nil {
			a.itemDetail, cmd = a.itemDetail.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScreenSettings:
		a.settings, cmd = a.settings.Update(msg)
		cmds = append(cmds, cmd)
	case ScreenLogs:
		a.logs, cmd = a.logs.Update(msg)
		cmds = append(cmds, cmd)
	case ScreenHelp:
		a.help, cmd = a.help.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

// View implements tea.Model
func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}

	var content string

	switch a.currentScreen {
	case ScreenDashboard:
		content = a.dashboard.View()
	case ScreenItems:
		content = a.items.View()
	case ScreenItemDetail:
		if a.itemDetail != nil {
			content = a.itemDetail.View()
		} else {
			content = "Item detail not available"
		}
	case ScreenSettings:
		content = a.settings.View()
	case ScreenLogs:
		content = a.logs.View()
	case ScreenHelp:
		content = a.help.View()
	default:
		content = "Unknown screen"
	}

	// Add navigation bar
	navBar := a.renderNavBar()

	// Combine navigation and content
	return lipgloss.JoinVertical(
		lipgloss.Left,
		navBar,
		content,
	)
}

// renderNavBar renders the navigation bar
func (a *App) renderNavBar() string {
	var tabs []string

	screens := []struct {
		screen Screen
		name   string
		key    string
	}{
		{ScreenDashboard, "Dashboard", "d"},
		{ScreenItems, "Media", "m"},
		{ScreenItemDetail, "Detail", ""},
		{ScreenSettings, "Settings", "s"},
		{ScreenLogs, "Logs", "l"},
		{ScreenHelp, "Help", "?"},
	}

	for _, s := range screens {
		// Skip item detail if not active
		if s.screen == ScreenItemDetail && a.currentScreen != ScreenItemDetail {
			continue
		}

		var style lipgloss.Style
		if s.screen == a.currentScreen {
			style = a.theme.ActiveTabStyle()
		} else {
			style = a.theme.TabStyle()
		}

		tabText := fmt.Sprintf("%s (%s)", s.name, s.key)
		tabs = append(tabs, style.Render(tabText))
	}

	navStyle := a.theme.NavBarStyle().Width(a.width)
	return navStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, tabs...))
}

// TestClient is a helper for testing API connectivity
type TestClient struct {
	client *api.Client
}

// NewTestClient creates a new test client
func NewTestClient(cfg *config.Config) *TestClient {
	return &TestClient{
		client: api.NewClient(cfg),
	}
}

// TestConnection tests the API connection
func (tc *TestClient) TestConnection() error {
	ctx := context.Background()
	_, err := tc.client.Health(ctx)
	return err
}
