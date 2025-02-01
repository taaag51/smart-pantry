# スマートパントリーシステム - アーキテクチャ設計

## システム構成図

```mermaid
graph TB
    subgraph "フロントエンド (React)"
        UI[UI Components]
        Hooks[Custom Hooks]
        Store[State Management]
    end

    subgraph "バックエンド (Go)"
        Router[Router]
        Controller[Controllers]
        Usecase[Usecases]
        Repository[Repositories]
        DB[(Database)]
    end

    UI --> Hooks
    Hooks --> Store
    Store --> |HTTP Requests| Router
    Router --> Controller
    Controller --> Usecase
    Usecase --> Repository
    Repository --> DB
```

## リクエストフロー

```mermaid
sequenceDiagram
    participant Browser as ブラウザ
    participant React as Reactアプリ
    participant Router as Echoルーター
    participant Controller as コントローラー
    participant Usecase as ユースケース
    participant Repository as リポジトリ
    participant DB as データベース

    Browser->>React: ユーザーアクション
    React->>Router: HTTPリクエスト
    Note over React,Router: JWT認証
    Router->>Controller: ルーティング
    Controller->>Usecase: ビジネスロジック実行
    Usecase->>Repository: データアクセス
    Repository->>DB: クエリ実行
    DB-->>Repository: データ
    Repository-->>Usecase: データ
    Usecase-->>Controller: 処理結果
    Controller-->>Router: レスポンス
    Router-->>React: HTTPレスポンス
    React-->>Browser: UI更新
```

## レイヤー構造

### 1. プレゼンテーション層 (React)

```mermaid
classDiagram
    class UIComponents {
        +FoodItem
        +FoodItemForm
        +RecipeSuggestions
        +Auth
    }

    class CustomHooks {
        +useFoodItemForm()
        +useAuth()
        +useMutateFoodItem()
    }

    class StateManagement {
        +Store
        +Actions
        +Mutations
    }

    UIComponents --> CustomHooks
    CustomHooks --> StateManagement
```

### 2. アプリケーション層 (Go)

```mermaid
classDiagram
    class Controllers {
        +FoodItemController
        +RecipeController
        +UserController
    }

    class Usecases {
        +FoodItemUsecase
        +RecipeUsecase
        +UserUsecase
    }

    class Repositories {
        +FoodItemRepository
        +UserRepository
    }

    Controllers --> Usecases
    Usecases --> Repositories
```

## セキュリティ設計

### 認証フロー

```mermaid
sequenceDiagram
    participant Client
    participant Auth
    participant JWT
    participant API

    Client->>Auth: ログインリクエスト
    Auth->>JWT: トークン生成
    JWT-->>Auth: JWTトークン
    Auth-->>Client: トークン返却
    Client->>API: APIリクエスト + JWT
    Note over API: トークン検証
    API-->>Client: レスポンス
```

## データモデル

### エンティティ関係図

```mermaid
erDiagram
    User ||--o{ FoodItem : "所有"
    FoodItem ||--o{ Recipe : "使用"
    User {
        uint id
        string email
        string password
    }
    FoodItem {
        uint id
        string title
        int quantity
        date expiry_date
        uint user_id
    }
    Recipe {
        uint id
        string title
        text description
        uint user_id
    }
```

## エラーハンドリング設計

```mermaid
flowchart TD
    A[クライアントエラー] --> B{エラー種別}
    B -->|バリデーションエラー| C[400 Bad Request]
    B -->|認証エラー| D[401 Unauthorized]
    B -->|権限エラー| E[403 Forbidden]
    B -->|未検出エラー| F[500 Internal Server Error]

    C --> G[エラーレスポンス]
    D --> G
    E --> G
    F --> G
```

## 開発環境構成

```mermaid
graph LR
    A[開発環境] --> B[フロントエンド]
    A --> C[バックエンド]
    B --> D[Node.js]
    B --> E[React]
    B --> F[TypeScript]
    C --> G[Go]
    C --> H[PostgreSQL]
    C --> I[Docker]
```

## テスト戦略

```mermaid
graph TD
    A[テスト戦略] --> B[ユニットテスト]
    A --> C[統合テスト]
    A --> D[E2Eテスト]

    B --> E[Go testing]
    B --> F[React Testing Library]
    C --> G[API テスト]
    D --> H[Cypress]
```

## デプロイメントフロー

```mermaid
graph LR
    A[Git Push] --> B[GitHub Actions]
    B --> C[テスト実行]
    C --> D{テスト結果}
    D -->|成功| E[ビルド]
    D -->|失敗| F[通知]
    E --> G[デプロイ]
```
