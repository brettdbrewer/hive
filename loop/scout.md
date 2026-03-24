# Scout Report — Iteration 217

## Gap: Agent runtime architecture undefined, no prompts written

The hive spec defines 22 roles but doesn't specify:
1. How agents physically run (processes, machines, triggers)
2. What tools agents have access to (CLI, API, filesystem)
3. How the pipeline orchestrator works
4. The actual system prompts

**This iteration:** Define the runtime architecture, then start writing prompts for the core pipeline agents.
