package db

import (
	"fmt"
	"log"
	"strings"

	editstack "github.com/haochend413/muninx/internal/app/editStack"
	"github.com/haochend413/muninx/internal/models"
)

// SyncData persists pending updates and deletes from editMap, then reloads the
// full thread tree from the database.  Creates are no longer tracked here;
// they are written immediately by CreateThread/CreateBranch/CreateNote.
func (d *DB) SyncData(
	threads []*models.Thread,
	editMap map[editstack.EditKey]*editstack.Edit) ([]*models.Thread, error) {
	if len(editMap) == 0 {
		return d.loadAll()
	}

	notePendingIDs := make([]uint, 0)
	noteDeleteIDs := make([]uint, 0)
	threadPendingIDs := make([]uint, 0)
	threadDeleteIDs := make([]uint, 0)
	branchPendingIDs := make([]uint, 0)
	branchDeleteIDs := make([]uint, 0)

	for key, edit := range editMap {
		id := key.ID
		switch edit.EditType {
		case editstack.UpdateNote:
			notePendingIDs = append(notePendingIDs, id)
		case editstack.DeleteNote:
			noteDeleteIDs = append(noteDeleteIDs, id)
		case editstack.UpdateThread:
			threadPendingIDs = append(threadPendingIDs, id)
		case editstack.DeleteThread:
			threadDeleteIDs = append(threadDeleteIDs, id)
		case editstack.UpdateBranch:
			branchPendingIDs = append(branchPendingIDs, id)
		case editstack.DeleteBranch:
			branchDeleteIDs = append(branchDeleteIDs, id)
		case editstack.None:
			// skip
		}
	}

	notePendingIDs = uniqueIDs(notePendingIDs)
	noteDeleteIDs = uniqueIDs(noteDeleteIDs)
	threadPendingIDs = uniqueIDs(threadPendingIDs)
	threadDeleteIDs = uniqueIDs(threadDeleteIDs)
	branchPendingIDs = uniqueIDs(branchPendingIDs)
	branchDeleteIDs = uniqueIDs(branchDeleteIDs)

	// Build O(1) lookup maps from the in-memory tree.
	threadsMap := make(map[uint]*models.Thread)
	branchesMap := make(map[uint]*models.Branch)
	notesMap := make(map[uint]*models.Note)
	for _, t := range threads {
		threadsMap[t.ID] = t
		for _, b := range t.Branches {
			branchesMap[b.ID] = b
			for _, n := range b.Notes {
				notesMap[n.ID] = n
			}
		}
	}

	for _, id := range threadPendingIDs {
		if t, ok := threadsMap[id]; ok {
			if err := d.persistThread(t); err != nil {
				return nil, fmt.Errorf("failed to update thread %d: %w", t.ID, err)
			}
		}
	}

	for _, id := range notePendingIDs {
		if n, ok := notesMap[id]; ok {
			if err := d.persistNote(n); err != nil {
				return nil, fmt.Errorf("failed to update note %d: %w", n.ID, err)
			}
		}
	}

	for _, id := range branchPendingIDs {
		if b, ok := branchesMap[id]; ok {
			if err := d.persistBranch(b); err != nil {
				return nil, fmt.Errorf("failed to update branch %d: %w", b.ID, err)
			}
		}
	}

	// Delete in reverse dependency order.
	if err := d.deleteNotes(noteDeleteIDs); err != nil {
		return nil, err
	}
	if err := d.deleteBranches(branchDeleteIDs); err != nil {
		return nil, err
	}
	if err := d.deleteThreads(threadDeleteIDs); err != nil {
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
	if err := d.Conn.Model(note).Association("Branches").Replace(note.Branches); err != nil {
		return err
	}
	if d.EmbedClient != nil && note.Content != "" {
		embedding, err := d.EmbedClient.Embed(note.Content)
		if err != nil {
			log.Printf("embedding skipped for note %d: %v", note.ID, err)
			return nil
		}
		if err := d.UpsertNoteEmbedding(note.ID, embedding); err != nil {
			log.Printf("upsert embedding skipped for note %d: %v", note.ID, err)
		}
	}
	return nil
}

func (d *DB) persistThread(thread *models.Thread) error {
	if thread == nil {
		return nil
	}
	if err := d.Conn.Save(thread).Error; err != nil {
		return err
	}
	return d.Conn.Model(thread).Association("Branches").Replace(thread.Branches)
}

func (d *DB) persistBranch(branch *models.Branch) error {
	if branch == nil {
		return nil
	}
	if err := d.Conn.Save(branch).Error; err != nil {
		return err
	}
	return d.Conn.Model(branch).Association("Notes").Replace(branch.Notes)
}

func (d *DB) deleteNotes(ids []uint) error {
	for _, id := range ids {
		if err := d.Conn.Delete(&models.Note{}, id).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) deleteThreads(ids []uint) error {
	for _, id := range ids {
		if err := d.Conn.Delete(&models.Thread{}, id).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) deleteBranches(ids []uint) error {
	for _, id := range ids {
		if err := d.Conn.Delete(&models.Branch{}, id).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) loadAll() ([]*models.Thread, error) {
	var dbThreads []*models.Thread
	if err := d.Conn.
		Preload("Branches.Notes.Branches").
		Order("created_at ASC").
		Find(&dbThreads).Error; err != nil {
		return nil, err
	}
	return dbThreads, nil
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
