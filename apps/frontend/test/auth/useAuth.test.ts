import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { ref } from 'vue'

// Mock dependencies
const mockRouterPush = vi.fn()
const mockFetch = vi.fn()
const mockLocalStorage = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}

// Mock Nuxt composables
vi.mock('#app', () => ({
  useRuntimeConfig: () => ({
    public: {
      backendUrl: 'http://localhost:8080'
    }
  }),
  useRouter: () => ({
    push: mockRouterPush
  }),
  $fetch: mockFetch
}))

// Mock global dependencies
global.$fetch = mockFetch
Object.defineProperty(window, 'localStorage', { value: mockLocalStorage })
Object.defineProperty(window, 'location', { value: { href: '' }, writable: true })

// Mock import.meta.client
Object.defineProperty(import.meta, 'client', {
  value: true,
  writable: true
})

describe('useAuth Composable', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockLocalStorage.getItem.mockReturnValue(null)
    mockRouterPush.mockResolvedValue(undefined)
    window.location.href = ''
    
    // Reset modules to avoid shared state
    vi.resetModules()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('Initial State', () => {
    it('should have correct default state', async () => {
      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()

      expect(auth.user.value).toBeNull()
      expect(auth.isAuthenticated.value).toBe(false)
      expect(auth.isLoading.value).toBe(false)
      expect(auth.error.value).toBeNull()
    })
  })

  describe('Sign In', () => {
    it('should initiate Discord OAuth flow successfully', async () => {
      const mockAuthUrl = 'https://discord.com/oauth2/authorize?client_id=123'
      mockFetch.mockResolvedValueOnce({
        auth_url: mockAuthUrl,
        state: 'test-state'
      })

      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()

      await auth.signIn('discord', { callbackUrl: '/dashboard' })

      expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/auth/login')
      expect(mockLocalStorage.setItem).toHaveBeenCalledWith('auth_callback_url', '/dashboard')
      expect(window.location.href).toBe(mockAuthUrl)
    })

    it('should reject unsupported providers', async () => {
      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()

      await expect(auth.signIn('google')).rejects.toThrow('Only Discord authentication is supported')
      expect(auth.error.value).toBe('Only Discord authentication is supported')
    })

    it('should handle API errors', async () => {
      mockFetch.mockRejectedValueOnce({
        data: { message: 'Server error' }
      })

      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()

      await expect(auth.signIn('discord')).rejects.toThrow('Server error')
      expect(auth.error.value).toBe('Server error')
    })
  })

  describe('OAuth Callback', () => {
    it('should handle successful callback', async () => {
      const mockUser = {
        id: '1',
        discord_user_id: '123456789',
        username: 'testuser',
        avatar: 'avatar.png',
        email: 'test@example.com',
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T00:00:00Z'
      }

      mockFetch.mockResolvedValueOnce({
        access_token: 'access-token',
        refresh_token: 'refresh-token',
        expires_in: 3600,
        user: mockUser
      })

      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()

      await auth.handleCallback('auth-code', 'state')

      expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/auth/callback', {
        method: 'GET',
        query: { code: 'auth-code', state: 'state' }
      })

      expect(auth.user.value).toEqual(mockUser)
      expect(auth.isAuthenticated.value).toBe(true)
      expect(mockLocalStorage.setItem).toHaveBeenCalledWith('access_token', 'access-token')
      expect(mockLocalStorage.setItem).toHaveBeenCalledWith('refresh_token', 'refresh-token')
    })

    it('should handle callback errors', async () => {
      mockFetch.mockRejectedValueOnce({
        data: { message: 'Invalid code' }
      })

      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()

      await expect(auth.handleCallback('invalid', 'state')).rejects.toThrow('Invalid code')
      expect(auth.error.value).toBe('Invalid code')
    })
  })

  describe('Token Management', () => {
    it('should restore auth state from localStorage', async () => {
      const mockUser = {
        id: '1',
        discord_user_id: '123456789',
        username: 'testuser',
        avatar: 'avatar.png',
        email: 'test@example.com',
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T00:00:00Z'
      }

      mockLocalStorage.getItem.mockImplementation((key) => {
        switch (key) {
          case 'access_token': return 'stored-token'
          case 'refresh_token': return 'stored-refresh'
          case 'user_data': return JSON.stringify(mockUser)
          default: return null
        }
      })

      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()
      
      auth.initializeAuth()

      expect(auth.user.value).toEqual(mockUser)
      expect(auth.isAuthenticated.value).toBe(true)
    })

    it('should handle corrupted localStorage data', async () => {
      mockLocalStorage.getItem.mockImplementation((key) => {
        switch (key) {
          case 'access_token': return 'stored-token'
          case 'refresh_token': return 'stored-refresh'
          case 'user_data': return 'invalid-json'
          default: return null
        }
      })

      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()
      
      auth.initializeAuth()

      expect(auth.user.value).toBeNull()
      expect(auth.isAuthenticated.value).toBe(false)
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('access_token')
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('refresh_token')
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('user_data')
    })
  })

  describe('Sign Out', () => {
    it('should call logout API and clear localStorage', async () => {
      // Set up authenticated state
      mockLocalStorage.getItem.mockImplementation((key) => {
        switch (key) {
          case 'access_token': return 'access-token'
          case 'refresh_token': return 'refresh-token'
          case 'user_data': return JSON.stringify({ id: '1', username: 'test' })
          default: return null
        }
      })

      mockFetch.mockResolvedValueOnce({ message: 'Logged out' })

      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()
      
      auth.initializeAuth()
      await auth.signOut()

      expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/auth/logout', {
        method: 'POST',
        headers: { Authorization: 'Bearer access-token' }
      })

      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('access_token')
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('refresh_token')
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('user_data')
    })

    it('should clear local state even if backend logout fails', async () => {
      mockLocalStorage.getItem.mockImplementation((key) => {
        switch (key) {
          case 'access_token': return 'access-token'
          default: return null
        }
      })

      mockFetch.mockRejectedValueOnce(new Error('Network error'))

      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()
      
      auth.initializeAuth()
      await auth.signOut()

      expect(auth.isAuthenticated.value).toBe(false)
      expect(auth.user.value).toBeNull()
    })
  })

  describe('Error Handling', () => {
    it('should clear errors', async () => {
      const { useAuth } = await import('~/composables/useAuth')
      const auth = useAuth()

      // Trigger an error
      await auth.signIn('google').catch(() => {})
      expect(auth.error.value).toBe('Only Discord authentication is supported')

      // Clear the error
      auth.clearError()
      expect(auth.error.value).toBeNull()
    })
  })
})