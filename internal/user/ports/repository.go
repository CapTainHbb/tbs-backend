package ports

import (
	"context"

	"github.com/captainhbb/tbs-backend/internal/user/domain"
)

//go:generate mockery --dir . --name Repository --structname MockRepository --filename mock_repository.go --output ./mock --outpkg mock
type Repository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

