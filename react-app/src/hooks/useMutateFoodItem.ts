import { AxiosError, AxiosResponse } from 'axios'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { FoodItem } from '../types'
import { useError } from './useError'
import axiosInstance, { getCsrfToken } from '../lib/axios'

/**
 * APIレスポンスの型定義
 */
interface ApiResponse {
  data: FoodItem
}

/**
 * APIエラーの型定義
 */
interface ApiError {
  message: string
}

/**
 * キャッシュ更新用のユーティリティ関数
 */
const updateFoodItemsCache = (
  queryClient: ReturnType<typeof useQueryClient>,
  updater: (prev: FoodItem[]) => FoodItem[]
) => {
  const previousFoodItems = queryClient.getQueryData<FoodItem[]>(['foodItems'])
  if (previousFoodItems) {
    queryClient.setQueryData(['foodItems'], updater(previousFoodItems))
  }
}

/**
 * エラーハンドリング用のユーティリティ関数
 */
const handleApiError = (
  err: AxiosError<ApiError>,
  defaultMessage: string,
  errorHandler: (message: string) => void
) => {
  console.error('API Error:', {
    response: err.response?.data,
    status: err.response?.status,
    headers: err.response?.headers,
  })
  const message = err.response?.data.message || defaultMessage
  errorHandler(message)
}

/**
 * 食材データの操作（作成・更新・削除）を管理するカスタムフック
 *
 * このフックは以下の機能を提供します：
 * - 食材の作成（createFoodItemMutation）
 * - 食材の更新（updateFoodItemMutation）
 * - 食材の削除（deleteFoodItemMutation）
 *
 * 各操作は自動的にキャッシュを更新し、エラーハンドリングを行います。
 */
export const useMutateFoodItem = () => {
  const queryClient = useQueryClient()
  const { switchErrorHandling } = useError()

  /**
   * 食材作成のミューテーション
   */
  const createFoodItemMutation = useMutation({
    mutationFn: async (foodItem: Omit<FoodItem, 'id'>) => {
      console.log('Creating food item with:', foodItem)
      await getCsrfToken()
      const headers = axiosInstance.defaults.headers
      console.log('Request headers:', headers)
      return await axiosInstance.post<ApiResponse>('/food-items', foodItem)
    },
    onSuccess: (res: AxiosResponse<ApiResponse>) => {
      console.log('Creation successful:', res.data)
      updateFoodItemsCache(queryClient, (prev) => [...prev, res.data.data])
    },
    onError: (err: AxiosError<ApiError>) => {
      handleApiError(err, '食材の作成に失敗しました', switchErrorHandling)
    },
  })

  /**
   * 食材更新のミューテーション
   */
  const updateFoodItemMutation = useMutation({
    mutationFn: async (foodItem: FoodItem) => {
      await getCsrfToken()
      return await axiosInstance.put<ApiResponse>(
        `/food-items/${foodItem.id}`,
        foodItem
      )
    },
    onSuccess: (res: AxiosResponse<ApiResponse>, variables: FoodItem) => {
      updateFoodItemsCache(queryClient, (prev) =>
        prev.map((item) => (item.id === variables.id ? res.data.data : item))
      )
    },
    onError: (err: AxiosError<ApiError>) => {
      handleApiError(err, '食材の更新に失敗しました', switchErrorHandling)
    },
  })

  /**
   * 食材削除のミューテーション
   */
  const deleteFoodItemMutation = useMutation({
    mutationFn: async (id: number) => {
      await getCsrfToken()
      return await axiosInstance.delete(`/food-items/${id}`)
    },
    onSuccess: (_, variables: number) => {
      updateFoodItemsCache(queryClient, (prev) =>
        prev.filter((item) => item.id !== variables)
      )
    },
    onError: (err: AxiosError<ApiError>) => {
      handleApiError(err, '食材の削除に失敗しました', switchErrorHandling)
    },
  })

  return {
    createFoodItemMutation,
    updateFoodItemMutation,
    deleteFoodItemMutation,
  }
}
