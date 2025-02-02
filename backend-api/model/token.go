package model

import "time"

// TokenPair は、アクセストークンとリフレッシュトークンのペアを表す
type TokenPair struct {
	AccessToken   string    `json:"accessToken"`
	RefreshToken  string    `json:"refreshToken"`
	AccessExpiry  time.Time `json:"accessExpiry"`
	RefreshExpiry time.Time `json:"refreshExpiry"`
}

// TokenClaims はJWTトークンのペイロードを表す
type TokenClaims struct {
	Email string `json:"email"`
	Type  string `json:"type"` // "access" または "refresh"
}

// RefreshTokenRequest はリフレッシュトークンのリクエストを表す
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// TokenResponse はトークンのレスポンスを表す
type TokenResponse struct {
	AccessToken  string    `json:"accessToken"`
	TokenType    string    `json:"tokenType"`
	ExpiresIn    int       `json:"expiresIn"` // 有効期限（秒）
	RefreshToken string    `json:"refreshToken,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
