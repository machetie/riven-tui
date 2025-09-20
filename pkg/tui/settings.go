package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"riven-tui/pkg/api"
)

// SettingsModel represents the settings screen
type SettingsModel struct {
	client   *api.Client
	ctx      context.Context
	width    int
	height   int
	loading  bool
	error    string
	settings interface{}

	// Auto-refresh
	lastUpdate time.Time
}

// settingsMsg represents messages for the settings screen
type settingsMsg struct {
	settings interface{}
	err      error
}

// NewSettingsModel creates a new settings model
func NewSettingsModel(client *api.Client, ctx context.Context) *SettingsModel {
	return &SettingsModel{
		client:  client,
		ctx:     ctx,
		loading: true,
	}
}

// SetSize sets the size of the settings screen
func (m *SettingsModel) SetSize(width, height int) {
	m.width = width
	m.height = height - 3 // Account for navigation bar
}

// Init implements tea.Model
func (m *SettingsModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchSettings(),
		m.autoRefresh(),
	)
}

// Update implements tea.Model
func (m *SettingsModel) Update(msg tea.Msg) (*SettingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case settingsMsg:
		m.loading = false
		if msg.err != nil {
			m.error = fmt.Sprintf("Failed to fetch settings: %v", msg.err)
		} else {
			m.settings = msg.settings
			m.error = ""
		}

	case refreshMsg:
		m.lastUpdate = time.Now()
		return m, tea.Batch(
			m.fetchSettings(),
			m.autoRefresh(),
		)

	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.loading = true
			return m, m.fetchSettings()
		}
	}

	return m, nil
}

// View implements tea.Model
func (m *SettingsModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.error != "" {
		return m.renderError()
	}

	return m.renderSettings()
}

// renderLoading renders the loading state
func (m *SettingsModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render("Loading settings...")
}

// renderError renders the error state
func (m *SettingsModel) renderError() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("196"))

	return style.Render(fmt.Sprintf("Error: %s\n\nPress 'r' to refresh", m.error))
}

// renderSettings renders the main settings view
func (m *SettingsModel) renderSettings() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Margin(0, 0, 2, 0).
		Render("⚙️ Riven Settings")

	var content string
	if m.settings != nil {
		// Try to display some actual settings data
		if settingsMap, ok := m.settings.(map[string]interface{}); ok {
			var settingsInfo []string
			settingsInfo = append(settingsInfo, "📊 Current Settings Overview:")
			settingsInfo = append(settingsInfo, "")

			// Display key settings categories
			categories := []string{"general", "scraping", "content", "downloaders", "indexers"}
			for _, category := range categories {
				if categoryData, exists := settingsMap[category]; exists {
					if categoryMap, ok := categoryData.(map[string]interface{}); ok {
						settingsInfo = append(settingsInfo, fmt.Sprintf("• %s: %d settings configured",
							strings.ToUpper(category[:1])+category[1:], len(categoryMap)))
					}
				}
			}

			settingsInfo = append(settingsInfo, "")
			settingsInfo = append(settingsInfo, "🔧 Available Actions:")
			settingsInfo = append(settingsInfo, "• View detailed settings (coming soon)")
			settingsInfo = append(settingsInfo, "• Export settings configuration")
			settingsInfo = append(settingsInfo, "• Validate current settings")
			settingsInfo = append(settingsInfo, "")
			settingsInfo = append(settingsInfo, "Note: Settings editing requires the Riven web interface.")
			settingsInfo = append(settingsInfo, "Press 'r' to refresh settings.")

			content = strings.Join(settingsInfo, "\n")
		} else {
			content = "Settings loaded successfully!\n\n" +
				"Note: Settings editing is read-only in this version.\n" +
				"Use the Riven web interface to modify settings.\n\n" +
				"Press 'r' to refresh settings."
		}
	} else {
		content = "No settings data available.\n\nPress 'r' to refresh."
	}

	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(2, 4).
		Margin(1, 0).
		Height(m.height - 8)

	body := lipgloss.JoinVertical(lipgloss.Left, title, contentStyle.Render(content))

	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(1, 2)

	return style.Render(body)
}

// fetchSettings fetches settings from the API
func (m *SettingsModel) fetchSettings() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		settings, err := m.client.GetAllSettings(m.ctx)
		return settingsMsg{settings: settings, err: err}
	})
}

// autoRefresh sets up auto-refresh
func (m *SettingsModel) autoRefresh() tea.Cmd {
	return tea.Tick(60*time.Second, func(t time.Time) tea.Msg {
		return refreshMsg{}
	})
}
