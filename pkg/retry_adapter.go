package pkg

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type IPostgresRetryAdapter interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type PostgresRetryAdapter struct {
	db              *pgx.Conn
	numberOfRetries int
}

func NewPostgresRetryAdapter(db *pgx.Conn, retryNumber int) *PostgresRetryAdapter {
	return &PostgresRetryAdapter{
		db:              db,
		numberOfRetries: retryNumber,
	}
}

func (p *PostgresRetryAdapter) Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
	for i := 0; i < p.numberOfRetries; i++ {
		commTag, err := p.db.Exec(ctx, sql, arguments...)
		if err == nil {
			return commTag, nil
		}
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
	}
	return rows, err
}
