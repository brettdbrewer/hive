package runner

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPipelineTreeFailureWritesDiagnostic(t *testing.T) {
	hiveDir := makeHiveDir(t, "# State\n", nil)

	pt := &PipelineTree{
		cfg: Config{HiveDir: hiveDir},
		phases: []Phase{
			{
				Name: "stub",
				Run: func(_ context.Context) error {
					return fmt.Errorf("injected failure")
				},
			},
		},
	}

	_ = pt.Execute(context.Background())

	path := filepath.Join(hiveDir, "loop", "diagnostics.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("diagnostics.jsonl not created: %v", err)
	}

	sc := bufio.NewScanner(strings.NewReader(string(data)))
	if !sc.Scan() {
		t.Fatal("diagnostics.jsonl is empty")
	}

	var e PhaseEvent
	if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
		t.Fatalf("invalid JSON: %v\ncontent: %s", err, sc.Bytes())
	}
	if e.Outcome != "failure" {
		t.Errorf("outcome: got %q, want %q", e.Outcome, "failure")
	}
	if e.Phase != "stub" {
		t.Errorf("phase: got %q, want %q", e.Phase, "stub")
	}
}
