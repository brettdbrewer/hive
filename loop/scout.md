Perfect. Now I have the full context. Let me write the Scout's gap report.

---

# GAP REPORT — Iteration 336

**Gap:** The Builder's artifact discipline during REVISE cycles is creating systematic audit trail corruption — `build.md` consistently describes stale commits instead of the actual fix commits, three iterations in a row.

**Evidence:**

1. **Critic explicitly documented the pattern** (Iteration 335 critique):
   - "This is the third consecutive iteration where `build.md` describes the wrong commit"
   - Iteration 335's `build.md` claims to describe commit `035dc32` but embeds diff stats from file `reflector_test.go`, which doesn't exist in that commit
   - The file actually describes commit `6279036` (the real work), not the file's stated commit

2. **Code inspection** — `pkg/runner/builder.go`:
   - When the Critic issues REVISE, the Builder applies fixes and commits new work
   - The Builder likely regenerates `build.md` from the original task context, not from the new commit
   - This causes `build.md` to retain the original (failed) commit hash while the diff stat includes the new (fixed) code

3. **Systemic impact** — CLAUDE.md Section "The artifacts ARE the loop":
   - "Every phase MUST write its artifact file...without them, the process didn't happen"
   - The audit trail is the mechanism for the Reflector to review what shipped (Lesson 68: "absence of feedback infrastructure is a system defect")
   - Lesson 43: "NEVER skip artifact writes" — this isn't skipping, it's writing false content, which is worse

4. **Related lessons flagged as critical**:
   - Lesson 68: "Feedback loop infrastructure is a critical path blocker...The loop depends on measurement to reflect on itself"
   - Lesson 70: "Loop artifact validation must check content completeness, not just file existence...Corrupted artifacts are worse than missing ones"

**Impact:**

- **Audit trail becomes untraceable** — future iterations can't determine what actually shipped or why. The Reflector's COVER section reads a corrupted history.
- **Loop can't self-correct** — The system measures success via artifacts. Corrupted measurement means the loop is blind to its own operation.
- **Violates VERIFIED invariant** — No visibility into which code changes correspond to which commits. Test coverage becomes unmappable to actual work.
- **Degrades post tool output** — `cmd/post` reads `build.md` to publish iteration summaries; corrupted metadata spreads to the public feed.

**Scope:**

| File | Issue | Lines |
|------|-------|-------|
| `pkg/runner/builder.go` | REVISE flow doesn't update `build.md` commit hash/diff to reflect the new fix commit | ~line 200-250 (estimate) |
| `pkg/hive/agentdef.go` | Builder agent prompt during REVISE — may not include instructions to regenerate artifact with new commit data | SystemPrompt builder section |
| `pkg/runner/runner.go` | REVISE verdict handling — needs to pass new commit context to Builder's artifact regeneration | ~line 180+ |
| Test coverage | No test verifies `build.md` content matches the commit it claims to describe | new test needed |

**Suggestion:**

Three coordinated fixes:

1. **Fix Builder REVISE flow**: When Critic issues REVISE, Builder fixes code → commits → **MUST update `build.md` with new commit hash, subject, timestamp, and regenerated diff stat**. Don't preserve the original commit reference.

2. **Add validation gate**: Before Builder considers REVISE complete, verify `build.md` contains the new commit hash (not the REVISE-triggering commit). Add this check to the Builder's internal post-fix validation.

3. **Add test coverage**: In `pkg/runner/builder_test.go` or `reflector_test.go`, add a regression test: simulate a REVISE cycle (fixture with a failed build → fix → re-commit), verify the resulting `build.md` describes the fix commit, not the original.

This is a structural issue in the Builder agent itself (lives in the hive repo), blocks loop integrity (Lesson 68 critical-path), and is explicitly flagged by the Critic for next-iteration fixing.