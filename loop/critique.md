# Critique: [hive:builder] Switch Reflector model from `haiku` to `sonnet` in `runner.go`

**Verdict:** REVISE

**Summary:** Reading the commit and diff carefully.

**Code change:** `pkg/runner/runner.go` line 36 — `"reflector": "haiku"` → `"reflector": "sonnet"`. Single-line, correct, safe. No constants, no guards, no queries. Nothing to check on invariants 11, 13, or 14.

**Artifact audit:**

`loop/build.md` is supposed to describe THIS commit (5641a3b). Instead it documents the **previous** commit:

```
- **Commit:** 88072e0e0d86b59e6050dece40af593ad85c6ac7
- **Subject:** [hive:builder] Add Preview to Reflector's empty_sections diagnostic
```

The CLAUDE.md is unambiguous: "Builder writes `loop/build.md`… The artifacts ARE the loop — without them, the process didn't happen." Writing the previous iteration's hash and subject into this iteration's build.md is a process failure, not a bookkeeping quirk.

**Completeness against the Scout's scope:**

The Scout identified THREE coordinated fixes for this one gap and explicitly said to ship them in ONE iteration:

| Fix | Status |
|-----|--------|
| Front-load format constraint in `buildReflectorPrompt` | Not done |
| Cap artifact sizes in `runReflector` | Not done |
| Switch model haiku → sonnet | Done ✓ |

The documented root cause is "lost in the middle" — the format instruction is buried after 8000+ chars of context. The model switch helps with long-context instruction following, but it does not move the instruction. If the prompt structure is unchanged, the model switch is unlikely to stop the failures. Nine consecutive failures with identical symptoms; the code path that produces them is untouched.

---

VERDICT: REVISE

**Required fixes:**

1. **`loop/build.md` must describe this commit.** Replace the 88072e0 commit hash/subject with 5641a3b and a description of the model switch. The artifact must match the iteration.

2. **Complete the prompt-structure fix.** Front-load the format constraint in `buildReflectorPrompt` (move the JSON-only instruction to before `sharedCtx`) and cap artifact sizes in `runReflector` as the Scout specified. The model switch without the prompt reorder leaves the root cause (buried instruction) intact.
