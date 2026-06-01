package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haochend413/muninx/sys"
)

type EmbedClient struct {
	BaseURL string
}

func NewEmbedClient(baseURL string) *EmbedClient {
	return &EmbedClient{BaseURL: baseURL}
}

func (c *EmbedClient) Embed(text string) ([]float32, error) {
	body, _ := json.Marshal(map[string]string{
		"text": text,
	})

	resp, err := http.Post(c.BaseURL+"/embed", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := fmt.Errorf("embedding server status: %d", resp.StatusCode)
		sys.LogError(err)
		return nil, err
	}

	var out struct {
		Embedding []float32 `json:"embedding"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out.Embedding, nil
}
