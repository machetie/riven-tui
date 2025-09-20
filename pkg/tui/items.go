package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"riven-tui/pkg/api"
	"riven-tui/pkg/models"
)

// ItemsModel represents the media items screen
type ItemsModel struct {
	client  *api.Client
	ctx     context.Context
	width   int
	height  int
	loading bool
	error   string

	// Data
	items  *models.ItemsResponse
	states *models.StateResponse

	// UI components
	table       table.Model
	searchInput textinput.Model

	// State
	currentPage int
	pageSize    int
	searchQuery string
	filterState string
	sortOrder   models.SortOrder
	showSearch  bool

	// Auto-refresh
	lastUpdate time.Time

	// Selection and actions
	selectedItems []string
	showActions   bool
}

// ItemsMsg represents messages for the items screen
type itemsMsg struct {
	items *models.ItemsResponse
	err   error
}

type statesMsg struct {
	states *models.StateResponse
	err    error
}

// NewItemsModel creates a new items model
func NewItemsModel(client *api.Client, ctx context.Context) *ItemsModel {
	// Create search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search media items..."
	searchInput.CharLimit = 100
	searchInput.Width = 50

	// Create table
	columns := []table.Column{
		{Title: "ID", Width: 8},
		{Title: "Title", Width: 40},
		{Title: "Type", Width: 8},
		{Title: "State", Width: 15},
		{Title: "Year", Width: 6},
		{Title: "Updated", Width: 12},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
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

	return &ItemsModel{
		client:      client,
		ctx:         ctx,
		loading:     true,
		table:       t,
		searchInput: searchInput,
		currentPage: 1,
		pageSize:    50,
		sortOrder:   models.SortDateDesc,
	}
}

// SetSize sets the size of the items screen
func (m *ItemsModel) SetSize(width, height int) {
	m.width = width
	m.height = height - 3 // Account for navigation bar

	// Update table size
	m.table.SetHeight(m.height - 10) // Leave space for search and pagination

	// Update search input width
	m.searchInput.Width = min(width-20, 80)
}

// Init implements tea.Model
func (m *ItemsModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchItems(),
		m.fetchStates(),
		m.autoRefresh(),
	)
}

// Update implements tea.Model
func (m *ItemsModel) Update(msg tea.Msg) (*ItemsModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case itemsMsg:
		m.loading = false
		if msg.err != nil {
			m.error = fmt.Sprintf("Failed to fetch items: %v", msg.err)
		} else {
			m.items = msg.items
			m.error = ""
			m.updateTable()
		}

	case statesMsg:
		if msg.err != nil {
			// States are optional, don't show error
		} else {
			m.states = msg.states
		}

	case refreshMsg:
		m.lastUpdate = time.Now()
		return m, tea.Batch(
			m.fetchItems(),
			m.fetchStates(),
			m.autoRefresh(),
		)

	case tea.KeyMsg:
		// Handle search mode
		if m.showSearch {
			switch msg.String() {
			case "enter":
				m.searchQuery = m.searchInput.Value()
				m.showSearch = false
				m.currentPage = 1
				return m, m.fetchItems()
			case "esc":
				m.showSearch = false
				m.searchInput.SetValue("")
				return m, nil
			}
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd
		}

		// Handle normal navigation
		switch msg.String() {
		case "/", "ctrl+f":
			m.showSearch = true
			m.searchInput.Focus()
			return m, textinput.Blink

		case "r":
			m.loading = true
			return m, m.fetchItems()

		case "n", "right":
			if m.items != nil && m.currentPage < m.items.TotalPages {
				m.currentPage++
				m.loading = true
				return m, m.fetchItems()
			}

		case "p", "left":
			if m.currentPage > 1 {
				m.currentPage--
				m.loading = true
				return m, m.fetchItems()
			}

		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if page, err := strconv.Atoi(msg.String()); err == nil && m.items != nil && page <= m.items.TotalPages {
				m.currentPage = page
				m.loading = true
				return m, m.fetchItems()
			}

		case "a":
			m.showActions = !m.showActions
			return m, nil

		case "f":
			// Toggle filter by state
			if m.filterState == "" {
				m.filterState = "Failed"
			} else if m.filterState == "Failed" {
				m.filterState = "Completed"
			} else if m.filterState == "Completed" {
				m.filterState = "Downloaded"
			} else {
				m.filterState = ""
			}
			m.currentPage = 1
			m.loading = true
			return m, m.fetchItems()

		case "s":
			// Cycle through sort orders
			switch m.sortOrder {
			case models.SortDateDesc:
				m.sortOrder = models.SortDateAsc
			case models.SortDateAsc:
				m.sortOrder = models.SortTitleAsc
			case models.SortTitleAsc:
				m.sortOrder = models.SortTitleDesc
			case models.SortTitleDesc:
				m.sortOrder = models.SortDateDesc
			}
			m.currentPage = 1
			m.loading = true
			return m, m.fetchItems()

		case "c":
			// Clear search and filters
			m.searchQuery = ""
			m.filterState = ""
			m.sortOrder = models.SortDateDesc
			m.currentPage = 1
			m.loading = true
			return m, m.fetchItems()

		case "enter":
			// Show item details
			if m.items != nil && len(m.items.Items) > 0 {
				selectedRow := m.table.Cursor()
				if selectedRow < len(m.items.Items) {
					item := m.items.Items[selectedRow]
					if itemID := getStringFromMap(item, "id", ""); itemID != "" {
						return m, showItemDetailCmd(itemID)
					}
				}
			}
		}

		// Update table
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m *ItemsModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.error != "" {
		return m.renderError()
	}

	return m.renderItems()
}

// renderLoading renders the loading state
func (m *ItemsModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	return style.Render("Loading media items...")
}

// renderError renders the error state
func (m *ItemsModel) renderError() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("196"))

	return style.Render(fmt.Sprintf("Error: %s\n\nPress 'r' to refresh", m.error))
}

// renderItems renders the main items view
func (m *ItemsModel) renderItems() string {
	var sections []string

	// Search bar
	if m.showSearch {
		searchStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1).
			Width(m.searchInput.Width + 2)
		sections = append(sections, searchStyle.Render(m.searchInput.View()))
	} else if m.searchQuery != "" {
		searchInfo := fmt.Sprintf("Search: %s (Press '/' to search again)", m.searchQuery)
		sections = append(sections, lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(searchInfo))
	}

	// Table
	sections = append(sections, m.table.View())

	// Status and pagination info
	if m.items != nil {
		var statusParts []string

		// Add filter info
		if m.filterState != "" {
			statusParts = append(statusParts, fmt.Sprintf("Filter: %s", m.filterState))
		}

		// Add sort info
		sortName := ""
		switch m.sortOrder {
		case models.SortDateDesc:
			sortName = "Date ↓"
		case models.SortDateAsc:
			sortName = "Date ↑"
		case models.SortTitleAsc:
			sortName = "Title ↑"
		case models.SortTitleDesc:
			sortName = "Title ↓"
		}
		statusParts = append(statusParts, fmt.Sprintf("Sort: %s", sortName))

		// Add pagination info
		statusParts = append(statusParts, fmt.Sprintf("Page %d/%d", m.currentPage, m.items.TotalPages))
		statusParts = append(statusParts, fmt.Sprintf("%d items", m.items.TotalItems))

		statusInfo := strings.Join(statusParts, " | ")
		sections = append(sections, lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(statusInfo))

		// Controls info
		controlsInfo := "Controls: [/] search [f] filter [s] sort [c] clear [n/p] page [a] actions [enter] details"
		sections = append(sections, lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Render(controlsInfo))
	}

	// Actions panel
	if m.showActions {
		actions := "Actions: [r]etry [R]eset [d]elete [p]ause [u]npause [ESC] close"
		sections = append(sections, lipgloss.NewStyle().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("255")).
			Padding(0, 1).
			Render(actions))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// updateTable updates the table with current items
func (m *ItemsModel) updateTable() {
	if m.items == nil || len(m.items.Items) == 0 {
		m.table.SetRows([]table.Row{})
		return
	}

	var rows []table.Row
	for _, item := range m.items.Items {
		// Extract fields from the map
		id := getStringFromMap(item, "id", "")
		title := getStringFromMap(item, "title", "Unknown")
		itemType := getStringFromMap(item, "type", "")
		state := getStringFromMap(item, "state", "")
		year := getStringFromMap(item, "year", "")
		updatedAt := getStringFromMap(item, "updated_at", "")

		// Format updated time
		if updatedAt != "" {
			if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
				updatedAt = t.Format("2006-01-02")
			}
		}

		rows = append(rows, table.Row{
			id,
			truncateString(title, 38),
			itemType,
			state,
			year,
			updatedAt,
		})
	}

	m.table.SetRows(rows)
}

// Helper functions
func getStringFromMap(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", val)
	}
	return defaultValue
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// fetchItems fetches items from the API
func (m *ItemsModel) fetchItems() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		params := &api.ItemsParams{
			Limit: models.IntPtr(m.pageSize),
			Page:  models.IntPtr(m.currentPage),
			Sort:  &m.sortOrder,
		}

		if m.searchQuery != "" {
			params.Search = models.StringPtr(m.searchQuery)
		}

		if m.filterState != "" {
			params.States = models.StringPtr(m.filterState)
		}

		items, err := m.client.GetItems(m.ctx, params)
		return itemsMsg{items: items, err: err}
	})
}

// fetchStates fetches available states from the API
func (m *ItemsModel) fetchStates() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		states, err := m.client.GetStates(m.ctx)
		return statesMsg{states: states, err: err}
	})
}

// autoRefresh sets up auto-refresh
func (m *ItemsModel) autoRefresh() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return refreshMsg{}
	})
}

// showItemDetailMsg represents a message to show item details
type showItemDetailMsg struct {
	itemID string
}

// showItemDetailCmd creates a command to show item details
func showItemDetailCmd(itemID string) tea.Cmd {
	return func() tea.Msg {
		return showItemDetailMsg{itemID: itemID}
	}
}
