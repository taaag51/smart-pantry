import { AxiosError, AxiosResponse } from 'axios'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { FoodItem } from '../types'
import { useError } from './useError'
import axiosInstance, { getCsrfToken } from '../lib/axios'

interface ApiError {
  message: string
}

export const useMutateFoodItem = () => {
  const queryClient = useQueryClient()
  const { switchErrorHandling } = useError()

  const createFoodItemMutation = useMutation({
    mutationFn: async (foodItem: Omit<FoodItem, 'id'>) => {
      await getCsrfToken()
      return await axiosInstance.post<{ data: FoodItem }>(
        '/food-items',
        foodItem
      )
    },
    onSuccess: (res: AxiosResponse<{ data: FoodItem }>) => {
      const previousFoodItems = queryClient.getQueryData<FoodItem[]>([
        'foodItems',
      ])
      if (previousFoodItems) {
        queryClient.setQueryData(
          ['foodItems'],
          [...previousFoodItems, res.data]
        )
      }
    },
    onError: (err: AxiosError<ApiError>) => {
      if (err.response?.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('食材の作成に失敗しました')
      }
    },
  })

  const updateFoodItemMutation = useMutation({
    mutationFn: async (foodItem: FoodItem) => {
      await getCsrfToken()
      return await axiosInstance.put<{ data: FoodItem }>(
        `/food-items/${foodItem.id}`,
        foodItem
      )
    },
    onSuccess: (
      res: AxiosResponse<{ data: FoodItem }>,
      variables: FoodItem
    ) => {
      const previousFoodItems = queryClient.getQueryData<FoodItem[]>([
        'foodItems',
      ])
      if (previousFoodItems) {
        queryClient.setQueryData(
          ['foodItems'],
          previousFoodItems.map((item) =>
            item.id === variables.id ? res.data : item
          )
        )
      }
    },
    onError: (err: AxiosError<ApiError>) => {
      if (err.response?.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('食材の更新に失敗しました')
      }
    },
  })

  const deleteFoodItemMutation = useMutation({
    mutationFn: async (id: number) => {
      await getCsrfToken()
      return await axiosInstance.delete(`/food-items/${id}`)
    },
    onSuccess: (_, variables: number) => {
      const previousFoodItems = queryClient.getQueryData<FoodItem[]>([
        'foodItems',
      ])
      if (previousFoodItems) {
        queryClient.setQueryData(
          ['foodItems'],
          previousFoodItems.filter((item) => item.id !== variables)
        )
      }
    },
    onError: (err: AxiosError<ApiError>) => {
      if (err.response?.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('食材の削除に失敗しました')
      }
    },
  })

  return {
    createFoodItemMutation,
    updateFoodItemMutation,
    deleteFoodItemMutation,
  }
}
