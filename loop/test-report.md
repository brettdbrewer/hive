# Test Report: 65d1e553 false completion — operational iteration

- **Date:** 2026-03-28
- **Build commit:** 388a9bb
- **Tests run:** 28 top-level in `cmd/post/` (25 existing + 5 new)
- **Result:** PASS

## What Was Tested

This iteration made no production code changes — it was operational (ran cmd/post manually, closed orphaned tasks). The regression pins from the d062e08 fix (TestCreateTaskSendsKindTask, TestAssertCritiqueNoTitle) were already present. This session closed the remaining coverage gaps.

### New tests added

**`TestEnsureSpaceExisting`**
Verifies ensureSpace returns nil when the API responds 200 — space exists, no POST.

**`TestEnsureSpaceCreates`**
Verifies ensureSpace POSTs to `/app/new` with `kind=community` when the space is missing (404 response). Was the only function that made two distinct HTTP calls with no test at all.

**`TestEnsureSpaceCreateError`**
Verifies ensureSpace returns an error when the create POST fails with HTTP 4xx.

**`TestSyncMindStateSuccess`**
Verifies syncMindState sends a PUT to `/api/mind-state` with the correct Authorization header, `key=loop_state`, and the full state content as the value.

**`TestSyncMindStateError`**
Verifies syncMindState returns an error when the server responds with HTTP 4xx.

### Coverage after this session

| Function | Tests |
|---|---|
| `buildTitle()` | `TestBuildTitle` (6 subtests) |
| `post()` | `TestPostCreatesDocument`, `TestBuildTitleExtractedOnPost` |
| `createTask()` | `TestCreateTaskSendsKindTask` |
| `ensureSpace()` | `TestEnsureSpaceExisting`, `TestEnsureSpaceCreates`, `TestEnsureSpaceCreateError` ✨ |
| `syncMindState()` | `TestSyncMindStateSuccess`, `TestSyncMindStateError` ✨ |
| `syncClaims()` | `TestSyncClaimsWritesFile`, `TestSyncClaimsEmptyDoesNotWrite`, `TestSyncClaimsAPIError`, `TestSyncClaimsClaimWithNoMetadata` |
| `assertScoutGap()` | `TestAssertScoutGapCreatesClaimNode`, `TestAssertScoutGapMissingFile`, `TestAssertScoutGapNoGapLine`, `TestAssertScoutGapAPIError`, `TestAssertScoutGapSendsAuthHeader` |
| `extractIterationFromScout()` | `TestExtractIterationFromScout` (3 subtests) |
| `extractGapTitle()` | `TestExtractGapTitle` (3 subtests) |
| `assertCritique()` | `TestAssertCritiqueCreatesClaimNode`, `TestAssertCritiqueMissingFile`, `TestAssertCritiqueNoTitle` |
| `extractCritiqueTitle()` | `TestExtractCritiqueTitle` (3 subtests) |
| `assertLatestReflection()` | `TestAssertLatestReflectionCreatesDocument`, `TestAssertLatestReflectionMissingFile` |
| `extractLatestReflection()` | `TestExtractLatestReflection`, `TestExtractLatestReflectionNoEntry` |

## Full Suite

```
ok  github.com/lovyou-ai/hive/cmd/post       0.566s  (28 top-level tests)
ok  github.com/lovyou-ai/hive/cmd/mcp-graph  (cached)
ok  github.com/lovyou-ai/hive/cmd/mcp-knowledge (cached)
ok  github.com/lovyou-ai/hive/pkg/api        (cached)
ok  github.com/lovyou-ai/hive/pkg/authority  (cached)
ok  github.com/lovyou-ai/hive/pkg/hive       (cached)
ok  github.com/lovyou-ai/hive/pkg/loop       (cached)
ok  github.com/lovyou-ai/hive/pkg/resources  (cached)
ok  github.com/lovyou-ai/hive/pkg/runner     (cached)
ok  github.com/lovyou-ai/hive/pkg/workspace  (cached)
```

All packages clean. No regressions.

## Remaining gap

`main()` itself — not directly testable (reads env vars, calls os.Exit). Covered indirectly by all the function-level tests above. No action needed.
