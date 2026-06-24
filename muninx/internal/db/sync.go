package db

import (
	"fmt"
	"strings"

	editstack "github.com/haochend413/muninx/internal/app/editStack"
	"github.com/haochend413/muninx/internal/models"
	"github.com/haochend413/muninx/sys"
	"gorm.io/gorm"
)

// SyncData persists pending updates and deletes from editMap, then reloads the
// full note list from the database.  Creates are no longer tracked here; they
// are written immediately by App.CreateNewNote.
func (d *DB) SyncData(
	notes []*models.Note,
	editMap map[editstack.EditKey]*editstack.Edit) ([]*models.Note, error) {
	if len(editMap) == 0 {
		return d.loadAll()
	}

	notePendingIDs := make([]uint, 0)
	noteDeleteIDs := make([]uint, 0)

	for key, edit := range editMap {
		id := key.ID
		switch edit.EditType {
		case editstack.UpdateNote:
			notePendingIDs = append(notePendingIDs, id)
		case editstack.DeleteNote:
			noteDeleteIDs = append(noteDeleteIDs, id)
		case editstack.None:
			// skip
		}
	}

	notePendingIDs = uniqueIDs(notePendingIDs)
	noteDeleteIDs = uniqueIDs(noteDeleteIDs)

	notesMap := make(map[uint]*models.Note, len(notes))
	for _, n := range notes {
		notesMap[n.ID] = n
	}

	for _, id := range notePendingIDs {
		if n, ok := notesMap[id]; ok {
			if err := d.persistNote(n); err != nil {
				wrappedErr := fmt.Errorf("failed to update note %d: %w", n.ID, err)
				sys.LogError(wrappedErr)
				return nil, wrappedErr
			}
		}
	}

	if err := d.deleteNotes(noteDeleteIDs); err != nil {
		return nil, err
	}

	return d.loadAll()
}

func (d *DB) persistNote(note *models.Note) error {
	if note == nil {
		return nil
	}
	note.Content = strings.TrimSpace(note.Content)
	if err := d.Conn.Save(note).Error; err != nil {
		return err
	}
	if d.EmbedClient != nil && note.Content != "" {
		embedding, err := d.EmbedClient.Embed(note.Content)
		if err != nil {
			sys.LogError(fmt.Errorf("embedding skipped for note %d: %v", note.ID, err))
			return nil
		}
		if err := d.UpsertNoteEmbedding(note.ID, embedding); err != nil {
			sys.LogError(fmt.Errorf("upsert embedding skipped for note %d: %v", note.ID, err))
		}
	}
	return nil
}

func (d *DB) deleteNotes(ids []uint) error {
	for _, id := range ids {
		if err := d.Conn.Delete(&models.Note{}, id).Error; err != nil {
			return err
		}
		// Explicit cleanup: this driver's DSN doesn't enable SQLite's
		// foreign_keys pragma, so the gorm "OnDelete:CASCADE" constraint on
		// Note.Commits is declared in the schema but not enforced by SQLite.
		if err := d.Conn.Where("note_id = ?", id).Delete(&models.NoteCommit{}).Error; err != nil {
			sys.LogError(fmt.Errorf("failed to delete commits for note %d: %v", id, err))
		}
		if err := d.DeleteNoteEmbedding(id); err != nil {
			sys.LogError(fmt.Errorf("failed to delete embedding for note %d: %v", id, err))
		}
	}
	return nil
}

func (d *DB) loadAll() ([]*models.Note, error) {
	var notes []*models.Note
	if err := d.Conn.
		Preload("Commits", func(db *gorm.DB) *gorm.DB {
			return db.Order("note_commits.id ASC")
		}).
		Order("created_at ASC").
		Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func uniqueIDs(ids []uint) []uint {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
