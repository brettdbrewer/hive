# Build Report — Iteration 228: Pipeline Mode

## What This Iteration Does

Adds `--pipeline` mode to cmd/hive: one command runs Scout → Builder → Critic in sequence. Phase 2 item 11.

## Files Changed

| File | What |
|------|------|
| `cmd/hive/main.go` | Added `--pipeline` flag, `runPipeline()` function. Runs Scout → Builder → Critic with shared context. |
| `pkg/runner/scout.go` | Fixed tick throttle bypass in one-shot mode (linter fix). |
| `pkg/runner/critic.go` | Fixed tick throttle bypass in one-shot mode (linter fix). |
| `pkg/runner/scout_test.go` | Added `TestScoutThrottleBypassInOneShot`. |
| `pkg/runner/critic_test.go` | Added `TestCriticThrottleBypassInOneShot`. |

## Pipeline E2E Test

```bash
go run ./cmd/hive --pipeline --repo ../site --space hive --agent-id ... --budget 5
```

```
[pipeline] ── scout ──     Created task ($0.05, 38s)
[pipeline] ── builder ──   Claimed, Operated ($0.77, 4m16s) → DONE, no changes
[pipeline] ── critic ──    Reviewed Policy commit → PASS ($0.32, 2m29s)
[pipeline] ── cycle complete ──
Total: $1.14, ~8 minutes
```

## Issue: Repo Mismatch

The Scout created a hive infrastructure task ("Add --pipeline mode to cmd/hive") but the Builder operates on the site repo. The Builder correctly Operated but found nothing to change in `../site`. The changes-required guard caught this and left the task in-progress.

**Root cause:** The Scout reads state.md from the hive repo and identifies hive gaps. But the Builder builds on the site repo. The Scout needs to understand which repo it's scouting for and create tasks accordingly.

**Fix:** Pass the target repo context to the Scout prompt. Or: run the pipeline with `--repo ../hive` for hive tasks, `--repo ../site` for site tasks.

## Build

- `go build ./...` ✓
- `go test ./...` ✓ (29 tests)
