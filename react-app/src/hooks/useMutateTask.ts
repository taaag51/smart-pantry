import { AxiosError, AxiosResponse } from 'axios'
import { useQueryClient, useMutation } from '@tanstack/react-query'
import { Task } from '../types'
import useStore from '../store'
import { useError } from '../hooks/useError'
import axiosInstance, { getCsrfToken } from '../lib/axios'

interface ApiError {
  message: string
}

export const useMutateTask = () => {
  const queryClient = useQueryClient()
  const { switchErrorHandling } = useError()
  const resetEditedTask = useStore((state) => state.resetEditedTask)

  const createTaskMutation = useMutation({
    mutationFn: async (
      task: Omit<Task, 'id' | 'created_at' | 'updated_at'>
    ) => {
      await getCsrfToken()
      return await axiosInstance.post<Task>('/tasks', task)
    },
    onSuccess: (res: AxiosResponse<Task>) => {
      const previousTasks = queryClient.getQueryData<Task[]>(['tasks'])
      if (previousTasks) {
        queryClient.setQueryData(['tasks'], [...previousTasks, res.data])
      }
      resetEditedTask()
    },
    onError: (err: AxiosError<ApiError>) => {
      if (err.response?.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('タスクの作成に失敗しました')
      }
    },
  })

  const updateTaskMutation = useMutation({
    mutationFn: async (task: Omit<Task, 'created_at' | 'updated_at'>) => {
      await getCsrfToken()
      return await axiosInstance.put<Task>(`/tasks/${task.id}`, {
        title: task.title,
      })
    },
    onSuccess: (res: AxiosResponse<Task>, variables) => {
      const previousTasks = queryClient.getQueryData<Task[]>(['tasks'])
      if (previousTasks) {
        queryClient.setQueryData<Task[]>(
          ['tasks'],
          previousTasks.map((task) =>
            task.id === variables.id ? res.data : task
          )
        )
      }
      resetEditedTask()
    },
    onError: (err: AxiosError<ApiError>) => {
      if (err.response?.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('タスクの更新に失敗しました')
      }
    },
  })

  const deleteTaskMutation = useMutation({
    mutationFn: async (id: number) => {
      await getCsrfToken()
      return await axiosInstance.delete(`/tasks/${id}`)
    },
    onSuccess: (_, variables) => {
      const previousTasks = queryClient.getQueryData<Task[]>(['tasks'])
      if (previousTasks) {
        queryClient.setQueryData<Task[]>(
          ['tasks'],
          previousTasks.filter((task) => task.id !== variables)
        )
      }
      resetEditedTask()
    },
    onError: (err: AxiosError<ApiError>) => {
      if (err.response?.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('タスクの削除に失敗しました')
      }
    },
  })

  return {
    createTaskMutation,
    updateTaskMutation,
    deleteTaskMutation,
  }
}
