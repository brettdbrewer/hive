# Build: Fix state.md — CAUSALITY items 1-2 and infra item 11 still marked open after iter 404

## Task

Items 1 and 2 of the CAUSALITY milestone were completed in iter 404 but remained open in state.md. Item 11 in infrastructure gaps was still marked NOT DEPLOYED. The Scout would re-deploy unnecessarily on the next pass (Lesson 208/209).

## What Was Built

Documentation-only change to `loop/state.md`:

1. **GATE 1 (deploy populateFormFromJSON)** — struck and moved to DONE section. Marked DONE (iter 404), tasks 0d617293 + c0a686e2.
2. **GATE 2 (observer fallback cause)** — struck and moved to DONE section. Marked DONE (iter 404), critique PASS 2f1e9125.
3. **GATE 3 promoted to GATE 1** — `assertClaim` wrapper (Lesson 167) is now the active gate.
4. **Stale "NOT deployed" notices** — two stale paragraphs at lines 49 and 49 removed (both said populateFormFromJSON was not deployed to production).
5. **Infrastructure gap item 11** — struck and marked DONE (iter 404), deployed to production.

## Verification

- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (14 packages)

ACTION: DONE
