package usecase

type CreateUserRequest struct {
	Username       string
	FirstName      string
	LastName       string
	Phone          string
	Email          string
	Password       string
	RepeatPassword string
	Role           string
}

type UpdateUserRequest struct {
	ID				int
	Username       string
	FirstName      string
	LastName       string
	Phone          string
	Email          string
	Role           string
}

