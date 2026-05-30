package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// searxng queries a self-hosted SearXNG instance's JSON API. The instance must
// have `json` listed under search.formats in settings.yml, else it returns 403.
type searxng struct{ base string }

func (s *searxng) Search(ctx context.Context, query string, count int) ([]Result, error) {
	q := url.Values{}
	q.Set("q", query)
	q.Set("format", "json")
	endpoint := s.base + "/search?" + q.Encode()

	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("searxng: HTTP %d (is `json` enabled in search.formats?)", resp.StatusCode)
	}

	var out struct {
		Results []struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Content string `json:"content"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("searxng: decode: %w", err)
	}
	limit := clampCount(count)
	res := make([]Result, 0, limit)
	for i, r := range out.Results {
		if i >= limit {
			break
		}
		res = append(res, Result{Title: r.Title, URL: r.URL, Snippet: r.Content})
	}
	return res, nil
}
