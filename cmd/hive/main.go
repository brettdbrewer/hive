// Command hive runs the agent-first hive runtime.
package main

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lovyou-ai/eventgraph/go/pkg/actor"
	"github.com/lovyou-ai/eventgraph/go/pkg/actor/pgactor"
	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/store/pgstore"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"

	"github.com/lovyou-ai/hive/pkg/hive"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	human := flag.String("human", "", "Human operator name (required)")
	idea := flag.String("idea", "", "Seed idea for agents to work on")
	storeDSN := flag.String("store", "", "Store DSN (postgres://... or empty for in-memory)")
	autoApprove := flag.Bool("yes", false, "Auto-approve all authority requests")
	repo := flag.String("repo", "", "Path to repo for Implementer's Operate (default: current dir)")
	flag.Parse()

	if *human == "" {
		return fmt.Errorf("usage: hive --human name [--store postgres://...] --idea 'description' [--repo path] [--yes]")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Resolve store DSN: flag > DATABASE_URL env > in-memory.
	dsn := *storeDSN
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}

	// Open shared pool for Postgres, or nil for in-memory.
	var pool *pgxpool.Pool
	if dsn != "" {
		fmt.Fprintf(os.Stderr, "Postgres: %s\n", dsn)
		var err error
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return fmt.Errorf("postgres: %w", err)
		}
		defer pool.Close()
	}

	s, err := openStore(ctx, pool)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}
	defer s.Close()

	actors, err := openActorStore(ctx, pool)
	if err != nil {
		return fmt.Errorf("actor store: %w", err)
	}

	// Register the human operator.
	if pool != nil {
		fmt.Fprintln(os.Stderr, "WARNING: CLI key derivation is insecure for persistent Postgres stores.")
	}
	humanID, err := registerHuman(actors, *human)
	if err != nil {
		return fmt.Errorf("register human: %w", err)
	}

	// Bootstrap the event graph if empty.
	if err := bootstrapGraph(s, humanID); err != nil {
		return fmt.Errorf("bootstrap graph: %w", err)
	}

	// Resolve repo path.
	repoPath := *repo
	if repoPath == "" {
		repoPath = "."
	}

	// Create the hive runtime.
	rt, err := hive.New(ctx, hive.Config{
		Store:       s,
		Actors:      actors,
		HumanID:     humanID,
		AutoApprove: *autoApprove,
		RepoPath:    repoPath,
	})
	if err != nil {
		return fmt.Errorf("runtime: %w", err)
	}

	// Register starter agents.
	for _, def := range hive.StarterAgents(*human) {
		if err := rt.Register(def); err != nil {
			return fmt.Errorf("register %s: %w", def.Name, err)
		}
	}

	// Run all agents.
	if err := rt.Run(ctx, *idea); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	// Print summary.
	count, _ := s.Count()
	fmt.Fprintf(os.Stderr, "Events recorded: %d\n", count)
	return nil
}

// ────────────────────────────────────────────────────────────────────
// Store helpers (ported from previous main.go)
// ────────────────────────────────────────────────────────────────────

func openStore(ctx context.Context, pool *pgxpool.Pool) (store.Store, error) {
	if pool == nil {
		fmt.Fprintln(os.Stderr, "Store: in-memory")
		return store.NewInMemoryStore(), nil
	}
	fmt.Fprintln(os.Stderr, "Store: postgres")
	return pgstore.NewPostgresStoreFromPool(ctx, pool)
}

func openActorStore(ctx context.Context, pool *pgxpool.Pool) (actor.IActorStore, error) {
	if pool == nil {
		fmt.Fprintln(os.Stderr, "Actor store: in-memory")
		return actor.NewInMemoryActorStore(), nil
	}
	fmt.Fprintln(os.Stderr, "Actor store: postgres")
	return pgactor.NewPostgresActorStoreFromPool(ctx, pool)
}

func registerHuman(actors actor.IActorStore, displayName string) (types.ActorID, error) {
	h := sha256.Sum256([]byte("human:" + displayName))
	priv := ed25519.NewKeyFromSeed(h[:])
	pub := priv.Public().(ed25519.PublicKey)

	pk, err := types.NewPublicKey([]byte(pub))
	if err != nil {
		return types.ActorID{}, fmt.Errorf("public key: %w", err)
	}

	a, err := actors.Register(pk, displayName, event.ActorTypeHuman)
	if err != nil {
		return types.ActorID{}, err
	}
	return a.ID(), nil
}

func bootstrapGraph(s store.Store, humanID types.ActorID) error {
	head, err := s.Head()
	if err != nil {
		return fmt.Errorf("check head: %w", err)
	}
	if head.IsSome() {
		return nil // already bootstrapped
	}

	fmt.Fprintln(os.Stderr, "Bootstrapping event graph...")
	registry := event.DefaultRegistry()
	bsFactory := event.NewBootstrapFactory(registry)

	signer := &bootstrapSigner{humanID: humanID}
	bootstrap, err := bsFactory.Init(humanID, signer)
	if err != nil {
		return fmt.Errorf("create genesis event: %w", err)
	}
	if _, err := s.Append(bootstrap); err != nil {
		return fmt.Errorf("append genesis event: %w", err)
	}
	fmt.Fprintln(os.Stderr, "Event graph bootstrapped.")
	return nil
}

type bootstrapSigner struct {
	humanID types.ActorID
}

func (b *bootstrapSigner) Sign(data []byte) (types.Signature, error) {
	h := sha256.Sum256([]byte("signer:" + b.humanID.Value()))
	priv := ed25519.NewKeyFromSeed(h[:])
	sig := ed25519.Sign(priv, data)
	return types.NewSignature(sig)
}
