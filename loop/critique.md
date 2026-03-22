# Critique — Iteration 21

## Verdict: APPROVED

## Trace

1. Scout identified that agents can't authenticate — entire vision blocked
2. Research explored hive agent capabilities and site API surface
3. Builder added api_keys table, SHA-256 hashing, Bearer token auth
4. Builder modified RequireAuth/OptionalAuth to check Bearer first
5. Builder added create/delete routes for key management
6. Built, pushed, deployed — both machines healthy

Sound chain. Minimal changes to existing code (auth middleware enhanced, not rewritten).

## Audit

**Correctness:** Key generation uses crypto/rand (same as session IDs). SHA-256 hash stored, not raw key. Bearer token checked before cookie — API clients never hit cookie logic. ✓

**Breakage:** RequireAuth/OptionalAuth are backward-compatible. If no Bearer header, behavior is identical to before. Existing session-based auth completely unaffected. ✓

**Security:**
- Raw key returned only once at creation, never stored ✓
- SHA-256 hashing prevents key recovery from database ✓
- Key creation requires session auth (can't bootstrap keys without logging in) ✓
- Delete validates user_id ownership ✓
- `lv_` prefix prevents accidental use of other tokens ✓

**Gaps (acceptable):**
- No rate limiting on API key usage. Fine for now — only agents and owner will use it.
- No key expiration. Keys live until deleted. Could add later.
- No UI for key management in settings page. Keys can only be created/deleted via API. Will add UI in a future iteration.
- No scoping (keys have full user permissions). Fine for owner-only usage.

## Observation

This is the most architecturally significant iteration since the unified graph product (iter 14). The previous 6 iterations (15-20) polished the site for human visitors. This one opens the door for non-human participants. The key insight is that Bearer token auth slots into the existing middleware with minimal changes — no new middleware chain, no separate auth path, just a `userFromBearer` check at the top of both wrappers.

The `lv_` prefix is a small but important detail: when an agent logs its API key in a config file, the prefix immediately identifies what it is.
