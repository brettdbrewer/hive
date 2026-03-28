# Build: Structural hardening: loop task parser should detect and reject meta-tasks programmatically (Lesson 137 level 2)

- **Commit:** 55d4214c7de113b18d70fda67357e56e95113afc
- **Subject:** [hive:builder] Structural hardening: loop task parser should detect and reject meta-tasks programmatically (Lesson 137 level 2)
- **Cost:** $0.3071
- **Timestamp:** 2026-03-28T15:45:12Z

## Task

Lesson 137 names two hardening levels for the meta-task defect: (1) prompt prohibition (done in 24924471) and (2) structural post-processing. Level 2 is not yet implemented. The loop task-command parser in pkg/loop/loop.go already parses /task create commands from agent output. Add a guard: before e...

## What Was Built

All 15 tests pass. Here's what was implemented:

**`pkg/loop/tasks.go`** — two additions:

1. `metaTaskPatterns` — a package-level slice of the four meta-task indicator phrases (`op=complete`, `close task`, `mark done`, `close the following`).

2. `isMetaTaskBody(title, description string) bool` — case-insensitive check over the combined title+description against those patterns.

3. Guard in `execTaskCreate` — called before `tasks.Create`. On match, logs a warning to stdout and returns a descriptive error, causing the command to be counted as failed in `executeTaskCommands`.

**`pkg/loop/tasks_test.go`** (new file):

- `TestIsMetaTaskBody` — 18 cases covering all four patterns, case-insensitivity, pattern in description vs title, and legitimate task descriptions.
- `TestParseTaskCommandsMetaTaskNotFiltered` — confirms `parseTaskCommands` still surfaces the command so the guard can fire (the filter is structural, not at parse time).

ACTION: DONE

## Diff Stat

```
commit 55d4214c7de113b18d70fda67357e56e95113afc
Author: hive <hive@lovyou.ai>
Date:   Sun Mar 29 02:45:12 2026 +1100

    [hive:builder] Structural hardening: loop task parser should detect and reject meta-tasks programmatically (Lesson 137 level 2)

 loop/budget-20260329.txt |  1 +
 pkg/loop/tasks.go        | 26 ++++++++++++++++++++++
 pkg/loop/tasks_test.go   | 58 ++++++++++++++++++++++++++++++++++++++++++++++++
 3 files changed, 85 insertions(+)
```
