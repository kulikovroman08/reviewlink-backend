package response

// Общие ошибки
const (
	ErrInvalidInput       = "invalid input"
	ErrUnauthorized       = "authentication required"
	ErrAtLeastOneField    = "at least one field must be provided"
	ErrInvalidCredentials = "invalid credentials"
	ErrInternalError      = "internal error"
)

// Users
const (
	ErrUserNotFound       = "user not found"
	ErrEmailAlreadyExists = "email already in use"
	ErrFailedSignup       = "failed to signup"
	ErrFailedLogin        = "login failed"
	ErrFailedGetUser      = "failed to get user"
	ErrFailedUpdateUser   = "failed to update user"
	ErrFailedDeleteUser   = "failed to delete user"
	ErrUserDeleted        = "user deleted"
)

// Tokens
const (
	ErrOnlyAdminCanGenerateTokens = "only admin can generate tokens"
	ErrFailedGenerateTokens       = "failed to generate tokens"
	ErrInvalidPlaceID             = "invalid place id"
	ErrInvalidToken               = "invalid token"
	ErrTokenExpired               = "token expired"
	ErrTokenAlreadyUsed           = "token already used"
)

// Places
const (
	ErrAccessDenied       = "access denied"
	ErrPlaceAlreadyExists = "place already exists"
	ErrInvalidPlaceData   = "invalid place data"
	ErrFailedCreatePlace  = "failed to create place"
	ErrPlaceNotFound      = "place not found"
)

// Reviews
const (
	ErrInvalidUserID      = "invalid user_id"
	ErrInvalidRating      = "invalid rating"
	ErrReviewNotFound     = "review not found"
	ErrFailedUpdateReview = "failed to update review"
)

// Admin
const (
	ErrFailedLoadStats = "failed to load stats"
)

// Bonuses
const (
	ErrNotEnoughPoints   = "not enough points"
	ErrFailedCreateBonus = "failed to create bonus"
	ErrBonusNotFound     = "bonus not found"
	ErrBonusAlreadyUsed  = "bonus already used"
)
