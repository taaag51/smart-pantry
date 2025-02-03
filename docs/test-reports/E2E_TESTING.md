# E2E テスト実装ガイド

## 概要

このドキュメントでは、Smart Pantry アプリケーションのエンドツーエンド（E2E）テストの実装と実行方法について説明します。

## 技術スタック

- Cypress: E2E テストフレームワーク
- TypeScript: テストコード記述言語
- React: フロントエンドフレームワーク
- Go: バックエンド API

## 環境構築

1. 必要なパッケージのインストール

```bash
yarn add -D cypress @types/cypress
```

2. Cypress の設定

- baseUrl: http://localhost:3000 (React アプリケーション)
- apiUrl: http://localhost:8080 (Go バックエンド API)
- スクリーンショット保存先: cypress/screenshots
- 動画保存先: cypress/videos

## テストの実行方法

### 開発時のテスト実行

1. バックエンド API の起動

```bash
cd backend-api
go run main.go
```

2. React アプリケーションの起動

```bash
cd react-app
yarn start
```

3. Cypress の起動

```bash
# GUIモード
yarn cypress open

# ヘッドレスモード
yarn cypress run
```

### CI/CD 環境でのテスト実行

GitHub Actions で自動的に E2E テストが実行されます（.github/workflows/e2e-test.yml 参照）

## テストケースの作成方法

### 1. テストファイルの配置

- テストファイルは `cypress/e2e/` ディレクトリに配置
- ファイル名は機能名に応じて `*.cy.ts` の形式で作成
  例: `auth.cy.ts`, `pantry.cy.ts`

### 2. テストケースの基本構造

```typescript
describe("テストスイート名", () => {
  beforeEach(() => {
    // テスト前の共通処理
    cy.clearLocalStorage();
    cy.visit("/");
  });

  it("テストケース名", () => {
    // テストの実装
  });
});
```

### 3. よく使用する Cypress コマンド

- ページ遷移: `cy.visit("/")`
- 要素の取得: `cy.get("セレクタ")`
- 要素のクリック: `cy.get("セレクタ").click()`
- テキスト入力: `cy.get("セレクタ").type("テキスト")`
- 要素の存在確認: `cy.get("セレクタ").should("exist")`
- URL の確認: `cy.url().should("include", "パス")`
- スクリーンショット: `cy.screenshot("ファイル名")`

### 4. カスタムコマンドの作成

共通で使用する操作はカスタムコマンドとして実装することを推奨します。

```typescript
// cypress/support/commands.ts
Cypress.Commands.add("login", (email: string, password: string) => {
  cy.get('input[name="email"]').type(email);
  cy.get('input[name="password"]').type(password);
  cy.get('button[type="submit"]').click();
});
```

## ベストプラクティス

1. テストの独立性

   - 各テストは他のテストに依存しないように実装
   - beforeEach でテスト環境をクリーンアップ

2. データ管理

   - テストデータは固定値を使用
   - 環境変数（cypress.config.js）で管理することを推奨

3. エラーハンドリング

   - 適切なアサーションを使用
   - タイムアウトの設定
   - リトライメカニズムの活用

4. 保守性
   - セレクタは data-testid 属性の使用を推奨
   - 共通処理はカスタムコマンドとして実装
   - 適切なコメントとドキュメンテーション

## トラブルシューティング

1. テストが不安定な場合

   - タイムアウトの設定を確認（cypress.config.js の defaultCommandTimeout）
   - 要素の待機処理を追加（cy.wait()の適切な使用）
   - ネットワークの状態を確認

2. テストが失敗する場合

   - スクリーンショットとビデオ録画を確認
   - エラーメッセージを詳細に確認
   - 環境変数の設定を確認

3. パフォーマンスの問題
   - 不要な wait を削除
   - テストの並列実行を検討
   - ビデオ録画を必要な場合のみ有効化

## Smart Pantry アプリケーション固有の注意点

1. 認証関連

   - テストユーザーのセットアップ
   - トークンの適切な管理
   - セッション状態のクリーンアップ

2. データ操作

   - 食品アイテムの追加/削除のテスト
   - レシピ提案機能のテスト
   - データベースの初期化

3. UI コンポーネント
   - レスポンシブデザインのテスト
   - モーダル/ポップアップの操作
   - フォーム入力の検証

## テスト結果の確認

1. レポート

   - テスト実行結果は cypress/reports に保存
   - スクリーンショットは cypress/screenshots に保存
   - ビデオ録画は cypress/videos に保存

2. CI/CD 連携
   - GitHub Actions での自動テスト
   - テスト結果のアーティファクト保存
   - ステータスバッジの表示

## 終わりに

このガイドは、Smart Pantry アプリケーションの E2E テスト実装の基本的な方針と手順をまとめたものです。
プロジェクトの成長に合わせて、適宜更新・拡張していくことを推奨します。
