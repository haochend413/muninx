package menu

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/ui/menu/menuinput"
	"github.com/haochend413/muninx/internal/ui/statusbar"
	"github.com/haochend413/muninx/internal/ui/styles"
	bTable "github.com/haochend413/muninx/internal/ui/table"
)

// recentWindow is how far back "recent" notes mode looks.
const recentWindow = 7 * 24 * time.Hour

// Messages sent to the root model.
type SelectNoteMsg struct{ NoteID uint }
type NewNoteRequestMsg struct{}
type DeleteNoteRequestMsg struct{ NoteID uint }
type SyncRequestMsg struct{}
type OpenQuitMsg struct{}

// statusbarMainTag is the single full-width status bar element. It shows
// whatever is currently relevant — for now, just the quit confirmation.
const statusbarMainTag = "main"

// InputMode determines what the input bar does with whatever is typed into
// it. SearchMode is the only mode today; more can be added later (e.g. a
// command mode) without changing how the bar itself is wired up.
type InputMode int

const (
	SearchMode InputMode = iota
)

type Model struct {
	app        *app.App
	table      bTable.Model
	input      menuinput.Model
	statusbar  statusbar.Model
	layout     Layout
	mode       InputMode
	recentOnly bool // when true and search is empty, show only notes edited within recentWindow

	// rowNoteIDs maps each currently visible table row to its note ID, in
	// the same order as the table's rows. The table cursor is a position in
	// this (possibly filtered) list, so any action keyed off the cursor must
	// resolve through rowNoteIDs rather than re-deriving an unfiltered list -
	// otherwise the cursor position silently points at the wrong note
	// whenever a filter (search or recent-only) is active.
	rowNoteIDs []uint
}

// SelectedNoteID returns the note ID for the row currently under the table
// cursor, accounting for whatever filter (search/recent-only) is active. ok
// is false if the table has no rows (e.g. empty filter results).
func (m Model) SelectedNoteID() (id uint, ok bool) {
	c := m.table.Cursor()
	if c < 0 || c >= len(m.rowNoteIDs) {
		return 0, false
	}
	return m.rowNoteIDs[c], true
}

func New(application *app.App) Model {
	menuCols := []bTable.Column{
		{Title: "ID", Width: 6},
		{Title: "Content", Width: 60},
		{Title: "Last Edited", Width: 16},
	}
	t := bTable.New(
		bTable.WithColumns(menuCols),
		bTable.WithFocused(true),
		bTable.WithHeight(20),
	)
	t.SetStyles(styles.FocusedTableStyle)

	mi := menuinput.New()
	mi.Placeholder = "Search or type a command..."
	mi.Focus()

	bar := statusbar.New(0, statusbar.DefaultHeight)
	bar.Register(statusbarMainTag, statusbar.Left, statusbar.ElemConfig{Align: statusbar.AlignLeft})

	return Model{
		app:       application,
		table:     t,
		input:     mi,
		statusbar: bar,
		mode:      SearchMode,
	}
}

func (m Model) Init() tea.Cmd { return nil }

// TickStatusbar lets the status bar process its own expiry/resize messages
// even while this view isn't the active one — a timed Signal (e.g. the
// "synced" flash) must still clear on schedule no matter which view the
// user switches to in the meantime.
func (m *Model) TickStatusbar(msg tea.Msg) {
	m.statusbar, _ = m.statusbar.Update(msg)
}

// ShowQuitPrompt displays the quit confirmation prompt on the status bar.
func (m *Model) ShowQuitPrompt() {
	m.statusbar.Signal(statusbarMainTag, statusbar.QuitPrompt())
}

// ClearQuitPrompt removes the quit confirmation prompt from the status bar.
func (m *Model) ClearQuitPrompt() {
	m.statusbar.Clear(statusbarMainTag)
}

// ShowSyncedMessage flashes a brief confirmation that pending changes were
// pushed to the database.
func (m *Model) ShowSyncedMessage() tea.Cmd {
	return m.statusbar.Signal(statusbarMainTag, statusbar.Synced())
}

// ToggleRecentOnly flips whether the table, when search is empty, shows
// only notes edited within the last 7 days or all notes.
func (m *Model) ToggleRecentOnly() {
	m.recentOnly = !m.recentOnly
}

// UpdateTable refreshes table rows and column widths using the current
// layout. In SearchMode, the input's value filters notes down to those
// containing it, and the Content column shows a snippet around the match
// instead of just the start of the note.
func (m *Model) UpdateTable() {
	l := m.layout
	cols := []bTable.Column{
		{Title: "ID", Width: l.TableIDWidth},
		{Title: "Content", Width: l.TableContentWidth},
		{Title: "", Width: l.TableColumnGap},
		{Title: "Last Edited", Width: l.TableTimeWidth},
	}

	var query string
	if m.mode == SearchMode {
		query = strings.TrimSpace(m.input.Value())
	}

	cutoff := time.Now().Add(-recentWindow)

	notes := m.app.GetDataMgr().GetAllNotesByIDDesc()
	rows := make([]bTable.Row, 0, len(notes))
	rowNoteIDs := make([]uint, 0, len(notes))
	for _, n := range notes {
		if query != "" {
			if !matchesQuery(n, query) {
				continue
			}
		} else if m.recentOnly && n.LastEdit.Before(cutoff) {
			continue
		}

		preview := n.Content
		if query != "" {
			preview = snippetAround(preview, query, l.TableContentWidth/2)
		}

		rows = append(rows, bTable.Row{
			fmt.Sprintf("%d", n.ID),
			preview,
			"",
			formatTimeAgo(n.LastEdit),
		})
		rowNoteIDs = append(rowNoteIDs, n.ID)
	}

	m.rowNoteIDs = rowNoteIDs
	m.table.SetColumns(cols)
	m.table.SetRows(rows)
	m.table.SetWidth(l.TableWidth)
	m.table.SetHeight(l.TableHeight)
}
