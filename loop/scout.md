# Scout Report — Iteration 208

## Gap: No claim op — can't self-assign tasks

**Source:** work-product-spec.md — "claim: Sets assignee_id = actor. Records Op. open → active."

**Current state:** Tasks can be assigned by the creator, but there's no way for someone to say "I'll take this." You have to ask the space owner to assign you. This blocks the market mechanism, agent self-assignment, and team claim workflows.

**What's needed:**
1. `claim` grammar op — sets assignee = actor, state → active, records op
2. "Claim" button on task cards (Board + List) for unassigned open tasks
3. Notification to task author when someone claims their task

**Why claim:** Every scale needs it. Solo dev claims from their backlog. Team member claims from the sprint. Agent claims from the work queue. Delivery driver claims a delivery. It's the self-assignment primitive.

**Approach:** Add handler case in the op switch. Update TaskCard with a claim button (shown when task is unassigned + open + user is logged in). HTMX swap for inline state update.
