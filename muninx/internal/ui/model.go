package ui

import (
	"fmt"
	"strconv"
	"time"

	bTable "github.com/haochend413/bubbles/v2/table"
	"github.com/haochend413/bubbles/v2/textarea_vim"
	"github.com/haochend413/lipgloss/v2"
	"github.com/haochend413/muninx/config"
	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/models"
	uiList "github.com/haochend413/muninx/internal/ui/list"
	"github.com/haochend413/muninx/internal/ui/openmenu/menuinput"
	"github.com/haochend413/muninx/internal/ui/styles"
	"github.com/haochend413/muninx/internal/ui/viewport"
	"github.com/haochend413/muninx/state"

	tea "charm.land/bubbletea/v2"
)

// ViewMode is which top-level screen is active.
type ViewMode int

const (
	MenuView        ViewMode = iota // opening menu
	WriteView                       // note editor
	QuitConfirmView                 // quit prompt
	FindNoteView                    // find-note overlay
)

// ApplicationView is kept as an alias so any remaining old references compile.
const ApplicationView = MenuView

// WriteFocus tracks which panel in WriteView has keyboard focus.
type WriteFocus int

const (
	WriteFocusTextArea WriteFocus = iota
	WriteFocusList
)

// FindFocus tracks which column in FindNoteView has keyboard focus.
type FindFocus int

const (
	FindFocusThreads FindFocus = iota
	FindFocusBranches
	FindFocusNotes
)

// tickMsg drives the once-per-second clock tick.
type tickMsg time.Time

// ---------- Related-note list item ----------

// RelatedNoteItem implements list.DefaultItem for the WriteView related-notes panel.
type RelatedNoteItem struct {
	NoteID  uint
	Content string
}

func (r RelatedNoteItem) FilterValue() string { return r.Content }
func (r RelatedNoteItem) Title() string       { return fmt.Sprintf("#%d", r.NoteID) }
func (r RelatedNoteItem) Description() string {
	if len(r.Content) > 100 {
		return r.Content[:97] + "..."
	}
	return r.Content
}

// ---------- Model ----------

type Model struct {
	app    *app.App
	Config *config.Config

	// Active view
	viewMode         ViewMode
	previousViewMode ViewMode

	// MenuView
	menuTable bTable.Model
	menuInput menuinput.Model

	// WriteView
	textArea    textarea_vim.Model
	relatedList uiList.Model
	writeFocus  WriteFocus

	// Hidden tables used only for DistributeState / CollectState.
	threadsTable  bTable.Model
	branchesTable bTable.Model
	notesTable    bTable.Model

	// FindNoteView
	findPreviousView  ViewMode
	findFocus         FindFocus
	findThreadsTable  bTable.Model
	findBranchesTable bTable.Model
	findNotesTable    bTable.Model
	findViewport      viewport.Model

	// Terminal dimensions
	width  int
	height int
	ready  bool
}

func colorPtr(c string) *string { return &c }

// NewModel builds the initial UI model.
func NewModel(application *app.App, cfg *config.Config, s *state.State) Model {
	if s == nil {
		s = state.DefaultState()
	}
	if cfg == nil {
		tmp := config.LoadOrCreateConfig()
		cfg = &tmp
	}

	// ---------- hidden state tables (not rendered) ----------
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

	// ---------- MenuView table ----------
	menuCols := []bTable.Column{
		{Title: "ID", Width: 6},
		{Title: "Content", Width: 60},
		{Title: "Last Edited", Width: 16},
	}
	menuTable := bTable.New(
		bTable.WithColumns(menuCols),
		bTable.WithFocused(true),
		bTable.WithHeight(20),
	)
	menuTable.SetStyles(styles.FocusedTableStyle)

	// ---------- MenuView input ----------
	mi := menuinput.New()
	mi.Placeholder = "Search or type a command..."

	// ---------- WriteView textarea ----------
	ta := textarea_vim.New()
	taStyles := ta.Styles()
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	taStyles.Focused.CursorLine = cursorStyle
	taStyles.Blurred.CursorLine = cursorStyle
	ta.SetStyles(taStyles)
	ta.Placeholder = "Start writing..."
	ta.SetWidth(60)
	ta.SetHeight(20)

	// ---------- WriteView related-notes list ----------
	delegate := uiList.NewDefaultDelegate()
	delegate.ShowDescription = true
	relatedList := uiList.New([]uiList.Item{}, delegate, 40, 20)
	relatedList.Title = "Related Notes"
	relatedList.SetShowHelp(false)
	relatedList.SetShowStatusBar(false)
	relatedList.SetShowFilter(false)
	relatedList.DisableQuitKeybindings()

	// ---------- FindNoteView tables ----------
	findTableCols := []bTable.Column{
		{Title: "ID", Width: 3},
		{Title: "Name", Width: 14},
	}
	findThreadsTable := bTable.New(
		bTable.WithColumns(findTableCols),
		bTable.WithFocused(true),
		bTable.WithHeight(20),
	)
	findBranchesTable := bTable.New(
		bTable.WithColumns(findTableCols),
		bTable.WithFocused(false),
		bTable.WithHeight(20),
	)
	findNotesTable := bTable.New(
		bTable.WithColumns(findTableCols),
		bTable.WithFocused(false),
		bTable.WithHeight(20),
	)
	findVP := viewport.New()

	m := Model{
		app:               application,
		Config:            cfg,
		viewMode:          MenuView,
		menuTable:         menuTable,
		menuInput:         mi,
		textArea:          ta,
		relatedList:       relatedList,
		writeFocus:        WriteFocusTextArea,
		threadsTable:      threadTable,
		branchesTable:     branchTable,
		notesTable:        noteTable,
		findThreadsTable:  findThreadsTable,
		findBranchesTable: findBranchesTable,
		findNotesTable:    findNotesTable,
		findViewport:      findVP,
		findFocus:         FindFocusThreads,
	}

	// Populate hidden state tables and restore cursor positions.
	m.updateThreadsTable()
	m.updateBranchesTable()
	m.updateNotesTable()
	m.DistributeState(&s.App)

	return m
}

// Init starts the clock tick.
func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// ---------- Menu table ----------

func (m *Model) updateMenuTable() {
	const idW = 6
	const timeW = 16
	const pad = 8 // border + padding overhead
	contentW := m.width - idW - timeW - pad
	if contentW < 10 {
		contentW = 10
	}

	cols := []bTable.Column{
		{Title: "ID", Width: idW},
		{Title: "Content", Width: contentW},
		{Title: "Last Edited", Width: timeW},
	}

	notes := m.app.GetDataMgr().GetAllNotesByIDDesc()
	rows := make([]bTable.Row, 0, len(notes))
	for _, n := range notes {
		preview := n.Content
		maxPrev := contentW - 3
		if maxPrev < 1 {
			maxPrev = 1
		}
		if len(preview) > maxPrev {
			preview = preview[:maxPrev] + "..."
		}
		timeStr := "—"
		if !n.LastEdit.IsZero() {
			timeStr = n.LastEdit.Format("06-01-02 15:04")
		}
		rows = append(rows, bTable.Row{
			fmt.Sprintf("%d", n.ID),
			preview,
			timeStr,
		})
	}

	m.menuTable.SetColumns(cols)
	m.menuTable.SetRows(rows)
	m.menuTable.SetWidth(m.width)
	tableH := m.height - 18 // header ~12 + input ~3 + help ~1 + margins
	if tableH < 3 {
		tableH = 3
	}
	m.menuTable.SetHeight(tableH)
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

// ---------- WriteView helpers ----------

// loadNoteIntoEditor sets up the WriteView for a given note and switches viewMode.
func (m *Model) loadNoteIntoEditor(note *models.Note) tea.Cmd {
	if note == nil {
		return nil
	}
	m.app.GetDataMgr().SwitchActiveThreadByID(note.ThreadID)
	if len(note.Branches) > 0 {
		m.app.GetDataMgr().SwitchActiveBranchByID(note.Branches[0].ID)
	}
	m.app.GetDataMgr().SwitchActiveNoteByID(note.ID)

	// Sync hidden table cursors for state persistence.
	if ptr := m.app.GetDataMgr().GetActiveThreadPtr(); ptr >= 0 {
		m.threadsTable.SetCursor(ptr)
	}
	if ptr := m.app.GetDataMgr().GetActiveBranchPtr(); ptr >= 0 {
		m.branchesTable.SetCursor(ptr)
	}
	if ptr := m.app.GetDataMgr().GetActiveNotePtr(); ptr >= 0 {
		m.notesTable.SetCursor(ptr)
	}

	m.textArea.SetValue(note.Content)
	cmd := m.textArea.Focus()
	m.writeFocus = WriteFocusTextArea
	m.loadRelatedNotes(note.ID)
	m.viewMode = WriteView
	return cmd
}

// loadRelatedNotes fetches semantic neighbours and populates the list.
func (m *Model) loadRelatedNotes(noteID uint) {
	related := m.app.FetchRelatedNotes(noteID, 10)
	items := make([]uiList.Item, 0, len(related))
	for _, n := range related {
		items = append(items, RelatedNoteItem{
			NoteID:  n.ID,
			Content: n.Content,
		})
	}
	m.relatedList.SetItems(items)
}

// saveCurrentNote persists the current textarea content to the active note.
func (m *Model) saveCurrentNote() {
	if m.app.GetCurrentNoteID() == 0 {
		return
	}
	spl := models.Superlink{
		ThreadID: int(m.app.GetCurrentThreadID()),
		BranchID: int(m.app.GetCurrentBranchID()),
		NoteID:   int(m.app.GetCurrentNoteID()),
	}
	m.app.SetCurrentNoteContent(m.textArea.Value(), &spl)
	m.app.SetCurrentNoteLastEdit()
	m.app.SetCurrentThreadLastEdit()
	m.app.IncrementCurrentThreadFrequency(nil)
	m.app.SetCurrentBranchLastEdit()
	m.app.IncrementCurrentBranchFrequency(nil)
}

// toggleWriteFocus flips focus between the textarea and the related-notes list.
func (m *Model) toggleWriteFocus() tea.Cmd {
	if m.writeFocus == WriteFocusTextArea {
		m.writeFocus = WriteFocusList
		m.textArea.Blur()
		return nil
	}
	m.writeFocus = WriteFocusTextArea
	return m.textArea.Focus()
}

// resizeComponents updates every component's dimensions on WindowSizeMsg.
func (m *Model) resizeComponents() {
	relatedW := (m.width * 40) / 100
	editorW := m.width - relatedW

	listW := relatedW - 4 // account for border (2) + padding (2)
	if listW < 10 {
		listW = 10
	}
	listH := m.height - 6
	if listH < 3 {
		listH = 3
	}
	m.relatedList.SetWidth(listW)
	m.relatedList.SetHeight(listH)

	taW := editorW - 4
	if taW < 10 {
		taW = 10
	}
	taH := m.height - 6
	if taH < 3 {
		taH = 3
	}
	m.textArea.SetWidth(taW)
	m.textArea.SetHeight(taH)

	m.menuInput.SetWidth(m.width - 6)
	m.resizeFindComponents()
}

// ---------- FindNoteView helpers ----------

const findTableInnerW = 17 // inner content width for each find table

// openFindOverlay switches to FindNoteView, saving the current view to return to.
func (m *Model) openFindOverlay() {
	m.findPreviousView = m.viewMode
	m.findFocus = FindFocusThreads
	// Reset table focus: threads gets focus, others lose it.
	m.findThreadsTable.Focus()
	m.findBranchesTable.Blur()
	m.findNotesTable.Blur()
	m.updateFindThreadsTable()
	m.updateFindBranchesForSelectedThread()
	m.updateFindViewport()
	m.resizeFindComponents()
	m.viewMode = FindNoteView
}

// resizeFindComponents updates find table and viewport sizes on window resize.
// Uses overlay dimensions so the panel fits inside the overlay margins.
func (m *Model) resizeFindComponents() {
	overlayW := m.width - 2*findOverlayMarginX
	if overlayW < 60 {
		overlayW = 60
	}
	overlayH := m.height - 2*findOverlayMarginY
	if overlayH < 10 {
		overlayH = 10
	}

	tableBoxW := findTableInnerW + 4 // inner + border(2) + padding(2)
	vpInnerW := overlayW - 3*tableBoxW - 4
	if vpInnerW < 10 {
		vpInnerW = 10
	}
	tableH := overlayH - 4
	if tableH < 3 {
		tableH = 3
	}

	colsNarrow := []bTable.Column{
		{Title: "ID", Width: 3},
		{Title: "Name", Width: findTableInnerW - 4},
	}
	m.findThreadsTable.SetColumns(colsNarrow)
	m.findBranchesTable.SetColumns(colsNarrow)
	m.findNotesTable.SetColumns(colsNarrow)
	m.findThreadsTable.SetWidth(findTableInnerW)
	m.findBranchesTable.SetWidth(findTableInnerW)
	m.findNotesTable.SetWidth(findTableInnerW)
	m.findThreadsTable.SetHeight(tableH)
	m.findBranchesTable.SetHeight(tableH)
	m.findNotesTable.SetHeight(tableH)

	m.findViewport.SetWidth(vpInnerW)
	m.findViewport.SetHeight(tableH)
}

func (m *Model) updateFindThreadsTable() {
	threads := m.app.GetDataMgr().GetThreads()
	rows := make([]bTable.Row, len(threads))
	for i, t := range threads {
		name := t.Name
		if name == "" {
			name = "(unnamed)"
		}
		maxW := findTableInnerW - 4
		if len(name) > maxW {
			name = name[:maxW-3] + "..."
		}
		rows[i] = bTable.Row{fmt.Sprintf("%d", t.ID), name}
	}
	m.findThreadsTable.SetRows(rows)
}

func (m *Model) updateFindBranchesForSelectedThread() {
	thread := m.selectedFindThread()
	if thread == nil {
		m.findBranchesTable.SetRows(nil)
		m.updateFindNotesForSelectedBranch()
		return
	}
	rows := make([]bTable.Row, len(thread.Branches))
	for i, b := range thread.Branches {
		name := b.Name
		if name == "" {
			name = "(unnamed)"
		}
		maxW := findTableInnerW - 4
		if len(name) > maxW {
			name = name[:maxW-3] + "..."
		}
		rows[i] = bTable.Row{fmt.Sprintf("%d", b.ID), name}
	}
	m.findBranchesTable.SetRows(rows)
	m.findBranchesTable.SetCursor(0)
	m.updateFindNotesForSelectedBranch()
}

func (m *Model) updateFindNotesForSelectedBranch() {
	branch := m.selectedFindBranch()
	if branch == nil {
		m.findNotesTable.SetRows(nil)
		return
	}
	rows := make([]bTable.Row, len(branch.Notes))
	for i, n := range branch.Notes {
		preview := n.Content
		maxW := findTableInnerW - 4
		if len(preview) > maxW {
			preview = preview[:maxW-3] + "..."
		}
		rows[i] = bTable.Row{fmt.Sprintf("%d", n.ID), preview}
	}
	m.findNotesTable.SetRows(rows)
	m.findNotesTable.SetCursor(0)
}

func (m *Model) updateFindViewport() {
	var content string
	switch m.findFocus {
	case FindFocusThreads:
		t := m.selectedFindThread()
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
	case FindFocusBranches:
		b := m.selectedFindBranch()
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
	case FindFocusNotes:
		n := m.selectedFindNote()
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
	m.findViewport.SetContent(content)
	m.findViewport.SetYOffset(0)
}

func (m *Model) selectedFindThread() *models.Thread {
	threads := m.app.GetDataMgr().GetThreads()
	cursor := m.findThreadsTable.Cursor()
	if cursor < 0 || cursor >= len(threads) {
		return nil
	}
	return threads[cursor]
}

func (m *Model) selectedFindBranch() *models.Branch {
	thread := m.selectedFindThread()
	if thread == nil {
		return nil
	}
	cursor := m.findBranchesTable.Cursor()
	if cursor < 0 || cursor >= len(thread.Branches) {
		return nil
	}
	return thread.Branches[cursor]
}

func (m *Model) selectedFindNote() *models.Note {
	branch := m.selectedFindBranch()
	if branch == nil {
		return nil
	}
	cursor := m.findNotesTable.Cursor()
	if cursor < 0 || cursor >= len(branch.Notes) {
		return nil
	}
	return branch.Notes[cursor]
}

// formatTimeAgo returns a human-readable relative time string.
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

// suppress unused warning — formatTimeAgo kept for future status bar use.
var _ = formatTimeAgo
var _ = colorPtr
