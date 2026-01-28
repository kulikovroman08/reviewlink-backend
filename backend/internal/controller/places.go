package controller

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/response"
)

// CreatePlace godoc
// @Summary      Создание места (только для админов)
// @Description  Создаёт заведение и привязывает его к текущему администратору
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

	ownerID := c.GetString("user_id")
	if ownerID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrUnauthorized})
		return
	}

	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrUnauthorized})
		return
	}

	var req dto.CreatePlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	place := model.Place{
		OwnerID: &ownerUUID,
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
// @Description  Возвращает список заведений, принадлежащих текущему администратору
// @Tags         admins
// @Produce      json
// @Success 200 {array} dto.PlaceResponse
// @Failure      401 {object} dto.ErrorResponse "unauthorized"
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

	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrUnauthorized})
		return
	}

	if _, err := uuid.Parse(adminID); err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrUnauthorized})
		return
	}

	places, err := h.PlaceService.GetPlacesByOwner(c.Request.Context(), adminID)
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

// ListPublicPlaces godoc
// @Summary      Публичный каталог заведений
// @Description  Возвращает список заведений из OSM (Overpass) и подмешивает мету из ReviewLink (rating, reviewsCount, hasOwner)
// @Tags         places
// @Produce      json
// @Param        city    query     string  true   "Город (обязателен)"
// @Param        search  query     string  false  "Поиск по имени (regex)"
// @Param        amenity query     string  false  "Фильтр по amenity (например cafe|restaurant)"
// @Param        limit   query     int     false  "Лимит (по умолчанию 50, максимум 200)"
// @Param        offset  query     int     false  "Смещение (MVP может игнорироваться)"
// @Success      200 {object} dto.PublicPlacesResponse
// @Failure      400 {object} dto.ErrorResponse "invalid input"
// @Failure      500 {object} dto.ErrorResponse "failed to load public places"
// @Router       /places/public [get]
func (h *Application) ListPublicPlaces(c *gin.Context) {
	var q dto.PublicPlacesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	// city обязателен
	if q.City == nil || *q.City == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	limit := 50
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}
	if limit > 200 {
		limit = 200
	}

	offset := 0
	if q.Offset != nil && *q.Offset > 0 {
		offset = *q.Offset
	}

	params := model.PublicPlacesParams{
		City:    *q.City,
		Search:  q.Search,
		Amenity: q.Amenity,
		Limit:   limit,
		Offset:  offset,
	}

	res, err := h.PlaceService.ListPublicPlaces(c.Request.Context(), params)
	if err != nil {
		log.Printf("ListPublicPlaces error: %+v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrFailedGetPlaces})
		return
	}

	items := make([]dto.PublicPlaceResponse, 0, len(res.Items))
	for _, it := range res.Items {
		dtoItem := dto.PublicPlaceResponse{
			Source:   it.Source,
			SourceID: it.SourceID,
			Name:     it.Name,
			Amenity:  it.Amenity,
			Address: dto.PublicPlaceAddress{
				City:        it.Address.City,
				Street:      it.Address.Street,
				HouseNumber: it.Address.HouseNumber,
				Postcode:    it.Address.Postcode,
				Country:     it.Address.Country,
				Display:     it.Address.Display,
			},
			Location: dto.PublicLatLon{
				Lat: it.Location.Lat,
				Lon: it.Location.Lon,
			},
		}

		if it.Reviewlink != nil {
			dtoItem.Reviewlink = &dto.PublicPlaceReviewlinkMeta{
				PlaceID:      it.Reviewlink.PlaceID.String(),
				Rating:       it.Reviewlink.Rating,
				ReviewsCount: it.Reviewlink.ReviewsCount,
				HasOwner:     it.Reviewlink.HasOwner,
			}
		}

		items = append(items, dtoItem)
	}

	c.JSON(http.StatusOK, dto.PublicPlacesResponse{
		Items: items,
		Total: res.Total,
	})
}

// EnsurePlaceFromPublic godoc
// @Summary      Создать/получить place_id по публичному источнику (OSM)
// @Description  По (source, source_id) находит place в ReviewLink или создаёт новый (owner_id = NULL) и возвращает place_id
// @Tags         places
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.EnsurePlaceFromPublicRequest  true  "source/source_id/name"
// @Success      200   {object}  dto.EnsurePlaceFromPublicResponse
// @Failure      400   {object}  dto.ErrorResponse "invalid input"
// @Failure      500   {object}  dto.ErrorResponse "failed to ensure place"
// @Router       /places/ensure_from_public [post]
func (h *Application) EnsurePlaceFromPublic(c *gin.Context) {
	var req dto.EnsurePlaceFromPublicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	source := strings.TrimSpace(req.Source)

	sourceID := strings.TrimSpace(req.SourceID)
	if sourceID == "" {
		sourceID = strings.TrimSpace(req.SourceId)
	}

	name := strings.TrimSpace(req.Name)

	if source == "" || sourceID == "" || name == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	placeID, err := h.PlaceService.EnsureFromPublic(c.Request.Context(), source, sourceID, name)
	if err != nil {
		log.Printf("EnsurePlaceFromPublic error: %+v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to ensure place"})
		return
	}

	c.JSON(http.StatusOK, dto.EnsurePlaceFromPublicResponse{PlaceID: placeID})
}
