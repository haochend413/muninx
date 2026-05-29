package ui

import (
	"fmt"
	"strconv"
	"time"

	bTable "github.com/haochend413/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/muninx/config"
	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/models"
	"github.com/haochend413/muninx/internal/ui/findnote"
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
	FindNoteView
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
	findPreviousView ViewMode

	// Hidden tables used only for DistributeState / CollectState.
	threadsTable  bTable.Model
	branchesTable bTable.Model
	notesTable    bTable.Model

	// Sub-models, one per view.
	menu        menu.Model
	write       write.Model
	quitConfirm quitconfirm.Model
	findNote    findnote.Model

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

	// Hidden state tables (not rendered, used for cursor-position persistence).
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

	branchColumns := []bTable.Column{
		{Title: "ID", Width: 4},
		{Title: "Time", Width: 16},
		{Title: "Name", Width: 40},
		{Title: "#Ns", Width: 2},
		{Title: "Flags", Width: 6},
	}
	branchTable := bTable.New(
		bTable.WithColumns(branchColumns),
		bTable.WithFocused(false),
		bTable.WithHeight(15),
	)

	threadColumns := []bTable.Column{
		{Title: "ID", Width: 4},
		{Title: "Time", Width: 16},
		{Title: "Name", Width: 40},
		{Title: "#Bs", Width: 2},
		{Title: "Flags", Width: 6},
	}
	threadTable := bTable.New(
		bTable.WithColumns(threadColumns),
		bTable.WithFocused(false),
		bTable.WithHeight(15),
	)

	m := Model{
		app:           application,
		Config:        cfg,
		viewMode:      MenuView,
		threadsTable:  threadTable,
		branchesTable: branchTable,
		notesTable:    noteTable,
		menu:          menu.New(application),
		write:         write.New(application),
		quitConfirm:   quitconfirm.New(),
		findNote:      findnote.New(application),
	}

	// Populate hidden state tables and restore cursor positions.
	m.updateThreadsTable()
	m.updateBranchesTable()
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

// ---------- Hidden state tables (for DistributeState / CollectState) ----------

func (m *Model) updateThreadsTable() {
	threads := m.app.GetThreadList()
	rows := make([]bTable.Row, len(threads))
	for i, t := range threads {
		name := t.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}
		idStr := fmt.Sprintf("%d", t.ID)
		timeStr := t.CreatedAt.Format("06-01-02 15:04")
		if t.ID == 0 {
			idStr = "P"
			timeStr = time.Now().Format("06-01-02 15:04")
		}
		rows[i] = bTable.Row{idStr, timeStr, name, strconv.Itoa(len(t.Branches)), ""}
	}
	m.threadsTable.SetRows(rows)
}

func (m *Model) updateBranchesTable() {
	branches := m.app.GetActiveBranchList()
	rows := make([]bTable.Row, len(branches))
	for i, b := range branches {
		name := b.Name
		if len(name) > 38 {
			name = name[:35] + "..."
		}
		idStr := fmt.Sprintf("%d", b.ID)
		timeStr := b.CreatedAt.Format("06-01-02 15:04")
		if b.ID == 0 {
			idStr = "P"
			timeStr = time.Now().Format("06-01-02 15:04")
		}
		rows[i] = bTable.Row{idStr, timeStr, name, strconv.Itoa(len(b.Notes)), ""}
	}
	m.branchesTable.SetRows(rows)
}

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

// loadNoteIntoEditor updates app state, syncs hidden table cursors, loads the
// note into WriteView, and switches to WriteView.
func (m *Model) loadNoteIntoEditor(note *models.Note) tea.Cmd {
	if note == nil {
		return nil
	}
	m.app.GetDataMgr().SwitchActiveThreadByID(note.ThreadID)
	if len(note.Branches) > 0 {
		m.app.GetDataMgr().SwitchActiveBranchByID(note.Branches[0].ID)
	}
	m.app.GetDataMgr().SwitchActiveNoteByID(note.ID)

	if ptr := m.app.GetDataMgr().GetActiveThreadPtr(); ptr >= 0 {
		m.threadsTable.SetCursor(ptr)
	}
	if ptr := m.app.GetDataMgr().GetActiveBranchPtr(); ptr >= 0 {
		m.branchesTable.SetCursor(ptr)
	}
	if ptr := m.app.GetDataMgr().GetActiveNotePtr(); ptr >= 0 {
		m.notesTable.SetCursor(ptr)
	}

	cmd := m.write.LoadNote(note)
	m.viewMode = WriteView
	return cmd
}

// handleNewNote creates a new note and opens it in the editor.
func (m *Model) handleNewNote() (tea.Model, tea.Cmd) {
	threads := m.app.GetThreadList()
	if len(threads) == 0 {
		m.app.CreateNewThread(nil)
		threads = m.app.GetThreadList()
	}
	if len(threads) == 0 {
		return m, nil
	}
	m.app.GetDataMgr().SwitchActiveThreadByID(threads[0].ID)

	branches := m.app.GetActiveBranchList()
	if len(branches) == 0 {
		m.app.CreateNewBranch(nil)
		branches = m.app.GetActiveBranchList()
	}
	if len(branches) == 0 {
		return m, nil
	}
	m.app.GetDataMgr().SwitchActiveBranchByID(branches[0].ID)

	m.app.CreateNewNote(nil)
	notes := m.app.GetActiveNoteList()
	if len(notes) == 0 {
		return m, nil
	}
	newNote := notes[len(notes)-1]
	cmd := m.loadNoteIntoEditor(newNote)
	return m, cmd
}

// openFindOverlay saves the current view and switches to FindNoteView.
func (m *Model) openFindOverlay() {
	m.findPreviousView = m.viewMode
	m.findNote.Open()
	m.viewMode = FindNoteView
}
