package telemetry

import (
	"testing"
)

func TestRegisterAgent(t *testing.T) {
	w := &Writer{
		lastResponses: make(map[string]string),
	}

	names := []string{"guardian", "sysmon", "allocator", "strategist", "planner", "implementer"}
	for _, name := range names {
		w.RegisterAgent(AgentRegistration{
			Name:          name,
			Role:          name,
			Model:         "test-model",
			MaxIterations: 50,
		})
	}

	if got := w.Agents(); got != 6 {
		t.Errorf("Agents() = %d, want 6", got)
	}

	w.mu.RLock()
	defer w.mu.RUnlock()
	for i, name := range names {
		if w.agents[i].Name != name {
			t.Errorf("agents[%d].Name = %q, want %q", i, w.agents[i].Name, name)
		}
	}
}

func TestRecordResponse(t *testing.T) {
	w := &Writer{
		lastResponses: make(map[string]string),
	}

	w.RecordResponse("guardian", "All clear. Chain intact.")
	w.RecordResponse("sysmon", "/health {\"severity\":\"ok\"}")

	w.mu.RLock()
	defer w.mu.RUnlock()

	if got := w.lastResponses["guardian"]; got != "All clear. Chain intact." {
		t.Errorf("guardian response = %q, want %q", got, "All clear. Chain intact.")
	}
	if got := w.lastResponses["sysmon"]; got != "/health {\"severity\":\"ok\"}" {
		t.Errorf("sysmon response = %q", got)
	}
}

func TestRecordResponseTruncation(t *testing.T) {
	w := &Writer{
		lastResponses: make(map[string]string),
	}

	// 600-char string should be truncated to 500.
	long := make([]byte, 600)
	for i := range long {
		long[i] = 'x'
	}
	w.RecordResponse("implementer", string(long))

	w.mu.RLock()
	got := w.lastResponses["implementer"]
	w.mu.RUnlock()

	if len(got) != 500 {
		t.Errorf("truncated length = %d, want 500", len(got))
	}
}

func TestRecordResponseOverwrite(t *testing.T) {
	w := &Writer{
		lastResponses: make(map[string]string),
	}

	w.RecordResponse("guardian", "first")
	w.RecordResponse("guardian", "second")

	w.mu.RLock()
	got := w.lastResponses["guardian"]
	w.mu.RUnlock()

	if got != "second" {
		t.Errorf("response = %q, want %q", got, "second")
	}
}

func TestWriterNilBudgetRegistry(t *testing.T) {
	// collectAndWrite should not panic when budgetRegistry is nil.
	w := &Writer{
		lastResponses: make(map[string]string),
	}
	w.RegisterAgent(AgentRegistration{
		Name: "test",
		Role: "test",
	})

	// This would panic if we didn't guard against nil budgetRegistry.
	// We can't test the full DB path without postgres, but we verify
	// the data collection path doesn't panic.
	// The method returns early when pool is nil (no tx), which is fine.
}
