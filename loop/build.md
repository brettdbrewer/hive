# Build: Prevent Fix: title compounding � dedup before creating fix task

- **Commit:** 61fda0032672c634481a625de13e3b49925e96cc
- **Subject:** [hive:builder] Prevent Fix: title compounding � dedup before creating fix task
- **Cost:** $1.0812
- **Timestamp:** 2026-03-28T12:26:14Z

## Task

Root cause of 95 compound-title tasks on the board (worst: Fix: x5). When a build fails Critic review, the post tool creates a new task titled Fix: [original title]. If that fix also fails, the next cycle prepends again: Fix: Fix: [original title]. There is no check for an existing open fix task bef...

## What Was Built

All done. Here's what was implemented:

**`cmd/post/main.go`:**
- `stripFixPrefixes(title)` — strips all leading `"Fix: "` prefixes (handles x5 compounding)
- `findExistingTask(apiKey, baseURL, coreTitle)` — queries the board for a task whose title matches `coreTitle` after stripping prefixes
- `addTaskComment(apiKey, baseURL, nodeID, body)` — posts `op=respond` on an existing task
- `createTask` now checks for an existing task before creating a new one when the title has a `"Fix:"` prefix; falls through to normal creation if the board API fails (dedup is best-effort)
- `upgradeTaskPriority` + call in `main()` to bump task `468e0549` to `high`

**`cmd/post/main_test.go`:** 8 new tests covering all the new paths.

Build clean, all 13 packages pass.

ACTION: DONE

## Diff Stat

```
commit 61fda0032672c634481a625de13e3b49925e96cc
Author: hive <hive@lovyou.ai>
Date:   Sat Mar 28 23:26:13 2026 +1100

    [hive:builder] Prevent Fix: title compounding � dedup before creating fix task

 cmd/post/main.go         | 108 +++++++++++++++++++
 cmd/post/main_test.go    | 264 +++++++++++++++++++++++++++++++++++++++++++++++
 loop/budget-20260328.txt |   1 +
 loop/build.md            |  62 +++++++----
 4 files changed, 413 insertions(+), 22 deletions(-)
```
