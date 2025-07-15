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
}

// Global auth state
const authState = ref<AuthState>({
  user: null,
  accessToken: null,
  refreshToken: null,
  isAuthenticated: false,
  isLoading: false
})

export const useAuth = () => {
  const config = useRuntimeConfig()
  const router = useRouter()

  // Computed properties
  const user = computed(() => authState.value.user)
  const isAuthenticated = computed(() => authState.value.isAuthenticated)
  const isLoading = computed(() => authState.value.isLoading)

  // Initialize auth state from localStorage
  const initializeAuth = () => {
    if (process.client) {
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
            isLoading: false
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
      isLoading: false
    }

    if (process.client) {
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
      isLoading: false
    }

    if (process.client) {
      localStorage.setItem('access_token', authResponse.access_token)
      localStorage.setItem('refresh_token', authResponse.refresh_token)
      localStorage.setItem('user_data', JSON.stringify(authResponse.user))
    }
  }

  // Login with Discord
  const signIn = async (provider: string = 'discord', options?: { callbackUrl?: string }) => {
    if (provider !== 'discord') {
      throw new Error('Only Discord authentication is supported')
    }

    authState.value.isLoading = true

    try {
      // Get Discord auth URL from backend
      const response = await $fetch<{ auth_url: string; state: string }>(`${config.public.backendUrl}/auth/login`)
      
      // Store callback URL for after authentication
      if (options?.callbackUrl && process.client) {
        localStorage.setItem('auth_callback_url', options.callbackUrl)
      }

      // Redirect to Discord OAuth
      if (process.client) {
        window.location.href = response.auth_url
      }
    } catch (error) {
      console.error('Failed to initiate Discord login:', error)
      authState.value.isLoading = false
      throw error
    }
  }

  // Handle OAuth callback
  const handleCallback = async (code: string, state: string) => {
    authState.value.isLoading = true

    try {
      const authResponse = await $fetch<AuthResponse>(`${config.public.backendUrl}/auth/callback`, {
        method: 'GET',
        query: { code, state }
      })

      storeAuth(authResponse)

      // Redirect to callback URL or dashboard
      const callbackUrl = process.client ? localStorage.getItem('auth_callback_url') : null
      if (process.client) {
        localStorage.removeItem('auth_callback_url')
      }

      await router.push(callbackUrl || '/dashboard')
    } catch (error) {
      console.error('Failed to handle OAuth callback:', error)
      authState.value.isLoading = false
      throw error
    }
  }

  // Refresh access token
  const refreshAccessToken = async () => {
    if (!authState.value.refreshToken) {
      throw new Error('No refresh token available')
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
    } catch (error) {
      console.error('Failed to refresh token:', error)
      clearAuth()
      throw error
    }
  }

  // Sign out
  const signOut = async () => {
    authState.value.isLoading = true

    try {
      if (authState.value.accessToken) {
        await $fetch(`${config.public.backendUrl}/auth/logout`, {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${authState.value.accessToken}`
          }
        })
      }
    } catch (error) {
      console.error('Failed to logout from backend:', error)
    } finally {
      clearAuth()
      await router.push('/login')
    }
  }

  // Get current user info
  const getCurrentUser = async () => {
    if (!authState.value.accessToken) {
      throw new Error('No access token available')
    }

    try {
      const response = await $fetch<{ user: User }>(`${config.public.backendUrl}/auth/me`, {
        headers: {
          Authorization: `Bearer ${authState.value.accessToken}`
        }
      })

      authState.value.user = response.user
      if (process.client) {
        localStorage.setItem('user_data', JSON.stringify(response.user))
      }

      return response.user
    } catch (error) {
      console.error('Failed to get current user:', error)
      throw error
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

    // Methods
    signIn,
    signOut,
    handleCallback,
    refreshAccessToken,
    getCurrentUser,
    apiRequest,
    initializeAuth,
    clearAuth
  }
}