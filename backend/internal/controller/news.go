package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
)

// GetNews godoc
// @Summary      Новости HoReCa
// @Description  Публичная лента новостей HoReCa (RSS РБК, с кешированием)
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Количество новостей (по умолчанию 3, максимум 20)"
// @Success      200     {object}  dto.NewsResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Router       /news [get]
func (h *Application) GetNews(c *gin.Context) {
	var q dto.NewsQuery
	_ = c.ShouldBindQuery(&q)

	limit := 0
	if q.Limit != nil {
		limit = *q.Limit
	}

	items, cachedAt, err := h.publicNewsService.ListNews(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusOK, dto.NewsResponse{
			Items: []dto.NewsItemResponse{},
		})
		return
	}

	resp := dto.NewsResponse{
		Items:    make([]dto.NewsItemResponse, 0, len(items)),
		CachedAt: cachedAt,
	}

	for _, it := range items {
		resp.Items = append(resp.Items, dto.NewsItemResponse{
			Title:       it.Title,
			Link:        it.Link,
			PublishedAt: it.PublishedAt,
			Source:      it.Source,
		})
	}

	c.JSON(http.StatusOK, resp)
}
