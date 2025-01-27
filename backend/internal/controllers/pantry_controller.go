package controllers

import (
	"net/http"
	"strconv"

	"smart-pantry/backend/internal/services"

	"github.com/labstack/echo/v4"
)

type PantryController struct {
	service *services.PantryService
}

func NewPantryController(service *services.PantryService) *PantryController {
	return &PantryController{service: service}
}

func (pc *PantryController) GetPantryItems(c echo.Context) error {
	items, err := pc.service.GetPantryItems()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get pantry items"})
	}
	return c.JSON(http.StatusOK, items)
}

func (pc *PantryController) CreatePantryItem(c echo.Context) error {
	var item services.PantryItemParams
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	createdItem, err := pc.service.CreatePantryItem(item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create pantry item"})
	}
	return c.JSON(http.StatusCreated, createdItem)
}

func (pc *PantryController) UpdatePantryItem(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid item ID"})
	}

	var item services.PantryItemParams
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	updatedItem, err := pc.service.UpdatePantryItem(id, item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update pantry item"})
	}
	return c.JSON(http.StatusOK, updatedItem)
}

func (pc *PantryController) DeletePantryItem(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid item ID"})
	}

	if err := pc.service.DeletePantryItem(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete pantry item"})
	}
	return c.NoContent(http.StatusNoContent)
}
