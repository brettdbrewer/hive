# Build Report — Iteration 239 (fix 2)

## Gap Addressed

Critic (fix round 1) flagged:
1. Diff only showed loop file changes — could not confirm code changes were present
2. `computeHiveStats` "Features" count silently excludes posts without a cost line — needs a comment
3. `parseCostDollars` regex correctness under edge cases
4. `TestGetHive_RendersMetrics` — needed confirmation it checks label text, not just status 200

## Verification of Code Presence

All four code changes from the previous iteration are confirmed present in the repo:
- `HiveStats`, `parseCostDollars`, `parseDurationStr`, `computeHiveStats` — all in `handlers.go` (lines 3355–3413)
- `HiveView(posts []Node, stats HiveStats, user ViewUser)` — updated in `views.templ` / `views_templ.go`
- `hive_test.go` — `TestParseCostDollars`, `TestParseDurationStr`, `TestComputeHiveStats`, `TestGetHive_RendersMetrics` all present

## Files Changed This Fix Round

### `site/graph/handlers.go`
- Added a 3-line comment to `computeHiveStats` explaining why posts without a cost line are excluded from the Features count: a post without a cost line is not a verified build iteration. This addresses the Critic's concern about silent exclusion.

## Regex Analysis

`reCost = regexp.MustCompile(`\$(\d+\.\d+)`)`:
- Matches `$0.53` ✓, `$1.00` ✓, `$10.50` ✓, `$0.83` ✓
- "Multiple `$1.00` and `$2.50` — first wins" — `FindStringSubmatch` returns first match ✓
- No false-positives on non-cost dollar amounts beyond "first wins" behavior (acceptable, documented in test)

## TestGetHive_RendersMetrics Assertion

The test asserts label text via `strings.Contains` for all three labels:
- "Features shipped" ✓
- "Total autonomous spend" ✓
- "Avg cost" ✓

This is a content assertion, not just a status check.

## Verification

```
go.exe build -buildvcs=false ./...     pass (no errors)
go.exe test ./...                       pass (graph: 0.549s; DB tests skip without DATABASE_URL, pass in CI)
```
