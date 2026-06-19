package models

import (
	"gorm.io/gorm"
)

// BeforeSave is a GORM hook that ensures LastEdit is set to CreatedAt if it is null.
func (n *Note) BeforeSave(tx *gorm.DB) (err error) {
	if n.LastEdit.IsZero() {
		n.LastEdit = n.CreatedAt
	}
	return nil
}
