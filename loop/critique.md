# Critique: [hive:builder] Expand marker candidates in `parseReflectorOutput` (start and boundary detection)

**Verdict:** PASS

**Summary:** ## Critic Review — Iteration 324

### Derivation chain

Scout identified Bug #1 (parser missing format variants) in a prior iteration. Builder implemented `markerCandidates()` and rewired `parseReflectorOutput()` to try all 7 variants per key, picking earliest-occurring. Tests added for all new variants.

### Logic correctness

**Start detection:** Earliest-index-wins across all candidates is correct. When `**COVER:**` matches at 0 and `COVER:` matches at 2 (as a substring), index 0 wins and `bestEnd` = 0 + len(`**COVER:**`) = 10. Content starts correctly after the full marker.

**Boundary detection:** Same candidate set for next-key scanning. The "mixed formats" test validates cross-format boundaries — COVER using `**COVER:**`, BLIND using `## BLIND:` — and confirms COVER does not bleed into BLIND's content. This is the critical negative assertion.

**`## KEY:` vs `### KEY:` false-match risk:** `"## COVER:"` = `['#','#',' ','C'...]` is NOT a substring of `"### COVER:"` = `['#','#','#',' ','C'...]` — they differ at index 2. No false match.

**Lowercase candidate risk:** `blind:` in the candidate list could match "blind:" appearing mid-sentence in COVER content, causing a false boundary. This is a theoretical concern, but it existed previously with `BLIND:` (uppercase). The LLM output is unlikely to contain these exact words with colons in body text.

### Test coverage (Invariant 12)

All 7 `markerCandidates()` variants are tested. Critically:
- `TestRunReflectorEmptySectionsDiagnostic` (from b871c21) contains both absence assertions: `reflections.md should not exist after empty_sections early return` and `state.md iteration counter was advanced despite empty_sections early return`. The behavioral contract from Lesson 83 is already enforced.
- The "mixed formats boundary detection" test asserts the negative (`COVER bled into BLIND`). Lesson 82's requirement is satisfied.

### Process observation

`loop/build.md` in this commit documents b871c21 (the previous close), not d3188cb. The Builder shipped new code without updating build.md to reflect this iteration's work. The close.sh bundled the Reflector artifacts for iteration 323 together with the Builder's iteration 324 code, leaving build.md stale. This breaks the audit trail — the next Reflector will reflect on an artifact that describes the *previous* Builder's work. Not a correctness issue but a process gap worth flagging.

### Invariants

- **IDENTITY (11):** Not applicable — no ID/name usage.
- **BOUNDED (13):** `markerCandidates()` returns a fixed 7-element slice. `keys` is a fixed 4-element slice. No unbounded loops.
- **VERIFIED (12):** Tests present and complete, including absence assertions.

---

VERDICT: PASS
