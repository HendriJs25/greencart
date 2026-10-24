package error

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrPasswordIncorrect      = errors.New("password incorrect")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)
