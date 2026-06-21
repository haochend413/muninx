package menu

import (
	"github.com/haochend413/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type keyMap struct {
	NewNote    key.Binding
	Select     key.Binding
	SyncDB     key.Binding
	DeleteNote key.Binding
	Quit       key.Binding
	NavUp      key.Binding
	NavDown    key.Binding
}

var keys = keyMap{
	NewNote:    key.NewBinding(key.WithKeys("N")),
	Select:     key.NewBinding(key.WithKeys("enter")),
	SyncDB:     key.NewBinding(key.WithKeys("ctrl+q")),
	DeleteNote: key.NewBinding(key.WithKeys("ctrl+d")),
	Quit:       key.NewBinding(key.WithKeys("ctrl+c")),
	NavUp:      key.NewBinding(key.WithKeys("up")),
	NavDown:    key.NewBinding(key.WithKeys("down")),
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Let the status bar process its own expiry messages regardless of
	// what else this Update call handles below.
	m.statusbar, _ = m.statusbar.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.layout = computeLayout(msg.Width, msg.Height)
		m.input.SetWidth(m.layout.InputWidth)
		m.statusbar.SetWidth(m.layout.WindowWidth)
		m.statusbar.SetElemWidth(statusbarMainTag, m.layout.WindowWidth)
		m.UpdateTable()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, func() tea.Msg { return OpenQuitMsg{} }
		case key.Matches(msg, keys.SyncDB):
			return m, func() tea.Msg { return SyncRequestMsg{} }
		case key.Matches(msg, keys.DeleteNote):
			return m, func() tea.Msg { return DeleteNoteRequestMsg{Index: m.table.Cursor()} }
		case key.Matches(msg, keys.NewNote):
			return m, func() tea.Msg { return NewNoteRequestMsg{} }
		case key.Matches(msg, keys.Select):
			return m, func() tea.Msg { return SelectNoteMsg{Index: m.table.Cursor()} }
		case key.Matches(msg, keys.NavUp), key.Matches(msg, keys.NavDown):
			// Arrow keys browse the (possibly filtered) table; every other
			// key belongs to the input bar below.
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			return m, cmd
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			m.UpdateTable()
			m.table.SetCursor(0)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}
