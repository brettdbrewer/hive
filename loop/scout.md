Perfect. I have all the context I need. Let me write the Scout report.

---

# Scout Report — Iteration 309

## Gap

The pipeline role has no handler in runTick, blocking autonomous phase coordination and feedback-driven recovery.

## Evidence

**Current state:**
- `pkg/runner/runner.go:165-186` — runTick's switch statement covers scout, architect, builder, critic, pm, observer, monitor, reflector. No "pipeline" case. Default falls through to "no handler" log.
- `pkg/runner/pipeline_tree.go` — does not exist.
- `pkg/runner/diagnostic.go` — PhaseEvent and appendDiagnostic() are complete (lines 13-52). Ready to emit failure signals.
- `pkg/runner/runner.go:17` — decision package already imported but unused.
- `eventgraph/go/pkg/decision/` — tree.go, evaluate.go, evolve.go exist and are tested.

**From state.md (lines 545-577):**
> The pipeline is still a sequential for-loop. Failures don't trigger branching. Root causes don't become tasks. Until this is wired, the PM optimizes blind.

**From commit e4643be:**
> Director mandate: engine before paint — PM must prioritize foundation

**Recent work (iterations 302–308):**
- PhaseEvent type added (diagnostic.go)
- appendDiagnostic() added to log failures
- diagnostics.jsonl infrastructure in place
- PM prompt updated to read failures and suggest fixes
- But: no tree to orchestrate phases, no branching on failure, no automatic recovery

## Impact

**For users (the PM agent):**
- The PM reads failure diagnostics but has no authority to branch or assign fixes
- Feedback loop is open: diagnostics go out, but no action comes back
- The pipeline runs to completion even when early phases fail, wasting tokens on doomed work

**For product:**
- Cannot implement autonomous failure recovery (Phase 2)
- Cannot build pattern detection that learns from repeated failures (Phase 2)
- Blocks higher-order operations: fixpoint awareness, irreversibility bounds, depth-first search

**For the hive:**
- The self-evolve invariant ("agents fix agents") cannot work if the pipeline can't fix itself

## Scope

**Files to create:**
- `pkg/runner/pipeline_tree.go` — PipelineTree type, Execute(ctx), phase nodes, failure branches
- `pkg/runner/pipeline_tree_test.go` — one test: tree with stub phase, verify appendDiagnostic writes on failure

**Files to modify:**
- `pkg/runner/runner.go` — add "pipeline" case to runTick (lines ~176), call PipelineTree.Execute(ctx)

**Existing files (do NOT touch):**
- eventgraph's decision package (already complete)
- Existing role handlers (scout, architect, builder, critic stay untouched for single-role runs)
- cmd/hive
- Agent definitions

## Suggestion

Implement Phase 1 exactly as described in state.md:

1. **PipelineTree** models one cycle as a decision tree with:
   - Root: "run-cycle"
   - Phase nodes: scout, architect, builder, critic
   - Each node has success condition (`PhaseEvent.Outcome == "success"`) and failure branch
   - On phase failure: emit PhaseEvent via appendDiagnostic, create a fix task via APIClient
   - Phase implementations delegate to existing `runScout`, `runArchitect`, `runBuilder`, `runCritic`

2. **Wire runTick** to use the tree for `role == "pipeline"`. Keep single-role dispatch untouched.

3. **Test once** with a minimal tree stub that returns failure, verify appendDiagnostic writes to diagnostics.jsonl.

**Scope boundary:** Wire-up only. Don't refactor existing role handlers. Don't add evolve.go pattern detection (that's Phase 2). One test, one decision path: success → continue, failure → diagnostic + task.