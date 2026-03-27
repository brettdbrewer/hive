# Test Report: False completion epidemic — child gate enforcement

**Date:** 2026-03-28

## What Was Tested

Builder added `ErrChildrenIncomplete` enforcement to `UpdateNodeState` and 422 responses in both completion handler paths (`handleOp` case `"complete"` and `handleNodeState`). Builder also wrote `TestUpdateNodeStateChildGate` covering the basic single-child case.

Tester added 5 targeted tests covering gaps in the Builder's coverage.

## Tests Run

### Pre-existing (Builder wrote)
**`TestUpdateNodeStateChildGate`** — `site/graph/store_test.go`
- Parent with one incomplete child → `ErrChildrenIncomplete`
- Complete child → parent now completable
- Result: **PASS** ✓

### Added by Tester

**`TestUpdateNodeStateChildGateLeafNode`** — `site/graph/store_test.go`
- Leaf node (zero children) → completes without error
- Critical: the fix must not block leaf nodes. `COUNT(*) = 0 → incomplete = 0` path confirmed.
- Result: **PASS** ✓

**`TestUpdateNodeStateChildGateMultipleChildren`** — `site/graph/store_test.go`
- 2 children, only 1 done → still blocked by `ErrChildrenIncomplete`
- Both done → parent completes
- Result: **PASS** ✓

**`TestUpdateNodeStateNonDoneSkipsGate`** — `site/graph/store_test.go`
- Parent with incomplete child → setting to `StateReview` (non-done) does NOT trigger gate
- Confirms the `if state == StateDone` guard is correctly scoped
- Result: **PASS** ✓

**`TestHandlerCompleteOpChildrenIncomplete`** — `site/graph/handlers_test.go`
- `POST /app/{slug}/op` with `{"op":"complete","node_id":"..."}` when parent has incomplete child
- Verifies handler returns **422 Unprocessable Entity** (not 500)
- Result: **PASS** ✓

**`TestHandlerNodeStateChildrenIncomplete`** — `site/graph/handlers_test.go`
- `POST /app/{slug}/node/{id}/state` with `state=done` when parent has incomplete child
- Verifies `handleNodeState` returns **422 Unprocessable Entity** (not 500)
- Result: **PASS** ✓

## Coverage Summary

| Path | Tested? |
|------|---------|
| `UpdateNodeState` single incomplete child → rejected | ✓ (Builder) |
| `UpdateNodeState` child completed → parent accepted | ✓ (Builder) |
| `UpdateNodeState` leaf node (no children) → accepted | ✓ (Tester) |
| `UpdateNodeState` multiple children, partial → rejected | ✓ (Tester) |
| `UpdateNodeState` multiple children, all done → accepted | ✓ (Tester) |
| `UpdateNodeState` non-done state skips gate | ✓ (Tester) |
| `handleOp` "complete" → 422 on `ErrChildrenIncomplete` | ✓ (Tester) |
| `handleNodeState` state=done → 422 on `ErrChildrenIncomplete` | ✓ (Tester) |

## Pre-existing Failures (not regressions)

`go test ./graph/` shows 9 pre-existing failures unrelated to this feature:
- Duplicate key violations from stale test DB data (slug conflicts)
- `TestReportsAndResolve`: SQL scan error on `Op` type (pre-existing schema mismatch)
- `TestReposts`: nil pointer dereference (pre-existing)
- `TestHivePage`: content string mismatch (pre-existing)

All gate-related tests isolated and confirmed passing. No regressions introduced.

## Verdict

**PASS.** The child gate enforces correctly across all tested scenarios. Both HTTP completion paths return 422 as specified. No regressions on new tests.
