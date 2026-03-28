# Fix: /hive diagnostics always empty in production — migrate to graph API

## What changed

Migrated `/hive/feed` diagnostics from local file (`loop/diagnostics.jsonl`) to the
graph database, eliminating the Fly.io production gap where `HIVE_REPO_PATH` was not
set and the file never shipped to the container.

## Files changed

### `site/graph/store.go`
- Added `hive_diagnostics` table to the auto-migration schema (`CREATE TABLE IF NOT EXISTS`)
- Added `AppendHiveDiagnostic(ctx, phase, outcome, costUSD, payload)` — inserts a phase event row
- Added `ListHiveDiagnostics(ctx, limit)` — queries the last N events newest-first

### `site/graph/handlers.go`
- Added `POST /api/hive/diagnostic` route (auth-protected via `writeWrap`)
- Added `handleHiveDiagnostic` handler — parses JSON body, calls `AppendHiveDiagnostic`
- Updated `handleHiveFeed` — queries DB first, falls back to local file (dev compatibility)
- Updated `handleHive` — same DB-first pattern for the initial page load
- Added `"io"` to imports

### `site/graph/hive_test.go`
- Added `TestPostHiveDiagnostic_StoresAndServes` — round-trip test: POST then GET /hive/feed
- Added `TestListHiveDiagnostics_Empty` — empty DB returns nil without error

### `hive/pkg/api/client.go`
- Added `PostDiagnostic(payload []byte)` — POSTs raw JSON to `/api/hive/diagnostic`

### `hive/pkg/runner/diagnostic.go`
- Updated `Runner.appendDiagnostic` to:
  1. Still write to `loop/diagnostics.jsonl` when `HiveDir` is set (local dev)
  2. Also POST via `APIClient.PostDiagnostic` when APIClient is set (production)

## Verification

```
site: go build -buildvcs=false ./...   pass
site: go test ./...                    pass (DB tests skip without DATABASE_URL)
hive: go build -buildvcs=false ./...   pass
hive: go test ./...                    pass (all runner tests pass)
```

## How production now works

1. Hive runner emits `PhaseEvent` → `Runner.appendDiagnostic` → POSTs JSON to `POST /api/hive/diagnostic` on lovyou.ai
2. Site stores row in `hive_diagnostics` table (auto-created on startup)
3. `GET /hive/feed` → `ListHiveDiagnostics` → renders phase timeline from DB

No `HIVE_REPO_PATH` or Fly volume required.
