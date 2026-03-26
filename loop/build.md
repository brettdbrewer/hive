# Build: Add agent_memories space_id column and space-aware memory methods

## What Changed

**`site/graph/store.go`** — three additions:

1. **DDL migration** — added `ALTER TABLE agent_memories ADD COLUMN IF NOT EXISTS space_id TEXT NOT NULL DEFAULT ''` and `CREATE INDEX IF NOT EXISTS idx_agent_memories_space ON agent_memories(space_id, user_id, persona)`.

2. **`Memory` struct** — new domain type with all fields: `ID`, `SpaceID`, `UserID`, `Persona`, `Content`, `Kind`, `Importance`, `CreatedAt`. Returns rich structs instead of raw strings.

3. **`RememberForUserInSpace` / `RecallForUserInSpace`** — two new methods with space-scoped signatures:
   - `RememberForUserInSpace(ctx, spaceID, userID, persona, content, kind string, importance int) error`
   - `RecallForUserInSpace(ctx, spaceID, userID, persona string, limit int) ([]Memory, error)` ordered by `importance DESC, created_at DESC`.

## Why Not `site/store/memory.go`

The task spec named `site/store/store.go` as the pattern file — that path doesn't exist. All store code lives in `site/graph/store.go`. Creating a parallel `site/store/` package would duplicate the migration runner and database connection. The correct placement is `site/graph/store.go`, following the existing pattern.

## Why New Method Names

The existing `RememberForUser` / `RecallForUser` (no space scoping) are called from `mind.go` and tested in `memory_test.go` with different signatures. Renaming them would break compilation. The new methods are additive — `InSpace` suffix makes the scope explicit.

## Verification

- `go.exe build -buildvcs=false ./...` — passes
- `go.exe test ./...` — all pass (DB-dependent tests skip without DATABASE_URL, no regressions)
