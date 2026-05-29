package findnote

import (
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/key"
	bTable "github.com/haochend413/bubbles/v2/table"
)

type keyMap struct {
	Close      key.Binding
	FocusLeft  key.Binding
	FocusRight key.Binding
	Select     key.Binding
}

var keys = keyMap{
	Close:      key.NewBinding(key.WithKeys("esc")),
	FocusLeft:  key.NewBinding(key.WithKeys("h", "left")),
	FocusRight: key.NewBinding(key.WithKeys("l", "right")),
	Select:     key.NewBinding(key.WithKeys("enter")),
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case bTable.MoveSelectMsg:
		switch m.focus {
		case FocusThreads:
			m.updateBranchesForSelectedThread()
			m.updateViewport()
		case FocusBranches:
			m.updateNotesForSelectedBranch()
			m.updateViewport()
		case FocusNotes:
			m.updateViewport()
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.layout = computeLayout(msg.Width, msg.Height)
		m.resizeComponents()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Close):
			return m, func() tea.Msg { return CloseMsg{} }
		case key.Matches(msg, keys.FocusLeft):
			m.MoveFocusLeft()
			return m, nil
		case key.Matches(msg, keys.FocusRight):
			m.MoveFocusRight()
			return m, nil
		case m.focus == FocusNotes && key.Matches(msg, keys.Select):
			note := m.selectedNote()
			if note == nil {
				return m, nil
			}
			fullNote := m.app.GetDataMgr().FindNoteByID(note.ID)
			if fullNote == nil {
				fullNote = note
			}
			capturedNote := fullNote
			return m, func() tea.Msg { return NoteSelectedMsg{Note: capturedNote} }
		}
	}

	// Forward unmatched messages to the focused table.
	var cmd tea.Cmd
	switch m.focus {
	case FocusThreads:
		m.threadsTable, cmd = m.threadsTable.Update(msg)
	case FocusBranches:
		m.branchesTable, cmd = m.branchesTable.Update(msg)
	case FocusNotes:
		m.notesTable, cmd = m.notesTable.Update(msg)
	}
	return m, cmd
}
