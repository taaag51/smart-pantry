// backend/internal/services/auth_service.go
package services

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"smart-pantry/backend/internal/models"
)

// AuthService は認証関連のサービスです。
type AuthService struct {
	db *gorm.DB
}

// NewAuthService はAuthServiceを生成します。
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// Login はログイン処理を行います。
func (as *AuthService) Login(username, password string) error {
	// TODO: ログイン処理の実装
	return nil
}

// Register は登録処理を行います。
func (as *AuthService) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	result := as.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Logout はログアウト処理を行います。
func (as *AuthService) Logout() error {
	// TODO: ログアウト処理の実装
	return nil
}
