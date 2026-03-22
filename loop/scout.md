# Scout Report — Iteration 36

## Map

Agent identity is visible in Feed, Chat, and Comments (violet badge + avatar). But the People lens and Activity lens don't distinguish agents from humans at all. Two of six lenses are identity-blind.

## Gap Type

Inconsistency (working feature not applied uniformly)

## The Gap

The People lens shows a flat table of names — agents and humans look identical. The Activity lens shows operations with actor names but no kind indicator. This undermines the "agents and humans are peers" vision — peers need to be visually distinguishable.

## Why This Gap Over Others

Agent identity was the differentiator (lesson 27). Five iterations (25-27) built identity infrastructure. But two lenses were never updated. This is low-effort, high-impact consistency work — the infrastructure exists (users table has `kind`), it's just not queried.

## What "Filled" Looks Like

People lens: agent members get violet avatar + "agent" badge (matching FeedCard and chatMessage patterns). Activity lens: ops from agents get violet avatar + "agent" pill (same pattern). Both resolved from the users table at query time via JOIN — no schema migration, no `RecordOp` signature change. Lesson 30 in action.
