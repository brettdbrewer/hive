# Test Report: Meta-task structural guard (Lesson 137 level 2)

- **Build:** 55d4214
- **Tester:** Tester agent
- **Date:** 2026-03-29

## What Was Tested

The Builder shipped `isMetaTaskBody` and the `execTaskCreate` guard, with 15 tests pre-existing (13 loop + 2 new). The Tester added 2 more tests to close coverage gaps.

## New Tests Added (`pkg/loop/tasks_test.go`)

### `TestExecTaskCreateRejectsMetaTask` (4 subtests)
**Gap it fills:** `isMetaTaskBody` was tested in isolation, but the guard wired into `execTaskCreate` had no direct test. This test exercises the full rejection path inside `execTaskCreate` — not just the pure function.

- All four patterns rejected from title field
- Case-insensitive rejection (Close Task, not just close task)
- Rejection when pattern is in description, not title
- Error message contains `"meta-task rejected"` in all cases

### `TestIsMetaTaskBodyTitleDescriptionJoin`
**Gap it fills:** Documents and verifies a subtle boundary behaviour — the join is `title + " " + description`, so a pattern can span the boundary (title=`"close the"`, description=`"following tasks"` → joined matches `"close the following"`). Also confirms unrelated fragment pairs don't false-positive.

## Results

| Test | Status |
|------|--------|
| TestIsMetaTaskBody (18 cases) | PASS |
| TestParseTaskCommandsMetaTaskNotFiltered | PASS |
| TestExecTaskCreateRejectsMetaTask/op=complete_in_title | PASS |
| TestExecTaskCreateRejectsMetaTask/close_task_in_title | PASS |
| TestExecTaskCreateRejectsMetaTask/mark_done_in_description | PASS |
| TestExecTaskCreateRejectsMetaTask/close_the_following_in_description | PASS |
| TestIsMetaTaskBodyTitleDescriptionJoin | PASS |
| All 13 pre-existing loop tests | PASS |

**Total: 17 tests, 0 failures**

## Coverage Notes

- `isMetaTaskBody` — all four patterns, both fields, case-insensitivity, boundary join: **covered**
- `execTaskCreate` rejection path — all four patterns: **covered**
- `execTaskCreate` happy path (reaching `tasks.Create`) — untestable without a real `TaskStore`; the `TestIsMetaTaskBody` negative cases provide the same guarantee at the unit level
- `parseTaskCommands` pass-through: **covered**

## @Critic

Ready for review. The `TestExecTaskCreateRejectsMetaTask` test is the substantive addition — it proves the guard fires at the correct point in `execTaskCreate`, not just that `isMetaTaskBody` returns the right bool.
