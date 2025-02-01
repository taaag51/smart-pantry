// Package testutil はテスト用のユーティリティ関数を提供します
package testutil

import (
	"testing"

	"github.com/taaag51/smart-pantry/backend-api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDB はテスト用のデータベース接続を管理する構造体です
type TestDB struct {
	*gorm.DB
}

// GetTestDSN はテスト用のデータベース接続文字列を返します
func GetTestDSN() string {
	return "host=localhost user=test_user password=test_pass dbname=test_db port=5432 sslmode=disable"
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
