package usecase_test

import (
	"context"
	"testing"

	"github.com/captainhbb/tbs-backend/internal/user/domain"
	"github.com/captainhbb/tbs-backend/internal/user/ports"
	portsMock "github.com/captainhbb/tbs-backend/internal/user/ports/mock"
	"github.com/captainhbb/tbs-backend/internal/user/usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          	string
		input         	usecase.CreateUserRequest
		mockSetup     	func(repo *portsMock.MockRepository)
		expectError 	bool
		expectedError   error
	}{
		{
			name: "successfully creates user",
			input: usecase.CreateUserRequest{
				Username:       "testuser1",
				FirstName:      "Hossein",
				LastName:       "Beiranvahd",
				Phone:          "+989399915084",
				Email:          "hossein1377075@gmail.com",
				Password:       "capitanhb12345",
				RepeatPassword: "capitanhb12345",
				Role:           "admin",
			},
			mockSetup: func(repo *portsMock.MockRepository) {
				repo.On("CreateUser", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						capturedArg := args.Get(1).(domain.User)
						require.Equal(t, "testuser1", capturedArg.Username)
					}).
					Return(domain.User{
						ID:             1,
						Username:       "testuser1",
						FirstName:      "Hossein",
						LastName:       "Beiranvahd",
						Phone:          "+989399915084",
						Email:          "hossein1377075@gmail.com",
						HashedPassword: "somerandomhashing",
						Role:           "admin",
					}, nil).Once()
			},
			expectError: false,
		},
		{
			name: "passwords do not match",
			input: usecase.CreateUserRequest{
				Username:       "captainhb",
				FirstName:      "Hossein",
				LastName:       "Beiranvahd",
				Phone:          "+989399915084",
				Email:          "hossein1377075@gmail.com",
				Password:       "abc123",
				RepeatPassword: "xyz123",
				Role:           "admin",
			},
			mockSetup:     	func(repo *portsMock.MockRepository) {}, // no DB call expected
			expectError: 	true,
			expectedError: 	usecase.ErrPasswordMismatch,
		},
		{
			name: "dublicate username",
			input: usecase.CreateUserRequest{
				Username:       "testuser2",
				FirstName:      "Hossein",
				LastName:       "Beiranvahd",
				Phone:          "+989399915084",
				Email:          "hossein1377075@gmail.com",
				Password:       "abc123",
				RepeatPassword: "abc123",
				Role:           "admin",
			},
			mockSetup: func(repo *portsMock.MockRepository) {
				repo.On("CreateUser", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						capturedArg := args.Get(1).(domain.User)
						require.Equal(t, "testuser2", capturedArg.Username)
					}).Return(domain.User{}, ports.ErrUsernameAlreadyExists).Once()
			},
			expectError: 	true,
			expectedError: 	usecase.ErrUsernameAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := portsMock.NewMockRepository(t)
			service := usecase.New(repo)

			ctx := context.Background()
			tt.mockSetup(repo)

			createdUser, err := service.CreateUser(ctx, tt.input)

			if tt.expectError {
				require.ErrorIs(t, err, tt.expectedError)
				require.Zero(t, createdUser.ID)
			} else {
				require.NoError(t, err)
				require.Greater(t, createdUser.ID, 0)
				require.Equal(t, tt.input.Username, createdUser.Username)
				require.Equal(t, tt.input.FirstName, createdUser.FirstName)
				require.Equal(t, tt.input.LastName, createdUser.LastName)
				require.Equal(t, tt.input.Phone, createdUser.Phone)
				require.Equal(t, tt.input.Email, createdUser.Email)
				require.Equal(t, tt.input.Role, createdUser.Role)
				require.NotEmpty(t, createdUser.HashedPassword)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name      			string
		input 				int
		mockSetup   		func(*portsMock.MockRepository)
		expectError 		bool
		expectedError  		error
	} {
		{
			name: "successfull",
			input: 1,
			mockSetup: func(mockRepo *portsMock.MockRepository) {
				mockRepo.On("GetUser", mock.Anything, 1).Return(domain.User{
					ID: 1,
					Username: "testuser1",
					FirstName: "hossein",
					LastName: "Beiranvand",
					Phone: "+989399915084",
					Email: "hossein1377075@gmail.com",
					HashedPassword: "hashedpassword",
					Role: "admin",
				}, nil)
			},
		},
		{
			name: "not found",
			input: 2,
			mockSetup: func(mockRepo *portsMock.MockRepository) {
				mockRepo.On("GetUser", mock.Anything, 2).Return(domain.User{}, ports.ErrUserNotFound)
			},
			expectError: true,
			expectedError: usecase.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		repo := portsMock.NewMockRepository(t)
		service := usecase.New(repo)
		ctx := context.Background()

		tt.mockSetup(repo)

		user, err := service.GetUser(ctx, tt.input)
		if tt.expectError {
			require.ErrorIs(t, tt.expectedError, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, 1, user.ID)
			require.NotEmpty(t, user.Username)
			require.NotEmpty(t, user.FirstName)
			require.NotEmpty(t, user.LastName)
			require.NotEmpty(t, user.Phone)
			require.NotEmpty(t, user.Email)
			require.NotEmpty(t, user.Role)
			require.NotEmpty(t, user.HashedPassword)
		}

		repo.AssertExpectations(t)
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name     		string
		input           usecase.UpdateUserRequest
		mockSetup       func(*portsMock.MockRepository)
		expectError     bool
		expectedError   error
	} {
		{
			name: "success",
			input: usecase.UpdateUserRequest{
				ID: 1,
				Username: "testuser1",
				FirstName: "Hossein",
				LastName: "Beiranvand",
				Phone: "+989399915084",
				Email: "hossein1377075@gmail.com",
				Role: "admin",
			},
			mockSetup: func(repo *portsMock.MockRepository) {
				repo.On("UpdateUser", mock.Anything, domain.User{
					ID: 1,
					Username: "testuser1",
					FirstName: "Hossein",
					LastName: "Beiranvand",
					Phone: "+989399915084",
					Email: "hossein1377075@gmail.com",
					Role: "admin",
				}).Run(func(args mock.Arguments) {
					capturedArg := args.Get(1).(domain.User)
					require.Equal(t, "testuser1", capturedArg.Username)
				}).Return(domain.User{
					ID: 1,
					Username: "testuser1",
					FirstName: "Hossein",
					LastName: "Beiranvand",
					Phone: "+989399915084",
					Email: "hossein1377075@gmail.com",
					Role: "admin",
					HashedPassword: "somehashedpassword",
				}, nil)
			},
			expectError: false,
		},
		{
			name: "username already exists",
			input: usecase.UpdateUserRequest{
				ID: 1,
				Username: "testuser2",
				FirstName: "Hossein",
				LastName: "Beiranvand",
				Phone: "+989399915084",
				Email: "hossein1377075@gmail.com",
				Role: "admin",
			},
			mockSetup: func(repo *portsMock.MockRepository) {
				repo.On("UpdateUser", mock.Anything, domain.User{
					ID: 1,
					Username: "testuser2",
					FirstName: "Hossein",
					LastName: "Beiranvand",
					Phone: "+989399915084",
					Email: "hossein1377075@gmail.com",
					Role: "admin",
				}).Run(func(args mock.Arguments) {
					capturedArg := args.Get(1).(domain.User)
					require.Equal(t, "testuser2", capturedArg.Username)
				}).Return(domain.User{}, ports.ErrUsernameAlreadyExists)
			},
			expectError: true,
			expectedError: usecase.ErrUsernameAlreadyExists,
		},
	}
	
	for _, tt := range tests {
		repo := portsMock.NewMockRepository(t)
		service := usecase.New(repo)
		
		ctx := context.Background()

		tt.mockSetup(repo)

		updatedUser, err := service.UpdateUser(ctx, tt.input)
		if tt.expectError {
			require.ErrorIs(t, err, tt.expectedError)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.input.ID, updatedUser.ID)
			require.Equal(t, tt.input.Username, updatedUser.Username)
			require.Equal(t, tt.input.FirstName, updatedUser.FirstName)
			require.Equal(t, tt.input.LastName, updatedUser.LastName)
			require.Equal(t, tt.input.Phone, updatedUser.Phone)
			require.Equal(t, tt.input.Email, updatedUser.Email)
			require.Equal(t, tt.input.Role, updatedUser.Role)
		}
		repo.AssertExpectations(t)
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name 			string
		input       	int
		mockSetup   	func(*portsMock.MockRepository)
		expectError 	bool
		expectedError   error
	} {
		{
			name: "success",
			input: 1,
			mockSetup: func(repo *portsMock.MockRepository) {
				repo.On("DeleteUser", mock.Anything, 1).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name: "error",
			input: 2,
			mockSetup: func(repo *portsMock.MockRepository) {
				repo.On("DeleteUser", mock.Anything, 2).Return(ports.ErrUserNotFound).Once()
			},
			expectError: true,
			expectedError: usecase.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		repo := portsMock.NewMockRepository(t)
		service := usecase.New(repo)
		
		ctx := context.Background()

		tt.mockSetup(repo)

		err := service.DeleteUser(ctx, tt.input)
		if tt.expectError {
			require.ErrorIs(t, err, tt.expectedError)
		} else {
			require.NoError(t, err)
		}

		repo.AssertExpectations(t)
	}
}