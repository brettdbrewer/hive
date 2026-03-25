# Build Report — Fix cmd/mcp-graph: security hardening + VERIFIED invariant

## Gap
Critic identified 5 issues in `cmd/mcp-graph/main.go` from commit 3fe5dc3:
1. URL injection in `toolSearch` — raw query param concatenation
2. Path injection in `toolGetNode` — raw nodeID in URL path
3. No HTTP client timeout — violates BOUNDED (invariant 13)
4. Unbounded `io.ReadAll` — no response body size cap
5. No tests — violates VERIFIED (invariant 12)

## Changes

### `cmd/mcp-graph/main.go`
- Added `net/url` and `time` imports
- Added `maxResponseBytes = 1 MiB` constant (BOUNDED)
- Added `client *http.Client` field to `server` struct; `newServer()` sets `Timeout: 30s` (BOUNDED)
- All `apiPost` calls use `url.PathEscape(space)` for path segment
- `toolSearch`: query param uses `url.QueryEscape(query)` — fixes URL injection
- `toolGetNode`: validates nodeID against `/?#` chars; uses `url.PathEscape` — fixes path injection
- `toolGetBoard`: uses `url.PathEscape(space)`
- `apiGet` / `apiPost`: replaced `http.DefaultClient` with `s.client`; replaced bare `io.ReadAll` with `io.ReadAll(io.LimitReader(..., maxResponseBytes))`

### `cmd/mcp-graph/main_test.go` (new — 18 tests)
- `spaceFor` default/override
- `toolIntend` missing-title validation, correct POST body
- `toolRespond` missing-field validation, correct POST body
- `toolSearch` missing-query validation, URL encoding of special chars
- `toolGetBoard` happy path
- `toolGetNode` missing-ID, path-traversal rejection, correct path encoding
- `apiGet` auth header, 4xx error, bounded read (2 MiB server → ≤1 MiB read)
- `apiPost` content-type header
- `okResult` / `errResult` helpers

## Build Results
```
go.exe build -buildvcs=false ./...         ✓ clean
go.exe test ./cmd/mcp-graph/... -v         ✓ 18/18 pass
```

---
# Previous: Agent Discovery Page + Chat Creation Flow

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
