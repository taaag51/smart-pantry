import { AxiosError } from 'axios'
import { useQuery } from '@tanstack/react-query'
import { FoodItem } from '../types'
import { useError } from './useError'
import axiosInstance, { getCsrfToken } from '../lib/axios'

interface ApiError {
  message: string
}

export const useQueryFoodItems = () => {
  const { switchErrorHandling } = useError()

  const getFoodItems = async () => {
    try {
      await getCsrfToken()
      const { data } = await axiosInstance.get<FoodItem[]>('/food-items', {
        withCredentials: true,
      })
      return data
    } catch (err) {
      if (err instanceof AxiosError) {
        const axiosError = err as AxiosError<ApiError>
        if (axiosError.response?.data.message) {
          switchErrorHandling(axiosError.response.data.message)
        } else {
          switchErrorHandling('食材の取得に失敗しました')
        }
      }
      throw err
    }
  }

  return useQuery({
    queryKey: ['foodItems'],
    queryFn: getFoodItems,
  })
}
