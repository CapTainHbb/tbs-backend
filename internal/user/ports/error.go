package ports

import "errors"


var (
	ErrUsernameAlreadyExists		= errors.New("username already exists")
	ErrUserNotFound					= errors.New("user not found")
)