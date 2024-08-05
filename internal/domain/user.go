package domain

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	Moderator = "moderator"
	Client    = "client"
)

var SessionUserID = uuid.New()

type User struct {
	UserID   uuid.UUID
	Mail     string
	Password string
	Role     string
}

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

type RegisterUserResponse struct {
	UserID uuid.UUID `json:"user_id"`
}

type LoginUserRequest struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

type DummyLoginRequest struct {
	UserType string `json:"user_type"`
}

type UserUsecase interface {
	Register(ctx context.Context, userReq *RegisterUserRequest, lg *zap.Logger) (RegisterUserResponse, error)
	Login(ctx context.Context, userReq *LoginUserRequest, lg *zap.Logger) (LoginUserResponse, error)
	DummyLogin(ctx context.Context, userType string, lg *zap.Logger) (LoginUserResponse, error)
}

type UserRepo interface {
	Create(ctx context.Context, user *User, lg *zap.Logger) error
	DeleteByID(ctx context.Context, id string, lg *zap.Logger) error
	Update(ctx context.Context, newUserData *User, lg *zap.Logger) error
	GetByID(ctx context.Context, id uuid.UUID, lg *zap.Logger) (User, error)
	GetAll(ctx context.Context, offset int, limit int, lg *zap.Logger) ([]User, error)
}
