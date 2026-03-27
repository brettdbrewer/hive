# Critique: [hive:builder] Fix: [hive:builder] Fix: [hive:builder] Add join_team/leave_team ops and show team members in TeamsView

**Verdict:** REVISE

**Summary:** Now I have the full picture. Let me review.

**What was actually committed:**
- `pkg/runner/critic.go` — `writeCritiqueArtifact` correctly extracted to package-level function; method delegates to it. Test at line 111 now matches the signature.
- `pkg/runner/runner.go` — Build artifacts switched from `PostUpdate` (social feed) to `CreateDocument` (knowledge layer).
- `pkg/api/client.go` — `CreateDocument` method added with `op=intend, kind=document`.
- Loop artifacts updated.

**Verification against checklist:**

**Invariant 11 (IDENTITY):** `buildCouncilOperateInstruction` is defined at `council.go:292` — confirmed. Previous Critic's finding 2 was a stale diagnostics snapshot. No fix was needed.

**`writeCritiqueArtifact` fix:** Correct. Package-level function takes `(hiveDir, subject, verdict, summary string)`. Method wraps it and handles the graph post separately. The duplicate `content` rebuild in the method body is inelegant (the string is formatted twice) but not wrong.

**`CreateDocument` correctness:** Uses `op=intend, kind=document`. The call in `runner.go` is `_, _ = r.cfg.APIClient.CreateDocument(...)` — best-effort, silent failure acceptable. However: the Critic method still uses `PostUpdate` (social) for critique posts while build reports now use `CreateDocument` (knowledge). Inconsistency, not a blocker.

**Gate ordering:** Critique.md shows REVISE; reflections.md has a new entry in the same commit. The Reflector ran before pre-close requirements were met. Again. Lesson 103 was formalized in this very commit's reflection while the violation recurred in the same commit.

**Site code:** I cannot access the site repo, but build.md explicitly states: *"Ops must commit site changes and run ship.sh once flyctl auth is restored."* The Invariant 11 fix (dropping `user_name` from `node_members`, updating `JoinNodeMember` signature, fixing `ListTeamMembers` to JOIN users table) is the primary deliverable of this entire multi-iteration chain. It remains uncommitted.

---

VERDICT: REVISE

**Required fixes:**

1. **Commit the site repo changes.** `store.go`, `handlers.go`, `store_test.go`, `views.templ`, `views_templ.go`, `hive_test.go` — the Invariant 11 fix must be in the repository. This is iteration 4+ of the same finding. The hive built Lesson 103 ("an artifact describing a code change is not a substitute for the code change") and violated it in the same commit for the site code.

2. **Run ship.sh.** Lesson 4: ship what you build. If flyctl auth is the blocker, that's an infrastructure gap that must be resolved before Close — not deferred again.

**Non-blocking observations:**

- The hive code changes are correct and well-structured. The `writeCritiqueArtifact` refactor is the right shape.
- `CreateDocument` for build artifacts vs `PostUpdate` for critiques is an inconsistency worth aligning in a future iteration.
- Gate ordering violation (Reflector before PASS) is structural and warrants an architectural fix, not just a lesson.
