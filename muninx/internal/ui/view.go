package ui

import (
	tea "charm.land/bubbletea/v2"
)

// View dispatches to the correct sub-view based on the active ViewMode.
func (m Model) View() tea.View {
	if !m.ready {
		v := tea.NewView("Initializing...")
		v.AltScreen = true
		v.MouseMode = tea.MouseModeCellMotion
		return v
	}

	var content string
	switch m.viewMode {
	case MenuView:
		content = m.menu.RenderContent()
	case WriteView:
		content = m.write.RenderContent()
	default:
		content = m.menu.RenderContent()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	// Capture the scroll wheel ourselves so it never falls through to the
	// terminal's own scrollback while we're in the alt screen.
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
