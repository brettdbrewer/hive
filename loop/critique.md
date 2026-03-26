# Critique

Commit: b652025717aec3db53007e67e2e280719f11a5f9
Verdict: PASS

I have everything I need. Here is my analysis:

---

## Critic Analysis

### What was built

`pipeline_tree.go` and `pipeline_tree_test.go` were created. The `PipelineTree` type is clean: it holds a slice of `Phase` structs, `Execute` runs them in order, and on the first failure it writes a `PhaseEvent` to `diagnostics.jsonl` via `appendDiagnostic` and short-circuits. The test constructs a `PipelineTree` directly with a stub failure phase, calls `Execute`, reads `diagnostics.jsonl`, and asserts `outcome=failure` and `phase=stub`. Test coverage is real and correct.

### Issue 1 — The "pipeline" case was NOT wired into `runTick` (correctness failure)

The scout scoped this explicitly:
> **Files to modify:** `pkg/runner/runner.go` — add "pipeline" case to runTick (lines ~176), call `PipelineTree.Execute(ctx)`

The build.md also listed it. But `runner.go` was not changed. `runTick`'s switch has no `"pipeline"` case. The `PipelineTree` type exists but is unreachable from the live pipeline. If `--role pipeline` is passed today, it falls through to the `default:` log line and does nothing.

This is the **primary deliverable of the iteration** — the entire point was to make the pipeline role use the decision tree for phase orchestration. The test exercises the `PipelineTree` struct directly, not through `runTick`, so it passes while the integration is broken.

### Issue 2 — `NewPipelineTree` phases always return `nil` (correctness concern, documented but worth flagging)

```go
{Name: "scout", Run: func(ctx context.Context) error { r.runScout(ctx); return nil }},
```

All four wrappers swallow any internal failures and return `nil`. The comment in the code acknowledges this ("Phase 2 work once the phase methods propagate errors up"), which is honest. But it means the failure-detection path in `Execute` is dead code for production phases — it can only be exercised via the test's stub. This is acceptable as scaffolding, but the runTick wiring gap means even the scaffolding isn't live.

### Issue 3 — Invariant 12 (VERIFIED): the test does not cover the `NewPipelineTree` constructor

The test bypasses `NewPipelineTree` entirely and constructs `PipelineTree` directly. `NewPipelineTree` is the real integration surface (it wires the four runner methods), but it has zero test coverage. A typo in one of those closures would be invisible.

### Checklist

| Check | Result |
|---|---|
| Identity (inv 11) | N/A — no JOINs or name-matching |
| Bounded (inv 13) | Loop over `pt.phases` is bounded by slice length |
| SQL injection / race conditions | None — no DB, no goroutines |
| Tests exist | Yes, but they don't exercise the wiring path (`NewPipelineTree`) |
| `runTick` wired | **NO** — missing `"pipeline"` case |

---

**VERDICT: REVISE**

**Required fix:**

Add the `"pipeline"` case to `runTick` in `pkg/runner/runner.go`:

```go
case "pipeline":
    if err := NewPipelineTree(r).Execute(ctx); err != nil {
        log.Printf("[pipeline] tick %d: %v", r.tick, err)
    }
```

Without this, the iteration's stated goal — "Wire runTick to use the tree for `role == pipeline`" — is undelivered, and `PipelineTree` is unreachable in production.
