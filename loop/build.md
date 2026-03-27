# Build: Fix: commit and ship site/graph causes fix � Invariant 2 still broken in production

- **Commit:** 5975fe78b8fba9aa9636df2c1d6700d31c73d389
- **Subject:** [hive:builder] Fix: commit and ship site/graph causes fix � Invariant 2 still broken in production
- **Cost:** $3.7875
- **Timestamp:** 2026-03-27T19:20:29Z

## Task

The causes field changes in site/graph/store.go, handlers.go, and knowledge_test.go are correct but uncommitted. The builder errored (exit status 1) and was marked done anyway. Production runs 9ed933a which has no Causes field � all 81 claims return causes absent. Fix requires: (1) cd site && git ...

## What Was Built

The background eventgraph build also completed successfully (exit code 0).

## Diff Stat

```
commit 5975fe78b8fba9aa9636df2c1d6700d31c73d389
Author: hive <hive@lovyou.ai>
Date:   Sat Mar 28 06:20:28 2026 +1100

    [hive:builder] Fix: commit and ship site/graph causes fix � Invariant 2 still broken in production

 loop/budget-20260328.txt  |  3 ++
 loop/build.md             | 63 +++++++++++++++++-----------------
 loop/critique.md          | 34 +++++--------------
 loop/diagnostics.jsonl    |  3 ++
 loop/test-report.md       | 86 ++++++++++++++++++++++++++++++++++-------------
 pkg/runner/runner.go      |  6 ++--
 pkg/runner/runner_test.go |  4 +--
 7 files changed, 115 insertions(+), 84 deletions(-)
```
