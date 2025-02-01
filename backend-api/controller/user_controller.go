// Package controller handles HTTP request processing and response generation
package controller

import (
	"backend-api/controller/response"
	"backend-api/errors"
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// IUserController defines the interface for user-related HTTP request handling
type IUserController interface {
	// SignUp handles user registration
	SignUp(c echo.Context) error
	// LogIn authenticates a user and creates a session
	LogIn(c echo.Context) error
	// LogOut terminates the current user session
	LogOut(c echo.Context) error
	// CsrfToken provides a CSRF token for form submission
	CsrfToken(c echo.Context) error
	// VerifyToken validates the current session token
	VerifyToken(c echo.Context) error
}

// userController implements IUserController interface
type userController struct {
	uu usecase.IUserUsecase
}

// NewUserController creates a new instance of IUserController
func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

// bindUser binds and validates user data from request
func bindUser(c echo.Context) (model.User, error) {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return user, errors.New(
			errors.ValidationError,
			"リクエストの形式が不正です",
			http.StatusBadRequest,
			err,
		)
	}
	return user, nil
}

// SignUp godoc
// @Summary Register a new user
// @Description Creates a new user account with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.User true "User registration details"
// @Success 201 {object} model.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /signup [post]
func (uc *userController) SignUp(c echo.Context) error {
	// Bind and validate user data
	user, err := bindUser(c)
	if err != nil {
		return response.Error(c, err)
	}

	// Call signup usecase
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, userRes, "ユーザーが正常に作成されました")
}

// LogIn godoc
// @Summary Authenticate a user
// @Description Authenticates user credentials and establishes a session
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.User true "User login credentials"
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /login [post]
func (uc *userController) LogIn(c echo.Context) error {
	// Bind and validate user data
	user, err := bindUser(c)
	if err != nil {
		return response.Error(c, err)
	}

	// Authenticate user and get token
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		return response.Error(c, err)
	}

	// Set authentication cookie
	response.SetCookie(c, response.NewAuthCookie(
		tokenString,
		time.Now().Add(24*time.Hour),
	))

	return response.Success(c, http.StatusOK, nil, "ログインに成功しました")
}

// LogOut godoc
// @Summary End user session
// @Description Terminates the current user session
// @Tags auth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /logout [post]
func (uc *userController) LogOut(c echo.Context) error {
	// Clear authentication cookie
	response.SetCookie(c, response.ClearAuthCookie())

	return response.Success(c, http.StatusOK, nil, "ログアウトしました")
}

// CsrfToken godoc
// @Summary Get CSRF token
// @Description Provides a CSRF token for form submission
// @Tags auth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /csrf [get]
func (uc *userController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return response.Success(c, http.StatusOK, map[string]string{"csrf_token": token}, "")
}

// VerifyToken godoc
// @Summary Verify JWT token
// @Description Validates the current session token
// @Tags auth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Security ApiKeyAuth
// @Router /verify [get]
func (uc *userController) VerifyToken(c echo.Context) error {
	// JWTミドルウェアによって既に検証されているため、
	// このエンドポイントに到達できた時点で有効なトークン
	return response.Success(c, http.StatusOK, nil, "トークンは有効です")
}
