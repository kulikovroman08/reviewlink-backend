package controller

import (
	"errors"
	"net/http"
	"strconv"

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

	var rating *int
	if r := c.Query("rating"); r != "" {
		if val, err := strconv.Atoi(r); err == nil {
			rating = &val
		} else {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
			return
		}
	}

	sort := c.Query("sort")
	if sort == "" {
		sort = "date_desc"
	}

	reviews, err := h.ReviewService.GetReviews(c.Request.Context(), placeID, rating, sort)
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
