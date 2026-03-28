# Test Report: Iteration 405 — runObserverReason fallback cause fix

- **Timestamp:** 2026-03-29
- **Build commit:** 231ec28 — Fix: Builder skipped primary gap / deploy site fix + observer.go fallback cause unshipped

## What Was Tested

Builder fixed `pkg/runner/observer.go:runObserverReason` to apply a `fallbackCauseID` when the LLM emits `TASK_CAUSE: none`, closing task c2ab9f11 and satisfying Invariant 2 (CAUSALITY).

### Pre-existing tests (all still pass)

| Test | File | Verdict |
|------|------|---------|
| `TestRunObserverReason_FallbackCause` | `observer_test.go` | PASS |
| `TestParseObserverTasksCauseID` (6 subtests) | `observer_test.go` | PASS |
| `TestParseObserverTasksTwoCauseIDs` | `observer_test.go` | PASS |
| `TestCausality_LoopTaskCommandPath` | `causality_test.go` | PASS |
| `TestCausality_DirectAPICallPath` | `causality_test.go` | PASS |
| `TestCausality_LoopTaskCommandPath_MultipleTasks` | `causality_test.go` | PASS |
| `TestCausality_CmdPostPath` | `causality_test.go` | PASS |

### New tests added (edge cases)

| Test | What it guards | Verdict |
|------|----------------|---------|
| `TestRunObserverReason_FallbackCause_WhenFallbackEmpty` | No panic when `TASK_CAUSE:none` AND `fallbackCauseID=""` (empty graph). Task still created. | PASS |
| `TestRunObserverReason_OwnCauseTakesPrecedence` | When LLM supplies a valid `TASK_CAUSE` node ID, fallback must not overwrite it. | PASS |

## Results

```
Full suite:  all packages PASS
pkg/runner:  PASS  (includes new edge cases)
pkg/loop:    PASS  (causality_test.go unchanged)
```

## Coverage Notes

- The three CAUSALITY code paths (loop task command, direct API call, cmd/post) are fully covered by `causality_test.go`.
- The `runObserverReason` fallback is now exercised in three states: fallback applied, fallback empty, task has own cause.
- `parseObserverTasks` has table-driven coverage for all sentinel values (`none`, `N/A`, empty, whitespace).
- The `Operate()` path of `runObserver` has no unit test (shells out to Claude CLI) — acceptable; functional coverage comes from live runs.

## @Critic

Tests written and passing. Ready for critique.
