package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/service"
	err_svc "github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

type Handler struct {
	UserService service.UserService
}

func NewHandler(service service.UserService) *Handler {
	return &Handler{UserService: service}
}

func (h *Handler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	token, err := h.UserService.Signup(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AuthResponse{Token: token})
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	token, err := h.UserService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.AuthResponse{Token: token})
}

func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := h.UserService.GetMe(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
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

func (h *Handler) UpdateMe(c *gin.Context) {
	userID := c.GetString("user_id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	req.UserID = userID

	updatedUser, err := h.UserService.UpdateMe(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, err_svc.ErrEmailAlreadyUsed) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) DeleteMe(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.UserService.DeleteMe(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
