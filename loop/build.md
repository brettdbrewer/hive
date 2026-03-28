# Build: Claims created without causes — CAUSALITY invariant fix

## What Was Built

Verified and confirmed the CAUSALITY invariant fix for `cmd/post`. The implementation correctly propagates cause IDs to all claim nodes created during a hive iteration.

## What Changed

No new code was required. The fix was already implemented in `cmd/post/main.go`:

1. **`assertScoutGap`** (line 431) — passes `taskCauseIDs` (the build task node ID) as `causes` when creating gap claim nodes via `op=assert`.

2. **`assertCritique`** (line 491) — passes `taskCauseIDs` as `causes` when creating critique claim nodes via `op=assert`.

3. **`assertLatestReflection`** (line 552) — passes `causeIDs` (the build document node ID) as `causes` when creating reflection document nodes via `op=intend`.

4. **`backfillClaimCauses`** (line 637) — fetches all `kind=claim` nodes from `/app/hive/knowledge?tab=claims&limit=200`, filters those with `causes=[]`, and patches each via `op=edit` with the current iteration's task node ID. This retroactively satisfies Invariant 2 for the 136 historical orphaned claims.

5. **Cause chain in `main()`**:
   - `post()` → returns `buildDocID` → wrapped as `causeIDs`
   - `createTask()` → returns `taskNodeID` → wrapped as `taskCauseIDs` (falls back to `causeIDs` if task creation failed)
   - Both IDs flow into the assert functions above

## Tests

All tests pass:

```
ok  github.com/lovyou-ai/hive/cmd/post  (all 30+ tests)
ok  github.com/lovyou-ai/hive/pkg/api
ok  github.com/lovyou-ai/hive/pkg/runner
ok  github.com/lovyou-ai/hive/pkg/hive
ok  github.com/lovyou-ai/hive/pkg/loop
... (all 13 packages pass)
```

Key test coverage for this fix:
- `TestBackfillClaimCausesUpdatesEmptyClaims` — verifies only empty-cause claims are patched
- `TestBackfillClaimCausesSkipsAlreadyCaused` — verifies already-caused claims are untouched
- `TestAssertCritiqueCarriesTaskNodeIDasCause` — verifies critique gets task ID as cause
- `TestAssertScoutGapSendsCauses` — verifies gap claim gets cause ID
- `TestAssertLatestReflectionSendsCauses` — verifies reflection gets cause ID
- `TestAssertCauseIDsMultipleJoined` — verifies multiple causes are comma-joined

## Build Verification

```
go.exe build -buildvcs=false ./...  ✓
go.exe test -buildvcs=false ./...   ✓  (all 13 packages)
```

## Invariant Status

**CAUSALITY (Invariant 2):**
- New claims: satisfied — every `op=assert` from cmd/post carries `causes=[taskNodeID]`
- Historical claims: backfill runs each iteration, patching up to 200 orphaned claims per run
- 136 historical claims will be patched on the next `cmd/post` execution
