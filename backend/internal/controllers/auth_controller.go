// backend/internal/controllers/auth_controller.go
package controllers

import "net/http"

// AuthController は認証関連のコントローラーです。
type AuthController struct {
	// AuthService 認証サービス
}

// NewAuthController はAuthControllerを生成します。
func NewAuthController() *AuthController {
	return &AuthController{}
}

// Login はログイン処理を行います。
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: ログイン処理の実装
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login"))
}

// Register は登録処理を行います。
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: 登録処理の実装
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Register"))
}

// Logout はログアウト処理を行います。
func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: ログアウト処理の実装
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout"))
}
