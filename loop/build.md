# Build Report — Iteration 40

## What Was Built

Logged-in redirect: `/` → `/app` for authenticated users.

**site/cmd/site/main.go**:
- Home route moved from early registration to after auth setup
- Wrapped with `readWrap` (OptionalAuth) to detect session
- If `user != nil && user.ID != "anonymous"`, redirect 303 to `/app`
- Anonymous visitors still see the landing page
- No-DB fallback: home route registered without auth (always shows landing)

## Files Changed

- `site/cmd/site/main.go` — 12 lines (route move + conditional redirect)
