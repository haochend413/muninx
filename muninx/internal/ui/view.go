package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
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
	case FindNoteView:
		// Render the background view, then composite the overlay on top of it.
		// Background views now fill the full screen, so the overlay is never clipped.
		var bg string
		switch m.findPreviousView {
		case WriteView:
			bg = m.write.RenderContent()
		default:
			bg = m.menu.RenderContent()
		}
		content = placeOverlay(m.findNote.MarginX(), m.findNote.MarginY(), m.findNote.RenderOverlay(), bg)
	default:
		content = m.menu.RenderContent()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// placeOverlay composites fg on top of bg, placing fg's top-left corner at (x, y).
func placeOverlay(x, y int, fg, bg string) string {
	fgLines := strings.Split(fg, "\n")
	bgLines := strings.Split(bg, "\n")

	result := make([]string, len(bgLines))
	copy(result, bgLines)

	for i, fgLine := range fgLines {
		row := y + i
		if row < 0 || row >= len(result) {
			continue
		}

		fgW := ansi.StringWidth(fgLine)
		bgLine := result[row]

		left := ansi.Truncate(bgLine, x, "")
		leftW := ansi.StringWidth(left)
		if leftW < x {
			left += strings.Repeat(" ", x-leftW)
		}

		right := ansi.TruncateLeft(bgLine, x+fgW, "")
		result[row] = left + fgLine + right
	}

	return strings.Join(result, "\n")
}
