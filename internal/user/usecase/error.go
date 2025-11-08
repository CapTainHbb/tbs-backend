package usecase

import "errors"

var (
	ErrPasswordMismatch 		= errors.New("passwords must be equal")
	ErrPasswordGeneration		= errors.New("failed to generate password")
	ErrUserNotFound 			= errors.New("user not found")
	ErrUsernameAlreadyExists	= errors.New("username already exists")
)