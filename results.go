package main

import (
	"encoding/json"

	"terva.sh/terva/packages/agent/ext"
	"terva.sh/terva/packages/agent/extproto"
)

// The supported host generates small correlation IDs; reserve 4 KiB for the
// result envelope and ID. Count serialized content, not bytes before escaping
// or base64 encoding. This does not lower the separate file-download limits.
const resultContentBudget = extproto.MaxFrameBytes - 4096

func resultFits(r ext.ToolResult) bool {
	blocks := make([]extproto.ContentBlock, len(r.Content))
	for i, c := range r.Content {
		blocks[i] = extproto.ContentBlock{Type: c.Type, Text: c.Text, MimeType: c.MimeType, Data: c.Data}
	}
	data, err := json.Marshal(blocks)
	return err == nil && len(data) <= resultContentBudget
}
func boundedToolResult(r ext.ToolResult) ext.ToolResult {
	if resultFits(r) {
		return r
	}
	return ext.TextErrorResult("Result exceeds the host message limit. Reduce max_chars or max_dimension; use web_fetch_raw for large text or inject:false with save_path for images. Any requested save may already have completed.")
}
func registerBoundedTool(e *ext.Extension, name, description string, schema json.RawMessage, handler ext.ToolHandler, options ...ext.ToolOption) {
	e.Tool(name, description, schema, func(args json.RawMessage) ext.ToolResult { return boundedToolResult(handler(args)) }, options...)
}
func registerBoundedCommand(e *ext.Extension, name, description string, handler ext.CommandHandler) {
	e.Command(name, description, func(args string) ext.Response {
		response := handler(args)
		// SDK Response has more fields than the cache command uses; marshaling it
		// conservatively includes all of them rather than undercounting the wire.
		data, err := json.Marshal(response)
		if err != nil || len(data) > resultContentBudget {
			return ext.Errorf("Cache response exceeds the host message limit; use /web-cache clear.")
		}
		return response
	})
}
