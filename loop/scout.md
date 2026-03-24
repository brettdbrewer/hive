# Scout Report — Iteration 191

## Gap: Follow users (Phase 2 item 2)

**Source:** social-spec.md Phase 2, board milestone. Maps to Subscribe grammar op.

**Current state:** No follow system exists. Profile pages show endorsements but no follow button. No follower/following counts. No "following" feed filter.

**What's needed:**
1. `follows` table — `follower_id, followed_id, created_at, PRIMARY KEY (follower_id, followed_id)`
2. Store methods: Follow, Unfollow, IsFollowing, CountFollowers, CountFollowing
3. Profile page: Follow/unfollow button (HTMX toggle), follower + following counts
4. Notification when someone follows you

**Scoping:** Feed filtering ("Following" tab on Feed) is a separate iteration. This iteration establishes the follow relation and makes it visible on profiles.

**Approach:** Follow the endorsement pattern — same table structure (from/to), same toggle handler, same HTMX swap. Profile page already has the endorsement button pattern to copy.

**From the spec:**
```
current_user.follows(user) → "Unfollow"
else → "Follow"
```

**Risk:** Low. New table, store methods, profile UI. Schema migration auto-applies via `CREATE TABLE IF NOT EXISTS`.
