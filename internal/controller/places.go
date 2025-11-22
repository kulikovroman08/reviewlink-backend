package controller

import (
	"errors"
	"net/http"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/response"
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
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})
		return
	}

	var req dto.CreatePlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
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
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: response.ErrPlaceAlreadyExists})

		case errors.Is(err, serviceErrors.ErrInvalidPlaceData):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidPlaceData})

		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrFailedCreatePlace})
		}
		return
	}

	resp := dto.CreatePlaceResponse{
		ID: createdPlace.ID.String(),
	}

	c.JSON(http.StatusCreated, resp)
}

// GetPlaces godoc
// @Summary      Получение списка мест (только для админов)
// @Description  Возвращает список всех заведений. Доступ только для роли **admin**.
// @Tags         admins
// @Produce      json
// @Success 200 {array} dto.PlaceResponse
// @Failure 403 {object} dto.ErrorResponse "access denied"
// @Failure 500 {object} dto.ErrorResponse "failed to load places"
// @Router /places [get]
// @Security     BearerAuth
func (h *Application) GetPlaces(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error: response.ErrAccessDenied,
		})
		return
	}

	places, err := h.PlaceService.GetAllPlaces(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: response.ErrFailedGetPlaces,
		})
		return
	}

	resp := make([]dto.PlaceResponse, 0, len(places))
	for _, p := range places {
		resp = append(resp, dto.PlaceResponse{
			ID:        p.ID.String(),
			Name:      p.Name,
			Address:   p.Address,
			CreatedAt: p.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}
