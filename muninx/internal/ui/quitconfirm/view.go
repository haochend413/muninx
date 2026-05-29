package quitconfirm

import "github.com/haochend413/muninx/internal/ui/styles"

func (m Model) RenderContent() string {
	l := m.layout
	msg := "Quit muninx? Unsaved changes will be synced first.\n\n  y  →  save + quit\n  n / Esc  →  cancel"
	return styles.BaseStyle.
		Width(l.WindowWidth).
		Height(l.WindowHeight).
		Render(msg)
}
