# Critique: [hive:builder] populateFormFromJSON silently drops causes array � CAUSALITY broken for all LLM-driven ops

**Verdict:** PASS

**Summary:** The code is in the file. Let me verify the fix is complete and correct.

**`populateFormFromJSON` (lines 520–554):** Correctly changed from `map[string]string` to `map[string]any` with a type switch. The `[]interface{}` arm joins string items as CSV — exactly what the server-side `strings.Split(causesStr, ",")` parser expects.

**`TestAssertOpMultipleCauses` (lines 244–339):** Posts `"causes":[cause1.ID, cause2.ID]` as a JSON array, then verifies:
1. The POST response returns both cause IDs (asserts `len == 2`, checks each ID in a set)
2. The subsequent GET `/knowledge?tab=claims` also returns both causes

Both verifications use `.ID` fields — Invariant 11 satisfied.

**Checklist:**
- Gap matches fix: JSON array causes were silently dropped → now joined as CSV → CAUSALITY invariant unblocked ✓
- No existing behavior broken: string fields handled identically ✓  
- Invariant 11 (IDs not names): test uses real node IDs throughout ✓
- Invariant 12 (VERIFIED): new test directly exercises the array path, 9/9 knowledge tests pass ✓
- Invariant 2 (CAUSALITY): structurally enforced again ✓
- Build clean ✓

One edge case worth noting: non-string items in a JSON array are silently dropped (e.g., `"causes":[1,2]`). Acceptable — cause IDs are always strings in this system.

VERDICT: PASS
