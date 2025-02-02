package response

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// レスポンス型の定義
type ErrorResponse struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type AuthCookie struct {
	Name     string
	Value    string
	Path     string
	Expires  time.Time
	HTTPOnly bool
	SameSite http.SameSite
}

func NewAuthCookie(token string, expires time.Time) *AuthCookie {
	return &AuthCookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HTTPOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func ClearAuthCookie() *AuthCookie {
	return &AuthCookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-24 * time.Hour),
		HTTPOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func SetCookie(c echo.Context, cookie *AuthCookie) {
	c.SetCookie(&http.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Expires:  cookie.Expires,
		HttpOnly: cookie.HTTPOnly,
		SameSite: cookie.SameSite,
	})
}

// HandleError は共通のエラーレスポンスを生成します
func HandleError(c echo.Context, status int, message string) error {
	return c.JSON(status, ErrorResponse{
		Message: message,
	})
}

// HandleErrorWithType は型情報を含むエラーレスポンスを生成します
func HandleErrorWithType(c echo.Context, status int, errorType string, message string) error {
	return c.JSON(status, ErrorResponse{
		Type:    errorType,
		Message: message,
	})
}

// HandleSuccess は共通の成功レスポンスを生成します
func HandleSuccess(c echo.Context, status int, message string) error {
	return c.JSON(status, SuccessResponse{
		Message: message,
	})
}

// HandleSuccessWithData はデータを含む成功レスポンスを生成します
func HandleSuccessWithData(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, SuccessResponse{
		Message: message,
		Data:    data,
	})
}
