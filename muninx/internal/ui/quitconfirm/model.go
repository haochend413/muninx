package quitconfirm

import tea "charm.land/bubbletea/v2"

// Messages sent to the root model.
type ConfirmMsg struct{}
type CancelMsg struct{}

type Model struct {
	layout Layout
}

func New() Model { return Model{} }

func (m Model) Init() tea.Cmd { return nil }
