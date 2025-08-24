package auth

import (
	"auth-service/internal/shared/utils"
	userModel "auth-service/internal/user"
	apperrors "auth-service/pkg/error"
	"auth-service/pkg/jwt"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	Register(ctx context.Context, param RegisterParam) error
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	RefreshToken(ctx context.Context, refreshToken string) (*Tokens, error)
}

type service struct {
	logger  *zap.Logger
	userSvc userModel.Service
}

func NewService(userSvc userModel.Service, logger *zap.Logger) Service {
	return &service{userSvc: userSvc, logger: logger}
}

func (s *service) Register(ctx context.Context, param RegisterParam) error {
	const op = "service.Register"

	exists, err := s.userSvc.IsUsernameOrEmailRegistered(ctx, param.Username, param.Email)
	if err != nil {
		return err
	} else if exists {
		return apperrors.ErrUserConflict.WithOp(op)
	}

	user := &userModel.CreateUserParam{
		Username: param.Username,
		Email:    param.Email,
		Password: param.Password,
	}

	err = s.userSvc.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	// TODO: Verify email (OTP)

	return nil
}

func (s *service) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	const op = "service.Login"

	user, err := s.userSvc.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	isValid, err := utils.IsPasswordMatch(*user.PasswordHash, password)
	if err != nil {
		return nil, apperrors.ErrInternalServerError.WithOp(op).Wrap(err)
	} else if !isValid {
		return nil, apperrors.ErrIncorrectPassword.WithOp(op)
	}

	accessToken, err := token.GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, apperrors.ErrInternalServerError.WithOp(op).Wrap(err)
	}

	refreshToken, err := token.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.ErrInternalServerError.WithOp(op).Wrap(err)
	}

	return &LoginResponse{
		Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, UserClaims{UserID: user.ID, Email: user.Email, Username: user.Username},
	}, nil
}

func (s *service) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	const op = "service.ChangePassword"

	user, err := s.userSvc.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	isValid, err := utils.IsPasswordMatch(*user.PasswordHash, oldPassword)
	if err != nil {
		return apperrors.ErrInternalServerError.WithOp(op).Wrap(err)
	} else if !isValid {
		return apperrors.ErrIncorrectPassword.WithOp(op)
	}

	param := &userModel.UpdateUserParam{
		Password: &newPassword,
	}

	return s.userSvc.UpdateUser(ctx, userID, param)
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*Tokens, error) {
	const op = "service.RefreshToken"

	claims, err := token.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, apperrors.ErrSessionExpired.WithOp(op).Wrap(err)
	}

	user, err := s.userSvc.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	err = token.InvalidateToken(ctx, refreshToken)
	if err != nil {
		s.logger.Warn("failed to blacklist old refresh token", zap.Error(err))
	}

	accessToken, err := token.GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, apperrors.ErrInternalServerError.WithOp(op).Wrap(err)
	}

	newRefreshToken, err := token.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.ErrInternalServerError.WithOp(op).Wrap(err)
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
