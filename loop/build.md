# Build: Observer reads /board for claim audit — misses 65 existing claims

## Gap
`buildPart2Instruction` only issued a curl to `GET /app/{slug}/board`, which returns `kind=task` nodes only. The Observer never fetched `/knowledge?tab=claims`, so every iteration it concluded "zero claims" and created false fix tasks.

## What Was Built

**`pkg/runner/observer.go` — `buildPart2Instruction`:**
- Added a second curl call: `GET /app/{slug}/knowledge?tab=claims&limit=50`
- Added audit item 6: "Claim integrity — claims with no body, no causes, or stuck in challenged state with no resolution"
- Added a note: "claims exist — do not report zero without checking"

**`pkg/runner/observer_test.go` — `TestBuildPart2Instruction`:**
- Added `wantClaimsURL bool` field to test cases
- Asserts claims URL present when apiKey is set, absent when apiKey is empty

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (pkg/runner: 3.736s)
