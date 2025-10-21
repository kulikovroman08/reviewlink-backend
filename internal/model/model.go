package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         string
	Points       int
	CreatedAt    time.Time
	IsDeleted    bool
}

type Place struct {
	ID        uuid.UUID
	Name      string
	Address   string
	CreatedAt time.Time
	IsDeleted bool
}

type ReviewToken struct {
	ID        uuid.UUID
	PlaceID   uuid.UUID
	Token     string
	IsUsed    bool
	ExpiresAt time.Time
}

type Review struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PlaceID   uuid.UUID
	TokenID   uuid.UUID
	Content   string
	Rating    int
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type ReviewFilter struct {
	Rating    int
	HasRating bool
	Sort      string
	FromDate  *time.Time
	ToDate    *time.Time
}

type GenerateTokensResult struct {
	Tokens []string
}

type AdminStats struct {
	TotalUsers    int
	TotalReviews  int
	AverageRating float64
	TotalBonuses  int
}

type LeaderboardEntry struct {
	ID           string
	Name         string
	ReviewsCount int
	AvgRating    float64
}

type LeaderboardFilter struct {
	SortBy     string
	MinRating  float64
	MinReviews int
}
