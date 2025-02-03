package controller

import (
	"log" // 追加
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/taaag51/smart-pantry/backend-api/controller/response"
	"github.com/taaag51/smart-pantry/backend-api/errors"
	"github.com/taaag51/smart-pantry/backend-api/model"
	"github.com/taaag51/smart-pantry/backend-api/usecase"
)

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error
	VerifyToken(c echo.Context) error
	RefreshToken(c echo.Context) error
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
		return response.HandleError(c, http.StatusBadRequest, "リクエストの形式が不正です")
	}

	_, err = uc.uu.SignUp(user)
	if err != nil {
		return response.HandleError(c, http.StatusInternalServerError, "ユーザーの登録に失敗しました")
	}

	return response.HandleSuccess(c, http.StatusCreated, "ユーザーが正常に作成されました")
}

func (uc *userController) LogIn(c echo.Context) error {
	log.Printf("ログインリクエスト - Method: %s, Headers: %v", c.Request().Method, c.Request().Header)
	user, err := bindUser(c)
	if err != nil {
		return response.HandleError(c, http.StatusBadRequest, "リクエストの形式が不正です")
	}

	tokenPair, err := uc.uu.Login(user)
	if err != nil {
		log.Printf("ログインエラー: %v", err)
		return response.HandleError(c, http.StatusUnauthorized, "メールアドレスまたはパスワードが正しくありません")
	}
	log.Printf("ログイン成功: %s", user.Email)

	// アクセストークンをCookieに設定
	response.SetCookie(c, response.NewAuthCookie(
		tokenPair.AccessToken,
		tokenPair.AccessExpiry,
	))

	// リフレッシュトークンをCookieに設定（HTTPOnly）
	response.SetCookie(c, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Expires:  tokenPair.RefreshExpiry,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return response.HandleSuccessWithData(c, http.StatusOK, "ログインに成功しました", &model.TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(time.Until(tokenPair.AccessExpiry).Seconds()),
		ExpiresAt:    tokenPair.AccessExpiry,
		RefreshToken: tokenPair.RefreshToken,
	})
}

func (uc *userController) LogOut(c echo.Context) error {
	// アクセストークンとリフレッシュトークンのCookieを削除
	response.SetCookie(c, response.ClearAuthCookie())
	response.SetCookie(c, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

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
	log.Printf("トークン検証リクエスト - Headers: %v", c.Request().Header)

	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		log.Printf("Authorizationヘッダーが見つかりません")
		return response.HandleError(c, http.StatusUnauthorized, "未認証")
	}

	if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
		log.Printf("不正なトークン形式: %s", tokenString)
		return response.HandleError(c, http.StatusUnauthorized, "未認証")
	}
	tokenString = tokenString[7:]
	log.Printf("検証するトークン: %s", tokenString)

	// トークンの検証のみを行い、新しいトークンは発行しない
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

func (uc *userController) RefreshToken(c echo.Context) error {
	// リフレッシュトークンをCookieから取得
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		return response.HandleError(c, http.StatusUnauthorized, "リフレッシュトークンが見つかりません。再度ログインしてください。")
	}

	if refreshCookie.Value == "" {
		return response.HandleError(c, http.StatusUnauthorized, "無効なリフレッシュトークンです。再度ログインしてください。")
	}

	// リフレッシュトークンを使用して新しいトークンペアを生成
	tokenPair, err := uc.uu.RefreshTokens(refreshCookie.Value)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			return response.HandleError(c, http.StatusUnauthorized, "リフレッシュトークンが期限切れです。再度ログインしてください。")
		} else if err == jwt.ErrSignatureInvalid {
			return response.HandleError(c, http.StatusUnauthorized, "無効なリフレッシュトークンです。再度ログインしてください。")
		} else {
			return response.HandleError(c, http.StatusUnauthorized, "トークンの更新に失敗しました。再度ログインしてください。")
		}
	}

	// 新しいアクセストークンをCookieに設定
	response.SetCookie(c, response.NewAuthCookie(
		tokenPair.AccessToken,
		tokenPair.AccessExpiry,
	))

	// 新しいリフレッシュトークンをCookieに設定
	response.SetCookie(c, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Expires:  tokenPair.RefreshExpiry,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return response.HandleSuccessWithData(c, http.StatusOK, "トークンを更新しました", &model.TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(time.Until(tokenPair.AccessExpiry).Seconds()),
		ExpiresAt:    tokenPair.AccessExpiry,
		RefreshToken: tokenPair.RefreshToken,
	})
}
