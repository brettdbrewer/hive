# Critique: [hive:builder] Switch Reflector to JSON output with text fallback

**Verdict:** PASS

**Summary:** ## Critique: Switch Reflector to JSON output with text fallback

**Derivation chain:** gap (REVISE — JSON parser inert because `buildReflectorPrompt` still requested text markers, AND fence-wrapped text-marker responses invisible to text-marker parser) → plan (normalize fences before both parse paths, mirror Architect pattern) → code (`normalizeReflectorResponse` added, `parseReflectorJSON` renamed, `parseReflectorOutput` calls normalize first) → tests (existing reflector test suite passes all variants).

---

### Code correctness

**`normalizeReflectorResponse`** — fence-stripping is correct. `TrimSpace` at the top means leading-space fence variants (` ``` `) are handled after trim. Closing-fence check works because content is trimmed before the suffix check. Single-line ` ```json ` with no newline correctly skips the opening-fence strip (no `\n` → `nl < 0` → no strip), which is harmless since there's no valid content anyway.

**`parseReflectorJSON`** — correctly delegates fence-stripping to the caller. The prose-preamble scan (`for i, ch := range content`) is bounded by content length. Flat object, `{"reflection":{...}}` wrapper, and preamble cases all handled.

**`parseReflectorOutput`** — normalize → JSON → text-marker. The key fix: fence-wrapped text-marker responses now reach the text-marker parser. This was the `empty_sections` root cause path that the prior fix missed. The pattern now mirrors the Architect fix correctly.

**Prior minor note stands:** `strings.HasSuffix(content, "```")` won't catch ` ``` ` with trailing spaces — but `TrimSpace` at the top of `normalizeReflectorResponse` handles this.

**Tests:** JSON flat object, wrapper, preamble, all text-marker variants covered. Build passes.

---

### Artifact issues (non-blocking)

**reflections.md** — the new entry is valid and well-formed (COVER/BLIND/ZOOM/FORMALIZE + Lesson 87), but it ends with a "Pending file updates" block that is LLM deliberation leaked into the artifact. This is the same failure mode that produced Lesson 86 and now Lesson 87. The loop is learning the lesson but still triggering it. Lesson 87 exists in reflections.md — the lesson is captured even if messily.

**Lesson 87 absent from state.md** — the reflections.md explicitly flags "Add Lesson 87 to state.md" as pending, but the state.md diff shows only the iteration counter was bumped. The Reflector's own lesson (close after verdict, finalize don't pose questions) wasn't codified into the lessons section it was written to protect. This is the Reflector failing to complete its own artifact pass.

The iteration counter is correct (330→331). The code fix is complete. The lesson omission from state.md is real but doesn't affect the production change.

---

### Invariants

- VERIFIED ✓ — test coverage exists for the change
- BOUNDED ✓ — scan loop bounded by content length
- IDENTITY ✓ — not applicable
- No invariant violations in the production code

---

VERDICT: PASS

Artifacts: clean up reflections.md (remove the "Pending file updates" block) and add Lesson 87 to state.md in the next artifact pass. Neither requires a build cycle — the code fix is correct and complete.
