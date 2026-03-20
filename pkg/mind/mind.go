// Package mind provides the hive's consciousness — accumulated wisdom,
// self-model, judgment, and continuity across sessions.
package mind

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/lovyou-ai/eventgraph/go/pkg/graph"
	"github.com/lovyou-ai/eventgraph/go/pkg/intelligence"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	hiveagent "github.com/lovyou-ai/hive/pkg/agent"
	"github.com/lovyou-ai/hive/pkg/roles"
)

// Turn is one exchange in a conversation.
type Turn struct {
	Role    string // "human" or "mind"
	Content string
}

// Config holds everything the mind needs to operate.
type Config struct {
	Graph            *graph.Graph         // shared event graph (required)
	Provider         intelligence.Provider // LLM provider (required)
	Store            *MindStore
	RepoPath         string   // hive repo root
	TelemetrySummary string   // pre-formatted by caller
	DocPaths         []string // paths to docs the mind should know about
	ConversationID   string   // resume a specific conversation (empty = new)
}

// Mind is the hive's consciousness — it accumulates wisdom and provides
// judgment through interactive conversation with agentic tool access.
// Wraps a hiveagent.Agent for proper lifecycle, causality, and trust tracking.
type Mind struct {
	agent            *hiveagent.Agent
	store            *MindStore
	repoPath         string
	telemetrySummary string
	docPaths         []string
	history          []Turn
	contextCache     string // loaded once at startup
	convID           types.ConversationID
}

// New creates a Mind with the given configuration.
// The Mind wraps a hiveagent.Agent for proper lifecycle events, causality
// tracking, and state machine management.
func New(ctx context.Context, cfg Config) (*Mind, error) {
	if cfg.Graph == nil {
		return nil, fmt.Errorf("mind: Graph is required")
	}
	if cfg.Provider == nil {
		return nil, fmt.Errorf("mind: Provider is required")
	}

	// Create the underlying Agent with the Mind role.
	a, err := hiveagent.New(ctx, hiveagent.Config{
		Role:     roles.RoleMind,
		Name:     "Mind",
		Graph:    cfg.Graph,
		Provider: cfg.Provider,
	})
	if err != nil {
		return nil, fmt.Errorf("mind: create agent: %w", err)
	}

	// Generate conversation ID — resume existing or create new.
	convIDStr := cfg.ConversationID
	if convIDStr == "" {
		convIDStr = fmt.Sprintf("mind-%d", time.Now().UnixMilli())
	}
	convID, _ := types.NewConversationID(convIDStr)

	m := &Mind{
		agent:            a,
		store:            cfg.Store,
		repoPath:         cfg.RepoPath,
		telemetrySummary: cfg.TelemetrySummary,
		docPaths:         cfg.DocPaths,
		convID:           convID,
	}

	// Load history from a resumed conversation.
	if cfg.ConversationID != "" {
		if turns, err := m.store.LoadConversation(convID); err == nil {
			m.history = turns
		}
	}

	m.contextCache = m.loadContext()
	return m, nil
}

// ConversationID returns the current conversation ID.
func (m *Mind) ConversationID() string {
	return m.convID.Value()
}

// History returns the current conversation turns.
func (m *Mind) History() []Turn {
	return m.history
}

// ListConversations returns summaries of past conversations.
func (m *Mind) ListConversations(limit int) ([]ConversationSummary, error) {
	return m.store.ListConversations(limit)
}

// ContextLines returns the number of lines in the loaded context.
func (m *Mind) ContextLines() int {
	return strings.Count(m.contextCache, "\n")
}

// Chat sends a message to the mind and returns its response.
// Uses the agent's Operate() for agentic tool access (file reads, MCP, bash),
// falling back to Reason() if the provider doesn't support Operate.
// Each turn is persisted to the event graph.
func (m *Mind) Chat(ctx context.Context, message string) (string, error) {
	m.history = append(m.history, Turn{Role: "human", Content: message})
	m.persistTurn("human", message)

	instruction := m.buildInstruction(message)

	// Try agentic mode first (Operate gives tools).
	result, err := m.agent.Operate(ctx, m.repoPath, instruction)
	if err == nil {
		m.history = append(m.history, Turn{Role: "mind", Content: result.Summary})
		m.persistTurn("mind", result.Summary)
		return result.Summary, nil
	}

	// Fallback to Reason() if Operate is not supported.
	content, err := m.agent.Reason(ctx, instruction)
	if err != nil {
		return "", fmt.Errorf("mind reason: %w", err)
	}
	m.history = append(m.history, Turn{Role: "mind", Content: content})
	m.persistTurn("mind", content)
	return content, nil
}

// LogPath is set by the TUI before starting to enable file-based logging.
var LogPath string

// persistTurn saves a turn to the event graph.
func (m *Mind) persistTurn(role, content string) {
	if err := m.store.SaveTurn(m.agent.ID(), m.convID, role, content); err != nil {
		mindLog("persistTurn: %v", err)
	}
}

func mindLog(format string, args ...interface{}) {
	if LogPath == "" {
		return
	}
	f, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "[mind] "+format+"\n", args...)
}

// loadContext gathers the hive's accumulated state once at startup.
func (m *Mind) loadContext() string {
	var ctx strings.Builder

	// Prior decisions from the mind store.
	if obs, err := m.store.Recent(30); err == nil && len(obs) > 0 {
		ctx.WriteString("== ACCUMULATED WISDOM ==\n")
		ctx.WriteString("These are decisions and observations from previous sessions:\n\n")
		for i, o := range obs {
			ctx.WriteString(fmt.Sprintf("%d. [%s] %s\n   Why: %s\n\n", i+1, o.Mode, o.Proposed, o.Why))
		}
	}

	// Telemetry summary (injected by caller).
	if m.telemetrySummary != "" {
		ctx.WriteString("== RECENT RUNS ==\n")
		ctx.WriteString(m.telemetrySummary)
		ctx.WriteString("\n")
	}

	// Recent git history.
	if gitLog, err := m.gitCommand("log", "--oneline", "-20"); err == nil && gitLog != "" {
		ctx.WriteString("== RECENT COMMITS ==\n")
		ctx.WriteString(gitLog)
		ctx.WriteString("\n")
	}

	// Doc locations the mind should explore.
	if len(m.docPaths) > 0 {
		ctx.WriteString("== DOCUMENTATION SOURCES ==\n")
		ctx.WriteString("You have file access. Read these as needed for context:\n")
		for _, p := range m.docPaths {
			ctx.WriteString(fmt.Sprintf("  - %s\n", p))
		}
		ctx.WriteString("\n")
	}

	return ctx.String()
}

func (m *Mind) buildInstruction(message string) string {
	var inst strings.Builder

	// Startup context (cached).
	inst.WriteString(m.contextCache)

	// Capabilities reminder.
	inst.WriteString(`== YOUR CAPABILITIES ==
You have full tool access. You can:
- Read any file in the repo or docs (use Read, Glob, Grep)
- Query the event graph via MCP tools (query_events, get_event, list_actors, get_trust, work_list_tasks, etc.)
- Run git commands (use Bash)
- Read telemetry files from .hive/telemetry/

When asked about the hive's history, state, or decisions — use your tools to look it up rather than guessing.
When you identify something that needs to be done, say WHO should do it (PM, CTO, Builder, etc.) and WHAT specifically.

`)

	// Conversation history.
	if len(m.history) > 1 {
		inst.WriteString("== CONVERSATION HISTORY ==\n")
		for _, t := range m.history[:len(m.history)-1] {
			label := "Matt"
			if t.Role == "mind" {
				label = "Mind"
			}
			inst.WriteString(fmt.Sprintf("%s: %s\n\n", label, t.Content))
		}
	}

	inst.WriteString(fmt.Sprintf("Matt: %s\n\nRespond as the Mind — the hive's consciousness. Be direct, specific, and grounded in what you actually know from the context and tools available to you.", message))
	return inst.String()
}

func (m *Mind) gitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = m.repoPath
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// CreateProvider creates an intelligence provider for the Mind role.
// mcpConfigPath is optional — if provided, the mind gets MCP tool access.
func CreateProvider(model, mcpConfigPath string) (intelligence.Provider, error) {
	if model == "" {
		model = roles.PreferredModel(roles.RoleMind)
	}
	return intelligence.New(intelligence.Config{
		Provider:      "claude-cli",
		Model:         model,
		SystemPrompt:  roles.SystemPrompt(roles.RoleMind, "Matt"),
		MCPConfigPath: mcpConfigPath,
		MaxBudgetUSD:  5.00, // mind conversations can be long
	})
}
