# Build Report — Iteration 193

## Repost (Propagate Grammar Op)

**Schema:**
- `reposts` table: `user_id, node_id, created_at, PRIMARY KEY (user_id, node_id)`
- Index on `node_id` for count queries

**Store:**
- `Repost(userID, nodeID)` — ON CONFLICT DO NOTHING (idempotent)
- `Unrepost(userID, nodeID)` — DELETE
- `HasReposted(userID, nodeID)` — EXISTS check
- `GetBulkRepostCounts(nodeIDs) map[string]int` — counts per node
- `GetBulkUserReposts(userID, nodeIDs) map[string]bool` — which nodes user reposted

**Handler:**
- `repost` grammar op: toggle (repost/unrepost), records op, notifies post author
- HTMX response: returns `repostButton` component for inline swap
- JSON response: `{"op": "repost", "reposted": true/false}`

**Feed handler:**
- Loads bulk repost counts + user repost state alongside endorsement data
- Passes both maps to FeedView

**Template:**
- `FeedView`: accepts repost maps
- `FeedCard`: accepts `repostCount int, reposted bool`
- `repostButton`: ↻ icon (arrows), emerald color when reposted, HTMX toggle
- Engagement bar order: replies → repost → quote → endorse

**Files changed:**
- `graph/store.go` — reposts table + 5 store methods
- `graph/handlers.go` — repost op + feed handler wiring
- `graph/views.templ` — FeedView, FeedCard, repostButton
