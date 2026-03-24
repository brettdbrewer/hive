# Scout Report — Iteration 201

## Gap: Work is defined too narrowly

**The problem:** We've been building "Linear competitor" — kanban board + task list. But Work at its most general is **organized activity toward outcomes**. A solo dev needs a todo list. A Fortune 500 needs departments, policies, compliance, resource allocation. A civilizational-scale system needs all of those across interconnected organizations. Our architecture (event graph, grammar ops, signed causal chains) can express ALL of this. We've been building 5% of the domain.

**What cognitive grammar reveals:** Work isn't a single product layer. Work is what happens when all 13 layers operate together on organized activity. The 13 layers already span the full domain — we just haven't composed them for Work the way we composed them for Social.

**Method:** Apply Need → Traverse → Derive to "organized activity" across all scales. Produce a converged Work spec that covers solo dev through civilizational, using the same Code Graph primitives.

**Output:** `hive/loop/work-general-spec.md` — the general Work specification, analogous to social-spec.md.

**This iteration produces spec, not code.**
