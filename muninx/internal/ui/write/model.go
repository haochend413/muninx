package write

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/key"
	"github.com/haochend413/bubbles/v2/textarea_vim"
	"github.com/haochend413/lipgloss/v2"
	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/models"
	"github.com/haochend413/muninx/internal/ui/viewport"
)

// Messages sent to the root model.
type BackToMenuMsg struct{}
type OpenFindNoteMsg struct{}
type OpenQuitMsg struct{}
type SyncRequestMsg struct{}
type OpenNoteMsg struct{ Note *models.Note }

// tickMsg drives the typewriter animation. gen guards against stale ticks.
type tickMsg struct{ gen int }

const tickInterval = 80 * time.Millisecond

func doTick(gen int) tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg{gen: gen}
	})
}

// Focus tracks which panel has keyboard focus.
type Focus int

const (
	FocusTextArea Focus = iota
	FocusRelated
)

// relatedNoteEntry holds the data needed to render one related note.
type relatedNoteEntry struct {
	content   string
	updatedAt time.Time
}

// Per-element styles. Colors are embedded directly in the viewport content so
// they survive independently of any outer lipgloss wrapper.
var (
	noteTextDimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	noteTextLitStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	timestampDimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	timestampLitStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type Model struct {
	app       *app.App
	textArea  textarea_vim.Model
	relatedVp viewport.Model
	focus     Focus
	layout    Layout

	// Typewriter state
	relatedNotes  []relatedNoteEntry // metadata for each note (content + time)
	relatedText   string             // plain concatenated text, used to count revealedChars
	revealedChars int
	tickGen       int
}

func New(application *app.App) Model {
	ta := textarea_vim.New()
	taStyles := ta.Styles()
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	taStyles.Focused.CursorLine = cursorStyle
	taStyles.Blurred.CursorLine = cursorStyle
	ta.SetStyles(taStyles)
	ta.Placeholder = "Start writing..."
	ta.ShowLineNumbers = false
	ta.Statusbar = nil
	// Disable vim mode switching — always stay in insert mode.
	ta.KeyMap.EnterViewMode = key.NewBinding()
	ta.KeyMap.EnterInsertMode = key.NewBinding()
	ta.SetWidth(60)
	ta.SetHeight(20)

	vp := viewport.New()
	vp.SoftWrap = true

	return Model{
		app:       application,
		textArea:  ta,
		relatedVp: vp,
		focus:     FocusTextArea,
	}
}

func (m Model) Init() tea.Cmd { return nil }

// applyLayout pushes layout dimensions into both panels.
func (m *Model) applyLayout() {
	m.textArea.SetWidth(m.layout.TextAreaWidth)
	m.textArea.SetHeight(m.layout.WindowHeight)
	m.relatedVp.SetWidth(m.layout.RelatedWidth)
	m.relatedVp.SetHeight(m.layout.WindowHeight)
}

// renderLines renders each line of s with style independently, preventing
// lipgloss from padding shorter lines to match the longest line (which would
// create blank rows in the viewport when soft-wrap is active).
func renderLines(s string, style lipgloss.Style) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = style.Render(l)
	}
	return strings.Join(lines, "\n")
}

// buildStyledContent constructs the viewport string with per-element ANSI colors.
// Note text and timestamp each get their own color; the focus state selects the
// brightness level.
func (m *Model) buildStyledContent() string {
	focused := m.focus == FocusRelated
	var textStyle, tsStyle lipgloss.Style
	if focused {
		textStyle = noteTextLitStyle
		tsStyle = timestampLitStyle
	} else {
		textStyle = noteTextDimStyle
		tsStyle = timestampDimStyle
	}

	remaining := m.revealedChars
	const sep = "\n\n"
	sepRunes := []rune(sep)
	var sb strings.Builder

	for i, entry := range m.relatedNotes {
		if remaining <= 0 {
			break
		}
		// Insert separator between notes (counts toward revealedChars).
		if i > 0 {
			if remaining >= len(sepRunes) {
				sb.WriteString(sep)
				remaining -= len(sepRunes)
			} else {
				sb.WriteString(string(sepRunes[:remaining]))
				remaining = 0
				break
			}
		}

		noteRunes := []rune(entry.content)
		if remaining >= len(noteRunes) {
			// Whole note revealed — append content then styled timestamp.
			sb.WriteString(renderLines(entry.content, textStyle))
			sb.WriteString("\n")
			ts := "(Last Updated at: unknown)"
			if !entry.updatedAt.IsZero() {
				ts = fmt.Sprintf("(Last Updated at: %s)", entry.updatedAt.Format("2006-01-02 15:04"))
			}
			sb.WriteString(tsStyle.Render(ts))
			remaining -= len(noteRunes)
		} else {
			// Partially revealed.
			sb.WriteString(renderLines(string(noteRunes[:remaining]), textStyle))
			remaining = 0
		}
	}

	return sb.String()
}

// updateRelatedViewport rebuilds the styled viewport content from the current
// reveal position.
func (m *Model) updateRelatedViewport() {
	m.relatedVp.SetContent(m.buildStyledContent())
}

// LoadNote sets up the editor for the given note and starts the typewriter.
func (m *Model) LoadNote(note *models.Note) tea.Cmd {
	if note == nil {
		return nil
	}
	m.textArea.SetValue(note.Content)
	focusCmd := m.textArea.Focus()
	m.focus = FocusTextArea
	m.loadRelatedNotes(note.ID)
	if m.relatedText != "" {
		return tea.Batch(focusCmd, doTick(m.tickGen))
	}
	return focusCmd
}

// SaveCurrentNote persists the textarea content to the active note.
func (m *Model) SaveCurrentNote() {
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

// trimEmptyLines removes lines whose content is entirely whitespace.
func trimEmptyLines(s string) string {
	lines := strings.Split(s, "\n")
	out := lines[:0]
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return strings.Join(out, "\n")
}

// loadRelatedNotes fetches related notes and builds the relatedNotes slice and
// plain relatedText string. Incrementing tickGen invalidates any in-flight tick.
func (m *Model) loadRelatedNotes(noteID uint) {
	related := m.app.FetchRelatedNotes(noteID, 10)
	currentContent := m.textArea.Value()

	m.relatedNotes = m.relatedNotes[:0]
	var textSb strings.Builder
	written := 0

	for _, n := range related {
		if n.ID == noteID || n.Content == currentContent {
			continue
		}
		// Choose the best available timestamp.
		t := n.UpdatedAt
		if !n.LastEdit.IsZero() {
			t = n.LastEdit
		}
		cleaned := trimEmptyLines(n.Content)
		if cleaned == "" {
			continue
		}
		if written > 0 {
			textSb.WriteString("\n\n")
		}
		textSb.WriteString(cleaned)
		m.relatedNotes = append(m.relatedNotes, relatedNoteEntry{
			content:   cleaned,
			updatedAt: t,
		})
		written++
	}

	m.relatedText = textSb.String()
	m.revealedChars = 0
	m.tickGen++
	m.relatedVp.SetContent("")
	m.relatedVp.SetYOffset(0)
}
