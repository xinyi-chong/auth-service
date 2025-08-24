package user

import (
	"auth-service/internal/shared/utils"
	apperrors "auth-service/pkg/error"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	IsUsernameOrEmailRegistered(ctx context.Context, username *string, email string) (bool, error)
	GetUser(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, param *CreateUserParam) error
	UpdateUser(ctx context.Context, id int, param *UpdateUserParam) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, filter *Filter) ([]User, error)
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) IsUsernameOrEmailRegistered(ctx context.Context, username *string, email string) (bool, error) {
	const op = "service.IsUsernameOrEmailRegistered"
	exists, err := s.repo.UsernameOrEmailExists(ctx, username, email)
	if err != nil {
		return false, apperrors.ErrInternalServerError.WithOp(op).Wrap(err)
	}
	return exists, nil
}

func (s *service) GetUser(ctx context.Context, id int) (*User, error) {
	const op = "service.GetUser"
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.HandleDbFindErr(err, apperrors.ModelUser).WithOp(op)
	}
	return user, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	const op = "service.GetUserByEmail"
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, apperrors.HandleDbFindErr(err, apperrors.ModelUser).WithOp(op)
	}
	return user, nil
}

func (s *service) CreateUser(ctx context.Context, param *CreateUserParam) error {
	const op = "service.CreateUser"

	hashedPassword, err := utils.HashPassword(param.Password)
	if err != nil {
		return apperrors.ErrInternalServerError.WithOp(op).Wrap(err)
	}

	user := &User{
		ID:           uuid.New().String(),
		Username:     param.Username,
		Email:        &param.Email,
		PasswordHash: &hashedPassword,
		IsActive:     true,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return apperrors.HandleDbCreateErr(err, apperrors.ModelUser).WithOp(op)
	}

	return nil
}

func (s *service) UpdateUser(ctx context.Context, id int, param *UpdateUserParam) error {
	const op = "service.UpdateUser"

	user := &User{}

	if param.Email != nil {
		user.Email = param.Email
	}

	if param.Password != nil {
		hashedPassword, err := utils.HashPassword(*param.Password)
		if err != nil {
			return apperrors.ErrInternalServerError.WithOp(op).Wrap(err)
		}
		user.PasswordHash = &hashedPassword
	}

	if err := s.repo.Update(ctx, id, user); err != nil {
		return apperrors.HandleDbUpdateErr(err, apperrors.ModelUser).WithOp(op)
	}

	return nil
}

func (s *service) DeleteUser(ctx context.Context, id int) error {
	const op = "service.DeleteUser"
	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.HandleDbDeleteErr(err, apperrors.ModelUser).WithOp(op)
	}
	return nil
}

func (s *service) ListUsers(ctx context.Context, filter *Filter) ([]User, error) {
	const op = "service.ListUsers"
	users, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, apperrors.ErrInternalServerError.WithOp(op).Wrap(err)
	}
	return users, nil
}
