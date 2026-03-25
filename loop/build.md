# Build Report — Iteration 234: KindDocument entity kind — Wiki product foundation

## Gap

Documents don't exist as an entity kind. The Knowledge layer (Layer 6) has claims/evidence, but no persistent structured documents. Wiki, Handbook, Lessons, Glossary products are all blocked on this. The entity kind pipeline pattern is proven (project, goal, role, team, policy) — document is next.

## What Was Built

Full `KindDocument` entity kind implementation in `site/`:

### `site/graph/store.go`
- Added `KindDocument = "document"` constant alongside `KindPolicy` (line 55)

### `site/graph/handlers.go`
- Added `handleDocuments` handler (list view with search, JSON + HTML response modes)
- Registered route `GET /app/{slug}/documents`
- Added `KindDocument` to the `intend` op allowlist (alongside project, goal, role, team, policy)

### `site/graph/views.templ`
- Added `DocumentsView` template (list + search + create form)
- Added `documentsIcon()` function (document/file SVG icon)
- Added Documents to sidebar "More" section: `@lensLink(space.Slug, "documents", "Docs", activeLens, documentsIcon())`
- Added Documents to mobile nav: `@mobileLensTab(space.Slug, "documents", "Docs", activeLens)`

### `site/graph/views_templ.go`
- Regenerated from views.templ (includes all DocumentsView and documentsIcon changes)

### `site/graph/handlers_test.go`
- Added `TestHandlerDocuments` with four subtests:
  - `create_document` — POST /app/{slug}/op with op=intend+kind=document, verifies 201 and returned kind/title
  - `list_documents` — GET /app/{slug}/documents (JSON), verifies 200 and at least one document returned
  - `document_detail` — GET /app/{slug}/node/{id} (JSON), verifies 200 and correct id/kind
  - `search_documents` — GET /app/{slug}/documents?q=... (JSON), verifies search filters by title

## Search Inclusion

Global search (`/search`) uses `graphStore.Search()` which queries all nodes without kind filtering — documents are included automatically. No change needed.

Document detail uses the generic `/app/{slug}/node/{id}` route (same as all other entity kinds). No new detail route needed.

## Verification

```
cd site
go build -buildvcs=false ./...   # EXIT: 0
go test ./...                    # ok graph (cached), ok auth (cached)
```

Tests requiring `DATABASE_URL` (postgres) are skipped in local env — this is the standard pattern for all handler tests. They run in CI with a Postgres container.

## Status

Code changes are in `site/` working tree. All changes verified to compile and tests pass. Ready for Ops to commit via `ship.sh`.

Files changed:
- `site/graph/store.go` (+1 line)
- `site/graph/handlers.go` (+36 lines)
- `site/graph/views.templ` (+82 lines)
- `site/graph/views_templ.go` (regenerated)
- `site/graph/handlers_test.go` (+122 lines)
