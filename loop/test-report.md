# Test Report — Close orphaned subtasks when parent completes

## What was tested
The cascade-close logic added in iteration 391: `cascadeCloseChildren` and the updated `UpdateNodeState` behavior in `site/graph/store.go`.

## New test added

### `TestCascadeDepthBoundary` (`site/graph/store_test.go`)
- **Gap addressed:** State.md lesson 179 called for a boundary test; previous test-report noted depth limit as "not tested".
- **What it tests:** `maxCascadeDepth = 50` is actually enforced. Calls `cascadeCloseChildren` directly at `depth=50` on a parent with one child — the recursive call hits `depth=51`, exceeds the cap, and must return an error containing "cascade depth exceeded".
- **Why not 51 real levels:** Since tests are in `package graph` (same package), we can call the unexported `cascadeCloseChildren` directly with a pre-set depth. Only 2 DB nodes needed.
- **Result:** PASS

## All cascade tests

| Test | Status |
|------|--------|
| `TestUpdateNodeStateChildGate` | PASS |
| `TestUpdateNodeStateChildGateLeafNode` | PASS |
| `TestUpdateNodeStateChildGateMultipleChildren` | PASS |
| `TestCascadeCloseChildrenDeep` | PASS |
| `TestUpdateNodeStateNonDoneSkipsGate` | PASS |
| `TestCascadeDepthBoundary` (new) | PASS |

## Pre-existing failure (not our code)

`TestReposts` panics on `sp.ID` (nil pointer) — `CreateSpace` fails on slug collision ("repost-test" already exists from a prior run) and the test ignores the error. Confirmed pre-existing on `main` before iteration 391's changes.

## Dead code flag

`ErrChildrenIncomplete` is declared in `site/graph/store.go:211` but is never returned anywhere — `cascadeCloseChildren` replaced the rejection path, the two handler callers were removed per build.md. State.md item 8 calls for an audit of external callers; no callers found in `site/`. Flagging for the Critic.

## Coverage notes

- `cascadeCloseChildren` happy path (2-level): TestUpdateNodeStateChildGate
- `cascadeCloseChildren` happy path (3-level recursive): TestCascadeCloseChildrenDeep
- `cascadeCloseChildren` multiple siblings: TestUpdateNodeStateChildGateMultipleChildren
- `cascadeCloseChildren` depth bound (51 > 50): **TestCascadeDepthBoundary** (new)
- `UpdateNodeState` non-done bypass: TestUpdateNodeStateNonDoneSkipsGate
- `UpdateNodeState` on nonexistent ID (ErrNotFound): not covered — pre-existing gap

@Critic — ready for review.
