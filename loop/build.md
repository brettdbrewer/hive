# Build: Add regression tests for JSON format and Preview field

- **Commit:** ce363bb8aae8232629d02d2f8302d0aa02abf417
- **Subject:** [hive:builder] Add regression tests for JSON format and Preview field
- **Cost:** $0.1749
- **Timestamp:** 2026-03-26T22:57:06Z

## Task

In `pkg/runner/architect_test.go`, add test cases to `TestParseArchitectSubtasks` covering: JSON array input `[{"title":"...","description":"...","priority":"high"}]`, `{"tasks":[...]}` wrapper, and a prose preamble followed by SUBTASK_TITLE markers (verifying preamble doesn't corrupt parsing). Also...

## Diff Stat

```
commit ce363bb8aae8232629d02d2f8302d0aa02abf417
Author: hive <hive@lovyou.ai>
Date:   Fri Mar 27 09:57:06 2026 +1100

    [hive:builder] Add regression tests for JSON format and Preview field

 loop/budget-20260327.txt     |  3 +++
 loop/build.md                | 47 ++++++++++++++++++++++++++----------
 loop/critique.md             | 57 +++++++++++++++-----------------------------
 loop/reflections.md          | 14 +++++++++++
 loop/state.md                |  2 +-
 pkg/runner/architect_test.go | 32 +++++++++++++++++++++++++
 6 files changed, 103 insertions(+), 52 deletions(-)
```
