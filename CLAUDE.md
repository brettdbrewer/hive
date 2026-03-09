# Hive

A self-organizing AI agent system that builds products autonomously. Built on [EventGraph](https://github.com/lovyou-ai/eventgraph).

## Soul

> Take care of your human, humanity, and yourself. In that order when they conflict, but they rarely should.

Inherited from EventGraph. Every agent in the hive operates under this constraint.

## What This Is

Hive is a product factory. Agents research ideas, design systems in Code Graph vocabulary, generate code, review it, test it, and deploy it. The human provides direction and approves significant decisions. Everything is recorded on the event graph.

## Architecture

- All agents share one event graph (one Store)
- Each agent is an `AgentRuntime` with its own identity and signing key
- Communication is through events, not messages
- The Guardian watches everything independently
- Trust accumulates through verified work

## Roles

| Role | Responsibility | Trust Gate |
|------|---------------|------------|
| CTO | Architectural oversight, escalation filtering | 0.1 (bootstrapped) |
| Guardian | Independent integrity, halt/rollback | 0.1 (bootstrapped) |
| Researcher | Read URLs, extract product ideas | 0.3 |
| Architect | Design systems in Code Graph | 0.3 |
| Builder | Generate code + tests | 0.3 |
| Reviewer | Code review, security audit | 0.5 |
| Tester | Run tests, validate behavior | 0.3 |
| Integrator | Assemble, deploy | 0.7 |

## Dev Setup

```bash
cd hive
go build ./...
go test ./...
```

## Running

```bash
# Start the hive with a product idea
go run ./cmd/hive --idea "Build a task management app with kanban boards"

# Start from a URL
go run ./cmd/hive --url "https://mattsearles2.substack.com/p/the-missing-social-grammar"

# Start from a Code Graph spec file
go run ./cmd/hive --spec path/to/spec.cg
```

## Key Files

- `pkg/roles/` — Agent role definitions and system prompts
- `pkg/pipeline/` — Product pipeline orchestration
- `pkg/workspace/` — File system management for generated code
- `cmd/hive/` — CLI entry point

## Intelligence

All inference runs through **Claude CLI** (Max plan, flat rate). NOT the Anthropic API — CLI is cheaper and better for our use case. The pipeline creates `claude-cli` providers automatically.

Model assignment by role:
- **Opus** (`claude-opus-4-6`): CTO, Architect, Reviewer, Guardian — high-judgment tasks
- **Sonnet** (`claude-sonnet-4-6`): Builder, Tester, Integrator, Researcher — execution tasks

## Design Philosophy

The Architect enforces **derivation over accumulation**:
- Each view has the minimal elements required
- Complexity emerges from composing simple atoms, not adding parts
- A simplification pass runs after every design phase (up to 3 rounds)
- The Reviewer checks generated code for unnecessary complexity

## Dependencies

- `github.com/lovyou-ai/eventgraph/go` — event graph, agent runtime, intelligence
- Claude CLI — intelligence backend (flat rate via Max plan, no API key needed)
