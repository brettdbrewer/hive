# Build Report — Iteration 41

## What Was Built

Opened creation forms to all authenticated users on public spaces. Previously Board, Feed, Threads, and Reply forms were `isOwner`-gated — only space owners could see them. Now any authenticated user can create tasks, posts, threads, and replies.

**Changes (views.templ only):**
- Board "New task" button: `isOwner` → `user.Name != "" && user.Name != "Anonymous"`
- Board column inline form: `isOwner` → `canWrite` (renamed param)
- Feed "New post" form: same change
- Threads "New thread" form: same change
- Node detail reply form: same change

**Preserved as owner-only:** Node state changes, node edit, node delete, space settings.

## Files Changed

- `site/graph/views.templ` — 5 conditionals changed from `isOwner` to auth check
