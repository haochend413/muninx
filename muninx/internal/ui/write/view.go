package write

import (
	"github.com/haochend413/lipgloss/v2"
	"github.com/haochend413/muninx/internal/ui/styles"
)

func (m Model) RenderContent() string {
	l := m.layout

	var leftStyle, rightStyle lipgloss.Style
	if m.focus == FocusTextArea {
		leftStyle = styles.BaseStyle
		rightStyle = styles.FocusedStyle
	} else {
		leftStyle = styles.FocusedStyle
		rightStyle = styles.BaseStyle
	}

	leftBox := leftStyle.
		Width(l.RelatedWidth).
		Height(l.InnerHeight).
		// BorderTitle("Related Notes").
		Render(m.relatedList.View())

	rightBox := rightStyle.
		Width(l.EditorWidth).
		Height(l.InnerHeight).
		BorderTitle("Editor").
		Render(m.textArea.View())

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	help := styles.HelpStyle.Render(
		"Tab: toggle focus  •  Ctrl+S: save  •  Ctrl+X: back to menu  •  Enter (related): switch note  •  Ctrl+C: quit",
	)

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, help)
}
