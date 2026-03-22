# Build Report — Iteration 50

Added `ResolveUserNames(ctx, ids)` to store — batch lookup of IDs to display names. Conversation templates accept `nameMap map[string]string`, display resolved names via `resolveName(tag, nameMap)` helper. Handlers resolve IDs before rendering.

Backfilled 30 nodes (author_id) and 30 ops (actor_id) from users table. Cleaned up conversation tags to contain only user IDs. Removed test data from production.
