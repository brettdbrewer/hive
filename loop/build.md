# Build: Wire Reflector as phase 5 in PipelineTree

**Gap:** Reflector was not wired as phase 5 in PipelineTree, so the loop never appended reflections.md or advanced state.md after Critic PASS.

**Change:** Added `"reflector"` as the fifth phase in `NewPipelineTree` (`pkg/runner/pipeline_tree.go`).

The phase temporarily sets `r.cfg.OneShot = true` to bypass the `tick%4` throttle gate in `runReflector`, calls `r.runReflector(ctx)`, then restores the original value. The `r.done = true` side-effect from one-shot mode is harmless — `PipelineTree.Execute` does not check `r.done`.

**Test added:** `TestNewPipelineTreeHasFivePhases` in `pkg/runner/pipeline_tree_test.go` — verifies all five phases (scout, architect, builder, critic, reflector) are present in order.

**Files changed:**
- `pkg/runner/pipeline_tree.go` — added reflector phase
- `pkg/runner/pipeline_tree_test.go` — added `TestNewPipelineTreeHasFivePhases`

**Build:** `go.exe build -buildvcs=false ./...` — clean
**Tests:** `go.exe test -buildvcs=false ./...` — all pass
