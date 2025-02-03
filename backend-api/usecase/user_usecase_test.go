package usecase

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taaag51/smart-pantry/backend-api/model"
)

func init() {
	// テスト用のシークレットキーを設定
	os.Setenv("SECRET", "test-secret-key")
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id uint) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

type MockUserValidator struct {
	mock.Mock
}

func (m *MockUserValidator) ValidateUser(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserValidator) ValidateLogin(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestSignUp(t *testing.T) {
	ur := new(MockUserRepository)
	uv := new(MockUserValidator)
	uu := NewUserUsecase(ur, uv)

	t.Run("正常系：ユーザー登録成功", func(t *testing.T) {
		user := model.User{
			Email:    "test@example.com",
			Password: "password123",
		}

		uv.On("ValidateUser", user).Return(nil)
		ur.On("CreateUser", &user).Return(user, nil)

		createdUser, err := uu.SignUp(user)

		assert.NoError(t, err)
		assert.Equal(t, user.Email, createdUser.Email)
		ur.AssertExpectations(t)
		uv.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	ur := new(MockUserRepository)
	uv := new(MockUserValidator)
	uu := NewUserUsecase(ur, uv)

	t.Run("正常系：ログイン成功", func(t *testing.T) {
		user := model.User{
			Email:    "test@example.com",
			Password: "password123",
		}

		uv.On("ValidateLogin", user).Return(nil)
		ur.On("GetUserByEmail", user.Email).Return(user, nil)

		tokenPair, err := uu.Login(user)

		assert.NoError(t, err)
		assert.NotEmpty(t, tokenPair.AccessToken)
		assert.NotEmpty(t, tokenPair.RefreshToken)
		assert.True(t, tokenPair.AccessExpiry.After(time.Now()))
		assert.True(t, tokenPair.RefreshExpiry.After(time.Now()))
		ur.AssertExpectations(t)
		uv.AssertExpectations(t)
	})
}

func TestVerifyToken(t *testing.T) {
	ur := new(MockUserRepository)
	uv := new(MockUserValidator)
	uu := NewUserUsecase(ur, uv)

	t.Run("正常系：トークン検証成功", func(t *testing.T) {
		email := "test@example.com"
		tokenString, _, err := generateToken(email, "access", 15)
		assert.NoError(t, err)

		token, err := uu.VerifyToken(tokenString)
		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, email, claims["email"])
	})

	t.Run("異常系：無効なトークン", func(t *testing.T) {
		_, err := uu.VerifyToken("invalid.token.string")
		assert.Error(t, err)
	})
}

func TestRefreshTokens(t *testing.T) {
	ur := new(MockUserRepository)
	uv := new(MockUserValidator)
	uu := NewUserUsecase(ur, uv)

	t.Run("正常系：トークン更新成功", func(t *testing.T) {
		email := "test@example.com"
		user := model.User{Email: email}

		refreshToken, _, err := generateToken(email, "refresh", 24*7)
		assert.NoError(t, err)

		uv.On("ValidateLogin", mock.AnythingOfType("model.User")).Return(nil)
		ur.On("GetUserByEmail", email).Return(user, nil)

		tokenPair, err := uu.RefreshTokens(refreshToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenPair.AccessToken)
		assert.NotEmpty(t, tokenPair.RefreshToken)
		assert.True(t, tokenPair.AccessExpiry.After(time.Now()))
		assert.True(t, tokenPair.RefreshExpiry.After(time.Now()))

		ur.AssertExpectations(t)
		uv.AssertExpectations(t)
	})

	t.Run("異常系：無効なリフレッシュトークン", func(t *testing.T) {
		_, err := uu.RefreshTokens("invalid.refresh.token")
		assert.Error(t, err)
	})
}

func TestMain(m *testing.M) {
	// テスト用の環境変数を設定
	os.Setenv("SECRET", "test-secret-key")

	// テストを実行
	code := m.Run()

	// テスト終了後に環境変数をクリア
	os.Unsetenv("SECRET")

	os.Exit(code)
}
