package response

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SetCookie はクッキーを設定する
func SetCookie(c echo.Context, cookie *http.Cookie) {
	http.SetCookie(c.Response(), cookie)
}

// NewAuthCookie は新しい認証用クッキーを作成する
func NewAuthCookie(token string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expires,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// ClearAuthCookie は認証用クッキーをクリアする
func ClearAuthCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// HandleSuccess は成功レスポンスを返す
func HandleSuccess(c echo.Context, status int, message string) error {
	return c.JSON(status, Response{
		Status:  "success",
		Message: message,
	})
}

// HandleSuccessWithData は成功レスポンスとデータを返す
func HandleSuccessWithData(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// HandleError はエラーレスポンスを返す
func HandleError(c echo.Context, status int, message string) error {
	return c.JSON(status, Response{
		Status:  "error",
		Message: message,
	})
}
