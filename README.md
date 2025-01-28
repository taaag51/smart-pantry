# Smart Pantry

食材管理とレシピ提案を行う Web アプリケーション

## 機能

1. 食材管理

   - 食材の登録・編集・削除
   - 数量管理
   - 賞味期限管理

2. 賞味期限トラッカー

   - 期限切れ間近の食材を自動検出
   - 視覚的なアラート表示（黄色：期限間近、赤：期限切れ）

3. AI レシピ提案
   - 期限切れ間近の食材を活用したレシピを自動提案
   - 栄養バランスを考慮したレシピ生成
   - 5 分ごとに自動更新

## 技術スタック

### バックエンド

- Go
- Echo Framework
- GORM
- Google Gemini API

### フロントエンド

- React
- TypeScript
- Material-UI
- React Query

## セットアップ

1. 環境変数の設定

```bash
cp .env.example .env
```

必要な環境変数：

- `GEMINI_API_KEY`: Google Gemini API のキー
- `SECRET`: JWT シークレットキー
- `API_DOMAIN`: API ドメイン

2. バックエンドの起動

```bash
cd go-rest-api
go mod download
go run migrate/migrate.go
go run main.go
```

3. フロントエンドの起動

```bash
cd react-todo
npm install
npm start
```

## テスト実行

```bash
# バックエンドのテスト
cd go-rest-api
go test ./...

# フロントエンドのテスト
cd react-todo
npm test
```

## API エンドポイント

### 食材管理

- GET `/food-items`: 食材一覧の取得
- GET `/food-items/:id`: 特定の食材の取得
- POST `/food-items`: 新規食材の登録
- PUT `/food-items/:id`: 食材情報の更新
- DELETE `/food-items/:id`: 食材の削除

### レシピ提案

- GET `/recipes/suggestions`: AI によるレシピ提案の取得

## 注意事項

- 賞味期限が 7 日以内の食材に対してアラートが表示されます
- レシピ提案は期限切れ間近の食材を優先的に使用します
- API リクエストには JWT 認証が必要です
