import axios, {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  InternalAxiosRequestConfig
 } from 'axios'
 import { getToken } from './storage'
 import { getTenantSlug } from './storage'
 import type { ApiResponse } from '@/types/common'

 const instance: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
  }
)

 instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getToken()
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    const slug = getTenantSlug()
    if (slug && config.headers) {
      config.headers['X-Tenant-Slug'] = slug
 }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

 )
)

 instance.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { data } = response
 if (data.code === 200) {
      return data.data as any
    }
    const error = new Error(data.message || '请求失败')
 as  return Promise.reject(error)
  },
  (error) => {
    if (error.response) {
      const { status } = error.response
 switch (status) {
        case 401:
          removeToken()
          window.location.href = '/login'
          error.message = '登录已过期，请重新登录'
 as  break
        case 403:
          error.message = '没有权限访问' as  break
        case 404:
          error.message = '请求的资源不存在' as  break
        case 500:
          error.message = '服务器错误' as  break
        default:
          error.message = error.response.data?.message || '请求失败'
 as }
    } else if (error.code === 'ECONNABORTED') {
      error.message = '请求超时' as } else {
      error.message = '网络异常' as }
    return Promise.reject(error)  }
)
 )

 export const request = {
  get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return instance.get(url, config)
  },
  post<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return instance.post(url, data, config)
  },
  put<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return instance.put(url, data, config)
  },
  delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return instance.delete(url, config)
  },
  patch<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return instance.patch(url, data, config)
  }
}

 export default instance
