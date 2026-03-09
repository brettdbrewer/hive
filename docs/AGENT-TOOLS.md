# Agent Tools & Autonomy

How agents get tools, use them, and act autonomously.

## Two Layers

### Layer 1: MCP Server (the hands)

An MCP (Model Context Protocol) server written in Go that exposes the event graph, actor store, and workspace as tools Claude CLI can call during reasoning.

When an agent reasons, it can call these tools mid-thought — Claude CLI handles the tool-call loop internally (call tool → get result → continue reasoning → call another tool → ...).

**Transport:** stdio (the MCP server is a Go binary, Claude CLI spawns it as a subprocess).

**Tools exposed:**

| Tool | Description | Reads/Writes |
|------|-------------|-------------|
| `query_events` | Query events by type, source, time range, limit | Read |
| `get_event` | Get a single event by ID with full detail | Read |
| `get_actor` | Look up an actor by ID | Read |
| `list_actors` | List actors with filters (type, status) | Read |
| `get_trust` | Get trust score between two actors | Read |
| `emit_event` | Record an event on the graph | Write |
| `read_file` | Read a file from the workspace | Read |
| `write_file` | Write a file to the workspace (stages in git) | Write |
| `list_files` | List files in a product directory | Read |
| `run_command` | Execute a shell command (sandboxed) | Write |
| `query_self` | Get own actor info, trust level, authority scope | Read |
| `query_human` | Get human operator info and preferences | Read |

Write tools require appropriate authority. The Guardian monitors all tool use.

**Architecture:**

```
Pipeline (Go)
  │
  ├── spawns Claude CLI with MCP config
  │     │
  │     ├── Claude CLI spawns MCP server (Go binary)
  │     │     │
  │     │     ├── query_events → store.Query()
  │     │     ├── get_actor → actors.Get()
  │     │     ├── emit_event → runtime.Emit()
  │     │     ├── read_file → workspace.ReadFile()
  │     │     └── ... other tools
  │     │
  │     └── Claude reasons, calls tools, gets results, continues
  │
  └── observes events on the graph
```

The MCP server and the pipeline share the same Store and IActorStore instances. The MCP server is a thin adapter — it translates MCP JSON-RPC calls into Go method calls on the store/actors/workspace.

### Layer 2: Agentic Loop (the brain)

The outer loop that gives agents sustained autonomy. An agent doesn't just respond to a prompt — it observes the world, decides what to do, acts, and observes again.

```
┌─────────────────────────────────────┐
│           AGENTIC LOOP              │
│                                     │
│  1. OBSERVE                         │
│     Query graph for new events      │
│     Check pending tasks             │
│     Read recent changes             │
│                                     │
│  2. REASON                          │
│     What needs doing?               │
│     What can I do with my tools?    │
│     What's beyond my authority?     │
│                                     │
│  3. ACT                             │
│     Call tools (MCP layer)          │
│     Emit events                     │
│     Write code                      │
│     Build new tools if needed       │
│                                     │
│  4. REFLECT                         │
│     Did it work?                    │
│     What changed on the graph?      │
│     Should I continue or escalate?  │
│                                     │
│  5. REPEAT or STOP                  │
│     Continue if work remains        │
│     Stop if quiescent               │
│     Escalate if uncertain           │
└─────────────────────────────────────┘
```

The AgentRuntime already has `RunTask` (Observe → Evaluate → Decide → Act → Learn) as a single pass. The agentic loop runs this repeatedly until:
- The task is complete (quiescence — no new events, nothing to do)
- The agent needs human approval (escalation)
- The Guardian halts the agent
- A budget/iteration limit is reached

**Key difference from the current pipeline:** Right now the pipeline orchestrates agents in a fixed sequence (research → design → build → ...). With the agentic loop, agents are self-directing — they observe the graph, identify what needs doing, and do it. The pipeline becomes a seed that kicks off work, then the agents take over.

## How They Work Together

1. Pipeline starts, registers agents, seeds initial work on the graph
2. CTO agent enters its agentic loop:
   - OBSERVE: reads the product idea from the graph
   - REASON: "I need to evaluate feasibility"
   - ACT: uses `query_events` to check prior work, emits feasibility assessment
   - REFLECT: "Now I need an architect"
   - ACT: emits a task for the Architect, escalates to human for agent spawn if needed
3. Architect agent enters its loop:
   - OBSERVE: picks up the task from the graph
   - REASON: "I need to design a Code Graph spec"
   - ACT: uses `read_file` to check existing specs, emits design
   - REFLECT: "Is this minimal? Let me simplify"
   - ACT: revises the design
4. Builder agent enters its loop:
   - OBSERVE: picks up the approved design
   - ACT: uses `write_file` to generate code, `run_command` to test
   - REFLECT: "Tests failing" → fixes → retests
5. Guardian watches all events continuously, halts if needed

Each agent runs its own loop. They communicate through events on the shared graph, not through direct messages.

## Self-Improvement

When an agent lacks a tool or skill:

1. Agent identifies the gap during REASON: "I need to deploy to fly.io but I don't have a deploy tool"
2. Agent emits a task: "Build a fly.io deployment tool"
3. CTO evaluates: is this a self-modification (changes to lovyou-ai/hive) or a new tool?
4. If self-mod: agent specs the change, submits PR, human approves
5. If new tool: agent builds it as an MCP tool extension, Guardian reviews
6. The new tool becomes available to all agents

This is how the hive grows its own capabilities. The MCP server's tool list isn't static — agents can extend it.

## Implementation Plan

### Phase 1: MCP Server (Tier 1.5)

Build the MCP server as a Go binary in `cmd/mcp-server/`:

```
cmd/mcp-server/
├── main.go          — stdio transport, JSON-RPC handler
├── tools.go         — tool definitions (name, description, schema)
└── handlers.go      — tool implementations (calls to store/actors/workspace)
```

Wire into the claude-cli provider:
- Add `MCPConfig` field to intelligence.Config
- Write `.mcp.json` before spawning Claude CLI
- Pass `--mcp-config` or let Claude CLI discover from `.mcp.json`

### Phase 2: Context Injection (immediate, before MCP)

Before MCP is built, improve context injection — before each prompt, inject:
- Recent events (last 20, already done via Memory)
- Actor list (who exists, their roles and trust levels)
- Pending tasks (what needs doing)
- Own identity (who am I, what's my authority scope)

This gives agents awareness without tool-use. MCP replaces and extends this.

### Phase 3: Agentic Loop (Tier 2)

Extend AgentRuntime with a `Loop` method:

```go
func (r *AgentRuntime) Loop(ctx context.Context, opts LoopConfig) error {
    for i := 0; i < opts.MaxIterations; i++ {
        // 1. Observe
        events, _ := r.Memory(opts.ObserveWindow)

        // 2. Reason + Act (Claude handles tool calls via MCP)
        _, response, err := r.Evaluate(ctx, "loop_iteration",
            buildLoopPrompt(events, opts))

        // 3. Reflect — check if work is done
        if isQuiescent(response) || shouldEscalate(response) {
            break
        }
    }
    return nil
}
```

The key insight: Claude CLI + MCP handles steps 2-3 internally (reason, call tools, get results, continue). The Go loop handles the outer cycle (observe new events, kick off next iteration, check stopping conditions).

### Phase 4: Self-Extending Tools (Tier 4)

Agents can add new MCP tools:
1. Agent writes a new tool handler in Go
2. Submits as PR to lovyou-ai/hive (or a plugins repo)
3. Guardian reviews for safety
4. Human approves
5. MCP server reloads with new tools

## Security

- **Read tools** require actor authentication (the agent's ActorID)
- **Write tools** require authority checks (is this agent authorized to emit this event type?)
- **run_command** is sandboxed — limited to the product workspace, no network access without approval
- **Self-modification tools** always require human approval (Required authority level)
- **Guardian monitors all tool calls** via events on the graph
- **Budget limits** prevent runaway tool use (max iterations, max tokens, max cost per loop)
