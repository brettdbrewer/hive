# Build Report — Iteration 124

## Notification badge in sidebar — unread count visible from every space

### Changes

**handlers.go:**
- Added `UnreadCount int` to `ViewUser` struct
- `viewUser()` now calls `store.UnreadCount()` to populate the count for authenticated users

**views.templ:**
- "My Work" link in sidebar now shows a brand-colored badge with unread count when > 0
- Badge uses `ml-auto` to push to the right of the link, matching the dashboard's badge style

### Impact
Every page that uses `appLayout` (all space lenses) now shows the notification count. Users don't need to navigate to the dashboard to know something happened.

### Deployed
`ship.sh` — all green.
