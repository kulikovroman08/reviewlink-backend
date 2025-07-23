package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kulikovroman08/reviewlink-backend/internal/auth"
	"github.com/kulikovroman08/reviewlink-backend/internal/auth/repository"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	UserRepo repository.UserRepository
}

func NewHandler(userRepo repository.UserRepository) *Handler {
	return &Handler{UserRepo: userRepo}
}

func GenerateJWT(user *auth.User) (string, error) {
	claims := auth.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}
	return tokenString, nil
}

func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	existingUser, err := h.UserRepo.FindByEmail(c.Request.Context(), req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("failed to check user existence", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("password hashing failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	newUser := &auth.User{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		Points:       0,
		CreatedAt:    time.Now(),
		IsDeleted:    false,
	}
	if err := h.UserRepo.CreateUser(c.Request.Context(), newUser); err != nil {
		slog.Error("failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
	}

	token, err := GenerateJWT(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, AuthResponse{Token: token})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("login bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user, err := h.UserRepo.FindByEmail(c.Request.Context(), req.Email)
	if err != nil || user == nil || user.IsDeleted {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if user.IsDeleted {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}

func (h *Handler) GetProfiel(c *gin.Context) {
	userID, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.UserRepo.FindByID(c.Request.Context(), userID.(string))
	if err != nil {
		slog.Error("failed to find user by ID", "user_id", userID, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":     userID,
		"name":   user.Name,
		"email":  user.Email,
		"role":   user.Role,
		"points": user.Points,
	})
}
func (h *Handler) CreateReview(c *gin.Context) {
	userID, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Text   string `json:"text"`
		Rating int    `json:"rating"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "review received",
		"user_id": userID,
		"text":    req.Text,
		"rating":  req.Rating,
	})
}
