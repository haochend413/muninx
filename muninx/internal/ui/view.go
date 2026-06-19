package ui

import (
	tea "charm.land/bubbletea/v2"
)

// View dispatches to the correct sub-view based on the active ViewMode.
func (m Model) View() tea.View {
	if !m.ready {
		v := tea.NewView("Initializing...")
		v.AltScreen = true
		return v
	}

	var content string
	switch m.viewMode {
	case MenuView:
		content = m.menu.RenderContent()
	case WriteView:
		content = m.write.RenderContent()
	case QuitConfirmView:
		content = m.quitConfirm.RenderContent()
	default:
		content = m.menu.RenderContent()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
