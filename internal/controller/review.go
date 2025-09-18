package controller

import (
	"errors"
	"net/http"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
)

// SubmitReview godoc
// @Summary      Отправка отзыва
// @Description  Авторизованный пользователь может оставить отзыв на место, используя одноразовый токен.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SubmitReviewRequest  true  "Данные отзыва"
// @Success      201
// @Failure      400      {object}  dto.ErrorResponse "invalid input или invalid token"
// @Failure      401      {object}  dto.ErrorResponse "invalid user_id"
// @Failure      403      {object}  dto.ErrorResponse "token already used или token expired"
// @Router       /reviews [post]
// @Security     BearerAuth
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

	err = h.ReviewService.SubmitReview(c.Request.Context(), review, req.Token)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrTokenExpired):
			c.JSON(http.StatusForbidden, gin.H{"error": "token expired"})

		case errors.Is(err, serviceErrors.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})

		case errors.Is(err, serviceErrors.ErrInvalidCredentials):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})

		}
		return
	}
	c.Status(http.StatusCreated)
}
