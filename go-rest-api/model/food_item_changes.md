# FoodItem モデルの変更点

## 既存の Task モデルからの変更点

### 共通のフィールド（変更なし）

- `ID` (uint): プライマリーキー
- `Title` (string): 名前（Task モデルの Title フィールドを流用）
- `CreatedAt` (time.Time): 作成日時
- `UpdatedAt` (time.Time): 更新日時
- `User` (User): ユーザーとの関連
- `UserId` (uint): ユーザー ID（外部キー）

### 新規追加フィールド

1. `Quantity` (int)

   - 食材の数量を管理
   - `gorm:"not null"` で必須項目として設定

2. `ExpiryDate` (time.Time)
   - 食材の賞味期限を管理
   - `gorm:"not null"` で必須項目として設定

## レスポンス構造体の変更点

### FoodItemResponse

- Task モデルの TaskResponse 構造体をベースに作成
- 新規フィールド（Quantity, ExpiryDate）を追加
- ユーザー情報は除外（TaskResponse と同様）

## マイグレーション

- `migrate.go` に `model.FoodItem{}` を追加
- 既存のマイグレーション機能を活用

## 注意点

- 既存の CRUD 処理を最大限活用
- ユーザーとの関連付けは Task モデルと同様の方式を採用
- 賞味期限による食材管理が可能な設計
