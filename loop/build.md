# Build Report — Document Edit Handler

## Gap
Documents had no dedicated edit endpoint. The CRUD loop was incomplete: create (via `/op intend`) and read (`/node/{id}`) existed but edit was missing.

## What Changed

### `site/graph/handlers.go`
- Added two routes: `GET /app/{slug}/document/{id}/edit` and `POST /app/{slug}/document/{id}/edit`
- Added `handleDocumentEdit` handler:
  - Uses `spaceFromRequest` (write access gate — owner or authenticated member of public space)
  - GET: fetches space + node, verifies `Kind == KindDocument`, renders `DocumentEditView`
  - POST: parses title+body form values, calls `store.UpdateNode`, redirects to node detail (or returns JSON if requested)
  - Returns 404 for non-members, non-documents, and missing nodes

### `site/graph/views.templ`
- Added `DocumentEditView` template: full-page edit form with title input + body textarea (16 rows, monospace, resizable), breadcrumb, Save/Cancel actions
- Added edit link in `NodeDetailView` for documents: appears in the metadata row for owners, links to `/app/{slug}/document/{id}/edit`

### `site/graph/handlers_test.go`
- Added `TestHandlerDocumentEdit` with three subtests:
  - `get_edit_form_member` — owner/member can GET the edit form (200)
  - `post_edit_member` — owner/member can POST an update; verifies title+body persist
  - `non_member_rejected` — different user on a private space gets 404

## Verification
- `templ generate` — no errors
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (graph: 0.532s)
