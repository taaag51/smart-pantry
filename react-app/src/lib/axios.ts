import axios from 'axios'

const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080',
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
})

// リクエストインターセプターを追加
axiosInstance.interceptors.request.use(
  (config) => {
    // CSRFトークンがある場合は設定
    const csrfToken = axiosInstance.defaults.headers.common['X-CSRF-Token']
    if (csrfToken) {
      config.headers['X-CSRF-Token'] = csrfToken
    }

    // JWTトークンがある場合は設定
    const token = localStorage.getItem('accessToken')
    if (token && config.url !== '/csrf') {
      config.headers['Authorization'] = `Bearer ${token}`
    }

    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// レスポンスインターセプターを追加
axiosInstance.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    console.error('Response error:', error)

    // ネットワークエラーの場合
    if (!error.response) {
      window.dispatchEvent(
        new CustomEvent('error', {
          detail: 'サーバーとの通信に失敗しました',
        })
      )
      return Promise.reject(error)
    }

    // 認証エラーの場合
    if (error.response.status === 401) {
      window.dispatchEvent(new CustomEvent('unauthorized'))
    }

    // CSRFエラーの場合
    if (error.response.data?.message?.includes('csrf')) {
      try {
        await getCsrfToken()
        // 元のリクエストを再試行
        const config = error.config
        return axiosInstance(config)
      } catch (retryError) {
        console.error('CSRF retry failed:', retryError)
      }
    }

    return Promise.reject(error)
  }
)

export const getCsrfToken = async () => {
  try {
    const response = await axiosInstance.get('/csrf')
    const token = response.headers['x-csrf-token']
    if (token) {
      axiosInstance.defaults.headers.common['X-CSRF-Token'] = token
      return token
    } else {
      console.error('No CSRF token in response headers')
      throw new Error('CSRFトークンの取得に失敗しました')
    }
  } catch (error) {
    console.error('Failed to get CSRF token:', error)
    throw error
  }
}

export default axiosInstance
