# Build Report — Iteration 201

## Work General Specification

**Output:** `hive/loop/work-general-spec.md`

**Method:** Applied cognitive grammar (Distinguish → Relate → Select → Compose) to "organized activity toward outcomes" across all scales: solo dev → small team → mid-size company → enterprise → civilizational.

**Key findings:**

1. **12 entity types** (not just tasks): Task, Project, Goal, Role, Team, Department, Policy, Process, Decision, Resource, Document, Organization. All are Nodes with different kinds. No new tables needed.

2. **6 modes of work** (not just Execute): Execute (do it), Organize (structure it), Govern (control it), Plan (direct it), Learn (improve it), Allocate (resource it). Scale determines which modes you use.

3. **Grammar ops are universal.** Intend creates a task OR a policy OR a goal. Review approves work OR audits compliance. Consent votes on proposals OR approves budgets. The mechanism is the same; the mode determines semantics.

4. **Architecture already supports this.** No new schema. Every entity is a Node. Every mutation is an Op. The work is in UI modes and entity kinds, not foundations.

5. **Convergence at pass 2.** Existing foothold: Governance lens (Govern), Knowledge lens (Learn), Dashboard (Plan), Membership (Organize). The architecture was designed for this — we just hadn't composed it yet.

**This iteration produced spec, not code.**
