# Build Report — Iteration 225: Builder Ships Code to Production

## What This Iteration Does

Fixes 3 critique issues from iter 224, then runs the builder on a real coding task. The hive ships its first autonomous code commit to production.

## Runner Fixes

### 1. Removed double role prompt
Provider no longer gets `SystemPrompt: rolePrompt`. Role prompt is only in the instruction (where it has task context). Saves ~500 tokens per call.

### 2. Recency tiebreak in task selection
`pickHighestPriority` now tiebreaks by `created_at` descending when tasks share the same priority. Newest tasks are preferred — most likely to be fresh assignments, not stale work.

### 3. Changes-required guard
After ACTION: DONE, builder checks `git status --porcelain`. If no files changed, the task stays in-progress with a comment explaining the issue. Prevents hollow completions.

### 4. Extracted `hasUncommittedChanges()` helper
Factored git status check out of `commitAndPush` into reusable helper.

## Builder E2E Result

```
[builder] working task 61f38992: Add Policy entity kind to the site
  ⏳ working done (2m49s)
[builder] Operate done (cost=$0.5325, tokens=31+7912)
[builder] action: DONE
[builder] committed and pushed: [hive:builder] Add Policy entity kind to the site
[builder] task 61f38992 DONE
```

**2m49s. $0.53. One Operate call. Real code committed and pushed.**

## Hive's Code Changes (autonomous)

| File | What |
|------|------|
| `site/graph/store.go` | `KindPolicy = "policy"` constant |
| `site/graph/handlers.go` | `handlePolicies` handler (34 lines) |
| `site/graph/views.templ` | `policiesIcon()`, sidebar + mobile nav, `PoliciesView` template (81 lines) |
| `site/graph/views_templ.go` | Generated |

## Human Fix (Critic catch)

The builder missed adding `KindPolicy` to the `intend` op's kind allowlist (line 1487 of handlers.go). Without this, creating a policy via the API would fall through to KindTask. Fixed manually and deployed.

## Files Changed (hive repo)

| File | What |
|------|------|
| `cmd/hive/main.go` | Removed SystemPrompt from provider config |
| `pkg/runner/runner.go` | Recency tiebreak, changes-required guard, hasUncommittedChanges helper |
| `pkg/runner/runner_test.go` | 2 new tests (recency tiebreak, priority-beats-recency) |

## Build

- `go build ./...` ✓ (hive + site)
- `go test ./...` ✓ (14 runner tests + all existing)
- `flyctl deploy --remote-only` ✓ — Policy entity live on lovyou.ai
