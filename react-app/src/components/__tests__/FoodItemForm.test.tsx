import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { FoodItemForm } from '../FoodItemForm'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'

// MSWサーバーのセットアップ
const server = setupServer(
  http.post('/api/food-items', () => {
    return HttpResponse.json(
      {
        id: 1,
        title: 'テスト食材',
        quantity: 1,
        expiryDate: '2025-02-05',
      },
      { status: 201 }
    )
  })
)

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
    mutations: {
      retry: false,
    },
  },
})

const renderWithProvider = (component: React.ReactElement) => {
  return render(
    <QueryClientProvider client={queryClient}>{component}</QueryClientProvider>
  )
}

describe('FoodItemForm', () => {
  it('必須フィールドのバリデーション', async () => {
    renderWithProvider(<FoodItemForm />)

    // 空のフォームを送信
    const submitButton = screen.getByRole('button', { name: '食材を追加' })
    fireEvent.click(submitButton)

    // エラーメッセージの確認
    await waitFor(() => {
      expect(screen.getByText('食材名は必須です')).toBeInTheDocument()
      expect(screen.getByText('数量は必須です')).toBeInTheDocument()
      expect(screen.getByText('賞味期限は必須です')).toBeInTheDocument()
    })
  })

  it('数量の最小値/最大値チェック', async () => {
    renderWithProvider(<FoodItemForm />)

    const quantityInput = screen.getByLabelText('数量入力')

    // 最小値未満
    await userEvent.type(quantityInput, '-1')
    expect(
      screen.getByText('数量は1以上である必要があります')
    ).toBeInTheDocument()

    // 最大値超過
    await userEvent.clear(quantityInput)
    await userEvent.type(quantityInput, '1001')
    expect(
      screen.getByText('数量は1000以下である必要があります')
    ).toBeInTheDocument()
  })

  it('賞味期限の日付形式チェック', async () => {
    renderWithProvider(<FoodItemForm />)

    const dateInput = screen.getByLabelText('賞味期限選択')

    // 過去の日付
    const pastDate = new Date()
    pastDate.setDate(pastDate.getDate() - 1)
    await userEvent.type(dateInput, pastDate.toISOString().split('T')[0])
    expect(
      screen.getByText('賞味期限は今日以降である必要があります')
    ).toBeInTheDocument()
  })

  it('正常なフォーム送信', async () => {
    renderWithProvider(<FoodItemForm />)

    // フォームに有効な値を入力
    await userEvent.type(screen.getByLabelText('食材名入力'), 'テスト食材')
    await userEvent.type(screen.getByLabelText('数量入力'), '1')

    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    await userEvent.type(
      screen.getByLabelText('賞味期限選択'),
      tomorrow.toISOString().split('T')[0]
    )

    // フォーム送信
    fireEvent.click(screen.getByRole('button', { name: '食材を追加' }))

    // 送信成功の確認
    await waitFor(() => {
      expect(screen.getByText('食材を追加しました')).toBeInTheDocument()
    })
  })

  it('APIエラー時の処理', async () => {
    // APIエラーをモック
    server.use(
      http.post('/api/food-items', () => {
        return HttpResponse.json(
          { message: 'Internal Server Error' },
          { status: 500 }
        )
      })
    )

    renderWithProvider(<FoodItemForm />)

    // フォームに有効な値を入力
    await userEvent.type(screen.getByLabelText('食材名入力'), 'テスト食材')
    await userEvent.type(screen.getByLabelText('数量入力'), '1')

    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    await userEvent.type(
      screen.getByLabelText('賞味期限選択'),
      tomorrow.toISOString().split('T')[0]
    )

    // フォーム送信
    fireEvent.click(screen.getByRole('button', { name: '食材を追加' }))

    // エラーメッセージの表示を確認
    await waitFor(() => {
      expect(screen.getByText(/エラーが発生しました/)).toBeInTheDocument()
    })
  })

  it('送信中の状態表示', async () => {
    // API遅延をシミュレート
    server.use(
      http.post('/api/food-items', async () => {
        await new Promise((resolve) => setTimeout(resolve, 100))
        return HttpResponse.json(
          {
            id: 1,
            title: 'テスト食材',
            quantity: 1,
            expiryDate: '2025-02-05',
          },
          { status: 201 }
        )
      })
    )

    renderWithProvider(<FoodItemForm />)

    // フォームに有効な値を入力
    await userEvent.type(screen.getByLabelText('食材名入力'), 'テスト食材')
    await userEvent.type(screen.getByLabelText('数量入力'), '1')

    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    await userEvent.type(
      screen.getByLabelText('賞味期限選択'),
      tomorrow.toISOString().split('T')[0]
    )

    // フォーム送信
    fireEvent.click(screen.getByRole('button', { name: '食材を追加' }))

    // 送信中の状態を確認
    expect(screen.getByRole('button', { name: '追加中...' })).toBeDisabled()

    // 送信完了後の状態を確認
    await waitFor(() => {
      expect(screen.getByRole('button', { name: '食材を追加' })).toBeEnabled()
    })
  })
})
