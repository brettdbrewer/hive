package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pruner deletes old telemetry data on a schedule.
type Pruner struct {
	pool     *pgxpool.Pool
	interval time.Duration
}

// NewPruner creates a pruner with the default 1-hour interval.
func NewPruner(pool *pgxpool.Pool) *Pruner {
	return &Pruner{
		pool:     pool,
		interval: 1 * time.Hour,
	}
}

// Start launches the pruning goroutine. It blocks until ctx is cancelled.
func (p *Pruner) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.prune(ctx)
		}
	}
}

func (p *Pruner) prune(ctx context.Context) {
	prunes := []struct {
		name  string
		query string
	}{
		{
			name:  "agent_snapshots",
			query: `DELETE FROM telemetry_agent_snapshots WHERE recorded_at < now() - interval '24 hours'`,
		},
		{
			name:  "hive_snapshots",
			query: `DELETE FROM telemetry_hive_snapshots WHERE recorded_at < now() - interval '7 days'`,
		},
		{
			name: "event_stream",
			query: `DELETE FROM telemetry_event_stream
				WHERE id NOT IN (SELECT id FROM telemetry_event_stream ORDER BY id DESC LIMIT 1000)`,
		},
	}

	for _, pr := range prunes {
		tag, err := p.pool.Exec(ctx, pr.query)
		if err != nil {
			fmt.Fprintf(os.Stderr, "telemetry prune %s: %v\n", pr.name, err)
			continue
		}
		if tag.RowsAffected() > 0 {
			fmt.Fprintf(os.Stderr, "telemetry prune %s: %d rows\n", pr.name, tag.RowsAffected())
		}
	}
}
