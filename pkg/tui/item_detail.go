package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"riven-tui/pkg/api"
	"riven-tui/pkg/models"
)

// ItemDetailModel represents the item detail screen
type ItemDetailModel struct {
	client  *api.Client
	ctx     context.Context
	width   int
	height  int
	loading bool
	error   string

	// Data
	itemID   string
	itemData map[string]interface{}
	streams  interface{}

	// UI state
	activeTab int // 0: Details, 1: Streams, 2: Actions

	// Streams table
	streamsTable table.Model

	// Auto-refresh
	lastUpdate time.Time
}

// ItemDetailMsg represents messages for the item detail screen
type itemDetailMsg struct {
	itemData map[string]interface{}
	err      error
}

type itemStreamsMsg struct {
	streams interface{}
	err     error
}

// NewItemDetailModel creates a new item detail model
func NewItemDetailModel(client *api.Client, ctx context.Context, itemID string) *ItemDetailModel {
	// Create streams table
	columns := []table.Column{
		{Title: "Title", Width: 50},
		{Title: "Quality", Width: 10},
		{Title: "Size", Width: 10},
		{Title: "Cached", Width: 8},
		{Title: "Rank", Width: 6},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &ItemDetailModel{
		client:       client,
		ctx:          ctx,
		itemID:       itemID,
		loading:      true,
		streamsTable: t,
	}
}

// SetSize sets the size of the item detail screen
func (m *ItemDetailModel) SetSize(width, height int) {
	m.width = width
	m.height = height - 3 // Account for navigation bar

	// Update streams table size
	m.streamsTable.SetHeight(m.height - 15) // Leave space for tabs and info
}

// Init implements tea.Model
func (m *ItemDetailModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchItemDetail(),
		m.fetchItemStreams(),
		m.autoRefresh(),
	)
}

// Update implements tea.Model
func (m *ItemDetailModel) Update(msg tea.Msg) (*ItemDetailModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case itemDetailMsg:
		m.loading = false
		if msg.err != nil {
			m.error = fmt.Sprintf("Failed to fetch item details: %v", msg.err)
		} else {
			m.itemData = msg.itemData
			m.error = ""
		}

	case itemStreamsMsg:
		if msg.err != nil {
			// Streams might not be available, don't show error
		} else {
			m.streams = msg.streams
			m.updateStreamsTable()
		}

	case refreshMsg:
		m.lastUpdate = time.Now()
		return m, tea.Batch(
			m.fetchItemDetail(),
			m.fetchItemStreams(),
			m.autoRefresh(),
		)

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.activeTab = (m.activeTab + 1) % 3
			return m, nil

		case "shift+tab":
			m.activeTab = (m.activeTab - 1 + 3) % 3
			return m, nil

		case "r":
			m.loading = true
			return m, tea.Batch(
				m.fetchItemDetail(),
				m.fetchItemStreams(),
			)

		case "1":
			m.activeTab = 0
			return m, nil

		case "2":
			m.activeTab = 1
			return m, nil

		case "3":
			m.activeTab = 2
			return m, nil
		}

		// Update streams table if on streams tab
		if m.activeTab == 1 {
			m.streamsTable, cmd = m.streamsTable.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m *ItemDetailModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.error != "" {
		return m.renderError()
	}

	return m.renderItemDetail()
}

// renderLoading renders the loading state
func (m *ItemDetailModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render("Loading item details...")
}

// renderError renders the error state
func (m *ItemDetailModel) renderError() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("196"))

	return style.Render(fmt.Sprintf("Error: %s\n\nPress 'r' to refresh", m.error))
}

// renderItemDetail renders the main item detail view
func (m *ItemDetailModel) renderItemDetail() string {
	var sections []string

	// Title
	title := "📺 Item Details"
	if m.itemData != nil {
		if itemTitle := getStringFromMap(m.itemData, "title", ""); itemTitle != "" {
			title = fmt.Sprintf("📺 %s", itemTitle)
		}
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Margin(0, 0, 1, 0)
	sections = append(sections, titleStyle.Render(title))

	// Tabs
	tabs := m.renderTabs()
	sections = append(sections, tabs)

	// Content based on active tab
	var content string
	switch m.activeTab {
	case 0:
		content = m.renderDetailsTab()
	case 1:
		content = m.renderStreamsTab()
	case 2:
		content = m.renderActionsTab()
	}
	sections = append(sections, content)

	// Controls
	controls := "Controls: [1-3] tabs [tab] next tab [r] refresh [esc] back"
	sections = append(sections, lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Render(controls))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderTabs renders the tab navigation
func (m *ItemDetailModel) renderTabs() string {
	var tabs []string

	tabNames := []string{"Details", "Streams", "Actions"}
	for i, name := range tabNames {
		style := lipgloss.NewStyle().
			Padding(0, 2).
			Margin(0, 1)

		if i == m.activeTab {
			style = style.
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230")).
				Bold(true)
		} else {
			style = style.
				Background(lipgloss.Color("240")).
				Foreground(lipgloss.Color("250"))
		}

		tabs = append(tabs, style.Render(fmt.Sprintf("%s (%d)", name, i+1)))
	}

	tabStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Padding(1, 0).
		Margin(0, 0, 1, 0)

	return tabStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, tabs...))
}

// renderDetailsTab renders the details tab content
func (m *ItemDetailModel) renderDetailsTab() string {
	if m.itemData == nil {
		return "No item data available."
	}

	var details []string

	// Basic info
	if id := getStringFromMap(m.itemData, "id", ""); id != "" {
		details = append(details, fmt.Sprintf("ID: %s", id))
	}
	if itemType := getStringFromMap(m.itemData, "type", ""); itemType != "" {
		details = append(details, fmt.Sprintf("Type: %s", itemType))
	}
	if state := getStringFromMap(m.itemData, "state", ""); state != "" {
		details = append(details, fmt.Sprintf("State: %s", state))
	}
	if year := getStringFromMap(m.itemData, "year", ""); year != "" {
		details = append(details, fmt.Sprintf("Year: %s", year))
	}

	// External IDs
	if tmdbId := getStringFromMap(m.itemData, "tmdb_id", ""); tmdbId != "" {
		details = append(details, fmt.Sprintf("TMDB ID: %s", tmdbId))
	}
	if tvdbId := getStringFromMap(m.itemData, "tvdb_id", ""); tvdbId != "" {
		details = append(details, fmt.Sprintf("TVDB ID: %s", tvdbId))
	}
	if imdbId := getStringFromMap(m.itemData, "imdb_id", ""); imdbId != "" {
		details = append(details, fmt.Sprintf("IMDB ID: %s", imdbId))
	}

	// Overview
	if overview := getStringFromMap(m.itemData, "overview", ""); overview != "" {
		details = append(details, "")
		details = append(details, "Overview:")
		details = append(details, overview)
	}

	content := strings.Join(details, "\n")

	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Height(m.height - 10)

	return contentStyle.Render(content)
}

// renderStreamsTab renders the streams tab content
func (m *ItemDetailModel) renderStreamsTab() string {
	if m.streams == nil {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Height(m.height - 10).
			Render("No streams data available.")
	}

	return m.streamsTable.View()
}

// renderActionsTab renders the actions tab content
func (m *ItemDetailModel) renderActionsTab() string {
	actions := []string{
		"Available Actions:",
		"",
		"• Retry item processing",
		"• Reset item state",
		"• Pause/unpause item",
		"• Remove item",
		"• Reindex item",
		"",
		"Note: Action implementation coming soon.",
	}

	content := strings.Join(actions, "\n")

	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Height(m.height - 10)

	return contentStyle.Render(content)
}

// updateStreamsTable updates the streams table with current data
func (m *ItemDetailModel) updateStreamsTable() {
	// TODO: Parse streams data and populate table
	// This would need to be implemented based on the actual streams data structure
	var rows []table.Row
	rows = append(rows, table.Row{"Sample Stream", "1080p", "2.5GB", "Yes", "100"})
	m.streamsTable.SetRows(rows)
}

// fetchItemDetail fetches item details from the API
func (m *ItemDetailModel) fetchItemDetail() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		itemData, err := m.client.GetItem(m.ctx, m.itemID, nil, models.BoolPtr(true))
		return itemDetailMsg{itemData: itemData, err: err}
	})
}

// fetchItemStreams fetches item streams from the API
func (m *ItemDetailModel) fetchItemStreams() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		if itemIDInt, err := strconv.Atoi(m.itemID); err == nil {
			streams, err := m.client.GetItemStreams(m.ctx, itemIDInt)
			return itemStreamsMsg{streams: streams, err: err}
		}
		return itemStreamsMsg{streams: nil, err: fmt.Errorf("invalid item ID")}
	})
}

// autoRefresh sets up auto-refresh
func (m *ItemDetailModel) autoRefresh() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return refreshMsg{}
	})
}
