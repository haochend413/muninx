package ui

import (
	"fmt"
	"time"

	"github.com/haochend413/muninx/internal/app/context"
	"github.com/haochend413/muninx/state"
)

// DistributeState restores cursor positions from saved state on startup.
func (m *Model) DistributeState(s *state.AppState) {
	if s == nil {
		return
	}

	threads := m.app.GetThreadList()
	if tc, ok := s.ThreadCursors[context.Default]; ok && int(tc) < len(threads) {
		m.threadsTable.SetCursor(int(tc))
		m.app.GetDataMgr().SwitchActiveThreadByID(threads[int(tc)].ID)
		m.updateBranchesTable()
	}

	branches := m.app.GetActiveBranchList()
	if bc, ok := s.BranchCursors[context.Default]; ok && int(bc) < len(branches) {
		m.branchesTable.SetCursor(int(bc))
		m.app.GetDataMgr().SwitchActiveBranchByID(branches[int(bc)].ID)
		m.updateNotesTable()
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
	s.App.ThreadCursors[context.Default] = uint(m.threadsTable.Cursor())
	s.App.BranchCursors[context.Default] = uint(m.branchesTable.Cursor())
	s.App.NoteCursors[context.Default] = uint(m.notesTable.Cursor())
	return s
}

func HelpText() string {
	return `muninx — note management tool

Views:
  Menu View  : Recent notes table. N=new note, Enter=open note, j/k=navigate.
  Write View : Left textarea (vim-like), right related notes. Tab=toggle focus,
               Ctrl+S=save, Ctrl+X=back, Enter (related)=switch note.
  Quit View  : y=quit+sync, n/Esc=cancel.

Global: Ctrl+Q=sync database, Ctrl+C=quit confirmation.`
}

func formatTimeAgo(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	d := time.Since(t)
	if d < time.Second {
		return "just now"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds ago", int(d.Minutes()), int(d.Seconds())-60*int(d.Minutes()))
	}
	if d < time.Hour*24 {
		return fmt.Sprintf("%dh%dm ago", int(d.Hours()), int(d.Minutes())-60*int(d.Hours()))
	}
	days := int(d.Hours() / 24)
	if days < 7 {
		return fmt.Sprintf("%dd%dh ago", days, int(d.Hours())-24*days)
	}
	return t.Format("01-02 15:04")
}

var _ = formatTimeAgo
