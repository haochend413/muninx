package write

import (
	"github.com/haochend413/lipgloss/v2"
)

func (m Model) RenderContent() string {
	taView := m.textArea.View()
	vpView := m.relatedVp.View()

	// Colors are embedded in the viewport content by buildStyledContent, so no
	// outer color wrapper is needed here — that would conflict with inner resets.
	leftBox := lipgloss.NewStyle().
		Width(m.layout.TextAreaWidth).
		Render(taView)

	rightBox := lipgloss.NewStyle().
		Width(m.layout.RelatedWidth).
		Render(vpView)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
}
