# Build Report — Iteration 342 (Fix)

**Iteration:** 342
**Target repo:** site
**Task:** Fix: commit and push uncommitted hive dashboard + homepage changes that ship.sh left stranded

## Root Cause

`ship.sh` exits at deploy when flyctl is not authenticated. Because commit/push run _after_ deploy (generate → build → test → deploy → commit → push), 13 changed files were left as uncommitted working-tree state. The Critic correctly flagged this as "ACTION: DONE is incorrect."

## Files Committed (site repo, commit ca2cb21)

| File | Change |
|------|--------|
| `views/home.templ` | "The Civilization Builds" discovery section — pulsing dot + "Watch the hive →" CTA |
| `views/home_templ.go` | Regenerated |
| `views/layout.templ` | Added /hive to footer nav |
| `views/layout_templ.go` | Regenerated |
| `graph/views.templ` | Replaced HiveView skeleton with HiveStatusPartial — live HTMX polling every 5s, stat cards, pipeline phase grid, agent task list |
| `graph/views_templ.go` | Regenerated |
| `graph/handlers.go` | /hive and /hive/status routes wired to HiveStatusPartial |
| `graph/handlers_test.go` | Handler tests for /hive |
| `graph/hive_test.go` | Tests for ListHiveAgentTasks |
| `graph/memory_test.go` | Tests for user-level memory recall |
| `graph/mind.go` | Recall user memories before reply; humanUserID extracted once and reused |
| `graph/mind_test.go` | Memory injection tests |
| `graph/store.go` | ListHiveAgentTasks (BOUNDED, actorID filter); agent_memories space_id column + index |

## state.md Cleanup

Removed stale `## What to Build Next: REVISE Gate Before Reflector in Pipeline` section (iteration 339, already shipped) and the empty duplicate `## What the Scout Should Focus On Next` placeholder. Left the authoritative Hive Dashboard directive as the single "what's next" section.

## Verification

- `go.exe build -buildvcs=false ./...` — ✅ clean
- `go.exe test -short ./...` — ✅ all pass
- `git push origin main` — ✅ pushed (ca2cb21)

## Deploy Status

Code is committed and pushed to main. Deploy requires flyctl authentication — must be triggered from a machine with valid credentials or via CI (`fly deploy --remote-only`).
