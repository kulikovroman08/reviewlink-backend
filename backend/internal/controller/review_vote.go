package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/response"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"
)

// VoteReview godoc
// @Summary      Оценка полезности отзыва
// @Description  Авторизованный пользователь может поставить голос "полезно" (+1) или "неполезно" (-1). Повторный такой же голос — noop.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        id       path   string               true  "Review ID"
// @Param        request  body   dto.VoteReviewRequest true "Vote value: 1 or -1"
// @Success      200
// @Failure      400 {object} dto.ErrorResponse "invalid input / invalid vote value"
// @Failure      401 {object} dto.ErrorResponse "invalid user_id / unauthorized"
// @Failure      403 {object} dto.ErrorResponse "review not found / access denied"
// @Failure      500 {object} dto.ErrorResponse "internal error"
// @Router       /reviews/{id}/vote [post]
// @Security     BearerAuth
func (h *Application) VoteReview(c *gin.Context) {
	reviewID := c.Param("id")

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrInvalidUserID})
		return
	}

	var req dto.VoteReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	err := h.ReviewVoteService.Vote(c.Request.Context(), reviewID, userID, req.Value)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidVoteValue):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})

		case errors.Is(err, serviceErrors.ErrAccessDenied):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})

		case errors.Is(err, serviceErrors.ErrReviewNotFound):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrReviewNotFound})

		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})
		}
		return
	}

	c.Status(http.StatusOK)
}

// DeleteReviewVote godoc
// @Summary      Удаление голоса с отзыва
// @Description  Авторизованный пользователь может удалить свой голос с отзыва. Если голоса не было — noop.
// @Tags         reviews
// @Produce      json
// @Param        id   path  string  true  "Review ID"
// @Success      204
// @Failure      401 {object} dto.ErrorResponse "invalid user_id / unauthorized"
// @Failure      403 {object} dto.ErrorResponse "review not found / access denied"
// @Failure      500 {object} dto.ErrorResponse "internal error"
// @Router       /reviews/{id}/vote [delete]
// @Security     BearerAuth
func (h *Application) DeleteReviewVote(c *gin.Context) {
	reviewID := c.Param("id")

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrInvalidUserID})
		return
	}

	err := h.ReviewVoteService.Unvote(c.Request.Context(), reviewID, userID)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrAccessDenied):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})

		case errors.Is(err, serviceErrors.ErrReviewNotFound):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrReviewNotFound})

		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
