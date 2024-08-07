package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/repo"
	"avito-test-task/internal/usecase"
	"avito-test-task/pkg"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func initUserEnv() (domain.UserUsecase, *zap.Logger, *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := "postgres://test-user:test-password@localhost:5431/test-db?sslmode=disable"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	retryAdapter := repo.NewPostgresRetryAdapter(pool, 3, time.Second)
	userRepo := repo.NewPostrgesUserRepo(pool, retryAdapter)
	userUsecase := usecase.NewUserUsecase(userRepo)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	return userUsecase, lg, pool
}

func TestRegisterNormalUser(t *testing.T) {
	userUsecase, lg, pool := initUserEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userReq := domain.RegisterUserRequest{
		Email:    "test@mail.ru",
		Password: "pass",
		UserType: "client",
	}

	_, err := userUsecase.Register(ctx, &userReq, lg)
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func TestRegisterBadMail(t *testing.T) {
	userUsecase, lg, pool := initUserEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userReq := domain.RegisterUserRequest{
		Email:    "",
		Password: "pass",
		UserType: "client",
	}

	_, err := userUsecase.Register(ctx, &userReq, lg)
	assert.Error(t, err)
}

func TestRegisterBadPassword(t *testing.T) {
	userUsecase, lg, pool := initUserEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userReq := domain.RegisterUserRequest{
		Email:    "test@mail.ru",
		Password: "",
		UserType: "client",
	}

	_, err := userUsecase.Register(ctx, &userReq, lg)
	assert.Error(t, err)
}

//func TestLoginNormal(t *testing.T) {
//
//}