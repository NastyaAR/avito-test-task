package usecase

import (
	"avito-test-task/internal/domain"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type HouseUsecase struct {
	houseRepo domain.HouseRepo
	dbTimeout time.Duration
}

func NewHouseUsecase(houseRepo *domain.HouseRepo, timeout time.Duration) *HouseUsecase {
	return &HouseUsecase{
		houseRepo: *houseRepo,
		dbTimeout: timeout,
	}
}

func (u *HouseUsecase) Create(ctx context.Context, req *domain.CreateHouseRequest, lg *zap.Logger) (domain.CreateHouseResponse, error) {
	lg.Info("user usecase: create")

	date := time.Now().String()

	house := domain.House{
		HouseID:         req.HomeID,
		Address:         req.Address,
		ConstructYear:   req.Year,
		Developer:       req.Developer,
		CreateHouseDate: date,
		UpdateFlatDate:  date,
	}

	err := u.houseRepo.Create(ctx, &house, lg)
	if err != nil {
		lg.Warn("user usecase: create error", zap.Error(err))
		return domain.CreateHouseResponse{}, fmt.Errorf("user usecase: create error: %v", err.Error())
	}

	houseResponse := domain.CreateHouseResponse{
		HomeID:    req.HomeID,
		Address:   req.Address,
		Year:      req.Year,
		Developer: req.Developer,
		CreatedAt: date,
		UpdateAt:  date,
	}

	return houseResponse, nil
}

func (u *HouseUsecase) GetFlatsByHouseID(ctx context.Context, req *domain.FlatsByHouseRequest, lg *zap.Logger) ([]domain.FlatsByHouseResponse, error) {
	if req == nil {
		lg.Warn("house usecase: get flats by house id error: nil request")
		return nil, errors.New("house usecase: get flats by house id error: nil request")
	}

	err := u.houseRepo.GetAll()
}
