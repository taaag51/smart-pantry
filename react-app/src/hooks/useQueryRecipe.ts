import { useQuery } from '@tanstack/react-query'
import axiosInstance from '../lib/axios'
import { Recipe } from '../types'

export const useQueryRecipe = () => {
  const getRecipes = async () => {
    try {
      const { data } = await axiosInstance.get<Recipe[]>('/recipes/suggestions')
      return data || []
    } catch (error) {
      console.error('Failed to fetch recipes:', error)
      return []
    }
  }

  return useQuery({
    queryKey: ['recipes'],
    queryFn: getRecipes,
    staleTime: 1000 * 60 * 5,
    initialData: [],
    retry: 1,
  })
}
