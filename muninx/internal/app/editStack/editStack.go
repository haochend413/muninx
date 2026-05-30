/*
EditMgr is a write-ahead ledger for deferred database operations.

Creates are no longer tracked here — they are written to the database
immediately by App.CreateNewThread/Branch/Note, which receive a real
autoincrement ID back from SQLite.  Only updates and deletes are deferred
and batched until SyncWithDatabase is called.

State machine per entity (only two states now):
  UpdateX  →  UpdateX  : idempotent, keep as update
  UpdateX  →  DeleteX  : upgrade to delete
  DeleteX  →  UpdateX  : error (cannot update a deleted entity)
  DeleteX  →  DeleteX  : error (double delete)
*/
package editstack

import (
	"fmt"

	"github.com/haochend413/muninx/internal/models"
)

type EditType = int

const (
	None         EditType = -1
	UpdateNote   EditType = 1
	DeleteNote   EditType = 2
	UpdateThread EditType = 4
	DeleteThread EditType = 6
	UpdateBranch EditType = 8
	DeleteBranch EditType = 10
)

// EntityType constants for EditKey
const (
	EntityNote   = "note"
	EntityBranch = "branch"
	EntityThread = "thread"
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

type NoteEdit struct {
	Link models.Superlink
}

type EditMgr struct {
	NoteEditStack []*NoteEdit
	EditStack     []*Edit
	EditMap       map[EditKey]*Edit
}

func getEntityType(tp EditType) string {
	switch tp {
	case UpdateNote, DeleteNote:
		return EntityNote
	case UpdateThread, DeleteThread:
		return EntityThread
	case UpdateBranch, DeleteBranch:
		return EntityBranch
	default:
		return ""
	}
}

func NewEditMgr() *EditMgr {
	return &EditMgr{
		NoteEditStack: make([]*NoteEdit, 0),
		EditStack:     make([]*Edit, 0),
		EditMap:       make(map[EditKey]*Edit),
	}
}

func AppendNoteEdit(stack []*NoteEdit, ne *NoteEdit) []*NoteEdit {
	for i, n := range stack {
		if ne.Link == n.Link {
			stack = append(stack[:i], stack[i+1:]...)
			break
		}
	}
	return append(stack, ne)
}

func (em *EditMgr) AddEdit(curr *Edit, spl *models.Superlink) error {
	em.EditStack = append(em.EditStack, curr)
	id := curr.ID
	tp := curr.EditType
	key := EditKey{EntityType: getEntityType(tp), ID: id}

	s := models.Superlink{ThreadID: -1, BranchID: -1, NoteID: -1}
	if spl != nil {
		s = *spl
	}
	ne := NoteEdit{Link: s}

	if edit, exists := em.EditMap[key]; exists {
		prev := edit.EditType
		switch tp {
		case UpdateNote:
			switch prev {
			case UpdateNote:
				em.NoteEditStack = AppendNoteEdit(em.NoteEditStack, &ne)
			case DeleteNote:
				return fmt.Errorf("cannot update note %d: already marked for deletion", id)
			}
		case DeleteNote:
			switch prev {
			case UpdateNote:
				em.EditMap[key].EditType = DeleteNote
			case DeleteNote:
				return fmt.Errorf("duplicate delete for note %d", id)
			}
		case UpdateThread:
			switch prev {
			case UpdateThread:
				// idempotent
			case DeleteThread:
				return fmt.Errorf("cannot update thread %d: already marked for deletion", id)
			}
		case DeleteThread:
			switch prev {
			case UpdateThread:
				em.EditMap[key].EditType = DeleteThread
			case DeleteThread:
				return fmt.Errorf("duplicate delete for thread %d", id)
			}
		case UpdateBranch:
			switch prev {
			case UpdateBranch:
				// idempotent
			case DeleteBranch:
				return fmt.Errorf("cannot update branch %d: already marked for deletion", id)
			}
		case DeleteBranch:
			switch prev {
			case UpdateBranch:
				em.EditMap[key].EditType = DeleteBranch
			case DeleteBranch:
				return fmt.Errorf("duplicate delete for branch %d", id)
			}
		}
	} else {
		em.EditMap[key] = &Edit{ID: id, EditType: tp}
		if spl != nil {
			em.NoteEditStack = AppendNoteEdit(em.NoteEditStack, &ne)
		}
	}
	return nil
}

func (em *EditMgr) Clear() {
	em.NoteEditStack = make([]*NoteEdit, 0)
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
