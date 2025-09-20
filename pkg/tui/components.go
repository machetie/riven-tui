package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoadingComponent represents a loading spinner with message
type LoadingComponent struct {
	spinner spinner.Model
	message string
	theme   Theme
}

// NewLoadingComponent creates a new loading component
func NewLoadingComponent(message string, theme Theme) LoadingComponent {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = theme.SpinnerStyle()

	return LoadingComponent{
		spinner: s,
		message: message,
		theme:   theme,
	}
}

// Init implements tea.Model
func (l LoadingComponent) Init() tea.Cmd {
	return l.spinner.Tick
}

// Update implements tea.Model
func (l LoadingComponent) Update(msg tea.Msg) (LoadingComponent, tea.Cmd) {
	var cmd tea.Cmd
	l.spinner, cmd = l.spinner.Update(msg)
	return l, cmd
}

// View implements tea.Model
func (l LoadingComponent) View() string {
	return fmt.Sprintf("%s %s", l.spinner.View(), l.message)
}

// ProgressComponent represents a progress bar with status
type ProgressComponent struct {
	progress progress.Model
	current  float64
	total    float64
	message  string
	theme    Theme
}

// NewProgressComponent creates a new progress component
func NewProgressComponent(message string, theme Theme) ProgressComponent {
	p := progress.New(progress.WithDefaultGradient())
	p.Width = 40

	return ProgressComponent{
		progress: p,
		message:  message,
		theme:    theme,
	}
}

// SetProgress updates the progress value
func (p *ProgressComponent) SetProgress(current, total float64) {
	p.current = current
	p.total = total
}

// View implements tea.Model
func (p ProgressComponent) View() string {
	percent := 0.0
	if p.total > 0 {
		percent = p.current / p.total
	}

	progressBar := p.progress.ViewAs(percent)
	status := fmt.Sprintf("%.0f/%.0f (%.1f%%)", p.current, p.total, percent*100)

	return fmt.Sprintf("%s\n%s %s", p.message, progressBar, status)
}

// ErrorComponent represents an error display with retry option
type ErrorComponent struct {
	error   string
	details string
	theme   Theme
	width   int
	height  int
}

// NewErrorComponent creates a new error component
func NewErrorComponent(err string, details string, theme Theme) ErrorComponent {
	return ErrorComponent{
		error:   err,
		details: details,
		theme:   theme,
	}
}

// SetSize sets the size of the error component
func (e *ErrorComponent) SetSize(width, height int) {
	e.width = width
	e.height = height
}

// View implements tea.Model
func (e ErrorComponent) View() string {
	title := e.theme.ErrorStyle().Render("⚠️ Error")
	
	errorMsg := e.theme.ErrorStyle().Render(e.error)
	
	var content []string
	content = append(content, title)
	content = append(content, "")
	content = append(content, errorMsg)
	
	if e.details != "" {
		content = append(content, "")
		content = append(content, e.theme.StatusStyle().Render("Details:"))
		content = append(content, e.theme.StatusStyle().Render(e.details))
	}
	
	content = append(content, "")
	content = append(content, e.theme.HelpStyle().Render("Press 'r' to retry or 'esc' to go back"))

	contentStr := strings.Join(content, "\n")

	style := lipgloss.NewStyle().
		Width(e.width).
		Height(e.height).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(e.theme.Error).
		Padding(2, 4)

	return style.Render(contentStr)
}

// StatusComponent represents a status message with auto-dismiss
type StatusComponent struct {
	message   string
	statusType StatusType
	theme     Theme
	createdAt time.Time
	duration  time.Duration
}

// StatusType represents different types of status messages
type StatusType int

const (
	StatusInfo StatusType = iota
	StatusSuccess
	StatusWarning
	StatusError
)

// NewStatusComponent creates a new status component
func NewStatusComponent(message string, statusType StatusType, theme Theme, duration time.Duration) StatusComponent {
	return StatusComponent{
		message:    message,
		statusType: statusType,
		theme:      theme,
		createdAt:  time.Now(),
		duration:   duration,
	}
}

// IsExpired checks if the status message has expired
func (s StatusComponent) IsExpired() bool {
	return time.Since(s.createdAt) > s.duration
}

// View implements tea.Model
func (s StatusComponent) View() string {
	var style lipgloss.Style
	var icon string

	switch s.statusType {
	case StatusSuccess:
		style = s.theme.SuccessStyle()
		icon = "✅"
	case StatusWarning:
		style = s.theme.WarningStyle()
		icon = "⚠️"
	case StatusError:
		style = s.theme.ErrorStyle()
		icon = "❌"
	default:
		style = s.theme.StatusStyle()
		icon = "ℹ️"
	}

	return style.Render(fmt.Sprintf("%s %s", icon, s.message))
}

// ConfirmationComponent represents a confirmation dialog
type ConfirmationComponent struct {
	title     string
	message   string
	confirmed bool
	theme     Theme
	width     int
	height    int
}

// NewConfirmationComponent creates a new confirmation component
func NewConfirmationComponent(title, message string, theme Theme) ConfirmationComponent {
	return ConfirmationComponent{
		title:   title,
		message: message,
		theme:   theme,
	}
}

// SetSize sets the size of the confirmation component
func (c *ConfirmationComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// Update handles key events for the confirmation dialog
func (c *ConfirmationComponent) Update(msg tea.KeyMsg) (bool, bool) { // returns (confirmed, dismissed)
	switch msg.String() {
	case "y", "Y", "enter":
		return true, true
	case "n", "N", "esc":
		return false, true
	}
	return false, false
}

// View implements tea.Model
func (c ConfirmationComponent) View() string {
	title := c.theme.TitleStyle().Render(c.title)
	message := c.theme.StatusStyle().Render(c.message)
	
	yesButton := c.theme.ActiveButtonStyle().Render("Yes (y)")
	noButton := c.theme.ButtonStyle().Render("No (n)")
	buttons := lipgloss.JoinHorizontal(lipgloss.Center, yesButton, noButton)
	
	content := lipgloss.JoinVertical(lipgloss.Center, title, "", message, "", buttons)

	style := lipgloss.NewStyle().
		Width(c.width).
		Height(c.height).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c.theme.BorderFocus).
		Padding(2, 4)

	return style.Render(content)
}

// ToastComponent represents a toast notification
type ToastComponent struct {
	message    string
	statusType StatusType
	theme      Theme
	createdAt  time.Time
	duration   time.Duration
	position   ToastPosition
}

// ToastPosition represents the position of toast notifications
type ToastPosition int

const (
	ToastTopRight ToastPosition = iota
	ToastTopLeft
	ToastBottomRight
	ToastBottomLeft
	ToastCenter
)

// NewToastComponent creates a new toast component
func NewToastComponent(message string, statusType StatusType, theme Theme, duration time.Duration, position ToastPosition) ToastComponent {
	return ToastComponent{
		message:    message,
		statusType: statusType,
		theme:      theme,
		createdAt:  time.Now(),
		duration:   duration,
		position:   position,
	}
}

// IsExpired checks if the toast has expired
func (t ToastComponent) IsExpired() bool {
	return time.Since(t.createdAt) > t.duration
}

// View implements tea.Model
func (t ToastComponent) View() string {
	var style lipgloss.Style
	var icon string

	switch t.statusType {
	case StatusSuccess:
		style = lipgloss.NewStyle().
			Background(t.theme.Success).
			Foreground(t.theme.Background)
		icon = "✅"
	case StatusWarning:
		style = lipgloss.NewStyle().
			Background(t.theme.Warning).
			Foreground(t.theme.Background)
		icon = "⚠️"
	case StatusError:
		style = lipgloss.NewStyle().
			Background(t.theme.Error).
			Foreground(t.theme.Background)
		icon = "❌"
	default:
		style = lipgloss.NewStyle().
			Background(t.theme.Primary).
			Foreground(t.theme.Background)
		icon = "ℹ️"
	}

	style = style.
		Padding(0, 2).
		Margin(0, 1).
		Border(lipgloss.RoundedBorder()).
		Bold(true)

	return style.Render(fmt.Sprintf("%s %s", icon, t.message))
}

// HelpComponent represents a help overlay
type HelpComponent struct {
	title    string
	bindings []KeyBinding
	theme    Theme
	width    int
	height   int
}

// KeyBinding represents a key binding with description
type KeyBinding struct {
	Key         string
	Description string
}

// NewHelpComponent creates a new help component
func NewHelpComponent(title string, bindings []KeyBinding, theme Theme) HelpComponent {
	return HelpComponent{
		title:    title,
		bindings: bindings,
		theme:    theme,
	}
}

// SetSize sets the size of the help component
func (h *HelpComponent) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// View implements tea.Model
func (h HelpComponent) View() string {
	title := h.theme.TitleStyle().Render(h.title)
	
	var bindings []string
	for _, binding := range h.bindings {
		key := h.theme.ActiveButtonStyle().Render(binding.Key)
		desc := h.theme.StatusStyle().Render(binding.Description)
		bindings = append(bindings, fmt.Sprintf("%s  %s", key, desc))
	}
	
	bindingsStr := strings.Join(bindings, "\n")
	
	footer := h.theme.HelpStyle().Render("Press 'esc' or '?' to close help")
	
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", bindingsStr, "", footer)

	style := lipgloss.NewStyle().
		Width(h.width).
		Height(h.height).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(h.theme.BorderFocus).
		Padding(2, 4)

	return style.Render(content)
}
