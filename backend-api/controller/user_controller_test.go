package controller

import (
	"backend-api/controller/response"
	"backend-api/errors"
	"backend-api/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// mockUserUsecase はユーザーユースケースのモック
type mockUserUsecase struct {
	mockSignUp      func(user model.User) (model.UserResponse, error)
	mockLogin       func(user model.User) (string, error)
	mockVerifyToken func(tokenString string) error
}

func (m *mockUserUsecase) SignUp(user model.User) (model.UserResponse, error) {
	return m.mockSignUp(user)
}

func (m *mockUserUsecase) Login(user model.User) (string, error) {
	return m.mockLogin(user)
}

func (m *mockUserUsecase) VerifyToken(tokenString string) error {
	return m.mockVerifyToken(tokenString)
}

// テストヘルパー関数
func setupTest(t *testing.T) (*echo.Echo, *httptest.ResponseRecorder) {
	e := echo.New()
	rec := httptest.NewRecorder()
	return e, rec
}

func assertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedCode int, expectedType string, expectedMessage string) {
	assert.Equal(t, expectedCode, rec.Code)
	var res response.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, expectedType, res.Type)
	assert.Equal(t, expectedMessage, res.Message)
}

func assertSuccessResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedCode int, expectedMessage string) {
	assert.Equal(t, expectedCode, rec.Code)
	var res response.SuccessResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, res.Message)
}

func TestUserController_SignUp(t *testing.T) {
	e, rec := setupTest(t)

	tests := []struct {
		name            string
		inputUser       model.User
		mockBehavior    func(*mockUserUsecase)
		expectedCode    int
		expectedType    string
		expectedMessage string
	}{
		{
			name: "正常系：ユーザー登録成功",
			inputUser: model.User{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockBehavior: func(m *mockUserUsecase) {
				m.mockSignUp = func(user model.User) (model.UserResponse, error) {
					return model.UserResponse{
						ID:    1,
						Email: user.Email,
					}, nil
				}
			},
			expectedCode:    http.StatusCreated,
			expectedMessage: "ユーザーが正常に作成されました",
		},
		{
			name: "異常系：既存のメールアドレス",
			inputUser: model.User{
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockBehavior: func(m *mockUserUsecase) {
				m.mockSignUp = func(user model.User) (model.UserResponse, error) {
					return model.UserResponse{}, errors.EmailExists
				}
			},
			expectedCode:    http.StatusBadRequest,
			expectedType:    string(errors.BusinessError),
			expectedMessage: errors.EmailExists.Message,
		},
		{
			name: "異常系：無効なメールアドレス",
			inputUser: model.User{
				Email:    "invalid-email",
				Password: "password123",
			},
			mockBehavior: func(m *mockUserUsecase) {
				m.mockSignUp = func(user model.User) (model.UserResponse, error) {
					return model.UserResponse{}, errors.InvalidEmail
				}
			},
			expectedCode:    http.StatusBadRequest,
			expectedType:    string(errors.ValidationError),
			expectedMessage: errors.InvalidEmail.Message,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mock := &mockUserUsecase{}
			tt.mockBehavior(mock)

			// コントローラーの作成
			uc := NewUserController(mock)

			// リクエストの準備
			jsonBody, _ := json.Marshal(tt.inputUser)
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)

			// テスト実行
			err := uc.SignUp(c)
			assert.NoError(t, err)

			// レスポンスの検証
			if tt.expectedType != "" {
				assertErrorResponse(t, rec, tt.expectedCode, tt.expectedType, tt.expectedMessage)
			} else {
				assertSuccessResponse(t, rec, tt.expectedCode, tt.expectedMessage)
			}

			// レコーダーをリセット
			rec.Body.Reset()
		})
	}
}

func TestUserController_Login(t *testing.T) {
	e, rec := setupTest(t)

	tests := []struct {
		name            string
		inputUser       model.User
		mockBehavior    func(*mockUserUsecase)
		expectedCode    int
		expectedType    string
		expectedMessage string
		checkCookie     bool
	}{
		{
			name: "正常系：ログイン成功",
			inputUser: model.User{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockBehavior: func(m *mockUserUsecase) {
				m.mockLogin = func(user model.User) (string, error) {
					return "test-jwt-token", nil
				}
			},
			expectedCode:    http.StatusOK,
			expectedMessage: "ログインに成功しました",
			checkCookie:     true,
		},
		{
			name: "異常系：認証失敗",
			inputUser: model.User{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockBehavior: func(m *mockUserUsecase) {
				m.mockLogin = func(user model.User) (string, error) {
					return "", errors.InvalidCredentials
				}
			},
			expectedCode:    http.StatusUnauthorized,
			expectedType:    string(errors.AuthenticationError),
			expectedMessage: errors.InvalidCredentials.Message,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mock := &mockUserUsecase{}
			tt.mockBehavior(mock)

			// コントローラーの作成
			uc := NewUserController(mock)

			// リクエストの準備
			jsonBody, _ := json.Marshal(tt.inputUser)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)

			// テスト実行
			err := uc.LogIn(c)
			assert.NoError(t, err)

			// レスポンスの検証
			if tt.expectedType != "" {
				assertErrorResponse(t, rec, tt.expectedCode, tt.expectedType, tt.expectedMessage)
			} else {
				assertSuccessResponse(t, rec, tt.expectedCode, tt.expectedMessage)
			}

			// クッキーの検証
			if tt.checkCookie {
				cookies := rec.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "token" {
						found = true
						assert.NotEmpty(t, cookie.Value)
						assert.True(t, cookie.Secure)
						assert.True(t, cookie.HttpOnly)
						assert.Equal(t, http.SameSiteNoneMode, cookie.SameSite)
						assert.True(t, cookie.Expires.After(time.Now()))
					}
				}
				assert.True(t, found, "認証クッキーが見つかりません")
			}

			// レコーダーをリセット
			rec.Body.Reset()
		})
	}
}

func TestUserController_Logout(t *testing.T) {
	e, rec := setupTest(t)

	t.Run("正常系：ログアウト成功", func(t *testing.T) {
		// コントローラーの作成
		mock := &mockUserUsecase{}
		uc := NewUserController(mock)

		// リクエストの準備
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		c := e.NewContext(req, rec)

		// テスト実行
		err := uc.LogOut(c)
		assert.NoError(t, err)

		// レスポンスの検証
		assertSuccessResponse(t, rec, http.StatusOK, "ログアウトしました")

		// クッキーの検証
		cookies := rec.Result().Cookies()
		found := false
		for _, cookie := range cookies {
			if cookie.Name == "token" {
				found = true
				assert.Empty(t, cookie.Value)
				assert.True(t, cookie.Expires.Before(time.Now()))
			}
		}
		assert.True(t, found, "クリアされた認証クッキーが見つかりません")
	})
}

func TestUserController_CsrfToken(t *testing.T) {
	e, rec := setupTest(t)

	t.Run("正常系：CSRFトークン取得", func(t *testing.T) {
		// コントローラーの作成
		mock := &mockUserUsecase{}
		uc := NewUserController(mock)

		// リクエストの準備
		req := httptest.NewRequest(http.MethodGet, "/csrf", nil)
		c := e.NewContext(req, rec)
		c.Set("csrf", "test-csrf-token")

		// テスト実行
		err := uc.CsrfToken(c)
		assert.NoError(t, err)

		// レスポンスの検証
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "test-csrf-token", data["csrf_token"])
	})
}

func TestUserController_VerifyToken(t *testing.T) {
	e, rec := setupTest(t)

	t.Run("正常系：トークン検証成功", func(t *testing.T) {
		// コントローラーの作成
		mock := &mockUserUsecase{}
		uc := NewUserController(mock)

		// リクエストの準備
		req := httptest.NewRequest(http.MethodGet, "/verify", nil)
		c := e.NewContext(req, rec)

		// テスト実行
		err := uc.VerifyToken(c)
		assert.NoError(t, err)

		// レスポンスの検証
		assertSuccessResponse(t, rec, http.StatusOK, "トークンは有効です")
	})
}
