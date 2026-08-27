import axios from 'axios'
import { type InternalAxiosRequestConfig, type AxiosResponse } from 'axios'
const request = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL,
  timeout: 6000,
})
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    config.headers['Content-Type'] = 'application/json'
    config.headers['Authorization'] = 'Bearer ' + localStorage.getItem('token') || ''
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)
request.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data
  },
  (error) => {
    return Promise.reject(error)
  },
)
export default request
