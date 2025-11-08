package ports

import (
	"context"

	"github.com/captainhbb/tbs-backend/internal/project/domain"
)

//go:generate mockery --dir . --name Repository --structname MockRepository --filename mock_repository.go --output ./mock --outpkg mock
type Repository interface {
	CreateProject(ctx context.Context, project domain.Project) (domain.Project, error)
	GetProject(ctx context.Context, id int) (domain.Project, error)
	UpdateProject(ctx context.Context, project domain.Project) (domain.Project, error)
	DeleteProject(ctx context.Context, id int) error
}