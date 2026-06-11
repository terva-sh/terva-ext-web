package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// searxng queries a self-hosted SearXNG instance's JSON API. The instance must
// have `json` listed under search.formats in settings.yml, else it returns 403.
type searxng struct {
	base   string
	client *http.Client
}

func (s *searxng) Search(ctx context.Context, query string, count int) ([]Result, error) {
	baseURL, err := url.Parse(s.base)
	if err != nil {
		return nil, fmt.Errorf("searxng: invalid base URL %q: %w", s.base, err)
	}
	baseURL = baseURL.JoinPath("search")
	q := url.Values{}
	q.Set("q", query)
	q.Set("format", "json")
	baseURL.RawQuery = q.Encode()
	endpoint := baseURL.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Limit response body to 1 MiB to prevent memory exhaustion.
	const maxBody = 1 << 20
	limited := io.LimitReader(resp.Body, maxBody)
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
	if err := json.NewDecoder(limited).Decode(&out); err != nil {
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
