package errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyUsed   = errors.New("email already used")
	ErrPlaceAlreadyExists = errors.New("place already exists")
	ErrInvalidPlaceData   = errors.New("invalid place data")
	ErrInvalidPlaceID     = errors.New("invalid place id")
	ErrValidation         = errors.New("validation error")
	ErrInvalidRating      = errors.New("invalid rating")
	ErrTokenRequired      = errors.New("token required")
	ErrTokenUsed          = errors.New("token is used")
	ErrAlreadyReviewed    = errors.New("review already submitted today for this place")
)
