import axios from 'axios'
import { useQuery } from '@tanstack/react-query'
import { FoodItem } from '../types'
import { useError } from './useError'

export const useQueryFoodItems = () => {
  const { switchErrorHandling } = useError()

  const getFoodItems = async () => {
    const { data } = await axios.get<FoodItem[]>(
      `${process.env.REACT_APP_API_URL}/food-items`
    )
    return data
  }

  return useQuery<FoodItem[], Error>({
    queryKey: ['foodItems'],
    queryFn: getFoodItems,
    onError: (err: any) => {
      if (err.response.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('食材の取得に失敗しました')
      }
    },
  })
}
