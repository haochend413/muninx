package menu

import (
	"github.com/haochend413/lipgloss/v2"
	"github.com/haochend413/muninx/internal/ui/styles"
)

func (m Model) RenderContent() string {
	o := HeaderOpts{
		MainTextColor: lipgloss.Color("255"),
		Width:         m.layout.WindowWidth,
		Height:        12,
	}
	header := renderHeader(styles.BaseStyle, o)

	tableBox := styles.FocusedStyle.Width(m.layout.WindowWidth).Border(lipgloss.Border{}, false, false, false, false).Render(m.table.View())

	inputBox := styles.BaseStyle.
		Width(m.layout.InputWidth).
		Render(m.input.View())

	help := styles.HelpStyle.Render(
		"N: new note  •  Enter: open note  •  j/k ↑↓: navigate  •  Ctrl+Q: sync  •  Ctrl+C: quit",
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, tableBox, inputBox, help)
}
