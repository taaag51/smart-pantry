// Package testutil はテスト用のユーティリティ関数を提供します
package testutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/taaag51/smart-pantry/backend-api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ParseJWTToken はJWTトークンをパースし、クレームを返します
func ParseJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}

// TestDB はテスト用のデータベース接続を管理する構造体です
type TestDB struct {
	*gorm.DB
}

// GetTestDSN はテスト用のデータベース接続文字列を返します
func GetTestDSN() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "test_user"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "test_password"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "smart_pantry_test"
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		host, user, password, dbName)
}

// NewTestDB は新しいテストデータベース接続を作成します
func NewTestDB(t *testing.T) *TestDB {
	db, err := gorm.Open(postgres.Open(GetTestDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("テストDB接続に失敗: %v", err)
	}

	// テストデータベースのマイグレーション
	err = db.AutoMigrate(&model.User{}, &model.FoodItem{})
	if err != nil {
		t.Fatalf("マイグレーション失敗: %v", err)
	}

	return &TestDB{db}
}

// Close はデータベース接続を閉じます
func (tdb *TestDB) Close() {
	sqlDB, err := tdb.DB.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}

// CreateTestUser はテストユーザーを作成し、そのIDを返します
func (tdb *TestDB) CreateTestUser(t *testing.T) (uint, error) {
	user := model.User{
		Email:    "test@example.com",
		Password: "testpass",
	}

	if err := tdb.Create(&user).Error; err != nil {
		t.Fatalf("テストユーザー作成失敗: %v", err)
		return 0, err
	}

	return user.ID, nil
}
