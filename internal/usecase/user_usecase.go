package usecase

import (
	"avito-test-task/internal/domain"
	"avito-test-task/pkg"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserUsecase struct {
	userRepo domain.UserRepo
}

func NewUserUsecase(userRepo *domain.UserRepo) *UserUsecase {
	return &UserUsecase{
		userRepo: *userRepo,
	}
}

func (u *UserUsecase) Register(ctx context.Context, userReq *domain.RegisterUserRequest, lg *zap.Logger) (domain.RegisterUserResponse, error) {
	lg.Info("user usecase: register", zap.String("user_type", userReq.UserType))

	if userReq.UserType != domain.Client && userReq.UserType != domain.Moderator {
		lg.Warn("user usecase: register error: bad role", zap.String("role", userReq.UserType))
		return domain.RegisterUserResponse{}, errors.New("user usecase: register error: bad role")
	}

	encryptedPassword, err := pkg.EncryptPassword(userReq.Password, lg)
	if err != nil {
		lg.Warn("user usecase: register error", zap.Error(err))
		return domain.RegisterUserResponse{}, fmt.Errorf("user usecase: register error: %v", err.Error())
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		lg.Warn("user usecase: register error", zap.Error(err))
		return domain.RegisterUserResponse{}, fmt.Errorf("user usecase: register error: %v", err.Error())
	}

	user := domain.User{
		UserID:   uuid,
		Mail:     userReq.Email,
		Password: encryptedPassword,
		Role:     userReq.UserType,
	}

	err = u.userRepo.Create(ctx, &user, lg)
	if err != nil {
		lg.Warn("user usecase: register error", zap.Error(err))
		return domain.RegisterUserResponse{}, fmt.Errorf("user usecase: register error: %v", err.Error())
	}

	return domain.RegisterUserResponse{uuid}, nil
}

func (u *UserUsecase) Login(ctx context.Context, userReq *domain.LoginUserRequest, lg *zap.Logger) (domain.LoginUserResponse, error) {
	lg.Info("user usecase: login")

	expectedUser, err := u.userRepo.GetByID(ctx, userReq.ID, lg)
	if err != nil {
		lg.Warn("user usescase: login error", zap.Error(err))
		return domain.LoginUserResponse{}, fmt.Errorf("user usecase: login error: %v", err.Error())
	}

	err = pkg.IsEqualPasswords(expectedUser.Password, userReq.Password)
	if err != nil {
		lg.Warn("user usecase: login error", zap.Error(err))
		return domain.LoginUserResponse{}, fmt.Errorf("user usecase: login error: %v", err.Error())
	}

	token, err := pkg.GenerateJWTToken(expectedUser.UserID, expectedUser.Role)
	if err != nil {
		lg.Warn("user usecase: login error", zap.Error(err))
		return domain.LoginUserResponse{}, fmt.Errorf("user usecase: login error: %v", err.Error())
	}

	return domain.LoginUserResponse{token}, nil
}

func (u *UserUsecase) DummyLogin(ctx context.Context, userType string, lg *zap.Logger) (domain.LoginUserResponse, error) {
	lg.Info("user usecase: dummy login")

	token, err := pkg.GenerateJWTToken(domain.SessionUserID, userType)
	if err != nil {
		lg.Warn("user usecase: login error", zap.Error(err))
		return domain.LoginUserResponse{}, fmt.Errorf("user usecase: login error: %v", err.Error())
	}

	return domain.LoginUserResponse{token}, nil
}
