package user

import "errors"

var (
	ErrEmailAlreadyUsed = errors.New("email already used")
	ErrUserDeleted      = errors.New("user is deleted")
)
