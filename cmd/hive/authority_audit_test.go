package main

import (
	"strings"
	"testing"

	"github.com/transpara-ai/eventgraph/go/pkg/event"
	"github.com/transpara-ai/eventgraph/go/pkg/store"
	"github.com/transpara-ai/eventgraph/go/pkg/types"

	"github.com/transpara-ai/hive/pkg/safety"
)

func TestAuthorityAuditEmitterRecordsAuthorityRequested(t *testing.T) {
	s := store.NewInMemoryStore()
	actorID := types.MustActorID("actor_test_authority_audit")
	if err := bootstrapGraph(s, actorID); err != nil {
		t.Fatalf("bootstrap graph: %v", err)
	}

	emitter := newAuthorityAuditEmitterForStore(s, actorID)
	if err := emitter.EmitAuthorityRequest(
		safety.ActionRepoCreate,
		safety.ApprovalRequired,
		"create repo test-service",
	); err != nil {
		t.Fatalf("emit authority request: %v", err)
	}

	page, err := s.ByType(event.EventTypeAuthorityRequested, 10, types.None[types.Cursor]())
	if err != nil {
		t.Fatalf("query authority requests: %v", err)
	}
	if len(page.Items()) != 1 {
		t.Fatalf("authority.requested count = %d, want 1", len(page.Items()))
	}
	content, ok := page.Items()[0].Content().(event.AuthorityRequestContent)
	if !ok {
		t.Fatalf("content type = %T, want event.AuthorityRequestContent", page.Items()[0].Content())
	}
	if content.Action != string(safety.ActionRepoCreate) {
		t.Fatalf("action = %q, want %q", content.Action, safety.ActionRepoCreate)
	}
	if content.Level != event.AuthorityLevelRequired {
		t.Fatalf("level = %q, want %q", content.Level, event.AuthorityLevelRequired)
	}
	if content.Actor != actorID {
		t.Fatalf("actor = %s, want %s", content.Actor, actorID)
	}
	if !strings.Contains(content.Justification, "test-service") {
		t.Fatalf("justification missing repo context: %q", content.Justification)
	}
	if content.Causes.Len() != 1 {
		t.Fatalf("causes len = %d, want 1", content.Causes.Len())
	}
}

func TestAuthorizeFinalPipelineSweepEmitsAuthorityRequest(t *testing.T) {
	s := store.NewInMemoryStore()
	actorID := types.MustActorID("actor_test_pipeline_authority")
	if err := bootstrapGraph(s, actorID); err != nil {
		t.Fatalf("bootstrap graph: %v", err)
	}
	emitter := newAuthorityAuditEmitterForStore(s, actorID)

	err := authorizeFinalPipelineSweep(map[string]string{
		"hive": "/tmp/hive",
		"site": "/tmp/site",
	}, "/tmp/hive", emitter)
	if err == nil {
		t.Fatal("expected cross-repo mutation authority error")
	}

	page, err := s.ByType(event.EventTypeAuthorityRequested, 10, types.None[types.Cursor]())
	if err != nil {
		t.Fatalf("query authority requests: %v", err)
	}
	if len(page.Items()) != 1 {
		t.Fatalf("authority.requested count = %d, want 1", len(page.Items()))
	}
	content := page.Items()[0].Content().(event.AuthorityRequestContent)
	if content.Action != string(safety.ActionRepoMutateCrossRepo) {
		t.Fatalf("action = %q, want %q", content.Action, safety.ActionRepoMutateCrossRepo)
	}
	if !strings.Contains(content.Justification, "repos=2") {
		t.Fatalf("justification missing repo count: %q", content.Justification)
	}
}

func TestAuthorizeIngestRepoBootstrapEmitsAuthorityRequests(t *testing.T) {
	s := store.NewInMemoryStore()
	actorID := types.MustActorID("actor_test_ingest_authority")
	if err := bootstrapGraph(s, actorID); err != nil {
		t.Fatalf("bootstrap graph: %v", err)
	}
	emitter := newAuthorityAuditEmitterForStore(s, actorID)

	err := authorizeIngestRepoBootstrap("new-service", "transpara-ai", "/tmp/new-service", emitter)
	if err == nil {
		t.Fatal("expected repo bootstrap authority error")
	}

	page, err := s.ByType(event.EventTypeAuthorityRequested, 10, types.None[types.Cursor]())
	if err != nil {
		t.Fatalf("query authority requests: %v", err)
	}
	if len(page.Items()) != 2 {
		t.Fatalf("authority.requested count = %d, want 2", len(page.Items()))
	}
	got := map[string]bool{}
	for _, ev := range page.Items() {
		content := ev.Content().(event.AuthorityRequestContent)
		got[content.Action] = true
	}
	for _, want := range []safety.ProtectedAction{
		safety.ActionRepoCreate,
		safety.ActionRepoPushDefaultBranch,
	} {
		if !got[string(want)] {
			t.Fatalf("missing authority request for %s", want)
		}
	}
}
