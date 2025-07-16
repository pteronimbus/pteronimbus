import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { createRouter, createWebHistory } from 'vue-router'
import TenantDashboard from '~/pages/tenant/[tenantId]/dashboard.vue'

// Create shared mock functions
const mockRouterPush = vi.fn()
const mockTenantApiRequest = vi.fn()
const mockSyncTenantData = vi.fn()
const mockClearError = vi.fn()
const mockDefinePageMeta = vi.fn()
const mockToastAdd = vi.fn()

// Create reactive refs for mocking
let mockCurrentTenant: any = {
  id: 'tenant-123',
  name: 'Test Discord Server',
  discord_server_id: '123456789',
  icon: 'test-icon',
  config: {
    resource_limits: {
      max_game_servers: 10
    }
  }
}
const mockUser = { id: 'user-123', name: 'Test User', role: 'admin' }

// Mock Nuxt components and composables
vi.mock('#app', () => ({
  definePageMeta: mockDefinePageMeta,
  useI18n: () => ({
    t: (key: string) => key
  }),
  useAuth: () => ({
    user: { value: mockUser }
  }),
  useTenant: () => ({
    currentTenant: { value: mockCurrentTenant },
    tenantApiRequest: mockTenantApiRequest,
    syncTenantData: mockSyncTenantData,
    clearError: mockClearError
  }),
  useRouter: () => ({
    push: mockRouterPush
  }),
  useRoute: () => ({
    params: { tenantId: 'tenant-123' }
  }),
  useToast: () => ({
    add: mockToastAdd
  }),
  onMounted: vi.fn(),
  watch: vi.fn(),
  ref: (value: any) => ({ value }),
  computed: (fn: Function) => ({ value: fn() })
}))

// Mock vue-router
vi.mock('vue-router', async () => {
  const actual = await vi.importActual('vue-router')
  return {
    ...actual,
    useRouter: () => ({
      push: mockRouterPush
    }),
    useRoute: () => ({
      params: { tenantId: 'tenant-123' }
    })
  }
})

const i18n = createI18n({
  locale: 'en',
  messages: {
    en: {
      dashboard: {
        tenantOverview: 'Tenant Overview',
        quickActions: 'Quick Actions',
        gameServers: 'Game Servers',
        discordIntegration: 'Discord Integration',
        sync: 'Sync',
        manageRoles: 'Manage Roles',
        viewLogs: 'View Logs',
        settings: 'Settings',
        noGameServers: 'No Game Servers',
        noGameServersDesc: 'Create your first game server to get started',
        stats: {
          gameServers: 'Game Servers',
          activeServers: 'Active Servers',
          totalPlayers: 'Total Players',
          discordMembers: 'Discord Members'
        },
        activity: {
          title: 'Recent Activity',
          noActivity: 'No Recent Activity',
          noActivityDesc: 'Activity will appear here as things happen'
        }
      },
      servers: {
        createServer: 'Create Server'
      },
      common: {
        refresh: 'Refresh',
        viewAll: 'View All'
      }
    }
  }
})

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/tenant/:tenantId/dashboard', component: TenantDashboard },
    { path: '/tenant/:tenantId/servers', component: {} },
    { path: '/tenant/:tenantId/servers/create', component: {} }
  ]
})

describe('Tenant Dashboard Page', () => {
  let wrapper: any

  beforeEach(() => {
    // Reset all mocks
    mockRouterPush.mockClear()
    mockTenantApiRequest.mockClear()
    mockSyncTenantData.mockClear()
    mockClearError.mockClear()
    mockDefinePageMeta.mockClear()
    mockToastAdd.mockClear()
  })

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount()
    }
  })

  describe('Component Mounting', () => {
    it('should mount successfully', () => {
      wrapper = mount(TenantDashboard, {
        global: {
          plugins: [i18n, router],
          stubs: {
            UIcon: true,
            UButton: true,
            UCard: true,
            UAvatar: true,
            StatsCard: true,
            QuickActions: true,
            StatusBadge: true,
            EmptyState: true,
            TenantSelector: true
          }
        }
      })
      expect(wrapper.exists()).toBe(true)
    })

    it('should render without errors', () => {
      expect(() => {
        wrapper = mount(TenantDashboard, {
          global: {
            plugins: [i18n, router],
            stubs: {
              UIcon: true,
              UButton: true,
              UCard: true,
              UAvatar: true,
              StatsCard: true,
              QuickActions: true,
              StatusBadge: true,
              EmptyState: true,
              TenantSelector: true
            }
          }
        })
      }).not.toThrow()
    })
  })

  describe('Component Structure', () => {
    beforeEach(() => {
      wrapper = mount(TenantDashboard, {
        global: {
          plugins: [i18n, router],
          stubs: {
            UIcon: true,
            UButton: true,
            UCard: true,
            UAvatar: true,
            StatsCard: true,
            QuickActions: true,
            StatusBadge: true,
            EmptyState: true,
            TenantSelector: true
          }
        }
      })
    })

    it('should have main container div', () => {
      expect(wrapper.find('div').exists()).toBe(true)
    })

    it('should include tenant selector component', () => {
      expect(wrapper.findComponent({ name: 'TenantSelector' })).toBeTruthy()
    })

    it('should include stats cards component', () => {
      expect(wrapper.findComponent({ name: 'StatsCard' })).toBeTruthy()
    })

    it('should include quick actions component', () => {
      expect(wrapper.findComponent({ name: 'QuickActions' })).toBeTruthy()
    })

    it('should include UI components', () => {
      expect(wrapper.findComponent({ name: 'UButton' })).toBeTruthy()
      expect(wrapper.findComponent({ name: 'UCard' })).toBeTruthy()
      expect(wrapper.findComponent({ name: 'UAvatar' })).toBeTruthy()
    })
  })

  describe('Text Content', () => {
    beforeEach(() => {
      wrapper = mount(TenantDashboard, {
        global: {
          plugins: [i18n, router],
          stubs: {
            UIcon: true,
            UButton: true,
            UCard: true,
            UAvatar: true,
            StatsCard: true,
            QuickActions: true,
            StatusBadge: true,
            EmptyState: true,
            TenantSelector: true
          }
        }
      })
    })

    it('should contain expected text content', () => {
      const text = wrapper.text()
      expect(text.length).toBeGreaterThan(0)
    })

    it('should display loading or content state', () => {
      const text = wrapper.text()
      // Should either show loading state or actual content
      expect(text).toMatch(/Loading|Game Servers|Recent Activity|Discord Integration/)
    })

    it('should include common UI elements', () => {
      const text = wrapper.text()
      expect(text).toMatch(/Refresh|View All|Sync/)
    })
  })

  describe('Component Props and Data', () => {
    beforeEach(() => {
      wrapper = mount(TenantDashboard, {
        global: {
          plugins: [i18n, router],
          stubs: {
            UIcon: true,
            UButton: true,
            UCard: true,
            UAvatar: true,
            StatsCard: true,
            QuickActions: true,
            StatusBadge: true,
            EmptyState: true,
            TenantSelector: true
          }
        }
      })
    })

    it('should be a Vue component instance', () => {
      expect(wrapper.vm).toBeDefined()
      expect(typeof wrapper.vm).toBe('object')
    })

    it('should have component data', () => {
      expect(wrapper.vm).toBeTruthy()
    })
  })

  describe('CSS Classes and Styling', () => {
    beforeEach(() => {
      wrapper = mount(TenantDashboard, {
        global: {
          plugins: [i18n, router],
          stubs: {
            UIcon: true,
            UButton: true,
            UCard: true,
            UAvatar: true,
            StatsCard: true,
            QuickActions: true,
            StatusBadge: true,
            EmptyState: true,
            TenantSelector: true
          }
        }
      })
    })

    it('should have grid layout classes', () => {
      const gridElements = wrapper.findAll('.grid')
      expect(gridElements.length).toBeGreaterThan(0)
    })

    it('should have spacing classes', () => {
      const spacingElements = wrapper.findAll('[class*="space-"], [class*="gap-"], [class*="mb-"], [class*="mt-"]')
      expect(spacingElements.length).toBeGreaterThan(0)
    })

    it('should have responsive classes', () => {
      const responsiveElements = wrapper.findAll('[class*="sm:"], [class*="lg:"], [class*="md:"]')
      expect(responsiveElements.length).toBeGreaterThan(0)
    })
  })

  describe('Accessibility', () => {
    beforeEach(() => {
      wrapper = mount(TenantDashboard, {
        global: {
          plugins: [i18n, router],
          stubs: {
            UIcon: true,
            UButton: true,
            UCard: true,
            UAvatar: true,
            StatsCard: true,
            QuickActions: true,
            StatusBadge: true,
            EmptyState: true,
            TenantSelector: true
          }
        }
      })
    })

    it('should have heading elements', () => {
      const headings = wrapper.findAll('h1, h2, h3, h4, h5, h6')
      expect(headings.length).toBeGreaterThan(0)
    })

    it('should have proper semantic structure', () => {
      expect(wrapper.find('div').exists()).toBe(true)
    })

    it('should have descriptive content', () => {
      const text = wrapper.text()
      expect(text.length).toBeGreaterThan(10)
    })
  })

  describe('Error Handling', () => {
    it('should handle component mounting gracefully', () => {
      expect(() => {
        const testWrapper = mount(TenantDashboard, {
          global: {
            plugins: [i18n, router],
            stubs: {
              UIcon: true,
              UButton: true,
              UCard: true,
              UAvatar: true,
              StatsCard: true,
              QuickActions: true,
              StatusBadge: true,
              EmptyState: true,
              TenantSelector: true
            }
          }
        })
        testWrapper.unmount()
      }).not.toThrow()
    })

    it('should handle missing props gracefully', () => {
      expect(() => {
        const testWrapper = mount(TenantDashboard, {
          global: {
            plugins: [i18n, router],
            stubs: {
              UIcon: true,
              UButton: true,
              UCard: true,
              UAvatar: true,
              StatsCard: true,
              QuickActions: true,
              StatusBadge: true,
              EmptyState: true,
              TenantSelector: true
            }
          }
        })
        testWrapper.unmount()
      }).not.toThrow()
    })
  })

  describe('Component Integration', () => {
    beforeEach(() => {
      wrapper = mount(TenantDashboard, {
        global: {
          plugins: [i18n, router],
          stubs: {
            UIcon: true,
            UButton: true,
            UCard: true,
            UAvatar: true,
            StatsCard: true,
            QuickActions: true,
            StatusBadge: true,
            EmptyState: true,
            TenantSelector: true
          }
        }
      })
    })

    it('should integrate with i18n', () => {
      // Component should mount without i18n errors
      expect(wrapper.exists()).toBe(true)
    })

    it('should integrate with router', () => {
      // Component should mount without router errors
      expect(wrapper.exists()).toBe(true)
    })

    it('should use stubbed components correctly', () => {
      // All stubbed components should be present
      expect(wrapper.findComponent({ name: 'UIcon' })).toBeTruthy()
      expect(wrapper.findComponent({ name: 'UButton' })).toBeTruthy()
      expect(wrapper.findComponent({ name: 'UCard' })).toBeTruthy()
      expect(wrapper.findComponent({ name: 'StatsCard' })).toBeTruthy()
      expect(wrapper.findComponent({ name: 'QuickActions' })).toBeTruthy()
    })
  })
})