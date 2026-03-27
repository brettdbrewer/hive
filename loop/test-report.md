# Test Report: Fix Architect Operate path causality (Invariant 2)

- **Build:** Fix: Architect Operate path does not thread causes — subtasks created via Operate() have no causes field
- **Commit:** 2abed27cb2f5d340ee227ac82108df3479f21ce4
- **Result:** PASS — all tests green

## What Was Tested

The fix added `milestoneID` as a 3rd parameter to `buildArchitectOperateInstruction`, injecting `,"causes":["<milestoneID>"]` into the curl template the LLM receives. Two production paths create tasks: the Operate path (IOperator) and the Reason/fallback path. Both were verified for causality. One edge case was missing coverage and was added.

## Tests Run

### Tests verified by Builder (confirmed passing)

| Test | File | Result |
|------|------|--------|
| `TestRunArchitectOperateInstructionIncludesCauses` | `pkg/runner/architect_test.go` | PASS |
| `TestRunArchitectSubtasksHaveCauses` | `pkg/runner/architect_test.go` | PASS |

**TestRunArchitectOperateInstructionIncludesCauses** — Core fix test. Sets up a milestone (`milestone-42`), uses `mockCaptureOperator` to intercept the instruction, asserts `"causes":["milestone-42"]` is present. Would have failed before the fix (milestoneID was not passed to `buildArchitectOperateInstruction`).

**TestRunArchitectSubtasksHaveCauses** — Reason/fallback path. Uses a real `mockProvider` that returns SUBTASK_ format, captures all HTTP POST bodies, asserts every `op=intend` request includes `"causes":["milestone-77"]`.

### New test (added by Tester)

| Test | File | Result |
|------|------|--------|
| `TestRunArchitectOperateInstructionNoCausesWhenNoMilestone` | `pkg/runner/architect_test.go` | PASS |

**TestRunArchitectOperateInstructionNoCausesWhenNoMilestone** — Edge case: Operate path triggered by scout-report fallback (no milestone on board). Verifies the curl template does NOT include `"causes"` at all — an empty `milestoneID` must produce empty `causesSuffix`, not `"causes":[""]`. This ensures Invariant 2 is respected symmetrically: causes are declared only when a causal parent exists.

### Full suite

```
ok  github.com/lovyou-ai/hive/cmd/mcp-graph       (cached)
ok  github.com/lovyou-ai/hive/cmd/mcp-knowledge   (cached)
ok  github.com/lovyou-ai/hive/cmd/post             (cached)
ok  github.com/lovyou-ai/hive/pkg/api              (cached)
ok  github.com/lovyou-ai/hive/pkg/authority        (cached)
ok  github.com/lovyou-ai/hive/pkg/hive             (cached)
ok  github.com/lovyou-ai/hive/pkg/loop             (cached)
ok  github.com/lovyou-ai/hive/pkg/resources        (cached)
ok  github.com/lovyou-ai/hive/pkg/runner           3.578s
ok  github.com/lovyou-ai/hive/pkg/workspace        (cached)
```

No regressions.

## Edge Cases Considered

- **Milestone present (Operate path)** — causes suffix injected into curl template. ✓ Tested.
- **No milestone (Operate path, scout-report fallback)** — no causes key in template, not even `"causes":[""]`. ✓ Tested (new).
- **Milestone present (Reason path)** — causes passed to `CreateTask`. ✓ Tested.
- **No milestone (Reason path)** — `causes` is nil, `CreateTask` called without it. Covered implicitly by existing parse tests.

## Coverage Notes

- `buildArchitectOperateInstruction`: both branches of `if milestoneID != ""` are now covered.
- `runArchitect` Operate path: both milestone and no-milestone cases covered.
- No untested code paths introduced.
