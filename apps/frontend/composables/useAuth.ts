import { ref, computed } from 'vue'

// JWT token decoder utility
const decodeJWT = (token: string): any => {
  try {
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
    }).join(''))
    return JSON.parse(jsonPayload)
  } catch (error) {
    console.error('Failed to decode JWT token:', error)
    return null
  }
}

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
  isSuperAdmin: boolean
}



// Global auth state
const authState = ref<AuthState>({
  user: null,
  accessToken: null,
  refreshToken: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
  isSuperAdmin: false
})

export const useAuth = () => {
  const config = useRuntimeConfig()
  const router = useRouter()

  // Computed properties
  const user = computed(() => authState.value.user)
  const isAuthenticated = computed(() => authState.value.isAuthenticated)
  const isLoading = computed(() => authState.value.isLoading)
  const error = computed(() => authState.value.error)
  const isSuperAdmin = computed(() => authState.value.isSuperAdmin)

  // Initialize auth state from localStorage
  const initializeAuth = () => {
    if (import.meta.client) {
      const accessToken = localStorage.getItem('access_token')
      const refreshToken = localStorage.getItem('refresh_token')
      const userData = localStorage.getItem('user_data')

      if (accessToken && refreshToken && userData) {
        try {
          const user = JSON.parse(userData)
          // Decode JWT to get super admin status from system roles
          const tokenPayload = decodeJWT(accessToken)
          const systemRoles = tokenPayload?.system_roles || []
          const isSuperAdmin = systemRoles.includes('superadmin')
          
          authState.value = {
            user,
            accessToken,
            refreshToken,
            isAuthenticated: true,
            isLoading: false,
            error: null,
            isSuperAdmin
          }
          console.log('Auth state initialized from localStorage:', { isAuthenticated: true, user: user.username, isSuperAdmin })
        } catch (error) {
          console.error('Failed to parse stored user data:', error)
          clearAuth()
        }
      } else {
        console.log('No stored auth data found')
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
      error: null,
      isSuperAdmin: false
    }

    if (import.meta.client) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user_data')
    }
  }

  // Store auth data
  const storeAuth = (authResponse: AuthResponse) => {
    // Decode JWT to get super admin status from system roles
    const tokenPayload = decodeJWT(authResponse.access_token)
    const systemRoles = tokenPayload?.system_roles || []
    const isSuperAdmin = systemRoles.includes('superadmin')
    
    authState.value = {
      user: authResponse.user,
      accessToken: authResponse.access_token,
      refreshToken: authResponse.refresh_token,
      isAuthenticated: true,
      isLoading: false,
      error: null,
      isSuperAdmin
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

      // Wait a tick to ensure state is properly set
      await nextTick()

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

    console.log('Making API request:', {
      url,
      method: options.method || 'GET',
      hasToken: !!authState.value.accessToken,
      tokenLength: authState.value.accessToken?.length
    })

    try {
      return await $fetch<T>(url, requestOptions)
    } catch (error: any) {
      console.error('API request failed:', {
        url,
        status: error.status,
        statusText: error.statusText,
        data: error.data
      })
      
      // If token expired, try to refresh
      if (error.status === 401 && authState.value.refreshToken) {
        console.log('Attempting token refresh...')
        try {
          await refreshAccessToken()
          // Retry the request with new token
          requestOptions.headers.Authorization = `Bearer ${authState.value.accessToken}`
          console.log('Retrying request with new token...')
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

  // Handle tokens from URL parameters (for backend redirect flow)
  const handleTokensFromUrl = async (accessToken: string, refreshToken: string, expiresIn: number) => {
    const authResponse: AuthResponse = {
      access_token: accessToken,
      refresh_token: refreshToken,
      expires_in: expiresIn,
      user: null as any // Will be set after getCurrentUser call
    }

    storeAuth(authResponse)
    
    // Get user info now that we have tokens
    const user = await getCurrentUser()
    return user
  }

  return {
    // State
    user: readonly(user),
    isAuthenticated: readonly(isAuthenticated),
    isLoading: readonly(isLoading),
    error: readonly(error),
    isSuperAdmin: readonly(isSuperAdmin),

    // Methods
    signIn,
    signOut,
    handleCallback,
    handleTokensFromUrl,
    refreshAccessToken,
    getCurrentUser,
    apiRequest,
    initializeAuth,
    clearAuth,
    clearError
  }
}