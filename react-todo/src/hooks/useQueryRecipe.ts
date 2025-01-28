import axios from 'axios'
import { useQuery } from '@tanstack/react-query'
import { useError } from './useError'

export const useQueryRecipe = () => {
  const { switchErrorHandling } = useError()

  const getRecipe = async () => {
    try {
      const { data } = await axios.get(
        `${process.env.REACT_APP_API_URL}/recipes/suggestions`,
        {
          withCredentials: true,
        }
      )
      return data.recipe as string
    } catch (err: any) {
      if (err.response.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('レシピの取得に失敗しました')
      }
      throw err
    }
  }

  return useQuery({
    queryKey: ['recipe'],
    queryFn: getRecipe,
    staleTime: 1000 * 60 * 5, // 5分間はキャッシュを使用
    refetchInterval: 1000 * 60 * 5, // 5分ごとに自動更新
    retry: false,
  })
}
