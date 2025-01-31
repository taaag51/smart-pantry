import { useState, useEffect } from 'react'
import axiosInstance from '../lib/axios'

export const useAuth = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(true)

  const checkAuth = async () => {
    try {
      await axiosInstance.get('/verify-token')
      setIsAuthenticated(true)
    } catch (error) {
      console.error('Token verification failed:', error)
      setIsAuthenticated(false)
      window.dispatchEvent(new CustomEvent('unauthorized'))
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    // 初回認証チェック
    checkAuth()

    // 定期的な認証チェック（5分ごと）
    const interval = setInterval(checkAuth, 5 * 60 * 1000)

    const handleUnauthorized = () => {
      setIsAuthenticated(false)
      setIsLoading(false)
    }

    const handleLoginSuccess = () => {
      setIsAuthenticated(true)
      setIsLoading(false)
    }

    const handleLogoutSuccess = () => {
      setIsAuthenticated(false)
      setIsLoading(false)
    }

    window.addEventListener('unauthorized', handleUnauthorized)
    window.addEventListener('login-success', handleLoginSuccess)
    window.addEventListener('logout-success', handleLogoutSuccess)

    return () => {
      clearInterval(interval)
      window.removeEventListener('unauthorized', handleUnauthorized)
      window.removeEventListener('login-success', handleLoginSuccess)
      window.removeEventListener('logout-success', handleLogoutSuccess)
    }
  }, [])

  return { isAuthenticated, isLoading }
}
