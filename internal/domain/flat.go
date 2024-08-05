package domain

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	CreatedStatus    = "created"
	ApprovedStatus   = "approved"
	DeclinedStatus   = "declined"
	ModeratingStatus = "on moderation"
	AnyStatus        = "any"
)

const DefaultEmptyFlatValue = -1

type Flat struct {
	ID          int
	HouseID     int
	UserID      uuid.UUID
	Price       int
	Rooms       int
	Status      string
	ModeratorID int
}

type UpdateFlat struct {
	ID      int
	HouseID int
	Values  map[string]any
}

type CreateFlatRequest struct {
	FlatID  int `json:"flat_id"`
	HouseID int `json:"house_id"`
	Price   int `json:"price"`
	Rooms   int `json:"rooms"`
}

type UpdateFlatRequest struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price,omitempty"`
	Rooms   int    `json:"rooms,omitempty"`
	Status  string `json:"status,omitempty"`
}

type CreateFlatResponse struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

type FlatUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, flatReq *CreateFlatRequest, lg *zap.Logger) (CreateFlatResponse, error)
	Update(ctx context.Context, userID uuid.UUID, newFlatData *UpdateFlatRequest, lg *zap.Logger) (CreateFlatResponse, error)
}

type FlatRepo interface {
	Create(ctx context.Context, flat *Flat, lg *zap.Logger) (Flat, error)
	DeleteByID(ctx context.Context, id int, lg *zap.Logger) error
	Update(ctx context.Context, newFlatData *Flat, lg *zap.Logger) (Flat, error)
	GetByID(ctx context.Context, id int, lg *zap.Logger) (Flat, error)
	GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]Flat, error)
}
