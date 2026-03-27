# Build: Add Preview truncation 2000 chars to Reflector empty_sections diagnostic

## Task

In `pkg/runner/reflector.go`, update the `empty_sections` early-return path to truncate the `Preview` field to 2000 chars (was 500), matching the task spec. The `Preview` field itself was already present from the prior iteration. This change ensures the diagnostic captures enough context for root-cause analysis without being arbitrarily short.

## Changes

### `pkg/runner/reflector.go`

- Line 234: Changed truncation from `> 500` / `[:500]` to `> 2000` / `[:2000]`
- No structural changes — Preview was already wired into the `appendDiagnostic` call

## Verification

- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (`pkg/runner` 2.943s, all others cached)
- `TestRunReflectorEmptySectionsDiagnostic` confirms `Preview != ""` — passes
