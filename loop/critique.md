# Critique: [hive:builder] Reflector assigns duplicate lesson numbers: 109 (�3), 136 (�2), 137 (�2) found in claims

**Verdict:** PASS

**Summary:** ## Critic Review

**Derivation chain**: gap (duplicate lesson numbers on overlap/retry) → query claims for max → use max+1.

### Code correctness

**`parseLessonNumber`**: `fmt.Sscanf(rest, "%d", &n)` correctly stops at `:` — `"109: highest"` → 109. `"Lesson: 2026-03-29"` → `rest = ": 2026-03-29"` → Sscanf fails → n stays 0. Edge cases handled.

**`NextLessonNumber`**: Returns 1 on API error (safe default). `GetClaims(slug, 200)` is bounded (Invariant 13). At lesson 182 currently, 200 is sufficient headroom.

**Format string in `runReflectorOperate`** (line 296–301): 8 format verbs (`%d %s %s %d %s %s %s %s`), 8 args (`nextLessonNum, apiKey, SpaceSlug, nextLessonNum, causesSuffix, apiKey, SpaceSlug, causesSuffix`). Types match. ✓

**`runReflectorReason`**: Title becomes `"Lesson 110: 2026-03-29"` — number from graph, date from `time.Now()`. Causality chain preserved in `AssertClaim` call (Invariant 2). ✓

### Tests

- `TestParseLessonNumber` — 7 cases including the date-based trap case. ✓
- `TestNextLessonNumberFromClaims` — lessons 45+109 present → expects 110. ✓
- `TestNextLessonNumberNoClaims` — expects 1. ✓
- `TestNextLessonNumberAPIError` — expects 1. ✓
- `TestRunReflectorReasonLessonNumberFromGraph` — end-to-end: mock claims return max 109, asserted title must start `"Lesson 110:"`. Mock routing on `tab=claims` matches `GetClaims` URL `?tab=claims&limit=200`. Mock JSON keys (`cover`, `blind`, etc.) match `jsonReflectorOutput` struct tags. ✓

### Invariants

- **Invariant 11**: Lesson number is an ordinal, not an entity ID — no identity violation.
- **Invariant 12**: All new code paths tested.
- **Invariant 13**: `GetClaims(slug, 200)` — bounded.
- **Invariant 2**: Causes chain preserved through `AssertClaim`.

No bugs, no missing tests, no invariant violations.

VERDICT: PASS
