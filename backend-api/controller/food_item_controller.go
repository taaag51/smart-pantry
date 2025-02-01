package controller

import (
	"backend-api/model"
	"backend-api/usecase"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

/**
 * 食材コントローラーのインターフェース
 */
type IFoodItemController interface {
	GetAllFoodItems(c echo.Context) error
	GetFoodItemById(c echo.Context) error
	CreateFoodItem(c echo.Context) error
	UpdateFoodItem(c echo.Context) error
	DeleteFoodItem(c echo.Context) error
}

/**
 * 食材コントローラーの構造体
 */
type foodItemController struct {
	fu usecase.IFoodItemUsecase
}

/**
 * APIレスポンスの構造体
 */
type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

/**
 * 食材コントローラーのコンストラクタ
 * @param fu 食材ユースケースのインターフェース
 * @return 食材コントローラーのインターフェース
 */
func NewFoodItemController(fu usecase.IFoodItemUsecase) IFoodItemController {
	return &foodItemController{fu}
}

/**
 * エラーハンドリング用のヘルパー関数
 */
func handleError(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, Response{
		Message: message,
	})
}

/**
 * 全ての食材を取得
 * @param c コンテキスト
 * @return エラー
 */
func (fc *foodItemController) GetAllFoodItems(c echo.Context) error {
	foodItems, err := fc.fu.GetAllFoodItems()
	if err != nil {
		return handleError(c, http.StatusInternalServerError, "食材の取得に失敗しました")
	}
	return c.JSON(http.StatusOK, Response{
		Data: foodItems,
	})
}

/**
 * IDによる食材の取得
 * @param c コンテキスト
 * @return エラー
 */
func (fc *foodItemController) GetFoodItemById(c echo.Context) error {
	id := c.Param("id")
	foodItemId, err := strconv.Atoi(id)
	if err != nil {
		return handleError(c, http.StatusBadRequest, "無効なID形式です")
	}

	foodItem, err := fc.fu.GetFoodItemById(uint(foodItemId))
	if err != nil {
		return handleError(c, http.StatusInternalServerError, "食材の取得に失敗しました")
	}
	return c.JSON(http.StatusOK, Response{
		Data: foodItem,
	})
}

/**
 * 食材の作成
 * @param c コンテキスト
 * @return エラー
 */
func (fc *foodItemController) CreateFoodItem(c echo.Context) error {
	foodItem := model.FoodItem{}
	if err := c.Bind(&foodItem); err != nil {
		return handleError(c, http.StatusBadRequest, "リクエスト形式が無効です")
	}

	// JWTトークンからユーザーIDを取得
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.MapClaims)
	userId := uint((*claims)["user_id"].(float64))
	foodItem.UserId = userId

	createdFoodItem, err := fc.fu.CreateFoodItem(foodItem)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, "食材の作成に失敗しました")
	}
	return c.JSON(http.StatusCreated, Response{
		Data:    createdFoodItem,
		Message: "Food item created successfully",
	})
}

/**
 * 食材の更新
 * @param c コンテキスト
 * @return エラー
 */
func (fc *foodItemController) UpdateFoodItem(c echo.Context) error {
	foodItem := model.FoodItem{}
	if err := c.Bind(&foodItem); err != nil {
		return handleError(c, http.StatusBadRequest, "リクエスト形式が無効です")
	}

	// IDの存在確認
	id := c.Param("id")
	foodItemId, err := strconv.Atoi(id)
	if err != nil {
		return handleError(c, http.StatusBadRequest, "無効なID形式です")
	}
	foodItem.ID = uint(foodItemId)

	if err := fc.fu.UpdateFoodItem(foodItem); err != nil {
		return handleError(c, http.StatusInternalServerError, "食材の更新に失敗しました")
	}
	return c.JSON(http.StatusOK, Response{
		Data:    foodItem,
		Message: "Food item updated successfully",
	})
}

/**
 * 食材の削除
 * @param c コンテキスト
 * @return エラー
 */
func (fc *foodItemController) DeleteFoodItem(c echo.Context) error {
	id := c.Param("id")
	foodItemId, err := strconv.Atoi(id)
	if err != nil {
		return handleError(c, http.StatusBadRequest, "無効なID形式です")
	}

	if err := fc.fu.DeleteFoodItem(uint(foodItemId)); err != nil {
		return handleError(c, http.StatusInternalServerError, "食材の削除に失敗しました")
	}
	return c.JSON(http.StatusOK, Response{
		Message: "Food item deleted successfully",
	})
}
