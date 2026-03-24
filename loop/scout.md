# Scout Report — Iteration 228

## Gap Identified

**The three roles work in isolation but have never run as a pipeline.** Lesson 55: "The autonomous loop is closed but untested as a pipeline." This is Phase 2 item 11: "Test: full Scout → Builder → Critic cycle."

Currently requires three separate commands. One command should run the full cycle.

## Plan

1. Add `--pipeline` mode to cmd/hive — runs Scout → Builder → Critic in sequence on the same repo
2. Scout creates a task → Builder claims and implements → Critic reviews the commit
3. All three share the same API client, budget, and repo path
4. Pipeline exits after one full cycle (like one-shot but for all three roles)
5. Run it and verify end-to-end: task appears on board, code committed, review posted

## Priority

**P0** — Phase 2 item 11. This proves the hive can operate autonomously in a single command.
