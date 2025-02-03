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

interface LoginResponse {
  message: string
  data: {
    accessToken: string
    tokenType: string
    expiresIn: number
    expiresAt: string
    refreshToken: string
  }
}

export const useMutateAuth = () => {
  const navigate = useNavigate()
  const resetEditedTask = useStore((state) => state.resetEditedTask)
  const setAuth = useStore((state) => state.setAuth)
  const clearAuth = useStore((state) => state.clearAuth)
  const { switchErrorHandling } = useError()

  const loginMutation = useMutation<
    AxiosResponse<LoginResponse>,
    AxiosError<ApiError>,
    Credential
  >({
    mutationFn: async (user: Credential) => {
      try {
        await getCsrfToken()
        return await axiosInstance.post<LoginResponse>('/login', user)
      } catch (err) {
        if (err instanceof Error) {
          if (err.message.includes('CSRF')) {
            throw new Error(
              'セキュリティトークンの取得に失敗しました。ページを再読み込みしてください。'
            )
          }
        }
        throw err
      }
    },
    onSuccess: (res) => {
      console.log('ログインレスポンス:', res.data)
      const { accessToken } = res.data.data
      if (accessToken) {
        console.log('アクセストークンを設定:', accessToken)
        localStorage.setItem('accessToken', accessToken)
        // Axiosのデフォルトヘッダーにトークンを設定
        axiosInstance.defaults.headers.common[
          'Authorization'
        ] = `Bearer ${accessToken}`
        setAuth(accessToken)
        // ログイン成功イベントを発火
        window.dispatchEvent(new CustomEvent('login-success'))
        navigate('/pantry', { replace: true })
      } else {
        console.error('No token in login response')
        switchErrorHandling('ログインに失敗しました。もう一度お試しください。')
      }
    },
    onError: (err: AxiosError<ApiError>) => {
      console.error('Login error:', err)
      if (err.message.includes('CSRF')) {
        switchErrorHandling(
          'セキュリティトークンの取得に失敗しました。ページを再読み込みしてください。'
        )
      } else if (err.response?.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling('ログインに失敗しました。もう一度お試しください。')
      }
    },
  })

  const registerMutation = useMutation<
    AxiosResponse,
    AxiosError<ApiError>,
    Credential
  >({
    mutationFn: async (user: Credential) => {
      try {
        await getCsrfToken()
        return await axiosInstance.post('/signup', user)
      } catch (err) {
        if (err instanceof Error) {
          if (err.message.includes('CSRF')) {
            throw new Error(
              'セキュリティトークンの取得に失敗しました。ページを再読み込みしてください。'
            )
          }
        }
        throw err
      }
    },
    onSuccess: () => {
      // 登録成功のみを処理し、ログインは別途行う
      console.log('Registration successful')
    },
    onError: (err: AxiosError<ApiError>) => {
      if (err.message.includes('CSRF')) {
        switchErrorHandling(
          'セキュリティトークンの取得に失敗しました。ページを再読み込みしてください。'
        )
      } else if (err.response?.data.message) {
        switchErrorHandling(err.response.data.message)
      } else {
        switchErrorHandling(
          'アカウント作成に失敗しました。もう一度お試しください。'
        )
      }
    },
  })

  const logoutMutation = useMutation<AxiosResponse, AxiosError<ApiError>, void>(
    {
      mutationFn: async () => {
        try {
          await getCsrfToken()
          return await axiosInstance.post('/logout')
        } catch (err) {
          if (err instanceof Error) {
            if (err.message.includes('CSRF')) {
              throw new Error(
                'セキュリティトークンの取得に失敗しました。ページを再読み込みしてください。'
              )
            }
          }
          throw err
        }
      },
      onSuccess: () => {
        resetEditedTask()
        clearAuth()
        navigate('/', { replace: true })
      },
      onError: (err: AxiosError<ApiError>) => {
        console.error('Logout error:', err)
        // エラーが発生しても、ローカルのクリーンアップは実行
        resetEditedTask()
        clearAuth()
        navigate('/', { replace: true })
      },
    }
  )

  return { loginMutation, registerMutation, logoutMutation }
}
