package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/taaag51/smart-pantry/backend-api/controller/response"
	"github.com/taaag51/smart-pantry/backend-api/errors"
	"github.com/taaag51/smart-pantry/backend-api/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// mockUserUsecase はユーザーユースケースのモック
type mockUserUsecase struct {
	mockSignUp      func(user model.User) (model.UserResponse, error)
	mockLogin       func(user model.User) (string, error)
	mockVerifyToken func(tokenString string) (*jwt.Token, error)
}

func (m *mockUserUsecase) SignUp(user model.User) (model.UserResponse, error) {
	return m.mockSignUp(user)
}

func (m *mockUserUsecase) Login(user model.User) (string, error) {
	return m.mockLogin(user)
}

func (m *mockUserUsecase) VerifyToken(tokenString string) (*jwt.Token, error) {
	return m.mockVerifyToken(tokenString)
}

// テストヘルパー関数
func setupTest(t *testing.T) (*echo.Echo, *httptest.ResponseRecorder) {
	e := echo.New()
	rec := httptest.NewRecorder()
	return e, rec
}

func assertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedCode int, expectedMessage string) {
	assert.Equal(t, expectedCode, rec.Code)
	var res response.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)
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
					return model.UserResponse{}, errors.New(errors.BusinessError, "メールアドレスが既に存在します", http.StatusBadRequest, nil)
				}
			},
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "ユーザーの登録に失敗しました",
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
			assertSuccessResponse(t, rec, tt.expectedCode, tt.expectedMessage)

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
					return "", errors.New(errors.AuthenticationError, "認証に失敗しました", http.StatusUnauthorized, nil)
				}
			},
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: "メールアドレスまたはパスワードが正しくありません",
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
			assertSuccessResponse(t, rec, tt.expectedCode, tt.expectedMessage)

			// クッキーの検証
			if tt.checkCookie {
				cookies := rec.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == "token" {
						found = true
						assert.NotEmpty(t, cookie.Value)
						assert.True(t, cookie.HttpOnly)
						assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
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

func TestUserController_VerifyToken(t *testing.T) {
	e, rec := setupTest(t)

	tests := []struct {
		name            string
		token           string
		mockBehavior    func(*mockUserUsecase)
		expectedCode    int
		expectedMessage string
	}{
		{
			name:  "正常系：トークン検証成功",
			token: "Bearer valid-token",
			mockBehavior: func(m *mockUserUsecase) {
				m.mockVerifyToken = func(tokenString string) (*jwt.Token, error) {
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
						"email": "test@example.com",
					})
					return token, nil
				}
			},
			expectedCode:    http.StatusOK,
			expectedMessage: "トークンは有効です",
		},
		{
			name:  "異常系：無効なトークン",
			token: "Bearer invalid-token",
			mockBehavior: func(m *mockUserUsecase) {
				m.mockVerifyToken = func(tokenString string) (*jwt.Token, error) {
					return nil, errors.New(errors.AuthenticationError, "無効なトークンです", http.StatusUnauthorized, nil)
				}
			},
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: "未認証",
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
			req := httptest.NewRequest(http.MethodGet, "/verify", nil)
			req.Header.Set("Authorization", tt.token)
			c := e.NewContext(req, rec)

			// テスト実行
			err := uc.VerifyToken(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assertSuccessResponse(t, rec, tt.expectedCode, tt.expectedMessage)

			// レコーダーをリセット
			rec.Body.Reset()
		})
	}
}
