# Build Report — Iteration 227: Scout Role

## What This Iteration Does

Implements the Scout role for the hive runtime (Phase 2, item 8). The Scout reads project state, identifies gaps, and creates concrete tasks on the board for the Builder. This closes the autonomous loop: Scout → Builder → Critic.

## Files Changed

| File | Lines | What |
|------|-------|------|
| `pkg/runner/scout.go` | 175 | New. Scout role: read state.md + git log + board → Reason() → create task. |
| `pkg/runner/scout_test.go` | 65 | New. 4 tests: parseScoutTask, parseScoutTaskDefaults, parseScoutTaskEmpty, buildScoutPrompt. |
| `pkg/runner/runner.go` | +2/-5 | Added HiveDir to Config. Removed scout stub. |
| `cmd/hive/main.go` | +1 | Pass HiveDir to runner config. |

## How It Works

1. **Every 8th tick** (~2 minutes), Scout checks agent's open task count
2. If agent has < 3 tasks, proceeds with scouting
3. Gathers context: `state.md` (os.ReadFile), `git log -20`, board summary via API
4. Calls `Reason()` (haiku, fast, cheap) with scouting prompt
5. Parses response for `TASK_TITLE:`, `TASK_PRIORITY:`, `TASK_DESCRIPTION:`
6. Creates task on lovyou.ai board via API

## E2E Test Result

```
[scout] tick 8: scouting (agent has 0/3 tasks)
  ⏳ thinking done (38s)
[scout] Reason done (cost=$0.0611)
[scout] no task found in response          ← first call: unstructured output
  ⏳ thinking done (33s)
[scout] Reason done (cost=$0.0185)
[scout] creating task: [high] Integrate Scout phase into hive runner Execute() path
[scout] created task 3d77ba43
[scout] cost summary: $0.0795 / $1.00 (calls=2)
```

Scout created a concrete task after 2 calls (~$0.08). First call didn't produce structured output; second did. Throttle correctly blocked when agent had 4 tasks (> 3 max).

## Build

- `go build ./...` ✓
- `go test ./...` ✓ (27 tests: 4 scout + 9 critic + 14 runner)
