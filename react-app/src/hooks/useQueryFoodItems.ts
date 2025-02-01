import { AxiosError } from 'axios'
import { useQuery } from '@tanstack/react-query'
import { FoodItem } from '../types'
import { useError } from './useError'
import axiosInstance from '../lib/axios'

export const useQueryFoodItems = () => {
  const { switchErrorHandling } = useError()

  const getFoodItems = async (): Promise<FoodItem[]> => {
    try {
      const { data } = await axiosInstance.get<FoodItem[]>('/food-items')
      if (!data) return []
      if (!Array.isArray(data)) {
        console.error('Received non-array data:', data)
        return []
      }
      return data
    } catch (err) {
      console.error('Failed to fetch food items:', err)
      if (err instanceof AxiosError) {
        if (err.response?.status === 401) {
          window.dispatchEvent(new CustomEvent('unauthorized'))
          return []
        }
        if (err.response?.data?.message) {
          switchErrorHandling(err.response.data.message)
        } else {
          switchErrorHandling('食材の取得に失敗しました')
        }
      }
      return []
    }
  }

  return useQuery({
    queryKey: ['foodItems'],
    queryFn: getFoodItems,
    staleTime: 1000 * 60 * 5, // 5分
    initialData: [],
    retry: 1,
    refetchOnWindowFocus: true,
    refetchOnMount: true,
    refetchOnReconnect: true,
    // 認証状態の変更時にデータを再取得
    enabled: true,
    refetchInterval: 0,
  })
}
