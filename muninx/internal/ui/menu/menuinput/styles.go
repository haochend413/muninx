package menuinput

import (
	"image/color"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/lipgloss/v2"
)

// DefaultStyles returns the default styles for focused and blurred states for
// the textarea.
func DefaultStyles(isDark bool) Styles {
	lightDark := lipgloss.LightDark(isDark)

	var s Styles
	s.Focused = StyleState{
		Placeholder: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Suggestion:  lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Prompt:      lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		Text:        lipgloss.NewStyle(),
	}
	s.Blurred = StyleState{
		Placeholder: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Suggestion:  lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Prompt:      lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		Text:        lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("245"), lipgloss.Color("7"))),
	}
	s.Cursor = CursorStyle{
		Color: lipgloss.Color("7"),
		Shape: tea.CursorBlock,
		Blink: true,
	}
	return s
}

// DefaultLightStyles returns the default styles for a light background.
func DefaultLightStyles() Styles {
	return DefaultStyles(false)
}

// DefaultDarkStyles returns the default styles for a dark background.
func DefaultDarkStyles() Styles {
	return DefaultStyles(true)
}

// Styles are the styles for the textarea, separated into focused and blurred
// states.
type Styles struct {
	Focused StyleState
	Blurred StyleState
	Cursor  CursorStyle
}

// StyleState holds the styles for a single focus state.
type StyleState struct {
	Text        lipgloss.Style
	Placeholder lipgloss.Style
	Suggestion  lipgloss.Style
	Prompt      lipgloss.Style
}

// CursorStyle is the style for real and virtual cursors.
type CursorStyle struct {
	Color      color.Color
	Shape      tea.CursorShape
	Blink      bool
	BlinkSpeed time.Duration
}
