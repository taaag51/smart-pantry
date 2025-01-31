import axios from 'axios'

const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
  withCredentials: true,
})

export const getCsrfToken = async () => {
  const { data } = await axios.get(`${process.env.REACT_APP_API_URL}/csrf`, {
    withCredentials: true,
  })
  axiosInstance.defaults.headers.common['X-CSRF-Token'] = data.csrf_token
}

export default axiosInstance
