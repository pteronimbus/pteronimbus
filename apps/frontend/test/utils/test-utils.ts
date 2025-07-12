// Test utilities for frontend testing
import { vi } from 'vitest'

// Mock router
export const mockRouter = {
  push: vi.fn(),
  go: vi.fn(),
  back: vi.fn(),
  forward: vi.fn(),
  currentRoute: {
    value: {
      path: '/',
      query: {},
      params: {}
    }
  }
}

// Mock i18n
export const mockI18n = {
  t: vi.fn((key: string) => key),
  locale: 'en',
  locales: ['en'],
  setLocale: vi.fn()
}

// Mock data generators
export const createMockServer = (overrides = {}) => ({
  id: 1,
  name: 'Test Server',
  game: 'Minecraft',
  status: 'online',
  players: '5/10',
  ip: '192.168.1.100',
  port: 25565,
  version: '1.20.4',
  uptime: '2h 30m',
  cpu: 45,
  memory: 60,
  createdAt: '2024-01-15',
  ...overrides
})

export const createMockUser = (overrides = {}) => ({
  id: 1,
  name: 'Test User',
  email: 'test@example.com',
  role: 'user',
  status: 'online',
  lastSeen: '2 minutes ago',
  serversAccess: 2,
  avatar: null,
  ...overrides
})

export const createMockPlayer = (overrides = {}) => ({
  id: 1,
  name: 'TestPlayer',
  server: 'Test Server',
  status: 'online',
  playtime: '2h 30m',
  lastSeen: '1 minute ago',
  avatar: null,
  ...overrides
})

export const createMockAlert = (overrides = {}) => ({
  id: 1,
  title: 'Test Alert',
  message: 'Test alert message',
  severity: 'warning',
  server: 'Test Server',
  timestamp: '2024-01-15T10:30:00Z',
  status: 'active',
  ...overrides
})

// Helper to wait for next tick
export const nextTick = () => new Promise(resolve => setTimeout(resolve, 0))

// Mock component resolver
export const mockResolveComponent = (name: string) => {
  const componentMap: Record<string, any> = {
    'UIcon': 'div',
    'UBadge': 'span',
    'UButton': 'button',
    'UCard': 'div',
    'UDropdownMenu': 'div',
    'UAvatar': 'div',
    'UInput': 'input',
    'USelect': 'select',
    'UTable': 'table',
    'UContainer': 'div',
    'UApp': 'div'
  }
  return componentMap[name] || 'div'
} 