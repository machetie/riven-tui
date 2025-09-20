package tui

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"riven-tui/pkg/api"
	"riven-tui/pkg/models"
)

// DashboardModel represents the dashboard screen
type DashboardModel struct {
	client  *api.Client
	ctx     context.Context
	width   int
	height  int
	loading bool
	error   string

	// Data
	stats    *models.StatsResponse
	services models.ServicesResponse
	rdUser   *models.RDUser

	// Auto-refresh
	lastUpdate time.Time
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel(client *api.Client, ctx context.Context) *DashboardModel {
	return &DashboardModel{
		client:  client,
		ctx:     ctx,
		loading: true,
	}
}

// SetSize sets the size of the dashboard
func (m *DashboardModel) SetSize(width, height int) {
	m.width = width
	m.height = height - 3 // Account for navigation bar
}

// Init implements tea.Model
func (m *DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchStats(),
		m.fetchServices(),
		m.fetchRDUser(),
		m.autoRefresh(),
	)
}

// Update implements tea.Model
func (m *DashboardModel) Update(msg tea.Msg) (*DashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case statsMsg:
		m.loading = false
		if msg.err != nil {
			m.error = fmt.Sprintf("Failed to fetch stats: %v", msg.err)
		} else {
			m.stats = msg.stats
			m.error = ""
		}

	case servicesMsg:
		if msg.err != nil {
			m.error = fmt.Sprintf("Failed to fetch services: %v", msg.err)
		} else {
			m.services = msg.services
		}

	case rdUserMsg:
		if msg.err != nil {
			// RD user might not be configured, don't show error
		} else {
			m.rdUser = msg.rdUser
		}

	case refreshMsg:
		m.lastUpdate = time.Now()
		return m, tea.Batch(
			m.fetchStats(),
			m.fetchServices(),
			m.fetchRDUser(),
			m.autoRefresh(),
		)
	}

	return m, nil
}

// View implements tea.Model
func (m *DashboardModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.error != "" {
		return m.renderError()
	}

	return m.renderDashboard()
}

// renderLoading renders the loading state
func (m *DashboardModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render("Loading dashboard...")
}

// renderError renders the error state
func (m *DashboardModel) renderError() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("196"))

	return style.Render(fmt.Sprintf("Error: %s\n\nPress 'r' to refresh", m.error))
}

// renderDashboard renders the main dashboard
func (m *DashboardModel) renderDashboard() string {
	var sections []string

	// Stats section
	if m.stats != nil {
		sections = append(sections, m.renderStats())
	}

	// Services section
	if m.services != nil {
		sections = append(sections, m.renderServices())
	}

	// RD User section
	if m.rdUser != nil {
		sections = append(sections, m.renderRDUser())
	}

	// Last update info
	sections = append(sections, m.renderLastUpdate())

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(1, 2)

	return style.Render(content)
}

// renderStats renders the statistics section
func (m *DashboardModel) renderStats() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Render("📊 System Statistics")

	statsStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Margin(1, 0)

	content := fmt.Sprintf(
		"Total Items: %d\n"+
			"Movies: %d | Shows: %d | Seasons: %d | Episodes: %d\n"+
			"Symlinks: %d | Incomplete: %d\n\n"+
			"States:\n"+
			"  Completed: %d | Downloaded: %d | Symlinked: %d\n"+
			"  Failed: %d | Paused: %d | Requested: %d",
		m.stats.TotalItems,
		m.stats.TotalMovies, m.stats.TotalShows, m.stats.TotalSeasons, m.stats.TotalEpisodes,
		m.stats.TotalSymlinks, m.stats.IncompleteItems,
		m.stats.States[models.StateCompleted],
		m.stats.States[models.StateDownloaded],
		m.stats.States[models.StateSymlinked],
		m.stats.States[models.StateFailed],
		m.stats.States[models.StatePaused],
		m.stats.States[models.StateRequested],
	)

	return lipgloss.JoinVertical(lipgloss.Left, title, statsStyle.Render(content))
}

// renderServices renders the services section
func (m *DashboardModel) renderServices() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Render("🔧 Services Status")

	servicesStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Margin(1, 0)

	var serviceList []string
	for service, status := range m.services {
		statusIcon := "❌"
		statusColor := lipgloss.Color("196")
		if status {
			statusIcon = "✅"
			statusColor = lipgloss.Color("46")
		}

		serviceList = append(serviceList, fmt.Sprintf(
			"%s %s",
			lipgloss.NewStyle().Foreground(statusColor).Render(statusIcon),
			service,
		))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, serviceList...)

	return lipgloss.JoinVertical(lipgloss.Left, title, servicesStyle.Render(content))
}

// renderRDUser renders the Real-Debrid user section
func (m *DashboardModel) renderRDUser() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Render("👤 Real-Debrid User")

	userStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Margin(1, 0)

	premiumDays := m.rdUser.Premium / (24 * 60 * 60) // Convert seconds to days

	content := fmt.Sprintf(
		"Username: %s\n"+
			"Type: %s\n"+
			"Points: %d\n"+
			"Premium Days Left: %d",
		m.rdUser.Username,
		m.rdUser.Type,
		m.rdUser.Points,
		premiumDays,
	)

	return lipgloss.JoinVertical(lipgloss.Left, title, userStyle.Render(content))
}

// renderLastUpdate renders the last update timestamp
func (m *DashboardModel) renderLastUpdate() string {
	if m.lastUpdate.IsZero() {
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Margin(1, 0)

	return style.Render(fmt.Sprintf("Last updated: %s", m.lastUpdate.Format("15:04:05")))
}

// Message types for async operations
type statsMsg struct {
	stats *models.StatsResponse
	err   error
}

type servicesMsg struct {
	services models.ServicesResponse
	err      error
}

type rdUserMsg struct {
	rdUser *models.RDUser
	err    error
}

// fetchStats fetches system statistics
func (m *DashboardModel) fetchStats() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		stats, err := m.client.GetStats(m.ctx)
		return statsMsg{stats: stats, err: err}
	})
}

// fetchServices fetches service status
func (m *DashboardModel) fetchServices() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		services, err := m.client.GetServices(m.ctx)
		return servicesMsg{services: services, err: err}
	})
}

// fetchRDUser fetches Real-Debrid user info
func (m *DashboardModel) fetchRDUser() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		rdUser, err := m.client.GetRDUser(m.ctx)
		return rdUserMsg{rdUser: rdUser, err: err}
	})
}

// autoRefresh sets up auto-refresh
func (m *DashboardModel) autoRefresh() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return refreshMsg{}
	})
}
