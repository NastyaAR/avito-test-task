package repo

import (
	"avito-test-task/internal/domain"
	"avito-test-task/pkg"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"time"
)

type PostgresFlatRepo struct {
	db           *pgxpool.Pool
	retryAdapter pkg.IPostgresRetryAdapter
}

func NewPostgresFlatRepo(db *pgxpool.Pool, retryAdapter pkg.IPostgresRetryAdapter) *PostgresFlatRepo {
	return &PostgresFlatRepo{
		db:           db,
		retryAdapter: retryAdapter,
	}
}

func (p *PostgresFlatRepo) Create(ctx context.Context, flat *domain.Flat, lg *zap.Logger) (domain.Flat, error) {
	lg.Info("postgres flat repo: create")

	var (
		createdFlat domain.Flat
	)

	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		lg.Warn("postgres flat repo: create error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: create error: %v", err.Error())
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("postgres flat repo: create error: %v", err.Error())
			}
		}
	}()

	query := `insert into flats(flat_id, house_id, user_id, price, rooms, status)
			values ($1, $2, $3, $4, $5, $6) 
			returning flat_id, house_id, user_id, price, rooms, status`
	err = tx.QueryRow(ctx, query, flat.ID, flat.HouseID, flat.UserID,
		flat.Price, flat.Rooms, domain.CreatedStatus).Scan(&createdFlat.ID,
		&createdFlat.HouseID, &createdFlat.UserID, &createdFlat.Price, &createdFlat.Rooms,
		&createdFlat.Status)
	if err != nil {
		lg.Warn("postgres flat repo: create error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: create error: %v", err.Error())
	}

	date := time.Now()
	query = `update houses set update_flat_date=$1 where house_id=$2`
	_, err = tx.Exec(ctx, query, date, createdFlat.HouseID)
	if err != nil {
		lg.Warn("postgres flat repo: create error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: create error: %v", err.Error())
	}

	if err = tx.Commit(ctx); err != nil {
		lg.Error("postgres flat repo: create error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: create error: %v", err.Error())
	}

	return createdFlat, nil
}

func (p *PostgresFlatRepo) DeleteByID(ctx context.Context, flatID int, houseID int, lg *zap.Logger) error {
	lg.Info("postgres flat repo: delete by id")

	query := `delete from flats where flat_id=$1 and house_id=$2`
	_, err := p.retryAdapter.Exec(ctx, query, flatID, houseID)
	if err != nil {
		lg.Warn("postgres flat repo: delete by id error", zap.Error(err))
		return fmt.Errorf("postgres flat repo: delete by id error: %v", err.Error())
	}

	return nil
}

func (p *PostgresFlatRepo) Update(ctx context.Context, moderatorID uuid.UUID, newFlatData *domain.Flat, lg *zap.Logger) (domain.Flat, error) {
	lg.Info("postgres flat repo: update")

	var (
		flat domain.Flat
	)

	query := `select flat_id, house_id, user_id, price, rooms, status 
	from update_status($1, $2, $3, $4)`

	err := p.retryAdapter.QueryRow(ctx, query, newFlatData.Status,
		newFlatData.ID, newFlatData.HouseID, moderatorID).Scan(&flat.ID, &flat.HouseID, &flat.UserID,
		&flat.Price, &flat.Rooms, &flat.Status)
	if err != nil {
		lg.Warn("postgres flat repo: update error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: update error: %v", err.Error())
	}

	return flat, nil
}

func (p *PostgresFlatRepo) GetByID(ctx context.Context, flatID int, houseID int, lg *zap.Logger) (domain.Flat, error) {
	var flat domain.Flat
	lg.Info("postgres flat repo: get by id")

	query := `select flat_id, house_id, user_id, price, rooms, status
	from flats where flat_id=$1 and house_id=$2`
	err := p.retryAdapter.QueryRow(ctx, query, flatID, houseID).Scan(&flat.ID, &flat.HouseID, &flat.UserID,
		&flat.Price, &flat.Rooms, &flat.Status)
	if err != nil {
		lg.Warn("postgres flat repo: get by id error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: get by id error: %v", err.Error())
	}

	return flat, nil
}

func (p *PostgresFlatRepo) GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]domain.Flat, error) {
	lg.Info("postgres flat repo: get all")

	query := `select flat_id, house_id, user_id, price, rooms, status from flats limit $1 offset $2`
	rows, err := p.retryAdapter.Query(ctx, query, limit, offset)
	if err != nil {
		lg.Warn("postgres flat repo: get all error", zap.Error(err))
		return nil, fmt.Errorf("postgres flat repo: get all error: %v", err.Error())
	}

	var (
		flats []domain.Flat
		flat  domain.Flat
	)
	for rows.Next() {
		err = rows.Scan(&flat.ID, &flat.HouseID, &flat.UserID,
			&flat.Price, &flat.Rooms, &flat.Status)
		if err != nil {
			lg.Warn("postgres flat repo: get all error: scan flat error", zap.Error(err))
			continue
		}
		flats = append(flats, flat)
	}

	return flats, err
}
