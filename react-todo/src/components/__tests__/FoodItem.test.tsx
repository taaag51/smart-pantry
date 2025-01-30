import { render, screen } from '@testing-library/react'
import { FoodItem } from '../FoodItem'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
})

const mockFoodItem = {
  id: 1,
  title: 'トマト',
  quantity: 3,
  expiryDate: new Date('2025-02-05'), // 1週間後
}

const renderWithProvider = (component: React.ReactElement) => {
  return render(
    <QueryClientProvider client={queryClient}>{component}</QueryClientProvider>
  )
}

describe('FoodItem', () => {
  it('食材の基本情報を表示', () => {
    renderWithProvider(<FoodItem foodItem={mockFoodItem} />)
    expect(screen.getByText('トマト')).toBeInTheDocument()
    expect(screen.getByText('数量: 3')).toBeInTheDocument()
    expect(screen.getByText(/2025年2月5日/)).toBeInTheDocument()
  })

  it('期限切れ間近の食材は警告表示', () => {
    const warningItem = {
      ...mockFoodItem,
      expiryDate: new Date(Date.now() + 3 * 24 * 60 * 60 * 1000), // 3日後
    }
    renderWithProvider(<FoodItem foodItem={warningItem} />)
    expect(screen.getByText(/期限切れまで3日/)).toBeInTheDocument()
  })

  it('期限切れの食材は赤色表示', () => {
    const expiredItem = {
      ...mockFoodItem,
      expiryDate: new Date(Date.now() - 24 * 60 * 60 * 1000), // 1日前
    }
    renderWithProvider(<FoodItem foodItem={expiredItem} />)
    const titleElement = screen.getByText('トマト')
    expect(titleElement).toHaveClass('text-red-600')
  })
})
