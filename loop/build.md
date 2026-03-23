# Build Report — Iteration 123

## Add dependency form — dropdown to create task dependencies from node detail

### Changes

**handlers.go:**
- `handleNodeDetail` now fetches space tasks (kind=task), filters out self and existing deps, passes to template for dropdown

**views.templ:**
- `NodeDetailView` signature extended with `spaceTasks []Node`
- Dependencies section now shows "Add dependency..." dropdown when authenticated user is viewing a task
- Dropdown shows task title + state label, excludes self and existing dependencies

### Files changed
- `graph/handlers.go` — space task fetch + filtering
- `graph/views.templ` — form + dropdown
- `graph/views_templ.go` — generated

### Deployed
`ship.sh` — all green. No 408 this time.
