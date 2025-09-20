package tui

import "github.com/charmbracelet/lipgloss"

// Theme represents a color theme for the TUI
type Theme struct {
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Accent      lipgloss.Color
	Background  lipgloss.Color
	Surface     lipgloss.Color
	Text        lipgloss.Color
	TextMuted   lipgloss.Color
	Success     lipgloss.Color
	Warning     lipgloss.Color
	Error       lipgloss.Color
	Border      lipgloss.Color
	BorderFocus lipgloss.Color
}

// DefaultTheme returns the default color theme
func DefaultTheme() Theme {
	return Theme{
		Primary:     lipgloss.Color("39"),  // Blue
		Secondary:   lipgloss.Color("62"),  // Purple
		Accent:      lipgloss.Color("205"), // Pink
		Background:  lipgloss.Color("235"), // Dark gray
		Surface:     lipgloss.Color("240"), // Medium gray
		Text:        lipgloss.Color("255"), // White
		TextMuted:   lipgloss.Color("243"), // Light gray
		Success:     lipgloss.Color("46"),  // Green
		Warning:     lipgloss.Color("226"), // Yellow
		Error:       lipgloss.Color("196"), // Red
		Border:      lipgloss.Color("240"), // Medium gray
		BorderFocus: lipgloss.Color("62"),  // Purple
	}
}

// DarkTheme returns a dark color theme
func DarkTheme() Theme {
	return Theme{
		Primary:     lipgloss.Color("75"),  // Light blue
		Secondary:   lipgloss.Color("141"), // Light purple
		Accent:      lipgloss.Color("213"), // Light pink
		Background:  lipgloss.Color("232"), // Very dark gray
		Surface:     lipgloss.Color("237"), // Dark gray
		Text:        lipgloss.Color("255"), // White
		TextMuted:   lipgloss.Color("245"), // Medium gray
		Success:     lipgloss.Color("82"),  // Light green
		Warning:     lipgloss.Color("220"), // Light yellow
		Error:       lipgloss.Color("203"), // Light red
		Border:      lipgloss.Color("237"), // Dark gray
		BorderFocus: lipgloss.Color("141"), // Light purple
	}
}

// LightTheme returns a light color theme
func LightTheme() Theme {
	return Theme{
		Primary:     lipgloss.Color("21"),  // Dark blue
		Secondary:   lipgloss.Color("54"),  // Dark purple
		Accent:      lipgloss.Color("161"), // Dark pink
		Background:  lipgloss.Color("255"), // White
		Surface:     lipgloss.Color("252"), // Light gray
		Text:        lipgloss.Color("16"),  // Black
		TextMuted:   lipgloss.Color("102"), // Dark gray
		Success:     lipgloss.Color("28"),  // Dark green
		Warning:     lipgloss.Color("136"), // Dark yellow
		Error:       lipgloss.Color("124"), // Dark red
		Border:      lipgloss.Color("250"), // Light gray
		BorderFocus: lipgloss.Color("54"),  // Dark purple
	}
}

// GetTheme returns a theme by name
func GetTheme(name string) Theme {
	switch name {
	case "dark":
		return DarkTheme()
	case "light":
		return LightTheme()
	default:
		return DefaultTheme()
	}
}

// Common style functions

// TitleStyle returns a style for titles
func (t Theme) TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary).
		Margin(0, 0, 1, 0)
}

// SubtitleStyle returns a style for subtitles
func (t Theme) SubtitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Secondary).
		Margin(0, 0, 1, 0)
}

// CardStyle returns a style for cards/panels
func (t Theme) CardStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Padding(1, 2).
		Margin(0, 0, 1, 0)
}

// FocusedCardStyle returns a style for focused cards/panels
func (t Theme) FocusedCardStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus).
		Padding(1, 2).
		Margin(0, 0, 1, 0)
}

// ButtonStyle returns a style for buttons
func (t Theme) ButtonStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Surface).
		Foreground(t.Text).
		Padding(0, 2).
		Margin(0, 1)
}

// ActiveButtonStyle returns a style for active/selected buttons
func (t Theme) ActiveButtonStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Primary).
		Foreground(t.Background).
		Padding(0, 2).
		Margin(0, 1).
		Bold(true)
}

// TabStyle returns a style for tabs
func (t Theme) TabStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Surface).
		Foreground(t.Text).
		Padding(0, 2).
		Margin(0, 1)
}

// ActiveTabStyle returns a style for active tabs
func (t Theme) ActiveTabStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Secondary).
		Foreground(t.Background).
		Padding(0, 2).
		Margin(0, 1).
		Bold(true)
}

// StatusStyle returns a style for status text
func (t Theme) StatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.TextMuted)
}

// SuccessStyle returns a style for success messages
func (t Theme) SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Success).
		Bold(true)
}

// WarningStyle returns a style for warning messages
func (t Theme) WarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Warning).
		Bold(true)
}

// ErrorStyle returns a style for error messages
func (t Theme) ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Error).
		Bold(true)
}

// LoadingStyle returns a style for loading messages
func (t Theme) LoadingStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true)
}

// HelpStyle returns a style for help text
func (t Theme) HelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.TextMuted).
		Italic(true)
}

// NavBarStyle returns a style for the navigation bar
func (t Theme) NavBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Background).
		Padding(1, 0)
}

// TableHeaderStyle returns a style for table headers
func (t Theme) TableHeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(t.Border).
		BorderBottom(true).
		Bold(true).
		Foreground(t.Primary)
}

// TableSelectedStyle returns a style for selected table rows
func (t Theme) TableSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Background).
		Background(t.Secondary).
		Bold(false)
}

// ProgressBarStyle returns a style for progress bars
func (t Theme) ProgressBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Primary)
}

// SpinnerStyle returns a style for spinners
func (t Theme) SpinnerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Primary)
}

// State-specific styles

// StateStyle returns a style for different item states
func (t Theme) StateStyle(state string) lipgloss.Style {
	switch state {
	case "Completed":
		return lipgloss.NewStyle().Foreground(t.Success).Bold(true)
	case "Failed":
		return lipgloss.NewStyle().Foreground(t.Error).Bold(true)
	case "Downloaded":
		return lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
	case "Paused":
		return lipgloss.NewStyle().Foreground(t.Warning).Bold(true)
	case "Ongoing":
		return lipgloss.NewStyle().Foreground(t.Accent).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(t.TextMuted)
	}
}

// PriorityStyle returns a style for different priority levels
func (t Theme) PriorityStyle(priority string) lipgloss.Style {
	switch priority {
	case "high":
		return lipgloss.NewStyle().Foreground(t.Error).Bold(true)
	case "medium":
		return lipgloss.NewStyle().Foreground(t.Warning).Bold(true)
	case "low":
		return lipgloss.NewStyle().Foreground(t.Success).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(t.TextMuted)
	}
}

// QualityStyle returns a style for different quality levels
func (t Theme) QualityStyle(quality string) lipgloss.Style {
	switch quality {
	case "2160p", "4K":
		return lipgloss.NewStyle().Foreground(t.Accent).Bold(true)
	case "1080p":
		return lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
	case "720p":
		return lipgloss.NewStyle().Foreground(t.Secondary).Bold(true)
	case "480p":
		return lipgloss.NewStyle().Foreground(t.Warning).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(t.TextMuted)
	}
}

// CachedStyle returns a style for cached/uncached status
func (t Theme) CachedStyle(cached bool) lipgloss.Style {
	if cached {
		return lipgloss.NewStyle().Foreground(t.Success).Bold(true)
	}
	return lipgloss.NewStyle().Foreground(t.TextMuted)
}
