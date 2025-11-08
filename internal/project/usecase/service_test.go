package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/captainhbb/tbs-backend/internal/project/domain"
	portsRepository "github.com/captainhbb/tbs-backend/internal/project/ports"
	portsMock "github.com/captainhbb/tbs-backend/internal/project/ports/mock"
	userUseCaseMock "github.com/captainhbb/tbs-backend/internal/user/usecase/mock"
	userDomain "github.com/captainhbb/tbs-backend/internal/user/domain"
	userUseCase "github.com/captainhbb/tbs-backend/internal/user/usecase"
	"github.com/captainhbb/tbs-backend/internal/project/usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)


func TestCreateProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		input usecase.CreateProjectRequest
		mockSetup func(repo *portsMock.MockRepository, userUserCase *userUseCaseMock.MockUserService)
		expectError bool
		expectedError error
	}{
		{
			name: "success",
			input: usecase.CreateProjectRequest{
				Name: "Test Project1",
				Description: "Test Description",
				StartDate: time.Now(),
				EndDate: time.Now().Add(time.Hour * 24 * 30),
				ProposedBudget: 1000000,
				Status: "active",
				OwnerID: 1,
			},
			mockSetup: func(repo *portsMock.MockRepository, userUserCase *userUseCaseMock.MockUserService) {
				userUserCase.On("GetUser", mock.Anything, 1).Return(userDomain.User{ID: 1}, nil)
				repo.On("CreateProject", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					capturedArg := args.Get(1).(domain.Project)
					require.Equal(t, "Test Project1", capturedArg.Name)
				}).Return(domain.Project{
					ID: 1,
					Name: "Test Project1",
					Description: "Test Description",
					StartDate: time.Now(),
					EndDate: time.Now().Add(time.Hour * 24 * 30),
					ProposedBudget: 1000000,
					Status: "active",
					OwnerID: 1,
				}, nil)
			},
			expectError: false,
		},
		{
			name: "owner user not found",
			input: usecase.CreateProjectRequest{
				Name:           "Test Project2",
				Description:    "Test Description",
				StartDate:      time.Now(),
				EndDate:        time.Now().Add(time.Hour * 24 * 30),
				ProposedBudget: 500000,
				Status:         "active",
				OwnerID:        2,
			},
			mockSetup: func(repo *portsMock.MockRepository, userUserCase *userUseCaseMock.MockUserService) {
				userUserCase.On("GetUser", mock.Anything, 2).Return(userDomain.User{}, userUseCase.ErrUserNotFound)
			},
			expectError: true,
			expectedError: usecase.ErrOwnerNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := portsMock.NewMockRepository(t)
			userServiceMock := userUseCaseMock.NewMockUserService(t)
			service := usecase.New(repoMock, userServiceMock)

			ctx := context.Background()

			tt.mockSetup(repoMock, userServiceMock)

			createdProject, err := service.CreateProject(ctx, tt.input)
			if tt.expectError {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Greater(t, createdProject.ID, 0)
				require.Equal(t, tt.input.Name, createdProject.Name)
				require.Equal(t, tt.input.Description, createdProject.Description)
				require.Equal(t, tt.input.StartDate.Format(time.RFC3339), createdProject.StartDate.Format(time.RFC3339))
				require.Equal(t, tt.input.EndDate.Format(time.RFC3339), createdProject.EndDate.Format(time.RFC3339))
				require.Equal(t, tt.input.ProposedBudget, createdProject.ProposedBudget)
				require.Equal(t, tt.input.Status, createdProject.Status)
				require.Equal(t, tt.input.OwnerID, createdProject.OwnerID)
			}

			userServiceMock.AssertExpectations(t)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestGetProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		input int
		mockSetup func(repo *portsMock.MockRepository)
		expectError bool
		expectedError error
	}{
		{
			name: "success",
			input: 1,
			mockSetup: func(repo *portsMock.MockRepository) {
				repo.On("GetProject", mock.Anything, 1).Return(domain.Project{
					ID: 1,
					Name: "Test Project1",
					Description: "Test Description",
					StartDate: time.Now(),
					EndDate: time.Now().Add(time.Hour * 24 * 30),
					ProposedBudget: 1000000,
					Status: "active",
					OwnerID: 1,
				}, nil)
			},
			expectError: false,
		},
		{
			name: "project not found",
			input: 1,
			mockSetup: func(repo *portsMock.MockRepository) {
				repo.On("GetProject", mock.Anything, 1).Return(domain.Project{}, portsRepository.ErrProjectNotFound)
			},
			expectError: true,
			expectedError: portsRepository.ErrProjectNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := portsMock.NewMockRepository(t)
			userUseCaseMock :=userUseCaseMock.NewMockUserService(t)
			service := usecase.New(repoMock, userUseCaseMock)

			ctx := context.Background()

			tt.mockSetup(repoMock)

			project, err := service.GetProject(ctx, tt.input)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.input, project.ID)
				require.NotEmpty(t, project.Name)
				require.NotEmpty(t, project.Description)
				require.NotEmpty(t, project.StartDate.Format(time.RFC3339))
				require.NotEmpty(t, project.EndDate.Format(time.RFC3339))
				require.NotEmpty(t, project.ProposedBudget)
				require.NotEmpty(t, project.Status)
				require.NotEmpty(t, project.OwnerID)
			}
			repoMock.AssertExpectations(t)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		input usecase.UpdateProjectRequest
		mockSetup func(repo *portsMock.MockRepository, userUserCase *userUseCaseMock.MockUserService)
		expectError bool
		expectedError error
	}{
		{
			name: "success",
			input: usecase.UpdateProjectRequest{
				ID: 1,
				Name: "Test Project1",
				Description: "Test Description",
				StartDate: time.Now(),
				EndDate: time.Now().Add(time.Hour * 24 * 30),
				ProposedBudget: 1000000,
				Status: "active",
				OwnerID: 1,
			},
			mockSetup: func(repo *portsMock.MockRepository, userUserCase *userUseCaseMock.MockUserService) {
				userUserCase.On("GetUser", mock.Anything, 1).Return(userDomain.User{ID: 1}, nil)
				repo.On("UpdateProject", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					capturedArg := args.Get(1).(domain.Project)
					require.Equal(t, 1, capturedArg.ID)
				}).Return(domain.Project{
					ID: 1,
					Name: "Test Project1",
					Description: "Test Description",
					StartDate: time.Now(),
					EndDate: time.Now().Add(time.Hour * 24 * 30),
					ProposedBudget: 1000000,
					Status: "active",
					OwnerID: 1,
				}, nil)
			},
			expectError: false,
		},
		{
			name: "user owner not found",
			input: usecase.UpdateProjectRequest{
				ID: 1,
				Name: "Project Without Valid Owner",
				Description: "Attempt update with missing owner",
				StartDate: time.Now(),
				EndDate: time.Now().Add(time.Hour * 24 * 30),
				ProposedBudget: 500000,
				Status: "active",
				OwnerID: 42, 
			},
			mockSetup: func(_ *portsMock.MockRepository, userUserCase *userUseCaseMock.MockUserService) {
				userUserCase.On("GetUser", mock.Anything, 42).Return(userDomain.User{}, userUseCase.ErrUserNotFound)
			},
			expectError: true,
			expectedError: usecase.ErrOwnerNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := portsMock.NewMockRepository(t)
			userUseCaseMock :=userUseCaseMock.NewMockUserService(t)
			service := usecase.New(repoMock, userUseCaseMock)

			ctx := context.Background()

			tt.mockSetup(repoMock, userUseCaseMock)

			updatedProject, err := service.UpdateProject(ctx, tt.input)
			if tt.expectError {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.input.ID, updatedProject.ID)
				require.Equal(t, tt.input.Name, updatedProject.Name)
				require.Equal(t, tt.input.Description, updatedProject.Description)
				require.Equal(t, tt.input.StartDate.Format(time.RFC3339), updatedProject.StartDate.Format(time.RFC3339))
				require.Equal(t, tt.input.EndDate.Format(time.RFC3339), updatedProject.EndDate.Format(time.RFC3339))
				require.Equal(t, tt.input.ProposedBudget, updatedProject.ProposedBudget)
				require.Equal(t, tt.input.Status, updatedProject.Status)
				require.Equal(t, tt.input.OwnerID, updatedProject.OwnerID)
			}

			userUseCaseMock.AssertExpectations(t)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestDeleteProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		input int
		mockSetup func(repo *portsMock.MockRepository)
		expectError bool
		expectedError error
	}{
		{
			name: "success",
			input: 1,
			mockSetup: func(repo *portsMock.MockRepository) {
				repo.On("DeleteProject", mock.Anything, 1).Return(nil)
			},
			expectError: false,
		},
		{
			name: "project not found",
			input: 1,
			mockSetup: func(repo *portsMock.MockRepository) {
				repo.On("DeleteProject", mock.Anything, 1).Return(portsRepository.ErrProjectNotFound)
			},
			expectError: true,
			expectedError: portsRepository.ErrProjectNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := portsMock.NewMockRepository(t)
			userUseCaseMock :=userUseCaseMock.NewMockUserService(t)
			service := usecase.New(repoMock, userUseCaseMock)

			ctx := context.Background()

			tt.mockSetup(repoMock)


			err := service.DeleteProject(ctx, tt.input)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			repoMock.AssertExpectations(t)
		})
	}
}