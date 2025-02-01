import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import useStore from '../store'
import axiosInstance from '../lib/axios'

export const useError = () => {
  const navigate = useNavigate()
  const resetEditedTask = useStore((state) => state.resetEditedTask)
  const [errorMessage, setErrorMessage] = useState<string>('')

  const getCsrfToken = async () => {
    try {
      const { data } = await axiosInstance.get('/csrf')
      if (data.csrf_token) {
        axiosInstance.defaults.headers.common['X-CSRF-Token'] = data.csrf_token
      }
    } catch (error) {
      console.error('Failed to get CSRF token:', error)
      setErrorMessage('サーバーとの通信に失敗しました')
    }
  }

  const switchErrorHandling = (msg: string) => {
    switch (msg) {
      case 'invalid csrf token':
        getCsrfToken()
        setErrorMessage('セッションが無効です。もう一度お試しください')
        break
      case 'invalid email or password':
        setErrorMessage('メールアドレスまたはパスワードが正しくありません')
        break
      case 'invalid or expired jwt':
        localStorage.removeItem('accessToken')
        setErrorMessage(
          'セッションの有効期限が切れました。再度ログインしてください'
        )
        resetEditedTask()
        navigate('/')
        break
      case 'missing or malformed jwt':
        localStorage.removeItem('accessToken')
        setErrorMessage('認証情報が無効です。再度ログインしてください')
        resetEditedTask()
        navigate('/')
        break
      case 'email already exists':
        setErrorMessage('このメールアドレスは既に登録されています')
        break
      case 'failed to create user':
        setErrorMessage('ユーザー登録に失敗しました')
        break
      case 'failed to generate token':
        setErrorMessage('ログインに失敗しました。もう一度お試しください')
        break
      case 'record not found':
        setErrorMessage('メールアドレスまたはパスワードが正しくありません')
        break
      default:
        setErrorMessage(msg)
    }
  }

  return { switchErrorHandling, errorMessage, setErrorMessage }
}
