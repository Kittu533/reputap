import type { ApiResponse, LoginRequest, TokenResponse, UserInfo } from '~/utils/api'

export const useAuth = () => {
  const accessToken = useCookie('reputap_access_token', { maxAge: 900 })
  const refreshToken = useCookie('reputap_refresh_token', { maxAge: 604800 })
  const user = useState<UserInfo | null>('auth_user', () => null)
  const isAuthenticated = computed(() => !!accessToken.value)

  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase as string

  async function login(credentials: LoginRequest): Promise<{ success: boolean; error?: string }> {
    try {
      const response = await $fetch<ApiResponse<TokenResponse>>(`${apiBase}/api/v1/auth/login`, {
        method: 'POST',
        body: credentials,
      })

      if (response.success && response.data) {
        accessToken.value = response.data.access_token
        refreshToken.value = response.data.refresh_token
        user.value = response.data.user
        return { success: true }
      }

      return { success: false, error: response.error?.message || 'Login failed' }
    }
    catch (error: any) {
      const message = error?.data?.error?.message || error?.message || 'Login failed'
      return { success: false, error: message }
    }
  }

  async function fetchUser(): Promise<void> {
    if (!accessToken.value) return

    try {
      const response = await $fetch<ApiResponse<UserInfo>>(`${apiBase}/api/v1/auth/me`, {
        headers: { Authorization: `Bearer ${accessToken.value}` },
      })

      if (response.success && response.data) {
        user.value = response.data
      }
    }
    catch {
      // Token might be expired, try refresh
      await refresh()
    }
  }

  async function refresh(): Promise<boolean> {
    if (!refreshToken.value) return false

    try {
      const response = await $fetch<ApiResponse<TokenResponse>>(`${apiBase}/api/v1/auth/refresh`, {
        method: 'POST',
        body: { refresh_token: refreshToken.value },
      })

      if (response.success && response.data) {
        accessToken.value = response.data.access_token
        refreshToken.value = response.data.refresh_token
        user.value = response.data.user
        return true
      }

      return false
    }
    catch {
      logout()
      return false
    }
  }

  function logout() {
    accessToken.value = null
    refreshToken.value = null
    user.value = null
    navigateTo('/login')
  }

  return {
    accessToken,
    refreshToken,
    user,
    isAuthenticated,
    login,
    fetchUser,
    refresh,
    logout,
  }
}
