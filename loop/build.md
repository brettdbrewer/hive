# Build Report — Iteration 28

## What Was Planned

Space previews on discover cards — node count, last activity timestamp, sorted by recent activity.

## What Was Built

**graph/store.go**: Added `SpaceWithStats` type (embeds `Space` + `NodeCount int` + `LastActivity *time.Time`). Changed `ListPublicSpaces` return type from `[]Space` to `[]SpaceWithStats`. Query uses `LEFT JOIN LATERAL` to compute per-space node count and max created_at. Sorting changed from `created_at DESC` to `COALESCE(last_at, created_at) DESC` — active spaces float to top.

**views/discover.templ**: Added `NodeCount` and `LastActivity` fields to `DiscoverSpace`. Updated `discoverCard` to show item count ("3 items") and relative time ("2h ago") below the description. Added `pluralize()` and `relativeTime()` helpers. Relative time shows: just now, Xm ago, Xh ago, Xd ago, or month/year for older content.

**cmd/site/main.go**: Updated `DiscoverSpace` mapping to pass `NodeCount` and `LastActivity` from store results.

4 files changed, deployed.

## What Works

- Discover cards show node count and last activity for spaces with content
- Spaces with zero nodes show no stats line (clean, not "0 items")
- Spaces sorted by most recent activity — active spaces appear first
- Relative timestamps render correctly (verified: "2 items", "1 item", timestamps showing)
- LATERAL JOIN uses existing `idx_nodes_space` index
