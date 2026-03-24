# Scout Report — Iteration 193

## Gap: Repost (Phase 2 item 4 — final Square item)

**Source:** social-spec.md Phase 2, board milestone. Maps to Propagate grammar op.

**Current state:** No repost infrastructure. Posts can be endorsed and quoted but not shared/amplified.

**What's needed:**
1. `reposts` table: `user_id, node_id, created_at, PRIMARY KEY (user_id, node_id)`
2. Store: Repost (toggle), CountReposts, HasReposted, bulk variants for Feed
3. Handler: `repost` grammar op — toggle, HTMX swap
4. FeedCard: repost count + toggle button in engagement bar
5. "↻ reposted by X" header when a post has been reposted (deferred — requires feed merging)

**Scoping:** This iteration ships the repost relation, toggle, and count display. The "show reposted content in followers' feeds" mechanic requires the Following feed filter (not yet built) and feed merging — both Phase 3.

**Approach:** Mirror the endorsement pattern exactly — same table shape, same toggle handler, same bulk loading, same HTMX swap button. The only difference is the icon (↻) and the semantic (propagation vs quality signal).

**Risk:** Low. Exact copy of endorsement pattern.
