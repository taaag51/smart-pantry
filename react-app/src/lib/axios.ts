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

    // CSRFエンドポイントでない場合は、Cookieに含まれるトークンを使用
    // Cookieは withCredentials: true で自動的に送信される

    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// レスポンスインターセプターを追加
axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    // ネットワークエラーの場合
    if (!error.response) {
      window.dispatchEvent(
        new CustomEvent('error', {
          detail: 'サーバーとの通信に失敗しました',
        })
      )
      return Promise.reject(error)
    }

    const originalRequest = error.config

    // 認証エラーの場合
    if (error.response.status === 401) {
      // リフレッシュトークンの取得試行回数を追跡
      if (!originalRequest._retryCount) {
        originalRequest._retryCount = 0
      }

      // 最大リトライ回数を超えた場合は、認証エラーとして処理
      if (
        originalRequest._retryCount >= 1 ||
        originalRequest.url === '/refresh-token'
      ) {
        window.dispatchEvent(
          new CustomEvent('unauthorized', {
            detail: 'セッションが期限切れです。再度ログインしてください。',
          })
        )
        return Promise.reject(error)
      }

      originalRequest._retryCount++

      try {
        // リフレッシュトークンを使用して新しいアクセストークンを取得
        const response = await axiosInstance.post('/refresh-token')
        const { data } = response.data

        if (data && data.accessToken) {
          // アクセストークンはCookieで自動的に設定されるため、
          // 元のリクエストを再試行するだけでよい
          return axiosInstance(originalRequest)
        } else {
          throw new Error('新しいアクセストークンの取得に失敗しました')
        }
      } catch (refreshError) {
        console.error('Token refresh failed:', refreshError)
        window.dispatchEvent(
          new CustomEvent('unauthorized', {
            detail: 'トークンの更新に失敗しました。再度ログインしてください。',
          })
        )
        return Promise.reject(error)
      }
    }

    // CSRFエラーの場合
    if (error.response.data?.message?.includes('csrf')) {
      try {
        await getCsrfToken()
        return axiosInstance(error.config)
      } catch (retryError) {
        console.error('CSRF retry failed:', retryError)
        return Promise.reject(retryError)
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
    }
    throw new Error('CSRFトークンの取得に失敗しました')
  } catch (error) {
    console.error('Failed to get CSRF token:', error)
    throw error
  }
}

export default axiosInstance
