# Build: Capture Operate summary in writeBuildArtifact

- **Commit:** (pending — not yet committed)
- **Subject:** [hive:builder] Capture Operate summary in writeBuildArtifact
- **Cost:** $0.0000
- **Timestamp:** 2026-03-27T00:00:00Z

## Task

In `pkg/runner/runner.go`, change `writeBuildArtifact(t api.Node, costUSD float64)` to accept a third `operateSummary string` parameter. Add a `## What Was Built` section to build.md that includes the summary (truncated to 2000 chars) between the metadata block and the diff stat. In `workTask`, pass `result.Summary` as the third arg at line 337. Add a test in `pkg/runner/runner_test.go` that verifies the summary content appears in build.md output.

## What Was Built

- Changed `writeBuildArtifact` signature to accept `operateSummary string` as third parameter
- Summary is truncated to 2000 chars if longer
- Added `## What Was Built` section written after `## Task` and before `## Diff Stat`
- Updated call site in `workTask` to pass `result.Summary`
- Added `TestBuildArtifactContainsSummary`: verifies the summary text appears under `## What Was Built`
- Added `TestBuildArtifactSummaryTruncated`: verifies summaries >2000 chars are truncated
- Updated `TestBuildArtifactWritten` to pass empty string for backward compatibility
- Build: `go.exe build -buildvcs=false ./...` — pass
- Tests: `go.exe test ./pkg/runner/...` — pass (3.587s)

## Diff Stat

```
pkg/runner/runner.go      | 8 ++++++--
pkg/runner/runner_test.go | 60 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
2 files changed, 66 insertions(+), 2 deletions(-)
```
