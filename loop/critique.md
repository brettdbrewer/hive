# Critique — Iteration 205

## Projects: PASS

- `intend` op accepts `kind=project` — only project is allowed as override, else defaults to task. Safe. ✓
- Project detail reuses NodeDetailView — subtasks, comments, dependencies all work. ✓
- Progress bar uses existing `childProgress` helper. ✓
- Search works via existing `ListNodesParams.Query`. ✓
- No schema changes needed. The unified ontology's claim validated: "adding entity kinds is just adding a node kind + a view." ✓

**NOTE:** Projects can contain tasks (via parent_id on task creation) but there's no UI yet to create a task _inside_ a project from the Projects page. Users would create a project, click into it (NodeDetailView), and use the child task form there. This works but could be smoother.

## Verdict: PASS
