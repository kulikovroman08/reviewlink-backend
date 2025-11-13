package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/response"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

// GetUserLeaderboard godoc
// @Summary      Получить топ пользователей
// @Description  Возвращает список пользователей, отсортированных по количеству отзывов и среднему рейтингу.
// @Tags         leaderboard
// @Accept       json
// @Produce      json
// @Param        limit        query     int     false  "Максимальное количество результатов (по умолчанию 10, максимум 100)"
// @Param        sort_by      query     string  false  "Сортировка: 'reviews' или 'rating' (по умолчанию 'reviews')"
// @Param        min_rating   query     number  false  "Минимальный средний рейтинг (по умолчанию 0)"
// @Param        min_reviews  query     int     false  "Минимальное количество отзывов (по умолчанию 0)"
// @Success      200  {array}   dto.LeaderboardEntry  "Список пользователей"
// @Failure      500  {object}  dto.ErrorResponse      "Внутренняя ошибка сервера"
// @Router       /leaderboard/users [get]
func (a *Application) GetUserLeaderboard(c *gin.Context) {
	limitStr := parseQueryParam(c, "limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	} else if limit <= 0 || limit > 100 {
		limit = 10
	}

	minRatingStr := parseQueryParam(c, "min_rating", "0")
	minRating, err := strconv.ParseFloat(minRatingStr, 64)
	if err != nil {
		minRating = 0
	} else if minRating < 0 {
		minRating = 0
	}

	minReviewsStr := parseQueryParam(c, "min_reviews", "0")
	minReviews, err := strconv.Atoi(minReviewsStr)
	if err != nil {
		minReviews = 0
	} else if minReviews < 0 {
		minReviews = 0
	}

	sortBy := parseQueryParam(c, "sort_by", "reviews")
	if sortBy != "reviews" && sortBy != "rating" {
		sortBy = "reviews"
	}

	filter := model.LeaderboardFilter{
		SortBy:     sortBy,
		MinRating:  minRating,
		MinReviews: minReviews,
	}

	entries, err := a.LeaderboardService.GetUserLeaderboard(c.Request.Context(), limit, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})
		return
	}

	result := make([]dto.LeaderboardEntry, len(entries))
	for i, entry := range entries {
		result[i] = dto.LeaderboardEntry{
			Rank:         i + 1,
			ID:           entry.ID,
			Name:         entry.Name,
			ReviewsCount: entry.ReviewsCount,
			AvgRating:    entry.AvgRating,
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetPlaceLeaderboard godoc
// @Summary      Получить топ заведений
// @Description  Возвращает список заведений, отсортированных по количеству отзывов и среднему рейтингу.
// @Tags         leaderboard
// @Accept       json
// @Produce      json
// @Param        limit        query     int     false  "Максимальное количество результатов (по умолчанию 10, максимум 100)"
// @Param        sort_by      query     string  false  "Сортировка: 'reviews' или 'rating' (по умолчанию 'reviews')"
// @Param        min_rating   query     number  false  "Минимальный средний рейтинг (по умолчанию 0)"
// @Param        min_reviews  query     int     false  "Минимальное количество отзывов (по умолчанию 0)"
// @Success      200  {array}   dto.LeaderboardEntry  "Список заведений"
// @Failure      500  {object}  dto.ErrorResponse      "Внутренняя ошибка сервера"
// @Router       /leaderboard/places [get]
func (a *Application) GetPlaceLeaderboard(c *gin.Context) {
	limitStr := parseQueryParam(c, "limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	} else if limit <= 0 || limit > 100 {
		limit = 10
	}

	minRatingStr := parseQueryParam(c, "min_rating", "0")
	minRating, err := strconv.ParseFloat(minRatingStr, 64)
	if err != nil {
		minRating = 0
	} else if minRating < 0 {
		minRating = 0
	}

	minReviewsStr := parseQueryParam(c, "min_reviews", "0")
	minReviews, err := strconv.Atoi(minReviewsStr)
	if err != nil {
		minReviews = 0
	} else if minReviews < 0 {
		minReviews = 0
	}

	sortBy := parseQueryParam(c, "sort_by", "reviews")
	if sortBy != "reviews" && sortBy != "rating" {
		sortBy = "reviews"
	}

	filter := model.LeaderboardFilter{
		SortBy:     sortBy,
		MinRating:  minRating,
		MinReviews: minReviews,
	}

	entries, err := a.LeaderboardService.GetPlaceLeaderboard(c.Request.Context(), limit, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})
		return
	}

	result := make([]dto.LeaderboardEntry, len(entries))
	for i, entry := range entries {
		result[i] = dto.LeaderboardEntry{
			Rank:         i + 1,
			ID:           entry.ID,
			Name:         entry.Name,
			ReviewsCount: entry.ReviewsCount,
			AvgRating:    entry.AvgRating,
		}
	}

	c.JSON(http.StatusOK, result)
}

func parseQueryParam(c *gin.Context, key string, defaultValue string) string {
	if value := c.Query(key); value != "" {
		return value
	}
	return defaultValue
}
