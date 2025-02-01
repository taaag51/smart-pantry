// Package errors provides custom error types and error handling utilities
package errors

import "net/http"

// ErrorType represents the type of error
type ErrorType string

const (
	// ValidationError indicates invalid input data
	ValidationError ErrorType = "VALIDATION_ERROR"
	// AuthenticationError indicates authentication failure
	AuthenticationError ErrorType = "AUTHENTICATION_ERROR"
	// DatabaseError indicates database operation failure
	DatabaseError ErrorType = "DATABASE_ERROR"
	// BusinessError indicates business rule violation
	BusinessError ErrorType = "BUSINESS_ERROR"
)

// AppError represents an application error
type AppError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	HTTPStatus int       `json:"-"`
	Err        error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(errType ErrorType, message string, httpStatus int, err error) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

// Common validation errors
var (
	InvalidEmail = New(
		ValidationError,
		"メールアドレスが無効です",
		http.StatusBadRequest,
		nil,
	)

	InvalidCredentials = New(
		AuthenticationError,
		"メールアドレスまたはパスワードが正しくありません",
		http.StatusUnauthorized,
		nil,
	)

	EmailExists = New(
		BusinessError,
		"このメールアドレスは既に登録されています",
		http.StatusBadRequest,
		nil,
	)

	HashPasswordError = New(
		BusinessError,
		"パスワードのハッシュ化に失敗しました",
		http.StatusInternalServerError,
		nil,
	)

	CreateUserError = New(
		DatabaseError,
		"ユーザーの作成に失敗しました",
		http.StatusInternalServerError,
		nil,
	)

	GenerateTokenError = New(
		AuthenticationError,
		"トークンの生成に失敗しました",
		http.StatusInternalServerError,
		nil,
	)
)

// AsAppError converts a standard error to an AppError if possible
func AsAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	// Default to internal server error
	return New(
		BusinessError,
		err.Error(),
		http.StatusInternalServerError,
		err,
	)
}

// GetHTTPStatus returns the HTTP status code for an error
func GetHTTPStatus(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}
