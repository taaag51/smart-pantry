package controller

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/taaag51/smart-pantry/backend-api/model"
	"github.com/taaag51/smart-pantry/backend-api/usecase"

	"github.com/taaag51/smart-pantry/backend-api/controller/response"
	"github.com/taaag51/smart-pantry/backend-api/errors"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error
	VerifyToken(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

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

func (uc *userController) SignUp(c echo.Context) error {
	user, err := bindUser(c)
	if err != nil {
		return response.HandleError(c, http.StatusInternalServerError, "ユーザーの登録に失敗しました")
	}

	_, err = uc.uu.SignUp(user)
	if err != nil {
		return response.HandleError(c, http.StatusInternalServerError, "ユーザーの登録に失敗しました")
	}

	return response.HandleSuccess(c, http.StatusCreated, "ユーザーが正常に作成されました")
}

func (uc *userController) LogIn(c echo.Context) error {
	user, err := bindUser(c)
	if err != nil {
		return response.HandleError(c, http.StatusBadRequest, "リクエストの形式が不正です")
	}

	tokenString, err := uc.uu.Login(user)
	if err != nil {
		return response.HandleError(c, http.StatusUnauthorized, "メールアドレスまたはパスワードが正しくありません")
	}

	response.SetCookie(c, response.NewAuthCookie(
		tokenString,
		time.Now().Add(24*time.Hour),
	))

	return response.HandleSuccessWithData(c, http.StatusOK, "ログインに成功しました", map[string]string{
		"token": tokenString,
	})
}

func (uc *userController) LogOut(c echo.Context) error {
	response.SetCookie(c, response.ClearAuthCookie())
	return response.HandleSuccess(c, http.StatusOK, "ログアウトしました")
}

func (uc *userController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	c.Response().Header().Set("X-CSRF-Token", token)
	return response.HandleSuccessWithData(c, http.StatusOK, "CSRFトークンを取得しました", map[string]string{
		"csrf_token": token,
	})
}

func (uc *userController) VerifyToken(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return response.HandleError(c, http.StatusUnauthorized, "未認証")
	}

	if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
		return response.HandleError(c, http.StatusUnauthorized, "未認証")
	}
	tokenString = tokenString[7:]

	token, err := uc.uu.VerifyToken(tokenString)
	if err != nil {
		return response.HandleError(c, http.StatusUnauthorized, "未認証")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return response.HandleError(c, http.StatusUnauthorized, "未認証")
	}

	return response.HandleSuccessWithData(c, http.StatusOK, "トークンは有効です", map[string]interface{}{
		"email": claims["email"],
	})
}
