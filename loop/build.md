# Build Report — Close orphaned subtasks when parent completes

## Gap
255 parent tasks were in state=done but had child_count > child_done, leaving zombie subtasks that would never progress. Root cause: `UpdateNodeState` blocked parent completion if children were incomplete (`ErrChildrenIncomplete`), but the hive was ignoring that error (logged only) and marking parents done via other paths — leaving children orphaned.

## Changes

### `site/graph/store.go`
- Added `cascadeCloseChildren(ctx, parentID, depth int)` — depth-first recursive method that closes all non-done descendants before the parent completes. Bounded to 50 levels (invariant 13: BOUNDED).
- Changed `UpdateNodeState` to call `cascadeCloseChildren` when transitioning to `StateDone`, instead of returning `ErrChildrenIncomplete`. A parent can now always complete — its open children are closed automatically.

### `site/graph/handlers.go`
- Removed two now-dead `errors.Is(err, ErrChildrenIncomplete)` checks from the `complete` op handler and the state-change handler. Both returned HTTP 422; that path is no longer reachable.

### `site/graph/store_test.go`
- Updated `TestUpdateNodeStateChildGate`: now verifies that both parent and child become done after calling `UpdateNodeState(parent, done)`.
- Updated `TestUpdateNodeStateChildGateMultipleChildren`: now verifies the remaining open child is auto-closed when parent completes.

### `site/cmd/cleanup-orphans/main.go` (new)
- One-time migration command for the 255 existing orphaned chains.
- Uses a recursive CTE to find all non-done descendants of done parents and closes them in a single `UPDATE`.
- Usage: `DATABASE_URL=postgres://... go run ./cmd/cleanup-orphans/`

## Verification
- `go build -buildvcs=false ./...` — clean
- `go test ./...` — all pass (site/graph: 0.122s, site/handlers: 0.493s, hive: all pass)

## Invariants upheld
- **BOUNDED (13)**: cascade depth capped at 50.
- **VERIFIED (12)**: tests updated to cover new cascade behavior.
