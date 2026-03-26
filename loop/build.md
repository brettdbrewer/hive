# Build: Wire pipeline role into runTick in runner.go

- **File changed:** `pkg/runner/runner.go` — added `case "pipeline": _ = NewPipelineTree(r).Execute(ctx)` between the `"architect"` and `"observer"` cases in `runTick`
- **Build:** `go.exe build -buildvcs=false ./...` — clean
- **Tests:** `go.exe test ./pkg/runner/...` — ok (1.132s)
- **Timestamp:** 2026-03-27
