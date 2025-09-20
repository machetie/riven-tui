package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"riven-tui/pkg/api"
	"riven-tui/pkg/models"
)

// LogsModel represents the logs screen
type LogsModel struct {
	client   *api.Client
	ctx      context.Context
	width    int
	height   int
	loading  bool
	error    string
	
	// Data
	logs     *models.LogsResponse
	viewport viewport.Model
}

// NewLogsModel creates a new logs model
func NewLogsModel(client *api.Client, ctx context.Context) *LogsModel {
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2)
	
	return &LogsModel{
		client:   client,
		ctx:      ctx,
		loading:  true,
		viewport: vp,
	}
}

// SetSize sets the size of the logs screen
func (m *LogsModel) SetSize(width, height int) {
	m.width = width
	m.height = height - 3 // Account for navigation bar
	
	// Update viewport size
	m.viewport.Width = width - 8  // Account for padding and borders
	m.viewport.Height = height - 10 // Account for title and controls
}

// Init implements tea.Model
func (m *LogsModel) Init() tea.Cmd {
	return m.fetchLogs()
}

// Update implements tea.Model
func (m *LogsModel) Update(msg tea.Msg) (*LogsModel, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.loading = true
			return m, m.fetchLogs()
		}
		
	case logsMsg:
		m.loading = false
		if msg.err != nil {
			m.error = fmt.Sprintf("Failed to fetch logs: %v", msg.err)
		} else {
			m.logs = msg.logs
			m.error = ""
			m.updateViewport()
		}
	}
	
	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	
	return m, cmd
}

// View implements tea.Model
func (m *LogsModel) View() string {
	if m.loading {
		return m.renderLoading()
	}
	
	if m.error != "" {
		return m.renderError()
	}
	
	return m.renderLogs()
}

// renderLoading renders the loading state
func (m *LogsModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)
	
	return style.Render("Loading logs...")
}

// renderError renders the error state
func (m *LogsModel) renderError() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("196"))
	
	return style.Render(fmt.Sprintf("Error: %s\n\nPress 'r' to refresh", m.error))
}

// renderLogs renders the logs viewer
func (m *LogsModel) renderLogs() string {
	var sections []string
	
	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Margin(0, 0, 1, 0).
		Render("📝 System Logs")
	
	sections = append(sections, title)
	
	// Viewport with logs
	sections = append(sections, m.viewport.View())
	
	// Controls
	controls := "Controls: ↑/↓ scroll | r refresh | j/k scroll line by line"
	controlsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Margin(1, 0)
	
	sections = append(sections, controlsStyle.Render(controls))
	
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(1, 2)
	
	return style.Render(content)
}

// updateViewport updates the viewport with current logs
func (m *LogsModel) updateViewport() {
	if m.logs == nil {
		return
	}
	
	// Join all log lines
	content := strings.Join(m.logs.Logs, "\n")
	
	m.viewport.SetContent(content)
	
	// Scroll to bottom to show latest logs
	m.viewport.GotoBottom()
}

// Message types
type logsMsg struct {
	logs *models.LogsResponse
	err  error
}

// fetchLogs fetches system logs
func (m *LogsModel) fetchLogs() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		logs, err := m.client.GetLogs(m.ctx)
		return logsMsg{logs: logs, err: err}
	})
}
