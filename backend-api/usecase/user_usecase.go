package usecase

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/taaag51/smart-pantry/backend-api/model"
	"github.com/taaag51/smart-pantry/backend-api/repository"
	"github.com/taaag51/smart-pantry/backend-api/validator"
)

var secretKey = []byte(os.Getenv("SECRET"))

type IUserUsecase interface {
	SignUp(user model.User) (model.User, error)
	Login(user model.User) (*model.TokenPair, error)
	VerifyToken(tokenString string) (*jwt.Token, error)
	RefreshTokens(refreshToken string) (*model.TokenPair, error)
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	if len(secretKey) == 0 {
		panic("SECRET environment variable is not set")
	}
	return &userUsecase{ur: ur, uv: uv}
}

func (uu *userUsecase) SignUp(user model.User) (model.User, error) {
	if err := uu.uv.ValidateUser(user); err != nil {
		return model.User{}, err
	}
	return uu.ur.CreateUser(&user)
}

func (uu *userUsecase) Login(user model.User) (*model.TokenPair, error) {
	if err := uu.uv.ValidateLogin(user); err != nil {
		return nil, err
	}

	storedUser, err := uu.ur.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	// パスワード検証は省略（既存の実装を使用）

	// アクセストークンの生成（15分有効）
	accessToken, accessExpiry, err := generateToken(storedUser.Email, "access", 15)
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンの生成（7日間有効）
	refreshToken, refreshExpiry, err := generateToken(storedUser.Email, "refresh", 24*7)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}, nil
}

func (uu *userUsecase) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// トークンの有効期限を確認
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	if time.Now().Unix() > int64(exp) {
		return nil, jwt.ErrTokenExpired
	}

	return token, nil
}

func (uu *userUsecase) RefreshTokens(refreshToken string) (*model.TokenPair, error) {
	// リフレッシュトークンの検証
	token, err := uu.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	// トークンタイプの確認
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, jwt.ErrInvalidKey
	}

	// メールアドレスの取得
	email, ok := claims["email"].(string)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	// ユーザーの存在確認
	user, err := uu.ur.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// 新しいトークンペアの生成
	return uu.Login(user)
}

// トークン生成のヘルパー関数
func generateToken(email, tokenType string, expiryHours int) (string, time.Time, error) {
	expiryTime := time.Now().Add(time.Duration(expiryHours) * time.Hour)

	claims := jwt.MapClaims{
		"email": email,
		"type":  tokenType,
		"exp":   expiryTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiryTime, nil
}
