package controller

import (
	"backend-api/mock"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRecipeController_GetRecipeSuggestions(t *testing.T) {
	// コントローラのセットアップ
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRecipeUsecase := mock.NewMockIRecipeUsecase(ctrl)
	recipeController := NewRecipeController(mockRecipeUsecase)

	// JWTトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(1),
	})

	tests := []struct {
		name           string
		setupAuth      func(c echo.Context)
		buildStubs     func()
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
		expectedStatus int
	}{
		{
			name: "正常系：レシピ提案の取得成功",
			setupAuth: func(c echo.Context) {
				c.Set("user", token)
			},
			buildStubs: func() {
				mockRecipeUsecase.EXPECT().
					GetRecipeSuggestions(uint(1)).
					Times(1).
					Return("おすすめレシピ：トマトパスタ\n材料：トマト、パスタ、オリーブオイル\n手順：...", nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				var response []string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response, 1)
				assert.Contains(t, response[0], "おすすめレシピ")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "異常系：認証トークンなし",
			setupAuth: func(c echo.Context) {
				// トークンを設定しない
			},
			buildStubs: func() {
				// 認証エラーの場合はusecaseは呼ばれない
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "異常系：レシピ生成エラー",
			setupAuth: func(c echo.Context) {
				c.Set("user", token)
			},
			buildStubs: func() {
				mockRecipeUsecase.EXPECT().
					GetRecipeSuggestions(uint(1)).
					Times(1).
					Return("", errors.New("レシピ生成エラー"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
				var response map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "レシピの提案に失敗しました", response["error"])
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "正常系：食材なしの場合のメッセージ",
			setupAuth: func(c echo.Context) {
				c.Set("user", token)
			},
			buildStubs: func() {
				mockRecipeUsecase.EXPECT().
					GetRecipeSuggestions(uint(1)).
					Times(1).
					Return("食材が登録されていません。食材を追加してからレシピを取得してください。", nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				var response []string
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response, 1)
				assert.Contains(t, response[0], "食材が登録されていません")
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echoのセットアップ
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/recipes/suggestions", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 認証設定
			tt.setupAuth(c)

			// スタブの構築
			tt.buildStubs()

			// テスト実行
			err := recipeController.GetRecipeSuggestions(c)
			if err != nil {
				// エラーハンドリングの検証
				he, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			}

			// レスポンスの検証
			tt.checkResponse(t, rec)
		})
	}
}
