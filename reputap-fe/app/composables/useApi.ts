import type { ApiResponse } from '~/utils/api'

export const useApi = () => {
  const { accessToken, refresh, logout } = useAuth()
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase as string

  async function apiFetch<T = any>(
    path: string,
    options: {
      method?: string
      body?: any
      query?: Record<string, any>
    } = {},
  ): Promise<ApiResponse<T>> {
    const headers: Record<string, string> = {}
    if (accessToken.value) {
      headers.Authorization = `Bearer ${accessToken.value}`
    }

    try {
      return await $fetch<ApiResponse<T>>(`${apiBase}${path}`, {
        method: options.method as any || 'GET',
        body: options.body,
        query: options.query,
        headers,
      })
    }
    catch (error: any) {
      // If 401, try refresh
      if (error?.statusCode === 401) {
        const refreshed = await refresh()
        if (refreshed) {
          headers.Authorization = `Bearer ${accessToken.value}`
          return await $fetch<ApiResponse<T>>(`${apiBase}${path}`, {
            method: options.method as any || 'GET',
            body: options.body,
            query: options.query,
            headers,
          })
        }
        else {
          logout()
        }
      }

      // Return error response
      if (error?.data) {
        return error.data as ApiResponse<T>
      }

      return {
        success: false,
        data: null as T,
        error: { code: 'NETWORK_ERROR', message: error?.message || 'Network error' },
      }
    }
  }

  return { apiFetch }
}
