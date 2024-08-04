package domain

import (
	"context"
	"go.uber.org/zap"
)

const (
	CreatedStatus    = "created"
	ApprovedStatus   = "approved"
	DeclinedStatus   = "declined"
	ModeratingStatus = "on moderation"
	AnyStatus        = "any"
)

type Flat struct {
	ID          int
	HouseID     int
	Price       int
	Rooms       int
	Status      string
	ModeratorID int
}

type CreateFlatRequest struct {
	HouseID int `json:"house_id"`
	Price   int `json:"price"`
	Rooms   int `json:"rooms"`
}

type UpdateFlatRequest struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

type CreateFlatResponse struct {
	ID      int    `json:"id required"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

type FlatUsecase interface {
	Create(ctx context.Context, flatReq *CreateFlatRequest, lg *zap.Logger) (CreateFlatResponse, error)
	Update(ctx context.Context, newFlatData *UpdateFlatRequest, lg *zap.Logger) (CreateFlatResponse, error)
}

type FlatRepo interface {
	Create(ctx context.Context, flat *Flat, lg *zap.Logger) (Flat, error)
	DeleteByID(ctx context.Context, id int, lg *zap.Logger) error
	Update(ctx context.Context, newFlatData *Flat, lg *zap.Logger) (Flat, error)
	GetByID(ctx context.Context, id int, lg *zap.Logger) (Flat, error)
	GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]Flat, error)
}
