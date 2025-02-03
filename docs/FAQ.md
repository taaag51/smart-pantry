# スマートパントリーシステム - FAQ

## 目次

1. [一般的な質問](#一般的な質問)
2. [技術的な質問（Go）](#技術的な質問go)
3. [技術的な質問（React）](#技術的な質問react)
4. [開発環境のセットアップ](#開発環境のセットアップ)
5. [トラブルシューティング](#トラブルシューティング)

## 一般的な質問

### Q: このアプリケーションは何ができますか？

A: スマートパントリーシステムは以下の機能を提供します：

- 食材の登録・管理（名前、数量、賞味期限）
- 賞味期限の管理と通知
- レシピの提案
- ユーザー認証（登録・ログイン）

### Q: 開発環境のセットアップに必要なものは何ですか？

A: 以下のツールが必要です：

- Go (1.20 以上)
- Node.js (18.0 以上)
- Docker Desktop
- PostgreSQL (Docker で提供)
- Visual Studio Code（推奨）

### Q: ローカルで動かすにはどうすればいいですか？

A: 以下の手順で起動できます：

1. バックエンド:

```bash
cd backend-api
docker-compose up -d  # DBの起動
go run main.go
```

2. フロントエンド:

```bash
cd react-app
npm install
npm start
```

## 技術的な質問（Go）

### Q: プロジェクトの構造はどうなっていますか？

A: クリーンアーキテクチャに基づいて、以下のレイヤーで構成されています：

- Controller: HTTP リクエストの受け付け
- Usecase: ビジネスロジックの実装
- Repository: データアクセス処理
- Model: ドメインモデル

### Q: エラーハンドリングはどのように実装されていますか？

A: 以下の方針でエラーハンドリングを実装しています：

- 各レイヤーで適切なエラー型を定義
- エラーの種類に応じた HTTP ステータスコードの返却
- クライアントへの適切なエラーメッセージの提供

### Q: データベースの操作はどのように行われていますか？

A: GORM を使用してデータベース操作を行っています：

```go
// 例：食材の取得
func (fr *foodItemRepository) GetAllFoodItems(foodItems *[]model.FoodItem) error {
    return fr.db.Find(foodItems).Error
}
```

### Q: 認証はどのように実装されていますか？

A: JWT を使用した認証を実装しています：

1. ログイン時に JWT トークンを発行
2. リクエスト時に Cookie からトークンを検証
3. ミドルウェアでの認証チェック

## 技術的な質問（React）

### Q: 状態管理はどのように行っていますか？

A: カスタムフックを使用して状態管理を実装しています：

```typescript
// 例：食材データの取得と更新
const { data, isLoading, error } = useQueryFoodItems();
const { createFoodItemMutation } = useMutateFoodItem();
```

### Q: コンポーネントの設計方針は？

A: 以下の方針でコンポーネントを設計しています：

- 機能ごとにコンポーネントを分割
- ロジックはカスタムフックに分離
- Material-UI を使用した UI 実装

### Q: API との通信はどのように実装されていますか？

A: axios と React Query を使用しています：

```typescript
// 例：食材データの取得
const { data } = useQuery("foodItems", async () => {
  const response = await axios.get("/api/food-items");
  return response.data;
});
```

### Q: CSS のスタイリングはどのように行っていますか？

A: Material-UI のスタイリングシステムを使用しています：

```typescript
// 例：スタイルの適用
<Box sx={{
    display: 'flex',
    justifyContent: 'space-between',
    p: 2
}}>
```

## 開発環境のセットアップ

### Q: Go の開発環境をセットアップするには？

A: 以下の手順で設定できます：

1. Go のインストール
2. 必要なパッケージのインストール：

```bash
go mod download
```

3. 環境変数の設定：

```bash
cp .env.example .env
# .envファイルを編集
```

### Q: React の開発環境をセットアップするには？

A: 以下の手順で設定できます：

1. Node.js のインストール
2. 依存パッケージのインストール：

```bash
npm install
```

3. 開発サーバーの起動：

```bash
npm start
```

## トラブルシューティング

### Q: `go run main.go` でエラーが発生する場合は？

A: 以下を確認してください：

1. Go のバージョンが 1.20 以上か
2. 必要なパッケージがすべてインストールされているか
3. .env ファイルが正しく設定されているか
4. PostgreSQL が起動しているか

### Q: `npm start` でエラーが発生する場合は？

A: 以下を確認してください：

1. Node.js のバージョンが 18.0 以上か
2. node_modules が正しくインストールされているか
3. package.json に必要な依存関係が含まれているか

### Q: API リクエストが失敗する場合は？

A: 以下を確認してください：

1. バックエンドサーバーが起動しているか
2. 正しい URL にリクエストを送っているか
3. JWT トークンが正しく設定されているか
4. CORS の設定が正しいか

### Q: テストが失敗する場合は？

A: 以下を確認してください：

1. テスト用のデータベースが正しく設定されているか
2. 必要なモックが正しく設定されているか
3. テスト環境の環境変数が正しく設定されているか
