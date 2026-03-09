// Package pipeline orchestrates the product build pipeline.
package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/intelligence"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	"github.com/lovyou-ai/hive/pkg/roles"
	"github.com/lovyou-ai/hive/pkg/workspace"
)

// Phase represents a stage in the product pipeline.
type Phase string

const (
	PhaseResearch  Phase = "research"
	PhaseDesign    Phase = "design"
	PhaseBuild     Phase = "build"
	PhaseReview    Phase = "review"
	PhaseTest      Phase = "test"
	PhaseIntegrate Phase = "integrate"
)

// ProductInput describes how a product idea enters the hive.
type ProductInput struct {
	Name        string // Product name (used for repo and directory)
	URL         string // Read from URL (Substack post, docs, etc.)
	Description string // Natural language description
	SpecFile    string // Path to a Code Graph spec file
}

// Pipeline orchestrates agents through the product build phases.
type Pipeline struct {
	store   store.Store
	ws      *workspace.Workspace
	product *workspace.Product // current product being built

	cto      *roles.Agent
	guardian *roles.Agent
	agents   map[roles.Role]*roles.Agent
}

// Config for creating a new pipeline.
type Config struct {
	Store   store.Store
	WorkDir string // Root directory for generated products
}

// New creates a pipeline and bootstraps the CTO and Guardian.
func New(ctx context.Context, cfg Config) (*Pipeline, error) {
	ws, err := workspace.New(cfg.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("workspace: %w", err)
	}

	p := &Pipeline{
		store:  cfg.Store,
		ws:     ws,
		agents: make(map[roles.Role]*roles.Agent),
	}

	// Bootstrap CTO first — architectural oversight (Opus)
	cto, err := p.ensureAgent(ctx, roles.RoleCTO, "cto")
	if err != nil {
		return nil, fmt.Errorf("bootstrap CTO: %w", err)
	}
	p.cto = cto

	// Bootstrap Guardian — independent integrity monitor (Opus)
	guardian, err := p.ensureAgent(ctx, roles.RoleGuardian, "guardian")
	if err != nil {
		return nil, fmt.Errorf("bootstrap Guardian: %w", err)
	}
	p.guardian = guardian

	return p, nil
}

// providerForRole creates an intelligence provider with the model appropriate for the role.
// Uses Claude CLI (flat rate via Max plan) — no API key needed.
func (p *Pipeline) providerForRole(role roles.Role) (intelligence.Provider, error) {
	model := roles.PreferredModel(role)
	return intelligence.New(intelligence.Config{
		Provider: "claude-cli",
		Model:    model,
	})
}

// ensureAgent creates an agent of the given role if it doesn't exist yet.
// Each role gets a provider with the appropriate model (Opus for judgment, Sonnet for execution).
func (p *Pipeline) ensureAgent(ctx context.Context, role roles.Role, name string) (*roles.Agent, error) {
	if agent, ok := p.agents[role]; ok {
		return agent, nil
	}
	provider, err := p.providerForRole(role)
	if err != nil {
		return nil, fmt.Errorf("provider for %s: %w", role, err)
	}
	agent, err := roles.NewAgent(ctx, roles.AgentConfig{
		Role:     role,
		Name:     name,
		Store:    p.store,
		Provider: provider,
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("  ↳ %s agent using %s\n", role, roles.PreferredModel(role))
	p.agents[role] = agent
	return agent, nil
}

// Run executes the full product pipeline for a given input.
func (p *Pipeline) Run(ctx context.Context, input ProductInput) error {
	// Initialize product repo
	name := input.Name
	if name == "" {
		name = "product"
	}
	product, err := p.ws.InitProduct(name)
	if err != nil {
		return fmt.Errorf("init product: %w", err)
	}
	p.product = product
	fmt.Printf("Product repo: %s → %s\n", product.Dir, product.Repo)

	fmt.Println("═══ Phase 1: Research ═══")
	spec, err := p.research(ctx, input)
	if err != nil {
		return fmt.Errorf("research: %w", err)
	}

	fmt.Println("═══ Phase 2: Design ═══")
	design, err := p.design(ctx, spec)
	if err != nil {
		return fmt.Errorf("design: %w", err)
	}

	fmt.Println("═══ Phase 2b: Simplify ═══")
	design, err = p.simplify(ctx, design)
	if err != nil {
		return fmt.Errorf("simplify: %w", err)
	}

	fmt.Println("═══ Phase 3: Build ═══")
	code, err := p.build(ctx, design)
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	fmt.Println("═══ Phase 4: Review ═══")
	err = p.review(ctx, code, design)
	if err != nil {
		return fmt.Errorf("review: %w", err)
	}

	fmt.Println("═══ Phase 5: Test ═══")
	err = p.test(ctx, code)
	if err != nil {
		return fmt.Errorf("test: %w", err)
	}

	fmt.Println("═══ Phase 6: Integrate ═══")
	err = p.integrate(ctx)
	if err != nil {
		return fmt.Errorf("integrate: %w", err)
	}

	fmt.Println("═══ Pipeline Complete ═══")
	return nil
}

// research gathers information about the product idea.
func (p *Pipeline) research(ctx context.Context, input ProductInput) (string, error) {
	var spec string

	if input.SpecFile != "" {
		// Read the spec file directly
		content, err := p.ws.ReadFile(input.SpecFile)
		if err != nil {
			return "", fmt.Errorf("read spec: %w", err)
		}
		spec = content
	} else {
		researcher, err := p.ensureAgent(ctx, roles.RoleResearcher, "researcher")
		if err != nil {
			return "", err
		}

		if input.URL != "" {
			_, evaluation, err := researcher.Runtime.Research(ctx, input.URL,
				"extract the product idea, key entities, features, and requirements. Output in Code Graph vocabulary where possible.")
			if err != nil {
				return "", fmt.Errorf("research URL: %w", err)
			}
			spec = evaluation
		} else if input.Description != "" {
			_, evaluation, err := researcher.Runtime.Evaluate(ctx, "product_idea", input.Description)
			if err != nil {
				return "", fmt.Errorf("evaluate idea: %w", err)
			}
			spec = evaluation
		}
	}

	// CTO evaluates feasibility
	_, ctoEval, err := p.cto.Runtime.Evaluate(ctx, "feasibility",
		fmt.Sprintf("Evaluate this product idea for feasibility. What agents are needed? What's the build sequence? Key risks?\n\n%s", spec))
	if err != nil {
		return "", fmt.Errorf("CTO evaluate: %w", err)
	}

	fmt.Printf("CTO Assessment:\n%s\n", ctoEval)
	return spec, nil
}

// design creates a full Code Graph spec from the product idea.
func (p *Pipeline) design(ctx context.Context, spec string) (string, error) {
	architect, err := p.ensureAgent(ctx, roles.RoleArchitect, "architect")
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf("%s\n\nDesign the full system architecture. Output a complete Code Graph spec. Remember: derive complexity from simple compositions. Each view should have the minimal elements needed — if a view feels heavy, decompose it. Elegant, simple, beautiful.\n\n%s",
		roles.SystemPrompt(roles.RoleArchitect), spec)

	_, design, err := architect.Runtime.Evaluate(ctx, "architecture", prompt)
	if err != nil {
		return "", fmt.Errorf("architect design: %w", err)
	}

	// CTO reviews the architecture — check for derivation and minimalism
	_, review, err := p.cto.Runtime.Evaluate(ctx, "architecture_review",
		fmt.Sprintf("Review this architecture. Check: Are views minimal? Is complexity derived from composition rather than accumulated? Are there any bloated entities or views that should be decomposed? Is it elegant and simple?\n\n%s", design))
	if err != nil {
		return "", fmt.Errorf("CTO review design: %w", err)
	}

	fmt.Printf("Architecture Review:\n%s\n", review)
	return design, nil
}

// simplify reviews the Code Graph spec and reduces it to its minimal form.
// The Architect re-examines every element: can views be composed from fewer parts?
// Can entities be split or merged? Can states be derived rather than declared?
// This runs in a loop until the Architect reports no further simplifications.
func (p *Pipeline) simplify(ctx context.Context, design string) (string, error) {
	architect, err := p.ensureAgent(ctx, roles.RoleArchitect, "architect")
	if err != nil {
		return "", err
	}

	const maxRounds = 3
	current := design

	for round := 1; round <= maxRounds; round++ {
		_, analysis, err := architect.Runtime.Evaluate(ctx, "simplify",
			fmt.Sprintf(`Review this Code Graph spec for simplification opportunities.

For each View: can it be composed from fewer elements? Are any elements redundant or derivable from others?
For each Entity: is it as small as possible? Should it be split or can properties be derived?
For each State machine: are there too many states? Can transitions be reduced?
For each Layout: does it have too many children? Can sub-views be composed instead?

If you find simplifications, output the REVISED spec with the changes applied.
If the spec is already minimal, respond with exactly: MINIMAL

Current spec:
%s`, current))
		if err != nil {
			return "", fmt.Errorf("simplify round %d: %w", round, err)
		}

		upper := strings.ToUpper(strings.TrimSpace(analysis))
		if upper == "MINIMAL" || strings.HasPrefix(upper, "MINIMAL") {
			fmt.Printf("Simplification complete after %d round(s) — spec is minimal.\n", round)
			return current, nil
		}

		fmt.Printf("Simplification round %d applied.\n", round)
		current = analysis
	}

	fmt.Printf("Simplification capped at %d rounds.\n", maxRounds)
	return current, nil
}

// build generates code from the design spec.
func (p *Pipeline) build(ctx context.Context, design string) (string, error) {
	builder, err := p.ensureAgent(ctx, roles.RoleBuilder, "builder")
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf("%s\n\nGenerate production-quality code from this specification. Include tests.\n\n%s",
		roles.SystemPrompt(roles.RoleBuilder), design)

	code, err := builder.Runtime.CodeWrite(ctx, prompt, "go")
	if err != nil {
		return "", fmt.Errorf("builder code: %w", err)
	}

	// Write code to product repo and commit
	if err := p.product.WriteFile("main.go", code); err != nil {
		return "", fmt.Errorf("write code: %w", err)
	}
	if err := p.product.Commit("feat: initial code generation from spec"); err != nil {
		return "", fmt.Errorf("commit code: %w", err)
	}

	fmt.Printf("Code generated and committed: %d bytes\n", len(code))
	return code, nil
}

// review checks code quality and spec compliance.
func (p *Pipeline) review(ctx context.Context, code string, design string) error {
	reviewer, err := p.ensureAgent(ctx, roles.RoleReviewer, "reviewer")
	if err != nil {
		return err
	}

	reviewEvt, review, err := reviewer.Runtime.CodeReview(ctx, code, "go")
	if err != nil {
		return fmt.Errorf("code review: %w", err)
	}

	// Check spec compliance
	_, specReview, err := reviewer.Runtime.Evaluate(ctx, "spec_compliance",
		fmt.Sprintf("%s\n\nDoes this code match the design spec?\n\nDesign:\n%s\n\nCode:\n%s",
			roles.SystemPrompt(roles.RoleReviewer), design, code))
	if err != nil {
		return fmt.Errorf("spec review: %w", err)
	}

	// Check for unnecessary complexity in the implementation
	_, simplicityReview, err := reviewer.Runtime.Evaluate(ctx, "simplicity_check",
		fmt.Sprintf(`Review this code for unnecessary complexity. Check:
- Are there components that could be derived from simpler compositions?
- Are there redundant abstractions or over-engineered patterns?
- Could any part be simpler while preserving the same behavior?
- Does the code match the minimal spec, or did the builder add extras?

Code:
%s`, code))
	if err != nil {
		return fmt.Errorf("simplicity review: %w", err)
	}

	fmt.Printf("Code Review:\n%s\n\nSpec Compliance:\n%s\n\nSimplicity:\n%s\n", review, specReview, simplicityReview)
	_ = reviewEvt
	return nil
}

// test runs tests and validates behavior.
func (p *Pipeline) test(ctx context.Context, code string) error {
	tester, err := p.ensureAgent(ctx, roles.RoleTester, "tester")
	if err != nil {
		return err
	}

	_, testEval, err := tester.Runtime.Evaluate(ctx, "test_analysis",
		fmt.Sprintf("%s\n\nAnalyze this code. What tests exist? What gaps are there? Write additional integration tests if needed.\n\n%s",
			roles.SystemPrompt(roles.RoleTester), code))
	if err != nil {
		return fmt.Errorf("test analysis: %w", err)
	}

	fmt.Printf("Test Analysis:\n%s\n", testEval)
	return nil
}

// integrate assembles and prepares for deployment.
func (p *Pipeline) integrate(ctx context.Context) error {
	integrator, err := p.ensureAgent(ctx, roles.RoleIntegrator, "integrator")
	if err != nil {
		return err
	}

	_, err = integrator.Runtime.Act(ctx, "integrate", "staging")
	if err != nil {
		return fmt.Errorf("integration: %w", err)
	}

	// Push to GitHub
	if err := p.product.Push(); err != nil {
		fmt.Printf("Push failed (may need manual push): %v\n", err)
	} else {
		fmt.Printf("Pushed to https://github.com/%s\n", p.product.Repo)
	}

	// Escalate to human for production approval
	humanID := types.MustActorID("actor_human_matt")
	_, err = integrator.Runtime.Escalate(ctx, humanID, "Product ready for human review before production deploy")
	if err != nil {
		return fmt.Errorf("escalate: %w", err)
	}

	fmt.Println("Product assembled and ready for human review.")
	return nil
}

// GuardianWatch runs the Guardian's monitoring loop.
// It checks recent events for policy violations.
func (p *Pipeline) GuardianWatch(ctx context.Context) error {
	events, err := p.guardian.Runtime.Memory(20)
	if err != nil {
		return fmt.Errorf("guardian memory: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	// Build summary of recent activity for the Guardian to evaluate
	var summary string
	for _, ev := range events {
		summary += fmt.Sprintf("[%s] %s: %s\n", ev.Type().Value(), ev.Source().Value(), ev.ID().Value())
	}

	_, eval, err := p.guardian.Runtime.Evaluate(ctx, "integrity_check",
		fmt.Sprintf("%s\n\nReview these recent events for policy violations, trust anomalies, or authority overreach:\n\n%s",
			roles.SystemPrompt(roles.RoleGuardian), summary))
	if err != nil {
		return fmt.Errorf("guardian evaluate: %w", err)
	}

	fmt.Printf("Guardian Report:\n%s\n", eval)

	// If the Guardian detects issues, emit an alert
	if containsAlert(eval) {
		_, err = p.guardian.Runtime.Emit(event.AgentEscalatedContent{
			AgentID:   p.guardian.Runtime.ID(),
			Authority: types.MustActorID("actor_human_matt"),
			Reason:    eval,
		})
		if err != nil {
			return fmt.Errorf("guardian alert: %w", err)
		}
	}

	return nil
}

// containsAlert checks if the Guardian's evaluation contains an alert keyword.
func containsAlert(eval string) bool {
	upper := strings.ToUpper(eval)
	for _, keyword := range []string{"HALT", "ALERT", "VIOLATION", "QUARANTINE"} {
		if strings.Contains(upper, keyword) {
			return true
		}
	}
	return false
}

// Store returns the shared event graph.
func (p *Pipeline) Store() store.Store {
	return p.store
}

// Agents returns all active agents.
func (p *Pipeline) Agents() map[roles.Role]*roles.Agent {
	return p.agents
}
