# Build: Fix: [hive:builder] Auth: helpful error messages and logging

## Gap Addressed

CAUSALITY GATE 1 (Lesson 167, Scout 406): empty `causeIDs` could reach the graph unvalidated via `assertScoutGap` and `assertCritique`. The fix adds `assertClaim` as a typed boundary that enforces non-empty causes before any HTTP call.

## What Was Built

### `hive/cmd/post/main.go`

**New function: `assertClaim`**

```go
func assertClaim(apiKey, baseURL string, causeIDs []string, kind, title, body string) (string, error)
```

- Enforces `len(causeIDs) > 0` — returns error with "Invariant 2: CAUSALITY" message if empty
- Posts `op=assert` with causes always set (no `if len(causeIDs) > 0` conditional)
- Returns `(nodeID string, error)` — callers can use the created node ID as a cause
- Guard fires **before** any HTTP I/O — zero network cost for invariant violations

**Refactored: `assertScoutGap`**

Removed inline HTTP posting logic. Now calls:
```go
assertClaim(apiKey, baseURL, causeIDs, "claim", gapTitle, body)
```
Returns the CAUSALITY error if `causeIDs` is empty.

**Refactored: `assertCritique`**

Same refactor — removed inline HTTP posting. Now calls:
```go
assertClaim(apiKey, baseURL, causeIDs, "claim", title, string(data))
```
Returns the CAUSALITY error if `causeIDs` is empty.

### `hive/cmd/post/main_test.go`

**New test: `TestAssertClaim_RejectsEmptyCauseIDs`**

Two subtests:
- `nil` — `assertClaim("key", srv.URL, nil, ...)` → error containing "CAUSALITY", HTTP server not called
- `empty_slice` — `assertClaim("key", srv.URL, []string{}, ...)` → error, HTTP server not called

**Updated tests (3):**

The existing tests `TestAssertScoutGapCreatesClaimNode`, `TestAssertScoutGapSendsAuthHeader`, and `TestAssertCritiqueCreatesClaimNode` called their functions with `nil` causeIDs (testing the happy path before CAUSALITY gate existed). Updated to pass `[]string{"cause-node-abc"}` / `[]string{"cause-id"}` / `[]string{"task-node-xyz"}`. `TestAssertScoutGapCreatesClaimNode` also now asserts `received["causes"] == "doc-node-abc"`.

## Build Results

```
go.exe build -buildvcs=false ./...   → OK
go.exe test -buildvcs=false ./...    → all pass (15 packages)
```

## Files Changed

- `hive/cmd/post/main.go` — new `assertClaim`, `assertScoutGap` + `assertCritique` refactored to use it
- `hive/cmd/post/main_test.go` — new `TestAssertClaim_RejectsEmptyCauseIDs`, 3 existing tests updated
