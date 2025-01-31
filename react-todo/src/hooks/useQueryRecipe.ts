import { useQuery } from '@tanstack/react-query'
import axios from 'axios'
import { Recipe } from '../types'

export const useQueryRecipe = () => {
  const getRecipes = async () => {
    try {
      const { data } = await axios.get<Recipe[]>(
        `${process.env.REACT_APP_API_URL}/recipes/suggestions`
      )
      return data || [] // データがない場合は空配列を返す
    } catch (error) {
      console.error('Failed to fetch recipes:', error)
      return [] // エラー時は空配列を返す
    }
  }

  return useQuery({
    queryKey: ['recipes'],
    queryFn: getRecipes,
    staleTime: 1000 * 60 * 5, // 5分間キャッシュ
    initialData: [], // 初期値として空配列を設定
    retry: 1, // リトライ回数を1回に制限
  })
}
