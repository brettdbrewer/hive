# Build: Observer audit: map loop artifacts to correct node kinds

## Task

All 491 hive board nodes are kind=task. The loop was only emitting `intend` op (which defaults to task). Critiques, reflections, and build reports were either missing from the graph or using the wrong kind.

## What Was Built

**`cmd/post/main.go`** — four changes:

1. **`createTask()`** — added explicit `"kind": "task"` to the intend payload. Functionally identical but now declares the kind rather than relying on the server default.

2. **`post()`** — changed from `op=express`/`kind=post`/`body` to `op=intend`/`kind=document`/`description`. Build reports are structured build records, not casual feed posts. They now appear in the Documents lens.

3. **`assertCritique()`** (new) — reads `loop/critique.md`, extracts the first heading as title, POSTs `op=assert`/`kind=claim` to persist the critique verdict as a searchable KindClaim node. Also added `extractCritiqueTitle()` helper.

4. **`assertLatestReflection()`** (new) — reads `loop/reflections.md`, extracts the first `##` section (the most recent entry), POSTs `op=intend`/`kind=document` with title `"Reflection: <date>"`. Also added `extractLatestReflection()` helper.

Both new functions are called from `main()` as non-fatal operations (same pattern as `assertScoutGap`).

**`cmd/post/main_test.go`** — updated and extended:

- Renamed `TestPostCreatesNode` → `TestPostCreatesDocument`; updated assertions to check `op=intend`, `kind=document`, `description` field.
- Updated `TestBuildTitleExtractedOnPost` to match new op/kind.
- Added `TestExtractCritiqueTitle` (3 subtests)
- Added `TestAssertCritiqueCreatesClaimNode`
- Added `TestAssertCritiqueMissingFile`
- Added `TestExtractLatestReflection`
- Added `TestExtractLatestReflectionNoEntry`
- Added `TestAssertLatestReflectionCreatesDocument`
- Added `TestAssertLatestReflectionMissingFile`

## Kind mapping after this build

| Artifact | Op | Kind | Lens |
|---|---|---|---|
| Scout gap | assert | claim | Knowledge |
| Build report | intend | document | Documents |
| Critique verdict | assert | claim | Knowledge |
| Latest reflection | intend | document | Documents |
| Board task | intend | task | Board |

## Diff stat

- `cmd/post/main.go`: +116 lines (new functions + helpers + main() calls)
- `cmd/post/main_test.go`: +170 lines (updated + new tests)

## Verification

```
go build -buildvcs=false ./cmd/post/   # OK
go test -buildvcs=false ./...           # all pass (23 tests in cmd/post)
```
