---

# SCOUT REPORT

## Gap: Autonomous Loop Cannot Close — Reflector Not in Pipeline

**Gap:** The PipelineTree (scout → architect → builder → critic) completes but never closes. Reflector exists and works, but runs only on the daemon ticker, not as part of the automated pipeline. Without closure, the iteration counter never advances, state.md's directive section stays stale, and PM can't generate new directions.

**Evidence:**
- `pkg/runner/pipeline_tree.go` has 4 phases: scout, architect, builder, critic (lines 52-57)
- `pkg/runner/reflector.go` is complete, reads scout/build/critique artifacts, appends reflections, advances iteration counter (lines 100-156)
- `pkg/runner/critic.go` already writes `loop/critique.md` (line 119) ✓
- Reflector is called only from runner's tick loop (`runner.go:184`), not from PipelineTree
- When pipeline runs (`--pipeline` mode), it halts after Critic — Reflector never executes
- Recent commits (iter 92ad958 onwards) wired failure detection but skipped the last phase

**Impact:** The autonomous loop stalls at closure. Scout-Architect-Builder-Critic work autonomously, but can't declare done and move forward. PM has no fresh directive to work from. Each pipeline run repeats the same tasks if Scout's directive is unchanged. The system proves it can build autonomously but not learn and re-prioritize autonomously.

**Scope:** One file: `pkg/runner/pipeline_tree.go`. Lines 52-59 (the phases list in `NewPipelineTree`). The Reflector's tick gate (`r.tick%4 != 0`) needs to be bypassed when called from pipeline mode.

**Suggestion:** Add Reflector as phase 5 in `NewPipelineTree`:
```go
{Name: "reflector", Run: func(ctx context.Context) error {
    // Bypass tick gate in pipeline mode (one-shot = pipeline)
    saved := r.cfg.OneShot
    r.cfg.OneShot = true
    r.runReflector(ctx)
    r.cfg.OneShot = saved
    return nil
}},
```

Then update pipeline_tree_test.go to verify the tree has 5 phases. Done criteria: `go test ./pkg/runner/...` passes, `NewPipelineTree(r)` includes "reflector" phase, pipeline run produces an advanced iteration counter in state.md.

---

**Ready for the Architect to create the implementation task.**