package usecase

import (
	"avito-test-task/internal/domain"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type FlatUsecase struct {
	flatRepo domain.FlatRepo
}

func NewFlatUsecase(flatRepo *domain.FlatRepo) *FlatUsecase {
	return &FlatUsecase{flatRepo: *flatRepo}
}

func (u *FlatUsecase) Create(ctx context.Context, flatReq *domain.CreateFlatRequest, lg *zap.Logger) (domain.CreateFlatResponse, error) {
	lg.Info("flat usecase: create")

	flat := domain.Flat{
		HouseID: flatReq.HouseID,
		Price:   flatReq.Price,
		Rooms:   flatReq.Rooms,
		Status:  domain.CreatedStatus,
	}

	createdFlat, err := u.flatRepo.Create(ctx, &flat, lg)
	if err != nil {
		lg.Warn("flat usecase repo: create error", zap.Error(err))
		return domain.CreateFlatResponse{}, fmt.Errorf("flat usecase: create error: %v", err.Error())
	}

	createdFlatResponse := domain.CreateFlatResponse{
		ID:      createdFlat.ID,
		HouseID: createdFlat.HouseID,
		Price:   createdFlat.Price,
		Rooms:   createdFlat.Rooms,
		Status:  createdFlat.Status,
	}

	return createdFlatResponse, nil
}

func (u *FlatUsecase) Update(ctx context.Context, newFlatData *domain.UpdateFlatRequest, lg *zap.Logger) (domain.CreateFlatResponse, error) {
	lg.Info("flat usecase: update")

	flat := domain.Flat{
		ID:      newFlatData.ID,
		HouseID: newFlatData.HouseID,
		Price:   newFlatData.Price,
		Rooms:   newFlatData.Rooms,
		Status:  newFlatData.Status,
	}

	updatedFlat, err := u.flatRepo.Update(ctx, &flat, lg)
	if err != nil {
		lg.Warn("flat usecase: update error", zap.Error(err))
		return domain.CreateFlatResponse{}, fmt.Errorf("flat usecase: update error: %v", err.Error())
	}

	updatedFlatResponse := domain.CreateFlatResponse{
		ID:      updatedFlat.ID,
		HouseID: updatedFlat.HouseID,
		Price:   updatedFlat.Price,
		Rooms:   updatedFlat.Rooms,
		Status:  updatedFlat.Status,
	}

	return updatedFlatResponse, nil

}
