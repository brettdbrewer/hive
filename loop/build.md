# Build Report — Iter 345 Fix: Critic Issues

**Date:** 2026-03-27
**Gap:** Three issues flagged by Critic on iter 345 (join_team/leave_team ops)

---

## Issues Fixed

### Issue 1 — Invariant 11: `user_name` removed from `node_members`

**Problem:** `user_name TEXT NOT NULL DEFAULT ''` was stored in `node_members` and captured at join time. Mutable display names captured at write time diverge from the user's current name.

**Fix:**
- Removed `user_name` column from `CREATE TABLE node_members` definition
- Added `ALTER TABLE node_members DROP COLUMN IF EXISTS user_name;` migration for existing databases
- Changed `JoinNodeMember(ctx, nodeID, userID, userName)` → `JoinNodeMember(ctx, nodeID, userID)` — no name stored
- Updated `ListTeamMembers` to `LEFT JOIN users u ON u.id = nm.user_id` and `COALESCE(u.name, nm.user_id)` — name resolved at query time from the authoritative users table
- Updated handler call in `handlers.go` to drop the `actor` (display name) argument
- Updated `store_test.go` to remove the `"Alice"` argument from `JoinNodeMember` calls

**Files changed:**
- `site/graph/store.go` — schema, `JoinNodeMember`, `ListTeamMembers`
- `site/graph/handlers.go` — `JoinNodeMember` call in `OpJoinTeam` case
- `site/graph/store_test.go` — `TestNodeMembership`

---

### Issue 2 — Duplicate heading in `state.md`

**Problem:** Lines 642–644 had two consecutive `## What the Scout Should Focus On Next` headings (one empty, one with stale Organize Mode content describing work already completed by iter 345).

**Fix:** Collapsed to a single heading. Updated the content to reflect the actual next focus: remaining Organize Mode tasks (assign_role/revoke_role, role badges, handler-level tests) rather than repeating the completed join_team/leave_team work.

**File changed:**
- `hive/loop/state.md`

---

### Issue 3 — Deploy

Deploy requires `flyctl` authentication not available in this session. Ops agent / human operator must run:

```bash
cd site && ./ship.sh "iter 345 fix: drop user_name from node_members, fix state.md duplicate heading"
```

---

## Verification

```
go.exe build -buildvcs=false ./...   → exit 0
go.exe test ./...                    → ok github.com/lovyou-ai/site/graph 0.108s
```

All tests pass including `TestNodeMembership`.
