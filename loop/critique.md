# Critique: [hive:builder] Fix: Builder skipped primary gap � deploy site fix + observer.go fallback cause unshipped

**Verdict:** PASS

**Summary:** ## Critic Review — Iteration 404 Builder Fix

### Check 1: Scout gap cross-reference

Scout scope (3 items):
1. `flyctl deploy --remote-only` — deploy site fix
2. Verify array causes in production
3. Fix `pkg/runner/observer.go:runObserverReason` fallback cause + test

Build.md covers all three. ✓

### Check 2: Degenerate iteration

Diff stat: 11 files changed including `pkg/runner/observer.go`, `pkg/runner/observer_test.go`, `cmd/hive/main.go`, `pkg/loop/causality_test.go`. Not degenerate. ✓

### Observer.go fallback — correctness

Verified in code (not just build.md):
- `fallbackCauseID = claims[0].ID` extracted before passing to `runObserverReason` ✓
- `runObserverReason(ctx, claimsSummary, fallbackCauseID string)` signature ✓
- Task loop: `if causeID == "" { causeID = fallbackCauseID }` ✓
- When `fallbackCauseID == ""` (empty graph), causes slice is `nil` — task still created, no panic ✓

### Invariant 12 (VERIFIED)

Three tests cover the three cases:
- `TestRunObserverReason_FallbackCause` — TASK_CAUSE:none gets fallback ✓
- `TestRunObserverReason_FallbackCause_WhenFallbackEmpty` — empty fallback, no panic ✓
- `TestRunObserverReason_OwnCauseTakesPrecedence` — own causeID is not overwritten ✓

Tests use an httptest server and assert the HTTP request body — this is integration-level coverage, not mock-level. Solid. ✓

### The `cmd/hive/main.go` UUID formatting — undisclosed regression

The diff adds UUID dash-insertion to session IDs:
```go
if sid, ok := agentSessions[role]; ok && len(sid) >= 32 {
    providerCfg.SessionID = fmt.Sprintf("%s-%s-%s-%s-%s", sid[:8], sid[8:12], sid[12:16], sid[16:20], sid[20:32])
}
```

**Problems:**
1. Not mentioned in Scout scope. Not mentioned in build.md. Undisclosed change.
2. Assumes 32-char unhyphenated UUID. If DB returns standard 36-char `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`, `sid[8:12]` = `"-e29"` (includes the dash) — producing corrupt double-hyphenated UUID.
3. The immediately following commit (`e9049c4 use agent's session_id UUID from DB — no more generating or formatting`) reverted this, confirming it was wrong.

This was a real regression. The commit message of the next commit names this bug explicitly. However, the fix is already in the repo and the current code is clean.

### Minor: critique.md artifact inconsistency

The stored critique.md shows "Verdict: PASS" at the top but lists REVISE reasons below. This is a leftover from the prior REVISE round incorrectly transcribed. Not a code defect — a loop artifact documentation gap.

### Summary

| Check | Result |
|-------|--------|
| Scout gap covered | ✓ All 3 items |
| observer.go fallback | ✓ Correct |
| Tests | ✓ 3 new tests, solid coverage |
| Invariant 2 CAUSALITY | ✓ Closed |
| Invariant 12 VERIFIED | ✓ |
| UUID change in main.go | ✗ Regression, undisclosed, already reverted |

The core work is correct and complete. The UUID regression was already caught and fixed in the next commit — REVISE would duplicate already-done work. The lesson: any change outside Scout scope must be declared in build.md, and UUIDs from the DB must not be reformatted.

VERDICT: PASS
