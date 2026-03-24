package runner

import "testing"

func TestParseScoutTask(t *testing.T) {
	input := `Based on the state, the next gap is adding the Decision entity kind.

TASK_TITLE: Add Decision entity kind to the site
TASK_PRIORITY: high
TASK_DESCRIPTION: Add KindDecision constant to store.go, handleDecisions handler to handlers.go, DecisionsView template to views.templ, sidebar+mobile nav entries, and add "decision" to the intend allowlist. Follow the entity pipeline pattern.`

	title, desc, priority := parseScoutTask(input)

	if title != "Add Decision entity kind to the site" {
		t.Errorf("title = %q", title)
	}
	if priority != "high" {
		t.Errorf("priority = %q", priority)
	}
	if desc == "" {
		t.Error("description is empty")
	}
	if !searchString(desc, "KindDecision") {
		t.Error("description should mention KindDecision")
	}
}

func TestParseScoutTaskDefaults(t *testing.T) {
	// Missing priority should default to medium.
	input := "TASK_TITLE: Fix something\nTASK_DESCRIPTION: Fix the thing."
	title, desc, priority := parseScoutTask(input)

	if title != "Fix something" {
		t.Errorf("title = %q", title)
	}
	if priority != "medium" {
		t.Errorf("priority = %q, want medium", priority)
	}
	if desc != "Fix the thing." {
		t.Errorf("desc = %q", desc)
	}
}

func TestParseScoutTaskEmpty(t *testing.T) {
	title, _, _ := parseScoutTask("No structured output here.")
	if title != "" {
		t.Errorf("expected empty title, got %q", title)
	}
}

func TestBuildScoutPrompt(t *testing.T) {
	prompt := buildScoutPrompt("state content", "git log", "board summary")

	if !searchString(prompt, "state content") {
		t.Error("prompt missing state")
	}
	if !searchString(prompt, "git log") {
		t.Error("prompt missing git log")
	}
	if !searchString(prompt, "board summary") {
		t.Error("prompt missing board")
	}
	if !searchString(prompt, "TASK_TITLE") {
		t.Error("prompt missing output format")
	}
}
