import axios from 'axios'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { FoodItem } from '../types'
import { useError } from './useError'

export const useMutateFoodItem = () => {
  const queryClient = useQueryClient()
  const { switchErrorHandling } = useError()

  const createFoodItemMutation = useMutation(
    (foodItem: Omit<FoodItem, 'id'>) =>
      axios.post<FoodItem>(
        `${process.env.REACT_APP_API_URL}/food-items`,
        foodItem
      ),
    {
      onSuccess: (res) => {
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
      onError: (err: any) => {
        if (err.response.data.message) {
          switchErrorHandling(err.response.data.message)
        } else {
          switchErrorHandling('食材の作成に失敗しました')
        }
      },
    }
  )

  const updateFoodItemMutation = useMutation(
    (foodItem: FoodItem) =>
      axios.put<FoodItem>(
        `${process.env.REACT_APP_API_URL}/food-items/${foodItem.id}`,
        foodItem
      ),
    {
      onSuccess: (res, variables) => {
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
      onError: (err: any) => {
        if (err.response.data.message) {
          switchErrorHandling(err.response.data.message)
        } else {
          switchErrorHandling('食材の更新に失敗しました')
        }
      },
    }
  )

  const deleteFoodItemMutation = useMutation(
    (id: number) =>
      axios.delete(`${process.env.REACT_APP_API_URL}/food-items/${id}`),
    {
      onSuccess: (_, variables) => {
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
      onError: (err: any) => {
        if (err.response.data.message) {
          switchErrorHandling(err.response.data.message)
        } else {
          switchErrorHandling('食材の削除に失敗しました')
        }
      },
    }
  )

  return {
    createFoodItemMutation,
    updateFoodItemMutation,
    deleteFoodItemMutation,
  }
}
