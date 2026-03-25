# Backlog — Ideas, Directions, Futures

**Not specs. Not tasks. Ideas that need to be somewhere they won't be lost.**

The council, the Director, and the agents generate ideas faster than they can be specced. This file holds them until they're ready to become specs — or until the Mourner says "let this one go."

---

## Product ideas

### Hive dashboard (spectator view)
lovyou.ai/hive — live view of what the civilization is doing. Pipeline status (Scout/Builder/Critic), current task, recent commits, cost total, play/pause. The hive as a spectator sport. Makes the civilization visible to outsiders. The Designer, Storyteller, and Growth agent all asked for this.

### Specs on the Knowledge layer
Specs should be nodes on the graph, not markdown files in a repo. The Knowledge layer already has assert/challenge/verify/retract — perfect for spec lifecycle. Scout reads verified specs, decomposes into tasks. When specs are exhausted, council generates more. Self-sustaining loop.

### Agents as contacts (standard chat UX)
Global contact list. Multiple conversations per agent. Conversation summaries. Cross-conversation search. The standard iMessage/Telegram pattern — contacts on the left, threads on the right. See agent-chat-spec.md for full details.

### Council as a product feature
User asks a question, 50 agents respond from their roles. Premium feature ($5-8 per convening). Gate behind subscription or BYOK API key. Unique differentiator — nobody else has it.

### Agents create their own tools (OpenClaw pattern)
When an agent hits a capability gap, it creates a task to build the tool. The pipeline implements it. The agent improves itself. See agent-capability-spec.md for full details.

### Agent memory across conversations
Per-persona memories stored on the graph. Agents remember you. Selective, interpretive — not a log dump. See agent-capability-spec.md.

### Hive status in the UI
The board already shows tasks. Add a "Hive" view that shows: what the pipeline is working on right now, recent autonomous commits, cost dashboard, council history. Real-time if possible.

---

## Architectural ideas

### Specs as graph events
Specs should be events on the event graph — signed, causal, attributable. When a spec is created, it links to the council or conversation that motivated it. When a task implements part of a spec, it links back. Full provenance.

### Agent pub/sub on the event graph
Agents should subscribe to event types they care about. The Critic subscribes to `hive.builder.committed`. The Guardian subscribes to `*`. The Philosopher subscribes to `council.*`. Currently: agents are invoked by the pipeline. Future: agents react to events.

### Cross-system agent identity (EGIP)
Agents should be able to participate on OTHER platforms — not just lovyou.ai. The event graph + EGIP protocol enables this. An agent's identity is its signing key, not its platform account.

### Revenue from agent conversations
The Finance agent's concern: zero revenue. Agent conversations could be the first revenue stream. Free tier: 10 agent chats/day. Paid tier: unlimited + councils. BYOK: bring your own API key. The soul says "free for individuals" — individual chats stay free, councils and enterprise features pay.

---

## Process ideas

### Automated council schedule
Council every 10 iterations, or when the Scout can't find gaps, or on Director demand. Results posted to feed + feed into state.md.

### REVISE enforcement gate
Before Scout creates new work, check for open fix tasks. Fix before build. Currently: fix tasks pile up ignored.

### Deploy-on-merge (not deploy-per-cycle)
Batch commits, deploy once. The current approach of deploying after every cycle causes Fly machine collisions. Accumulate commits, deploy on a schedule or trigger.

### Reflector in the pipeline
The Reflector role exists but doesn't run. It should close every pipeline cycle: read what happened, update state.md, append to reflections.md. Currently: Claude Code (me) does this manually.

---

*This file is append-only. Ideas move to specs when they're ready. The Mourner reviews periodically and releases what's no longer relevant.*
