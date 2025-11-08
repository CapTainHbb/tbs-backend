package usecase

import (
	"context"

	"github.com/captainhbb/tbs-backend/internal/project/domain"
	"github.com/captainhbb/tbs-backend/internal/project/ports"
	userUseCase "github.com/captainhbb/tbs-backend/internal/user/usecase"
)


type ProjectService interface {
	CreateProject(ctx context.Context, project CreateProjectRequest) (domain.Project, error)
	GetProject(ctx context.Context, id int) (domain.Project, error)
	UpdateProject(ctx context.Context, project UpdateProjectRequest) (domain.Project, error)
	DeleteProject(ctx context.Context, id int) error
}

type projectService struct {
	repo ports.Repository
	userService userUseCase.UserService
}

func New(repo ports.Repository, userService userUseCase.UserService) ProjectService {
	return &projectService{
		repo: repo,
		userService: userService,
	}
}

func(s *projectService) CreateProject(ctx context.Context, createProjectRequest CreateProjectRequest) (domain.Project, error) {
	_, err := s.userService.GetUser(ctx, createProjectRequest.OwnerID)
	switch err {
	case userUseCase.ErrUserNotFound:
		return domain.Project{}, ErrOwnerNotFound
	}
	
	project := domain.Project{
		Name: createProjectRequest.Name,
		Description: createProjectRequest.Description,
		StartDate: createProjectRequest.StartDate,
		EndDate: createProjectRequest.EndDate,
		ProposedBudget: createProjectRequest.ProposedBudget,
		Status: createProjectRequest.Status,
		OwnerID: createProjectRequest.OwnerID,
	}
	return s.repo.CreateProject(ctx, project)
}

func(s *projectService) GetProject(ctx context.Context, id int) (domain.Project, error) {
	return s.repo.GetProject(ctx, id)
}

func(s *projectService) UpdateProject(ctx context.Context, updateProjectRequest UpdateProjectRequest) (domain.Project, error) {
	_, err := s.userService.GetUser(ctx, updateProjectRequest.OwnerID)
	switch err {
	case userUseCase.ErrUserNotFound:
		return domain.Project{}, ErrOwnerNotFound
	}

	project := domain.Project{
		ID: updateProjectRequest.ID,
		Name: updateProjectRequest.Name,
		Description: updateProjectRequest.Description,
		StartDate: updateProjectRequest.StartDate,
		EndDate: updateProjectRequest.EndDate,
		ProposedBudget: updateProjectRequest.ProposedBudget,
	}
	return s.repo.UpdateProject(ctx, project)
}

func(s *projectService) DeleteProject(ctx context.Context, id int) error {
	return s.repo.DeleteProject(ctx, id)
}