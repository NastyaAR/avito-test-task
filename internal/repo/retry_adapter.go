package repo

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type IPostgresRetryAdapter interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type PostgresRetryAdapter struct {
	db              *pgxpool.Pool
	numberOfRetries int
	sleepTimeMs     time.Duration
}

func NewPostgresRetryAdapter(db *pgxpool.Pool, retryNumber int, sleepTimeMs time.Duration) *PostgresRetryAdapter {
	return &PostgresRetryAdapter{
		db:              db,
		numberOfRetries: retryNumber,
		sleepTimeMs:     sleepTimeMs,
	}
}

func (p *PostgresRetryAdapter) Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
	for i := 0; i < p.numberOfRetries; i++ {
		commTag, err := p.db.Exec(ctx, sql, arguments...)
		if err == nil {
			return commTag, nil
		}
		time.Sleep(p.sleepTimeMs)
	}
	return pgconn.CommandTag{}, err
}

func (p *PostgresRetryAdapter) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	var rows pgx.Rows
	for i := 0; i < p.numberOfRetries; i++ {
		rows, err := p.db.Query(ctx, sql, args...)
		if err == nil {
			rows.Next()
			return rows
		}
		time.Sleep(p.sleepTimeMs)
	}
	return rows
}

func (p *PostgresRetryAdapter) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	var (
		rows pgx.Rows
		err  error
	)
	for i := 0; i < p.numberOfRetries; i++ {
		rows, err = p.db.Query(ctx, sql, args...)
		if err == nil {
			return rows, nil
		}
		time.Sleep(p.sleepTimeMs)
	}
	return rows, err
}
