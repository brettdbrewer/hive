# Build Report — Iteration 233: Auto-answer KindQuestion with document grounding

## Gap
Q&A loop was incomplete: questions were created but agents never answered them. The loop closes by having the Mind auto-answer any KindQuestion on creation, grounded in the space's KindDocument nodes.

## Changes

### `site/graph/store.go`
- Added `ListDocumentContext(ctx, spaceID) ([]Node, error)` — queries KindDocument nodes in a space, LIMIT 10 (BOUNDED invariant).

### `site/graph/mind.go`
- Added `OnQuestionAsked(spaceID, spaceSlug string, question *Node)` — async handler. Finds first agent, queries document context, builds grounded prompt, calls Claude, creates KindComment answer attributed to the agent.
- Added `buildQuestionAnswerPrompt(question *Node, docs []Node) string` — builds SOUL + ROLE + SPACE DOCUMENTS (up to 10 docs, 1000 chars each) + QUESTION prompt.

### `site/graph/handlers.go`
- Extended `express` case to support `kind=question`: creates KindQuestion, triggers `go h.mind.OnQuestionAsked(...)`, redirects to question detail. Standard `express` (no kind) unchanged.

### `site/graph/views.templ` + `views_templ.go` (generated)
- **QuestionsView form**: changed `op=intend` → `op=express`, `name=description` → `name=body`.
- **QuestionDetailView**: agent answers styled with violet badge (`bg-violet-950/20`, "agent" pill with dot). Human answers use standard surface card. Section header: "Answers (N)".
- Bug fix: answer form used `node_id` but `respond` handler expects `parent_id`. Fixed.

### `site/graph/handlers_test.go`
- Added `TestHandlerExpressQuestion`: two sub-tests verifying express+kind=question creates KindQuestion, and bare express still creates KindPost.

### `site/graph/mind_test.go`
- Added `TestBuildQuestionAnswerPrompt`: verifies doc context injection in prompt (with/without docs).
- Added `TestMindOnQuestionAsked_NoAgent`: verifies graceful no-op when no agent exists.

## Verification
- `templ generate` ✓ (15 updates)
- `go.exe build -buildvcs=false ./...` ✓
- `go.exe test ./...` ✓ (graph: 0.546s, all pass)
