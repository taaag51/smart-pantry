import { AxiosError } from 'axios'
import { useError } from './useError'
import { useQuery } from '@tanstack/react-query'
import { Task } from '../types'
import axiosInstance, { getCsrfToken } from '../lib/axios'

interface ApiError {
  message: string
}

export const useQueryTasks = () => {
  const { switchErrorHandling } = useError()

  const getTasks = async () => {
    try {
      await getCsrfToken()
      const { data } = await axiosInstance.get<Task[]>('/tasks', {
        withCredentials: true,
      })
      return data
    } catch (err) {
      if (err instanceof AxiosError) {
        const axiosError = err as AxiosError<ApiError>
        if (axiosError.response?.data.message) {
          switchErrorHandling(axiosError.response.data.message)
        } else {
          switchErrorHandling('タスクの取得に失敗しました')
        }
      }
      throw err
    }
  }

  return useQuery({
    queryKey: ['tasks'],
    queryFn: getTasks,
    staleTime: Infinity,
  })
}
