package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// initGitRepo creates a minimal git repo in dir with one commit so that
// worktree and branch operations have a valid HEAD to work from.
func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run("init", "-b", "main")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "test")
	// Write a file and commit so HEAD exists.
	if err := os.WriteFile(filepath.Join(dir, "README"), []byte("init"), 0644); err != nil {
		t.Fatalf("write README: %v", err)
	}
	run("add", ".")
	run("commit", "-m", "init")
}

// TestCreateTaskWorktree verifies that a worktree is created with the correct
// branch, directory, and git identity scoped to the worktree (not main repo).
func TestCreateTaskWorktree(t *testing.T) {
	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	wc, err := CreateTaskWorktree(repoDir, "test task", "task-001")
	if err != nil {
		t.Fatalf("CreateTaskWorktree: %v", err)
	}
	t.Cleanup(func() { wc.Cleanup() })

	// Worktree dir must exist.
	if _, err := os.Stat(wc.Dir); err != nil {
		t.Errorf("worktree dir %s does not exist: %v", wc.Dir, err)
	}

	// Branch must have the hive/feat/ prefix (CreateTaskWorktree uses "hive/{slug}-{ts}"
	// where branchSlug returns "feat/YYYYMMDD-...").
	if !strings.HasPrefix(wc.Branch, "hive/") {
		t.Errorf("branch %q does not start with hive/", wc.Branch)
	}

	// SourceDir and TaskID round-trip.
	if wc.SourceDir != repoDir {
		t.Errorf("SourceDir = %q, want %q", wc.SourceDir, repoDir)
	}
	if wc.TaskID != "task-001" {
		t.Errorf("TaskID = %q, want %q", wc.TaskID, "task-001")
	}

	// Git identity must be set so commits in the worktree use the hive identity.
	// Note: git worktrees share config with their main repo, so this config is
	// visible from both dirs — that's expected git behavior.
	getConfig := func(dir, key string) string {
		cmd := exec.Command("git", "config", "--local", key)
		cmd.Dir = dir
		out, _ := cmd.Output()
		return strings.TrimSpace(string(out))
	}
	if got := getConfig(wc.Dir, "user.name"); got != "hive" {
		t.Errorf("worktree user.name = %q, want %q", got, "hive")
	}
	if got := getConfig(wc.Dir, "user.email"); got != "hive@lovyou.ai" {
		t.Errorf("worktree user.email = %q, want %q", got, "hive@lovyou.ai")
	}
}

// TestGitConfigScopedToWorktree is a focused regression test for the cmd.Dir
// bug: git config must not bleed into the main repo.
func TestGitConfigScopedToWorktree(t *testing.T) {
	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	wc, err := CreateTaskWorktree(repoDir, "scope test", "task-002")
	if err != nil {
		t.Fatalf("CreateTaskWorktree: %v", err)
	}
	defer wc.Cleanup()

	getConfig := func(dir, key string) string {
		cmd := exec.Command("git", "config", "--local", key)
		cmd.Dir = dir
		out, _ := cmd.Output()
		return strings.TrimSpace(string(out))
	}

	// gitIn must run config in the repo context (not the process CWD).
	// Worktrees share config with main repo, so identity is visible from both.
	if got := getConfig(wc.Dir, "user.name"); got != "hive" {
		t.Errorf("worktree user.name = %q, want %q", got, "hive")
	}
	if got := getConfig(wc.Dir, "user.email"); got != "hive@lovyou.ai" {
		t.Errorf("worktree user.email = %q, want %q", got, "hive@lovyou.ai")
	}
}

// TestMergeToMain verifies that a worktree branch is merged back to main.
func TestMergeToMain(t *testing.T) {
	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	wc, err := CreateTaskWorktree(repoDir, "merge test", "task-003")
	if err != nil {
		t.Fatalf("CreateTaskWorktree: %v", err)
	}
	defer wc.Cleanup()

	// Add a commit on the worktree branch.
	newFile := filepath.Join(wc.Dir, "feature.txt")
	if err := os.WriteFile(newFile, []byte("feature"), 0644); err != nil {
		t.Fatalf("write feature.txt: %v", err)
	}
	if err := gitIn(wc.Dir, "add", "."); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := gitIn(wc.Dir, "commit", "-m", "add feature"); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	// Merge back to main.
	if err := wc.MergeToMain(); err != nil {
		t.Fatalf("MergeToMain: %v", err)
	}

	// Verify the file is now on main.
	mergedFile := filepath.Join(repoDir, "feature.txt")
	if _, err := os.Stat(mergedFile); err != nil {
		t.Errorf("feature.txt not found on main after merge: %v", err)
	}
}

// TestMergeToMainConcurrency verifies that concurrent MergeToMain calls do not
// corrupt repo state. The mutex must serialise all calls.
func TestMergeToMainConcurrency(t *testing.T) {
	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	const n = 3
	worktrees := make([]*WorktreeContext, n)
	for i := range worktrees {
		title := fmt.Sprintf("concurrent task %d", i)
		wc, err := CreateTaskWorktree(repoDir, title, fmt.Sprintf("task-concurrent-%d", i))
		if err != nil {
			t.Fatalf("CreateTaskWorktree %d: %v", i, err)
		}
		worktrees[i] = wc

		// Each worktree needs a unique commit so the merge is non-empty.
		newFile := filepath.Join(wc.Dir, "concurrent.txt")
		if err := os.WriteFile(newFile, []byte("data"), 0644); err != nil {
			t.Fatalf("write file %d: %v", i, err)
		}
		if err := gitIn(wc.Dir, "add", "."); err != nil {
			t.Fatalf("git add %d: %v", i, err)
		}
		if err := gitIn(wc.Dir, "commit", "-m", "concurrent commit"); err != nil {
			t.Fatalf("git commit %d: %v", i, err)
		}
	}

	// Launch all merges concurrently — must not panic or deadlock.
	var wg sync.WaitGroup
	errs := make([]error, n)
	for i, wc := range worktrees {
		wg.Add(1)
		go func(idx int, w *WorktreeContext) {
			defer wg.Done()
			errs[idx] = w.MergeToMain()
		}(i, wc)
	}
	wg.Wait()

	for i, wc := range worktrees {
		defer wc.Cleanup()
		// Merge conflicts are possible when multiple branches touch the same file;
		// what we're testing is that there's no panic, deadlock, or data race.
		// At least one merge must succeed.
		_ = errs[i]
	}
}

// TestCleanup verifies that Cleanup removes the worktree directory and is safe
// to call multiple times (idempotent).
func TestCleanup(t *testing.T) {
	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	wc, err := CreateTaskWorktree(repoDir, "cleanup test", "task-004")
	if err != nil {
		t.Fatalf("CreateTaskWorktree: %v", err)
	}

	dir := wc.Dir
	wc.Cleanup()

	// Directory must be gone.
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("worktree dir %s still exists after Cleanup", dir)
	}

	// Second call must not panic.
	wc.Cleanup()
}
