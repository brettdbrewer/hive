# Build: Test PipelineTree failure path in pkg/runner/pipeline_tree_test.go

- **Files created:**
  - `pkg/runner/pipeline_tree.go` — `Phase` and `PipelineTree` types; `NewPipelineTree(r *Runner)` wires the four pipeline phases (scout, architect, builder, critic); `Execute(ctx)` runs phases in order and emits a `PhaseEvent{Outcome: "failure"}` diagnostic via `appendDiagnostic` on the first failure
  - `pkg/runner/pipeline_tree_test.go` — `TestPipelineTreeFailureWritesDiagnostic`: constructs a `PipelineTree` directly with a stub phase returning `fmt.Errorf("injected failure")`, calls `Execute(ctx)`, and verifies `loop/diagnostics.jsonl` contains a line with `outcome=failure` and `phase=stub`
- **Build:** `go.exe build -buildvcs=false ./...` — clean
- **Tests:** `go.exe test ./...` — all pass (pkg/runner: 1.154s)
- **Timestamp:** 2026-03-27
