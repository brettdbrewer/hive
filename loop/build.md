# Build: Update Critic prompt — enforce Scout-gap cross-reference and degenerate-iteration REVISE

## Gap Addressed
Lessons 168, 171, 197, 200, 201 converged on a single root cause: the Critic prompt was not enforcing two critical checks, allowing it to ratify drift across 16 consecutive iterations.

## What Changed

### `agents/critic.md`
- Added **Scout gap cross-reference** (Lessons 168/171) to "What You Produce": Critic must REVISE when `build.md` does not explicitly reference the open gap from `scout.md`.
- Added **Degenerate iteration** (Lesson 200) to "What You Produce": Critic must REVISE when the diff only contains `loop/` artifact changes.
- Added two corresponding anti-patterns to reinforce both rules.

### `loop/critic-prompt.txt`
- Added `SCOUT GAP CROSS-REFERENCE` section with explicit REVISE instruction.
- Added `DEGENERATE ITERATION CHECK` section with explicit REVISE instruction.
- Both sections reference the lessons that motivated them.

### `pkg/runner/critic.go`
- `buildReviewPrompt`: Added `scoutContent` and `buildContent` parameters; both checks are items 1 and 2 in the review checklist.
- Operate path instruction: Added `Required Checks (Lessons 168/171/200)` section with both rules.
- Reason path call site: Loads `loop/scout.md` and `loop/build.md` via new `loadLoopArtifact` helper and passes them into the prompt.
- New `isDegenerateIteration(diff)` helper: detects when every file in a diff is under `loop/`.
- New `loadLoopArtifact(hiveDir, name)` helper: reads a loop artifact, capped at 3000 bytes.

### `pkg/runner/critic_test.go`
- Updated `TestBuildReviewPrompt` signature and added assertions for both new checklist items.
- Added `TestBuildReviewPromptWithArtifacts`: verifies scout/build content appears in prompt.
- Added `TestIsDegenerateIteration`: covers all-loop, mixed, empty, and no-loop cases.

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (pkg/runner: 3.724s)
