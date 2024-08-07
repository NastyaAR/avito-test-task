package usecase

import (
	"avito-test-task/internal/domain"
	"avito-test-task/pkg"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/mail"
)

type UserUsecase struct {
	userRepo domain.UserRepo
}

func NewUserUsecase(userRepo domain.UserRepo) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidUserType(userType string) bool {
	return userType == domain.Client || userType == domain.Moderator
}

func (u *UserUsecase) Register(ctx context.Context, userReq *domain.RegisterUserRequest, lg *zap.Logger) (domain.RegisterUserResponse, error) {
	lg.Info("user usecase: register", zap.String("user_type", userReq.UserType))

	if userReq == nil {
		lg.Warn("user usecase: register error: bad nil request")
		return domain.RegisterUserResponse{},
			fmt.Errorf("user user: register error: %w", domain.ErrUser_BadRequest)
	}

	if !isValidUserType(userReq.UserType) {
		lg.Warn("user usecase: register error: bad role", zap.String("role", userReq.UserType))
		return domain.RegisterUserResponse{},
			fmt.Errorf("user usecase: register error: %w", domain.ErrUser_BadType)
	}

	if !isValidEmail(userReq.Email) {
		lg.Warn("user usecase: register error: bad mail", zap.String("mail", userReq.Email))
		return domain.RegisterUserResponse{},
			fmt.Errorf("user usecase: register error: %w", domain.ErrUser_BadMail)
	}

	if userReq.Password == "" {
		lg.Warn("user usecase: register error: bad empty password")
		return domain.RegisterUserResponse{},
			fmt.Errorf("user usecase: register error: %w", domain.ErrUser_BadPassword)
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

	if userReq == nil {
		lg.Warn("user usecase: login error: bad nil request")
		return domain.LoginUserResponse{},
			fmt.Errorf("user usecase: login error: %w", domain.ErrUser_BadRequest)
	}

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

	if !isValidUserType(userType) {
		lg.Warn("user usecase: dummy login error: bad role", zap.String("role", userType))
		return domain.LoginUserResponse{},
			fmt.Errorf("user usecase: dummy login error: %w", domain.ErrUser_BadType)
	}

	token, err := pkg.GenerateJWTToken(domain.SessionUserID, userType)
	if err != nil {
		lg.Warn("user usecase: login error", zap.Error(err))
		return domain.LoginUserResponse{}, fmt.Errorf("user usecase: login error: %v", err.Error())
	}

	return domain.LoginUserResponse{token}, nil
}
