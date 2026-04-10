package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// Pool wraps pgxpool for dependency injection.
type Pool struct {
	*pgxpool.Pool
}

// New creates and verifies a connection pool.
func New(ctx context.Context, dsn string) (*Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Verify connectivity
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	log.Info().Msg("database connected")
	return &Pool{pool}, nil
}

// RefreshLatestPrices refreshes the materialised view concurrently.
func (p *Pool) RefreshLatestPrices(ctx context.Context) error {
	_, err := p.Exec(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY latest_prices")
	return err
}
