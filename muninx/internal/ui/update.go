package ui

import (
	"log"

	bTable "github.com/haochend413/bubbles/v2/table"
	"github.com/haochend413/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/haochend413/muninx/state"
	"github.com/haochend413/muninx/sys"
)

// ---------- Key maps ----------

type menuKeyMap struct {
	NewNote  key.Binding // shift-N — create new note in thread[0]/branch[0]
	Select   key.Binding // Enter  — open selected note
	SyncDB   key.Binding // Ctrl+Q — sync to database
	FindNote key.Binding // Ctrl+F — open find-note overlay
	Quit     key.Binding // Ctrl+C — open quit confirmation
}

var menuKeys = menuKeyMap{
	NewNote:  key.NewBinding(key.WithKeys("N")),
	Select:   key.NewBinding(key.WithKeys("enter")),
	SyncDB:   key.NewBinding(key.WithKeys("ctrl+q")),
	FindNote: key.NewBinding(key.WithKeys("ctrl+f")),
	Quit:     key.NewBinding(key.WithKeys("ctrl+c")),
}

type writeKeyMap struct {
	ToggleFocus  key.Binding // Tab      — switch between editor and list
	Save         key.Binding // Ctrl+S   — save current note
	Back         key.Binding // Ctrl+X   — save and return to menu
	BackFromList key.Binding // Esc      — return to menu when list is focused
	SelectNote   key.Binding // Enter    — switch to selected related note (list only)
	SyncDB       key.Binding // Ctrl+Q   — save then sync
	FindNote     key.Binding // Ctrl+F   — open find-note overlay
	Quit         key.Binding // Ctrl+C   — open quit confirmation
}

var writeKeys = writeKeyMap{
	ToggleFocus:  key.NewBinding(key.WithKeys("tab")),
	Save:         key.NewBinding(key.WithKeys("ctrl+s")),
	Back:         key.NewBinding(key.WithKeys("ctrl+x")),
	BackFromList: key.NewBinding(key.WithKeys("esc")),
	SelectNote:   key.NewBinding(key.WithKeys("enter")),
	SyncDB:       key.NewBinding(key.WithKeys("ctrl+q")),
	FindNote:     key.NewBinding(key.WithKeys("ctrl+f")),
	Quit:         key.NewBinding(key.WithKeys("ctrl+c")),
}

type findKeyMap struct {
	Close      key.Binding // Esc        — close overlay
	FocusLeft  key.Binding // h / left   — move focus left
	FocusRight key.Binding // l / right  — move focus right
	Select     key.Binding // Enter      — open note (notes table only)
}

var findKeys = findKeyMap{
	Close:      key.NewBinding(key.WithKeys("esc")),
	FocusLeft:  key.NewBinding(key.WithKeys("h", "left")),
	FocusRight: key.NewBinding(key.WithKeys("l", "right")),
	Select:     key.NewBinding(key.WithKeys("enter")),
}

type quitKeyMap struct {
	Confirm key.Binding // y
	Reject  key.Binding // n or Esc
}

var quitKeys = quitKeyMap{
	Confirm: key.NewBinding(key.WithKeys("y")),
	Reject:  key.NewBinding(key.WithKeys("n", "esc")),
}

// ---------- Update ----------

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.resizeComponents()
		m.updateMenuTable()
		return m, nil

	case tickMsg:
		return m, tick()

	case bTable.MoveSelectMsg:
		if m.viewMode == FindNoteView {
			switch m.findFocus {
			case FindFocusThreads:
				m.updateFindBranchesForSelectedThread()
				m.updateFindViewport()
			case FindFocusBranches:
				m.updateFindNotesForSelectedBranch()
				m.updateFindViewport()
			case FindFocusNotes:
				m.updateFindViewport()
			}
		}
		return m, nil

	case tea.KeyMsg:
		switch m.viewMode {
		case MenuView:
			switch {
			case key.Matches(msg, menuKeys.Quit):
				m.previousViewMode = MenuView
				m.viewMode = QuitConfirmView
				return m, nil

			case key.Matches(msg, menuKeys.SyncDB):
				m.app.SyncWithDatabase()
				m.updateMenuTable()
				return m, nil

			case key.Matches(msg, menuKeys.FindNote):
				m.openFindOverlay()
				return m, nil

			case key.Matches(msg, menuKeys.NewNote):
				return m.handleNewNote()

			case key.Matches(msg, menuKeys.Select):
				return m.handleMenuSelect()
			}
			// unmatched keys fall through to component update

		case WriteView:
			switch {
			case key.Matches(msg, writeKeys.Quit):
				m.previousViewMode = WriteView
				m.viewMode = QuitConfirmView
				return m, nil

			case key.Matches(msg, writeKeys.SyncDB):
				m.saveCurrentNote()
				m.app.SyncWithDatabase()
				m.updateMenuTable()
				return m, nil

			case key.Matches(msg, writeKeys.Save):
				m.saveCurrentNote()
				return m, nil

			case key.Matches(msg, writeKeys.Back):
				m.saveCurrentNote()
				m.textArea.Blur()
				m.viewMode = MenuView
				m.updateMenuTable()
				switchToEnglish()
				return m, nil

			// ESC only goes back to menu when the list is focused
			// (when textarea is focused, ESC is handled by textarea_vim internally)
			case m.writeFocus == WriteFocusList && key.Matches(msg, writeKeys.BackFromList):
				m.saveCurrentNote()
				m.viewMode = MenuView
				m.updateMenuTable()
				switchToEnglish()
				return m, nil

			case key.Matches(msg, writeKeys.FindNote):
				m.saveCurrentNote()
				m.openFindOverlay()
				return m, nil

			case key.Matches(msg, writeKeys.ToggleFocus):
				cmd = m.toggleWriteFocus()
				return m, cmd

			case m.writeFocus == WriteFocusList && key.Matches(msg, writeKeys.SelectNote):
				cmd = m.handleRelatedNoteSelect()
				return m, cmd
			}
			// unmatched keys fall through to component update

		case FindNoteView:
			switch {
			case key.Matches(msg, findKeys.Close):
				m.viewMode = m.findPreviousView
				return m, nil

			case key.Matches(msg, findKeys.FocusLeft):
				m.moveFindFocusLeft()
				return m, nil

			case key.Matches(msg, findKeys.FocusRight):
				m.moveFindFocusRight()
				return m, nil

			case m.findFocus == FindFocusNotes && key.Matches(msg, findKeys.Select):
				cmd = m.handleFindNoteSelect()
				return m, cmd
			}
			// unmatched keys fall through to component update (j/k navigation)

		case QuitConfirmView:
			switch {
			case key.Matches(msg, quitKeys.Confirm):
				m.saveCurrentNote()
				m.app.SyncWithDatabase()
				if s := m.CollectState(); s != nil {
					if err := state.SaveState(m.Config.StateFilePath, s); err != nil {
						log.Printf("error saving state: %v", err)
					}
				}
				return m, tea.Quit

			case key.Matches(msg, quitKeys.Reject):
				m.viewMode = m.previousViewMode
				return m, nil
			}
		}
	}

	// Forward unhandled messages to the focused component.
	switch m.viewMode {
	case MenuView:
		m.menuTable, cmd = m.menuTable.Update(msg)
		cmds = append(cmds, cmd)

	case WriteView:
		if m.writeFocus == WriteFocusTextArea {
			m.textArea, cmd = m.textArea.Update(msg)
		} else {
			m.relatedList, cmd = m.relatedList.Update(msg)
		}
		cmds = append(cmds, cmd)

	case FindNoteView:
		switch m.findFocus {
		case FindFocusThreads:
			m.findThreadsTable, cmd = m.findThreadsTable.Update(msg)
		case FindFocusBranches:
			m.findBranchesTable, cmd = m.findBranchesTable.Update(msg)
		case FindFocusNotes:
			m.findNotesTable, cmd = m.findNotesTable.Update(msg)
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// ---------- Menu actions ----------

func (m *Model) handleMenuSelect() (tea.Model, tea.Cmd) {
	notes := m.app.GetDataMgr().GetAllNotesByIDDesc()
	cursor := m.menuTable.Cursor()
	if cursor < 0 || cursor >= len(notes) {
		return m, nil
	}
	cmd := m.loadNoteIntoEditor(notes[cursor])
	return m, cmd
}

func (m *Model) handleNewNote() (tea.Model, tea.Cmd) {
	// Ensure first thread exists.
	threads := m.app.GetThreadList()
	if len(threads) == 0 {
		m.app.CreateNewThread(nil)
		threads = m.app.GetThreadList()
	}
	if len(threads) == 0 {
		return m, nil
	}
	m.app.GetDataMgr().SwitchActiveThreadByID(threads[0].ID)

	// Ensure first branch of that thread exists.
	branches := m.app.GetActiveBranchList()
	if len(branches) == 0 {
		m.app.CreateNewBranch(nil)
		branches = m.app.GetActiveBranchList()
	}
	if len(branches) == 0 {
		return m, nil
	}
	m.app.GetDataMgr().SwitchActiveBranchByID(branches[0].ID)

	// Create the note.
	m.app.CreateNewNote(nil)
	notes := m.app.GetActiveNoteList()
	if len(notes) == 0 {
		return m, nil
	}
	newNote := notes[len(notes)-1]
	cmd := m.loadNoteIntoEditor(newNote)
	return m, cmd
}

// ---------- Write actions ----------

func (m *Model) handleRelatedNoteSelect() tea.Cmd {
	item := m.relatedList.SelectedItem()
	if item == nil {
		return nil
	}
	ri, ok := item.(RelatedNoteItem)
	if !ok {
		return nil
	}
	// Try to look up the note in the data manager's global index.
	note := m.app.GetDataMgr().FindNoteByID(ri.NoteID)
	if note == nil {
		// Note not in local cache — just load content without navigation.
		m.textArea.SetValue(ri.Content)
		return m.textArea.Focus()
	}
	m.saveCurrentNote()
	return m.loadNoteIntoEditor(note)
}

// ---------- FindNote actions ----------

func (m *Model) moveFindFocusLeft() {
	switch m.findFocus {
	case FindFocusBranches:
		m.findBranchesTable.Blur()
		m.findThreadsTable.Focus()
		m.findFocus = FindFocusThreads
		m.updateFindViewport()
	case FindFocusNotes:
		m.findNotesTable.Blur()
		m.findBranchesTable.Focus()
		m.findFocus = FindFocusBranches
		m.updateFindViewport()
	}
}

func (m *Model) moveFindFocusRight() {
	switch m.findFocus {
	case FindFocusThreads:
		m.findThreadsTable.Blur()
		m.findBranchesTable.Focus()
		m.findFocus = FindFocusBranches
		m.updateFindViewport()
	case FindFocusBranches:
		m.findBranchesTable.Blur()
		m.findNotesTable.Focus()
		m.findFocus = FindFocusNotes
		m.updateFindViewport()
	}
}

func (m *Model) handleFindNoteSelect() tea.Cmd {
	note := m.selectedFindNote()
	if note == nil {
		return nil
	}
	// Resolve the full note from the global index (has Branches loaded).
	fullNote := m.app.GetDataMgr().FindNoteByID(note.ID)
	if fullNote == nil {
		fullNote = note
	}
	m.saveCurrentNote()
	return m.loadNoteIntoEditor(fullNote)
}

// switchToEnglish switches the macOS input method to English when leaving the editor.
func switchToEnglish() {
	if id, err := sys.InputMethodID(sys.InputMethodEnglish); err == nil {
		_ = sys.SwitchInputMethod(id)
	}
}
