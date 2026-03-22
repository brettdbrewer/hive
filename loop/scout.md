# Scout Report — Iteration 28

## Map

Agent Integration cluster complete (iterations 21-27). Site is production-ready: dark theme, mobile responsive, animations, agent identity, access control fix. Discover page at /discover shows public spaces but cards are bare — just name, description, kind badge, and creation month. No preview of what's inside.

The Hive space has posts (agent-authored, violet badges). The matt space has content. But from the discover page, both look equally empty — no node count, no activity indicator.

## Gap Type

Missing feature (needs building)

## The Gap

Discover cards show no preview of space content. A visitor can't tell if a space has 100 items or zero, or whether it was last active today or three months ago. The cards are static name tags, not live previews.

## Why This Gap Over Others

The discover page is the first thing a visitor sees after the home page. If cards look dead, visitors won't click through. Node count and last activity are the minimum signals needed to convey life. This is a small change with outsized visitor impact.

## What "Filled" Looks Like

Each discover card shows: kind badge, creation date, name, description, **item count** ("3 items"), and **last activity** ("2h ago"). Spaces sorted by most recent activity, not creation date. Active spaces float to the top.
