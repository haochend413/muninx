package openmenu

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/haochend413/muninx/internal/ui/table"
)

// RecentNoteRow holds the display data for one row in the recent-notes table.
type RecentNoteRow struct {
	NoteID     uint
	Preview    string
	LastEdited time.Time
}

// RenderRecentNotesTable renders a full-width table of recently edited notes.
// Columns: NoteID | Content preview | Last updated time.
func RenderRecentNotesTable(rows []RecentNoteRow, width int) string {
	const idWidth = 6
	const timeWidth = 16
	const colPad = 4 // approximate padding/border overhead per column

	contentWidth := width - idWidth - timeWidth - colPad*3
	if contentWidth < 10 {
		contentWidth = 10
	}

	cols := []table.Column{
		{Title: "ID", Width: idWidth},
		{Title: "Content", Width: contentWidth},
		{Title: "Last Updated", Width: timeWidth},
	}

	tableRows := make([]table.Row, len(rows))
	for i, r := range rows {
		tableRows[i] = table.Row{
			fmt.Sprintf("%d", r.NoteID),
			previewText(r.Preview, contentWidth),
			formatTime(r.LastEdited),
		}
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(tableRows),
		table.WithHeight(len(rows)+1),
		table.WithWidth(width),
	)

	return t.View()
}

// previewText returns the first few words of s, truncated to maxLen characters.
func previewText(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if maxLen <= 0 {
		return s
	}
	if len(s) <= maxLen {
		return s
	}
	cut := s[:maxLen-3]
	if i := strings.LastIndexFunc(cut, unicode.IsSpace); i > 0 {
		cut = cut[:i]
	}
	return cut + "..."
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.Format("06-01-02 15:04")
}
