package ui

import (
	"fmt"
	"time"

	bTable "github.com/haochend413/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/muninx/config"
	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/models"
	"github.com/haochend413/muninx/internal/ui/menu"
	"github.com/haochend413/muninx/internal/ui/quitconfirm"
	"github.com/haochend413/muninx/internal/ui/write"
	"github.com/haochend413/muninx/state"
)

// ViewMode is which top-level screen is active.
type ViewMode int

const (
	MenuView        ViewMode = iota
	WriteView
	QuitConfirmView
)

// ApplicationView is kept as an alias so any remaining old references compile.
const ApplicationView = MenuView

// tickMsg drives the once-per-second clock tick.
type tickMsg time.Time

type Model struct {
	app    *app.App
	Config *config.Config

	viewMode         ViewMode
	previousViewMode ViewMode

	// Hidden table used only for DistributeState / CollectState.
	notesTable bTable.Model

	// Sub-models, one per view.
	menu        menu.Model
	write       write.Model
	quitConfirm quitconfirm.Model

	width  int
	height int
	ready  bool
}

func NewModel(application *app.App, cfg *config.Config, s *state.State) Model {
	if s == nil {
		s = state.DefaultState()
	}
	if cfg == nil {
		tmp := config.LoadOrCreateConfig()
		cfg = &tmp
	}

	// Hidden state table (not rendered, used for cursor-position persistence).
	noteColumns := []bTable.Column{
		{Title: "ID", Width: 4},
		{Title: "Time", Width: 16},
		{Title: "Content", Width: 40},
		{Title: "Flags", Width: 6},
	}
	noteTable := bTable.New(
		bTable.WithColumns(noteColumns),
		bTable.WithFocused(false),
		bTable.WithHeight(15),
	)

	m := Model{
		app:         application,
		Config:      cfg,
		viewMode:    MenuView,
		notesTable:  noteTable,
		menu:        menu.New(application),
		write:       write.New(application),
		quitConfirm: quitconfirm.New(),
	}

	// Populate hidden state table and restore cursor position.
	m.updateNotesTable()
	m.DistributeState(&s.App)

	return m
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// ---------- Hidden state table (for DistributeState / CollectState) ----------

func (m *Model) updateNotesTable() {
	notes := m.app.GetActiveNoteList()
	rows := make([]bTable.Row, len(notes))
	for i, n := range notes {
		content := n.Content
		if len(content) > 38 {
			content = content[:35] + "..."
		}
		idStr := fmt.Sprintf("%d", n.ID)
		timeStr := n.CreatedAt.Format("06-01-02 15:04")
		if n.ID == 0 {
			idStr = "P"
			timeStr = time.Now().Format("06-01-02 15:04")
		}
		rows[i] = bTable.Row{idStr, timeStr, content, ""}
	}
	m.notesTable.SetRows(rows)
}

// loadNoteIntoEditor updates app state, syncs hidden table cursor, loads the
// note into WriteView, and switches to WriteView.
func (m *Model) loadNoteIntoEditor(note *models.Note) tea.Cmd {
	if note == nil {
		return nil
	}
	m.app.GetDataMgr().SwitchActiveNoteByID(note.ID)

	if ptr := m.app.GetDataMgr().GetActiveNotePtr(); ptr >= 0 {
		m.notesTable.SetCursor(ptr)
	}

	cmd := m.write.LoadNote(note)
	m.viewMode = WriteView
	return cmd
}

// handleNewNote creates a new note and opens it in the editor.
func (m *Model) handleNewNote() (tea.Model, tea.Cmd) {
	m.app.CreateNewNote()
	notes := m.app.GetActiveNoteList()
	if len(notes) == 0 {
		return m, nil
	}
	newNote := notes[len(notes)-1]
	cmd := m.loadNoteIntoEditor(newNote)
	return m, cmd
}
