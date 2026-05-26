package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/haochend413/lipgloss/v2"
	"github.com/haochend413/muninx/internal/ui/openmenu"
	"github.com/haochend413/muninx/internal/ui/styles"
)

const (
	findOverlayMarginX = 4 // terminal columns reserved on each horizontal side
	findOverlayMarginY = 3 // terminal rows reserved on each vertical side
)

// View dispatches to the correct sub-view based on the active ViewMode.
func (m Model) View() tea.View {
	if !m.ready {
		v := tea.NewView("Initializing...")
		v.AltScreen = true
		return v
	}
	switch m.viewMode {
	case MenuView:
		return m.menuView()
	case WriteView:
		return m.writeView()
	case QuitConfirmView:
		return m.quitView()
	case FindNoteView:
		return m.findView()
	default:
		return m.menuView()
	}
}

// menuView renders: ASCII logo + recent-notes table + text input.
func (m Model) menuView() tea.View {
	o := openmenu.Opts{
		MainTextColor: lipgloss.Color("255"),
		Width:         m.width,
		Height:        12,
	}
	header := openmenu.Render(styles.BaseStyle, o)

	tableBox := styles.FocusedStyle.Render(m.menuTable.View())

	inputBox := styles.BaseStyle.
		Width(m.width - 6).
		Render(m.menuInput.View())

	help := styles.HelpStyle.Render(
		"N: new note  •  Enter: open note  •  j/k ↑↓: navigate  •  Ctrl+Q: sync  •  Ctrl+C: quit",
	)

	content := lipgloss.JoinVertical(lipgloss.Left, header, tableBox, inputBox, help)
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// writeView renders: related-notes list (left) + textarea (right).
func (m Model) writeView() tea.View {
	relatedW := (m.width * 40) / 100
	editorW := m.width - relatedW

	var leftStyle, rightStyle lipgloss.Style
	if m.writeFocus == WriteFocusTextArea {
		leftStyle = styles.BaseStyle
		rightStyle = styles.FocusedStyle
	} else {
		leftStyle = styles.FocusedStyle
		rightStyle = styles.BaseStyle
	}

	innerH := m.height - 6
	if innerH < 3 {
		innerH = 3
	}

	leftBox := leftStyle.
		Width(relatedW - 4).
		Height(innerH).
		BorderTitle("Related Notes").
		Render(m.relatedList.View())

	rightBox := rightStyle.
		Width(editorW - 4).
		Height(innerH).
		BorderTitle("Editor").
		Render(m.textArea.View())

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	help := styles.HelpStyle.Render(
		"Tab: toggle focus  •  Ctrl+S: save  •  Ctrl+X: back to menu  •  Enter (related): switch note  •  Ctrl+C: quit",
	)

	content := lipgloss.JoinVertical(lipgloss.Left, mainContent, help)
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// quitView renders the quit-confirmation prompt.
func (m Model) quitView() tea.View {
	msg := "Quit muninx? Unsaved changes will be synced first.\n\n  y  →  save + quit\n  n / Esc  →  cancel"
	v := tea.NewView(styles.BaseStyle.Render(msg))
	v.AltScreen = true
	return v
}

// findView renders the four-column FindNote overlay on top of the background view.
func (m Model) findView() tea.View {
	overlayW := m.width - 2*findOverlayMarginX
	overlayH := m.height - 2*findOverlayMarginY
	if overlayW < 60 {
		overlayW = 60
	}
	if overlayH < 10 {
		overlayH = 10
	}

	innerH := overlayH - 4
	if innerH < 3 {
		innerH = 3
	}

	// rendered box width = inner + border(2) + padding(2)
	tableBoxW := findTableInnerW + 4

	var threadStyle, branchStyle, notesStyle lipgloss.Style
	switch m.findFocus {
	case FindFocusThreads:
		threadStyle = styles.FocusedStyle
		branchStyle = styles.BaseStyle
		notesStyle = styles.BaseStyle
	case FindFocusBranches:
		threadStyle = styles.BaseStyle
		branchStyle = styles.FocusedStyle
		notesStyle = styles.BaseStyle
	case FindFocusNotes:
		threadStyle = styles.BaseStyle
		branchStyle = styles.BaseStyle
		notesStyle = styles.FocusedStyle
	}

	vpInnerW := overlayW - 3*tableBoxW - 4
	if vpInnerW < 10 {
		vpInnerW = 10
	}

	threadBox := threadStyle.
		Width(findTableInnerW).
		Height(innerH).
		BorderTitle("Threads").
		Render(m.findThreadsTable.View())

	branchBox := branchStyle.
		Width(findTableInnerW).
		Height(innerH).
		BorderTitle("Branches").
		Render(m.findBranchesTable.View())

	notesBox := notesStyle.
		Width(findTableInnerW).
		Height(innerH).
		BorderTitle("Notes").
		Render(m.findNotesTable.View())

	vpBox := styles.BaseStyle.
		Width(vpInnerW).
		Height(innerH).
		BorderTitle("Content").
		Render(m.findViewport.View())

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, threadBox, branchBox, notesBox, vpBox)

	help := styles.HelpStyle.Render(
		"h/l ←→: switch column  •  j/k ↑↓: navigate  •  Enter (Notes): open note  •  Esc: close",
	)

	overlayPanel := lipgloss.JoinVertical(lipgloss.Left, mainContent, help)

	// Render background and composite overlay on top.
	var bgContent string
	switch m.findPreviousView {
	case WriteView:
		bgContent = m.writeView().Content
	default:
		bgContent = m.menuView().Content
	}

	composed := placeOverlay(findOverlayMarginX, findOverlayMarginY, overlayPanel, bgContent)
	v := tea.NewView(composed)
	v.AltScreen = true
	return v
}

// placeOverlay composites fg on top of bg, placing fg's top-left at (x, y).
// Both strings are multi-line ANSI-coded terminal output.
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

		// Left part: visible characters before column x.
		left := ansi.Truncate(bgLine, x, "")
		leftW := ansi.StringWidth(left)
		if leftW < x {
			left += strings.Repeat(" ", x-leftW)
		}

		// Right part: visible characters after column x+fgW.
		right := ansi.TruncateLeft(bgLine, x+fgW, "")

		result[row] = left + fgLine + right
	}

	return strings.Join(result, "\n")
}
