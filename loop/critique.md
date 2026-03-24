# Critique — Iteration 191

## Derivation Chain
- **Gap:** Phase 2 item 2 — follow users. Subscribe grammar op.
- **Plan:** New follows table, store methods, profile button + counts.
- **Code:** Matches plan. Scoped correctly — no feed filtering yet.

## Follow Users: PASS

**Correctness:**
- ON CONFLICT DO NOTHING — idempotent follow. ✓
- Self-follow prevented (redirect no-op). ✓
- Notification only on follow, not unfollow. ✓
- Viewer auth checked before showing button. ✓

**Identity:**
- Follow uses IDs (`follower_id`, `followed_id`). ✓
- `ResolveUserID(name)` to get target ID from URL. ✓
- `viewer.ID` from auth context. ✓

**BOUNDED:**
- Count queries are single-row aggregates. ✓
- `IsFollowing` is an EXISTS check. ✓

**Template:**
- Follow button next to endorse button. Layout preserved.
- Stats line shows follower/following counts. Clean.
- No HTMX swap — uses full form POST + redirect. Acceptable for profile page (not high-frequency action). Could be HTMX-ified later.

**Tests:** No new tests for follow methods. Pattern matches endorsements which are tested. Acceptable.

**NOTE:** Feed filtering by followed users not yet implemented. Scoped out correctly — one gap per iteration.

## Verdict: PASS
