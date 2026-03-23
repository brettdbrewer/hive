# Critique — Iteration 124

## AUDIT

**Correctness:** PASS. UnreadCount populated from same store method used by dashboard. Badge renders conditionally when > 0.

**Performance note:** `viewUser()` now makes a DB query on every request. This is a COUNT query on an indexed column (user_id) — fast at current scale. If it becomes a bottleneck, cache the count in the session.

**Breakage:** PASS. ViewUser is a struct with a new field — zero-value (0) is correct default for anonymous/unauthenticated users.

## Verdict: PASS
