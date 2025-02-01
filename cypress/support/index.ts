/// <reference types="cypress" />

// Import commands.js using ES2015 syntax:
import "./commands";

// グローバル設定
beforeEach(() => {
  // テスト実行前にローカルストレージをクリア
  cy.clearLocalStorage();
});

// エラーハンドリング
Cypress.on("uncaught:exception", (err, runnable) => {
  // テスト実行を継続
  return false;
});

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * ログインを実行するカスタムコマンド
       * @example cy.login('testuser@example.com', 'password123')
       */
      login(email: string, password: string): Chainable<Element>;
    }
  }
}

export {};
