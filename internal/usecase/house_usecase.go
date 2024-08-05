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

func NewHouseUsecase(houseRepo domain.HouseRepo, dbTimeout time.Duration) *HouseUsecase {
	return &HouseUsecase{
		houseRepo: houseRepo,
		dbTimeout: dbTimeout,
	}
}

func (u *HouseUsecase) Create(ctx context.Context, req *domain.CreateHouseRequest, lg *zap.Logger) (domain.CreateHouseResponse, error) {
	lg.Info("user usecase: create")

	date := time.Now()

	house := domain.House{
		HouseID:         req.HomeID,
		Address:         req.Address,
		ConstructYear:   req.Year,
		Developer:       req.Developer,
		CreateHouseDate: date,
		UpdateFlatDate:  date,
	}

	house, err := u.houseRepo.Create(ctx, &house, lg)
	if err != nil {
		lg.Warn("user usecase: create error", zap.Error(err))
		return domain.CreateHouseResponse{}, fmt.Errorf("user usecase: create error: %v", err.Error())
	}

	houseResponse := domain.CreateHouseResponse{
		HomeID:    house.HouseID,
		Address:   house.Address,
		Year:      house.ConstructYear,
		Developer: house.Developer,
		CreatedAt: house.CreateHouseDate.Format(time.DateTime),
		UpdateAt:  house.UpdateFlatDate.Format(time.DateTime),
	}

	return houseResponse, nil
}

func (u *HouseUsecase) GetFlatsByHouseID(ctx context.Context, id int, status string, lg *zap.Logger) (domain.FlatsByHouseResponse, error) {
	if id < 0 {
		lg.Warn("house usecase: get flats by house id error: nil request")
		return domain.FlatsByHouseResponse{}, errors.New("house usecase: get flats by house id error: nil request")
	}

	if !IsCorrectFlatStatus(status) {
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

func (uc *HouseUsecase) SubscribeByID(ctx context.Context, id int, req *domain.SubscribeOnHouseRequest, lg *zap.Logger) error {
	lg.Info("house usecase: subscribe by id")

	err := uc.houseRepo.SubscribeByID(ctx, id, req.Mail, lg)
	if err != nil {
		lg.Warn("house usecase: subscribe by id", zap.Error(err))
		return fmt.Errorf("house usecase: subscribe by id: %v", err.Error())
	}

	return nil
}

func (uc *HouseUsecase) Subscribing(done chan bool, frequency time.Duration, lg *zap.Logger) {
	for {
		select {
		case <-done:
			lg.Warn("house usecase: subscribing goroutine exited")
			return
		default:
			lg.Info("house usecase: subscribing goroutine working")
			ctx, cancel := context.WithTimeout(context.Background(), uc.dbTimeout)
			defer cancel()

			reqs, err := uc.
			if err != nil {
				lg.Warnf("pay adapter: paying error: %v", err.Error())
			}

			for _, req := range reqs {

			}
			time.Sleep(frequency)
		}
	}
}
