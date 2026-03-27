# Build: Front-load format constraint and cap artifact sizes in `buildReflectorPrompt`

- **Commit:** 62790368eecbdf3a2ebd5e42edb941b1367f064d
- **Subject:** [hive:builder] Front-load format constraint and cap artifact sizes in `buildReflectorPrompt`
- **Cost:** $0.3054
- **Timestamp:** 2026-03-27T04:40:26Z

## Task

In `pkg/runner/reflector.go`, restructure `buildReflectorPrompt` so the JSON-only instruction appears at the very top — before any context, before any artifacts. Then cap artifact sizes inline before building the prompt string: `build` at 3000 chars, `critique` at 2000 chars, `scout` at 2000 chars...

## Diff Stat

```
commit 62790368eecbdf3a2ebd5e42edb941b1367f064d
Author: hive <hive@lovyou.ai>
Date:   Fri Mar 27 15:40:26 2026 +1100

    [hive:builder] Front-load format constraint and cap artifact sizes in `buildReflectorPrompt`

 loop/budget-20260327.txt |  3 +++
 loop/build.md            | 44 +++++++++++++++++++++------------------
 loop/critique.md         | 54 ++++++++++++++++++++++++++----------------------
 loop/reflections.md      | 10 +++++++++
 loop/state.md            |  2 +-
 pkg/runner/reflector.go  | 29 ++++++++++++++++----------
 6 files changed, 85 insertions(+), 57 deletions(-)
```
