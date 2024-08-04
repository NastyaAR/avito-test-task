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
}

func NewHouseUsecase(houseRepo *domain.HouseRepo) *HouseUsecase {
	return &HouseUsecase{
		houseRepo: *houseRepo,
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

func isCorrectFlatStatus(status string) bool {
	return status == domain.CreatedStatus || status == domain.ApprovedStatus ||
		status == domain.DeclinedStatus || status == domain.AnyStatus
}

func (u *HouseUsecase) GetFlatsByHouseID(ctx context.Context, id int, status string, lg *zap.Logger) (domain.FlatsByHouseResponse, error) {
	if id < 0 {
		lg.Warn("house usecase: get flats by house id error: nil request")
		return domain.FlatsByHouseResponse{}, errors.New("house usecase: get flats by house id error: nil request")
	}

	if !isCorrectFlatStatus(status) {
		lg.Warn("house usecase: get flats by house id error: bad status", zap.String("status", status))
		return domain.FlatsByHouseResponse{}, errors.New("house usecase: get flats by house id error: bad status")
	}

	flats, err := u.houseRepo.GetFlatsByHouseID(ctx, id, lg)
	if err != nil {
		lg.Warn("house usecase: get flats by house id error", zap.Error(err))
		return domain.FlatsByHouseResponse{}, fmt.Errorf("house usecase: get flats by house id error: %v", err.Error())
	}

	var (
		flatsArr   []domain.SingleFlatResponse
		singleFlat domain.SingleFlatResponse
	)
	for _, flat := range flats {
		if flat.Status == status || status == domain.AnyStatus {
			singleFlat = domain.SingleFlatResponse{
				ID:      flat.ID,
				HouseID: flat.HouseID,
				Price:   flat.Price,
				Rooms:   flat.Rooms,
				Status:  flat.Status,
			}
		}

		flatsArr = append(flatsArr, singleFlat)
	}

	return domain.FlatsByHouseResponse{flatsArr}, nil
}
