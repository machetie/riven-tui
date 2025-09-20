package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"riven-tui/pkg/api"
	"riven-tui/pkg/models"
)

// EventsModel represents the real-time events screen
type EventsModel struct {
	client   *api.Client
	ctx      context.Context
	width    int
	height   int
	loading  bool
	error    string

	// Events data
	events      []models.Event
	eventsTable table.Model
	maxEvents   int

	// Event streaming
	eventChan   chan models.Event
	streaming   bool
	
	// Auto-refresh
	lastUpdate time.Time
}

// eventMsg represents messages for the events screen
type eventMsg struct {
	event models.Event
}

type eventsErrorMsg struct {
	err error
}

// NewEventsModel creates a new events model
func NewEventsModel(client *api.Client, ctx context.Context) *EventsModel {
	// Create events table
	columns := []table.Column{
		{Title: "Time", Width: 20},
		{Title: "Type", Width: 15},
		{Title: "Message", Width: 50},
		{Title: "Data", Width: 30},
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

	return &EventsModel{
		client:      client,
		ctx:         ctx,
		loading:     false,
		eventsTable: t,
		maxEvents:   100, // Keep last 100 events
		eventChan:   make(chan models.Event, 50),
	}
}

// SetSize sets the size of the events screen
func (m *EventsModel) SetSize(width, height int) {
	m.width = width
	m.height = height - 3 // Account for navigation bar
	
	// Update events table size
	m.eventsTable.SetHeight(m.height - 8) // Leave space for title and controls
}

// Init implements tea.Model
func (m *EventsModel) Init() tea.Cmd {
	return tea.Batch(
		m.startEventStreaming(),
		m.listenForEvents(),
	)
}

// Update implements tea.Model
func (m *EventsModel) Update(msg tea.Msg) (*EventsModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case eventMsg:
		m.addEvent(msg.event)
		m.updateEventsTable()

	case eventsErrorMsg:
		m.error = fmt.Sprintf("Event streaming error: %v", msg.err)
		m.streaming = false
		// Try to restart streaming after a delay
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			// Restart event streaming
			m.error = ""
			return m, tea.Batch(
				m.startEventStreaming(),
				m.listenForEvents(),
			)

		case "c":
			// Clear events
			m.events = []models.Event{}
			m.updateEventsTable()
			return m, nil

		case "s":
			// Toggle streaming
			if m.streaming {
				m.streaming = false
			} else {
				return m, tea.Batch(
					m.startEventStreaming(),
					m.listenForEvents(),
				)
			}
		}

		// Update events table
		m.eventsTable, cmd = m.eventsTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m *EventsModel) View() string {
	if m.error != "" {
		return m.renderError()
	}

	return m.renderEvents()
}

// renderError renders the error state
func (m *EventsModel) renderError() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("196"))

	return style.Render(fmt.Sprintf("Error: %s\n\nPress 'r' to retry", m.error))
}

// renderEvents renders the main events view
func (m *EventsModel) renderEvents() string {
	var sections []string

	// Title with streaming status
	streamingStatus := "🔴 Disconnected"
	if m.streaming {
		streamingStatus = "🟢 Live"
	}
	
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Margin(0, 0, 1, 0).
		Render(fmt.Sprintf("📡 Real-time Events %s", streamingStatus))
	sections = append(sections, title)

	// Events count and controls
	eventCount := fmt.Sprintf("Events: %d/%d", len(m.events), m.maxEvents)
	controls := "Controls: [r] restart [c] clear [s] toggle streaming [↑/↓] navigate"
	
	info := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(eventCount),
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(0, 0, 0, 4).Render(controls),
	)
	sections = append(sections, info)

	// Events table
	sections = append(sections, m.eventsTable.View())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// addEvent adds a new event to the list
func (m *EventsModel) addEvent(event models.Event) {
	// Add to beginning of slice (newest first)
	m.events = append([]models.Event{event}, m.events...)
	
	// Trim to max events
	if len(m.events) > m.maxEvents {
		m.events = m.events[:m.maxEvents]
	}
	
	m.lastUpdate = time.Now()
}

// updateEventsTable updates the events table with current data
func (m *EventsModel) updateEventsTable() {
	var rows []table.Row
	
	for _, event := range m.events {
		// Parse timestamp
		timestamp := event.Timestamp
		if t, err := time.Parse(time.RFC3339, event.Timestamp); err == nil {
			timestamp = t.Format("15:04:05")
		}
		
		// Format event type
		eventType := event.Type
		if len(eventType) > 12 {
			eventType = eventType[:12] + "..."
		}
		
		// Format message
		message := event.Message
		if message == "" {
			message = "No message"
		}
		if len(message) > 47 {
			message = message[:47] + "..."
		}
		
		// Format data
		dataStr := ""
		if event.Data != nil && len(event.Data) > 0 {
			var dataParts []string
			for key, value := range event.Data {
				dataParts = append(dataParts, fmt.Sprintf("%s:%v", key, value))
			}
			dataStr = strings.Join(dataParts, ", ")
			if len(dataStr) > 27 {
				dataStr = dataStr[:27] + "..."
			}
		}
		
		rows = append(rows, table.Row{timestamp, eventType, message, dataStr})
	}
	
	m.eventsTable.SetRows(rows)
}

// startEventStreaming starts the event streaming goroutine
func (m *EventsModel) startEventStreaming() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		m.streaming = true
		
		// Start streaming in a goroutine
		go func() {
			err := m.client.StreamEventsSSE(m.ctx, m.eventChan)
			if err != nil {
				// Send error message
				select {
				case <-m.ctx.Done():
					return
				default:
					// This would need to be handled differently in a real implementation
					// For now, we'll just mark streaming as false
					m.streaming = false
				}
			}
		}()
		
		return nil
	})
}

// listenForEvents listens for events from the event channel
func (m *EventsModel) listenForEvents() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		select {
		case event := <-m.eventChan:
			return eventMsg{event: event}
		case <-m.ctx.Done():
			return nil
		case <-time.After(100 * time.Millisecond):
			// Timeout to prevent blocking
			return m.listenForEvents()
		}
	})
}
