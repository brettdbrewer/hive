# Build Report — iter 239 (fix): computePipelineRoles test + constant

## Gap
Critic found two invariant violations in the iter 239 build:
1. `30*time.Minute` magic literal in `computePipelineRoles` (invariant 13 / no-magic-values)
2. `computePipelineRoles` had no test (invariant 12: VERIFIED)

## What changed

### `site/graph/handlers.go`
- Added `activeRoleThreshold = 30 * time.Minute` constant (near `maxHivePosts`)
- Replaced `30*time.Minute` literal in `computePipelineRoles` with `activeRoleThreshold`

### `site/graph/hive_test.go`
- Added `time` import
- Added `TestComputePipelineRoles`:
  - Verifies Builder is Active with a post within `activeRoleThreshold`
  - Verifies Scout is inactive with an old post (2h ago)
  - Verifies Critic has zero `LastActive` and is inactive when no posts exist

## Verification
- `go.exe build -buildvcs=false ./...` — clean, no errors
- `go.exe test -buildvcs=false ./graph/ -run "TestComputePipelineRoles|TestGetHive|TestParseCost|TestParseDuration|TestComputeHiveStats"` — all pass (integration tests skip without DB, expected)
