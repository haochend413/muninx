package iobuf

import (
	"slices"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/Munina/internal/app/iobuf/memoization"
	"github.com/haochend413/bubbles/v2/cursor"
	"github.com/haochend413/bubbles/v2/key"
)

// Update is the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focus {
		m.virtualCursor.Blur()
		return m, nil
	}

	// Used to determine if the cursor should blink.
	oldRow, oldCol := m.cursorLineNumber(), m.col

	var cmds []tea.Cmd

	if m.value[m.row] == nil {
		m.value[m.row] = make([]rune, 0)
	}

	if m.MaxHeight > 0 && m.MaxHeight != m.cache.Capacity() {
		m.cache = memoization.NewMemoCache[line, [][]rune](m.MaxHeight)
	}

	switch msg := msg.(type) {
	case tea.PasteMsg:
		m.insertRunesFromUserInput([]rune(msg.Content))
	case tea.KeyPressMsg:
		// Handle mode switching first

		// Only allow editing in INSERT mode
		if !m.InsertMode {
			// VIEW mode - navigation only
			switch {
			case key.Matches(msg, m.KeyMap.EnterInsertMode):
				m.SetInsertMode()
				m.updateWordCount()
				return m, nil

			case key.Matches(msg, m.KeyMap.LineEnd):
				m.CursorEnd()
			case key.Matches(msg, m.KeyMap.LineStart):
				m.CursorStart()
			case key.Matches(msg, m.KeyMap.CharacterForward):
				m.characterRight()
			case key.Matches(msg, m.KeyMap.LineNext):
				m.setCursorLineRelative(+1)
			case key.Matches(msg, m.KeyMap.WordForward):
				m.wordRight()
			case key.Matches(msg, m.KeyMap.CharacterBackward):
				m.characterLeft(false)
			case key.Matches(msg, m.KeyMap.LinePrevious):
				m.setCursorLineRelative(-1)
			case key.Matches(msg, m.KeyMap.WordBackward):
				m.wordLeft()
			case key.Matches(msg, m.KeyMap.InputBegin):
				m.MoveToBegin()
			case key.Matches(msg, m.KeyMap.InputEnd):
				m.MoveToEnd()
			case key.Matches(msg, m.KeyMap.PageUp):
				m.PageUp()
			case key.Matches(msg, m.KeyMap.PageDown):
				m.PageDown()
			}
			// Skip all editing keys in VIEW mode
			break
		}

		// INSERT mode - full editing capabilities
		switch {
		case key.Matches(msg, m.KeyMap.EnterViewMode):
			m.SetViewMode()
			m.updateWordCount()
			return m, nil
		case key.Matches(msg, m.KeyMap.DeleteAfterCursor):
			m.col = clamp(m.col, 0, len(m.value[m.row]))
			if m.col >= len(m.value[m.row]) {
				m.mergeLineBelow(m.row)
				m.updateWordCount()
				break
			}
			m.deleteAfterCursor()
			m.updateWordCount()
		case key.Matches(msg, m.KeyMap.DeleteBeforeCursor):
			m.col = clamp(m.col, 0, len(m.value[m.row]))
			if m.col <= 0 {
				m.mergeLineAbove(m.row)
				m.updateWordCount()
				break
			}
			m.deleteBeforeCursor()
			m.updateWordCount()
		case key.Matches(msg, m.KeyMap.DeleteCharacterBackward):
			m.col = clamp(m.col, 0, len(m.value[m.row]))
			if m.col <= 0 {
				m.mergeLineAbove(m.row)
				m.updateWordCount()
				break
			}
			if len(m.value[m.row]) > 0 {
				m.value[m.row] = append(m.value[m.row][:max(0, m.col-1)], m.value[m.row][m.col:]...)
				if m.col > 0 {
					m.SetCursorColumn(m.col - 1)
				}
			}
			m.updateWordCount()
		case key.Matches(msg, m.KeyMap.DeleteCharacterForward):
			if len(m.value[m.row]) > 0 && m.col < len(m.value[m.row]) {
				m.value[m.row] = slices.Delete(m.value[m.row], m.col, m.col+1)
			}
			if m.col >= len(m.value[m.row]) {
				m.mergeLineBelow(m.row)
			}
			m.updateWordCount()
		case key.Matches(msg, m.KeyMap.DeleteWordBackward):
			if m.col <= 0 {
				m.mergeLineAbove(m.row)
				m.updateWordCount()
				break
			}
			m.deleteWordLeft()
			m.updateWordCount()
		case key.Matches(msg, m.KeyMap.DeleteWordForward):
			m.col = clamp(m.col, 0, len(m.value[m.row]))
			if m.col >= len(m.value[m.row]) {
				m.mergeLineBelow(m.row)
				m.updateWordCount()
				break
			}
			m.deleteWordRight()
			m.updateWordCount()
		case key.Matches(msg, m.KeyMap.InsertNewline):
			if m.MaxHeight > 0 && len(m.value) >= m.MaxHeight {
				return m, nil
			}
			m.col = clamp(m.col, 0, len(m.value[m.row]))
			m.splitLine(m.row, m.col)
			m.updateWordCount()
		case key.Matches(msg, m.KeyMap.LineEnd):
			m.CursorEnd()
		case key.Matches(msg, m.KeyMap.LineStart):
			m.CursorStart()
		case key.Matches(msg, m.KeyMap.CharacterForward):
			m.characterRight()
		case key.Matches(msg, m.KeyMap.LineNext):
			m.CursorDown()
		case key.Matches(msg, m.KeyMap.WordForward):
			m.wordRight()
		case key.Matches(msg, m.KeyMap.Paste):
			return m, Paste
		case key.Matches(msg, m.KeyMap.CharacterBackward):
			m.characterLeft(false /* insideLine */)
		case key.Matches(msg, m.KeyMap.LinePrevious):
			m.CursorUp()
		case key.Matches(msg, m.KeyMap.WordBackward):
			m.wordLeft()
		case key.Matches(msg, m.KeyMap.InputBegin):
			m.MoveToBegin()
		case key.Matches(msg, m.KeyMap.InputEnd):
			m.MoveToEnd()
		case key.Matches(msg, m.KeyMap.PageUp):
			m.PageUp()
		case key.Matches(msg, m.KeyMap.PageDown):
			m.PageDown()
		case key.Matches(msg, m.KeyMap.LowercaseWordForward):
			m.lowercaseRight()
		case key.Matches(msg, m.KeyMap.UppercaseWordForward):
			m.uppercaseRight()
		case key.Matches(msg, m.KeyMap.CapitalizeWordForward):
			m.capitalizeRight()
		case key.Matches(msg, m.KeyMap.TransposeCharacterBackward):
			m.transposeLeft()

		default:
			m.insertRunesFromUserInput([]rune(msg.Text))
			m.updateWordCount()
		}

	case pasteMsg:
		m.insertRunesFromUserInput([]rune(msg))
		m.updateWordCount()

	case pasteErrMsg:
		m.Err = msg
	}

	// Make sure we set the content of the viewport before updating it.
	view := m.view()
	m.viewport.SetContent(view)
	vp, cmd := m.viewport.Update(msg)
	m.viewport = &vp
	cmds = append(cmds, cmd)

	if m.useVirtualCursor {
		m.virtualCursor, cmd = m.virtualCursor.Update(msg)

		// If the cursor has moved, reset the blink state. This is a small UX
		// nuance that makes cursor movement obvious and feel snappy.
		newRow, newCol := m.cursorLineNumber(), m.col
		if (newRow != oldRow || newCol != oldCol) && m.virtualCursor.Mode() == cursor.CursorBlink {
			m.virtualCursor.IsBlinked = false
			cmd = m.virtualCursor.Blink()
		}
		cmds = append(cmds, cmd)
	}

	m.repositionView()

	return m, tea.Batch(cmds...)
}
