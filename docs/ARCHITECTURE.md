# 認証システムのアーキテクチャ改善案

## 現状の課題

### 1. 認証フローの問題

- トークン検証時の無限ループ
- トークン更新ロジックの複雑さ
- 認証状態管理の重複

### 2. セキュリティの課題

- トークンの有効期限管理が不明確
- リフレッシュトークンの仕組みが未実装
- CSRF トークンと JWT の関係が不明確

## 改善案

### 1. 認証フローの改善

#### フロントエンド（React）

```typescript
// 改善後のaxiosインターセプター
axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (!error.response || error.response.status !== 401) {
      return Promise.reject(error);
    }

    const originalRequest = error.config;
    if (originalRequest._retry) {
      // すでにリトライ済みの場合は、ログアウト処理を実行
      localStorage.removeItem("accessToken");
      window.dispatchEvent(new CustomEvent("unauthorized"));
      return Promise.reject(error);
    }

    try {
      originalRequest._retry = true;
      const response = await axiosInstance.post("/refresh-token");
      const { accessToken } = response.data;

      if (accessToken) {
        localStorage.setItem("accessToken", accessToken);
        axiosInstance.defaults.headers.common[
          "Authorization"
        ] = `Bearer ${accessToken}`;
        originalRequest.headers["Authorization"] = `Bearer ${accessToken}`;
        return axiosInstance(originalRequest);
      }
    } catch (refreshError) {
      localStorage.removeItem("accessToken");
      window.dispatchEvent(new CustomEvent("unauthorized"));
      return Promise.reject(refreshError);
    }
  }
);
```

#### バックエンド（Go）

```go
// トークン管理の改善
type TokenPair struct {
    AccessToken    string
    RefreshToken   string
    AccessExpiry   time.Time
    RefreshExpiry  time.Time
}

// トークン生成の改善
func (uc *userController) GenerateTokenPair(user model.User) (*TokenPair, error) {
    accessToken, accessExpiry, err := uc.generateAccessToken(user)
    if err != nil {
        return nil, err
    }

    refreshToken, refreshExpiry, err := uc.generateRefreshToken(user)
    if err != nil {
        return nil, err
    }

    return &TokenPair{
        AccessToken:    accessToken,
        RefreshToken:   refreshToken,
        AccessExpiry:   accessExpiry,
        RefreshExpiry:  refreshExpiry,
    }, nil
}
```

### 2. アーキテクチャの改善点

#### トークン管理の一元化

- アクセストークン（短期）とリフレッシュトークン（長期）の分離
- トークンの状態管理をカスタムフックに集約
- トークンの更新ロジックの明確化

#### エラーハンドリングの改善

- エラー種別の明確な定義
- 一貫性のあるエラーレスポンス
- ユーザーフレンドリーなエラーメッセージ

#### セキュリティの強化

- トークンのローテーション
- リフレッシュトークンの再利用検知
- 適切なトークン有効期限の設定

### 3. 実装手順

1. バックエンドの改善

   - リフレッシュトークンエンドポイントの実装
   - トークン検証ロジックの改善
   - エラーハンドリングの統一

2. フロントエンドの改善

   - 認証状態管理の一元化
   - トークン更新ロジックの実装
   - エラーハンドリングの改善

3. セキュリティ強化
   - CSRF トークンの適切な管理
   - トークンのセキュアな保存
   - 適切な CORS 設定

### 4. 新しい認証フロー

1. ログイン時

   - ユーザー認証
   - アクセストークンとリフレッシュトークンの発行
   - セキュアなトークン保存

2. API リクエスト時

   - アクセストークンの検証
   - 必要に応じてリフレッシュトークンによる更新
   - エラー時の適切なハンドリング

3. ログアウト時
   - トークンの無効化
   - セッションのクリーンアップ
   - 状態のリセット

## 推奨設定

### トークンの有効期限

- アクセストークン: 15 分
- リフレッシュトークン: 7 日

### セキュリティヘッダー

```go
e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
    XSSProtection:         "1; mode=block",
    ContentTypeNosniff:    "nosniff",
    XFrameOptions:         "SAMEORIGIN",
    HSTSMaxAge:           31536000,
    HSTSPreloadEnabled:    true,
}))
```

### CORS 設定の最適化

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000", os.Getenv("FRONTEND_URL")},
    AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
    ExposeHeaders:    []string{"X-CSRF-Token"},
    AllowCredentials: true,
    MaxAge:           3600,
}))
```

## 移行計画

1. バックエンド更新

   - リフレッシュトークンの実装
   - トークン検証の改善
   - エラーハンドリングの統一

2. フロントエンド更新

   - 認証フックの改善
   - トークン管理の一元化
   - エラーハンドリングの実装

3. テストとデプロイ
   - 単体テストの追加
   - E2E テストの実装
   - 段階的なデプロイ
