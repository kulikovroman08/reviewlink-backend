package dto

import (
	"time"

	"github.com/google/uuid"
)

// Структура входного запроса
type SignupRequest struct {
	Name     string `json:"name" binding:"required,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	UserID   string  `json:"-"`
	Name     *string `json:"name" binding:"omitempty,max=100"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=6"`
}

type UserResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Points int    `json:"points"`
}
type CreatePlaceRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type CreatePlaceResponse struct {
	ID string `json:"id"`
}

type SubmitReviewRequest struct {
	Token   string    `json:"token" binding:"required"`
	PlaceID uuid.UUID `json:"place_id" binding:"required"`
	Rating  int       `json:"rating" binding:"required,min=1,max=5"`
	Content string    `json:"content"`
}

type UpdateReviewRequest struct {
	Content string `json:"content"`
	Rating  int    `json:"rating"`
}

type ReviewResponse struct {
	Rating    int       `json:"rating"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type GenerateTokensRequest struct {
	PlaceID string `json:"place_id" binding:"required,uuid"`
	Count   int    `json:"count" binding:"required,min=1,max=100"`
}

type GenerateTokensResponse struct {
	Tokens []string `json:"tokens"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type AdminStatsResponse struct {
	TotalUsers    int     `json:"total_users"`
	TotalReviews  int     `json:"total_reviews"`
	AverageRating float64 `json:"average_rating"`
	TotalBonuses  int     `json:"total_bonuses"`
}

type LeaderboardEntry struct {
	Rank         int     `json:"rank"`
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	ReviewsCount int     `json:"reviews_count"`
	AvgRating    float64 `json:"avg_rating"`
}

type BonusRedeemResponse struct {
	ID             string     `json:"id"`
	PlaceID        string     `json:"place_id"`
	RewardType     string     `json:"reward_type"`
	RequiredPoints int        `json:"required_points"`
	QRToken        string     `json:"qr_token"`
	IsUsed         bool       `json:"is_used"`
	UsedAt         *time.Time `json:"used_at,omitempty"`
}

type BonusRedeemRequest struct {
	PlaceID    string `json:"place_id" binding:"required,uuid"`
	RewardType string `json:"reward_type" binding:"required,oneof=free_coffee free_meal discount_10"`
}

type BonusResponse struct {
	ID             string     `json:"id"`
	PlaceID        string     `json:"place_id"`
	RewardType     string     `json:"reward_type"`
	RequiredPoints int        `json:"required_points"`
	QRToken        string     `json:"qr_token"`
	IsUsed         bool       `json:"is_used"`
	UsedAt         *time.Time `json:"used_at,omitempty"`
}

type BonusValidateRequest struct {
	QRToken string `json:"qr_token" binding:"required"`
}

type BonusValidateResponse struct {
	Status string `json:"status"`
}

type UserRestrictionResponse struct {
	RestrictionType string    `json:"restriction_type"`
	Reason          string    `json:"reason"`
	ExpiresAt       time.Time `json:"expires_at"`
}

type UserStatsResponse struct {
	TotalReviews  int     `json:"total_reviews"`
	AvgRating     float64 `json:"avg_rating"`
	Points        int     `json:"points"`
	BonusesActive int     `json:"bonuses_active"`
	BonusesUsed   int     `json:"bonuses_used"`
}
