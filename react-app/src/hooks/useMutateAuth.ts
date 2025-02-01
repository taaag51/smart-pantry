import { AxiosError, AxiosResponse } from 'axios'
import { useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import useStore from '../store'
import { Credential } from '../types'
import { useError } from '../hooks/useError'
import axiosInstance, { getCsrfToken } from '../lib/axios'

interface ApiError {
  message: string
}

export const useMutateAuth = () => {
  const navigate = useNavigate()
  const resetEditedTask = useStore((state) => state.resetEditedTask)
  const { switchErrorHandling } = useError()

  const loginMutation = useMutation<
    AxiosResponse,
    AxiosError<ApiError>,
    Credential
  >({
    mutationFn: async (user: Credential) => {
      await getCsrfToken()
      return await axiosInstance.post('/login', user)
    },
    onSuccess: () => {
      window.dispatchEvent(new CustomEvent('login-success'))
      navigate('/pantry', { replace: true })
    },
    onError: (err: AxiosError<ApiError>) => {
      console.error('Login error:', err)
      if (err.response?.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('ログインに失敗しました')
      }
    },
  })

  const registerMutation = useMutation<
    AxiosResponse,
    AxiosError<ApiError>,
    Credential
  >({
    mutationFn: async (user: Credential) => {
      await getCsrfToken()
      return await axiosInstance.post('/signup', user)
    },
    onError: (err: AxiosError<ApiError>) => {
      if (err.response?.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('アカウント作成に失敗しました')
      }
    },
  })

  const logoutMutation = useMutation<AxiosResponse, AxiosError<ApiError>, void>(
    {
      mutationFn: async () => {
        await getCsrfToken()
        return await axiosInstance.post('/logout')
      },
      onSuccess: () => {
        window.dispatchEvent(new CustomEvent('logout-success'))
        localStorage.removeItem('accessToken')
        resetEditedTask()
        navigate('/login', { replace: true })
      },
      onError: (err: AxiosError<ApiError>) => {
        if (err.response?.data.message) {
          switchErrorHandling(err.response.data.message)
        } else {
          switchErrorHandling('ログアウトに失敗しました')
        }
      },
    }
  )

  return { loginMutation, registerMutation, logoutMutation }
}
