package app

import (
	"log"
	"sync"
	"time"

	"github.com/haochend413/muninx/internal/app/data"
	editstack "github.com/haochend413/muninx/internal/app/editStack"
	"github.com/haochend413/muninx/internal/app/embedder"
	"github.com/haochend413/muninx/internal/clients"
	"github.com/haochend413/muninx/internal/db"
	"github.com/haochend413/muninx/internal/models"
	"github.com/haochend413/muninx/state"
)

// App encapsulates application logic and state.
type App struct {
	db       *db.DB
	dataMgr  *data.DataMgr
	editMgr  *editstack.EditMgr
	embedder *embedder.Embedder
	Synced   bool
	mutex    sync.Mutex
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

// GetNoteEditStack returns the note edit stack.
func (a *App) GetNoteEditStack() []*editstack.NoteEdit {
	return a.editMgr.NoteEditStack
}

// loadData loads the full thread tree from the database.
func (a *App) loadData() {
	threads, err := a.db.SyncData(
		[]*models.Thread{},
		make(map[editstack.EditKey]*editstack.Edit),
	)
	if err != nil {
		log.Panic(err)
	}
	a.dataMgr = data.NewDataMgr(threads)
}

// CreateNewThread creates a thread in the database immediately and adds it to
// the in-memory tree.
func (a *App) CreateNewThread(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := &models.Thread{Name: ""}
	thread.CreatedAt = time.Now()
	thread.UpdatedAt = time.Now()

	if err := a.db.CreateThread(thread); err != nil {
		log.Printf("Error creating thread: %v", err)
		return
	}
	a.Synced = false
	a.dataMgr.AddThread(thread)
}

// CreateNewBranch creates a branch under the active thread immediately in the
// database and adds it to the in-memory tree.
func (a *App) CreateNewBranch(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.dataMgr.GetActiveThread()
	if thread == nil {
		log.Printf("Cannot create branch: no active thread")
		return
	}

	branch := &models.Branch{Name: ""}
	branch.CreatedAt = time.Now()
	branch.UpdatedAt = time.Now()
	branch.ThreadID = thread.ID

	if err := a.db.CreateBranch(branch); err != nil {
		log.Printf("Error creating branch: %v", err)
		return
	}
	a.Synced = false
	a.dataMgr.AddBranch(branch)
}

// CreateNewNote creates a note under the active branch immediately in the
// database (getting a real autoincrement ID) and adds it to the in-memory tree.
func (a *App) CreateNewNote(link *models.Superlink) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	thread := a.dataMgr.GetActiveThread()
	branch := a.dataMgr.GetActiveBranch()
	if thread == nil {
		log.Printf("Cannot create note: no active thread")
		return
	}
	if branch == nil {
		log.Printf("Cannot create note: no active branch")
		return
	}

	note := &models.Note{Content: ""}
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	note.ThreadID = thread.ID
	note.Branches = []*models.Branch{branch}

	if err := a.db.CreateNote(note); err != nil {
		log.Printf("Error creating note: %v", err)
		return
	}
	a.Synced = false
	a.dataMgr.AddNote(note)
}

func (a *App) GetThreadList() []*models.Thread {
	if a == nil {
		log.Panic("null app")
	}
	return a.dataMgr.GetThreads()
}

func (a *App) GetActiveBranchList() []*models.Branch {
	if a == nil {
		log.Panic("null app")
	}
	return a.dataMgr.GetActiveBranchList()
}

func (a *App) GetActiveNoteList() []*models.Note {
	if a == nil {
		log.Panic("null app")
	}
	return a.dataMgr.GetActiveNoteList()
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

	threads := a.dataMgr.GetThreads()
	editMapCopy := make(map[editstack.EditKey]*editstack.Edit)
	for k, v := range a.editMgr.EditMap {
		editMapCopy[k] = v
	}

	threadID := a.dataMgr.GetActiveThreadID()
	branchID := a.dataMgr.GetActiveBranchID()
	noteID := a.dataMgr.GetActiveNoteID()

	updatedThreads, err := a.db.SyncData(threads, editMapCopy)
	if err != nil {
		log.Printf("Error syncing with database: %v", err)
		return
	}

	a.dataMgr.RefreshDataByID(updatedThreads, &threadID, &branchID, &noteID)
	a.editMgr.ClearOnSync()
	a.Synced = true
}
