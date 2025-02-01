// ***********************************************************
// This example support/index.js is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
// ***********************************************************

// Import commands.js using ES2015 syntax:
import './commands'

// グローバル設定
beforeEach(() => {
  // テスト実行前にローカルストレージをクリア
  cy.clearLocalStorage()
})

// エラーハンドリング
Cypress.on('uncaught:exception', (err, runnable) => {
  // テスト実行を継続
  return false
})

// カスタムコマンドの型定義
declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * ログインを実行するカスタムコマンド
       * @example cy.login('testuser@example.com', 'password123')
       */
      login(email: string, password: string): Chainable<Element>
    }
  }
}