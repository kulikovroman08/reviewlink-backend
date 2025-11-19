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
	ErrPlaceNotFound      = errors.New("place not found")
	ErrInvalidRating      = errors.New("invalid rating value")
	ErrReviewNotFound     = errors.New("review not found or access denied")
	ErrNotEnoughPoints    = errors.New("not enough points to redeem bonus")
	ErrBonusCreateFail    = errors.New("failed to create bonus")
	ErrBonusNotFound      = errors.New("bonus not found")
	ErrBonusAlreadyUsed   = errors.New("bonus already used")
	ErrTooManyReviews     = errors.New("too many reviews today")
)
