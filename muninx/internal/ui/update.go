package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/ui/findnote"
	"github.com/haochend413/muninx/internal/ui/menu"
	"github.com/haochend413/muninx/internal/ui/quitconfirm"
	"github.com/haochend413/muninx/internal/ui/write"
	statePkg "github.com/haochend413/muninx/state"
	"github.com/haochend413/muninx/sys"
)

// syncDoneMsg is sent when a background sync completes normally.
type syncDoneMsg struct{}

// quitSyncDoneMsg is sent when the quit-path sync completes.
type quitSyncDoneMsg struct{}

// syncCmd runs SyncWithDatabase in a goroutine so it never blocks the event loop.
func syncCmd(a *app.App) tea.Cmd {
	return func() tea.Msg {
		a.SyncWithDatabase()
		return syncDoneMsg{}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global messages handled before view-mode dispatch.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		var wc1, wc2, wc3, wc4 tea.Cmd
		m.menu, wc1 = m.menu.Update(msg)
		m.write, wc2 = m.write.Update(msg)
		m.findNote, wc3 = m.findNote.Update(msg)
		m.quitConfirm, wc4 = m.quitConfirm.Update(msg)
		return m, tea.Batch(wc1, wc2, wc3, wc4)

	case tickMsg:
		return m, tick()

	// --- Async sync completion ---

	case syncDoneMsg:
		m.menu.UpdateTable()
		return m, nil

	case quitSyncDoneMsg:
		return m, tea.Quit

	// --- Messages from menu sub-model ---

	case menu.SelectNoteMsg:
		notes := m.app.GetDataMgr().GetAllNotesByIDDesc()
		if msg.Index >= 0 && msg.Index < len(notes) {
			cmd := m.loadNoteIntoEditor(notes[msg.Index])
			return m, cmd
		}
		return m, nil

	case menu.NewNoteRequestMsg:
		return m.handleNewNote()

	case menu.SyncRequestMsg:
		return m, syncCmd(m.app)

	case menu.OpenFindNoteMsg:
		m.openFindOverlay()
		return m, nil

	case menu.OpenQuitMsg:
		m.previousViewMode = m.viewMode
		m.viewMode = QuitConfirmView
		return m, nil

	// --- Messages from write sub-model ---

	case write.BackToMenuMsg:
		m.viewMode = MenuView
		m.menu.UpdateTable()
		switchToEnglish()
		return m, nil

	case write.OpenFindNoteMsg:
		m.openFindOverlay()
		return m, nil

	case write.OpenQuitMsg:
		m.previousViewMode = m.viewMode
		m.viewMode = QuitConfirmView
		return m, nil

	case write.SyncRequestMsg:
		return m, syncCmd(m.app)

	case write.OpenNoteMsg:
		cmd := m.loadNoteIntoEditor(msg.Note)
		return m, cmd

	// --- Messages from findnote sub-model ---

	case findnote.NoteSelectedMsg:
		m.write.SaveCurrentNote()
		cmd := m.loadNoteIntoEditor(msg.Note)
		return m, cmd

	case findnote.CloseMsg:
		m.viewMode = m.findPreviousView
		return m, nil

	// --- Messages from quitconfirm sub-model ---

	case quitconfirm.ConfirmMsg:
		// Save note and collect state synchronously (fast, no I/O), then sync DB in background.
		m.write.SaveCurrentNote()
		s := m.CollectState()
		a := m.app
		cfg := m.Config
		return m, func() tea.Msg {
			a.SyncWithDatabase()
			if s != nil {
				if err := statePkg.SaveState(cfg.StateFilePath, s); err != nil {
					sys.LogError(fmt.Errorf("error saving state: %v", err))
				}
			}
			return quitSyncDoneMsg{}
		}

	case quitconfirm.CancelMsg:
		m.viewMode = m.previousViewMode
		return m, nil
	}

	// Delegate unhandled messages to the active sub-model.
	var cmd tea.Cmd
	switch m.viewMode {
	case MenuView:
		m.menu, cmd = m.menu.Update(msg)
	case WriteView:
		m.write, cmd = m.write.Update(msg)
	case FindNoteView:
		m.findNote, cmd = m.findNote.Update(msg)
	case QuitConfirmView:
		m.quitConfirm, cmd = m.quitConfirm.Update(msg)
	}
	return m, cmd
}

func switchToEnglish() {
	if id, err := sys.InputMethodID(sys.InputMethodEnglish); err == nil {
		_ = sys.SwitchInputMethod(id)
	}
}
