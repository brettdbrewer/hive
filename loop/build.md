# Build Report — Iteration 280 Fix: reflections.md artifact cleanup

## Gap
Critic review of commit 5f85ef8 found three defects in `loop/reflections.md`. No code defects — the daemon error recovery logic was correct per Critic. All defects were artifact integrity issues.

## Changes

### `loop/reflections.md`

**Defect 1 (conversational artifact):** Removed "Do you want me to append this to `loop/reflections.md` and update `loop/state.md` to iteration 280?" — a Reflector meta-question appended verbatim to the append-only audit log.

**Defect 2 (stray code fence):** Removed ` ``` ` with no matching opener that broke markdown structure.

**Defect 3 (content mismatch):** The two defects above were masking the correct iteration 280 reflection at the end of the file. After cleanup, the final entry correctly describes what was actually built this iteration (daemon error recovery), not the budget tracker from iteration 279. Also removed trailing "Shall I append this..." artifact.

The correct iteration 280 COVER (already present, now visible):
> Builder implemented error recovery for daemon loop (consecFailures counter, backoff retry logic), fixed budget tracking in RunCouncil, added nil-safety guards. Daemon infrastructure now has fault tolerance.

## Verification

- `go.exe build -buildvcs=false ./...` — ✓ clean
- `go.exe test ./...` — ✓ all pass
- `reflections.md` — no stray fences, no conversational artifacts, correct iteration 280 COVER

---

# Build Report — Daemon error recovery + status file

## Task
Add error recovery to daemon loop so pipeline failures don't halt continuous operation.

## Changes

### `cmd/hive/main.go` — `runDaemon()`

- Added constants `daemonMaxConsecFailures = 3` and `daemonBackoffInterval = 5 * time.Minute`.
- Added `consecFailures` counter alongside the existing `cycle` counter.
- On `runPipeline()` error:
  - Increments `consecFailures`.
  - Logs prominently with `████` delimiters showing `N/3 consecutive`.
  - Writes `cycle=N error: ...` to `loop/daemon.status`.
  - If `consecFailures >= 3`, returns a wrapped error (halt with clear message).
  - Otherwise sleeps `daemonBackoffInterval` (5 min) before retrying, respecting ctx cancellation.
  - Uses `continue` to skip the normal interval wait.
- On success:
  - Resets `consecFailures` to 0.
  - Writes `cycle=N ok` to `loop/daemon.status`.
  - Sleeps the normal `interval` before the next cycle.
- Added `writeDaemonStatus(path, line)` helper — writes one line + newline, logs a warning on failure (non-fatal).
- Status file path: `<hiveDir>/loop/daemon.status` (resolved via `findHiveDir()`).

## Verification

- `go.exe build -buildvcs=false ./...` — clean, no errors.
- `go.exe test ./...` — all pass.

---

# Build Report — Fix RunCouncil budget tracking + nil-safety

## Gap
Two issues from Critic review of the daily budget tracker:

1. `RunCouncil` accumulated `totalCost` but never wrote it to the daily budget file — cross-invocation spend from council sessions was invisible to the ceiling check.
2. `Spent()` and `Remaining()` lacked nil guards, breaking the "Safe to call on nil" contract advertised in the comments.

## Changes

### `pkg/runner/council.go`
- Added `NewDailyBudget(cfg.HiveDir).Record(totalCost)` immediately after `wg.Wait()`.
- Council sessions now contribute to cross-invocation spend tracking like all other runner methods.

### `pkg/runner/budget.go`
- Added nil guard to `Spent()`: returns `0` if receiver is nil.
- Added nil guard to `Remaining()`: returns `ceiling` if receiver is nil (consistent with "no spend known" semantics).
- Updated doc comments to advertise nil-safety on both methods.

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test -buildvcs=false ./...` — all pass (pkg/runner: 0.447s)

---

# Previous Build Report — Daily Budget Tracker

## Task
Add file-backed daily budget tracker to `pkg/runner/budget.go`

## Files Changed

### Created: `pkg/runner/budget.go`
- `DailyBudget` type with `dir` (hiveDir) and `date` (YYYYMMDD) fields
- `NewDailyBudget(hiveDir string) *DailyBudget` constructor
- `Record(amount float64)` — appends float as a line to `loop/budget-YYYYMMDD.txt`
- `Spent() float64` — reads and sums all lines in today's file
- `Remaining(ceiling float64) float64` — returns `max(0, ceiling - Spent())`

### Created: `pkg/runner/budget_test.go`
- `TestDailyBudgetRoundTrip` — Record, Spent, Remaining round-trip including over-ceiling clamp
- `TestDailyBudgetPersistence` — two DailyBudget instances against the same tmpdir to verify persistence across process restarts

### Modified: `pkg/runner/runner.go`
- Added `dailyBudget *DailyBudget` field to `Runner`
- `New()` initialises `dailyBudget` from `cfg.HiveDir`
- `Run()` loop: after the in-process `IsOverBudget()` check, reads `dailyBudget.Remaining(ceiling)` — if ≤ 0, logs and sleeps one interval then continues (does not stop the process; resets naturally the next day)
- `workTask()`: calls `r.dailyBudget.Record(result.Usage.CostUSD)` alongside the existing `r.cost.Record()`

## Build Results

```
go.exe build -buildvcs=false ./...   -> OK
go.exe test ./...                    -> OK (pkg/runner 0.649s, all pass)
```

---

# Previous Build Report — Iteration 271 Fix: Artifact Cleanup

**Date:** 2026-03-26
**Gap:** Critic review of commit 10d26045f857 found 6 issues with loop artifacts

---

## What Was Fixed

### Issue 1 — Reflections artifact corrupted (CRITICAL)

**Problem:** LLM meta-commentary was committed verbatim into `loop/reflections.md`, including phrases like "Ready to append to loop/reflections.md and update loop/state.md to iteration 270?" and "Shall I proceed with writing these updates?" appended after legitimate reflection entries.

**Fix:** Rewrote `loop/reflections.md` from line 2320 onwards. Removed all LLM dialog, duplicate reflection entries, empty template blocks, and code-fence artifacts. Retained the 5 canonical reflection entries covering iterations 262–270 (Lessons 64–67). File is now 2379 lines (was 2528).

**Files changed:** `hive/loop/reflections.md`

### Issue 4 — Lessons 65 and 66 not added to state.md

**Problem:** Lessons 65 and 66 were formalized in `reflections.md` but never added to `state.md`'s canonical lessons list. Scout reads state.md first — lessons invisible there are invisible to the loop.

**Fix:** Added lessons 64, 65, 66, and 67 to `state.md` after lesson 60. All four lessons from the recent escalation-binding reflection cluster are now in the canonical list.

**Files changed:** `hive/loop/state.md`

### Issue 2 — build.md was missing

**Problem:** The Builder phase had no artifact. CLAUDE.md requires every phase to write its file.

**Fix:** This file.

---

## Issues Not Fixed (with explanation)

### Issue 3 — Commit title contradicts diff

The title "Verify and patch Knowledge tab routing and templates" is now in git history and cannot be changed. The title was misleading — noted in the audit trail via reflections.

### Issue 5 — Iteration numbering gap (268 to 270)

The state.md iteration counter is in the past. State.md now reads "Iteration 271" from a subsequent update. The gap happened; it's in the append-only reflections log as a lesson learned. Cannot retroactively renumber.

### Issue 6 — No Knowledge tab verification occurred

The escalation from iteration 265 (Knowledge tab routing and template verification) remains unverified. DATABASE_URL is not set in the Builder environment — integration tests cannot run. This is Lesson 65 in action: escalations without matching infrastructure remain unverifiable. The site build passes and unit tests pass. The Knowledge tab was shipped in a prior iteration and the code is in place.

---

## Build Results

```
go.exe build -buildvcs=false ./...   -> OK (no errors)
go.exe test ./...                    -> OK (all pass: auth, graph packages)
```

---

## Files Changed

| File | Change |
|------|--------|
| `hive/loop/reflections.md` | Removed 149 lines of corrupted content (LLM dialogs, duplicates, empty blocks); now 2379 lines with 5 canonical reflection entries |
| `hive/loop/state.md` | Added lessons 64-67 to canonical lessons list |
| `hive/loop/build.md` | Created (this file) |
