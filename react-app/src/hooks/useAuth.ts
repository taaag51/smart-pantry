import { useState, useEffect, useCallback } from 'react'
import axiosInstance from '../lib/axios'

export const useAuth = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(true)

  const checkAuth = useCallback(async () => {
    const token = localStorage.getItem('accessToken')

    if (!token) {
      setIsAuthenticated(false)
      setIsLoading(false)
      return
    }

    try {
      // トークンの検証のみを行い、新しいトークンは取得しない
      await axiosInstance.get('/verify-token')
      setIsAuthenticated(true)
    } catch (error) {
      console.error('Token verification failed:', error)
      // 401エラーの場合は、axiosインターセプターでリフレッシュトークンを使用して
      // 新しいアクセストークンの取得を試みる
      // 失敗した場合は、インターセプターで unauthorized イベントが発火される
    } finally {
      setIsLoading(false)
    }
  }, [])

  const clearAuth = useCallback(() => {
    localStorage.removeItem('accessToken')
    setIsAuthenticated(false)
    setIsLoading(false)
  }, [])

  useEffect(() => {
    // 初回認証チェック
    checkAuth()

    // 定期的な認証チェック（15分ごと - アクセストークンの有効期限に合わせる）
    const interval = setInterval(checkAuth, 15 * 60 * 1000)

    const handleUnauthorized = () => {
      clearAuth()
    }

    const handleLoginSuccess = () => {
      setIsAuthenticated(true)
      setIsLoading(false)
    }

    const handleLogoutSuccess = () => {
      clearAuth()
    }

    // イベントリスナーの設定
    window.addEventListener('unauthorized', handleUnauthorized)
    window.addEventListener('login-success', handleLoginSuccess)
    window.addEventListener('logout-success', handleLogoutSuccess)

    // クリーンアップ
    return () => {
      clearInterval(interval)
      window.removeEventListener('unauthorized', handleUnauthorized)
      window.removeEventListener('login-success', handleLoginSuccess)
      window.removeEventListener('logout-success', handleLogoutSuccess)
    }
  }, [checkAuth, clearAuth])

  return {
    isAuthenticated,
    isLoading,
    checkAuth,
    clearAuth,
  }
}
