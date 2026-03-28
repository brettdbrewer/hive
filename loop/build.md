# Build: Fix: Validate LLM-generated cause IDs in Observer before posting

## Task

Validate LLM-generated cause IDs in Observer before posting nodes. After parsing LLM response and before posting any node, check that each cause ID in the response actually exists on the graph. Call `NodeExists` for each; if not found, log a warning and replace with the fallback cause. Prevents Observer from silently posting nodes with dangling/hallucinated cause IDs (Lesson 170).

## What Was Built

### 1. `pkg/api/client.go` — Added `NodeExists(slug, id string) bool`

New method on `Client` that validates whether a node ID exists on the graph before it is used as a cause. Issues a `GET /app/{slug}/node/{id}?format=json` and returns `true` only on HTTP 200. Network errors and 404s return `false`. Drains the response body for connection reuse.

### 2. `pkg/runner/observer.go` — Validate LLM cause IDs in `runObserverReason`

When `t.causeID` is non-empty (LLM provided an ID), `NodeExists` is called before using it. If the node does not exist on the graph, a warning is logged and `fallbackCauseID` is used instead. This closes Lesson 170: ghost IDs from LLM hallucination no longer propagate silently.

```
if causeID == "" {
    causeID = fallbackCauseID          // was already: fallback for "none"
} else if r.cfg.APIClient != nil {
    if !r.cfg.APIClient.NodeExists(...) {
        log.Printf("[observer] warning: LLM cause ID %q not found on graph; using fallback", causeID)
        causeID = fallbackCauseID       // new: fallback for hallucinated IDs
    }
}
```

### 3. `pkg/runner/observer_test.go` — `TestRunObserverReason_HallucinatedCauseIDGetsReplaced`

New test. Server returns 404 for `GET .../node/ghost-node-does-not-exist`. Asserts that the `CreateTask` request uses `fallbackCauseID`, not the ghost ID. Completes the test matrix for `runObserverReason` cause-validation paths.

## Verification

- `go.exe build -buildvcs=false ./pkg/api/...` — OK
- `go.exe build -buildvcs=false ./pkg/runner/...` — OK
- `go.exe test -buildvcs=false ./pkg/runner/... ./pkg/api/... ./pkg/loop/...` — all pass

```
ok  github.com/lovyou-ai/hive/pkg/runner  6.303s
ok  github.com/lovyou-ai/hive/pkg/api     1.004s
ok  github.com/lovyou-ai/hive/pkg/loop    1.156s
```

ACTION: DONE
