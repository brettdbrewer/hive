# Build: Auth: helpful error messages and logging

## Scout Gap Cross-reference

Scout 406 gap: `assertClaim` wrapper in `cmd/post` (CAUSALITY GATE 1). This build addresses a different high-priority task: auth failure visibility. No conflict — Scout gap remains open for next iteration.

## What Was Built

### `site/auth/auth.go`

**New imports:** `html`, `net/url` for error page rendering and URL encoding.

**New routes registered:**
- `GET /auth/error` — user-facing error page
- `GET /auth/status` — debug endpoint (no secrets)

**`handleLogin`:** Log now includes `state` prefix (first 8 chars) and `host` for correlation.

**`handleCallback`:** Comprehensive overhaul:
1. Checks `?error=` from Google before state validation — redirects to `/auth/error?code=<errCode>` (e.g. `access_denied`)
2. State mismatch/missing cookie → redirect to `/auth/error?code=invalid_state` (was: `http.Error` 400)
3. Token exchange failure → redirect to `/auth/error?code=exchange_failed` + logs error
4. Userinfo fetch/decode failure → redirect to `/auth/error?code=userinfo_failed` + logs error
5. Upsert user failure → redirect to `/auth/error?code=user_create_failed`, logs email
6. Session create failure → redirect to `/auth/error?code=session_failed`
7. Success path: logs `auth: user upserted id=... email=...` and `auth: login success user=... email=... session=...`

**`handleAuthError`:** Styled HTML error page with per-code messages:
- `access_denied` → "Sign-in was cancelled or blocked. Your organisation may restrict third-party sign-in."
- `invalid_state` → "Your sign-in session expired. Please try signing in again."
- `exchange_failed` → "Could not complete sign-in with Google... Try using an API key instead."
- default → "Authentication failed. Please try again."
- Shows error code and "Try again" link back to `/auth/login`.

**`handleStatus`:** Returns JSON with `oauth_configured`, `redirect_url`, `secure`. No client secret exposed.

### `site/auth/auth_test.go`

New tests (no DB required):
- `TestCallbackInvalidState` — state mismatch → redirect to `/auth/error?code=invalid_state`
- `TestCallbackExpiredState` — missing cookie → redirect to `/auth/error?code=invalid_state`
- `TestCallbackGoogleError` — `?error=access_denied` → redirect to `/auth/error?code=access_denied`
- `TestAuthErrorPage` — 4 subtests: access_denied message, invalid_state message, no-code default, retry link
- `TestAuthStatus` — oauth_configured=true, redirect_url contains host, no secrets in body

DB-required test:
- `TestConcurrentSessions` — two active sessions for same user both resolve correctly (skips without DATABASE_URL)

## Build Results

```
site: go.exe build -buildvcs=false ./...   → OK
site: go.exe test ./...                    → auth OK, graph OK, handlers OK
hive: go.exe build -buildvcs=false ./...   → OK
hive: go.exe test ./...                    → all 11 packages pass
```

## Files Changed

- `site/auth/auth.go` — logging, error redirects, `handleAuthError`, `handleStatus`, new routes
- `site/auth/auth_test.go` — 6 new test functions (5 no-DB, 1 DB-required)
