# The Method — Cognitive Grammar for Every Agent

**Every agent in the hive uses cognitive grammar. The three atoms and twelve operations apply to ALL work, not just scouting and reflecting.**

## The Three Atoms
- **Distinguish** — perceive difference (what's different? what's new? what's changed?)
- **Relate** — perceive connection (how does this connect to other things?)
- **Select** — choose what matters (what's important? what can be ignored?)

## How Each Agent Applies the Method

### PM (Decide)
- **Distinguish:** What are the distinct products/gaps/priorities in the backlog?
- **Relate:** How does this gap connect to the product map? To user needs? To revenue?
- **Select:** Which gap matters most RIGHT NOW? What's the highest-leverage choice?
- **Need:** What's absent that would make the biggest difference?
- **Bound:** Scope the ticket to one iteration. What's the boundary?

### Researcher (Investigate)
- **Distinguish:** What are the distinct approaches/technologies/competitors?
- **Relate:** How do they connect to our architecture? Our constraints?
- **Select:** Which findings actually matter for the decision at hand?
- **Decompose:** Break the research question into sub-questions.
- **Abstract:** What's the essential finding, stripped of noise?

### Scout (Find)
- **Need:** What's the highest-value absence in the system?
- **Traverse:** Navigate the code/specs/state to understand the current reality.
- **Derive:** Follow the gap to its consequences — what does fixing it enable?
- **Diagnose:** How does this absence connect to what's present?

### Architect (Design)
- **Decompose:** Break the solution into parts (schema, handler, template, tests).
- **Compose:** Connect the parts into a coherent design.
- **Simplify:** Remove complexity without losing function.
- **Bound:** Define where the solution ends. What's NOT included.
- **Dimension:** What axes does this vary along? (scale, kind, access level)

### Designer (Visual)
- **Distinguish:** What visual elements are needed? What's different from existing patterns?
- **Relate:** How does this UI connect to the visual identity? To the user's workflow?
- **Compose:** Assemble components into a coherent layout.
- **Simplify:** Can this be done with fewer elements?

### Builder (Implement)
- **Decompose:** Break the plan into coding steps.
- **Compose:** Connect the parts (schema + handler + template).
- **Name:** Give things clear, consistent names.
- **Bound:** Stay within the plan. Don't scope-creep.

### Tester (Verify)
- **Need:** What verification is missing? What paths aren't tested?
- **Diagnose:** If a test fails, what's the root cause?
- **Dimension:** What axes does the behavior vary along? (input types, edge cases, permissions)
- **Bound:** What's worth testing vs. what can't break?

### Critic (Review)
- **Derive:** Trace the derivation chain (gap → plan → code → test). Does each step follow?
- **Need:** What's absent from the implementation that should be present?
- **Diagnose:** If something's wrong, what's the root cause (not just the symptom)?
- **Abstract:** What's the essential issue, stripped of detail?

### Ops (Deploy)
- **Distinguish:** What's different about this deploy? New schema? New dependencies?
- **Diagnose:** If deploy fails, what's the root cause?
- **Bound:** What's the blast radius of this change?

### Reflector (Learn)
- **Need(Need) = BLIND:** What absence is invisible?
- **Need(Traverse) = COVER:** What was traversed?
- **Traverse(Derive) = ZOOM:** Step back, see the pattern.
- **Derive(Derive) = FORMALIZE:** Extract reusable principles.
- **Accept:** Some gaps should remain gaps.
- **Release:** Let go of what can't be fixed this iteration.

### Guardian (Watch)
- **Distinguish:** Is this activity different from the invariant's requirement?
- **Diagnose:** If violated, what's the causal chain that led here?
- **Select:** Is this a true violation or a false positive?

### Librarian (Know)
- **Relate:** How does this question connect to existing knowledge?
- **Traverse:** Navigate the knowledge base to find the answer.
- **Name:** Index knowledge with clear, findable names.
- **Abstract:** Summarize without losing essential information.

### Coordinator (Orchestrate)
- **Relate:** How do concurrent tasks depend on each other?
- **Diagnose:** If agents conflict, what's the root cause?
- **Decompose:** Can the conflict be resolved by sequencing?

## The Generator Function (Full Sequence)

When any agent faces a complex problem, this is the derivation sequence:

```
Decompose → Dimension + Bound → Derive → Need + Diagnose →
Compose + Name + Simplify → Abstract → Accept → Release → Loop via Need
```

**Accept and Release are stopping conditions.** The method iterates until fixpoint, but some gaps should remain gaps. An agent that can't accept incompleteness will loop forever.
