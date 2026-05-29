package quitconfirm

import (
	"github.com/haochend413/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type keyMap struct {
	Confirm key.Binding
	Reject  key.Binding
}

var keys = keyMap{
	Confirm: key.NewBinding(key.WithKeys("y")),
	Reject:  key.NewBinding(key.WithKeys("n", "esc")),
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.layout = computeLayout(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Confirm):
			return m, func() tea.Msg { return ConfirmMsg{} }
		case key.Matches(msg, keys.Reject):
			return m, func() tea.Msg { return CancelMsg{} }
		}
	}
	return m, nil
}
