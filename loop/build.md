# Build: cmd/post: dedup loop header tasks (Iteration N) on board

## Scout Gap Addressed

Scout 406 gap: missing typed `assertClaim` wrapper in `cmd/post` (CAUSALITY GATE 1, Lesson 167). This iteration addresses a prerequisite: duplicate "Iteration N" and "Target repo" tasks accumulating on the board on every loop run, which is a board hygiene blocker. The dedup guard was only firing for "Fix:"-prefixed titles — all other titles bypassed the board check entirely.

## What Was Built

### `hive/cmd/post/main.go` — unconditional dedup in `createTask`

Removed the `coreTitle != title &&` guard from the dedup check, so `findExistingTask` is called for **all** non-empty titles, not just "Fix:"-prefixed ones.

Before:
```go
coreTitle := stripFixPrefixes(title)
if coreTitle != title && coreTitle != "" {
```

After:
```go
coreTitle := stripFixPrefixes(title)
if coreTitle != "" {
```

### `hive/cmd/post/main_test.go` — updated `TestCreateTaskNoDedup`

Updated the test to reflect the new unconditional dedup behavior:
- Board is now queried for all titles (test no longer asserts board must not be queried)
- Test verifies: board queried → empty result → new task still created via `op=intend`

## Build Results

```
go.exe build -buildvcs=false ./...   → OK
go.exe test ./...                    → all 11 packages pass
```

## Files Changed

- `hive/cmd/post/main.go` — removed `coreTitle != title &&` guard (1 line change)
- `hive/cmd/post/main_test.go` — updated `TestCreateTaskNoDedup` to match new behavior
