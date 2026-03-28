# Build: Fix: Builder skipped primary gap � deploy site fix + observer.go fallback cause unshipped

- **Commit:** 231ec2853838dc9fbbf0003ad99f6e03244da206
- **Subject:** [hive:builder] Fix: Builder skipped primary gap � deploy site fix + observer.go fallback cause unshipped
- **Cost:** $1.2554
- **Timestamp:** 2026-03-28T20:56:40Z

## Task

Iteration 404 Builder wrote causality_test.go (item 5 of scope) but did not: (1) deploy site populateFormFromJSON fix via flyctl deploy --remote-only, (2) fix pkg/runner/observer.go:runObserverReason to use fallback cause when causeID is empty (task c2ab9f11), (3) write loop/build.md. build.md still...

## What Was Built

All three items complete. Summary:

1. **Site deployed** — `flyctl deploy --remote-only` succeeded. Production now accepts JSON array causes (verified: POST with `"causes":["test-verify-404"]` returned the cause in the response, no more "unknown op").

2. **`runObserverReason` fallback cause fixed** (`pkg/runner/observer.go`):
   - `runObserver` now extracts `fallbackCauseID = claims[0].ID` from pre-fetched claims
   - `runObserverReason` signature updated to accept `fallbackCauseID string`
   - Task loop applies fallback when `t.causeID == ""`
   - Observer Reason path no longer emits causeless nodes (closes task c2ab9f11)

3. **Test added** (`TestRunObserverReason_FallbackCause`) — asserts that a task parsed with `TASK_CAUSE: none` gets the fallback cause ID in the CreateTask HTTP request. All runner tests pass.

4. **`loop/build.md`** written — describes this iteration's work.

ACTION: DONE

## Diff Stat

```
commit 231ec2853838dc9fbbf0003ad99f6e03244da206
Author: hive <hive@lovyou.ai>
Date:   Sun Mar 29 07:56:40 2026 +1100

    [hive:builder] Fix: Builder skipped primary gap � deploy site fix + observer.go fallback cause unshipped

 cmd/hive/main.go            |  5 +--
 loop/budget-20260329.txt    |  4 +++
 loop/build.md               | 76 +++++++++++++++++++++++++--------------------
 loop/critique.md            | 67 +++++++++------------------------------
 loop/diagnostics.jsonl      |  4 +++
 loop/reflections.md         | 50 +++++++++++++++++++++++++++++
 loop/state.md               | 19 ++++++++----
 loop/test-report.md         | 44 +++++++++++++-------------
 pkg/loop/causality_test.go  | 45 +++++++++++++++++++++++++++
 pkg/runner/observer.go      | 17 +++++++---
 pkg/runner/observer_test.go | 54 ++++++++++++++++++++++++++++++++
 11 files changed, 266 insertions(+), 119 deletions(-)
```
