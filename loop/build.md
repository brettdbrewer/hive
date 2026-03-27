# Build: Fix Reflector prompt structure — front-load format constraint and cap artifacts

- **Commit:** (pending)
- **Cost:** (this iteration)
- **Timestamp:** 2026-03-27

## What changed

### `pkg/runner/reflector.go`

1. **`truncateArtifact`** — new helper that caps a string at `max` bytes, appending `\n... (truncated)` if cut.

2. **`buildReflectorPrompt`** — moved the output format constraint to the top of the prompt, before all artifact context. Previously it was buried after 8000+ chars of `sharedCtx` + artifacts, causing "lost in the middle" failures where the LLM entered essay mode instead of returning JSON.

3. **`runReflector`** — caps artifacts before building the prompt:
   - `scout.md` → 2000 bytes
   - `build.md` → 3000 bytes
   - `critique.md` → 2000 bytes
   - `sharedCtx` → 4000 bytes

### `pkg/runner/reflector_test.go`

4. **`TestBuildReflectorPrompt`** — added assertion that "Return ONLY" appears before the first artifact section header (guards against regression where format gets re-buried).

5. **`TestTruncateArtifact`** — four cases: short string unchanged, exact limit unchanged, over-limit with marker, empty string.

## Root cause addressed

Scout report identified three coordinated fixes. The previous iteration (5641a3b) did only the model switch (haiku → sonnet). This iteration completes the other two: prompt reordering and artifact capping. Together they address the documented root cause — format instruction buried after 8000+ chars of context — that caused nine consecutive `empty_sections` failures.
