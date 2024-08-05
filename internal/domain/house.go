package domain

import (
	"context"
	"go.uber.org/zap"
	"time"
)

type House struct {
	HouseID         int
	Address         string
	ConstructYear   int
	Developer       string
	CreateHouseDate time.Time
	UpdateFlatDate  time.Time
}

type CreateHouseRequest struct {
	HomeID    int    `json:"id"`
	Address   string `json:"address"`
	Year      int    `json:"year"`
	Developer string `json:"developer"`
}

type CreateHouseResponse struct {
	HomeID    int    `json:"id"`
	Address   string `json:"address"`
	Year      int    `json:"year"`
	Developer string `json:"developer"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
}

type FlatsByHouseRequest struct {
	ID int `json:"id"`
}

type FlatsByHouseResponse struct {
	Flats []SingleFlatResponse `json:"flats"`
}

type SingleFlatResponse struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

type SubscribeOnHouseRequest struct {
	Mail string `json:"email"`
}

type HouseUsecase interface {
	Create(ctx context.Context, req *CreateHouseRequest, lg *zap.Logger) (CreateHouseResponse, error)
	GetFlatsByHouseID(ctx context.Context, id int, status string, lg *zap.Logger) (FlatsByHouseResponse, error)
	SubscribeByID(ctx context.Context, id int, req *SubscribeOnHouseRequest, lg *zap.Logger) error
}

type HouseRepo interface {
	Create(ctx context.Context, house *House, lg *zap.Logger) (House, error)
	DeleteByID(ctx context.Context, id int, lg *zap.Logger) error
	Update(ctx context.Context, newHouseData *House, lg *zap.Logger) error
	GetByID(ctx context.Context, id int, lg *zap.Logger) (House, error)
	GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]House, error)
	GetFlatsByHouseID(ctx context.Context, id int, lg *zap.Logger) ([]Flat, error)
	SubscribeByID(ctx context.Context, id int, email string, lg *zap.Logger) error
}
