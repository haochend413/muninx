package ui

import (
	"github.com/haochend413/muninx/internal/app/context"
	"github.com/haochend413/muninx/state"
)

// DistributeState restores cursor positions from saved state on startup.
func (m *Model) DistributeState(s *state.AppState) {
	if s == nil {
		return
	}

	notes := m.app.GetActiveNoteList()
	if nc, ok := s.NoteCursors[context.Default]; ok && int(nc) < len(notes) {
		m.notesTable.SetCursor(int(nc))
		m.app.GetDataMgr().SwitchActiveNoteByID(notes[int(nc)].ID)
	}
}

// CollectState gathers current cursor positions for persistence on quit.
func (m Model) CollectState() *state.State {
	s := state.DefaultState()
	s.App.NoteCursors[context.Default] = uint(m.notesTable.Cursor())
	return s
}
