# Scout Report — Iteration 183

## Gap Identified

**The Social layer lacks message reactions** — the single most impactful missing feature across all four planned social modes (Chat, Rooms, Square, Forum).

## Evidence

- Competitive analysis: Every major platform (Slack, Discord, Twitter, Reddit, iMessage) has reactions/emoji on messages. It's table stakes.
- Current state: Messages have no reaction support. The only acknowledgment mechanism is reply or endorse (which is heavyweight).
- Impact: Reactions reduce noise (no more "thanks!", "lol", "+1" messages), provide lightweight social signal, and are the foundation for the Acknowledge grammar operation.
- The social-spec.md (Phase 1, item 1) identifies reactions as the first Chat Foundation feature.

## What to Build

1. **reactions table** — (node_id, user_id, emoji, created_at) with compound PK
2. **react/unreact grammar op** — new `react` op in the ops table
3. **UI: hover action bar on messages** — show reaction picker on hover
4. **UI: reaction badges below messages** — grouped emoji with counts, clickable to toggle
5. **Quick reactions** — most-used emoji shortcuts (thumbs up, heart, eyes, fire, checkmark)

## Why This First

Message reactions are:
- **Cross-cutting** — benefits Chat, Rooms, Square, and Forum modes
- **Foundation** — the Acknowledge grammar operation needs a UI primitive
- **High visibility** — immediately makes conversations feel alive and modern
- **Low risk** — additive, no existing behavior changes
- **Small scope** — can ship in one iteration

## Constraints

- Must work on the existing `nodes` message model (reactions reference node_id)
- Must render in both full and compact (grouped) message views
- Must update in real-time via existing HTMX polling mechanism
- Emoji set: start with 6 quick reactions, expandable later
