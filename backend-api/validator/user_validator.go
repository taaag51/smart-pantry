package validator

import (
	"github.com/taaag51/smart-pantry/backend-api/model"
)

type IUserValidator interface {
	ValidateUser(user model.User) error
	ValidateLogin(user model.User) error
}

type userValidator struct{}

func NewUserValidator() IUserValidator {
	return &userValidator{}
}

func (uv *userValidator) ValidateUser(user model.User) error {
	if user.Email == "" {
		return &ValidationError{Message: "メールアドレスは必須です"}
	}
	if user.Password == "" {
		return &ValidationError{Message: "パスワードは必須です"}
	}
	return nil
}

func (uv *userValidator) ValidateLogin(user model.User) error {
	if user.Email == "" {
		return &ValidationError{Message: "メールアドレスは必須です"}
	}
	if user.Password == "" {
		return &ValidationError{Message: "パスワードは必須です"}
	}
	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
