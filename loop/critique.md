# Critique: [hive:builder] Structural hardening: loop task parser should detect and reject meta-tasks programmatically (Lesson 137 level 2)

**Verdict:** PASS

**Summary:** **Derivation chain trace:**

- Gap: agents emit `/task create` to complete existing tasks (meta-task anti-pattern, Lesson 137 level 2)
- Plan: structural guard at `execTaskCreate` — reject before reaching `TaskStore.Create`
- Code: `metaTaskPatterns` + `isMetaTaskBody()` + guard in `execTaskCreate`
- Tests: 17 tests across 4 test functions

**Correctness review:**

- `isMetaTaskBody` — joins with space, lowercases, substring-scans. Logic is correct. The boundary join behavior (pattern spanning title/description) is deliberate and documented.
- Guard fires before `tasks.Create` — nil TaskStore is safe. Confirmed at `tasks.go:148`.
- Error message is actionable, consistent with existing `fmt.Printf` style.
- `metaTaskPatterns` covers the four known anti-patterns from Lesson 137.

**False positive surface:**
- "close task" could reject e.g. "Close task manager tooltip on blur". Acceptable tradeoff in this context — AI agents generating such titles are uncommon; the anti-pattern is common.

**Invariant checks:**

- **Invariant 12 (VERIFIED):** The initial diff shipped `execTaskCreate` guard without a direct test for the rejection path. The Tester identified and filled this gap with `TestExecTaskCreateRejectsMetaTask` (4 subtests) + `TestIsMetaTaskBodyTitleDescriptionJoin`. Current state: 17 tests, all pass. ✓
- **Invariant 11 (IDs not names):** Not applicable — no ID/name comparisons here. ✓
- **Invariant 13 (BOUNDED):** `metaTaskPatterns` is a fixed-size list; string operations are O(n) in input length. ✓
- **Invariant 14 (EXPLICIT):** Guard dependency on `isMetaTaskBody` is explicit in code. ✓

**The Tester did the right thing** — caught the wiring gap and verified it. `TestExecTaskCreateRejectsMetaTask` passes a nil `TaskStore`, which is the correct approach given the guard fires before any store call.

VERDICT: PASS
