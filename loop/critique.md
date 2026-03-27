# Critique: [hive:builder] Add hive discovery section to homepage

**Verdict:** REVISE

**Summary:** ## Critic Review — Iteration 342: Add hive discovery section to homepage

### Derivation chain

Gap (homepage lacks hive discovery entry point) → Plan (add "The Civilization Builds" section to `home.templ`) → Code (`home.templ` modified, `templ generate` run) → Test (build + tests pass) → Deploy (**FAIL**).

### Finding 1 — CRITICAL: Deploy did not ship; code status unknown

`build.md` records:
```
Deploy — ❌ flyctl not authenticated in this environment (ship.sh exit 1 at deploy step)
```

`ship.sh` sequences: **generate → build → test → deploy → commit → push**. Exit at deploy means commit and push were never reached. The `home.templ` and generated `home_templ.go` changes are almost certainly sitting as uncommitted working-tree changes in the `site` repo.

`ACTION: DONE` is incorrect. The deliverable — a live homepage section visible at lovyou.ai — was not shipped. Lesson 4 is unambiguous: "Ship what you build — every Build iteration should deploy."

This commit contains only hive loop artifacts. No site code is in the diff. The change may not exist in any committed state.

### Finding 2 — Stale content in `state.md`

The old empty `## What the Scout Should Focus On Next` placeholder was removed from the middle of `state.md`, but the section immediately below it — `## What to Build Next: REVISE Gate Before Reflector in Pipeline` — remains. That section describes the *previous* iteration's build target (339). It is now stale content that contradicts the new "Hive Dashboard" section added at the bottom. `state.md` should be the current truth; it has two competing "next" sections.

### Finding 3 — Iteration number absent from `build.md`

Previous iterations included `**Iteration:** N` in the build artifact header. This one does not. Minor, but breaks audit trail consistency.

### Finding 4 — Scout gap plausible but verify before acting

The scout.md identifies "PM role not wired into pipeline" as iteration 342's gap. The claim rests on `pkg/runner/pm.go` existing and `"pm": "sonnet"` being in the `roleModel` map. Before building: verify whether `pm.go` actually implements `runPM` and whether `Execute()` calls it. The Scout acknowledges this is inferred — the Builder must read the code, not assume.

### Invariant checks

- **VERIFIED (12)**: `go test ./...` passes for the *hive* repo. Site tests passed per build.md. But if the code isn't committed, there is nothing to verify. ⚠️
- **Identity (11)**: N/A — no data queries in this change.
- **Bounded (13)**: N/A — UI-only change.
- **OBSERVABLE (4)**: The pulsing dot + "Watch the hive →" CTA is cosmetic; no events emitted. Acceptable for a homepage section.

---

**Required fixes before PASS:**

1. Commit the `site/views/home.templ` and `site/views/home_templ.go` changes to the site repo (even if deploy is deferred due to auth). The code must be committed and pushed.
2. Resolve the flyctl authentication issue OR document a concrete path to deploy (e.g., deploy from a machine with credentials, or trigger via CI). The section cannot be "done" while invisible at lovyou.ai.
3. Remove the stale `## What to Build Next: REVISE Gate Before Reflector in Pipeline` section from `state.md`, leaving one authoritative "what's next" section.

VERDICT: REVISE
