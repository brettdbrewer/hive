# Test Report: Prevent Fix: title compounding

- **Iteration:** 354
- **Timestamp:** 2026-03-28

## What Was Tested

The Builder added 8 tests covering the new dedup logic. I found 4 uncovered cases and added tests for them.

### Gaps found and filled

| Function | Missing coverage | New test |
|----------|-----------------|----------|
| `upgradeTaskPriority` | Zero tests — new function, completely untested | `TestUpgradeTaskPrioritySendsEditOp`, `TestUpgradeTaskPriorityAPIError` |
| `findExistingTask` | Empty `coreTitle` early-return path (no API call) | `TestFindExistingTaskEmptyCoreTitle` |
| `createTask` dedup | Single "Fix: " prefix against bare-title existing task — the most common real-world case | `TestCreateTaskDeduplicatesSingleFixPrefix` |

### Builder's 8 tests (all pass)

- `TestStripFixPrefixes` — strips x1, x2, x3 prefixes; no-op on clean title
- `TestAddTaskCommentSendsRespondOp` — op=respond with correct node_id and body
- `TestAddTaskCommentAPIError` — HTTP 403 returns error
- `TestFindExistingTaskMatchesCoreTitle` — board node titled "Fix: X" matches coreTitle "X"
- `TestFindExistingTaskNoMatch` — unrelated board node returns empty ID
- `TestCreateTaskDeduplicatesFixTask` — double-Fix: input, board has single-Fix: → comment, return existing ID
- `TestCreateTaskNoDedup` — non-Fix: title skips board query entirely
- `TestCreateTaskDeduplicatesBoardAPIError` — board 500 falls through to normal task creation

### My 4 additional tests (all pass)

- `TestUpgradeTaskPrioritySendsEditOp` — verifies op=edit, node_id, priority fields
- `TestUpgradeTaskPriorityAPIError` — HTTP 403 returns error
- `TestFindExistingTaskEmptyCoreTitle` — returns ("", nil) with zero API calls
- `TestCreateTaskDeduplicatesSingleFixPrefix` — single "Fix: X" input, bare "X" task exists → comment

## Results

```
ok  github.com/lovyou-ai/hive/cmd/post  0.807s
```

All 13 packages pass. 0 failures.

## Coverage notes

- `stripFixPrefixes`: fully covered
- `findExistingTask`: covered (match, no-match, empty input, API error via syncClaims)
- `addTaskComment`: covered (success, API error)
- `createTask` dedup path: covered (single-prefix, double-prefix, board API error fallback, non-Fix passthrough)
- `upgradeTaskPriority`: covered (success, API error)

The `upgradeTaskPriority` gap was the critical miss — new function, zero tests. The single-prefix dedup case was the most likely real-world regression point. Both are now pinned.
