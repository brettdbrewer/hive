# Build: Fix Critic issues — Agent Memory Phase 4 (importance clamp + deploy)

## What Changed

### `site/graph/mind.go`

**Bug fixed — importance range guard in `extractAndSaveMemories()`:**

Added a clamp to ensure importance stays within 1-5 (matching the LLM prompt's instruction) before calling `RememberForPersona()`. The LLM output is an external boundary (BOUNDED invariant 13) and must be validated before storage.

```go
// Clamp importance to 1-5 (matches prompt instructions; LLM output is an external boundary).
if e.Importance < 1 {
    e.Importance = 1
} else if e.Importance > 5 {
    e.Importance = 5
}
```

Previously, out-of-range values (0, -1, 10, etc.) passed straight to `RememberForPersona()`. The store normalises values outside 1-10 to 5, but the prompt contracts a 1-5 range — so e.g. a raw value of `8` would have been stored unchecked.

### `site/graph/memory_test.go`

**Test added — `TestImportanceClampLogic` (no DB required):**

Pure logic test verifying the clamp function covers all cases: below 1 → 1, above 5 → 5, in-range stays unchanged. Runs in CI without DATABASE_URL.

**Test added — `TestExtractMemoriesImportanceClamp` (requires DB):**

Integration test verifying the clamped value (raw=10 → clamped=5) is accepted by the store and recalled correctly.

## Files Changed

- `site/graph/mind.go` — importance clamp in `extractAndSaveMemories()`
- `site/graph/memory_test.go` — `TestImportanceClampLogic` + `TestExtractMemoriesImportanceClamp`

## Verification

- `go build -buildvcs=false ./...` — clean, no errors
- `go test ./...` — all pass (TestImportanceClampLogic passes without DB)
- `./ship.sh "iter 234: fix importance range guard + add clamp test"` — deployed successfully
  - Fly.io machines were in stuck "replacing" state from a prior concurrent deploy; destroyed both and redeployed fresh
  - Deployment confirmed: https://lovyou-ai.fly.dev/
