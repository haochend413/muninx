package menu

import (
	"github.com/haochend413/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type keyMap struct {
	NewNote  key.Binding
	Select   key.Binding
	SyncDB   key.Binding
	FindNote key.Binding
	Quit     key.Binding
}

var keys = keyMap{
	NewNote:  key.NewBinding(key.WithKeys("N")),
	Select:   key.NewBinding(key.WithKeys("enter")),
	SyncDB:   key.NewBinding(key.WithKeys("ctrl+q")),
	FindNote: key.NewBinding(key.WithKeys("ctrl+f")),
	Quit:     key.NewBinding(key.WithKeys("ctrl+c")),
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.layout = computeLayout(msg.Width, msg.Height)
		m.input.SetWidth(m.layout.InputWidth)
		m.UpdateTable()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, func() tea.Msg { return OpenQuitMsg{} }
		case key.Matches(msg, keys.SyncDB):
			return m, func() tea.Msg { return SyncRequestMsg{} }
		case key.Matches(msg, keys.FindNote):
			return m, func() tea.Msg { return OpenFindNoteMsg{} }
		case key.Matches(msg, keys.NewNote):
			return m, func() tea.Msg { return NewNoteRequestMsg{} }
		case key.Matches(msg, keys.Select):
			return m, func() tea.Msg { return SelectNoteMsg{Index: m.table.Cursor()} }
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}
