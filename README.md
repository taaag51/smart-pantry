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

## セットアップ手順

### 1. 環境変数の設定

`.env`ファイルを作成し、以下の環境変数を設定：

```bash
cp .env.example .env
```

```env
# バックエンド設定
PORT=8080
POSTGRES_USER=postgres
POSTGRES_PW=postgres
POSTGRES_DB=smart_pantry
POSTGRES_PORT=5432
POSTGRES_HOST=localhost
API_DOMAIN=localhost
SECRET=your-secret-key

# Gemini API設定
GEMINI_API_KEY=your-gemini-api-key
```

### 2. Gemini API キーの取得

1. [Google Cloud Console](https://console.cloud.google.com/)にアクセス
2. プロジェクトを作成（または既存のプロジェクトを選択）
3. [API とサービス] > [認証情報]に移動
4. [認証情報を作成] > [API キー]を選択
5. 作成された API キーを`.env`の`GEMINI_API_KEY`にコピー

### 3. データベースのセットアップ

PostgreSQL を Docker で起動：

```bash
cd go-rest-api
docker-compose up -d
```

マイグレーションの実行：

```bash
go run migrate/migrate.go
```

### 4. バックエンドの起動

```bash
cd go-rest-api
go mod download
go run main.go
```

### 5. フロントエンドの起動

```bash
cd react-todo
npm install
npm start
```

## 開発環境

### バックエンド

- Go 1.23.5
- Echo Framework
- GORM
- PostgreSQL 15
- Google Gemini API

### フロントエンド

- React 18
- TypeScript 4.9
- Material-UI 5
- React Query
- date-fns

## データベース構造

### FoodItems テーブル

```sql
CREATE TABLE food_items (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
```

### Users テーブル

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## API エンドポイント

### 認証

- POST `/signup`: ユーザー登録
- POST `/login`: ログイン
- POST `/logout`: ログアウト
- GET `/csrf`: CSRF トークン取得

### 食材管理

- GET `/food-items`: 食材一覧の取得
- GET `/food-items/:id`: 特定の食材の取得
- POST `/food-items`: 新規食材の登録
- PUT `/food-items/:id`: 食材情報の更新
- DELETE `/food-items/:id`: 食材の削除

### レシピ提案

- GET `/recipes/suggestions`: AI によるレシピ提案の取得

## テスト実行

```bash
# バックエンドのテスト
cd go-rest-api
go test ./...

# フロントエンドのテスト
cd react-todo
npm test
```

## 注意事項

- 賞味期限が 7 日以内の食材に対してアラートが表示されます
- レシピ提案は期限切れ間近の食材を優先的に使用します
- API リクエストには JWT 認証が必要です
- Gemini API の利用には課金が発生する可能性があります
