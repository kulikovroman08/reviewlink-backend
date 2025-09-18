package controller

import (
	"errors"
	"net/http"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
)

func (h *Application) Signup(c *gin.Context) {
	var req dto.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	token, err := h.UserService.Signup(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrEmailAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to signup"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{Token: token})
}

func (h *Application) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	token, err := h.UserService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})

		case errors.Is(err, serviceErrors.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{Token: token})
}

func (h *Application) GetUser(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.UserService.GetUser(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		}
		return
	}

	resp := dto.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		Points: user.Points,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Application) UpdateUser(c *gin.Context) {
	userID := c.GetString("user_id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	req.UserID = userID

	hasUpdate := false
	if req.Name != nil {
		hasUpdate = true
	}
	if req.Email != nil {
		hasUpdate = true
	}
	if req.Password != nil {
		hasUpdate = true
	}
	if !hasUpdate {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one field must be provided"})
		return
	}

	var user model.User
	user.ID = req.UserID
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	var password string
	if req.Password != nil {
		password = *req.Password
	}

	updatedUser, err := h.UserService.UpdateUser(c.Request.Context(), user, password)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})

		case errors.Is(err, serviceErrors.ErrEmailAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (h *Application) DeleteUser(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.UserService.DeleteUser(c.Request.Context(), userID); err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
