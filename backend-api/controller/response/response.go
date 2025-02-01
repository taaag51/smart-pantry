// Package response provides common HTTP response handling utilities for controllers
package response

import (
	"backend-api/errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Error sends a standardized error response
func Error(c echo.Context, err error) error {
	appErr := errors.AsAppError(err)
	return c.JSON(appErr.HTTPStatus, ErrorResponse{
		Type:    string(appErr.Type),
		Message: appErr.Message,
		Code:    appErr.HTTPStatus,
	})
}

// Success sends a standardized success response
func Success(c echo.Context, status int, data interface{}, message string) error {
	return c.JSON(status, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// BadRequest is a helper for 400 Bad Request responses
func BadRequest(c echo.Context, message string) error {
	return Error(c, errors.New(errors.ValidationError, message, http.StatusBadRequest, nil))
}

// Unauthorized is a helper for 401 Unauthorized responses
func Unauthorized(c echo.Context, message string) error {
	return Error(c, errors.New(errors.AuthenticationError, message, http.StatusUnauthorized, nil))
}

// NotFound is a helper for 404 Not Found responses
func NotFound(c echo.Context, message string) error {
	return Error(c, errors.New(errors.BusinessError, message, http.StatusNotFound, nil))
}

// InternalServerError is a helper for 500 Internal Server Error responses
func InternalServerError(c echo.Context, err error) error {
	return Error(c, errors.New(errors.BusinessError, "内部サーバーエラーが発生しました", http.StatusInternalServerError, err))
}

// CookieConfig contains configuration for HTTP cookies
type CookieConfig struct {
	Name     string
	Value    string
	Expires  time.Time
	Path     string
	Domain   string
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}

// DefaultCookieConfig returns the default cookie configuration
func DefaultCookieConfig() CookieConfig {
	return CookieConfig{
		Path:     "/",
		Domain:   "localhost", // TODO: Make this configurable
		Secure:   true,
		HTTPOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
}

// SetCookie sets an HTTP cookie with the given configuration
func SetCookie(c echo.Context, config CookieConfig) {
	cookie := new(http.Cookie)
	cookie.Name = config.Name
	cookie.Value = config.Value
	cookie.Expires = config.Expires
	cookie.Path = config.Path
	cookie.Domain = config.Domain
	cookie.Secure = config.Secure
	cookie.HttpOnly = config.HTTPOnly
	cookie.SameSite = config.SameSite
	c.SetCookie(cookie)
}

// NewAuthCookie creates a new authentication cookie
func NewAuthCookie(token string, expires time.Time) CookieConfig {
	config := DefaultCookieConfig()
	config.Name = "token"
	config.Value = token
	config.Expires = expires
	return config
}

// ClearAuthCookie creates a cookie configuration that will clear the auth cookie
func ClearAuthCookie() CookieConfig {
	config := DefaultCookieConfig()
	config.Name = "token"
	config.Value = ""
	config.Expires = time.Now()
	return config
}
