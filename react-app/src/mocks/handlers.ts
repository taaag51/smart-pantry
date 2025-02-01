import { http, HttpResponse } from 'msw'

export const handlers = [
  // 認証関連のエンドポイント
  http.post('/login', async () => {
    return HttpResponse.json(
      {
        message: 'ログインに成功しました',
        user: {
          id: 1,
          email: 'test@example.com',
          name: 'テストユーザー',
        },
      },
      { status: 200 }
    )
  }),

  http.post('/logout', () => {
    return HttpResponse.json({ message: 'ログアウトしました' }, { status: 200 })
  }),

  http.get('/csrf', () => {
    return HttpResponse.json({ csrf_token: 'test-csrf-token' }, { status: 200 })
  }),

  // 食材追加のエンドポイント
  http.post('/api/food-items', async () => {
    return HttpResponse.json(
      {
        id: 1,
        title: 'テスト食材',
        quantity: 1,
        expiryDate: '2025-02-05',
      },
      { status: 201 }
    )
  }),

  // 食材一覧取得のエンドポイント
  http.get('/api/food-items', () => {
    return HttpResponse.json([
      {
        id: 1,
        title: 'テスト食材1',
        quantity: 2,
        expiryDate: '2025-02-05',
      },
      {
        id: 2,
        title: 'テスト食材2',
        quantity: 3,
        expiryDate: '2025-02-10',
      },
    ])
  }),

  // レシピ生成のエンドポイント
  http.post('/api/recipes/generate', () => {
    return HttpResponse.json({
      recipe: `
        【レシピ名】
        テスト食材のテストレシピ

        【材料（2人分）】
        - テスト食材1 2個
        - テスト食材2 3個

        【作り方】
        1. テスト手順1
        2. テスト手順2
        3. テスト手順3

        【栄養バランス】
        バランスの良い一品です。
      `,
    })
  }),
]
