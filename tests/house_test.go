package tests

import (
	"avito-test-task/internal/domain"
	"avito-test-task/internal/ports"
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

func initHouseEnv() (domain.HouseUsecase, *zap.Logger, *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := "postgres://test-user:test-password@localhost:5431/test-db?sslmode=disable"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err.Error())
	}

	retryAdapter := repo.NewPostgresRetryAdapter(pool, 3, time.Second)
	notifyRepo := repo.NewPostgresNotifyRepo(pool, retryAdapter)
	notifySender := ports.NewSender()

	done := make(chan bool, 1)
	defer func() {
		done <- true
	}()

	lg, _ := pkg.CreateLogger("../log.log", "prod")
	houseRepo := repo.NewPostgresHouseRepo(pool, retryAdapter)
	houseUsecase := usecase.NewHouseUsecase(houseRepo, notifySender, notifyRepo,
		done, time.Second, time.Second, lg)

	return houseUsecase, lg, pool
}

func TestCreateHouseNormal(t *testing.T) {
	houseUsecase, lg, pool := initHouseEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := domain.CreateHouseRequest{
		HomeID:    0,
		Address:   "address",
		Year:      1000,
		Developer: "dev",
	}

	resp, err := houseUsecase.Create(ctx, &req, lg)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	expected := domain.CreateHouseResponse{
		HomeID:    3,
		Address:   "address",
		Year:      1000,
		Developer: "dev",
		CreatedAt: "",
		UpdateAt:  "",
	}
	t.Log(resp)

	if expected.HomeID != resp.HomeID || expected.Address != resp.Address ||
		expected.Year != resp.Year || expected.Developer != resp.Developer {
		assert.Fail(t, "not same")
	}
}

func TestCreateBadHouseID(t *testing.T) {
	houseUsecase, lg, pool := initHouseEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := domain.CreateHouseRequest{
		HomeID:    0,
		Address:   "address",
		Year:      -1000,
		Developer: "dev",
	}

	_, err := houseUsecase.Create(ctx, &req, lg)
	assert.Error(t, err)
}

func TestGetFlatsByID(t *testing.T) {
	houseUsecase, lg, pool := initHouseEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := houseUsecase.GetFlatsByHouseID(ctx, 1, "created", lg)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	expectedFlat := domain.SingleFlatResponse{
		ID:      10,
		HouseID: 1,
		Price:   100,
		Rooms:   2,
		Status:  "created",
	}
	flats := make([]domain.SingleFlatResponse, 0)
	flats = append(flats, expectedFlat)

	expected := domain.FlatsByHouseResponse{Flats: flats}
	assert.Equal(t, expected, resp)
}

func TestGetFlatsByIDNullFlats(t *testing.T) {
	houseUsecase, lg, pool := initHouseEnv()
	initDB("")
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := houseUsecase.GetFlatsByHouseID(ctx, 2, "created", lg)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	expected := domain.FlatsByHouseResponse{Flats: nil}
	assert.Equal(t, expected, resp)
}
