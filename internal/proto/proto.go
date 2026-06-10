// Package proto implements the zot extension wire protocol (newline-delimited
// JSON over stdio) directly, so this extension depends on nothing but the Go
// standard library — no coupling to the zot module or its version.
//
// It implements the subset a tool-providing extension needs: the hello
// handshake, register_tool, the ready sentinel, tool_call dispatch,
// tool_result replies, and graceful shutdown. See docs/extensions.md in the
// zot repo for the full protocol. Field names here mirror zot's extproto
// package exactly (e.g. is_error, mime_type, snake_case throughout).
package proto

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// ProtocolVersion is the zot extension protocol version this package speaks. It
// is sent by the host in hello_ack; a mismatch is logged (not fatal) so the
// extension keeps working against minor host changes while surfacing a drift.
const ProtocolVersion = 1

// Result is a tool handler's reply. Text is sent back to the model as a single
// text content block; IsError marks the call as failed.
type Result struct {
	Text    string
	IsError bool
}

// Text builds a successful Result.
func Text(s string) Result { return Result{Text: s} }

// Errorf builds a failed Result with a formatted message.
func Errorf(format string, a ...any) Result {
	return Result{Text: fmt.Sprintf(format, a...), IsError: true}
}

// ToolHandler runs when the model invokes a registered tool. args is the raw
// JSON object the model produced; the handler validates it.
type ToolHandler func(args json.RawMessage) Result

type toolDef struct {
	name        string
	description string
	schema      json.RawMessage
	handler     ToolHandler
}

// Host carries the hello_ack fields this extension cares about.
type Host struct {
	ProtocolVersion int
	ZotVersion      string
	DataDir         string
	ExtensionDir    string
	Provider        string
	Model           string
	CWD             string
}

// Extension is one tool-providing extension. Construct with New, register
// tools, then call Run.
type Extension struct {
	name    string
	version string

	in      io.Reader
	out     io.Writer
	writeMu sync.Mutex

	mu    sync.Mutex
	tools []toolDef
	host  Host
}

// New constructs an Extension that talks to zot over stdin/stdout.
func New(name, version string) *Extension {
	return &Extension{name: name, version: version, in: os.Stdin, out: os.Stdout}
}

// Tool registers an LLM-callable tool. Call before Run. schema is a JSON Schema
// object (same shape Anthropic/OpenAI accept).
func (e *Extension) Tool(name, description string, schema json.RawMessage, h ToolHandler) {
	e.mu.Lock()
	e.tools = append(e.tools, toolDef{name, description, schema, h})
	e.mu.Unlock()
}

// Host returns the info zot sent in hello_ack. Zero value until the handshake
// completes (which happens before any tool_call).
func (e *Extension) Host() Host {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.host
}

// Logf writes a debug line to stderr, which zot captures to
// $ZOT_HOME/logs/ext-<name>.log. Never write to stdout — that's the wire.
func (e *Extension) Logf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "["+e.name+"] "+format+"\n", a...)
}

func (e *Extension) send(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	e.writeMu.Lock()
	defer e.writeMu.Unlock()
	_, _ = e.out.Write(append(b, '\n'))
}

// Run sends the hello + registrations, then serves tool calls until zot closes
// stdin or sends shutdown. Blocks until then.
func (e *Extension) Run() error {
	e.send(map[string]any{
		"type": "hello", "name": e.name, "version": e.version,
		"capabilities": []string{"tools"},
	})
	e.mu.Lock()
	tools := append([]toolDef(nil), e.tools...)
	e.mu.Unlock()
	for _, t := range tools {
		e.send(map[string]any{
			"type": "register_tool", "name": t.name,
			"description": t.description, "schema": t.schema,
		})
	}
	e.send(map[string]any{"type": "ready"})

	sc := bufio.NewScanner(e.in)
	sc.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	for sc.Scan() {
		var f struct {
			Type            string          `json:"type"`
			ID              string          `json:"id"`
			Name            string          `json:"name"`
			Args            json.RawMessage `json:"args"`
			ProtocolVersion int             `json:"protocol_version"`
			ZotVersion      string          `json:"zot_version"`
			DataDir         string          `json:"data_dir"`
			ExtensionDir    string          `json:"extension_dir"`
			Provider        string          `json:"provider"`
			Model           string          `json:"model"`
			CWD             string          `json:"cwd"`
		}
		if err := json.Unmarshal(sc.Bytes(), &f); err != nil {
			e.Logf("bad frame from host: %v", err)
			continue
		}
		switch f.Type {
		case "hello_ack":
			e.mu.Lock()
			e.host = Host{
				ProtocolVersion: f.ProtocolVersion,
				ZotVersion:      f.ZotVersion,
				DataDir:         f.DataDir,
				ExtensionDir:    f.ExtensionDir,
				Provider:        f.Provider,
				Model:           f.Model,
				CWD:             f.CWD,
			}
			e.mu.Unlock()
			if f.ProtocolVersion != 0 && f.ProtocolVersion != ProtocolVersion {
				e.Logf("warning: host speaks protocol_version %d but this extension implements %d (zot %s); proceeding, but behavior may drift",
					f.ProtocolVersion, ProtocolVersion, f.ZotVersion)
			}
		case "tool_call":
			h := e.handlerFor(f.Name)
			if h == nil {
				e.sendToolResult(f.ID, Errorf("no handler for tool %q", f.Name))
				continue
			}
			// Own goroutine so a slow fetch doesn't block other calls.
			go func(id string, args json.RawMessage) {
				defer func() {
					if r := recover(); r != nil {
						e.sendToolResult(id, Errorf("panic: %v", r))
					}
				}()
				e.sendToolResult(id, h(args))
			}(f.ID, f.Args)
		case "shutdown":
			e.send(map[string]any{"type": "shutdown_ack"})
			return nil
		}
	}
	return sc.Err()
}

func (e *Extension) handlerFor(name string) ToolHandler {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, t := range e.tools {
		if t.name == name {
			return t.handler
		}
	}
	return nil
}

func (e *Extension) sendToolResult(id string, r Result) {
	e.send(map[string]any{
		"type":     "tool_result",
		"id":       id,
		"content":  []map[string]any{{"type": "text", "text": r.Text}},
		"is_error": r.IsError,
	})
}
