package controller

import (
	"encoding/json"
	"errors"
	"go-rest-api/mock"
	"go-rest-api/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestFoodItemController_GetAllFoodItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFoodItemUsecase := mock.NewMockIFoodItemUsecase(ctrl)
	foodItemController := NewFoodItemController(mockFoodItemUsecase)

	tests := []struct {
		name           string
		buildStubs     func()
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
		expectedStatus int
	}{
		{
			name: "正常系：食材一覧の取得成功",
			buildStubs: func() {
				mockFoodItemUsecase.EXPECT().
					GetAllFoodItems().
					Times(1).
					Return([]model.FoodItem{
						{
							ID:         1,
							Title:      "りんご",
							Quantity:   5,
							ExpiryDate: time.Now().Add(24 * time.Hour),
							UserId:     1,
							CreatedAt:  time.Now(),
							UpdatedAt:  time.Now(),
						},
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)

				foodItems, ok := response.Data.([]interface{})
				assert.True(t, ok)
				assert.Len(t, foodItems, 1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "異常系：内部エラー",
			buildStubs: func() {
				mockFoodItemUsecase.EXPECT().
					GetAllFoodItems().
					Times(1).
					Return(nil, errors.New("internal error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "internal error", response.Message)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/food-items", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.buildStubs()

			err := foodItemController.GetAllFoodItems(c)
			assert.NoError(t, err)

			tt.checkResponse(t, rec)
		})
	}
}

func TestFoodItemController_GetFoodItemById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFoodItemUsecase := mock.NewMockIFoodItemUsecase(ctrl)
	foodItemController := NewFoodItemController(mockFoodItemUsecase)

	tests := []struct {
		name           string
		id             string
		buildStubs     func()
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
		expectedStatus int
	}{
		{
			name: "正常系：食材の取得成功",
			id:   "1",
			buildStubs: func() {
				mockFoodItemUsecase.EXPECT().
					GetFoodItemById(uint(1)).
					Times(1).
					Return(&model.FoodItem{
						ID:         1,
						Title:      "りんご",
						Quantity:   5,
						ExpiryDate: time.Now().Add(24 * time.Hour),
						UserId:     1,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)

				foodItem, ok := response.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, float64(1), foodItem["id"])
				assert.Equal(t, "りんご", foodItem["title"])
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "異常系：不正なID",
			id:   "invalid",
			buildStubs: func() {
				// 不正なIDの場合はusecaseは呼ばれない
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid ID format", response.Message)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/food-items/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			tt.buildStubs()

			err := foodItemController.GetFoodItemById(c)
			assert.NoError(t, err)

			tt.checkResponse(t, rec)
		})
	}
}

func TestFoodItemController_CreateFoodItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFoodItemUsecase := mock.NewMockIFoodItemUsecase(ctrl)
	foodItemController := NewFoodItemController(mockFoodItemUsecase)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(1),
	})

	tests := []struct {
		name           string
		body           string
		buildStubs     func()
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
		expectedStatus int
	}{
		{
			name: "正常系：食材の作成成功",
			body: `{"title":"りんご","quantity":5,"expiry_date":"2024-02-01T00:00:00Z"}`,
			buildStubs: func() {
				mockFoodItemUsecase.EXPECT().
					CreateFoodItem(gomock.Any()).
					Times(1).
					Return(&model.FoodItem{
						ID:         1,
						Title:      "りんご",
						Quantity:   5,
						ExpiryDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
						UserId:     1,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Food item created successfully", response.Message)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "異常系：不正なリクエストボディ",
			body: `{"title": 123}`, // 不正な型
			buildStubs: func() {
				// 不正なボディの場合はusecaseは呼ばれない
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid request format", response.Message)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/food-items", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user", token)

			tt.buildStubs()

			err := foodItemController.CreateFoodItem(c)
			assert.NoError(t, err)

			tt.checkResponse(t, rec)
		})
	}
}

func TestFoodItemController_UpdateFoodItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFoodItemUsecase := mock.NewMockIFoodItemUsecase(ctrl)
	foodItemController := NewFoodItemController(mockFoodItemUsecase)

	tests := []struct {
		name           string
		id             string
		body           string
		buildStubs     func()
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
		expectedStatus int
	}{
		{
			name: "正常系：食材の更新成功",
			id:   "1",
			body: `{"title":"更新済みりんご","quantity":3,"expiry_date":"2024-02-01T00:00:00Z"}`,
			buildStubs: func() {
				mockFoodItemUsecase.EXPECT().
					UpdateFoodItem(gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Food item updated successfully", response.Message)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "異常系：不正なID",
			id:   "invalid",
			body: `{"title":"りんご","quantity":5}`,
			buildStubs: func() {
				// 不正なIDの場合はusecaseは呼ばれない
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid ID format", response.Message)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/food-items/:id", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			tt.buildStubs()

			err := foodItemController.UpdateFoodItem(c)
			assert.NoError(t, err)

			tt.checkResponse(t, rec)
		})
	}
}

func TestFoodItemController_DeleteFoodItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFoodItemUsecase := mock.NewMockIFoodItemUsecase(ctrl)
	foodItemController := NewFoodItemController(mockFoodItemUsecase)

	tests := []struct {
		name           string
		id             string
		buildStubs     func()
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
		expectedStatus int
	}{
		{
			name: "正常系：食材の削除成功",
			id:   "1",
			buildStubs: func() {
				mockFoodItemUsecase.EXPECT().
					DeleteFoodItem(uint(1)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Food item deleted successfully", response.Message)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "異常系：不正なID",
			id:   "invalid",
			buildStubs: func() {
				// 不正なIDの場合はusecaseは呼ばれない
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
				var response Response
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid ID format", response.Message)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/api/food-items/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			tt.buildStubs()

			err := foodItemController.DeleteFoodItem(c)
			assert.NoError(t, err)

			tt.checkResponse(t, rec)
		})
	}
}
