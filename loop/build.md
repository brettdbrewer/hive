# Build: Scout phase must assert gap as claim node — currently writes only flat file

## Task

The Scout phase writes `loop/scout.md` but never creates a graph node. Each iteration Scout identifies exactly one gap — that gap should become a searchable KindClaim on the graph via `op=assert`. Without this, gaps are invisible to other agents, unsearchable via knowledge MCP, and lost when `scout.md` is overwritten next iteration.

## What Was Built

**The bridge:** `cmd/post/main.go` now calls `assertScoutGap()` at the end of each run. It reads `loop/scout.md`, parses the `**Gap:** ` line for the gap title and extracts the iteration number from the header, then POSTs `op=assert` to `/app/hive/op` creating a `KindClaim` node. The claim title is the gap text; the body includes both the iteration number and the gap title.

**Also fixed:** `pkg/runner/pipeline_state.go` was missing `"time"` in its import block (pre-existing compile error from the pipeline state machine commit).

**Files changed:**
- `cmd/post/main.go` — added `assertScoutGap()`, `extractIterationFromScout()`, `extractGapTitle()`; call from `main()`
- `cmd/post/main_test.go` — 4 new tests: `TestExtractGapTitle`, `TestExtractIterationFromScout`, `TestAssertScoutGapCreatesClaimNode`, `TestAssertScoutGapMissingFile`
- `pkg/runner/pipeline_state.go` — added `"time"` import (pre-existing bug)

**All tests pass.** After the next `close.sh` run, each Scout gap becomes a permanent `KindClaim` node on the hive graph and is searchable via `knowledge_search`.

ACTION: DONE
