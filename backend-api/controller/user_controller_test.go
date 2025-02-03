package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taaag51/smart-pantry/backend-api/model"
)

type mockUserUsecase struct {
	mock.Mock
}

func (m *mockUserUsecase) SignUp(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *mockUserUsecase) Login(user model.User) (*model.TokenPair, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TokenPair), args.Error(1)
}

func (m *mockUserUsecase) VerifyToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func (m *mockUserUsecase) RefreshTokens(refreshToken string) (*model.TokenPair, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TokenPair), args.Error(1)
}

func TestSignUp(t *testing.T) {
	e := echo.New()
	mock := new(mockUserUsecase)
	controller := NewUserController(mock)

	t.Run("正常系：ユーザー登録成功", func(t *testing.T) {
		user := model.User{
			Email:    "test@example.com",
			Password: "password123",
		}

		mock.On("SignUp", user).Return(user, nil).Once()

		jsonStr := `{"email":"test@example.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, controller.SignUp(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			var response struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "success", response.Status)
			assert.Equal(t, "ユーザーが正常に作成されました", response.Message)
		}
	})
}

func TestLogin(t *testing.T) {
	e := echo.New()
	mock := new(mockUserUsecase)
	controller := NewUserController(mock)

	t.Run("正常系：ログイン成功", func(t *testing.T) {
		user := model.User{
			Email:    "test@example.com",
			Password: "password123",
		}

		tokenPair := &model.TokenPair{
			AccessToken:   "access_token",
			RefreshToken:  "refresh_token",
			AccessExpiry:  time.Now().Add(15 * time.Minute),
			RefreshExpiry: time.Now().Add(7 * 24 * time.Hour),
		}

		mock.On("Login", user).Return(tokenPair, nil).Once()

		jsonStr := `{"email":"test@example.com","password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(jsonStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, controller.LogIn(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response struct {
				Status  string `json:"status"`
				Message string `json:"message"`
				Data    struct {
					AccessToken string    `json:"accessToken"`
					TokenType   string    `json:"tokenType"`
					ExpiresIn   int       `json:"expiresIn"`
					ExpiresAt   time.Time `json:"expiresAt"`
				} `json:"data"`
			}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "success", response.Status)
			assert.Equal(t, "ログインに成功しました", response.Message)
			assert.Equal(t, tokenPair.AccessToken, response.Data.AccessToken)
			assert.Equal(t, "Bearer", response.Data.TokenType)
		}
	})
}

func TestVerifyToken(t *testing.T) {
	e := echo.New()
	mock := new(mockUserUsecase)
	controller := NewUserController(mock)

	t.Run("正常系：トークン検証成功", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": "test@example.com",
			"exp":   time.Now().Add(time.Hour).Unix(),
		})

		mock.On("VerifyToken", "valid_token").Return(token, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/verify-token", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, controller.VerifyToken(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response struct {
				Status  string `json:"status"`
				Message string `json:"message"`
				Data    struct {
					Email string `json:"email"`
				} `json:"data"`
			}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "success", response.Status)
			assert.Equal(t, "トークンは有効です", response.Message)
			assert.Equal(t, "test@example.com", response.Data.Email)
		}
	})
}

func TestRefreshToken(t *testing.T) {
	e := echo.New()
	mock := new(mockUserUsecase)
	controller := NewUserController(mock)

	t.Run("正常系：トークン更新成功", func(t *testing.T) {
		refreshToken := "valid_refresh_token"
		newTokenPair := &model.TokenPair{
			AccessToken:   "new_access_token",
			RefreshToken:  "new_refresh_token",
			AccessExpiry:  time.Now().Add(15 * time.Minute),
			RefreshExpiry: time.Now().Add(7 * 24 * time.Hour),
		}

		mock.On("RefreshTokens", refreshToken).Return(newTokenPair, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/refresh-token", nil)
		cookie := &http.Cookie{
			Name:  "refresh_token",
			Value: refreshToken,
		}
		req.AddCookie(cookie)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, controller.RefreshToken(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var response struct {
				Status  string `json:"status"`
				Message string `json:"message"`
				Data    struct {
					AccessToken string    `json:"accessToken"`
					TokenType   string    `json:"tokenType"`
					ExpiresIn   int       `json:"expiresIn"`
					ExpiresAt   time.Time `json:"expiresAt"`
				} `json:"data"`
			}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "success", response.Status)
			assert.Equal(t, "トークンを更新しました", response.Message)
			assert.Equal(t, newTokenPair.AccessToken, response.Data.AccessToken)
			assert.Equal(t, "Bearer", response.Data.TokenType)
		}
	})
}
