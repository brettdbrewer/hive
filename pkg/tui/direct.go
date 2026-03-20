package tui

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lovyou-ai/hive/pkg/mind"
)

// DirectBackend implements Backend by calling the mind directly (in-process).
// This is the local dev path. Production uses an HTTP backend instead.
type DirectBackend struct {
	mind         *mind.Mind
	mcpConnected bool
	model        string
}

// NewDirectBackend creates a backend backed by a local mind instance.
func NewDirectBackend(m *mind.Mind, mcpConnected bool, model string) *DirectBackend {
	return &DirectBackend{
		mind:         m,
		mcpConnected: mcpConnected,
		model:        model,
	}
}

func (b *DirectBackend) Chat(_ context.Context, agent, message string) (string, error) {
	// Suppress stdout during mind.Chat() — the eventgraph intelligence layer
	// uses fmt.Printf for progress heartbeats which corrupts bubbletea's alt screen.
	origStdout := os.Stdout
	if devNull, err := os.Open(os.DevNull); err == nil {
		os.Stdout = devNull
		defer func() {
			os.Stdout = origStdout
			devNull.Close()
		}()
	}
	return b.mind.Chat(context.Background(), message)
}

func (b *DirectBackend) History(_ context.Context) ([]ChatMessage, error) {
	turns := b.mind.History()
	msgs := make([]ChatMessage, len(turns))
	for i, t := range turns {
		msgs[i] = ChatMessage{Role: t.Role, Content: t.Content}
	}
	return msgs, nil
}

func (b *DirectBackend) ListConversations(_ context.Context, limit int) ([]ConversationInfo, error) {
	summaries, err := b.mind.ListConversations(limit)
	if err != nil {
		return nil, err
	}
	infos := make([]ConversationInfo, len(summaries))
	for i, s := range summaries {
		infos[i] = ConversationInfo{
			ID:        s.ID,
			Preview:   s.Preview,
			TurnCount: s.TurnCount,
			StartedAt: formatTimeAgo(s.StartedAt),
		}
	}
	return infos, nil
}

func (b *DirectBackend) ResumeConversation(_ context.Context, id string) ([]ChatMessage, error) {
	// TODO: create a new mind instance with ConversationID set.
	// For now, return empty — full resume requires re-creating the mind.
	return nil, fmt.Errorf("resume not yet implemented")
}

func (b *DirectBackend) ListAgents(_ context.Context) ([]AgentInfo, error) {
	return []AgentInfo{
		{Role: "mind", Name: "Mind", Status: "ready"},
		{Role: "pm", Name: "PM", Status: "idle"},
		{Role: "cto", Name: "CTO", Status: "idle"},
		{Role: "guardian", Name: "Guardian", Status: "idle"},
		{Role: "architect", Name: "Architect", Status: "idle"},
		{Role: "builder", Name: "Builder", Status: "idle"},
	}, nil
}

func (b *DirectBackend) ListTasks(_ context.Context) ([]TaskInfo, error) {
	return nil, nil
}

func (b *DirectBackend) RecentActivity(_ context.Context, limit int) ([]ActivityEntry, error) {
	return nil, nil
}

func (b *DirectBackend) Info(_ context.Context) (SystemInfo, error) {
	return SystemInfo{
		MCPConnected: b.mcpConnected,
		Model:        b.model,
	}, nil
}

func formatTimeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return t.Format("Jan 2")
	}
}
