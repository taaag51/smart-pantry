package controllers

import (
	"net/http"
	"strconv"

	"smart-pantry/backend/internal/models"
	"smart-pantry/backend/internal/services"

	"github.com/labstack/echo/v4"
)

// PantryController struct for pantry related controllers
type PantryController struct {
	pantryService *services.PantryService
	recipeService *services.RecipeService
}

// NewPantryController creates a new PantryController
func NewPantryController(ps *services.PantryService, rs *services.RecipeService) *PantryController {
	return &PantryController{
		pantryService: ps,
		recipeService: rs,
	}
}

// GetPantryItems handles the request to get all pantry items.
func (pc *PantryController) GetPantryItems(c echo.Context) error {
	items, err := pc.pantryService.GetPantryItems()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch pantry items"})
	}
	return c.JSON(http.StatusOK, items)
}

// SuggestRecipe suggests a recipe based on pantry items.
// @Router /pantry/recipe [get]
func (pc *PantryController) SuggestRecipe(c echo.Context) error {
	items, err := pc.pantryService.GetPantryItems()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve pantry items"})
	}

	ingredients := []string{}
	for _, item := range items {
		ingredients = append(ingredients, item.Name)
	}

	recipe, err := pc.recipeService.SuggestRecipe(ingredients)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to suggest recipe"})
	}
	return c.JSON(http.StatusOK, map[string]string{"recipe": recipe})
}

// GetPantryItem handles the request to get a pantry item by ID.
func (pc *PantryController) GetPantryItem(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid pantry item ID"})
	}

	item, err := pc.pantryService.GetPantryItem(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch pantry item"})
	}
	if item == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Pantry item not found"})
	}
	return c.JSON(http.StatusOK, item)
}

// AddPantryItem handles the request to add a new pantry item.
func (pc *PantryController) AddPantryItem(c echo.Context) error {
	var pantryItem models.PantryItem
	if err := c.Bind(&pantryItem); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	createdItem, err := pc.pantryService.CreatePantryItem(&pantryItem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create pantry item"})
	}
	return c.JSON(http.StatusCreated, createdItem)
}

// UpdatePantryItem handles the request to update an existing pantry item.
func (pc *PantryController) UpdatePantryItem(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid pantry item ID"})
	}

	var updatedItem models.PantryItem
	if err := c.Bind(&updatedItem); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
	}

	err := pc.pantryService.UpdatePantryItem(uint(id), &updatedItem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update pantry item"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Pantry item updated successfully"})
}

// DeletePantryItem handles the request to delete a pantry item.
func (pc *PantryController) DeletePantryItem(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid pantry item ID"})
	}

	if err := pc.pantryService.DeletePantryItem(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete pantry item"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Pantry item deleted successfully"})
}
