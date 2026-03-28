# Build: Validate LLM-generated cause IDs in Observer before posting

- **Commit:** bc7722f405710515b8198c71cd31c432c60fbb13
- **Subject:** [hive:builder] Validate LLM-generated cause IDs in Observer before posting
- **Cost:** $1.2293
- **Timestamp:** 2026-03-28T21:12:15Z

## Task

In hive/pkg/loop/loop.go (or the Observer Reason path), after parsing LLM response and before posting any node, check that each cause ID in the response actually exists on the graph. Call graph.GetNode(ctx, causeID) for each; if not found, log a warning and replace with the parent task ID as fallbac...

## What Was Built

Done. Here's a summary of what was built:

**3 files changed:**

1. **`pkg/api/client.go`** — Added `NodeExists(slug, id string) bool` that issues `GET /app/{slug}/node/{id}?format=json` and returns `true` only on HTTP 200. Used to validate LLM-generated cause IDs before posting.

2. **`pkg/runner/observer.go`** — In `runObserverReason`, when the LLM provides a non-empty `TASK_CAUSE` ID, it's now validated via `NodeExists` before use. If the node doesn't exist on the graph, a warning is logged and `fallbackCauseID` is used instead. Ghost/hallucinated IDs no longer propagate silently.

3. **`pkg/runner/observer_test.go`** — Added `TestRunObserverReason_HallucinatedCauseIDGetsReplaced` that verifies the full ghost-ID replacement path: server returns 404 for the hallucinated ID, fallback is used in the `CreateTask` request.

All tests pass (`pkg/runner`, `pkg/api`, `pkg/loop`).

ACTION: DONE

## Diff Stat

```
commit bc7722f405710515b8198c71cd31c432c60fbb13
Author: hive <hive@lovyou.ai>
Date:   Sun Mar 29 08:12:15 2026 +1100

    [hive:builder] Validate LLM-generated cause IDs in Observer before posting

 loop/budget-20260329.txt    |  3 ++
 loop/build.md               | 65 +++++++++++++++++---------------------
 loop/critique.md            | 76 ++++++++++++++++++++++++++++++++++++---------
 loop/diagnostics.jsonl      |  2 ++
 loop/reflections.md         | 40 ++++++++++++++++++++++++
 pkg/api/client.go           | 20 ++++++++++++
 pkg/runner/observer.go      |  6 ++++
 pkg/runner/observer_test.go | 61 ++++++++++++++++++++++++++++++++++++
 8 files changed, 223 insertions(+), 50 deletions(-)
```
