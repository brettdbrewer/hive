# Build Report — Agent Personas Infrastructure Phase 1

## Gap
Create `agent_personas` table, seed 8 starter personas at startup, and update `buildSystemPrompt()` in `mind.go` to use persona prompts when a `role:<name>` tag is present on the conversation.

## Status: ALREADY COMPLETE

All three components were implemented in a prior iteration. Task verified, no changes required.

### 1. `agent_personas` table — `site/graph/store.go:381`
```sql
CREATE TABLE IF NOT EXISTS agent_personas (
    id TEXT PRIMARY KEY, name TEXT UNIQUE, display TEXT,
    description TEXT, category TEXT, prompt TEXT,
    model TEXT, active BOOLEAN, created_at TIMESTAMPTZ
);
```
Supporting `agent_memories` table with `kind`, `source_id`, `importance` columns also present (lines 393–406).

### 2. Startup seeding — `site/cmd/site/main.go:215`
`graphStore.SeedAgentPersonas(context.Background())` is called at startup.
`SeedAgentPersonas` in `site/graph/personas.go` reads 50+ embedded `personas/*.md` files and upserts all of them with category/model/active metadata.

### 3. `buildSystemPrompt()` routing — `site/graph/mind.go:404`
Checks for `role:` tag → calls `s.GetAgentPersona(ctx, role)` → uses `persona.Prompt` instead of `mindSoul`. Also injects persona memories and saves new memories per conversation.

### Store methods verified present
- `UpsertAgentPersona`, `GetAgentPersona`, `GetAgentPersonaForConversation`
- `GetAgentPersonasForConversations` (batch, avoids N+1)
- `ListAgentPersonas` (active only, ordered by category/display)
- `RememberForPersona` / `RecallForPersona` (per-persona memory)

## Verification
- `go build -buildvcs=false ./...` — exit 0
- `go test ./...` — all pass (graph, auth)

## Files Changed
None — infrastructure was already complete.
