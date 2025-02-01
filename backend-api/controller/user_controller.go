/*
Package controller handles HTTP request processing and response generation.
This package defines the IUserController interface and implements user-related HTTP request handling.
*/
package controller

import (
	"backend-api/controller/response"
	"backend-api/errors"
	"backend-api/model"
	"backend-api/usecase"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

/*
IUserController defines the interface for user-related HTTP request handling.
It includes methods for user registration, authentication, session management, and CSRF token handling.
*/
type IUserController interface {
	/*
		SignUp handles user registration.
		It binds user data from the request, validates it, and calls the use case to create a new user.
		Returns a success response or an error if the registration fails.
	*/
	SignUp(c echo.Context) error
	/*
		LogIn authenticates a user and creates a session.
		It binds user data from the request, validates it, and calls the use case to log in the user.
		On success, it sets an authentication cookie and returns a success response.
	*/
	LogIn(c echo.Context) error
	/*
		LogOut terminates the current user session.
		It clears the authentication cookie and returns a success response.
	*/
	LogOut(c echo.Context) error
	/*
		CsrfToken provides a CSRF token for form submission.
		It retrieves the CSRF token from the context and returns a success response.
	*/
	CsrfToken(c echo.Context) error
	/*
		VerifyToken validates the current session token.
		This method is called after the JWT middleware has already validated the token.
		It returns a success response indicating the token is valid.
	*/
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
		return handleError(c, http.StatusInternalServerError, "ユーザーの登録に失敗しました")
	}

	// Call signup usecase
	_, err = uc.uu.SignUp(user)
	if err != nil {
		return handleError(c, http.StatusInternalServerError, "ユーザーの登録に失敗しました")
	}

	return handleError(c, http.StatusCreated, "ユーザーが正常に作成されました")
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
		return handleError(c, http.StatusBadRequest, "リクエストの形式が不正です")
	}

	// Authenticate user and get token
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		return handleError(c, http.StatusUnauthorized, "ログインに失敗しました")
	}
	if err != nil {
		return handleError(c, http.StatusUnauthorized, "ログインに失敗しました")
	}
	if err != nil {
		return handleError(c, http.StatusUnauthorized, "ログインに失敗しました")
	}

	// Set authentication cookie
	response.SetCookie(c, response.NewAuthCookie(
		tokenString,
		time.Now().Add(24*time.Hour),
	))

	return handleError(c, http.StatusOK, "ログインに成功しました")
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

	return handleError(c, http.StatusOK, "ログアウトしました")
}

// CsrfToken godoc
// @Summary Get CSRF token
// @Description Provides a CSRF token for form submission
// @Tags auth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Router /csrf [get]
func (uc *userController) CsrfToken(c echo.Context) error {
	_ = c.Get("csrf").(string)
	return handleError(c, http.StatusOK, "CSRFトークンを取得しました")
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
	return handleError(c, http.StatusOK, "トークンは有効です")
}
