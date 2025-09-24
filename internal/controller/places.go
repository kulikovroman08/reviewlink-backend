package controller

import (
	"errors"
	"net/http"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
)

// CreatePlace godoc
// @Summary      Создание места (только для админов)
// @Description  Эндпоинт доступен только пользователям с ролью **admin**.
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreatePlaceRequest  true  "Данные для создания места"
// @Success      201      {object}  dto.CreatePlaceResponse
// @Failure 400 {object} dto.ErrorResponse "invalid input / invalid place data"
// @Failure 403 {object} dto.ErrorResponse "access denied"
// @Failure 409 {object} dto.ErrorResponse "place already exists"
// @Failure 500 {object} dto.ErrorResponse "failed to create place"
// @Router       /places [post]
// @Security     BearerAuth
func (h *Application) CreatePlace(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: dto.ErrAccessDenied})
		return
	}

	var req dto.CreatePlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrInvalidInput})
		return
	}

	place := model.Place{
		Name:    req.Name,
		Address: req.Address,
	}

	createdPlace, err := h.PlaceService.CreatePlace(c.Request.Context(), place)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrPlaceAlreadyExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: dto.ErrPlaceAlreadyExists})

		case errors.Is(err, serviceErrors.ErrInvalidPlaceData):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrInvalidPlaceData})

		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: dto.ErrFailedCreatePlace})
		}
		return
	}

	resp := dto.CreatePlaceResponse{
		ID: createdPlace.ID.String(),
	}

	c.JSON(http.StatusCreated, resp)
}
