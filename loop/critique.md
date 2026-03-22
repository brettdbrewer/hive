# Critique — Iteration 36

## Verdict: APPROVED

## Trace

1. Scout identified: People and Activity lenses don't distinguish agents from humans
2. Builder added `ActorKind` to Op via JOIN against users table (no schema migration)
3. People and Activity templates updated with violet avatar + badge (same pattern as Feed/Chat)
4. Deployed and healthy

## Audit

**Correctness:**
- LEFT JOIN handles ops from actors with no user record (COALESCE to 'human'). ✓
- Agent detection resolved from users table, not from op content. Lesson 30 applied. ✓
- Visual treatment identical across all six lenses: Feed, Chat, Comments, People, Activity, Board. ✓
- Member kind populated from first op seen (stable — all ops from same actor have same kind). ✓

**Gaps:**
- **JOIN performance at scale**: Every ListOps and ListNodeOps query now JOINs against users. Fine for current traffic. At scale, consider denormalizing `actor_kind` onto the ops table or caching.
- **Board lens doesn't show agent badges**: Task cards show author but not author_kind badges. Visually consistent since tasks are usually human-authored, but worth noting.
- **Name collision**: JOIN is `users.name = ops.actor`. If two users share a name (unlikely but possible), the JOIN is ambiguous. A future iteration should use actor IDs throughout.

## DUAL

Clean iteration. No mid-build feedback needed. The JOIN pattern avoided a large refactor (changing RecordOp signature + 10 call sites + schema migration) in favor of resolving from the identity system at query time. The right trade-off for now — performance optimization can come later if needed.
