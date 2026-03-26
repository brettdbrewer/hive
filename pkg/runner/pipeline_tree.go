package runner

import (
	"context"
	"fmt"
)

// Phase is a single pipeline phase that can succeed or fail.
type Phase struct {
	Name string
	Run  func(ctx context.Context) error
}

// PipelineTree orchestrates a sequence of phases, emitting diagnostics on failure.
type PipelineTree struct {
	cfg    Config
	phases []Phase
}

// NewPipelineTree creates a PipelineTree wired to r's phase implementations.
// Each phase delegates to the corresponding Runner method. Those methods do not
// return errors today; the wrappers always succeed. Real failure detection is
// Phase 2 work once the phase methods propagate errors up.
func NewPipelineTree(r *Runner) *PipelineTree {
	return &PipelineTree{
		cfg: r.cfg,
		phases: []Phase{
			{Name: "scout", Run: func(ctx context.Context) error { r.runScout(ctx); return nil }},
			{Name: "architect", Run: func(ctx context.Context) error { r.runArchitect(ctx); return nil }},
			{Name: "builder", Run: func(ctx context.Context) error { r.runBuilder(ctx); return nil }},
			{Name: "critic", Run: func(ctx context.Context) error { r.runCritic(ctx); return nil }},
		},
	}
}

// Execute runs each phase in order. On the first failure it emits a PhaseEvent
// diagnostic and returns the error; subsequent phases are skipped.
func (pt *PipelineTree) Execute(ctx context.Context) error {
	for _, phase := range pt.phases {
		if err := phase.Run(ctx); err != nil {
			_ = appendDiagnostic(pt.cfg.HiveDir, PhaseEvent{
				Phase:   phase.Name,
				Outcome: "failure",
				Error:   err.Error(),
			})
			return fmt.Errorf("phase %s failed: %w", phase.Name, err)
		}
	}
	return nil
}
