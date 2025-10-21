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
// @Summary      Get top users leaderboard
// @Description  Returns a ranked list of users based on number of reviews and average rating.
// @Tags         leaderboard
// @Accept       json
// @Produce      json
// @Param        limit        query     int     false  "Maximum number of results (default 10, max 100)"
// @Param        sort_by      query     string  false  "Sorting method: 'reviews' or 'rating' (default 'reviews')"
// @Param        min_rating   query     number  false  "Minimum average rating filter (default 0)"
// @Param        min_reviews  query     int     false  "Minimum number of reviews filter (default 0)"
// @Success      200  {array}   dto.LeaderboardEntry  "List of user leaderboard entries"
// @Failure      500  {object}  dto.ErrorResponse      "Internal server error"
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
// @Summary      Get top places leaderboard
// @Description  Returns a ranked list of places based on number of reviews and average rating.
// @Tags         leaderboard
// @Accept       json
// @Produce      json
// @Param        limit        query     int     false  "Maximum number of results (default 10, max 100)"
// @Param        sort_by      query     string  false  "Sorting method: 'reviews' or 'rating' (default 'reviews')"
// @Param        min_rating   query     number  false  "Minimum average rating filter (default 0)"
// @Param        min_reviews  query     int     false  "Minimum number of reviews filter (default 0)"
// @Success      200  {array}   dto.LeaderboardEntry  "List of place leaderboard entries"
// @Failure      500  {object}  dto.ErrorResponse      "Internal server error"
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
