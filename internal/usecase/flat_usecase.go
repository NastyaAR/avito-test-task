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

func NewFlatUsecase(flatRepo domain.FlatRepo) *FlatUsecase {
	return &FlatUsecase{flatRepo: flatRepo}
}

func (u *FlatUsecase) Create(ctx context.Context, flatReq *domain.CreateFlatRequest, lg *zap.Logger) (domain.CreateFlatResponse, error) {
	lg.Info("flat usecase: create")

	if flatReq == nil {
		lg.Warn("flat usecase: create error: bad flat request: nil")
		return domain.CreateFlatResponse{}, fmt.Errorf("flat usecase: create error: nil flat request")
	}

	if flatReq.FlatID < 1 {
		lg.Warn("flat usecase: create error: bad flat id", zap.Int("flat_id", flatReq.FlatID))
		return domain.CreateFlatResponse{}, fmt.Errorf("flat usecase: create error: bad flat id")
	}

	if flatReq.HouseID < 1 {
		lg.Warn("flat usecase: create error: bad house id", zap.Int("house_id", flatReq.HouseID))
		return domain.CreateFlatResponse{}, fmt.Errorf("flat usecase: create error: bad house id")
	}

	if flatReq.Rooms < 1 {
		lg.Warn("flat usecase: create error: bad rooms", zap.Int("rooms", flatReq.Rooms))
		return domain.CreateFlatResponse{}, fmt.Errorf("flat usecase: create error: bad rooms")
	}

	if flatReq.Price < 0 {
		lg.Warn("flat usecase: create error: bad price", zap.Int("price", flatReq.Price))
		return domain.CreateFlatResponse{}, fmt.Errorf("flat usecase: create error: bad price")
	}

	flat := domain.Flat{
		ID:      flatReq.FlatID,
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

	if newFlatData == nil {
		lg.Warn("flat usecase: update error: bad newFlatData = nil")
		return domain.CreateFlatResponse{}, fmt.Errorf("flat usecase: update error: bad newFlatData = nil")
	}

	if newFlatData.ID < 1 {
		lg.Warn("flat usecase: update error: bad flat id", zap.Int("flat_id", newFlatData.ID))
		return domain.CreateFlatResponse{}, fmt.Errorf("flat usecase: update error: bad flat id")
	}

	if newFlatData.HouseID < 1 {
		lg.Warn("flat usecase: update error: bad house id", zap.Int("house_id", newFlatData.HouseID))
		return domain.CreateFlatResponse{}, fmt.Errorf("flat usecase: update error: bad house id")
	}

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
