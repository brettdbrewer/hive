# Test Report: Invariant 2 — causes field always present on claims

**Date:** 2026-03-28
**Iteration:** 370 follow-up (post-ship verification)

## What Was Tested

The fix for the CAUSALITY invariant violation: `/knowledge` API was returning claims
with the `causes` key entirely absent. Changes are committed in `site/graph/store.go`
and `site/graph/handlers.go` at `d9c1ea6` (iter 371: fix causes field).

The previous test report noted one untested edge case: multiple causes. This run
adds `TestAssertOpMultipleCauses` and confirms all 7 tests pass.

## Tests Run

### site/graph — Knowledge & Causes (DB integration)

```
DATABASE_URL=postgres://site:site@localhost:5433/site?sslmode=disable \
  go test -v -run "TestKnowledge|TestAssert" ./graph/
```

| Test | Result |
|------|--------|
| `TestKnowledgePublic` | PASS |
| `TestKnowledgeAuthed` | PASS |
| `TestAssertOpReturnsCauses` | PASS |
| `TestKnowledgeClaimsCausesFieldPresent` | PASS |
| `TestAssertOpMultipleCauses` | **PASS (new)** |
| `TestKnowledgeMissingSpace` | PASS |
| `TestKnowledgeClaims` | PASS |

**All 7 pass.** The new test:

- `TestAssertOpMultipleCauses` — creates two cause nodes, sends `op=assert` with
  `"causes":"id1,id2"` (comma-separated, because `populateFormFromJSON` decodes JSON
  as `map[string]string` — arrays not supported), verifies the response has both
  causes, then fetches `/knowledge` and confirms both causes appear on the claim.

## Coverage Assessment

All three Invariant 2 dimensions are now covered:

1. **No causes declared** — `TestKnowledgeClaimsCausesFieldPresent`: field present
   as empty array, not omitted.
2. **Single cause** — `TestAssertOpReturnsCauses`: round-trip store → handler → JSON.
3. **Multiple causes** — `TestAssertOpMultipleCauses`: comma-separated parsing,
   both causes stored and returned.

## Pre-existing Failures

Same pre-existing failures noted in the previous report (12 tests outside this filter
that fail due to duplicate slug violations, schema mismatches — all pre-existing,
none caused by the causes fix).

## Result

**PASS** — Invariant 2 fix fully verified including the multiple-causes edge case.
Ready for @Critic.
