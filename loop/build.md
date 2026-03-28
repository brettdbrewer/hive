# Build: Fix: hardcoded one-time upgradeTaskPriority call permanently in main()

- **Commit:** 077c4c1decbc4789720e69e86277bf3b8ecf0347
- **Subject:** [hive:builder] Fix: hardcoded one-time upgradeTaskPriority call permanently in main()
- **Cost:** $0.1181
- **Timestamp:** 2026-03-28T12:45:30Z

## Task

The retroactive priority upgrade for task 468e0549 (lines 160-165 in cmd/post/main.go) fires on every invocation of cmd/post � not just once. A magic string ("468e0549") is permanently embedded in the production execution path. This violates the BOUNDED invariant (infinite scope for a one-time ope...

## What Was Built

All tests pass. Removed the hardcoded `upgradeTaskPriority("468e0549", ...)` call from `main()`. The `upgradeTaskPriority` function itself is kept — it's a valid general-purpose helper with test coverage.

ACTION: DONE

## Diff Stat

```
commit 077c4c1decbc4789720e69e86277bf3b8ecf0347
Author: hive <hive@lovyou.ai>
Date:   Sat Mar 28 23:45:30 2026 +1100

    [hive:builder] Fix: hardcoded one-time upgradeTaskPriority call permanently in main()

 cmd/post/main.go         | 7 -------
 loop/budget-20260328.txt | 1 +
 2 files changed, 1 insertion(+), 7 deletions(-)
```
