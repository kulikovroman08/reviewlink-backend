package controller

import (
	"errors"
	"net/http"

	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/response"
)

// GenerateTokens godoc
// @Summary      Генерация токенов (только для админов)
// @Description  Эндпоинт доступен только пользователям с ролью **admin**.
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        request  body      dto.GenerateTokensRequest  true  "Данные для генерации токенов"
// @Success      200      {object}  dto.GenerateTokensResponse
// @Failure      400      {object}  dto.ErrorResponse "invalid input"
// @Failure      403      {object}  dto.ErrorResponse "only admin can generate tokens"
// @Failure      500      {object}  dto.ErrorResponse "failed to generate tokens"
// @Router       /admin/tokens [post]
// @Security     BearerAuth
func (a *Application) GenerateTokens(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrOnlyAdminCanGenerateTokens})
		return
	}

	var req dto.GenerateTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	resp, err := a.TokenService.GenerateTokens(c.Request.Context(), req.PlaceID, req.Count)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidPlaceID):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidPlaceID})

		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrFailedGenerateTokens})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
