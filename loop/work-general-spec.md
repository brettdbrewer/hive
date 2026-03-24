# Work: General Specification

**Work as organized activity toward outcomes, at every scale. Derived via cognitive grammar from the Code Graph primitives.**

Matt Searles + Claude · March 2026

---

## The Insight

Work isn't a product layer. Work is what happens when all 13 layers operate together on organized activity. A kanban board is one view of one mode of one scale of work. The architecture — event graph, grammar ops, signed causal chains — can express the full domain without modification.

The question isn't "how do we build a task tracker." It's "what are all the things organized activity requires, and how do they compose from our primitives?"

---

## Distinguish: The Entities of Work

Applied across scales (solo dev → team → company → enterprise → civilizational):

| Entity | What it is | Node kind | Solo dev | Mid-size co. | Fortune 500 | Civilizational |
|--------|-----------|-----------|----------|-------------|-------------|----------------|
| **Task** | Atomic work unit | `task` | Fix bug | Ship feature | Compliance audit item | Deliver food to region |
| **Project** | Scoped work collection | `project` | My app | Q2 roadmap | Enterprise migration | City food program |
| **Goal** | Desired outcome | `goal` | Ship v1 | 10k users | $1B revenue | Reduce hunger 20% |
| **Role** | Capability + responsibility | `role` | "Dev" | "Senior Engineer" | "VP Engineering" | "Regional Coordinator" |
| **Team** | Group around a function | `team` | — | "Backend" | "Platform Infra" | "Logistics Fleet 7" |
| **Department** | Organizational unit | `department` | — | "Engineering" | "Global Engineering" | "Operations" |
| **Policy** | Rule governing behavior | `policy` | — | "Code review required" | "SOC2 compliance" | "Food safety protocol" |
| **Process** | Repeatable sequence | `process` | Deploy script | Onboarding flow | Incident response | Supply chain pipeline |
| **Decision** | Choice with rationale | `decision` | Use Postgres | Hire candidate A | Enter market X | Allocate $10M to region |
| **Resource** | Anything consumed | `resource` | Time | Budget, headcount | CapEx, OpEx, compute | Regional GDP allocation |
| **Document** | Knowledge artifact | `document` | README | Spec, ADR | Contract, handbook | Treaty, regulation |
| **Organization** | Legal/structural entity | `organization` | — | The company | Subsidiary / BU | Nation, NGO, consortium |

**Key insight:** Every entity is a Node on the graph. Every mutation is an Op. The grammar operations apply uniformly. Creating a Policy is Intend. Approving it is Consent. Reviewing it is Review. Challenging it is Challenge. The SAME operations that manage tasks also manage departments, policies, goals, and decisions.

---

## Relate: How Entities Connect

```
Organization ──contains──→ Department ──contains──→ Team ──contains──→ Role ──filled_by──→ User/Agent
     │                          │                       │
     └──governed_by──→ Policy   └──owns──→ Project      └──works_on──→ Task
                          │                    │                          │
                          └──requires──→ Process │                        └──toward──→ Goal
                                               └──decomposes──→ Task
                                                                  │
                                               Resource ──consumed_by──┘
                                               Document ──supports──┘
                                               Decision ──affects──┘
```

On our graph, these are all parent-child relationships (parent_id), dependencies (node_deps), or tag associations (tags[]). No new schema needed.

---

## Select: The Six Modes of Work

Like Social has Chat/Rooms/Square/Forum, Work has six modes. Each mode is a lens on the same graph data, optimized for a different aspect of organized activity.

| Mode | What it does | Replaces | Primary Ops | Who uses it |
|------|-------------|----------|-------------|-------------|
| **Execute** | Do the work | Linear, Jira, Asana | Intend, Claim, Complete, Review | Everyone |
| **Organize** | Structure the work | BambooHR, Workday, Org charts | Delegate, Assign, Scope | Managers, HR, Admins |
| **Govern** | Control the work | Policy tools, compliance, GRC | Consent, Propose, Vote, Enforce | Leaders, Compliance, Legal |
| **Plan** | Direct the work | Lattice, OKR tools, Roadmaps | Intend, Decompose, Prioritize | PMs, Directors, Strategy |
| **Learn** | Improve the work | Confluence, Notion, Retrospectives | Reflect, Assert, Challenge, Endorse | Everyone (after work) |
| **Allocate** | Resource the work | Spreadsheets, ERP, Budgeting | Claim, Prioritize, Scope, Consent | Finance, Operations, Leads |

### Scale determines which modes you use:

| Scale | Modes used |
|-------|-----------|
| **Solo dev** | Execute (tasks + agent) |
| **Small team** | Execute + Plan |
| **Mid-size company** | Execute + Plan + Organize + Learn |
| **Enterprise** | All six |
| **Civilizational** | All six, across interconnected organizations |

The product grows with the user. A solo dev sees a kanban board. When they add a team, Organize appears. When they need compliance, Govern appears. The modes aren't feature gates — they emerge from what entities exist in the graph.

---

## Grammar Operation Coverage Matrix

How each grammar op manifests across the six Work modes:

| Operation | Execute | Organize | Govern | Plan | Learn | Allocate |
|-----------|---------|----------|--------|------|-------|----------|
| **Intend** | Create task | Create role | Draft policy | Set goal | Start retro | Request budget |
| **Decompose** | Split into subtasks | Break into sub-teams | Break into provisions | Goal → milestones | Issue → root causes | Budget → line items |
| **Assign** | Assign to person | Fill role | Assign reviewer | Own a goal | Assign action item | Assign budget owner |
| **Claim** | Self-assign | Volunteer for role | — | Own a milestone | — | Claim resource |
| **Prioritize** | Set urgency | — | Set enforcement level | Rank goals | Rank findings | Rank allocations |
| **Complete** | Mark done | — | Ratify policy | Mark goal met | Close retro | Close budget period |
| **Review** | Approve work | Review performance | Audit compliance | Review progress | Peer review | Audit spend |
| **Delegate** | Assign to agent | Appoint manager | Delegate authority | Delegate ownership | — | Delegate budget |
| **Consent** | — | — | Vote on policy | — | — | Approve allocation |
| **Scope** | Define permissions | Define responsibilities | Define jurisdiction | Define success criteria | — | Define spending limits |
| **Handoff** | Transfer task | Succession | Transfer authority | Transfer ownership | — | Transfer budget |
| **Endorse** | Vouch for quality | Recommend person | Support policy | Endorse direction | Validate finding | — |
| **Block/Unblock** | Dependency | — | Veto/lift veto | Blocker/resolution | — | Freeze/unfreeze |
| **Reflect** | — | — | — | — | Post-mortem | — |
| **Assert/Challenge** | — | — | — | — | Claim/dispute finding | — |

**Every operation exists in the same grammar.** The mode determines the semantics, not the mechanism.

---

## The Entity-Mode Matrix

Each entity is primarily managed by one mode but visible to all:

| Entity | Primary mode | Secondary modes |
|--------|-------------|-----------------|
| Task | Execute | Plan (as goal child), Allocate (resource consumption) |
| Project | Plan | Execute (contains tasks), Allocate (has budget) |
| Goal | Plan | Learn (retrospective), Execute (measured by tasks) |
| Role | Organize | Govern (authority), Execute (who does work) |
| Team | Organize | Execute (group work), Plan (team goals) |
| Department | Organize | Govern (policies), Allocate (budget) |
| Policy | Govern | Organize (who it applies to), Learn (review effectiveness) |
| Process | Govern | Execute (follow it), Learn (improve it) |
| Decision | Govern | Plan (strategic), Learn (evaluate outcomes) |
| Resource | Allocate | Execute (consumed), Plan (forecasted) |
| Document | Learn | Govern (policies as docs), Plan (specs as docs) |
| Organization | Organize | Govern (charter), Allocate (top-level budget) |

---

## Views per Mode

Each mode has views optimized for its core activity:

### Execute (what we partially have)
- **Board** — kanban columns by state ✓ (have)
- **List** — sortable table ✓ (have, iter 200)
- **Detail** — task + discussion + subtasks ✓ (have)
- **Triage** — inbox-zero for incoming work (spec exists, not built)
- **Timeline** — Gantt-like dependency visualization (spec exists, not built)

### Organize (new)
- **Org Chart** — visual tree of departments → teams → roles → people
- **Directory** — searchable people/agent list with roles, teams, skills
- **Role Editor** — define roles with responsibilities, permissions, scope

### Govern (partially exists as Governance lens)
- **Policies** — list of active policies with status, jurisdiction, enforcement
- **Decisions** — decision log with rationale, alternatives, authority chain
- **Compliance** — audit trail dashboard, policy violation tracking
- **Approvals** — pending consent requests across all entities

### Plan (new)
- **Goals** — OKR tree (objective → key results → tasks)
- **Roadmap** — time-based view of projects and milestones
- **Metrics** — progress dashboards, velocity, completion rates

### Learn (partially exists as Knowledge lens)
- **Retrospectives** — structured post-mortems with action items
- **Knowledge Base** — claims, evidence, institutional memory ✓ (have as Knowledge)
- **Decisions** — shared with Govern, different lens (learning from outcomes)

### Allocate (new)
- **Budget** — resource allocation with approval chains
- **Capacity** — who has bandwidth, who's overloaded
- **Forecast** — projected resource needs based on plan

---

## Implementation Strategy

**Phase 1: Foundation** — the entity types already exist as node kinds. Most grammar ops already work. Ship the missing Execute ops (claim, review, handoff) and the Organize mode basics (roles, teams as node kinds).

**Phase 2: Govern** — policies and decisions as first-class entities. Consent already exists. Approval chains already exist via the governance lens. Connect them to work entities.

**Phase 3: Plan** — goals and OKRs. Decomposition already works. Add goal-to-task linking and progress aggregation.

**Phase 4: Learn** — retrospectives as structured reflection. Knowledge claims already exist. Add retro-specific workflow.

**Phase 5: Allocate** — resource tracking. This is the most novel mode. Budget as a node kind with approval chains.

**Key architectural principle:** No new tables needed. Every entity is a Node with a kind. Every operation is an Op. The grammar is the API. The views are lenses on the same data. This is why the architecture works at every scale — it doesn't know or care whether the node is a "task" or a "policy" or a "department." It's all events on a graph.

---

## Convergence Analysis

**Pass 1 — Need (what's absent):**
- No entity types beyond task, post, thread, conversation, comment, claim, proposal. Missing: project, goal, role, team, department, policy, process, decision, resource, document, organization.
- No mode beyond Execute. Organize, Govern (partially), Plan, Learn (partially), Allocate all absent.
- No cross-entity views (e.g., "show me all tasks owned by this team toward this goal under this policy").

**Pass 2 — Traverse (what exists that we missed):**
- Governance lens already has Propose + Vote → Govern mode has a foothold.
- Knowledge lens already has Assert + Challenge → Learn mode has a foothold.
- Spaces already function as proto-organizations.
- Membership already functions as proto-team membership.
- The dashboard already shows cross-space work → Plan mode has a foothold.

**Derive (what follows):**
- Each new entity kind requires: a create form, a detail view, and a lens/list to browse them.
- The grammar ops don't change — they already work on any node kind.
- The modes are UI organizations, not data model changes.
- Scale emerges from entity composition, not from feature flags.

**Fixpoint reached at pass 2.** The architecture already supports this. The work is in the UI modes and entity kinds, not in the foundations.

---

## What Changes vs. the Existing work-product-spec.md

The existing spec is correct but narrow. It covers Execute mode thoroughly (12 ops, 4 views, state machine, decomposition, dependencies). This general spec:

1. **Widens** — 6 modes, not 1. 12 entity types, not 1.
2. **Scales** — same primitives serve solo dev and Fortune 500.
3. **Connects** — entities reference each other (task → goal → department → policy).
4. **Emerges** — modes appear as entities are created, not as features are enabled.

The existing spec becomes the Execute mode section of this general spec. Nothing is discarded.
