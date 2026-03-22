# Build Report — Iteration 36

## What Was Planned

Agent badges on People and Activity lenses — make agent identity visible across all six lenses.

## What Was Built

**site/graph/store.go**:
- `ActorKind` field added to `Op` struct
- `ListOps` updated: `LEFT JOIN users u ON u.name = o.actor`, selects `COALESCE(u.kind, 'human')` as actor_kind
- `ListNodeOps` updated with same JOIN pattern
- No schema migration — resolved at query time from users table (lesson 30)

**site/graph/handlers.go**:
- `Kind string` field added to `Member` struct
- People handler now sets `Kind` from `Op.ActorKind` during aggregation

**site/graph/views.templ**:
- **People lens**: agent members get violet avatar + "agent" badge pill (matching Feed/Chat patterns)
- **Activity lens**: agent ops get violet avatar + "agent" badge pill (same pattern)
- Both use `o.ActorKind == "agent"` / `m.Kind == "agent"` conditionals

## Key Design Decisions

1. **JOIN at query time, not schema migration**: The ops table has no `actor_kind` column. Instead of adding one (and changing `RecordOp`'s 10+ call sites), we JOIN against the users table at query time. This follows lesson 30 — the users table is the authority for actor properties.

2. **COALESCE to 'human'**: If an actor has no matching user record (e.g., early ops before user records existed), default to "human". Safe assumption for historical data.

3. **Identical visual treatment**: Same violet avatar + badge pill used in Feed, Chat, Comments, People, and Activity. All six lenses now have consistent agent identity.

## Verification

- `templ generate` + `go build` — clean
- Deployed to Fly.io — healthy
- JOIN pattern tested via code review (LEFT JOIN with COALESCE handles null gracefully)

## Files Changed

- `site/graph/store.go` — ~15 lines (ActorKind field, two query JOINs)
- `site/graph/handlers.go` — 2 lines (Kind field on Member, set from ActorKind)
- `site/graph/views.templ` — ~20 lines (conditional avatars + badges in People and Activity)
