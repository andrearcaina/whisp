package db

import (
	"context"
	"fmt"

	"github.com/andrearcaina/whisp/internal/config"
	"github.com/andrearcaina/whisp/internal/db/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

func New(cfg *config.Config) (*Database, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := generated.New(pool)

	return &Database{
		pool:    pool,
		queries: queries,
	}, nil
}

func (db *Database) GetQueries() *generated.Queries {
	return db.queries
}

func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *Database) GetPool() *pgxpool.Pool {
	return db.pool
}
