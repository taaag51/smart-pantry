import { render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RecipeSuggestions } from '../RecipeSuggestions'
import { useQueryRecipe } from '../../hooks/useQueryRecipe'

jest.mock('../../hooks/useQueryRecipe')
const mockUseQueryRecipe = useQueryRecipe as jest.MockedFunction<
  typeof useQueryRecipe
>

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

const renderWithProvider = (component: React.ReactElement) => {
  return render(
    <QueryClientProvider client={queryClient}>{component}</QueryClientProvider>
  )
}

describe('RecipeSuggestions', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('ローディング中はローディングインジケータを表示', () => {
    mockUseQueryRecipe.mockReturnValue({
      data: undefined,
      isLoading: true,
      error: null,
      isError: false,
    } as any)

    renderWithProvider(<RecipeSuggestions />)
    expect(screen.getByRole('progressbar')).toBeInTheDocument()
  })

  it('エラー時はエラーメッセージを表示', () => {
    mockUseQueryRecipe.mockReturnValue({
      data: undefined,
      isLoading: false,
      error: new Error('エラーが発生しました'),
      isError: true,
    } as any)

    renderWithProvider(<RecipeSuggestions />)
    expect(screen.getByText(/レシピの取得に失敗しました/)).toBeInTheDocument()
  })

  it('データがない場合は適切なメッセージを表示', () => {
    mockUseQueryRecipe.mockReturnValue({
      data: undefined,
      isLoading: false,
      error: null,
      isError: false,
    } as any)

    renderWithProvider(<RecipeSuggestions />)
    expect(
      screen.getByText(/期限切れ間近の食材がありません/)
    ).toBeInTheDocument()
  })

  it('レシピデータを正しく表示', async () => {
    const mockRecipe = 'テストレシピの内容'
    mockUseQueryRecipe.mockReturnValue({
      data: mockRecipe,
      isLoading: false,
      error: null,
      isError: false,
    } as any)

    renderWithProvider(<RecipeSuggestions />)

    await waitFor(() => {
      expect(screen.getByText('おすすめレシピ')).toBeInTheDocument()
      expect(screen.getByText(mockRecipe)).toBeInTheDocument()
    })
  })
})
