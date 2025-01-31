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

    // cookieベースの認証を使用するため、ここでの追加のトークン設定は不要

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
    const { data } = await axiosInstance.get('/csrf')
    if (data.csrf_token) {
      axiosInstance.defaults.headers.common['X-CSRF-Token'] = data.csrf_token
    } else {
      console.error('No CSRF token in response')
    }
  } catch (error) {
    console.error('Failed to get CSRF token:', error)
    throw error
  }
}

export default axiosInstance
