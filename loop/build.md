# Build: Auth: ship email magic link for Workspace-blocked users

- **Commit:** 2dcb0264f37bdb4249b132b1f00f4760d6cf6734
- **Subject:** [hive:builder] Auth: ship email magic link for Workspace-blocked users
- **Cost:** $0.6489
- **Timestamp:** 2026-03-29T14:53:18Z

## Task

Corporate Google Workspace users (e.g. transpara.com) cannot log in because their admin blocks OAuth app approvals. Email magic link is the fallback. This is user-blocking.

**Target repo:** site

**TASK 1 � Magic link request handler**
Add `POST /auth/magic-link` in `site/handlers/auth.go` (or ne...

## What Was Built

Build and tests pass. Here's what was done:

**Already implemented (prior iteration):**
- `magic_link_tokens` DB table in `migrate()`
- `POST /auth/magic-link/request` — generates token, stores hash, logs link
- `GET /auth/magic-link/verify` — validates token, creates session, redirects to `/app`
- Full test suite (happy path, expired, used, invalid, idempotent)

**Added this iteration (`auth/auth.go`):**
- `GET /auth/google` — new route for the OAuth redirect (renamed from `handleLogin` to `handleGoogleOAuth`)
- `GET /auth/login` — new login page with Google button + collapsible "Use email instead" section; email form posts to `/auth/magic-link/request` with a JS-enhanced inline confirmation ("Check your email") on success

ACTION: DONE

## Diff Stat

```
commit 2dcb0264f37bdb4249b132b1f00f4760d6cf6734
Author: Matt Searles <matt.s.searles@gmail.com>
Date:   Mon Mar 30 01:53:18 2026 +1100

    [hive:builder] Auth: ship email magic link for Workspace-blocked users

 auth/auth.go  | 93 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++-
 loop/build.md | 40 +++++++++++++++++--------
 2 files changed, 120 insertions(+), 13 deletions(-)
```
