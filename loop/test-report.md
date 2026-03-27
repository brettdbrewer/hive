# Test Report: Observer audit — kind=task fix

- **Date:** 2026-03-28
- **Build commit:** d062e08
- **Tests run:** 25 in `cmd/post/` (23 existing + 2 new)
- **Result:** PASS

## What Was Tested

### New tests added

**`TestCreateTaskSendsKindTask`**
The critical regression pin for this iteration. `createTask()` was changed to include
explicit `"kind": "task"` in the intend payload. Without this, the board kind assignment
depended on server defaults rather than the client explicitly requesting it. This test
captures the intend request and asserts `kind == "task"` is present.

**`TestAssertCritiqueNoTitle`**
Verifies `assertCritique()` returns an error (mentioning "critique title") when
`critique.md` exists but contains no markdown heading. Parallel to the existing
`TestAssertScoutGapNoGapLine` for the scout path.

### Coverage confirmed for all new/changed functions

| Function | Tests |
|---|---|
| `createTask()` — kind=task field | `TestCreateTaskSendsKindTask` (new) |
| `post()` — op=intend, kind=document | `TestPostCreatesDocument`, `TestBuildTitleExtractedOnPost` |
| `assertCritique()` | `TestAssertCritiqueCreatesClaimNode`, `TestAssertCritiqueMissingFile`, `TestAssertCritiqueNoTitle` (new) |
| `assertLatestReflection()` | `TestAssertLatestReflectionCreatesDocument`, `TestAssertLatestReflectionMissingFile` |
| `extractCritiqueTitle()` | `TestExtractCritiqueTitle` (3 cases) |
| `extractLatestReflection()` | `TestExtractLatestReflection`, `TestExtractLatestReflectionNoEntry` |

## Gap Found

The Builder wrote 23 tests but missed the primary regression case:
`createTask()` had no test verifying `kind=task` was in the payload. This was the
core fix of the iteration — the 491 board nodes lacked explicit kind because the
field wasn't being sent. `TestCreateTaskSendsKindTask` pins it.

## Full Suite

```
ok  github.com/lovyou-ai/hive/cmd/post  0.560s  (25 tests)
ok  github.com/lovyou-ai/hive/pkg/runner  3.522s
```

All packages clean. No regressions.
