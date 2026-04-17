# Debt: JSONB Expression Index for `hasSiteOpReceived`

**Status:** Tracked. Not a correctness issue; defer until site-op volume crosses ~10k anchored rows.
**Opened:** 2026-04-17
**Owner:** unassigned
**Origin:** PR [#73](https://github.com/transpara-ai/hive/pull/73) (event-driven bridge migration, Prompt 4 of 6) — flagged during code review.

## Problem

`pkg/reconciliation/watermark.go:hasSiteOpReceived` is the idempotency gate that prevents the reconciliation ticker from re-anchoring a site op the webhook already wrote to the chain:

```sql
SELECT EXISTS(
    SELECT 1 FROM events
    WHERE event_type = 'site.op.received'
      AND content_json->'external_ref'->>'id' = $1
)
```

`idx_events_type` on `events(event_type)` narrows the first filter fast. The JSONB path extraction `content_json->'external_ref'->>'id'` is **not** indexed. At low anchored-row counts this is invisible. Once `site.op.received` accumulates ~10k rows, every reconciliation cycle does a linear JSONB scan over the full `site.op.received` subset for each op being checked — O(ops_in_cycle × anchored_rows).

A 60s cycle that lists 100 ops and scans 50k anchored rows means 5M JSONB extractions per minute. The planner has nowhere to go; it has to read every row.

## Fix

Partial expression index, scoped to the row type the query actually hits:

```sql
CREATE INDEX IF NOT EXISTS idx_events_site_op_received_ref
    ON events((content_json->'external_ref'->>'id'))
    WHERE event_type = 'site.op.received';
```

The partial predicate keeps the index small (only `site.op.received` rows). The expression is exactly what the query extracts, so the planner picks it up without changes to `hasSiteOpReceived`.

## Where it lives

DDL goes in `pkg/reconciliation/schema.go` alongside `reconciliation_state`, wrapped in `CREATE INDEX IF NOT EXISTS` so existing deployments pick it up on next `EnsureTables` call — no migration step needed.

## Acceptance

- [ ] Partial expression index added to `pkg/reconciliation/schema.go`.
- [ ] `EXPLAIN ANALYZE` on `hasSiteOpReceived` shows an index scan (not a seq scan or bitmap heap scan over the full `site.op.received` subset) when the subset has >1k rows.
- [ ] No change to the `hasSiteOpReceived` query text — the planner picks the index automatically.

## When to do it

Watch for either trigger:

- `site.op.received` row count > 10k, or
- reconciliation cycle duration p95 > 1s in telemetry.

Either signal means the scan is starting to dominate cycle cost.

## Why not now

The PR #73 reviewer surfaced three CRITICAL findings, all hallucinated. This was the one legitimate perf follow-up. Fixing it now has no observable benefit and adds a startup DDL step on every deploy — deferring until there's actual load justifies the change.
