/*
EditMgr is a write-ahead ledger for deferred database operations.

Creates are no longer tracked here — they are written to the database
immediately by App.CreateNewNote, which receives a real autoincrement ID
back from SQLite.  Only updates and deletes are deferred and batched until
SyncWithDatabase is called.

State machine per note (only two states now):
  UpdateNote  →  UpdateNote  : idempotent, keep as update
  UpdateNote  →  DeleteNote  : upgrade to delete
  DeleteNote  →  UpdateNote  : error (cannot update a deleted entity)
  DeleteNote  →  DeleteNote  : error (double delete)
*/
package editstack

import (
	"fmt"

	"github.com/haochend413/muninx/sys"
)

type EditType = int

const (
	None       EditType = -1
	UpdateNote EditType = 1
	DeleteNote EditType = 2
)

// EntityType constants for EditKey
const (
	EntityNote = "note"
)

// EditKey is a composite key for EditMap to avoid ID collisions between entity types.
type EditKey struct {
	EntityType string
	ID         uint
}

type Edit struct {
	ID         uint
	EditType   EditType
	Additional *uint
}

type EditMgr struct {
	EditStack []*Edit
	EditMap   map[EditKey]*Edit
}

func NewEditMgr() *EditMgr {
	return &EditMgr{
		EditStack: make([]*Edit, 0),
		EditMap:   make(map[EditKey]*Edit),
	}
}

func (em *EditMgr) AddEdit(curr *Edit) error {
	em.EditStack = append(em.EditStack, curr)
	id := curr.ID
	tp := curr.EditType
	key := EditKey{EntityType: EntityNote, ID: id}

	if edit, exists := em.EditMap[key]; exists {
		prev := edit.EditType
		switch tp {
		case UpdateNote:
			switch prev {
			case UpdateNote:
				// idempotent
			case DeleteNote:
				err := fmt.Errorf("cannot update note %d: already marked for deletion", id)
				sys.LogError(err)
				return err
			}
		case DeleteNote:
			switch prev {
			case UpdateNote:
				em.EditMap[key].EditType = DeleteNote
			case DeleteNote:
				err := fmt.Errorf("duplicate delete for note %d", id)
				sys.LogError(err)
				return err
			}
		}
	} else {
		em.EditMap[key] = &Edit{ID: id, EditType: tp}
	}
	return nil
}

func (em *EditMgr) Clear() {
	em.EditStack = make([]*Edit, 0)
	em.EditMap = make(map[EditKey]*Edit)
}

func (em *EditMgr) ClearOnSync() {
	em.EditStack = make([]*Edit, 0)
	em.EditMap = make(map[EditKey]*Edit)
}

func (em *EditMgr) RemoveEdit(entityType string, id uint) {
	delete(em.EditMap, EditKey{EntityType: entityType, ID: id})
}

func (em *EditMgr) GetEdit(entityType string, id uint) (*Edit, bool) {
	edit, exists := em.EditMap[EditKey{EntityType: entityType, ID: id}]
	return edit, exists
}
