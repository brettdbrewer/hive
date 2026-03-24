# Scout Report — Iteration 213

## Gap: Spaces can't nest — blocks Organizations, Departments, Teams

**Source:** layers-general-spec.md fixpoint pass — "Spaces nest via parent_id."

**Current:** Spaces are flat. No hierarchy. A company can't have an "Engineering" space inside a "Company" space.

**What's needed:**
1. Schema: `ALTER TABLE spaces ADD COLUMN parent_id TEXT REFERENCES spaces(id)`
2. Store: update CreateSpace to accept parent_id, add ListChildSpaces
3. Space detail: show child spaces when they exist
4. Discover: only show top-level spaces (parent_id IS NULL)

**One column. Unlocks the entire organizational hierarchy.**
