package usecase

import (
	"context"

	"github.com/captainhbb/tbs-backend/internal/user/domain"
	"github.com/captainhbb/tbs-backend/internal/user/ports"
	hash "github.com/captainhbb/tbs-backend/pkg/hash"
)

//go:generate mockery --dir . --name UserService --structname MockUserService --filename mock_user_service.go --output ./mock --outpkg mock
type UserService interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	UpdateUser(ctx context.Context, user UpdateUserRequest) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	repo   ports.Repository
}

func New(repo ports.Repository) UserService {
	return &userService{
		repo: repo,
	}
}

func(s *userService) CreateUser(ctx context.Context, createUserRequest CreateUserRequest) (domain.User, error) {
	if createUserRequest.Password != createUserRequest.RepeatPassword {
		return domain.User{}, ErrPasswordMismatch
	}

	hashedPassword, err := hash.HashPassword(createUserRequest.Password)
	if err != nil {
		return domain.User{}, ErrPasswordGeneration
	}

	user := domain.User{
		Username: createUserRequest.Username,
		FirstName: createUserRequest.FirstName,
		LastName: createUserRequest.LastName,
		Email: createUserRequest.Email,
		Phone: createUserRequest.Phone,
		Role: createUserRequest.Role,
		HashedPassword: hashedPassword,
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	switch err {
	case ports.ErrUsernameAlreadyExists:
		return domain.User{}, ErrUsernameAlreadyExists
	}

	return createdUser, err
}

func(s *userService) GetUser(ctx context.Context, id int) (domain.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	switch err {
	case ports.ErrUserNotFound:
		return domain.User{}, ErrUserNotFound
	case ports.ErrUsernameAlreadyExists:
		return domain.User{}, ErrUsernameAlreadyExists
	}
	return user, err
}

func(s *userService) UpdateUser(ctx context.Context, user UpdateUserRequest) (domain.User, error) {
	updatedUserDomain := domain.User{
		ID: user.ID,
		Username: user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Phone: user.Phone,
		Email: user.Email,
		Role: user.Role,
	}

	updatedUser, err := s.repo.UpdateUser(ctx, updatedUserDomain)
	switch err {
	case ports.ErrUserNotFound:
		return domain.User{}, ErrUserNotFound
	case ports.ErrUsernameAlreadyExists:
		return domain.User{}, ErrUsernameAlreadyExists
	}
	return updatedUser, err
}

func(s *userService) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.DeleteUser(ctx, id)
	switch err {
	case ports.ErrUserNotFound:
		return ErrUserNotFound
	}
	return err
}

