package embedder

import (
	"github.com/haochend413/muninx/internal/clients"
	"github.com/haochend413/muninx/internal/db"
	"github.com/haochend413/muninx/internal/models"
)

type Embedder struct {
	client *clients.EmbedClient
	db     *db.DB
}

func NewEmbedder(embclient *clients.EmbedClient, db *db.DB) *Embedder {
	return &Embedder{client: embclient, db: db}
}

func (e *Embedder) EmbedNote(noteID uint, text string) error {
	vec, err := e.client.Embed(text)
	if err != nil {
		return err
	}
	return e.db.UpsertNoteEmbedding(noteID, vec)
}

// ReEmbedAll re-embeds every provided note. No-ops silently if the embed client
// is unavailable.
func (e *Embedder) ReEmbedAll(notes []*models.Note) {
	if e.client == nil {
		return
	}
	for _, n := range notes {
		if n.Content != "" {
			e.EmbedNote(n.ID, n.Content) //nolint:errcheck — best-effort
		}
	}
}

func (e *Embedder) FetchRelated(noteID uint, k int) []models.Note {
	embedding, err := e.db.GetNoteEmbedding(noteID)
	if err != nil {
		return nil
	}
	if embedding == nil {
		return nil
	}
	notes, err := e.db.SearchRelatedNotes(embedding, k)
	if err != nil {
		return nil
	}

	return notes
}
