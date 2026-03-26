# Critique

Commit: b705a3d92e3b8514a9748de3e04d6b1f2515f7ce
Verdict: REVISE

Now I have a complete picture. Let me summarize my findings.

---

## Analysis

### What the commit does
Adds `PhaseEvent` type and `appendDiagnostic` method in `diagnostic.go`, then calls it from `workTask` at two points:
1. `runner.go:282` — Operate error
2. `runner.go:303` — `verifyBuild()` failure after DONE

### Issue 1: Missing test for the build verification failure path (Invariant 12 — VERIFIED)

The commit title explicitly claims to cover **both** "Operate error" and "build verification failure". The test `TestWorkTaskOperateErrorWritesDiagnostic` covers path 1 (Operate error). **There is no test for path 2** (build verification failure at `runner.go:303`).

Path 2 requires:
- A mock operator that returns a DONE response (not an error)
- A `verifyBuild()` call that fails (e.g., `go.exe` not found, or the workdir has no Go files)

The test file ends at line 310 — no such test exists. The invariant is unambiguous: "If the Critic can't point to a test that covers the change, REVISE."

### Issue 2: state.md corruption

`loop/state.md` — the "What the Scout Should Focus On Next" section was replaced with `` `/` `` (a single backtick and forward slash). The section heading remains, the content is gone. This looks like a truncated write from a prior agent. Not a code bug, but it's a corrupted loop artifact that will mislead future iterations.

### What's correct

- `appendDiagnostic` is tested end-to-end via the Operate error test (file is written, content is verified) ✓
- `PhaseEvent` fields are sensible (`cost_usd` = cumulative cost at point of failure, which is correct for both error paths) ✓  
- File open/write is append-with-create, correct for a `.jsonl` sink ✓
- Silent skip when `HiveDir == ""` is appropriate for test isolation ✓
- `Timestamp` is set server-side, not caller-provided — correct ✓

---

VERDICT: REVISE

**Required fix:** Add `TestWorkTaskBuildFailureWritesDiagnostic` — a test where the mock operator returns a `ACTION: DONE` response successfully (no Operate error), but `verifyBuild()` fails because the repo has no buildable Go code. Verify that `diagnostics.jsonl` is written with `phase=builder` and an error message. The pattern is already established in `TestWorkTaskOperateErrorWritesDiagnostic` — follow the same shape.
