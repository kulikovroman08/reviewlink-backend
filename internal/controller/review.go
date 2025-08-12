package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
)

func (a *Application) SubmitReview(c *gin.Context) {
	var req dto.SubmitReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
		return
	}

	err = a.ReviewService.SubmitReview(c.Request.Context(), userID.String(), req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "token already used"):
			c.JSON(http.StatusForbidden, gin.H{"error": "token already used"})

		case strings.Contains(err.Error(), "token expired"):
			c.JSON(http.StatusForbidden, gin.H{"error": "token expired"})

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})

		}
		return
	}
	c.Status(http.StatusCreated)
}
