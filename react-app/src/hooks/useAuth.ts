import { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import axiosInstance from '../lib/axios'
import useStore from '../store'

export const useAuth = () => {
  const navigate = useNavigate()
  const { isAuthenticated, setAuth, clearAuth } = useStore()
  const [isLoading, setIsLoading] = useState(true)

  const checkAuth = useCallback(async () => {
    const token = localStorage.getItem('accessToken')

    if (!token) {
      clearAuth()
      setIsLoading(false)
      return
    }

    try {
      await axiosInstance.get('/verify-token')
      setAuth(token)
    } catch (error) {
      console.error('Token verification failed:', error)
      clearAuth()
      navigate('/', { replace: true })
    } finally {
      setIsLoading(false)
    }
  }, [clearAuth, setAuth, navigate])

  useEffect(() => {
    checkAuth()

    const interval = setInterval(checkAuth, 15 * 60 * 1000)

    const handleUnauthorized = () => {
      clearAuth()
      navigate('/', { replace: true })
    }

    window.addEventListener('unauthorized', handleUnauthorized)

    return () => {
      clearInterval(interval)
      window.removeEventListener('unauthorized', handleUnauthorized)
    }
  }, [checkAuth, clearAuth, navigate])

  return {
    isAuthenticated,
    isLoading,
    checkAuth,
  }
}
