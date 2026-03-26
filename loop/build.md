# Build: Expand marker candidates in `parseReflectorOutput` (start and boundary detection)

- **Commit:** d3188cb97786338c62a54a578c5a82608bf5634b
- **Subject:** [hive:builder] Expand marker candidates in `parseReflectorOutput` (start and boundary detection)
- **Cost:** $0.2920
- **Timestamp:** 2026-03-26T22:02:29Z

## Task

In `pkg/runner/reflector.go`, rewrite the marker-detection loop to try all format variants per key: `**KEY:**`, `**KEY**:`, `**KEY** :`, `### KEY:`, `## KEY:`, `KEY:`, and `strings.ToLower(key)+":"`. Pick the earliest-occurring match. Critically, also expand the end-of-section boundary detection (li...

## Diff Stat

```
commit d3188cb97786338c62a54a578c5a82608bf5634b
Author: hive <hive@lovyou.ai>
Date:   Fri Mar 27 09:02:29 2026 +1100

    [hive:builder] Expand marker candidates in `parseReflectorOutput` (start and boundary detection)

 loop/budget-20260327.txt     |  2 ++
 loop/build.md                | 37 ++++++++++++---------
 loop/reflections.md          | 10 ++++++
 loop/state.md                |  2 +-
 pkg/runner/reflector.go      | 48 ++++++++++++++++++---------
 pkg/runner/reflector_test.go | 78 ++++++++++++++++++++++++++++++++++++++++++++
 6 files changed, 146 insertions(+), 31 deletions(-)
```
