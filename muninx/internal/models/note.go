package models

import (
	"time"

	"gorm.io/gorm"
)

// NoteCommit is a persisted snapshot of a note's content, recorded as a
// patch from the previous snapshot (or from empty, for the first commit).
type NoteCommit struct {
	gorm.Model
	NoteID     uint   // Foreign key - commit belongs to a single note
	Patch      string // raw patch, unparsed.
	CommitTime time.Time
	Archived   bool `gorm:"default:false"`
}

// Note represents a note entity
type Note struct {
	gorm.Model
	Content        string
	CheckedContent string // latest committed content
	Diff           string
	LastEdit       time.Time
	Highlight      bool          `gorm:"default:false"`
	Private        bool          `gorm:"default:false"`
	Frequency      int           `gorm:"not null;default:0"`
	Commits        []*NoteCommit `gorm:"constraint:OnDelete:CASCADE;"`
	Deleted        bool          `gorm:"-"` // pending deletion, not yet synced
	Archived       bool
}
