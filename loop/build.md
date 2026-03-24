# Build Report — Iteration 191

## Follow Users

**Schema:**
- `follows` table: `follower_id, followed_id, created_at, PRIMARY KEY (follower_id, followed_id)`
- Index on `followed_id` for follower count queries

**Store:**
- `Follow(followerID, followedID)` — ON CONFLICT DO NOTHING (idempotent)
- `Unfollow(followerID, followedID)` — DELETE
- `IsFollowing(followerID, followedID)` — EXISTS check
- `CountFollowers(userID)` — COUNT where followed_id = user
- `CountFollowing(userID)` — COUNT where follower_id = user

**Profile page:**
- `UserProfile` struct: added `Followers int`, `Following int`, `IsFollowing bool`
- Follow/unfollow button — form POST to `/user/{name}/follow`, redirects back
- Stats line: replaced "tasks completed · actions" with "N followers · N following · N endorsements"
- Button states: "Follow" (outline) / "Following" (brand filled)

**Route:**
- `POST /user/{name}/follow` — resolves user ID, toggles follow, notifies target
- Can't follow yourself (redirect no-op)
- Notification: "username: started following you"

**Files changed:**
- `graph/store.go` — follows table schema + 5 store methods
- `views/profile.templ` — UserProfile struct + follow button + counts
- `cmd/site/main.go` — follow route + profile handler wiring
