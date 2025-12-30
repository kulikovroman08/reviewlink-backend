package controller

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/response"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"
)

// ReplyToReview godoc
// @Summary      Ответ на отзыв (только для админов)
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        id       path      string                 true  "Review ID"
// @Param        request  body      dto.ReplyToReviewRequest true "Reply content"
// @Success      201  {object}  dto.ReviewReplyResponse
// @Failure      400  {object}  dto.ErrorResponse "invalid input"
// @Failure      401  {object}  dto.ErrorResponse "unauthorized"
// @Failure      403  {object}  dto.ErrorResponse "access denied"
// @Failure      404  {object}  dto.ErrorResponse "review not found"
// @Failure      409  {object}  dto.ErrorResponse "reply already exists"
// @Failure      500  {object}  dto.ErrorResponse "internal error"
// @Router       /admin/reviews/{id}/reply [post]
// @Security     BearerAuth
func (h *Application) ReplyToReview(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})
		return
	}

	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrInvalidUserID})
		return
	}

	reviewID := c.Param("id")
	if _, err := uuid.Parse(reviewID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	var req dto.ReplyToReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	reply, err := h.ReviewReplyService.ReplyToReview(c.Request.Context(), adminID, reviewID, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrAccessDenied):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})
		case errors.Is(err, serviceErrors.ErrReviewNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: response.ErrReviewNotFound})
		case errors.Is(err, serviceErrors.ErrReplyAlreadyExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: response.ErrReplyAlreadyExists})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})
		}
		return
	}

	c.JSON(http.StatusCreated, dto.ReviewReplyResponse{
		Content:   reply.Content,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	})
}

// UpdateReviewReply godoc
// @Summary      Обновить ответ на отзыв (только для админов)
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        id       path      string                 true  "Review ID"
// @Param        request  body      dto.ReplyToReviewRequest true "Reply content"
// @Success      200  {object}  dto.ReviewReplyResponse
// @Failure      400  {object}  dto.ErrorResponse "invalid input"
// @Failure      401  {object}  dto.ErrorResponse "unauthorized"
// @Failure      403  {object}  dto.ErrorResponse "access denied"
// @Failure      404  {object}  dto.ErrorResponse "reply or review not found"
// @Failure      500  {object}  dto.ErrorResponse "internal error"
// @Router       /admin/reviews/{id}/reply [put]
// @Security     BearerAuth
func (h *Application) UpdateReviewReply(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})
		return
	}

	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrInvalidUserID})
		return
	}

	reviewID := c.Param("id")
	if _, err := uuid.Parse(reviewID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	var req dto.ReplyToReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	reply, err := h.ReviewReplyService.UpdateReply(c.Request.Context(), adminID, reviewID, req.Content)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrAccessDenied):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})
		case errors.Is(err, serviceErrors.ErrReviewNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: response.ErrReviewNotFound})
		case errors.Is(err, serviceErrors.ErrReplyNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: response.ErrReplyNotFound})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ReviewReplyResponse{
		Content:   reply.Content,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	})
}
