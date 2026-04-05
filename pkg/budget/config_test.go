package budget

import (
	"testing"
)

// --- envInt tests ---

func TestEnvInt_UnsetReturnsFallback(t *testing.T) {
	// No env var set — should return fallback.
	got := envInt("ALLOCATOR_TEST_UNSET_INT", 42)
	if got != 42 {
		t.Errorf("envInt(unset) = %d, want 42", got)
	}
}

func TestEnvInt_ValidOverride(t *testing.T) {
	t.Setenv("ALLOCATOR_TEST_INT", "99")
	got := envInt("ALLOCATOR_TEST_INT", 1)
	if got != 99 {
		t.Errorf("envInt(99) = %d, want 99", got)
	}
}

func TestEnvInt_NegativeValue(t *testing.T) {
	t.Setenv("ALLOCATOR_TEST_NEG", "-5")
	got := envInt("ALLOCATOR_TEST_NEG", 10)
	if got != -5 {
		t.Errorf("envInt(-5) = %d, want -5", got)
	}
}

func TestEnvInt_ZeroValue(t *testing.T) {
	t.Setenv("ALLOCATOR_TEST_ZERO", "0")
	got := envInt("ALLOCATOR_TEST_ZERO", 50)
	if got != 0 {
		t.Errorf("envInt(0) = %d, want 0", got)
	}
}

func TestEnvInt_InvalidReturnsFallback(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"non-numeric", "abc"},
		{"float", "3.14"},
		{"whitespace", " 10 "},
		{"empty-after-set", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ALLOCATOR_TEST_BAD_INT", tt.value)
			got := envInt("ALLOCATOR_TEST_BAD_INT", 77)
			if got != 77 {
				t.Errorf("envInt(%q) = %d, want fallback 77", tt.value, got)
			}
		})
	}
}

// --- envFloat tests ---

func TestEnvFloat_UnsetReturnsFallback(t *testing.T) {
	got := envFloat("ALLOCATOR_TEST_UNSET_FLOAT", 3.14)
	if got != 3.14 {
		t.Errorf("envFloat(unset) = %f, want 3.14", got)
	}
}

func TestEnvFloat_ValidOverride(t *testing.T) {
	t.Setenv("ALLOCATOR_TEST_FLOAT", "55.5")
	got := envFloat("ALLOCATOR_TEST_FLOAT", 1.0)
	if got != 55.5 {
		t.Errorf("envFloat(55.5) = %f, want 55.5", got)
	}
}

func TestEnvFloat_IntegerString(t *testing.T) {
	// Integer strings should parse as valid floats.
	t.Setenv("ALLOCATOR_TEST_FLOAT_INT", "100")
	got := envFloat("ALLOCATOR_TEST_FLOAT_INT", 1.0)
	if got != 100.0 {
		t.Errorf("envFloat(100) = %f, want 100.0", got)
	}
}

func TestEnvFloat_NegativeValue(t *testing.T) {
	t.Setenv("ALLOCATOR_TEST_NEG_FLOAT", "-2.5")
	got := envFloat("ALLOCATOR_TEST_NEG_FLOAT", 10.0)
	if got != -2.5 {
		t.Errorf("envFloat(-2.5) = %f, want -2.5", got)
	}
}

func TestEnvFloat_ZeroValue(t *testing.T) {
	t.Setenv("ALLOCATOR_TEST_ZERO_FLOAT", "0")
	got := envFloat("ALLOCATOR_TEST_ZERO_FLOAT", 50.0)
	if got != 0.0 {
		t.Errorf("envFloat(0) = %f, want 0.0", got)
	}
}

func TestEnvFloat_InvalidReturnsFallback(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"non-numeric", "xyz"},
		{"whitespace", " 10.5 "},
		{"empty-after-set", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ALLOCATOR_TEST_BAD_FLOAT", tt.value)
			got := envFloat("ALLOCATOR_TEST_BAD_FLOAT", 88.8)
			if got != 88.8 {
				t.Errorf("envFloat(%q) = %f, want fallback 88.8", tt.value, got)
			}
		})
	}
}

// --- LoadConfig partial override tests ---

func TestLoadConfig_PartialOverride(t *testing.T) {
	// Set only a few env vars — the rest should be defaults.
	t.Setenv("ALLOCATOR_BUDGET_FLOOR", "50")
	t.Setenv("ALLOCATOR_DAILY_CAP_WARNING_PCT", "75")

	cfg := LoadConfig()
	def := DefaultConfig()

	if cfg.BudgetFloor != 50 {
		t.Errorf("BudgetFloor = %d, want 50", cfg.BudgetFloor)
	}
	if cfg.DailyCapWarningPct != 75.0 {
		t.Errorf("DailyCapWarningPct = %f, want 75.0", cfg.DailyCapWarningPct)
	}
	// Unset fields should remain at defaults.
	if cfg.StabilizationWindow != def.StabilizationWindow {
		t.Errorf("StabilizationWindow = %d, want default %d", cfg.StabilizationWindow, def.StabilizationWindow)
	}
	if cfg.AgentCooldown != def.AgentCooldown {
		t.Errorf("AgentCooldown = %d, want default %d", cfg.AgentCooldown, def.AgentCooldown)
	}
	if cfg.BudgetCeiling != def.BudgetCeiling {
		t.Errorf("BudgetCeiling = %d, want default %d", cfg.BudgetCeiling, def.BudgetCeiling)
	}
	if cfg.ConcentrationPct != def.ConcentrationPct {
		t.Errorf("ConcentrationPct = %f, want default %f", cfg.ConcentrationPct, def.ConcentrationPct)
	}
}

func TestLoadConfig_EachEnvVarIndependently(t *testing.T) {
	// Table-driven: set one env var at a time, verify only that field changes.
	tests := []struct {
		envKey   string
		envValue string
		check    func(Config) bool
		desc     string
	}{
		{"ALLOCATOR_STABILIZATION_WINDOW", "25", func(c Config) bool { return c.StabilizationWindow == 25 }, "StabilizationWindow=25"},
		{"ALLOCATOR_AGENT_COOLDOWN", "30", func(c Config) bool { return c.AgentCooldown == 30 }, "AgentCooldown=30"},
		{"ALLOCATOR_GLOBAL_COOLDOWN", "12", func(c Config) bool { return c.GlobalCooldown == 12 }, "GlobalCooldown=12"},
		{"ALLOCATOR_BUDGET_FLOOR", "5", func(c Config) bool { return c.BudgetFloor == 5 }, "BudgetFloor=5"},
		{"ALLOCATOR_BUDGET_CEILING", "999", func(c Config) bool { return c.BudgetCeiling == 999 }, "BudgetCeiling=999"},
		{"ALLOCATOR_INITIAL_SPAWN_BUDGET", "100", func(c Config) bool { return c.InitialSpawnBudget == 100 }, "InitialSpawnBudget=100"},
		{"ALLOCATOR_CONCENTRATION_PCT", "60.5", func(c Config) bool { return c.ConcentrationPct == 60.5 }, "ConcentrationPct=60.5"},
		{"ALLOCATOR_EXHAUSTION_WARNING_PCT", "70", func(c Config) bool { return c.ExhaustionWarningPct == 70.0 }, "ExhaustionWarningPct=70"},
		{"ALLOCATOR_IDLE_THRESHOLD_PCT", "20.5", func(c Config) bool { return c.IdleThresholdPct == 20.5 }, "IdleThresholdPct=20.5"},
		{"ALLOCATOR_MARGINAL_THRESHOLD_PCT", "3", func(c Config) bool { return c.MarginalThresholdPct == 3.0 }, "MarginalThresholdPct=3"},
		{"ALLOCATOR_DAILY_CAP_WARNING_PCT", "80", func(c Config) bool { return c.DailyCapWarningPct == 80.0 }, "DailyCapWarningPct=80"},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Setenv(tt.envKey, tt.envValue)
			cfg := LoadConfig()
			if !tt.check(cfg) {
				t.Errorf("LoadConfig() with %s=%s: check failed, got %+v", tt.envKey, tt.envValue, cfg)
			}
		})
	}
}
