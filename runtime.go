package main

import (
	"encoding/json"
	"sync"

	"terva-ext-web/internal/config"
	"terva-ext-web/internal/fetch"
	"terva-ext-web/internal/search"
	"terva.sh/terva/packages/agent/ext"
)

// Each invocation retains one immutable runtime. Old calls may finish, but
// their cache belongs only to their runtime, never to a new configuration.
type webRuntime struct {
	fetcher                *fetch.Client
	provider               search.Provider
	configErr, providerErr error
}

type runtimeStore struct {
	mu      sync.Mutex
	current *webRuntime
	last    string
	read    func() ext.Config
	notify  func(string)
}

func (s *runtimeStore) snapshot() *webRuntime {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Read under the lock so a delayed config callback cannot replace a newer
	// snapshot with an older event payload. The SDK owns the current map.
	values := s.read()
	encoded, _ := json.Marshal(values)
	key := string(encoded)
	if s.current != nil && key == s.last {
		return s.current
	}
	s.last = key
	cfg, err := config.Resolve(values)
	if err != nil {
		s.notify("web configuration rejected: " + err.Error())
		if s.current == nil || s.current.configErr != nil {
			s.current = &webRuntime{configErr: err, providerErr: err}
		}
		return s.current
	}
	f := fetch.New(cfg, fetch.ParseAllowList(cfg.AllowLocalHosts))
	provider, providerErr := search.New(cfg, f.HTTPClient())
	s.current = &webRuntime{fetcher: f, provider: provider, providerErr: providerErr}
	return s.current
}
