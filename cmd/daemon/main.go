// Command daemon runs background agents as long-running watchers.
//
// Each background agent (Guardian, Librarian, Coordinator, etc.) gets a goroutine
// that polls for relevant events and resumes the agent's Claude CLI session
// when triggered.
//
// Pipeline agents (Scout, Builder, Critic, etc.) are NOT run by the daemon —
// they're triggered by cmd/loop. The daemon handles continuous monitoring.
//
// Usage:
//
//	cd hive && go run ./cmd/daemon/ --agents guardian,librarian
//	cd hive && go run ./cmd/daemon/ --all-background
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// BackgroundAgent defines a daemon-managed agent.
type BackgroundAgent struct {
	Name         string
	PollInterval time.Duration
	PollURL      string // lovyou.ai API endpoint to poll
	Filter       func(event map[string]any) bool
}

var backgroundAgents = map[string]BackgroundAgent{
	"guardian": {
		Name:         "guardian",
		PollInterval: 30 * time.Second,
		PollURL:      "/app/hive/activity",
		Filter:       func(e map[string]any) bool { return true }, // watches everything
	},
	"librarian": {
		Name:         "librarian",
		PollInterval: 10 * time.Second,
		PollURL:      "/app/hive/conversations",
		Filter: func(e map[string]any) bool {
			// Triggered by messages in #questions or @librarian mentions
			return true
		},
	},
	"accountant": {
		Name:         "accountant",
		PollInterval: 5 * time.Minute,
		PollURL:      "/app/hive/activity",
		Filter: func(e map[string]any) bool {
			// Triggered by iteration completions
			op, _ := e["op"].(string)
			return op == "express" // iteration posts
		},
	},
	"coordinator": {
		Name:         "coordinator",
		PollInterval: 1 * time.Minute,
		PollURL:      "/app/hive/board",
		Filter: func(e map[string]any) bool {
			return true // watches task state changes
		},
	},
	"maintainer": {
		Name:         "maintainer",
		PollInterval: 10 * time.Minute,
		PollURL:      "/app/hive/activity",
		Filter: func(e map[string]any) bool {
			op, _ := e["op"].(string)
			return op == "complete" // watches for completed tasks
		},
	},
	"security": {
		Name:         "security",
		PollInterval: 5 * time.Minute,
		PollURL:      "/app/hive/activity",
		Filter: func(e map[string]any) bool {
			op, _ := e["op"].(string)
			return op == "express" || op == "respond" // watches code-related ops
		},
	},
}

func main() {
	agents := flag.String("agents", "", "Comma-separated list of background agents to run")
	allBackground := flag.Bool("all-background", false, "Run all background agents")
	apiBase := flag.String("api", "https://lovyou.ai", "lovyou.ai API base URL")
	flag.Parse()

	apiKey := os.Getenv("LOVYOU_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "LOVYOU_API_KEY required")
		os.Exit(1)
	}

	// Determine which agents to run.
	var toRun []BackgroundAgent
	if *allBackground {
		for _, a := range backgroundAgents {
			toRun = append(toRun, a)
		}
	} else if *agents != "" {
		for _, name := range strings.Split(*agents, ",") {
			name = strings.TrimSpace(name)
			if a, ok := backgroundAgents[name]; ok {
				toRun = append(toRun, a)
			} else {
				fmt.Fprintf(os.Stderr, "WARNING: unknown agent %q (valid: guardian, librarian, accountant, coordinator, maintainer, security)\n", name)
			}
		}
	} else {
		fmt.Fprintln(os.Stderr, "usage: daemon --agents guardian,librarian  OR  daemon --all-background")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "=== HIVE DAEMON ===\n")
	fmt.Fprintf(os.Stderr, "Running %d background agents:\n", len(toRun))
	for _, a := range toRun {
		fmt.Fprintf(os.Stderr, "  • %s (poll every %s)\n", a.Name, a.PollInterval)
	}
	fmt.Fprintln(os.Stderr)

	// Set up graceful shutdown.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM,
	)
	defer stop()

	// Run each agent in a goroutine.
	var wg sync.WaitGroup
	for _, agent := range toRun {
		wg.Add(1)
		go func(a BackgroundAgent) {
			defer wg.Done()
			runDaemon(ctx, a, *apiBase, apiKey)
		}(agent)
	}

	// Wait for shutdown signal.
	<-ctx.Done()
	fmt.Fprintln(os.Stderr, "\nShutting down...")
	wg.Wait()
	fmt.Fprintln(os.Stderr, "Daemon stopped.")
}

// runDaemon polls for events and triggers the agent when relevant events occur.
func runDaemon(ctx context.Context, agent BackgroundAgent, apiBase, apiKey string) {
	ticker := time.NewTicker(agent.PollInterval)
	defer ticker.Stop()

	lastSeen := time.Now()
	fmt.Fprintf(os.Stderr, "[%s] started, polling %s every %s\n", agent.Name, agent.PollURL, agent.PollInterval)

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "[%s] stopping\n", agent.Name)
			return
		case <-ticker.C:
			events := pollEvents(agent, apiBase, apiKey, lastSeen)
			if len(events) > 0 {
				fmt.Fprintf(os.Stderr, "[%s] %d new events, triggering\n", agent.Name, len(events))
				triggerAgent(agent.Name, events)
				lastSeen = time.Now()
			}
		}
	}
}

// pollEvents fetches new events from the lovyou.ai API.
func pollEvents(agent BackgroundAgent, apiBase, apiKey string, since time.Time) []map[string]any {
	url := fmt.Sprintf("%s%s?format=json&after=%s", apiBase, agent.PollURL, since.Format(time.RFC3339))
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	var result map[string]any
	if json.Unmarshal(body, &result) != nil {
		return nil
	}

	// Extract events from response (format depends on endpoint).
	// For activity: {"ops": [...]}
	// For conversations: {"conversations": [...]}
	// For board: {"nodes": [...]}
	var events []map[string]any
	for _, key := range []string{"ops", "conversations", "nodes"} {
		if items, ok := result[key].([]any); ok {
			for _, item := range items {
				if m, ok := item.(map[string]any); ok {
					if agent.Filter(m) {
						events = append(events, m)
					}
				}
			}
		}
	}
	return events
}

// triggerAgent resumes the agent's Claude CLI session with the event context.
func triggerAgent(agentName string, events []map[string]any) {
	// Build a summary of the events.
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("You have %d new events to process:\n\n", len(events)))
	for i, e := range events {
		if i >= 10 {
			sb.WriteString(fmt.Sprintf("... and %d more\n", len(events)-10))
			break
		}
		data, _ := json.Marshal(e)
		sb.WriteString(string(data))
		sb.WriteString("\n")
	}

	// Check for existing session.
	sessionID := getSessionID(agentName)

	var args []string
	if sessionID != "" {
		args = []string{"--resume", sessionID, "--print", "--message", sb.String()}
	} else {
		// First run — load full context.
		prompt := loadFullPrompt(agentName)
		args = []string{"--print", "-n", "hive-" + agentName, "--system-prompt", prompt, "--message", sb.String()}
		saveSessionID(agentName, "hive-"+agentName)
	}

	cmd := exec.Command("claude", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout // let agent output be visible

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] ERROR: %v\n", agentName, err)
	}
}

// loadFullPrompt assembles the full initial prompt for a background agent.
func loadFullPrompt(agentName string) string {
	var sb strings.Builder

	context, _ := os.ReadFile(filepath.Join("agents", "CONTEXT.md"))
	sb.WriteString(string(context))
	sb.WriteString("\n\n")

	agentPrompt, _ := os.ReadFile(filepath.Join("agents", agentName+".md"))
	sb.WriteString(string(agentPrompt))
	sb.WriteString("\n\n")

	method, _ := os.ReadFile(filepath.Join("agents", "METHOD.md"))
	sb.WriteString(string(method))

	return sb.String()
}

// Session management (shared with cmd/loop).
func sessionFile(agent string) string {
	return filepath.Join("agents", ".sessions", agent)
}

func getSessionID(agent string) string {
	data, err := os.ReadFile(sessionFile(agent))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func saveSessionID(agent, sessionID string) {
	os.MkdirAll(filepath.Join("agents", ".sessions"), 0755)
	os.WriteFile(sessionFile(agent), []byte(sessionID), 0644)
}
