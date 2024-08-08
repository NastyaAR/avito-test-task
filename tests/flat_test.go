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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func initDB(connString string) {
	m, _ := migrate.New(
		"file:///../test_migrations",
		"postgres://test-user:test-password@localhost:5431/test-db?sslmode=disable")
	err := m.Force(20240806143730)
	fmt.Println(err)
	err = m.Down()
	fmt.Println(err)
	err = m.Up()
	fmt.Println(err)
}

func initFlatEnv() (domain.FlatUsecase, *zap.Logger, *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := "postgres://test-user:test-password@localhost:5431/test-db?sslmode=disable"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	retryAdapter := repo.NewPostgresRetryAdapter(pool, 3, time.Second)
	flatRepo := repo.NewPostgresFlatRepo(pool, retryAdapter)
	flatUsecase := usecase.NewFlatUsecase(flatRepo)
	lg, _ := pkg.CreateLogger("../log.log", "prod")

	return flatUsecase, lg, pool
}

func TestCreateNormalFlat(t *testing.T) {
	flatUsecase, lg, pool := initFlatEnv()
	initFlatEnv()
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

	assert.Equal(t, expected, flat)
}

func TestCreateHouseNotExist(t *testing.T) {
	flatUsecase, lg, pool := initFlatEnv()
	initFlatEnv()
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
	flatUsecase, lg, pool := initFlatEnv()
	initFlatEnv()
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
	flatUsecase, lg, pool := initFlatEnv()
	initFlatEnv()
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

func TestUpdateBadIDFlat(t *testing.T) {
	flatUsecase, lg, pool := initFlatEnv()
	initFlatEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	modID, _ := uuid.Parse("019126ee-2b7d-758e-bb22-fe2e45b2db23")
	flatReq := domain.UpdateFlatRequest{
		ID:      0,
		HouseID: 1,
		Status:  "on moderation",
	}

	_, err := flatUsecase.Update(ctx, modID, &flatReq, lg)
	assert.Error(t, err)
}
