# Scout Report — Iteration 225

## Gap Identified

**The runtime proved the plumbing works (iter 224). Now prove it can ship code.**

Iter 224's builder completed a design task (no code changes). The runtime has never shipped a real commit to production. Three critique issues remain:

1. **Double role prompt** — Provider created with `SystemPrompt: rolePrompt` AND runner prepends rolePrompt in instruction. Wastes tokens, doubles the context.
2. **No recency tiebreak** — When multiple tasks share priority, first from API wins. Builder should prefer newest tasks (most likely to be fresh assignments).
3. **No "changes required" check** — Builder marks DONE even when Operate produces no file changes. Should require at least one changed file for implementation tasks.

After fixes, re-run the builder on the Policy entity task (high priority, assigned, created today). This is the acid test: can the runtime ship a real commit?

## Plan

1. Fix double role prompt — remove SystemPrompt from provider config, keep role prompt only in instruction
2. Add recency tiebreak — sort by priority first, then by created_at descending (newest first)
3. Add "changes required" guard — after DONE, if no git changes detected and task is implementation, leave in-progress
4. Re-run builder on Policy entity task
5. Verify: task claimed → code changed → build passes → commit pushed → task closed

## Priority

**P0** — This completes the runtime proof. Iter 224 proved the flow. Iter 225 proves the output.
