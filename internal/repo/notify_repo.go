package repo

import (
	"avito-test-task/internal/domain"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PostgresNotifyRepo struct {
	db *pgx.Conn
}

func NewPostgresNotifyRepo(pg *pgx.Conn) *PostgresNotifyRepo {
	return &PostgresNotifyRepo{db: pg}
}

func (p *PostgresNotifyRepo) GetNoSendNotifies(ctx context.Context, lg *zap.Logger) ([]domain.Notify, error) {
	lg.Info("postgres notify repo: get no send notifies")

	query := `select * from new_flats_outbox where status=$1`
	rows, err := p.db.Query(ctx, query, domain.NoSendedNotifyStatus)
	if err != nil {
		lg.Warn("postgres notify repo: get no send notifies", zap.Error(err))
		return nil, fmt.Errorf("postgres notify error: get no send notifies: %v", err.Error())
	}

	var (
		notifies []domain.Notify
		notify   domain.Notify
	)

	for rows.Next() {
		err = rows.Scan(&notify.ID, &notify.FlatID, &notify.HouseID, &notify.UserMail, &notify.Status)
		if err != nil {
			lg.Warn("postgres notify repo: get no send notify error: scan notify error")
			continue
		}
		notifies = append(notifies, notify)
	}

	return notifies, err
}

func (p *PostgresNotifyRepo) SendNotifyByID(ctx context.Context, id int, lg *zap.Logger) error {
	lg.Info("postgres notify repo: send notify by id")

	query := `update new_flats_outbox set status=$1`
	_, err := p.db.Exec(ctx, query, domain.SendedNotifyStatus)
	if err != nil {
		lg.Warn("postgres notify repo: send notify by id error", zap.Error(err))
		return fmt.Errorf("postgres notify repo: send notify by id error: %v", err.Error())
	}

	return nil
}
