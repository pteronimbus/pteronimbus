import { ref, computed } from 'vue'

interface User {
  id: string
  discord_user_id: string
  username: string
  avatar: string
  email: string
  created_at: string
  updated_at: string
}

interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: User
}

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}



// Global auth state
const authState = ref<AuthState>({
  user: null,
  accessToken: null,
  refreshToken: null,
  isAuthenticated: false,
  isLoading: false,
  error: null
})

export const useAuth = () => {
  const config = useRuntimeConfig()
  const router = useRouter()

  // Computed properties
  const user = computed(() => authState.value.user)
  const isAuthenticated = computed(() => authState.value.isAuthenticated)
  const isLoading = computed(() => authState.value.isLoading)
  const error = computed(() => authState.value.error)

  // Initialize auth state from localStorage
  const initializeAuth = () => {
    if (import.meta.client) {
      const accessToken = localStorage.getItem('access_token')
      const refreshToken = localStorage.getItem('refresh_token')
      const userData = localStorage.getItem('user_data')

      if (accessToken && refreshToken && userData) {
        try {
          const user = JSON.parse(userData)
          authState.value = {
            user,
            accessToken,
            refreshToken,
            isAuthenticated: true,
            isLoading: false,
            error: null
          }
        } catch (error) {
          console.error('Failed to parse stored user data:', error)
          clearAuth()
        }
      }
    }
  }

  // Clear auth state
  const clearAuth = () => {
    authState.value = {
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,
      error: null
    }

    if (import.meta.client) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user_data')
    }
  }

  // Store auth data
  const storeAuth = (authResponse: AuthResponse) => {
    authState.value = {
      user: authResponse.user,
      accessToken: authResponse.access_token,
      refreshToken: authResponse.refresh_token,
      isAuthenticated: true,
      isLoading: false,
      error: null
    }

    if (import.meta.client) {
      localStorage.setItem('access_token', authResponse.access_token)
      localStorage.setItem('refresh_token', authResponse.refresh_token)
      localStorage.setItem('user_data', JSON.stringify(authResponse.user))
    }
  }

  // Set error state
  const setError = (errorMessage: string) => {
    authState.value.error = errorMessage
    authState.value.isLoading = false
  }

  // Clear error state
  const clearError = () => {
    authState.value.error = null
  }

  // Login with Discord
  const signIn = async (provider: string = 'discord', options?: { callbackUrl?: string }) => {
    if (provider !== 'discord') {
      const errorMsg = 'Only Discord authentication is supported'
      setError(errorMsg)
      throw new Error(errorMsg)
    }

    authState.value.isLoading = true
    clearError()

    try {
      // Get Discord auth URL from backend
      const response = await $fetch<{ auth_url: string; state: string }>(`${config.public.backendUrl}/auth/login`)
      
      // Store callback URL for after authentication
      if (options?.callbackUrl && import.meta.client) {
        localStorage.setItem('auth_callback_url', options.callbackUrl)
      }

      // Redirect to Discord OAuth
      if (import.meta.client) {
        window.location.href = response.auth_url
      }
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to initiate Discord login'
      console.error('Failed to initiate Discord login:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    }
  }

  // Handle OAuth callback
  const handleCallback = async (code: string, state: string) => {
    authState.value.isLoading = true
    clearError()

    try {
      const authResponse = await $fetch<AuthResponse>(`${config.public.backendUrl}/auth/callback`, {
        method: 'GET',
        query: { code, state }
      })

      storeAuth(authResponse)

      // Redirect to callback URL or dashboard
      const callbackUrl = import.meta.client ? localStorage.getItem('auth_callback_url') : null
      if (import.meta.client) {
        localStorage.removeItem('auth_callback_url')
      }

      await router.push(callbackUrl || '/dashboard')
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to handle OAuth callback'
      console.error('Failed to handle OAuth callback:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    }
  }

  // Refresh access token
  const refreshAccessToken = async () => {
    if (!authState.value.refreshToken) {
      const errorMsg = 'No refresh token available'
      setError(errorMsg)
      throw new Error(errorMsg)
    }

    try {
      const authResponse = await $fetch<AuthResponse>(`${config.public.backendUrl}/auth/refresh`, {
        method: 'POST',
        body: {
          refresh_token: authState.value.refreshToken
        }
      })

      storeAuth(authResponse)
      return authResponse.access_token
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to refresh token'
      console.error('Failed to refresh token:', error)
      clearAuth()
      setError(errorMsg)
      throw new Error(errorMsg)
    }
  }

  // Sign out
  const signOut = async () => {
    authState.value.isLoading = true
    clearError()

    try {
      if (authState.value.accessToken) {
        await $fetch(`${config.public.backendUrl}/auth/logout`, {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${authState.value.accessToken}`
          }
        })
      }
    } catch (error: any) {
      console.error('Failed to logout from backend:', error)
      // Don't set error state for logout failures, just continue with local cleanup
    } finally {
      clearAuth()
      await router.push('/login')
    }
  }

  // Get current user info
  const getCurrentUser = async () => {
    if (!authState.value.accessToken) {
      const errorMsg = 'No access token available'
      setError(errorMsg)
      throw new Error(errorMsg)
    }

    try {
      const response = await $fetch<{ user: User }>(`${config.public.backendUrl}/auth/me`, {
        headers: {
          Authorization: `Bearer ${authState.value.accessToken}`
        }
      })

      authState.value.user = response.user
      if (import.meta.client) {
        localStorage.setItem('user_data', JSON.stringify(response.user))
      }

      return response.user
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to get current user'
      console.error('Failed to get current user:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    }
  }

  // API request with automatic token refresh
  const apiRequest = async <T>(url: string, options: any = {}) => {
    if (!authState.value.accessToken) {
      throw new Error('No access token available')
    }

    const requestOptions = {
      ...options,
      headers: {
        ...options.headers,
        Authorization: `Bearer ${authState.value.accessToken}`
      }
    }

    try {
      return await $fetch<T>(url, requestOptions)
    } catch (error: any) {
      // If token expired, try to refresh
      if (error.status === 401 && authState.value.refreshToken) {
        try {
          await refreshAccessToken()
          // Retry the request with new token
          requestOptions.headers.Authorization = `Bearer ${authState.value.accessToken}`
          return await $fetch<T>(url, requestOptions)
        } catch (refreshError) {
          console.error('Failed to refresh token:', refreshError)
          clearAuth()
          await router.push('/login')
          throw refreshError
        }
      }
      throw error
    }
  }

  return {
    // State
    user: readonly(user),
    isAuthenticated: readonly(isAuthenticated),
    isLoading: readonly(isLoading),
    error: readonly(error),

    // Methods
    signIn,
    signOut,
    handleCallback,
    refreshAccessToken,
    getCurrentUser,
    apiRequest,
    initializeAuth,
    clearAuth,
    clearError
  }
}