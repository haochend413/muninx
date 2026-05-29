package menu

import (
	"fmt"

	bTable "github.com/haochend413/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/ui/menu/menuinput"
	"github.com/haochend413/muninx/internal/ui/styles"
)

// Messages sent to the root model.
type SelectNoteMsg struct{ Index int }
type NewNoteRequestMsg struct{}
type SyncRequestMsg struct{}
type OpenFindNoteMsg struct{}
type OpenQuitMsg struct{}

type Model struct {
	app    *app.App
	table  bTable.Model
	input  menuinput.Model
	layout Layout
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

	return Model{
		app:   application,
		table: t,
		input: mi,
	}
}

func (m Model) Init() tea.Cmd { return nil }

// UpdateTable refreshes table rows and column widths using the current layout.
func (m *Model) UpdateTable() {
	l := m.layout
	cols := []bTable.Column{
		{Title: "ID", Width: l.TableIDWidth},
		{Title: "Content", Width: l.TableContentWidth},
		{Title: "Last Edited", Width: l.TableTimeWidth},
	}

	notes := m.app.GetDataMgr().GetAllNotesByIDDesc()
	rows := make([]bTable.Row, 0, len(notes))
	for _, n := range notes {
		preview := n.Content
		maxPrev := l.TableContentWidth - 3
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

	m.table.SetColumns(cols)
	m.table.SetRows(rows)
	m.table.SetWidth(l.TableWidth)
	m.table.SetHeight(l.TableHeight)
}
