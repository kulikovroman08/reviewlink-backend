package controller

import (
	"net/http"
	"strings"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
)

func (h *Application) SubmitReview(c *gin.Context) {
	var req dto.SubmitReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
		return
	}

	review := model.Review{
		UserID:  userID,
		PlaceID: req.PlaceID,
		Content: req.Content,
		Rating:  req.Rating,
	}

	err = h.Service.SubmitReview(c.Request.Context(), review, req.Token)
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
