package findnote

import (
	"github.com/haochend413/lipgloss/v2"
	"github.com/haochend413/muninx/internal/ui/styles"
)

func (m Model) RenderOverlay() string {
	l := m.layout

	var threadStyle, branchStyle, notesStyle lipgloss.Style
	switch m.focus {
	case FocusThreads:
		threadStyle = styles.FocusedStyle
		branchStyle = styles.BaseStyle
		notesStyle = styles.BaseStyle
	case FocusBranches:
		threadStyle = styles.BaseStyle
		branchStyle = styles.FocusedStyle
		notesStyle = styles.BaseStyle
	case FocusNotes:
		threadStyle = styles.BaseStyle
		branchStyle = styles.BaseStyle
		notesStyle = styles.FocusedStyle
	}

	// Width includes border(2) + padding(2) so content fits without word-wrap.
	tableBoxW := l.TableInnerWidth + 4
	vpBoxW := l.ViewportWidth + 4

	threadBox := threadStyle.
		Width(tableBoxW).
		Height(l.TableHeight).
		BorderTitle("Threads").
		Render(m.threadsTable.View())

	branchBox := branchStyle.
		Width(tableBoxW).
		Height(l.TableHeight).
		BorderTitle("Branches").
		Render(m.branchesTable.View())

	notesBox := notesStyle.
		Width(tableBoxW).
		Height(l.TableHeight).
		BorderTitle("Notes").
		Render(m.notesTable.View())

	vpBox := styles.BaseStyle.
		Width(vpBoxW).
		Height(l.TableHeight).
		BorderTitle("Content").
		Render(m.vp.View())

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, threadBox, branchBox, notesBox, vpBox)

	help := styles.HelpStyle.Render(
		"h/l ←→: switch column  •  j/k ↑↓: navigate  •  Enter (Notes): open note  •  Esc: close",
	)

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, help)
}
