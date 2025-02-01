describe("認証機能のE2Eテスト", () => {
  beforeEach(() => {
    // テスト実行前にローカルストレージをクリア
    cy.clearLocalStorage();
    cy.visit("/");
  });

  it("ログインフローの検証", () => {
    // メールアドレスとパスワードの入力
    cy.get('input[name="email"]').type("tani@example.com");
    cy.get('input[name="password"]').type("password123");

    // ログインボタンのクリック
    cy.get('button[type="submit"]').click();

    // ログイン後の画面遷移を確認
    cy.url().should("include", "/pantry");

    // スクリーンショットを撮影
    cy.screenshot("login-success");
  });

  it("バリデーションエラーの検証", () => {
    // 空の状態でログインボタンをクリック
    cy.get('button[type="submit"]').should("be.disabled");

    // メールアドレスのみ入力
    cy.get('input[name="email"]').type("tani@example.com");
    cy.get('button[type="submit"]').should("be.disabled");

    // パスワードのみ入力
    cy.get('input[name="email"]').clear();
    cy.get('input[name="password"]').type("password123");
    cy.get('button[type="submit"]').should("be.disabled");

    // スクリーンショットを撮影
    cy.screenshot("validation-error");
  });

  it("新規登録→ログインフローの検証", () => {
    // 新規登録モードに切り替え
    cy.get("button").contains("SwapHoriz").click();

    // メールアドレスとパスワードの入力
    cy.get('input[name="email"]').type("new-user@example.com");
    cy.get('input[name="password"]').type("newpassword123");

    // 登録ボタンのクリック
    cy.get('button[type="submit"]').click();

    // 登録後の画面遷移を確認
    cy.url().should("include", "/pantry");

    // スクリーンショットを撮影
    cy.screenshot("registration-success");
  });

  it("エラーメッセージの表示確認", () => {
    // 不正なログイン情報で試行
    cy.get('input[name="email"]').type("invalid@example.com");
    cy.get('input[name="password"]').type("wrongpassword");
    cy.get('button[type="submit"]').click();

    // エラーメッセージの表示を確認
    cy.get(".MuiAlert-root").should("be.visible");

    // スクリーンショットを撮影
    cy.screenshot("error-message");
  });
});
