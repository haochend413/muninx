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
	"github.com/haochend413/muninx/internal/ui/statusbar"
	"github.com/haochend413/muninx/internal/ui/viewport"
)

// Messages sent to the root model.
type BackToMenuMsg struct{}
type OpenQuitMsg struct{}
type SyncRequestMsg struct{}
type OpenNoteMsg struct{ Note *models.Note }

// TickMsg drives the typewriter animation. Gen guards against stale ticks.
// Exported so the root model can intercept it and always route it here,
// keeping the animation alive regardless of which view is currently active.
type TickMsg struct{ Gen int }

const tickInterval = 80 * time.Millisecond

// statusbarMainTag is the single full-width status bar element. It shows
// whatever is currently relevant — for now, just the quit confirmation.
const statusbarMainTag = "main"

func doTick(gen int) tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return TickMsg{Gen: gen}
	})
}

// Focus tracks which panel has keyboard focus.
type Focus int

const (
	FocusTextArea Focus = iota
	FocusRelated
)

// RightPanelMode controls what the right-hand panel shows while it has
// focus: related notes (the default), or commit history — a diff for some
// position in the note's commit stack, navigable with left/right arrows.
type RightPanelMode int

const (
	RelatedNotesMode RightPanelMode = iota
	CommitHistoryMode
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
	app            *app.App
	textArea       textarea_vim.Model
	relatedVp      viewport.Model
	statusbar      statusbar.Model
	focus          Focus
	rightPanelMode RightPanelMode
	commitIndex    int // position in commit history mode; len(commits) means "next commit preview"
	layout         Layout

	// Typewriter state
	relatedNotes  []relatedNoteEntry // metadata for each note (content + time)
	relatedText   string             // plain concatenated text, used to count revealedChars
	revealedChars int
	tickGen       int

	// loadedNoteID is the note actually opened into the textarea via
	// LoadNote, or 0 if none has been opened yet this session. The app
	// always has *some* active note by default (the first one loaded from
	// the DB), even before the user opens anything in the write view - so
	// SaveCurrentNote must gate on this instead of on app.GetCurrentNoteID,
	// or it would overwrite that default-active note's content with the
	// textarea's empty starting value on a quit-without-opening-anything.
	loadedNoteID uint
}

func New(application *app.App) Model {
	ta := textarea_vim.New()
	taStyles := ta.Styles()
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	taStyles.Focused.CursorLine = cursorStyle
	// Match the related panel's dim/lit colors, so whichever side has focus
	// is obviously the brighter one. The cursor's line is styled separately
	// from the rest of the text, so it needs the same dim treatment while
	// blurred — otherwise it stays bright regardless of focus.
	taStyles.Blurred.CursorLine = noteTextDimStyle
	taStyles.Focused.Text = noteTextLitStyle
	taStyles.Blurred.Text = noteTextDimStyle
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

	bar := statusbar.New(0, statusbar.DefaultHeight)
	bar.Register(statusbarMainTag, statusbar.Left, statusbar.ElemConfig{Align: statusbar.AlignLeft})

	return Model{
		app:       application,
		textArea:  ta,
		relatedVp: vp,
		statusbar: bar,
		focus:     FocusTextArea,
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

// applyLayout pushes layout dimensions into both panels and the status bar.
func (m *Model) applyLayout() {
	m.textArea.SetWidth(m.layout.TextAreaWidth)
	m.textArea.SetHeight(m.layout.ContentHeight)
	m.relatedVp.SetWidth(m.layout.RelatedWidth)
	m.relatedVp.SetHeight(m.layout.ContentHeight)
	m.statusbar.SetWidth(m.layout.WindowWidth)
	m.statusbar.SetElemWidth(statusbarMainTag, m.layout.WindowWidth)
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

// ShowSavedMessage flashes a brief confirmation that the note's modified
// content was saved.
func (m *Model) ShowSavedMessage() tea.Cmd {
	return m.statusbar.Signal(statusbarMainTag, statusbar.Saved())
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

// refreshRightPanel re-renders the right panel for whichever mode is
// currently active.
func (m *Model) refreshRightPanel() {
	switch m.rightPanelMode {
	case CommitHistoryMode:
		if m.commitIndex >= m.app.GetCommitCount() {
			m.relatedVp.SetContent(m.app.GetNextCommitPreview())
		} else {
			m.relatedVp.SetContent(m.app.GetCommitDiffAt(m.commitIndex))
		}
	default:
		m.updateRelatedViewport()
	}
}

// ToggleCommitHistoryView switches the right panel between related notes
// and commit history. Entering commit history starts at the preview of
// the next commit (saving the textarea's current value first, so the
// preview reflects what's actually been typed rather than a stale Content
// field); left/right arrows then step backward/forward through the note's
// commit stack.
func (m *Model) ToggleCommitHistoryView() {
	if m.rightPanelMode == CommitHistoryMode {
		m.rightPanelMode = RelatedNotesMode
	} else {
		m.SaveCurrentNote()
		m.commitIndex = m.app.GetCommitCount()
		m.rightPanelMode = CommitHistoryMode
	}
	m.refreshRightPanel()
}

// CommitHistoryLeft steps to the previous (older) position in the commit
// stack, if any.
func (m *Model) CommitHistoryLeft() {
	if m.commitIndex > 0 {
		m.commitIndex--
		m.refreshRightPanel()
	}
}

// CommitHistoryRight steps to the next (newer) position in the commit
// stack, up to the next-commit preview, if any.
func (m *Model) CommitHistoryRight() {
	if m.commitIndex < m.app.GetCommitCount() {
		m.commitIndex++
		m.refreshRightPanel()
	}
}

// CommitNote pushes a commit for the active note, flashes a confirmation,
// and keeps the right panel in sync if it's currently showing a commit
// diff or preview. No flash if there was nothing to commit.
func (m *Model) CommitNote() tea.Cmd {
	if !m.app.CommitNote() {
		return nil
	}
	if m.rightPanelMode == CommitHistoryMode {
		m.commitIndex = m.app.GetCommitCount()
		m.refreshRightPanel()
	}
	return m.statusbar.Signal(statusbarMainTag, statusbar.Committed())
}

// OmitCommit undoes the active note's most recent commit, flashes a
// confirmation, and keeps the right panel in sync if it's currently
// showing a commit diff or preview. No flash if there was nothing to omit.
func (m *Model) OmitCommit() tea.Cmd {
	if !m.app.OmitCommit() {
		return nil
	}
	if m.rightPanelMode == CommitHistoryMode {
		if count := m.app.GetCommitCount(); m.commitIndex > count {
			m.commitIndex = count
		}
		m.refreshRightPanel()
	}
	return m.statusbar.Signal(statusbarMainTag, statusbar.CommitOmitted())
}

// RefreshRelatedNotes re-fetches related notes for the currently open note and
// restarts the typewriter animation. Call this after a re-embed completes so
// the related panel reflects the updated embeddings.
func (m *Model) RefreshRelatedNotes() tea.Cmd {
	noteID := m.app.GetCurrentNoteID()
	if noteID == 0 {
		return nil
	}
	m.loadRelatedNotes(noteID)
	if m.relatedText != "" {
		return doTick(m.tickGen)
	}
	return nil
}

// LoadNote sets up the editor for the given note and starts the typewriter.
func (m *Model) LoadNote(note *models.Note) tea.Cmd {
	if note == nil {
		return nil
	}
	m.loadedNoteID = note.ID
	m.textArea.SetValue(note.Content)
	focusCmd := m.textArea.Focus()
	m.focus = FocusTextArea
	m.rightPanelMode = RelatedNotesMode
	m.commitIndex = 0
	m.loadRelatedNotes(note.ID)
	if m.relatedText != "" {
		return tea.Batch(focusCmd, doTick(m.tickGen))
	}
	return focusCmd
}

// SaveCurrentNote persists the textarea content to the active note.
func (m *Model) SaveCurrentNote() {
	if m.loadedNoteID == 0 || m.app.GetCurrentNoteID() != m.loadedNoteID {
		return
	}
	m.app.SetCurrentNoteContent(m.textArea.Value())
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
