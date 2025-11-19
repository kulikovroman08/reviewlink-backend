package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/response"
)

// GetStats godoc
// @Summary      Получение общей статистики (только для админов)
// @Description  Возвращает агрегированные данные: количество пользователей, отзывов, средний рейтинг и количество бонусов.
// @Tags         admins
// @Produce      json
// @Success      200  {object}  dto.AdminStatsResponse
// @Failure      403  {object}  dto.ErrorResponse "access denied"
// @Failure      500  {object}  dto.ErrorResponse "failed to load stats"
// @Router       /admin/stats [get]
// @Security     BearerAuth
func (h *Application) GetStats(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})
		return
	}

	stats, err := h.AdminService.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrFailedLoadStats})
		return
	}

	resp := dto.AdminStatsResponse{
		TotalUsers:    stats.TotalUsers,
		TotalReviews:  stats.TotalReviews,
		AverageRating: stats.AverageRating,
		TotalBonuses:  stats.TotalBonuses,
	}

	c.JSON(http.StatusOK, resp)
}
