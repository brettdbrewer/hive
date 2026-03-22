# Build Report — Iteration 21

## What I planned

Add API key authentication so agents can programmatically interact with lovyou.ai. This is the foundation for "humans and agents, building together."

## What I built

Changes to 1 file (auth/auth.go) in the site repo.

### API key data model
- `api_keys` table: `id`, `name`, `key_hash` (SHA-256, unique), `user_id` (FK → users), `created_at`
- `APIKey` struct for metadata (no raw key in struct — only returned at creation)

### Key generation
- `createAPIKey(ctx, userID, name)` → returns raw key `lv_` + 64 hex chars
- SHA-256 hash stored, raw key never persisted
- `lv_` prefix makes keys identifiable in logs/configs

### Authentication flow
- `userFromBearer(r)` — extracts `Authorization: Bearer <key>` header, hashes key, looks up user
- `RequireAuth` now checks Bearer token first, falls back to session cookie
- `OptionalAuth` same — Bearer token → cookie → anonymous
- Seamless: existing browser users unaffected, agents use Bearer token

### API key management routes
- `POST /auth/api-keys` — create key (requires session auth), returns JSON `{"key": "lv_...", "name": "..."}`
- `POST /auth/api-keys/{id}/delete` — revoke key (owner only)
- `ListAPIKeys(ctx, userID)` — public method for UI integration (returns metadata only)

### Security
- Raw key shown exactly once (at creation), never stored
- SHA-256 hash comparison for lookups
- Delete requires both key ID and matching user_id (can't delete others' keys)
- API key creation requires session auth (must be logged in via browser to generate keys)

## Verification

- `go build -o /tmp/site.exe ./cmd/site/` — success
- Committed and pushed to main
- Deployed to Fly.io — both machines healthy
