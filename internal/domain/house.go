package domain

import (
	"context"
	"go.uber.org/zap"
)

type House struct {
	HouseID         int
	Address         string
	ConstructYear   int
	Developer       string
	CreateHouseDate string
	UpdateFlatDate  string
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

type HouseUsecase interface {
	Create(ctx context.Context, req *CreateHouseRequest, lg *zap.Logger) (CreateHouseResponse, error)
	GetFlatsByHouseID(ctx context.Context, req *FlatsByHouseRequest, lg *zap.Logger) ([]FlatsByHouseResponse, error)
}

type HouseRepo interface {
	Create(ctx context.Context, house *House, lg *zap.Logger) error
	DeleteByID(ctx context.Context, id int, lg *zap.Logger) error
	Update(ctx context.Context, newHouseData *House, lg *zap.Logger) error
	GetByID(ctx context.Context, id int, lg *zap.Logger) (House, error)
	GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]House, error)
	GetFlatsByHouseID(ctx context.Context, id int, lg *zap.Logger) ([]Flat, error)
}
