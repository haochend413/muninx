package app

import (
	"sync"
	"time"

	"github.com/haochend413/muninx/internal/app/data"
	editstack "github.com/haochend413/muninx/internal/app/editStack"
	"github.com/haochend413/muninx/internal/app/embedder"
	"github.com/haochend413/muninx/internal/clients"
	"github.com/haochend413/muninx/internal/db"
	"github.com/haochend413/muninx/internal/models"
	"github.com/haochend413/muninx/state"
	"github.com/haochend413/muninx/sys"
)

// App encapsulates application logic and state.
type App struct {
	db           *db.DB
	dataMgr      *data.DataMgr
	editMgr      *editstack.EditMgr
	embedder     *embedder.Embedder
	Synced       bool
	mutex        sync.Mutex
	deletedStack []uint // IDs pending deletion, most-recent-last, for UndoDelete
}

// NewApp creates a new App, loading all data from the database.
func NewApp(dbConn *db.DB, AppState *state.AppState, embedClient *clients.EmbedClient) *App {
	app := &App{
		db:       dbConn,
		dataMgr:  &data.DataMgr{},
		editMgr:  editstack.NewEditMgr(),
		embedder: embedder.NewEmbedder(embedClient, dbConn),
		Synced:   true,
	}
	app.loadData()
	return app
}

// GetDataMgr returns the data manager.
func (a *App) GetDataMgr() *data.DataMgr {
	return a.dataMgr
}

// GetEditMap returns the current edit map.
func (a *App) GetEditMap() map[editstack.EditKey]*editstack.Edit {
	return a.editMgr.EditMap
}

// loadData loads all notes from the database.
func (a *App) loadData() {
	notes, err := a.db.SyncData(
		[]*models.Note{},
		make(map[editstack.EditKey]*editstack.Edit),
	)
	if err != nil {
		sys.LogError(err)
		panic(err)
	}
	a.dataMgr = data.NewDataMgr(notes)
}

// CreateNewNote creates a note in the database immediately (getting a real
// autoincrement ID) and adds it to the in-memory list.
func (a *App) CreateNewNote() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := &models.Note{Content: ""}
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	if err := a.db.CreateNote(note); err != nil {
		sys.LogError(err)
		return
	}
	a.Synced = false
	a.dataMgr.AddNote(note)
}

// GetActiveNoteList returns all non-deleted notes.
func (a *App) GetActiveNoteList() []*models.Note {
	all := a.dataMgr.GetNotes()
	visible := make([]*models.Note, 0, len(all))
	for _, n := range all {
		if !n.Deleted {
			visible = append(visible, n)
		}
	}
	return visible
}

// ReEmbedAllNotes re-embeds every note in the dataset. Runs synchronously;
// call from a goroutine (tea.Cmd) to avoid blocking the UI.
func (a *App) ReEmbedAllNotes() {
	if a.embedder == nil {
		return
	}
	a.embedder.ReEmbedAll(a.dataMgr.GetAllNotesByIDDesc())
}

// FetchRelatedNotes returns the k most semantically related notes for a given noteID.
func (a *App) FetchRelatedNotes(noteID uint, k int) []models.Note {
	if a.embedder == nil {
		return nil
	}
	return a.embedder.FetchRelated(noteID, k)
}

// SyncWithDatabase flushes pending updates and deletes to the database.
func (a *App) SyncWithDatabase() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	notes := a.dataMgr.GetNotes()
	editMapCopy := make(map[editstack.EditKey]*editstack.Edit)
	for k, v := range a.editMgr.EditMap {
		editMapCopy[k] = v
	}

	noteID := a.dataMgr.GetActiveNoteID()

	updatedNotes, err := a.db.SyncData(notes, editMapCopy)
	if err != nil {
		sys.LogError(err)
		return
	}

	a.dataMgr.RefreshDataByID(updatedNotes, &noteID)
	a.editMgr.ClearOnSync()
	a.deletedStack = nil
	a.Synced = true
}
