package usecase

import (
	"backend-api/model"
	"backend-api/testutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// モックリポジトリの定義
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(user *model.User, email string) error {
	args := m.Called(user, email)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	// モックユーザーデータの設定
	if mockUser, ok := args.Get(1).(*model.User); ok && mockUser != nil {
		*user = *mockUser
	}
	return nil
}

// モックバリデータの定義
type MockUserValidator struct {
	mock.Mock
}

func (m *MockUserValidator) UserValidate(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestUserUsecase_SignUp(t *testing.T) {
	// テスト用の環境変数設定
	os.Setenv("SECRET", "test_secret_key")
	defer os.Unsetenv("SECRET")

	tests := []struct {
		name       string
		user       model.User
		setupMocks func(*MockUserRepository, *MockUserValidator)
		wantErr    bool
		errMessage string
	}{
		{
			name: "正常な新規登録",
			user: model.User{
				Email:    "test@example.com",
				Password: "validPassword123",
			},
			setupMocks: func(ur *MockUserRepository, uv *MockUserValidator) {
				uv.On("UserValidate", mock.Anything).Return(nil)
				ur.On("GetUserByEmail", mock.Anything, "test@example.com").Return(assert.AnError) // ユーザーが存在しない
				ur.On("CreateUser", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "重複メールアドレス",
			user: model.User{
				Email:    "existing@example.com",
				Password: "validPassword123",
			},
			setupMocks: func(ur *MockUserRepository, uv *MockUserValidator) {
				uv.On("UserValidate", mock.Anything).Return(nil)
				ur.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(nil)
			},
			wantErr:    true,
			errMessage: "email already exists",
		},
		{
			name: "バリデーションエラー",
			user: model.User{
				Email:    "invalid-email",
				Password: "short",
			},
			setupMocks: func(ur *MockUserRepository, uv *MockUserValidator) {
				uv.On("UserValidate", mock.Anything).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			ur := new(MockUserRepository)
			uv := new(MockUserValidator)
			tt.setupMocks(ur, uv)

			// ユースケースの作成
			uu := NewUserUsecase(ur, uv)

			// テスト実行
			_, err := uu.SignUp(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
			} else {
				assert.NoError(t, err)
			}

			// モックの検証
			ur.AssertExpectations(t)
			uv.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_Login(t *testing.T) {
	// テスト用の環境変数設定
	os.Setenv("SECRET", "test_secret_key")
	defer os.Unsetenv("SECRET")

	// テスト用のハッシュ化されたパスワード
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctPassword123"), 10)

	tests := []struct {
		name       string
		user       model.User
		setupMocks func(*MockUserRepository, *MockUserValidator)
		wantErr    bool
		errMessage string
	}{
		{
			name: "正常なログイン",
			user: model.User{
				Email:    "test@example.com",
				Password: "correctPassword123",
			},
			setupMocks: func(ur *MockUserRepository, uv *MockUserValidator) {
				uv.On("UserValidate", mock.Anything).Return(nil)
				ur.On("GetUserByEmail", mock.Anything, "test@example.com").Return(nil, &model.User{
					ID:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				})
			},
			wantErr: false,
		},
		{
			name: "存在しないユーザー",
			user: model.User{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMocks: func(ur *MockUserRepository, uv *MockUserValidator) {
				uv.On("UserValidate", mock.Anything).Return(nil)
				ur.On("GetUserByEmail", mock.Anything, "nonexistent@example.com").Return(assert.AnError)
			},
			wantErr:    true,
			errMessage: "invalid email or password",
		},
		{
			name: "パスワード不一致",
			user: model.User{
				Email:    "test@example.com",
				Password: "wrongPassword123",
			},
			setupMocks: func(ur *MockUserRepository, uv *MockUserValidator) {
				uv.On("UserValidate", mock.Anything).Return(nil)
				ur.On("GetUserByEmail", mock.Anything, "test@example.com").Return(nil, &model.User{
					ID:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				})
			},
			wantErr:    true,
			errMessage: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			ur := new(MockUserRepository)
			uv := new(MockUserValidator)
			tt.setupMocks(ur, uv)

			// ユースケースの作成
			uu := NewUserUsecase(ur, uv)

			// テスト実行
			token, err := uu.Login(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// JWTトークンの検証
				claims, err := testutil.ParseJWTToken(token)
				assert.NoError(t, err)
				assert.Equal(t, uint(1), uint(claims["user_id"].(float64)))
				assert.Equal(t, "test@example.com", claims["email"].(string))

				// 有効期限の検証
				exp := time.Unix(int64(claims["exp"].(float64)), 0)
				assert.True(t, exp.After(time.Now()))
				assert.True(t, exp.Before(time.Now().Add(13*time.Hour))) // 12時間 + バッファ
			}

			// モックの検証
			ur.AssertExpectations(t)
			uv.AssertExpectations(t)
		})
	}
}
