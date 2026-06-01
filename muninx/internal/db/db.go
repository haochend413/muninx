package db

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/haochend413/muninx/internal/clients"
	"github.com/haochend413/muninx/internal/models"
	"github.com/haochend413/muninx/sys"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB wraps the GORM database connection
type DB struct {
	Conn        *gorm.DB
	EmbedClient *clients.EmbedClient
}

// NewDB initializes a new database connection and migrates schema
func NewDB(path string, embedClient *clients.EmbedClient) (*DB, error) {
	sqlite_vec.Auto()

	// if not exist, create all dirs
	_, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		// Config file doesn't exist, create directory and config file with defaults
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			sys.LogError(err)
			return nil, err
		}

	}

	conn, err := gorm.Open(sqlite.Open(path+"?_journal_mode=WAL&_busy_timeout=5000"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate schema
	err = conn.AutoMigrate(&models.Note{}, &models.Thread{}, &models.Branch{})
	if err != nil {
		return nil, err
	}

	database := &DB{
		Conn:        conn,
		EmbedClient: embedClient,
	}

	if err := database.InitVectorTable(); err != nil {
		return nil, err
	}

	return database, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	sqlDB, err := d.Conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CreateThread inserts a new thread immediately and sets thread.ID from the
// database-assigned autoincrement value.
func (d *DB) CreateThread(thread *models.Thread) error {
	return d.Conn.Omit("Branches").Create(thread).Error
}

// CreateBranch inserts a new branch immediately and sets branch.ID.
func (d *DB) CreateBranch(branch *models.Branch) error {
	return d.Conn.Omit("Notes").Create(branch).Error
}

// CreateNote inserts a new note immediately, sets note.ID, and writes the
// branch_notes join-table row so the association is persisted from the start.
func (d *DB) CreateNote(note *models.Note) error {
	note.Content = strings.TrimSpace(note.Content)
	if err := d.Conn.Omit("Branches").Create(note).Error; err != nil {
		return err
	}
	return d.Conn.Model(note).Association("Branches").Replace(note.Branches)
}
