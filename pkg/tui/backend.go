// Package tui provides a terminal user interface for the hive.
package tui

import "context"

// Backend abstracts the hive's services for the TUI.
// Direct implementation wraps the mind; HTTP implementation talks to lovyou.ai.
type Backend interface {
	// Chat sends a message to an agent and returns the response.
	Chat(ctx context.Context, agent, message string) (string, error)

	// History returns the conversation history for the current session.
	History(ctx context.Context) ([]ChatMessage, error)

	// ListConversations returns summaries of past conversations.
	ListConversations(ctx context.Context, limit int) ([]ConversationInfo, error)

	// ResumeConversation switches to a previous conversation.
	ResumeConversation(ctx context.Context, id string) ([]ChatMessage, error)

	// ListAgents returns the available agents and their status.
	ListAgents(ctx context.Context) ([]AgentInfo, error)

	// ListTasks returns work graph tasks.
	ListTasks(ctx context.Context) ([]TaskInfo, error)

	// RecentActivity returns recent agent-to-agent interactions.
	RecentActivity(ctx context.Context, limit int) ([]ActivityEntry, error)

	// Info returns system status (connection, event count, etc.).
	Info(ctx context.Context) (SystemInfo, error)
}

// ConversationInfo describes a past conversation for listing.
type ConversationInfo struct {
	ID        string
	Preview   string // first human message, truncated
	TurnCount int
	StartedAt string // human-readable
}

// ChatMessage is one message in a conversation (used by Backend).
type ChatMessage struct {
	Role    string // "human", "mind", "error", "system"
	Content string
}

// AgentInfo describes an agent visible in the sidebar.
type AgentInfo struct {
	Role   string // "mind", "cto", "pm", etc.
	Name   string // display name
	Status string // "ready", "thinking", "idle"
}

// TaskInfo describes a work graph task.
type TaskInfo struct {
	ID          string
	Title       string
	Status      string // "open", "in_progress", "done"
	AssignedTo  string
	Description string
}

// ActivityEntry is one agent-to-agent interaction.
type ActivityEntry struct {
	From    string // role of sender
	To      string // role of receiver
	Summary string // what happened
	Ago     string // human-readable time ago
}

// SystemInfo holds system-wide status.
type SystemInfo struct {
	MCPConnected bool
	EventCount   int
	TaskCount    int
	CostUSD      float64
	Model        string
}
