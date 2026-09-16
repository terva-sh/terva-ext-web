package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// tavily calls the Tavily Search API (https://docs.tavily.com). Auth is a
// bearer token; results already carry summarized, agent-ready content.
type tavily struct {
	key    string
	client *http.Client
}

func (t *tavily) Search(ctx context.Context, q Query) (results []Result, resultErr error) {
	// A provider response or transport error can echo its authorization input.
	// Redact the configured key before anything reaches tools or diagnostics.
	defer func() {
		if t.key == "" {
			return
		}
		redact := func(s string) string { return strings.ReplaceAll(s, t.key, "[redacted]") }
		if resultErr != nil {
			resultErr = credentialSafeError{cause: resultErr, message: redact(resultErr.Error())}
		}
		for i := range results {
			results[i].Title = redact(results[i].Title)
			results[i].URL = redact(results[i].URL)
			results[i].Snippet = redact(results[i].Snippet)
			results[i].Published = redact(results[i].Published)
		}
	}()
	payload := map[string]any{
		"query":        q.Text,
		"max_results":  clampCount(q.Count),
		"search_depth": "basic",
	}
	if q.Depth == "advanced" {
		payload["search_depth"] = "advanced"
	}
	if q.Freshness != "" {
		payload["time_range"] = q.Freshness
	}
	if len(q.IncludeDomains) > 0 {
		payload["include_domains"] = q.IncludeDomains
	}
	if len(q.ExcludeDomains) > 0 {
		payload["exclude_domains"] = q.ExcludeDomains
	}
	body, _ := json.Marshal(payload)
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
		hint := ""
		switch resp.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			hint = " (invalid or missing API key — check TAVILY_API_KEY / tavily_api_key)"
		case http.StatusTooManyRequests:
			hint = " (Tavily rate or credit limit — wait before retrying)"
		}
		return nil, fmt.Errorf("tavily: HTTP %d%s", resp.StatusCode, hint)
	}

	var out struct {
		Results []struct {
			Title         string `json:"title"`
			URL           string `json:"url"`
			Content       string `json:"content"`
			PublishedDate string `json:"published_date"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("tavily: decode: invalid search response")
	}
	res := make([]Result, 0, len(out.Results))
	for _, r := range out.Results {
		res = append(res, Result{Title: r.Title, URL: r.URL, Snippet: r.Content, Published: r.PublishedDate})
	}
	return res, nil
}

// Keep error identity (timeouts/SSRF) while keeping the printable form safe.
type credentialSafeError struct {
	cause   error
	message string
}

func (e credentialSafeError) Error() string { return e.message }
func (e credentialSafeError) Unwrap() error { return e.cause }
