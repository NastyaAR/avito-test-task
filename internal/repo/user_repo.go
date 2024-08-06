package repo

import (
	"avito-test-task/internal/domain"
	"avito-test-task/pkg"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type PostgresUserRepo struct {
	db           *pgxpool.Pool
	retryAdapter pkg.IPostgresRetryAdapter
}

func NewPostrgesUserRepo(db *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (p *PostgresUserRepo) Create(ctx context.Context, user *domain.User, lg *zap.Logger) error {
	lg.Info("create user", zap.String("user_id", user.UserID.String()))

	query := `insert into users(user_id, mail, password, role) values ($1, $2, $3, $4)`
	_, err := p.retryAdapter.Exec(ctx, query, user.UserID, user.Mail, user.Password, user.Role)
	if err != nil {
		lg.Warn("postgres create user error", zap.Error(err))
		return err
	}

	return nil
}

func (p *PostgresUserRepo) DeleteByID(ctx context.Context, id string, lg *zap.Logger) error {
	lg.Info("delete user", zap.String("user_id", id))

	query := `delete from users where id=$1`
	_, err := p.retryAdapter.Exec(ctx, query, id)
	if err != nil {
		lg.Warn("postgres delete user error", zap.Error(err))
		return err
	}

	return nil
}

func (p *PostgresUserRepo) Update(ctx context.Context, newUserData *domain.User, lg *zap.Logger) error {
	lg.Info("update user", zap.String("user_id", newUserData.UserID.String()))

	query := `update users set user_id=$1,	
			mail=$2,
			password=$3,
			role=$4`
	_, err := p.retryAdapter.Exec(ctx, query, newUserData.UserID, newUserData.Mail,
		newUserData.Password, newUserData.Role)
	if err != nil {
		lg.Warn("postgres update user error", zap.Error(err))
		return err
	}

	return nil
}

func (p *PostgresUserRepo) GetByID(ctx context.Context, id uuid.UUID, lg *zap.Logger) (domain.User, error) {
	var user domain.User
	lg.Info("get user by id", zap.String("user_id", id.String()))

	query := `select * from users where user_id=$1`
	err := p.retryAdapter.QueryRow(ctx, query, id).Scan(&user.UserID, &user.Mail, &user.Password, &user.Role)
	if err != nil {
		lg.Warn("postgres get by id user error", zap.Error(err))
		return domain.User{}, err
	}

	return user, nil
}

func (p *PostgresUserRepo) GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]domain.User, error) {
	lg.Info("get users", zap.Int("offset", offset), zap.Int("limit", limit))

	query := `select * from users limit $1 offset $2`
	rows, err := p.retryAdapter.Query(ctx, query, limit, offset)

	var (
		users []domain.User
		user  domain.User
	)
	for rows.Next() {
		err = rows.Scan(&user.UserID, &user.Mail, &user.Password, &user.Role)
		if err != nil {
			lg.Warn("postgres user get all error: scan user error")
			continue
		}
		users = append(users, user)
	}

	return users, err
}
