package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	UserRepo UserRepository
}

func NewHandler(userRepo UserRepository) *Handler {
	return &Handler{UserRepo: userRepo}
}

func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	existingUser, err := h.UserRepo.FindByEmail(c.Request.Context(), req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	newUser := &User{
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  newUser.ID,
		"role": newUser.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("supersecret")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, AuthResponse{Token: tokenString})
}
