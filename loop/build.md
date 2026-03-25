# Build Report — Agent Discovery Page + Chat Creation Flow

## What Was Built

Phase 2 of agent-chat-spec: the `/agents` discovery page and `POST /agents/{name}/chat` conversation creation flow.

## Findings

The core implementation was already complete from Phase 1 infrastructure work:
- `GET /agents` route — lists active personas grouped by category via `views.AgentsPage`
- `GET /agents/{name}` route — individual agent profile page
- `POST /agents/{name}/chat` route — creates a conversation with `role:{name}` tag
- `views.AgentCategoryGroup` + `views.AgentPersonaItem` types in `agents_templ.go`
- Persona cards with name, description, category badge, Profile + Chat buttons

**The one gap:** the POST handler was routing conversations into the **demo space** (`graph.DemoSpaceSlug`) instead of the dedicated **agents space** (`graph.AgentsSpaceSlug`). The spec explicitly calls for `lovyou.ai/app/agents` as the home for all agent conversations.

## Changes Made

### `site/cmd/site/main.go`
- Updated `POST /agents/{name}/chat` handler to use `graph.AgentsSpaceSlug` instead of `graph.DemoSpaceSlug`
- Updated error message from "demo space not available" to "agents space not available"
- Renamed local variable `demoSpace` → `agentsSpace` for clarity
- The `agents` space is guaranteed to exist: `graphStore.EnsureAgentsSpace()` runs at startup (line 214)

## Verification

```
go.exe build -buildvcs=false ./...   → exit 0
go.exe test ./...                    → all pass (auth, graph packages)
```

## Route Summary

| Route | Auth | Effect |
|-------|------|--------|
| `GET /agents` | optional | Shows all active personas grouped by category |
| `GET /agents/{name}` | optional | Shows agent profile with full prompt rendered as HTML |
| `POST /agents/{name}/chat` | required | Creates conversation in `/app/agents` space with `role:{name}` tag, redirects to chat |
