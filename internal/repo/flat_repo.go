package repo

import (
	"avito-test-task/internal/domain"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PostgresFlatRepo struct {
	db *pgx.Conn
}

func NewPostgresFlatRepo(db *pgx.Conn) *PostgresFlatRepo {
	return &PostgresFlatRepo{db: db}
}

func (p *PostgresFlatRepo) Create(ctx context.Context, flat *domain.Flat, lg *zap.Logger) (domain.Flat, error) {
	lg.Info("postgres flat repo: create")

	var createdFlat domain.Flat
	query := `insert into flats(flat_id, house_id, price, rooms, status)
			values ($1, $2, $3, $4, $5) returning *`
	err := p.db.QueryRow(ctx, query, flat.ID, flat.HouseID,
		flat.Price, flat.Rooms, domain.CreatedStatus).Scan(&createdFlat.ID,
		&createdFlat.HouseID, &createdFlat.Price, &createdFlat.Rooms,
		&createdFlat.Status)
	if err != nil {
		lg.Warn("postgres flat repo: create error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: create error: %v", err.Error())
	}

	return createdFlat, nil
}

func (p *PostgresFlatRepo) DeleteByID(ctx context.Context, id int, lg *zap.Logger) error {
	lg.Info("postgres flat repo: delete by id")

	query := `delete from flats where flat_id=$1`
	_, err := p.db.Exec(ctx, query, id)
	if err != nil {
		lg.Warn("postgres flat repo: delete by id error", zap.Error(err))
		return fmt.Errorf("postgres flat repo: delete by id error: %v", err.Error())
	}

	return nil
}

func (p *PostgresFlatRepo) Update(ctx context.Context, newFlatData *domain.Flat, lg *zap.Logger) (Flat, error) {
	lg.Info("postgres flat repo: update")

	var flat domain.Flat
	query := `update flats set flat_id=$1,
								house_id=$2,
                 				price=$3,
                 				rooms=$4,
                 				status=$5
				returning *`
	err := p.db.QueryRow(ctx, query, newFlatData.ID, newFlatData.HouseID,
		newFlatData.Price, newFlatData.Rooms,
		newFlatData.Status).Scan(&flat.ID, &flat.HouseID,
		&flat.Price, &flat.Rooms, &flat.Status)
	if err != nil {
		lg.Warn("postgres flat repo: update error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: update error: %v", err.Error())
	}

	return flat, nil
}

func (p *PostgresFlatRepo) GetByID(ctx context.Context, id int, lg *zap.Logger) (domain.Flat, error) {
	var flat domain.Flat
	lg.Info("postgres flat repo: get by id")

	query := `select * from flats where flat_id=$1`
	err := p.db.QueryRow(ctx, query, id).Scan(&flat.ID, &flat.HouseID,
		&flat.Price, &flat.Rooms,
		&flat.Status, &flat.ModeratorID)
	if err != nil {
		lg.Warn("postgres flat repo: get by id error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: get by id error: %v", err.Error())
	}

	return flat, nil
}

func (p *PostgresFlatRepo) GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]domain.Flat, error) {
	lg.Info("postgres flat repo: get all")

	query := `select * from flats limit $1 offset $2`
	rows, err := p.db.Query(ctx, query, limit, offset)
	if err != nil {
		lg.Warn("postgres flat repo: get all error", zap.Error(err))
		return nil, fmt.Errorf("postgres flat repo: get all error: %v", err.Error())
	}

	var (
		flats []domain.Flat
		flat  domain.Flat
	)
	for rows.Next() {
		err = rows.Scan(&flat.ID, &flat.HouseID,
			&flat.Price, &flat.Rooms,
			&flat.Status, &flat.ModeratorID)
		if err != nil {
			lg.Warn("postgres flat repo: get all error: scan flat error", zap.Error(err))
			continue
		}
		flats = append(flats, flat)
	}

	return flats, err
}
