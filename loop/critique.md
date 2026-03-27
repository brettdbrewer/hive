# Critique: [hive:builder] Fix: Architect Operate path does not thread causes � subtasks created via Operate() have no causes field

**Verdict:** PASS

**Summary:** The implementation and tests are correct. Let me do a final check on the instruction template to confirm the JSON is well-formed when `causesSuffix` is injected.

The template at line 236:
```
'{"op":"intend","kind":"task","title":"<TITLE>","description":"<DESCRIPTION>","priority":"high"%s}'
```

With `causesSuffix = ,"causes":["milestone-42"]` this produces:
```json
{"op":"intend","kind":"task","title":"...","description":"...","priority":"high","causes":["milestone-42"]}
```

Valid JSON. ✓

Six `%s` in the `fmt.Sprintf` at line 247: `context`, `apiKey`, `spaceSlug`, `causesSuffix`, `apiKey`, `spaceSlug` — count matches. ✓

Test coverage:
- `TestRunArchitectOperateInstructionIncludesCauses` — Operate path with milestone, asserts `"causes":["milestone-42"]` present ✓
- `TestRunArchitectSubtasksHaveCauses` — Reason/fallback path, asserts causes field on each subtask ✓  
- `TestRunArchitectOperateInstructionNoCausesWhenNoMilestone` — Operate path with no milestone, asserts no `"causes"` key (bonus test not mentioned in build.md but present) ✓

Invariant 2 (CAUSALITY): satisfied — both Operate and Reason paths now wire causes. ✓  
Invariant 11 (IDENTITY): `milestone.ID` used, not `milestone.Title`. ✓  
Invariant 12 (VERIFIED): all paths have direct test coverage. ✓

The fix is precise, the tests are complete, and there's an additional test (`NoCausesWhenNoMilestone`) that wasn't required but correctly validates the empty-milestone guard. The build.md and critique.md accurately reflect the state of the code.

VERDICT: PASS
