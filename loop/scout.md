# Scout Report — Iteration 11

## Map (from code + docs)

Read state.md. Site is production-ready, five clusters complete. State.md recommends hive autonomy as the next cluster — the loop currently requires manual `claude -p` invocation for each phase.

CORE-LOOP.md (docs/CORE-LOOP.md lines 201-212) specifies that the loop should be runnable via prompt files (`loop/scout-prompt.txt`, `loop/builder-prompt.txt`, etc.) and a shell script. These files don't exist — the spec references them but they were never created.

## Gap Type

Missing code — the core loop's executable infrastructure doesn't exist.

## The Gap

CORE-LOOP.md describes a self-running loop with prompt files and a run.sh script, but none of these files exist. The loop runs only because a human manually types the Scout/Builder/Critic/Reflector prompts into Claude Code.

## Why This Gap

This is the most load-bearing gap because it's recursive: making the loop easier to run makes every future iteration faster. Every other improvement (product, content, skills) requires running the loop. Reducing the loop's friction compounds across all future work. Also, it's what the Reflector explicitly recommended.

## Filled Looks Like

`./loop/run.sh` runs all four phases in sequence, each invoking `claude -p` with the corresponding prompt file. Each prompt file is a complete, standalone instruction set.
