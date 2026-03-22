# Build Report — Iteration 11

## What I planned

Create the five files referenced in CORE-LOOP.md: four phase prompt files and a run.sh orchestrator.

## What I built

1. **loop/scout-prompt.txt** — Complete Scout phase instructions. Reads state.md first, then explores codebase/infra. Produces MAP, EXPLORE, CATALOG, NEED. Writes scout.md.

2. **loop/builder-prompt.txt** — Complete Builder phase instructions. Reads scout.md, plans, builds, commits, pushes. Handles deploys if site changes. Writes build.md.

3. **loop/critic-prompt.txt** — Complete Critic phase instructions. Reads scout.md + build.md + git diff. TRACE + AUDIT. Writes critique.md with APPROVED/REVISE verdict.

4. **loop/reflector-prompt.txt** — Complete Reflector phase instructions. COVER/BLIND/ZOOM/FORMALIZE assessments. Appends to reflections.md, updates state.md. Commits and pushes.

5. **loop/run.sh** — Shell orchestrator. Supports `./loop/run.sh [scout|builder|critic|reflector|all]`. Uses `claude -p "$(cat prompt-file)"` for each phase. Made executable.

## What works

All five files exist and are syntactically correct. run.sh is executable. The prompts match the loop structure described in CORE-LOOP.md. Each prompt is self-contained — it tells the agent what to read, what to do, and what to write.

## Key finding

The prompt files distill what was previously implicit knowledge (carried in conversation context) into explicit, reusable instructions. This is the first step toward hive autonomy — the loop can now be triggered by a single command instead of manual prompting.
