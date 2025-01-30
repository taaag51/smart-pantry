import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { FoodList } from '../FoodList'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useQueryFoodItems } from '../../hooks/useQueryFoodItems'
import { useMutateFoodItem } from '../../hooks/useMutateFoodItem'

jest.mock('../../hooks/useQueryFoodItems')
jest.mock('../../hooks/useMutateFoodItem')

const mockUseQueryFoodItems = useQueryFoodItems as jest.MockedFunction<
  typeof useQueryFoodItems
>
const mockUseMutateFoodItem = useMutateFoodItem as jest.MockedFunction<
  typeof useMutateFoodItem
>

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

const mockFoodItems = [
  {
    id: 1,
    title: 'トマト',
    quantity: 3,
    expiryDate: new Date('2025-02-05'),
  },
  {
    id: 2,
    title: 'きゅうり',
    quantity: 2,
    expiryDate: new Date('2025-02-10'),
  },
]

const renderWithProvider = (component: React.ReactElement) => {
  return render(
    <QueryClientProvider client={queryClient}>{component}</QueryClientProvider>
  )
}

describe('FoodList', () => {
  beforeEach(() => {
    mockUseQueryFoodItems.mockReturnValue({
      data: mockFoodItems,
      isLoading: false,
      error: null,
    } as any)

    mockUseMutateFoodItem.mockReturnValue({
      createFoodItemMutation: {
        mutate: jest.fn(),
      },
    } as any)
  })

  it('食材一覧を表示', () => {
    renderWithProvider(<FoodList />)
    expect(screen.getByText('トマト')).toBeInTheDocument()
    expect(screen.getByText('きゅうり')).toBeInTheDocument()
  })

  it('新しい食材を追加できる', async () => {
    const mockCreateMutate = jest.fn()
    mockUseMutateFoodItem.mockReturnValue({
      createFoodItemMutation: {
        mutate: mockCreateMutate,
      },
    } as any)

    renderWithProvider(<FoodList />)

    // フォームに入力
    fireEvent.change(screen.getByLabelText('食材名'), {
      target: { value: 'なす' },
    })
    fireEvent.change(screen.getByLabelText('数量'), {
      target: { value: '4' },
    })
    const today = new Date()
    fireEvent.change(screen.getByLabelText('賞味期限'), {
      target: { value: today.toISOString().split('T')[0] },
    })

    // フォームを送信
    fireEvent.submit(screen.getByRole('button', { name: '追加' }))

    await waitFor(() => {
      expect(mockCreateMutate).toHaveBeenCalledWith({
        title: 'なす',
        quantity: 4,
        expiryDate: expect.any(Date),
      })
    })
  })

  it('ローディング中はローディングインジケータを表示', () => {
    mockUseQueryFoodItems.mockReturnValue({
      data: undefined,
      isLoading: true,
      error: null,
    } as any)

    renderWithProvider(<FoodList />)
    expect(screen.getByRole('progressbar')).toBeInTheDocument()
  })

  it('エラー時はエラーメッセージを表示', () => {
    mockUseQueryFoodItems.mockReturnValue({
      data: undefined,
      isLoading: false,
      error: new Error('エラーが発生しました'),
    } as any)

    renderWithProvider(<FoodList />)
    expect(screen.getByText(/エラーが発生しました/)).toBeInTheDocument()
  })
})
