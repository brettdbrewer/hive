# Build: Fix loop artifact failures from iter 312 Critic review

## What changed

### 1. `loop/reflections.md`
- Removed the empty reflection entry (all four sections blank) committed in iter 312.
- Removed trailing cruft (planning notes and "Should I proceed?" text) that leaked into the file after the last FORMALIZE section.

### 2. `loop/state.md`
- Added Lessons 73–77 to the lessons list (items 69–73), which were formalized in prior reflections but never propagated to state.md, violating Lesson 72 in each case.

## Files changed
- `loop/reflections.md` — remove empty entry + cruft
- `loop/state.md` — add Lessons 73–77 to lessons list
- `loop/build.md` — this file
