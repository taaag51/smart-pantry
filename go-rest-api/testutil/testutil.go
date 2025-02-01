package testutil

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

// TestDB はテスト用のデータベース接続を保持する構造体
type TestDB struct {
	DB *sql.DB
}

// NewTestDB はテスト用のデータベース接続を作成する
func NewTestDB(t *testing.T) *TestDB {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("データベース接続エラー: %v", err)
	}

	// コネクション確認
	if err := db.Ping(); err != nil {
		t.Fatalf("データベースPingエラー: %v", err)
	}

	return &TestDB{DB: db}
}

// CleanupDB はテストデータベースをクリーンアップする
func (tdb *TestDB) CleanupDB(t *testing.T) {
	tables := []string{"food_items", "recipes", "users"}
	for _, table := range tables {
		_, err := tdb.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Errorf("テーブル %s のクリーンアップに失敗: %v", table, err)
		}
	}
}

// Close はデータベース接続を閉じる
func (tdb *TestDB) Close() {
	if err := tdb.DB.Close(); err != nil {
		log.Printf("データベース接続のクローズに失敗: %v", err)
	}
}

// GenerateTestJWT はテスト用のJWTトークンを生成する
func GenerateTestJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	return token.SignedString(secretKey)
}

// ParseJWTToken はJWTトークンをパースしてクレームを返す
func ParseJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("トークンのパースに失敗: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("無効なトークン")
}

// CreateTestUser はテスト用のユーザーを作成する
func (tdb *TestDB) CreateTestUser(t *testing.T) (uint, string) {
	var userID uint
	err := tdb.DB.QueryRow(`
		INSERT INTO users (email, password_hash, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`, "test@example.com", "hashedpassword").Scan(&userID)

	if err != nil {
		t.Fatalf("テストユーザーの作成に失敗: %v", err)
	}

	token, err := GenerateTestJWT(userID)
	if err != nil {
		t.Fatalf("JWTトークンの生成に失敗: %v", err)
	}

	return userID, token
}

// CreateTestFoodItem はテスト用の食材データを作成する
func (tdb *TestDB) CreateTestFoodItem(t *testing.T, userID uint) uint {
	var itemID uint
	err := tdb.DB.QueryRow(`
		INSERT INTO food_items (name, quantity, expiry_date, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`, "テスト食材", 1, time.Now().AddDate(0, 0, 7), userID).Scan(&itemID)

	if err != nil {
		t.Fatalf("テスト食材の作成に失敗: %v", err)
	}

	return itemID
}

// AssertHTTPStatus はHTTPステータスコードを検証する
func AssertHTTPStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("HTTPステータスコードが一致しません: got %d, want %d", got, want)
	}
}

// AssertErrorMessage はエラーメッセージを検証する
func AssertErrorMessage(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("エラーメッセージが一致しません: got %q, want %q", got, want)
	}
}
