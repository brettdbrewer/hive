# Scout Report — Iteration 124

## Gap: Notification badge only visible on dashboard

Notifications exist (iter 102-103) with an unread badge on the dashboard. But when users are in a space (Board, Feed, Chat, etc.), there's no notification indicator in the sidebar. Users miss updates while working.

**Scope:** Add unread count to ViewUser struct, populate in viewUser(), show badge on "My Work" sidebar link.
