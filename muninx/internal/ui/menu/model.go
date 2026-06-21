package menu

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/ui/menu/menuinput"
	"github.com/haochend413/muninx/internal/ui/statusbar"
	"github.com/haochend413/muninx/internal/ui/styles"
	bTable "github.com/haochend413/muninx/internal/ui/table"
)

// Messages sent to the root model.
type SelectNoteMsg struct{ Index int }
type NewNoteRequestMsg struct{}
type DeleteNoteRequestMsg struct{ Index int }
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
	app       *app.App
	table     bTable.Model
	input     menuinput.Model
	statusbar statusbar.Model
	layout    Layout
	mode      InputMode
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

	notes := m.app.GetDataMgr().GetAllNotesByIDDesc()
	rows := make([]bTable.Row, 0, len(notes))
	for _, n := range notes {
		if query != "" && !matchesQuery(n, query) {
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
	}

	m.table.SetColumns(cols)
	m.table.SetRows(rows)
	m.table.SetWidth(l.TableWidth)
	m.table.SetHeight(l.TableHeight)
}
