package repo

import (
	"avito-test-task/internal/domain"
	"context"
	"database/sql"
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

	var (
		createdFlat domain.Flat
		moderatorId sql.NullInt32
	)
	query := `insert into flats(flat_id, house_id, price, rooms, status)
			values ($1, $2, $3, $4, $5) returning *`
	err := p.db.QueryRow(ctx, query, flat.ID, flat.HouseID,
		flat.Price, flat.Rooms, domain.CreatedStatus).Scan(&createdFlat.ID,
		&createdFlat.HouseID, &createdFlat.Price, &createdFlat.Rooms,
		&createdFlat.Status, &moderatorId)
	if err != nil {
		lg.Warn("postgres flat repo: create error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: create error: %v", err.Error())
	}
	if !moderatorId.Valid {
		createdFlat.ModeratorID = 0
	} else {
		createdFlat.ModeratorID = int(moderatorId.Int32)
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

func (p *PostgresFlatRepo) updatePrice(ctx context.Context, tx pgx.Tx, flatData *domain.Flat, lg *zap.Logger) error {
	lg.Info("postgres flat repo: update price")
	query := `update flats set price=$1 where flat_id=$2 and house_id=$3`
	_, err := tx.Exec(ctx, query, flatData.Price, flatData.ID, flatData.HouseID)
	if err != nil {
		lg.Warn("postgres flat repo: update price error", zap.Error(err))
		return fmt.Errorf("postgres flat repo: update price error: %v", err.Error())
	}
	return nil
}

func (p *PostgresFlatRepo) updateStatus(ctx context.Context, tx pgx.Tx, flatData *domain.Flat, lg *zap.Logger) error {
	lg.Info("postgres flat repo: update status")
	query := `update flats set status=$1 where flat_id=$2 and house_id=$3`
	_, err := tx.Exec(ctx, query, flatData.Status, flatData.ID, flatData.HouseID)
	if err != nil {
		lg.Warn("postgres flat repo: update status error", zap.Error(err))
		return fmt.Errorf("postgres flat repo: update status error: %v", err.Error())
	}
	return nil
}

func (p *PostgresFlatRepo) updateRooms(ctx context.Context, tx pgx.Tx, flatData *domain.Flat, lg *zap.Logger) error {
	lg.Info("postgres flat repo: update rooms")
	query := `update flats set rooms=$1 where flat_id=$2 and house_id=$3`
	_, err := tx.Exec(ctx, query, flatData.Rooms, flatData.ID, flatData.HouseID)
	if err != nil {
		lg.Warn("postgres flat repo: update rooms error", zap.Error(err))
		return fmt.Errorf("postgres flat repo: update rooms error: %v", err.Error())
	}
	return nil
}

func (p *PostgresFlatRepo) Update(ctx context.Context, newFlatData *domain.Flat, lg *zap.Logger) (domain.Flat, error) {
	lg.Info("postgres flat repo: update")

	var (
		flat domain.Flat
	)
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		lg.Warn("postgres flat repo: update error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: update error: %v", err.Error())
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("postgres flat repo: update error: %v", err.Error())
			}
		}
	}()

	if newFlatData.Price != domain.DefaultEmptyFlatValue {
		err = p.updatePrice(ctx, tx,
			&domain.Flat{ID: newFlatData.ID, HouseID: newFlatData.HouseID, Price: newFlatData.Price}, lg)
	}
	if newFlatData.Rooms != domain.DefaultEmptyFlatValue {
		err = p.updateRooms(ctx, tx,
			&domain.Flat{ID: newFlatData.ID, HouseID: newFlatData.HouseID, Rooms: newFlatData.Rooms}, lg)
	}
	if newFlatData.Status != domain.AnyStatus {
		err = p.updateStatus(ctx, tx,
			&domain.Flat{ID: newFlatData.ID, HouseID: newFlatData.HouseID, Status: newFlatData.Status}, lg)
	}

	query := `select flat_id, house_id, price, rooms, status from flats
	where flat_id=$1 and house_id=$2`

	err = p.db.QueryRow(ctx, query, newFlatData.ID, newFlatData.HouseID).Scan(&flat.ID, &flat.HouseID,
		&flat.Price, &flat.Rooms, &flat.Status)
	if err != nil {
		lg.Warn("postgres flat repo: update error", zap.Error(err))
		return domain.Flat{}, fmt.Errorf("postgres flat repo: update error: %v", err.Error())
	}

	if err = tx.Commit(ctx); err != nil {
		lg.Error("postgres flat repo: update error", zap.Error(err))
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
