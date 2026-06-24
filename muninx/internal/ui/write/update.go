package write

import (
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/key"
)

type keyMap struct {
	Save                key.Binding
	Back                key.Binding
	SyncDB              key.Binding
	Quit                key.Binding
	ToggleFocus         key.Binding
	PushCommit          key.Binding
	OmitCommit          key.Binding
	ToggleCommitHistory key.Binding
	CommitHistoryLeft   key.Binding
	CommitHistoryRight  key.Binding
}

var keys = keyMap{
	Save:                key.NewBinding(key.WithKeys("ctrl+s")),
	Back:                key.NewBinding(key.WithKeys("ctrl+x", "esc")),
	SyncDB:              key.NewBinding(key.WithKeys("ctrl+q")),
	Quit:                key.NewBinding(key.WithKeys("ctrl+c")),
	ToggleFocus:         key.NewBinding(key.WithKeys("tab")),
	PushCommit:          key.NewBinding(key.WithKeys("ctrl+p")),
	OmitCommit:          key.NewBinding(key.WithKeys("ctrl+o")),
	ToggleCommitHistory: key.NewBinding(key.WithKeys("h")),
	CommitHistoryLeft:   key.NewBinding(key.WithKeys("left")),
	CommitHistoryRight:  key.NewBinding(key.WithKeys("right")),
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Let the status bar process its own expiry messages regardless of
	// what else this Update call handles below.
	m.statusbar, _ = m.statusbar.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.layout = computeLayout(msg.Width, msg.Height)
		m.applyLayout()
		return m, nil

	case TickMsg:
		if msg.Gen != m.tickGen {
			return m, nil // stale tick from a previous note load
		}
		runes := []rune(m.relatedText)
		if m.revealedChars < len(runes) {
			m.revealedChars++
			// Keep advancing in the background, but don't clobber the
			// commit diff if that's what's currently displayed.
			if m.rightPanelMode == RelatedNotesMode {
				m.updateRelatedViewport()
				// Auto-scroll to bottom while animating, unless user is browsing.
				if m.focus != FocusRelated {
					m.relatedVp.GotoBottom()
				}
			}
			return m, doTick(m.tickGen)
		}
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
			return m, m.ShowSavedMessage()

		case key.Matches(msg, keys.Back):
			m.SaveCurrentNote()
			m.textArea.Blur()
			return m, func() tea.Msg { return BackToMenuMsg{} }

		case key.Matches(msg, keys.PushCommit):
			m.SaveCurrentNote()
			return m, m.CommitNote()

		case key.Matches(msg, keys.OmitCommit):
			return m, m.OmitCommit()

		case key.Matches(msg, keys.ToggleCommitHistory) && m.focus == FocusRelated:
			m.ToggleCommitHistoryView()
			return m, nil

		case key.Matches(msg, keys.CommitHistoryLeft) && m.rightPanelMode == CommitHistoryMode && m.focus == FocusRelated:
			m.CommitHistoryLeft()
			return m, nil

		case key.Matches(msg, keys.CommitHistoryRight) && m.rightPanelMode == CommitHistoryMode && m.focus == FocusRelated:
			m.CommitHistoryRight()
			return m, nil

		case key.Matches(msg, keys.ToggleFocus):
			if m.focus == FocusTextArea {
				m.focus = FocusRelated
				m.textArea.Blur()
			} else {
				m.focus = FocusTextArea
				m.refreshRightPanel() // refresh to dim colors
				return m, m.textArea.Focus()
			}
			m.refreshRightPanel() // refresh to active colors
			return m, nil
		}
	}

	// Delegate unmatched messages to the focused panel.
	var cmd tea.Cmd
	if m.focus == FocusTextArea {
		m.textArea, cmd = m.textArea.Update(msg)
	} else {
		m.relatedVp, cmd = m.relatedVp.Update(msg)
	}
	return m, cmd
}
