package usecase

import (
	"avito-test-task/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
	"time"
)

type HouseUsecase struct {
	houseRepo    domain.HouseRepo
	notifySender domain.NotifySender
	notifyRepo   domain.NotifyRepo
}

func NewHouseUsecase(houseRepo domain.HouseRepo, notifySender domain.NotifySender,
	notifyRepo domain.NotifyRepo, done chan bool, freq time.Duration, timeout time.Duration, lg *zap.Logger) *HouseUsecase {
	houseUsecase := HouseUsecase{
		houseRepo:    houseRepo,
		notifySender: notifySender,
		notifyRepo:   notifyRepo,
	}

	go houseUsecase.Notifying(done, freq, timeout, lg)

	return &houseUsecase
}

func (u *HouseUsecase) Create(ctx context.Context, req *domain.CreateHouseRequest, lg *zap.Logger) (domain.CreateHouseResponse, error) {
	lg.Info("house usecase: create")

	if req == nil {
		lg.Warn("house usecase: create error: bad request = nil")
		return domain.CreateHouseResponse{},
			fmt.Errorf("house usecase: create error: %w", domain.ErrHouse_BadRequest)
	}

	if req.Year < 0 {
		lg.Warn("house usecase: create errorL bad house year", zap.Int("year", req.Year))
		return domain.CreateHouseResponse{},
			fmt.Errorf("house usecase: create error: %w", domain.ErrHouse_BadYear)
	}

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

func parallelFlatFilter(flats []domain.Flat, lg *zap.Logger) domain.FlatsByHouseResponse {
	var (
		flatsArr []domain.SingleFlatResponse
	)

	lenOfPart := len(flats) / 3
	parts := make([][]domain.Flat, 0)
	for i := 0; i < 2; i++ {
		parts = append(parts, flats[i*lenOfPart:(i+1)*lenOfPart])
	}
	parts = append(parts, flats[2*lenOfPart:])

	mtx := sync.Mutex{}
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func(n int, part []domain.Flat, wg *sync.WaitGroup) {
			defer wg.Done()
			for j := 0; j < len(part); j++ {
				singleFlat := domain.SingleFlatResponse{
					ID:      parts[n][j].ID,
					HouseID: parts[n][j].HouseID,
					Price:   parts[n][j].Price,
					Rooms:   parts[n][j].Rooms,
					Status:  parts[n][j].Status,
				}
				mtx.Lock()
				flatsArr = append(flatsArr, singleFlat)
				mtx.Unlock()
			}
		}(i, parts[i], &wg)
	}
	wg.Wait()
	return domain.FlatsByHouseResponse{flatsArr}
}

func usualFlatFilter(flats []domain.Flat) domain.FlatsByHouseResponse {
	var (
		flatsArr []domain.SingleFlatResponse
	)
	for _, flat := range flats {
		singleFlat := domain.SingleFlatResponse{
			ID:      flat.ID,
			HouseID: flat.HouseID,
			Price:   flat.Price,
			Rooms:   flat.Rooms,
			Status:  flat.Status,
		}
		flatsArr = append(flatsArr, singleFlat)
	}

	return domain.FlatsByHouseResponse{flatsArr}
}

func (u *HouseUsecase) GetFlatsByHouseID(ctx context.Context, id int, status string, lg *zap.Logger) (domain.FlatsByHouseResponse, error) {
	if id < 0 {
		lg.Warn("house usecase: get flats by house id error: bad id", zap.Int("house_id", id))
		return domain.FlatsByHouseResponse{},
			fmt.Errorf("house usecase: get flats by house id error: %w", domain.ErrHouse_BadID)
	}

	if !IsCorrectFlatStatus(status) {
		lg.Warn("house usecase: get flats by house id error: bad status", zap.String("status", status))
		return domain.FlatsByHouseResponse{},
			fmt.Errorf("house usecase: get flats by house id error: %w", domain.ErrFlat_BadStatus)
	}

	flats, err := u.houseRepo.GetFlatsByHouseID(ctx, id, status, lg)
	if err != nil {
		lg.Warn("house usecase: get flats by house id error", zap.Error(err))
		return domain.FlatsByHouseResponse{}, fmt.Errorf("house usecase: get flats by house id error: %v", err.Error())
	}

	var flatsResponse domain.FlatsByHouseResponse
	if len(flats) < domain.FlatThreshhold {
		flatsResponse = usualFlatFilter(flats)
	} else {
		flatsResponse = parallelFlatFilter(flats, lg)
	}

	return flatsResponse, nil
}

func (uc *HouseUsecase) SubscribeByID(ctx context.Context, id int, userID uuid.UUID, lg *zap.Logger) error {
	lg.Info("house usecase: subscribe by id")

	if id < 0 {
		lg.Warn("house usecase: subscribe by id error: bad id", zap.Int("house_id", id))
		return fmt.Errorf("house usecase: suscribe by id error: %w", domain.ErrHouse_BadID)
	}

	err := uc.houseRepo.SubscribeByID(ctx, id, userID, lg)
	if err != nil {
		lg.Warn("house usecase: subscribe by id", zap.Error(err))
		return fmt.Errorf("house usecase: subscribe by id: %v", err.Error())
	}

	return nil
}

func (uc *HouseUsecase) Notifying(done chan bool, frequency time.Duration, timeout time.Duration, lg *zap.Logger) {
	for {
		select {
		case <-done:
			lg.Warn("house usecase: subscribing goroutine exited")
			return
		default:
			lg.Info("house usecase: subscribing goroutine working")
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			notifies, err := uc.notifyRepo.GetNoSendNotifies(ctx, lg)
			if err != nil {
				lg.Warn("house usecase: notifying error", zap.Error(err))
			}

			for _, notify := range notifies {
				msg := fmt.Sprintf("New flat with number %d in house %d!", notify.FlatID, notify.HouseID)
				err = uc.notifySender.SendEmail(ctx, notify.UserMail, msg)
				if err != nil {
					lg.Warn("house usecase: notifying error: send email error", zap.Error(err))
					continue
				} else {
					err = uc.notifyRepo.SendNotifyByID(ctx, notify.ID, lg)
				}
			}
			time.Sleep(frequency)
		}
	}
}
