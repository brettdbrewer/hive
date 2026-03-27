package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// EventListener receives op events from the site's webhook and dispatches
// to the right agent. This is the pub/sub backbone — events arrive, agents react.
type EventListener struct {
	runner  *Runner
	ctx     context.Context
	port    string
}

// OpEvent is the JSON payload from the site's webhook.
type OpEvent struct {
	ID        string          `json:"id"`
	SpaceID   string          `json:"space_id"`
	NodeID    string          `json:"node_id"`
	NodeTitle string          `json:"node_title"`
	Actor     string          `json:"actor"`
	ActorID   string          `json:"actor_id"`
	ActorKind string          `json:"actor_kind"`
	Op        string          `json:"op"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

// StartEventListener starts an HTTP server that receives webhook events
// from the site and dispatches to agents.
func StartEventListener(ctx context.Context, r *Runner, port string) error {
	el := &EventListener{runner: r, ctx: ctx, port: port}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /event", el.handleEvent)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	log.Printf("[pubsub] listening on :%s", port)
	return server.ListenAndServe()
}

func (el *EventListener) handleEvent(w http.ResponseWriter, r *http.Request) {
	var event OpEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	// Don't react to our own events (avoid infinite loops).
	if event.ActorKind == "agent" {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Printf("[pubsub] event: op=%s actor=%s node=%s", event.Op, event.Actor, event.NodeTitle)

	// Dispatch based on op type and context.
	go el.dispatch(event)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"received"}`)
}

func (el *EventListener) dispatch(event OpEvent) {
	switch event.Op {
	case "intend":
		// Task created — if assigned, notify the assigned agent.
		log.Printf("[pubsub] task created: %s", event.NodeTitle)

	case "assign":
		// Task assigned to an agent — that agent should pick it up.
		log.Printf("[pubsub] task assigned: %s", event.NodeTitle)
		// The assigned agent runs as builder to work the task.
		// Future: invoke the specific agent by ID.

	case "complete":
		// Task completed — Reflector should note it, PM should check board.
		log.Printf("[pubsub] task completed: %s", event.NodeTitle)

	case "respond":
		// Message in conversation — agent should reply if it's a participant.
		log.Printf("[pubsub] message: %s", event.NodeTitle)

	case "assert":
		// Claim asserted — relevant agents may want to verify or challenge.
		log.Printf("[pubsub] claim: %s", event.NodeTitle)

	case "express":
		// Post created — social notification.
		if strings.HasPrefix(event.NodeTitle, "Scout Report:") {
			log.Printf("[pubsub] scout report published — Architect should decompose")
		}
	}
}
