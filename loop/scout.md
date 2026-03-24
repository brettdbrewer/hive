# Scout Report — Iteration 205

## Gap: No Project entity — tasks are flat, ungrouped

**Source:** unified-spec.md entity table. work-general-spec.md. Every scale needs projects.

**Current state:** Tasks live in a flat list per space. Subtasks exist (parent_id) but there's no higher-level grouping. A space with 50 tasks is a wall of cards. No way to say "these 12 tasks are the Auth project" and "those 8 are the API project."

**What's needed:**
1. `project` node kind — a container for tasks with its own title, description, status
2. Tasks link to projects via parent_id (already works — subtask mechanism)
3. Projects view — list of projects with task counts, progress bars
4. Board grouped by project (optional filter)

**Why this first:** Project is the first new entity kind from the unified spec. It proves "adding entity kinds is just adding a node kind + a view." No schema changes. No new ops — Intend creates a project just like it creates a task. The grammar is kind-agnostic.

**Approach:**
- Add `KindProject = "project"` constant
- Add a "Projects" page: `/app/{slug}/projects` — lists project nodes with child task counts
- Project detail: just NodeDetailView (already works for any node kind)
- Board: add project filter dropdown

**Risk:** Low. One new constant, one new handler, one new template. Existing infrastructure handles the rest.
