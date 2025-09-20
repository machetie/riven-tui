package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpModel represents the help screen
type HelpModel struct {
	width  int
	height int
	keys   KeyMap
}

// NewHelpModel creates a new help model
func NewHelpModel(keys KeyMap) *HelpModel {
	return &HelpModel{
		keys: keys,
	}
}

// SetSize sets the size of the help screen
func (m *HelpModel) SetSize(width, height int) {
	m.width = width
	m.height = height - 3 // Account for navigation bar
}

// Init implements tea.Model
func (m *HelpModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *HelpModel) Update(msg tea.Msg) (*HelpModel, tea.Cmd) {
	return m, nil
}

// View implements tea.Model
func (m *HelpModel) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Margin(0, 0, 2, 0).
		Render("❓ Help & Keyboard Shortcuts")
	
	var sections []string
	
	// Keyboard shortcuts section
	shortcutsTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("33")).
		Margin(1, 0, 0, 0).
		Render("Keyboard Shortcuts")
	
	shortcutsContent := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Margin(1, 0).
		Render(
			"Navigation:\n" +
				"  ↑/k, ↓/j, ←/h, →/l  Navigate\n" +
				"  Enter               Select/Confirm\n" +
				"  Esc                 Back/Cancel\n\n" +
				"Application:\n" +
				"  d                   Dashboard\n" +
				"  m                   Media Items\n" +
				"  s                   Settings\n" +
				"  l                   Logs\n" +
				"  ?                   Help\n" +
				"  r                   Refresh\n" +
				"  q                   Quit",
		)
	
	sections = append(sections, shortcutsTitle, shortcutsContent)
	
	// Additional information
	infoTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("33")).
		Margin(2, 0, 0, 0).
		Render("Screen-Specific Controls")
	
	infoContent := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Margin(1, 0).
		Render(
			"Dashboard:\n" +
				"  • Auto-refreshes every 30 seconds\n" +
				"  • Press 'r' to manually refresh\n\n" +
				"Media Items:\n" +
				"  • 'n' - Next page\n" +
				"  • 'p' - Previous page\n" +
				"  • 'r' - Refresh items\n\n" +
				"Logs:\n" +
				"  • ↑/↓ - Scroll through logs\n" +
				"  • 'j'/'k' - Scroll line by line\n" +
				"  • 'r' - Refresh logs",
		)
	
	sections = append(sections, infoTitle, infoContent)
	
	// About section
	aboutTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("33")).
		Margin(2, 0, 0, 0).
		Render("About")
	
	aboutContent := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Margin(1, 0).
		Render(
			"Riven TUI - Terminal User Interface for Riven Media Management\n" +
				"Built with Bubble Tea framework\n" +
				"https://github.com/charmbracelet/bubbletea",
		)
	
	sections = append(sections, aboutTitle, aboutContent)
	
	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{title}, sections...)...)
	
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(1, 2)
	
	return style.Render(content)
}
