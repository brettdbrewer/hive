# Build: Fix commit subject: strip [hive:*] prefix from task title in commitAndPush

- **Commit:** 9695f290ebf4540222cde1de2e7517ec7078751e
- **Subject:** [hive:builder] Fix commit subject: strip [hive:*] prefix from task title in commitAndPush
- **Cost:** $0.3441
- **Timestamp:** 2026-03-27T08:35:10Z

## Task

In pkg/runner/runner.go, fix commitAndPush() around line 525. The current msg := fmt.Sprintf("[hive:%s] %s", r.cfg.Role, t.Title) embeds whatever is in t.Title — which can contain prior [hive:builder] prefixes from the board, producing recursive nesting like [hive:builder] Fix: [hive:builder] Fix:...

## What Was Built

Done. Two files changed:

- **`pkg/runner/runner.go`**: Added `stripHivePrefix` helper that loops stripping `[hive:xxx]` prefixes, and updated `commitAndPush` to use it.
- **`pkg/runner/runner_test.go`**: Added `TestStripHivePrefix` covering no-prefix, single-prefix, and double-nested cases.

Build and tests both clean.

ACTION: DONE

## Diff Stat

```
commit 9695f290ebf4540222cde1de2e7517ec7078751e
Author: hive <hive@lovyou.ai>
Date:   Fri Mar 27 19:35:10 2026 +1100

    [hive:builder] Fix commit subject: strip [hive:*] prefix from task title in commitAndPush

 loop/budget-20260327.txt  |  3 +++
 loop/build.md             | 63 +++++++++++---------------------------------
 loop/critique.md          | 67 +++++++++++++++++++++++++++++++----------------
 loop/reflections.md       | 10 +++++++
 loop/state.md             |  2 +-
 pkg/runner/runner.go      | 15 ++++++++++-
 pkg/runner/runner_test.go | 20 ++++++++++++++
 7 files changed, 108 insertions(+), 72 deletions(-)
```
