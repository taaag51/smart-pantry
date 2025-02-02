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
    token: string
  }
}

export const useMutateAuth = () => {
  const navigate = useNavigate()
  const resetEditedTask = useStore((state) => state.resetEditedTask)
  const { switchErrorHandling } = useError()

  const loginMutation = useMutation<
    AxiosResponse<LoginResponse>,
    AxiosError<ApiError>,
    Credential
  >({
    mutationFn: async (user: Credential) => {
      try {
        // CSRFトークンを取得
        await getCsrfToken()
        // ログインリクエスト
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
    onSuccess: async (res) => {
      const token = res.data.data.token
      if (token) {
        localStorage.setItem('accessToken', token)
        // 認証状態を更新
        window.dispatchEvent(new CustomEvent('login-success'))
        // トークンの設定とイベントの反映を待つ
        await new Promise((resolve) => setTimeout(resolve, 100))
        axiosInstance.defaults.headers.common[
          'Authorization'
        ] = `Bearer ${token}`
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
        // CSRFトークンを取得
        await getCsrfToken()
        // アカウント登録リクエスト
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
          // CSRFトークンを取得
          await getCsrfToken()
          // ログアウトリクエスト
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
      onSuccess: async () => {
        localStorage.removeItem('accessToken')
        resetEditedTask()
        window.dispatchEvent(new CustomEvent('logout-success'))
        // 状態の更新が反映されるのを待ってからナビゲーション
        await new Promise((resolve) => setTimeout(resolve, 0))
        navigate('/', { replace: true })
      },
      onError: async (err: AxiosError<ApiError>) => {
        console.error('Logout error:', err)
        // エラーが発生しても、ローカルのクリーンアップは実行
        localStorage.removeItem('accessToken')
        resetEditedTask()
        window.dispatchEvent(new CustomEvent('logout-success'))
        // 状態の更新が反映されるのを待ってからナビゲーション
        await new Promise((resolve) => setTimeout(resolve, 0))
        navigate('/', { replace: true })
      },
    }
  )

  return { loginMutation, registerMutation, logoutMutation }
}
