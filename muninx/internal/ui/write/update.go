package write

import (
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/key"
)

type keyMap struct {
	ToggleFocus  key.Binding
	Save         key.Binding
	Back         key.Binding
	BackFromList key.Binding
	SelectNote   key.Binding
	SyncDB       key.Binding
	FindNote     key.Binding
	Quit         key.Binding
}

var keys = keyMap{
	ToggleFocus:  key.NewBinding(key.WithKeys("tab")),
	Save:         key.NewBinding(key.WithKeys("ctrl+s")),
	Back:         key.NewBinding(key.WithKeys("ctrl+x")),
	BackFromList: key.NewBinding(key.WithKeys("esc")),
	SelectNote:   key.NewBinding(key.WithKeys("enter")),
	SyncDB:       key.NewBinding(key.WithKeys("ctrl+q")),
	FindNote:     key.NewBinding(key.WithKeys("ctrl+f")),
	Quit:         key.NewBinding(key.WithKeys("ctrl+c")),
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.layout = computeLayout(msg.Width, msg.Height)
		m.applyLayout()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, func() tea.Msg { return OpenQuitMsg{} }

		case key.Matches(msg, keys.SyncDB):
			m.SaveCurrentNote()
			return m, func() tea.Msg { return SyncRequestMsg{} }

		case key.Matches(msg, keys.Save):
			m.SaveCurrentNote()
			return m, nil

		case key.Matches(msg, keys.Back):
			m.SaveCurrentNote()
			m.textArea.Blur()
			return m, func() tea.Msg { return BackToMenuMsg{} }

		// ESC only goes back when the list is focused; textarea handles ESC internally.
		case m.focus == FocusList && key.Matches(msg, keys.BackFromList):
			m.SaveCurrentNote()
			return m, func() tea.Msg { return BackToMenuMsg{} }

		case key.Matches(msg, keys.FindNote):
			m.SaveCurrentNote()
			return m, func() tea.Msg { return OpenFindNoteMsg{} }

		case key.Matches(msg, keys.ToggleFocus):
			cmd := m.ToggleFocus()
			return m, cmd

		case m.focus == FocusList && key.Matches(msg, keys.SelectNote):
			cmd := m.handleRelatedNoteSelect()
			return m, cmd
		}
	}

	var cmd tea.Cmd
	if m.focus == FocusTextArea {
		m.textArea, cmd = m.textArea.Update(msg)
	} else {
		m.relatedList, cmd = m.relatedList.Update(msg)
	}
	return m, cmd
}

func (m *Model) handleRelatedNoteSelect() tea.Cmd {
	item := m.relatedList.SelectedItem()
	if item == nil {
		return nil
	}
	ri, ok := item.(RelatedNoteItem)
	if !ok {
		return nil
	}
	note := m.app.GetDataMgr().FindNoteByID(ri.NoteID)
	if note == nil {
		m.textArea.SetValue(ri.Content)
		return m.textArea.Focus()
	}
	m.SaveCurrentNote()
	return func() tea.Msg { return OpenNoteMsg{Note: note} }
}
