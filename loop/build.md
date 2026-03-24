# Build Report — Iteration 205

## Projects — First New Entity Kind

**Store:**
- `KindProject = "project"` constant added
- No new store methods — existing `ListNodes(kind=project)` and `CreateNode(kind=project)` work as-is

**Handler:**
- `handleProjects` — lists project nodes for a space, search support
- `intend` op: now accepts `kind` form field. `kind=project` creates a project; defaults to task.
- Route: `GET /app/{slug}/projects`

**Template:**
- `ProjectsView` — project list with title, description, task progress (done/total + bar), status badge
- Create form: hidden `kind=project` input, title + description
- `projectsIcon` — folder SVG icon

**Sidebar + Mobile:**
- "Projects" added between Board and Feed in both desktop sidebar and mobile nav

**The proof:** Adding a new entity kind required:
- 1 constant (1 line)
- 1 handler (~30 lines, copypaste of handleChangelog pattern)
- 1 template (~80 lines)
- 1 line in `intend` op (accept kind param)
- 0 schema changes, 0 new store methods

The grammar is genuinely kind-agnostic.
