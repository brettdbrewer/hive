# Scout Report — Iteration 227

## Gap Identified

**The hive can build and review, but it can't decide what to build.** Without a Scout, a human must create and assign tasks. The Scout closes the autonomous loop:

```
Scout creates tasks → Builder implements → Critic reviews → repeat
```

This is Phase 2 item 8 from hive-runtime-spec.md. With Builder (224-225) and Critic (226) done, Scout is the final piece that enables fully autonomous operation.

## Plan

Implement `runScout` in the runner:

1. Read `state.md` from the hive repo (os.ReadFile — no LLM needed for reading)
2. Read recent git log from the target repo (what was recently shipped)
3. Read the current board (what tasks already exist)
4. Call `Reason()` (haiku, fast, cheap) with: state context + git log + board summary + "what's the next concrete, implementable gap?"
5. Parse the response for: task title, description, priority
6. Create the task on the board via API
7. Assign it to the hive agent

The Scout runs every ~8 ticks (~2 minutes). It only creates a task if the board has fewer than 3 open tasks assigned to the agent (prevents task pile-up).

## Priority

**P0** — This is the last piece for autonomous operation. Builder + Critic + Scout = self-sustaining development loop.
