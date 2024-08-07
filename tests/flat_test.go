package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/repo"
	"avito-test-task/internal/usecase"
	"avito-test-task/pkg"
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func initDB(connString string) {
	m, err := migrate.New(
		"file:///home/nastya/avito/migrations",
		"postgres://test-user:test-password@localhost:5431/test-db?sslmode=disable")
	fmt.Println(err)
	m.Down()
	m.Up()
}

func initEnv() (domain.FlatUsecase, *zap.Logger, *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := "postgres://test-user:test-password@localhost:5431/test-db?sslmode=disable"
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	retryAdapter := pkg.NewPostgresRetryAdapter(pool, 3)
	flatRepo := repo.NewPostgresFlatRepo(pool, retryAdapter)
	flatUsecase := usecase.NewFlatUsecase(flatRepo)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	return flatUsecase, lg, pool
}

func TestCreateNormalFlat(t *testing.T) {
	flatUsecase, lg, pool := initEnv()
	initEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID, _ := uuid.Parse("019126ee-2b7d-758e-bb22-fe2e45b2db22")
	flatReq := domain.CreateFlatRequest{
		FlatID:  1,
		HouseID: 1,
		Price:   1000,
		Rooms:   2,
	}

	flat, err := flatUsecase.Create(ctx, userID, &flatReq, lg)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}

	expected := domain.CreateFlatResponse{
		ID:      1,
		HouseID: 1,
		Price:   1000,
		Rooms:   2,
		Status:  domain.CreatedStatus,
	}

	assert.Same(t, &expected, &flat)
}

func TestCreateHouseNotExist(t *testing.T) {
	flatUsecase, lg, pool := initEnv()
	initEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID, _ := uuid.Parse("019126ee-2b7d-758e-bb22-fe2e45b2db22")
	flatReq := domain.CreateFlatRequest{
		FlatID:  1,
		HouseID: 10,
		Price:   1000,
		Rooms:   2,
	}

	_, err := flatUsecase.Create(ctx, userID, &flatReq, lg)

	assert.Error(t, err)
}

func TestCreateUserNotExists(t *testing.T) {
	flatUsecase, lg, pool := initEnv()
	initEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID := uuid.New()
	flatReq := domain.CreateFlatRequest{
		FlatID:  1,
		HouseID: 10,
		Price:   1000,
		Rooms:   2,
	}

	_, err := flatUsecase.Create(ctx, userID, &flatReq, lg)

	assert.Error(t, err)
}

func TestCreateBadFlatID(t *testing.T) {
	flatUsecase, lg, pool := initEnv()
	initEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID, _ := uuid.Parse("019126ee-2b7d-758e-bb22-fe2e45b2db22")
	flatReq := domain.CreateFlatRequest{
		FlatID:  0,
		HouseID: 1,
		Price:   1000,
		Rooms:   2,
	}

	_, err := flatUsecase.Create(ctx, userID, &flatReq, lg)

	assert.Error(t, err)
}