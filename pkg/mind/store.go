package mind

import (
	"fmt"
	"strings"
	"time"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// Observation is a recorded CTO decision: what was proposed and why.
type Observation struct {
	Proposed string
	Why      string
	Mode     string
}

// MindStore persists CTO observations as signed events on the event graph.
type MindStore struct {
	store   store.Store
	factory *event.EventFactory
	signer  event.Signer
}

// NewMindStore creates a MindStore backed by the given event graph infrastructure.
func NewMindStore(s store.Store, factory *event.EventFactory, signer event.Signer) *MindStore {
	return &MindStore{store: s, factory: factory, signer: signer}
}

// Save records a CTO observation as a mind.observation.created event.
func (ms *MindStore) Save(source types.ActorID, obs Observation, causes []types.EventID, convID types.ConversationID) error {
	content := ObservationCreatedContent{
		Proposed: obs.Proposed,
		Why:      obs.Why,
		Mode:     obs.Mode,
	}
	ev, err := ms.factory.Create(EventTypeObservationCreated, source, content, causes, convID, ms.store, ms.signer)
	if err != nil {
		return fmt.Errorf("create observation event: %w", err)
	}
	if _, err := ms.store.Append(ev); err != nil {
		return fmt.Errorf("append observation event: %w", err)
	}
	return nil
}

// Recent returns the last n observations, newest first.
func (ms *MindStore) Recent(n int) ([]Observation, error) {
	page, err := ms.store.ByType(EventTypeObservationCreated, n, types.None[types.Cursor]())
	if err != nil {
		return nil, fmt.Errorf("query observations: %w", err)
	}
	obs := make([]Observation, 0, len(page.Items()))
	for _, ev := range page.Items() {
		c, ok := ev.Content().(ObservationCreatedContent)
		if !ok {
			continue
		}
		obs = append(obs, Observation{
			Proposed: c.Proposed,
			Why:      c.Why,
			Mode:     c.Mode,
		})
	}
	return obs, nil
}

// SaveTurn persists a single conversation turn to the event graph.
func (ms *MindStore) SaveTurn(source types.ActorID, convID types.ConversationID, role, content string) error {
	causes, err := ms.headCause()
	if err != nil {
		return fmt.Errorf("get head for cause: %w", err)
	}
	c := TurnCreatedContent{Role: role, Content: content}
	ev, err := ms.factory.Create(EventTypeTurnCreated, source, c, causes, convID, ms.store, ms.signer)
	if err != nil {
		return fmt.Errorf("create turn event: %w", err)
	}
	if _, err := ms.store.Append(ev); err != nil {
		return fmt.Errorf("append turn event: %w", err)
	}
	return nil
}

// headCause returns the current head event ID as a single-element cause slice.
func (ms *MindStore) headCause() ([]types.EventID, error) {
	head, err := ms.store.Head()
	if err != nil {
		return nil, err
	}
	if head.IsNone() {
		return nil, fmt.Errorf("store has no head event (not bootstrapped)")
	}
	return []types.EventID{head.Unwrap().ID()}, nil
}

// LoadConversation returns all turns for a conversation, oldest first.
func (ms *MindStore) LoadConversation(convID types.ConversationID) ([]Turn, error) {
	page, err := ms.store.ByConversation(convID, 1000, types.None[types.Cursor]())
	if err != nil {
		return nil, fmt.Errorf("query conversation: %w", err)
	}
	var turns []Turn
	for _, ev := range page.Items() {
		c, ok := ev.Content().(TurnCreatedContent)
		if !ok {
			continue
		}
		turns = append(turns, Turn{Role: c.Role, Content: c.Content})
	}
	// ByConversation returns newest-first; reverse for chronological order.
	for i, j := 0, len(turns)-1; i < j; i, j = i+1, j-1 {
		turns[i], turns[j] = turns[j], turns[i]
	}
	return turns, nil
}

// ConversationSummary describes one past conversation for listing.
type ConversationSummary struct {
	ID        string
	StartedAt time.Time
	TurnCount int
	Preview   string // first human message, truncated
}

// ListConversations returns summaries of recent conversations, newest first.
func (ms *MindStore) ListConversations(limit int) ([]ConversationSummary, error) {
	page, err := ms.store.ByType(EventTypeTurnCreated, 500, types.None[types.Cursor]())
	if err != nil {
		return nil, fmt.Errorf("query turns: %w", err)
	}

	// Group by conversation ID, track first human message and count.
	type convState struct {
		startedAt time.Time
		count     int
		preview   string
	}
	seen := make(map[string]*convState)
	var order []string // preserve insertion order

	for _, ev := range page.Items() {
		c, ok := ev.Content().(TurnCreatedContent)
		if !ok {
			continue
		}
		cid := ev.ConversationID().Value()
		if cid == "" {
			continue
		}
		state, exists := seen[cid]
		if !exists {
			state = &convState{startedAt: ev.Timestamp().Value()}
			seen[cid] = state
			order = append(order, cid)
		}
		state.count++
		if c.Role == "human" && state.preview == "" {
			preview := c.Content
			if len(preview) > 80 {
				preview = preview[:80] + "..."
			}
			state.preview = preview
		}
	}

	// Build summaries, newest first (reverse order).
	var summaries []ConversationSummary
	for i := len(order) - 1; i >= 0; i-- {
		cid := order[i]
		state := seen[cid]
		summaries = append(summaries, ConversationSummary{
			ID:        cid,
			StartedAt: state.startedAt,
			TurnCount: state.count,
			Preview:   state.preview,
		})
		if len(summaries) >= limit {
			break
		}
	}
	return summaries, nil
}

// FormatPriorDecisions formats a slice of observations for injection into CTO prompts.
// Returns an empty string if obs is empty.
func FormatPriorDecisions(obs []Observation) string {
	if len(obs) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\nPRIOR DECISIONS (from previous sessions — build on these, avoid repeating them):\n")
	for i, o := range obs {
		sb.WriteString(fmt.Sprintf("  %d. [%s] %s\n     Impact: %s\n", i+1, o.Mode, o.Proposed, o.Why))
	}
	return sb.String()
}
