# Invariant Derivation from Cognitive Grammar

The system's invariants are not arbitrary rules — they're structural consequences of the three cognitive atoms. Each atom implies requirements on any system that processes information correctly.

## The Three Atoms

| Atom | What it does | System requirement |
|------|-------------|-------------------|
| **Distinguish** | Perceive difference | Entities must be distinguishable |
| **Relate** | Perceive connection | Connections must be traceable |
| **Select** | Choose what matters | Scope must be bounded |

## Derivation

### From Distinguish (perceive difference)

If a system can't distinguish between two entities, it can't operate on them correctly. This generates:

- **IDENTITY (11)** — Entities referenced by immutable IDs, never mutable display values. Two users named "Matt" must still be distinguishable.
- **TRANSPARENT (7)** — Users know when talking to agents. Humans and agents must be distinguishable to the humans interacting with them.
- **INTEGRITY (3)** — Events are signed and hash-chained. Each event is distinguishable from forgeries.

**Test:** For any two entities in the system, can the system always tell them apart? If not, IDENTITY is violated.

### From Relate (perceive connection)

If a system can't trace how entities are connected, it can't reason about causality or dependencies. This generates:

- **CAUSALITY (2)** — Every event has declared causes. The connection between cause and effect is explicit.
- **EXPLICIT (14)** — Dependencies are declared, not inferred. If A requires B, the relationship is in the code.
- **OBSERVABLE (4)** — All operations emit events. Activity is connected to the event graph, not invisible.

**Test:** For any entity in the system, can you trace its connections to other entities? If a dependency exists but isn't declared, EXPLICIT is violated.

### From Select (choose what matters)

If a system can't bound its operations, it will consume unbounded resources or produce unbounded output. This generates:

- **BUDGET (1)** — Never exceed token budget. Resource consumption is selected and bounded.
- **BOUNDED (13)** — Every operation has defined scope. Queries return bounded results. Loops terminate.
- **RESERVE (10)** — Maintain 7-day runway. The boundary between operating and not-operating is defined.
- **MARGIN (9)** — Never work at a loss. Revenue must exceed cost (selection of sustainable operations).

**Test:** For any operation in the system, is there a defined limit? If a query could return unbounded results, BOUNDED is violated.

## Invariants Not Derived from Atoms

Some invariants come from the system's values, not its cognitive structure:

- **SELF-EVOLVE (5)** — Agents fix agents, not humans. This is an architectural decision about who does what.
- **DIGNITY (6)** — Agents are entities with rights. This is an ethical commitment.
- **CONSENT (8)** — No data use without permission. This is a legal/ethical commitment.
- **VERIFIED (12)** — No code ships without tests. This is a process discipline.

These invariants are not less important — they're just not derivable from the cognitive grammar alone. They come from the soul ("take care of your human, humanity, and yourself") and from lessons learned.

## Completeness Check

Every invariant should trace to either:
1. A cognitive atom (structural requirement), or
2. The soul (ethical requirement), or
3. A lesson learned (empirical requirement)

If an invariant can't be traced to any of these, it's arbitrary and should be questioned.
If a structural requirement exists but no invariant covers it, there's a gap.

## Application

The Critic uses this derivation to check changes systematically:
- **Distinguish** → Can entities be told apart? (IDENTITY, TRANSPARENT, INTEGRITY)
- **Relate** → Are connections traceable? (CAUSALITY, EXPLICIT, OBSERVABLE)
- **Select** → Are operations bounded? (BUDGET, BOUNDED, RESERVE, MARGIN)

This is not a checklist to memorize — it's a generative framework. When the Critic encounters a pattern it hasn't seen before, it can derive the relevant invariant from the atoms rather than waiting for a human to add it to the list.
