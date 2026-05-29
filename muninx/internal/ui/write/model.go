package write

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/textarea_vim"
	"github.com/haochend413/lipgloss/v2"
	"github.com/haochend413/muninx/internal/app"
	"github.com/haochend413/muninx/internal/models"
	uiList "github.com/haochend413/muninx/internal/ui/list"
)

// Messages sent to the root model.
type BackToMenuMsg struct{}
type OpenFindNoteMsg struct{}
type OpenQuitMsg struct{}
type SyncRequestMsg struct{}
type OpenNoteMsg struct{ Note *models.Note }

// Focus tracks which panel has keyboard focus.
type Focus int

const (
	FocusTextArea Focus = iota
	FocusList
)

// RelatedNoteItem implements list.DefaultItem for the related-notes panel.
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

type Model struct {
	app         *app.App
	textArea    textarea_vim.Model
	relatedList uiList.Model
	focus       Focus
	layout      Layout
}

func New(application *app.App) Model {
	ta := textarea_vim.New()
	taStyles := ta.Styles()
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	taStyles.Focused.CursorLine = cursorStyle
	taStyles.Blurred.CursorLine = cursorStyle
	ta.SetStyles(taStyles)
	ta.Placeholder = "Start writing..."
	ta.SetWidth(60)
	ta.SetHeight(20)

	delegate := uiList.NewDefaultDelegate()
	delegate.ShowDescription = true
	relatedList := uiList.New([]uiList.Item{}, delegate, 40, 20)
	relatedList.Title = "Related Notes"
	relatedList.SetShowHelp(false)
	relatedList.SetShowStatusBar(false)
	relatedList.SetShowFilter(false)
	relatedList.DisableQuitKeybindings()

	return Model{
		app:         application,
		textArea:    ta,
		relatedList: relatedList,
		focus:       FocusTextArea,
	}
}

func (m Model) Init() tea.Cmd { return nil }

// applyLayout pushes the current layout dimensions into all components.
func (m *Model) applyLayout() {
	l := m.layout
	m.relatedList.SetWidth(l.ListWidth)
	m.relatedList.SetHeight(l.ListHeight)
	m.textArea.SetWidth(l.TextAreaWidth)
	m.textArea.SetHeight(l.TextAreaHeight)
}

// LoadNote sets up the editor for the given note.
func (m *Model) LoadNote(note *models.Note) tea.Cmd {
	if note == nil {
		return nil
	}
	m.textArea.SetValue(note.Content)
	cmd := m.textArea.Focus()
	m.focus = FocusTextArea
	m.loadRelatedNotes(note.ID)
	return cmd
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

// ToggleFocus flips focus between the textarea and the related-notes list.
func (m *Model) ToggleFocus() tea.Cmd {
	if m.focus == FocusTextArea {
		m.focus = FocusList
		m.textArea.Blur()
		return nil
	}
	m.focus = FocusTextArea
	return m.textArea.Focus()
}
