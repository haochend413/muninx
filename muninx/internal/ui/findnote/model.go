package findnote

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	bTable "github.com/haochend413/bubbles/v2/table"
	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/models"
	"github.com/haochend413/muninx/internal/ui/viewport"
)

const tableInnerW = 17

// Messages sent to the root model.
type NoteSelectedMsg struct{ Note *models.Note }
type CloseMsg struct{}

// Focus tracks which column has keyboard focus.
type Focus int

const (
	FocusThreads Focus = iota
	FocusBranches
	FocusNotes
)

type Model struct {
	app           *app.App
	focus         Focus
	threadsTable  bTable.Model
	branchesTable bTable.Model
	notesTable    bTable.Model
	vp            viewport.Model
	layout        Layout
}

func New(application *app.App) Model {
	cols := []bTable.Column{
		{Title: "ID", Width: 3},
		{Title: "Name", Width: 14},
	}
	threadsTable := bTable.New(
		bTable.WithColumns(cols),
		bTable.WithFocused(true),
		bTable.WithHeight(20),
	)
	branchesTable := bTable.New(
		bTable.WithColumns(cols),
		bTable.WithFocused(false),
		bTable.WithHeight(20),
	)
	notesTable := bTable.New(
		bTable.WithColumns(cols),
		bTable.WithFocused(false),
		bTable.WithHeight(20),
	)
	vp := viewport.New()

	return Model{
		app:           application,
		focus:         FocusThreads,
		threadsTable:  threadsTable,
		branchesTable: branchesTable,
		notesTable:    notesTable,
		vp:            vp,
	}
}

func (m Model) Init() tea.Cmd { return nil }

// Open resets the overlay state when it is opened.
func (m *Model) Open() {
	m.focus = FocusThreads
	m.threadsTable.Focus()
	m.branchesTable.Blur()
	m.notesTable.Blur()
	m.updateThreadsTable()
	m.updateBranchesForSelectedThread()
	m.updateViewport()
	m.resizeComponents()
}

// resizeComponents applies the current layout dimensions to all table and viewport components.
func (m *Model) resizeComponents() {
	l := m.layout
	m.threadsTable.SetColumns(l.TableColumns)
	m.branchesTable.SetColumns(l.TableColumns)
	m.notesTable.SetColumns(l.TableColumns)
	m.threadsTable.SetWidth(l.TableInnerWidth)
	m.branchesTable.SetWidth(l.TableInnerWidth)
	m.notesTable.SetWidth(l.TableInnerWidth)
	m.threadsTable.SetHeight(l.TableHeight)
	m.branchesTable.SetHeight(l.TableHeight)
	m.notesTable.SetHeight(l.TableHeight)
	m.vp.SetWidth(l.ViewportWidth)
	m.vp.SetHeight(l.TableHeight)
}

func (m *Model) updateThreadsTable() {
	threads := m.app.GetDataMgr().GetThreads()
	rows := make([]bTable.Row, len(threads))
	for i, t := range threads {
		name := t.Name
		if name == "" {
			name = "(unnamed)"
		}
		maxW := tableInnerW - 4
		if len(name) > maxW {
			name = name[:maxW-3] + "..."
		}
		rows[i] = bTable.Row{fmt.Sprintf("%d", t.ID), name}
	}
	m.threadsTable.SetRows(rows)
}

func (m *Model) updateBranchesForSelectedThread() {
	thread := m.selectedThread()
	if thread == nil {
		m.branchesTable.SetRows(nil)
		m.updateNotesForSelectedBranch()
		return
	}
	rows := make([]bTable.Row, len(thread.Branches))
	for i, b := range thread.Branches {
		name := b.Name
		if name == "" {
			name = "(unnamed)"
		}
		maxW := tableInnerW - 4
		if len(name) > maxW {
			name = name[:maxW-3] + "..."
		}
		rows[i] = bTable.Row{fmt.Sprintf("%d", b.ID), name}
	}
	m.branchesTable.SetRows(rows)
	m.branchesTable.SetCursor(0)
	m.updateNotesForSelectedBranch()
}

func (m *Model) updateNotesForSelectedBranch() {
	branch := m.selectedBranch()
	if branch == nil {
		m.notesTable.SetRows(nil)
		return
	}
	rows := make([]bTable.Row, len(branch.Notes))
	for i, n := range branch.Notes {
		preview := n.Content
		maxW := tableInnerW - 4
		if len(preview) > maxW {
			preview = preview[:maxW-3] + "..."
		}
		rows[i] = bTable.Row{fmt.Sprintf("%d", n.ID), preview}
	}
	m.notesTable.SetRows(rows)
	m.notesTable.SetCursor(0)
}

func (m *Model) updateViewport() {
	var content string
	switch m.focus {
	case FocusThreads:
		t := m.selectedThread()
		if t == nil {
			content = "(no thread selected)"
		} else {
			name := t.Name
			if name == "" {
				name = "(unnamed)"
			}
			summary := t.Summary
			if summary == "" {
				summary = "(no summary)"
			}
			content = fmt.Sprintf("Thread #%d\n\nName: %s\n\nSummary:\n%s\n\nBranches: %d\nFrequency: %d",
				t.ID, name, summary, len(t.Branches), t.Frequency)
		}
	case FocusBranches:
		b := m.selectedBranch()
		if b == nil {
			content = "(no branch selected)"
		} else {
			name := b.Name
			if name == "" {
				name = "(unnamed)"
			}
			summary := b.Summary
			if summary == "" {
				summary = "(no summary)"
			}
			content = fmt.Sprintf("Branch #%d\n\nName: %s\n\nSummary:\n%s\n\nNotes: %d\nFrequency: %d",
				b.ID, name, summary, len(b.Notes), b.Frequency)
		}
	case FocusNotes:
		n := m.selectedNote()
		if n == nil {
			content = "(no note selected)"
		} else {
			editTime := "—"
			if !n.LastEdit.IsZero() {
				editTime = n.LastEdit.Format("2006-01-02 15:04")
			}
			content = fmt.Sprintf("Note #%d  (last edit: %s)\n\n%s", n.ID, editTime, n.Content)
		}
	}
	m.vp.SetContent(content)
	m.vp.SetYOffset(0)
}

func (m *Model) selectedThread() *models.Thread {
	threads := m.app.GetDataMgr().GetThreads()
	cursor := m.threadsTable.Cursor()
	if cursor < 0 || cursor >= len(threads) {
		return nil
	}
	return threads[cursor]
}

func (m *Model) selectedBranch() *models.Branch {
	thread := m.selectedThread()
	if thread == nil {
		return nil
	}
	cursor := m.branchesTable.Cursor()
	if cursor < 0 || cursor >= len(thread.Branches) {
		return nil
	}
	return thread.Branches[cursor]
}

func (m *Model) selectedNote() *models.Note {
	branch := m.selectedBranch()
	if branch == nil {
		return nil
	}
	cursor := m.notesTable.Cursor()
	if cursor < 0 || cursor >= len(branch.Notes) {
		return nil
	}
	return branch.Notes[cursor]
}

// MoveFocusLeft shifts keyboard focus one column to the left.
func (m *Model) MoveFocusLeft() {
	switch m.focus {
	case FocusBranches:
		m.branchesTable.Blur()
		m.threadsTable.Focus()
		m.focus = FocusThreads
		m.updateViewport()
	case FocusNotes:
		m.notesTable.Blur()
		m.branchesTable.Focus()
		m.focus = FocusBranches
		m.updateViewport()
	}
}

// MoveFocusRight shifts keyboard focus one column to the right.
func (m *Model) MoveFocusRight() {
	switch m.focus {
	case FocusThreads:
		m.threadsTable.Blur()
		m.branchesTable.Focus()
		m.focus = FocusBranches
		m.updateViewport()
	case FocusBranches:
		m.branchesTable.Blur()
		m.notesTable.Focus()
		m.focus = FocusNotes
		m.updateViewport()
	}
}

// MarginX returns the horizontal overlay margin used by the root for compositing.
func (m Model) MarginX() int { return overlayMarginX }

// MarginY returns the vertical overlay margin used by the root for compositing.
func (m Model) MarginY() int { return overlayMarginY }
