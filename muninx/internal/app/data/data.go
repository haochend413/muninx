package data

import (
	"sort"

	"github.com/haochend413/muninx/internal/models"
)

// DataMgr keeps the in-memory list of all notes, exposing the active note and
// handling switching/CRUD bookkeeping. It is the layer above db: a single
// source of truth for what notes the app currently knows about.
type DataMgr struct {
	notes         []*models.Note
	activeNotePtr int
	activeNoteID  uint
	noteIndexByID map[uint]int
}

func NewDataMgr(notes []*models.Note) *DataMgr {
	dm := &DataMgr{
		notes:         notes,
		activeNotePtr: 0,
		noteIndexByID: make(map[uint]int),
	}
	dm.rebuildNoteIndex()
	if len(notes) > 0 {
		dm.activeNoteID = notes[0].ID
	}
	return dm
}

func (dm *DataMgr) rebuildNoteIndex() {
	dm.noteIndexByID = make(map[uint]int, len(dm.notes))
	for i, n := range dm.notes {
		dm.noteIndexByID[n.ID] = i
	}
}

// GetNotes returns all known notes.
func (dm *DataMgr) GetNotes() []*models.Note {
	return dm.notes
}

func (dm *DataMgr) GetActiveNote() *models.Note {
	if len(dm.notes) == 0 || dm.activeNotePtr < 0 || dm.activeNotePtr >= len(dm.notes) {
		return nil
	}
	return dm.notes[dm.activeNotePtr]
}

// GetActiveNotePtr returns the current note pointer
func (dm *DataMgr) GetActiveNotePtr() int {
	return dm.activeNotePtr
}

// GetActiveNoteID returns the current note ID.
func (dm *DataMgr) GetActiveNoteID() uint {
	return dm.activeNoteID
}

// RefreshDataByID updates datamgr with a new note list and restores the
// active note by ID.
func (dm *DataMgr) RefreshDataByID(notes []*models.Note, noteID *uint) {
	dm.notes = notes
	dm.rebuildNoteIndex()

	if len(dm.notes) == 0 {
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		return
	}

	if noteID != nil && dm.SwitchActiveNoteByID(*noteID) {
		return
	}

	dm.SwitchActiveNote(0)
}

// SwitchActiveNote switches to a different note by cursor position.
func (dm *DataMgr) SwitchActiveNote(cursor int) {
	if len(dm.notes) == 0 {
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		return
	}
	if cursor < 0 || cursor >= len(dm.notes) {
		return
	}
	_ = dm.SwitchActiveNoteByID(dm.notes[cursor].ID)
}

// SwitchActiveNoteByID switches active note by stable note ID.
func (dm *DataMgr) SwitchActiveNoteByID(noteID uint) bool {
	if len(dm.notes) == 0 {
		dm.activeNotePtr = 0
		dm.activeNoteID = 0
		return false
	}

	idx, ok := dm.noteIndexByID[noteID]
	if !ok {
		return false
	}

	dm.activeNotePtr = idx
	dm.activeNoteID = noteID
	return true
}

// AddNote adds a note to the note list without switching to it.
func (dm *DataMgr) AddNote(n *models.Note) {
	if n == nil {
		return
	}
	dm.notes = append(dm.notes, n)
	dm.noteIndexByID[n.ID] = len(dm.notes) - 1
}

// FindNoteByID finds a note by ID in O(1). Returns notes pending deletion too.
func (dm *DataMgr) FindNoteByID(id uint) *models.Note {
	idx, ok := dm.noteIndexByID[id]
	if !ok {
		return nil
	}
	return dm.notes[idx]
}

// GetAllNotesByIDDesc returns all non-deleted notes sorted by ID descending.
func (dm *DataMgr) GetAllNotesByIDDesc() []*models.Note {
	notes := make([]*models.Note, 0, len(dm.notes))
	for _, n := range dm.notes {
		if !n.Deleted {
			notes = append(notes, n)
		}
	}
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].ID > notes[j].ID
	})
	return notes
}
