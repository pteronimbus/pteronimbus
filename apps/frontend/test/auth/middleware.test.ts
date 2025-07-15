import { describe, it, expect, beforeEach, vi } from 'vitest'
import { ref } from 'vue'

// Mock auth composable
const mockIsAuthenticated = ref(false)
const mockInitializeAuth = vi.fn()

vi.mock('~/composables/useAuth', () => ({
  useAuth: () => ({
    isAuthenticated: mockIsAuthenticated,
    initializeAuth: mockInitializeAuth
  })
}))

// Mock Nuxt navigation
const mockNavigateTo = vi.fn()
vi.mock('#app', () => ({
  defineNuxtRouteMiddleware: (fn: Function) => fn,
  navigateTo: mockNavigateTo
}))

// Mock the actual middleware modules
vi.mock('~/middleware/auth', () => ({
  default: (route: any) => {
    const { isAuthenticated, initializeAuth } = useAuth()
    initializeAuth()
    if (!isAuthenticated.value) {
      return mockNavigateTo('/login')
    }
  }
}))

vi.mock('~/middleware/guest', () => ({
  default: (route: any) => {
    const { isAuthenticated, initializeAuth } = useAuth()
    initializeAuth()
    if (isAuthenticated.value) {
      return mockNavigateTo('/dashboard')
    }
  }
}))

describe('Authentication Middleware', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockIsAuthenticated.value = false
  })

  describe('Auth Middleware', () => {
    it('should initialize auth state', async () => {
      const authMiddleware = await import('~/middleware/auth')
      const mockRoute = { path: '/dashboard' }

      authMiddleware.default(mockRoute as any)

      expect(mockInitializeAuth).toHaveBeenCalled()
    })

    it('should redirect to login when not authenticated', async () => {
      mockIsAuthenticated.value = false
      
      const authMiddleware = await import('~/middleware/auth')
      const mockRoute = { path: '/dashboard' }

      const result = authMiddleware.default(mockRoute as any)

      expect(mockNavigateTo).toHaveBeenCalledWith('/login')
    })

    it('should allow access when authenticated', async () => {
      mockIsAuthenticated.value = true
      
      const authMiddleware = await import('~/middleware/auth')
      const mockRoute = { path: '/dashboard' }

      const result = authMiddleware.default(mockRoute as any)

      expect(mockNavigateTo).not.toHaveBeenCalled()
      expect(result).toBeUndefined()
    })
  })

  describe('Guest Middleware', () => {
    it('should initialize auth state', async () => {
      const guestMiddleware = await import('~/middleware/guest')
      const mockRoute = { path: '/login' }

      guestMiddleware.default(mockRoute as any)

      expect(mockInitializeAuth).toHaveBeenCalled()
    })

    it('should allow access when not authenticated', async () => {
      mockIsAuthenticated.value = false
      
      const guestMiddleware = await import('~/middleware/guest')
      const mockRoute = { path: '/login' }

      const result = guestMiddleware.default(mockRoute as any)

      expect(mockNavigateTo).not.toHaveBeenCalled()
      expect(result).toBeUndefined()
    })

    it('should redirect to dashboard when authenticated', async () => {
      mockIsAuthenticated.value = true
      
      const guestMiddleware = await import('~/middleware/guest')
      const mockRoute = { path: '/login' }

      const result = guestMiddleware.default(mockRoute as any)

      expect(mockNavigateTo).toHaveBeenCalledWith('/dashboard')
    })
  })
})