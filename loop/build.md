# Build Report — Fix title compounding + PRMode config infrastructure

## What Changed

### `pkg/runner/critic.go` — fixTitle deduplication
Replaced the `if`-based double-prefix guard with `strings.TrimPrefix`:

```go
// Before
func fixTitle(subject string) string {
    if strings.HasPrefix(subject, "Fix: ") {
        return subject
    }
    return "Fix: " + subject
}

// After
func fixTitle(subject string) string {
    return "Fix: " + strings.TrimPrefix(subject, "Fix: ")
}
```

Idempotent: any subject already starting with "Fix: " gets stripped before the single prefix is prepended.

### No changes needed
- `PRMode bool` in `pkg/runner/runner.go` Config struct: already present (line 58)
- `--pr` flag in `cmd/hive/main.go`: already present (line 65), wired to `Config.PRMode`

## Verification

```
go.exe build -buildvcs=false ./...   → clean
go.exe test ./...                    → ok pkg/runner (0.687s), all others cached/ok
```
