package app

import (
	"errors"
	"os"
	"time"

	editstack "github.com/haochend413/muninx/internal/app/editStack"
	"github.com/haochend413/muninx/internal/models"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/haochend413/muninx/sys"
)

// current_note.go provides a controlled interface for accessing and modifying
// the currently selected note, with proper edit tracking and synchronization.

// =============================================================================
// Helper: Get current note with safety checks
// =============================================================================

func (a *App) getCurrentNote() *models.Note {
	if a.dataMgr == nil {
		sys.LogError(errors.New("Critical error: dataMgr is nil - app not properly initialized"))
		os.Exit(1)
	}
	return a.dataMgr.GetActiveNote()
}

// =============================================================================
// Getters - Read-only access to current note properties
// =============================================================================

// GetCurrentNoteContent returns the content of the current note, or empty string if none selected
func (a *App) GetCurrentNoteContent() string {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return ""
	}
	return note.Content
}

// GetCurrentNoteID returns the ID of the current note, or 0 if none selected
func (a *App) GetCurrentNoteID() uint {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return 0
	}
	return note.ID
}

// GetCurrentNoteHighlight returns whether the current note is highlighted
func (a *App) GetCurrentNoteHighlight() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return false
	}
	return note.Highlight
}

// GetCurrentNotePrivate returns whether the current note is private
func (a *App) GetCurrentNotePrivate() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return false
	}
	return note.Private
}

// GetCurrentNoteFrequency returns the edit count of the current note
func (a *App) GetCurrentNoteFrequency() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return 0
	}
	return note.Frequency
}

// GetCurrentNoteLastEdit returns the last edit timestamp of the current note.
func (a *App) GetCurrentNoteLastEdit() time.Time {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return time.Time{}
	}
	return note.LastEdit
}

// GetCurrentNoteUpdatedAt returns when the current note was last modified
func (a *App) GetCurrentNoteUpdatedAt() time.Time {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return time.Time{}
	}
	return note.UpdatedAt
}

// GetCurrentNoteCreatedAt returns when the current note was created
func (a *App) GetCurrentNoteCreatedAt() time.Time {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return time.Time{}
	}
	return note.CreatedAt
}

// HasCurrentNote checks if a note is currently selected
func (a *App) HasCurrentNote() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.getCurrentNote() != nil
}

// =============================================================================
// Setters - Controlled modification with edit tracking
// =============================================================================

// SetCurrentNoteContent updates the current note's content with edit tracking.
// This function should automatically update diff. Maybe controlled with a signal.
func (a *App) SetCurrentNoteContent(content string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return
	}

	// No-op if content hasn't changed
	if note.Content == content {
		return
	}

	note.Content = content // update newest content
	note.Frequency++
	note.LastEdit = time.Now()
	a.Synced = false

	edit := &editstack.Edit{ID: note.ID, EditType: editstack.UpdateNote}
	if err := a.editMgr.AddEdit(edit); err != nil {
		sys.LogError(err)
	}
}

func (a *App) CommitCurrentNoteChanges() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return
	}

	// avoid empty commits
	if note.CheckedContent == note.Content {
		return
	}

	dmp := diffmatchpatch.New()
	commit := models.NoteCommit{Patch: dmp.PatchToText(dmp.PatchMake(note.CheckedContent, note.Content)), CommitTime: time.Now()}
	note.Commits = append(note.Commits, &commit)
	note.Diff = dmp.DiffPrettyText((dmp.DiffMain(note.CheckedContent, note.Content, false)))
	note.CheckedContent = note.Content
}

// SetCurrentNoteLastEdit updates the LastEdit timestamp of the current note to the current time.
// Ensures the timestamp is not set to a past time.
func (a *App) SetCurrentNoteLastEdit() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return
	}

	// dont set it backwards
	if time.Now().Before(note.LastEdit) {
		return
	}
	note.LastEdit = time.Now()
}

// ToggleCurrentNoteHighlight toggles the highlight status of the current note
func (a *App) ToggleCurrentNoteHighlight() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return
	}

	note.Highlight = !note.Highlight
	note.UpdatedAt = time.Now()
	a.Synced = false

	edit := &editstack.Edit{ID: note.ID, EditType: editstack.UpdateNote}
	if err := a.editMgr.AddEdit(edit); err != nil {
		sys.LogError(err)
	}
}

// ToggleCurrentNotePrivate toggles the private status of the current note
func (a *App) ToggleCurrentNotePrivate() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return
	}

	note.Private = !note.Private
	note.UpdatedAt = time.Now()
	a.Synced = false

	edit := &editstack.Edit{ID: note.ID, EditType: editstack.UpdateNote}
	if err := a.editMgr.AddEdit(edit); err != nil {
		sys.LogError(err)
	}
}

// DeleteNoteByID marks a note as deleted without removing it from memory:
// it's hidden from rendering immediately, but stays in the dataset (and the
// database) so it can be restored with UndoDelete until the next sync.
func (a *App) DeleteNoteByID(id uint) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.dataMgr.FindNoteByID(id)
	if note == nil || note.Deleted {
		return
	}

	edit := &editstack.Edit{ID: id, EditType: editstack.DeleteNote}
	if err := a.editMgr.AddEdit(edit); err != nil {
		sys.LogError(err)
		return
	}

	note.Deleted = true
	a.deletedStack = append(a.deletedStack, id)
	a.Synced = false
}

// UndoDelete restores the most recently deleted note, if it hasn't been
// synced away yet, and returns it. Returns nil if there's nothing to undo.
func (a *App) UndoDelete() *models.Note {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if len(a.deletedStack) == 0 {
		return nil
	}
	id := a.deletedStack[len(a.deletedStack)-1]
	a.deletedStack = a.deletedStack[:len(a.deletedStack)-1]

	note := a.dataMgr.FindNoteByID(id)
	if note == nil {
		return nil
	}

	note.Deleted = false
	a.editMgr.RemoveEdit(editstack.EntityNote, id)
	return note
}
