package main

import (
	"strings"
	"testing"

	"terva.sh/terva/packages/agent/ext"
)

func TestResultBudgetAccountsForEscapingAndBase64(t *testing.T) {
	normal := ext.TextResult("normal result")
	if got := boundedToolResult(normal); got.IsError || got.Content[0].Text != "normal result" {
		t.Fatal("normal result changed")
	}
	for name, r := range map[string]ext.ToolResult{
		"text":    ext.TextResult(strings.Repeat("x", resultContentBudget+1)),
		"escaped": ext.TextResult(strings.Repeat("\x00", resultContentBudget/2)),
		"image":   {Content: []ext.ToolContent{ext.ImageBytes("image/png", make([]byte, 3<<20))}},
	} {
		t.Run(name, func(t *testing.T) {
			if resultFits(r) {
				t.Fatal("oversized serialized result accepted")
			}
			got := boundedToolResult(r)
			if !got.IsError || !resultFits(got) {
				t.Fatal("fallback is not a bounded error")
			}
		})
	}
}
