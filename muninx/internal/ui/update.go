package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/key"

	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/ui/menu"
	"github.com/haochend413/muninx/internal/ui/write"
	statePkg "github.com/haochend413/muninx/state"
	"github.com/haochend413/muninx/sys"
)

// reEmbedDoneMsg is sent when a full re-embed completes.
type reEmbedDoneMsg struct{}

var globalKeys = struct {
	ReEmbed key.Binding
	Undo    key.Binding
}{
	ReEmbed: key.NewBinding(key.WithKeys("ctrl+r")),
	Undo:    key.NewBinding(key.WithKeys("ctrl+z")),
}

var quitConfirmKeys = struct {
	Confirm key.Binding
	Cancel  key.Binding
}{
	Confirm: key.NewBinding(key.WithKeys("ctrl+c")),
	Cancel:  key.NewBinding(key.WithKeys("n", "esc")),
}

func reEmbedCmd(a *app.App) tea.Cmd {
	return func() tea.Msg {
		a.ReEmbedAllNotes()
		return reEmbedDoneMsg{}
	}
}

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
	// Both views' status bars must process every message regardless of
	// which view is active: a timed Signal (e.g. the post-sync "synced"
	// flash) fires on both bars at once, and only the active view's own
	// Update would otherwise see the expiry message — leaving the other
	// bar's signal stuck forever if the user switches views in between.
	m.menu.TickStatusbar(msg)
	m.write.TickStatusbar(msg)

	// While confirming quit, every keypress is gated here: only y/n/esc do
	// anything, and nothing reaches the active sub-model or the global keys
	// below. Non-key messages (resize, ticks) still flow through normally.
	if m.confirmingQuit {
		if kMsg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(kMsg, quitConfirmKeys.Confirm):
				return m.confirmQuit()
			case key.Matches(kMsg, quitConfirmKeys.Cancel):
				m.cancelQuit()
				return m, nil
			default:
				return m, nil
			}
		}
	}

	// Intercept global keys before delegating to sub-models.
	if kMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(kMsg, globalKeys.ReEmbed) {
			return m, reEmbedCmd(m.app)
		}
		if key.Matches(kMsg, globalKeys.Undo) {
			if m.app.UndoDelete() != nil {
				m.menu.UpdateTable()
			}
			return m, nil
		}
	}

	// Global messages handled before view-mode dispatch.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		var wc1, wc2 tea.Cmd
		m.menu, wc1 = m.menu.Update(msg)
		m.write, wc2 = m.write.Update(msg)
		return m, tea.Batch(wc1, wc2)

	case tickMsg:
		return m, tick()

	// Always forward typewriter ticks to the write model so the animation
	// continues regardless of which view is currently active.
	case write.TickMsg:
		var cmd tea.Cmd
		m.write, cmd = m.write.Update(msg)
		return m, cmd

	// --- Async sync completion ---

	case syncDoneMsg:
		m.menu.UpdateTable()
		cmd1 := m.menu.ShowSyncedMessage()
		cmd2 := m.write.ShowSyncedMessage()
		return m, tea.Batch(cmd1, cmd2)

	case reEmbedDoneMsg:
		cmd := m.write.RefreshRelatedNotes()
		return m, cmd

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

	case menu.DeleteNoteRequestMsg:
		notes := m.app.GetDataMgr().GetAllNotesByIDDesc()
		if msg.Index >= 0 && msg.Index < len(notes) {
			m.app.DeleteNoteByID(notes[msg.Index].ID)
			m.menu.UpdateTable()
		}
		return m, nil

	case menu.SyncRequestMsg:
		return m, syncCmd(m.app)

	case menu.OpenQuitMsg:
		m.confirmingQuit = true
		m.menu.ShowQuitPrompt()
		return m, nil

	// --- Messages from write sub-model ---

	case write.BackToMenuMsg:
		m.viewMode = MenuView
		m.menu.UpdateTable()
		switchToEnglish()
		return m, nil

	case write.OpenQuitMsg:
		m.confirmingQuit = true
		m.write.ShowQuitPrompt()
		return m, nil

	case write.SyncRequestMsg:
		return m, syncCmd(m.app)

	case write.OpenNoteMsg:
		cmd := m.loadNoteIntoEditor(msg.Note)
		return m, cmd
	}

	// Delegate unhandled messages to the active sub-model.
	var cmd tea.Cmd
	switch m.viewMode {
	case MenuView:
		m.menu, cmd = m.menu.Update(msg)
	case WriteView:
		m.write, cmd = m.write.Update(msg)
	}
	return m, cmd
}

// confirmQuit saves the current note and state synchronously, then runs the
// database sync in the background before quitting.
func (m Model) confirmQuit() (tea.Model, tea.Cmd) {
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
}

// cancelQuit clears the quit prompt from whichever view's status bar is
// currently showing it.
func (m *Model) cancelQuit() {
	m.confirmingQuit = false
	switch m.viewMode {
	case WriteView:
		m.write.ClearQuitPrompt()
	default:
		m.menu.ClearQuitPrompt()
	}
}

func switchToEnglish() {
	if id, err := sys.InputMethodID(sys.InputMethodEnglish); err == nil {
		_ = sys.SwitchInputMethod(id)
	}
}
