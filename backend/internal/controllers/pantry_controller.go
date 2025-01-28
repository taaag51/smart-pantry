package controllers

import (
	"net/http"
	"strconv"

	"smart-pantry/backend/internal/services"

	"github.com/labstack/echo/v4"
)

// PantryController handles pantry-related API endpoints
type PantryController struct {
	PantryService services.PantryService
}

// NewPantryController creates a new PantryController
func NewPantryController(ps services.PantryService) *PantryController {
	return &PantryController{
		PantryService: ps,
	}
}

// GetPantryItems retrieves all pantry items
func (pc *PantryController) GetPantryItems(c echo.Context) error {
	items, err := pc.PantryService.GetPantryItems()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

// CreatePantryItem adds a new pantry item
func (pc *PantryController) CreatePantryItem(c echo.Context) error {
	var params services.PantryItemParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	item, err := pc.PantryService.CreatePantryItem(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, item)
}

// UpdatePantryItem updates an existing pantry item
func (pc *PantryController) UpdatePantryItem(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}

	var params services.PantryItemParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	item, err := pc.PantryService.UpdatePantryItem(id, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, item)
}

// DeletePantryItem removes a pantry item
func (pc *PantryController) DeletePantryItem(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}

	if err := pc.PantryService.DeletePantryItem(id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}
