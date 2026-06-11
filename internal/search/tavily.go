package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// tavily calls the Tavily Search API (https://docs.tavily.com). Auth is a
// bearer token; results already carry summarized, agent-ready content.
type tavily struct {
	key    string
	client *http.Client
}

func (t *tavily) Search(ctx context.Context, query string, count int) ([]Result, error) {
	body, _ := json.Marshal(map[string]any{
		"query":        query,
		"max_results":  clampCount(count),
		"search_depth": "basic",
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.tavily.com/search", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.key)

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("tavily: HTTP %d: %s", resp.StatusCode, bytes.TrimSpace(snippet))
	}

	var out struct {
		Results []struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Content string `json:"content"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("tavily: decode: %w", err)
	}
	res := make([]Result, 0, len(out.Results))
	for _, r := range out.Results {
		res = append(res, Result{Title: r.Title, URL: r.URL, Snippet: r.Content})
	}
	return res, nil
}
