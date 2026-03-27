# Build Report — Iter 339: Show cost and duration on hive build log entries

## Gap
Build log rows showed title + timestamp but not cost or duration — visitors couldn't see what each iteration cost at a glance.

## Changes

### `site/graph/handlers.go`
Added two package-level helpers callable from templ templates:
- `hiveCostStr(n Node) string` — calls `parseCostDollars(n.Body)`; returns `"$0.42"` if cost > 0, else `""`
- `hiveDurationStr(n Node) string` — calls `parseDurationStr(n.Body)`; returns duration string or `""`

### `site/graph/views.templ`
In `HiveStatusPartial`, updated the build log row layout:
- Wrapped `<time>` in a `<div class="flex items-center gap-2 flex-shrink-0">` container
- After `<time>`, conditionally renders a brand-accent pill badge for cost (`$0.42`)
- After cost badge, conditionally renders a muted text span for duration (`3m28s`)

### `site/graph/hive_test.go`
Added two new unit tests:
- `TestHiveCostStr` — covers post with cost, post with no cost, and $0.00 (returns empty)
- `TestHiveDurationStr` — covers post with duration, post without, and empty body

### Generated
- `site/graph/views_templ.go` — regenerated via `templ generate`

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (graph: 0.084s)
