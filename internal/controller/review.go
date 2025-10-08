package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/response"
)

// SubmitReview godoc
// @Summary      Отправка отзыва
// @Description  Авторизованный пользователь может оставить отзыв на место, используя одноразовый токен.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SubmitReviewRequest  true  "Данные отзыва"
// @Success      201
// @Failure 400 {object} dto.ErrorResponse "invalid input"
// @Failure 401 {object} dto.ErrorResponse "invalid user_id / invalid token"
// @Failure 403 {object} dto.ErrorResponse "token expired / token already used"
// @Failure 500 {object} dto.ErrorResponse "internal error"
// @Router       /reviews [post]
// @Security     BearerAuth
func (h *Application) SubmitReview(c *gin.Context) {
	var req dto.SubmitReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrInvalidUserID})
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
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrTokenExpired})

		case errors.Is(err, serviceErrors.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrInvalidToken})

		case errors.Is(err, serviceErrors.ErrInvalidCredentials):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidCredentials})

		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})

		}
		return
	}
	c.Status(http.StatusCreated)
}

// GetReviews godoc
// @Summary      Просмотр отзывов по заведению
// @Description  Получение списка отзывов по placeID с фильтрацией и сортировкой
// @Tags         places
// @Accept       json
// @Produce      json
// @Param        id     path      string  true   "Place ID"
// @Param        rating query     int     false  "Фильтр по рейтингу (1-5)"
// @Param        sort   query     string  false  "Сортировка: date_asc или date_desc"
// @Success 200 {array} dto.ReviewResponse
// @Failure 400 {object} dto.ErrorResponse "invalid input"
// @Failure 404 {object} dto.ErrorResponse "place not found"
// @Failure 500 {object} dto.ErrorResponse "internal error"
// @Router /places/{id}/reviews [get]
func (h *Application) GetReviews(c *gin.Context) {
	placeID := c.Param("id")

	filter, err := parseReviewFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	reviews, err := h.ReviewService.GetReviews(c.Request.Context(), placeID, filter)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrPlaceNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: response.ErrPlaceNotFound})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})
		}
		return
	}

	resp := make([]dto.ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		resp = append(resp, dto.ReviewResponse{
			Rating:    r.Rating,
			Content:   r.Content,
			CreatedAt: r.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func parseReviewFilter(c *gin.Context) (model.ReviewFilter, error) {
	var f model.ReviewFilter

	if r := c.Query("rating"); r != "" {
		val, err := strconv.Atoi(r)
		if err != nil {
			return f, fmt.Errorf("invalid rating: %w", err)
		}
		f.Rating = val
		f.HasRating = true
	}

	f.Sort = c.DefaultQuery("sort", "date_desc")

	if from := c.Query("from"); from != "" {
		t, err := time.Parse("2006-01-02", from)
		if err != nil {
			return f, fmt.Errorf("invalid from date: %w", err)
		}
		f.FromDate = &t
	}

	if to := c.Query("to"); to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err != nil {
			return f, fmt.Errorf("invalid to date: %w", err)
		}
		f.ToDate = &t
	}

	return f, nil
}

// UpdateReview godoc
// @Summary      Редактирование отзыва
// @Description  Автор отзыва может изменить контент и рейтинг
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Review ID"
// @Param        request body   dto.UpdateReviewRequest true "Данные для обновления отзыва"
// @Success      200  {object}  dto.MessageResponse "review updated successfully"
// @Failure      400  {object}  dto.ErrorResponse "invalid input or rating"
// @Failure      401  {object}  dto.ErrorResponse "invalid user_id / unauthorized"
// @Failure      403  {object}  dto.ErrorResponse "review not found or not author"
// @Failure      500  {object}  dto.ErrorResponse "failed to update review"
// @Router       /reviews/{id} [put]
// @Security     BearerAuth
func (h *Application) UpdateReview(c *gin.Context) {
	var req dto.UpdateReviewRequest

	reviewID := c.Param("id")

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrInvalidUserID})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	if req.Content == "" && req.Rating == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrAtLeastOneField})
		return
	}

	err = h.ReviewService.UpdateReview(c.Request.Context(), reviewID, userID.String(), req.Content, req.Rating)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidRating):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidRating})

		case errors.Is(err, serviceErrors.ErrReviewNotFound):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrReviewNotFound})

		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrFailedUpdateReview})
		}
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "review updated successfully"})
}
