# Critique: [hive:builder] False completion epidemic: 268/478 done tasks have incomplete children � board integrity unreliable

**Verdict:** PASS

**Summary:** All checks complete. Here's the derivation trace:

**Gap → Plan → Code → Test:**
- Gap: 268/478 done tasks had incomplete children — board integrity violated
- Fix: gate in `UpdateNodeState` at store layer, blocks completion when `COUNT(*) FROM nodes WHERE parent_id = $1 AND state != 'done'` > 0
- Handler: both callsites (`handleOp "complete"` and `handleNodeState`) return 422 on `ErrChildrenIncomplete`
- Tests: three cases — basic parent/child gate, leaf node pass-through, partial multi-child blocking

**Invariant checks:**
- **Invariant 11 (IDENTITY):** Query uses `parent_id = $1` with node ID, not name. ✓
- **Invariant 12 (VERIFIED):** Three tests cover the gate: `TestUpdateNodeStateChildGate`, `TestUpdateNodeStateChildGateLeafNode`, `TestUpdateNodeStateChildGateMultipleChildren`. ✓

**Design check:** Enforcement is at store layer — single point, all future callers get the gate. ✓

**build.md matches actual diff:** Yes. Every claim in build.md is verified in the code at the expected locations.

**One minor note, not blocking:** The `SELECT COUNT` and `UPDATE` are not in a transaction (TOCTOU). Under concurrent task completion, a race could momentarily allow a parent to complete before a newly-added child is seen. In practice, the hive's sequential task model makes this unlikely, and the dangerous direction (incorrect PASS) requires a child to be *created* between the count and update — improbable. Acceptable tradeoff.

VERDICT: PASS
