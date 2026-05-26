package models

import (
	"time"

	"gorm.io/gorm"
)

// we will get to save commits for the note.
type NoteCommit struct {
	Patch      string // raw patch, unparsed.
	CommitTime time.Time
}

// Note represents a note entity
type Note struct {
	gorm.Model
	Content        string
	CheckedContent string // latest committed content
	Diff           string
	LastEdit       time.Time
	Highlight      bool      `gorm:"default:false"`
	Private        bool      `gorm:"default:false"`
	Frequency      int       `gorm:"not null;default:0"`
	Branches       []*Branch `gorm:"many2many:branch_notes;constraint:OnDelete:CASCADE;"`
	Commits        []*NoteCommit `gorm:"-"` // not yet persisted
	ThreadID       uint // Foreign key - note belongs to a single thread
}
