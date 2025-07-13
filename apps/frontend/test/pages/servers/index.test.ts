import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { createRouter, createWebHistory } from 'vue-router'
import ServersIndex from '~/pages/servers/index.vue'

// Mock Nuxt components and composables
vi.mock('#app', () => ({
  definePageMeta: vi.fn(),
  useI18n: () => ({
    t: (key: string) => key
  }),
  useRouter: () => ({
    push: vi.fn()
  })
}))

vi.mock('vue', async () => {
  const actual = await vi.importActual('vue')
  return {
    ...actual,
    resolveComponent: (name: string) => name
  }
})

const i18n = createI18n({
  locale: 'en',
  messages: {
    en: {
      servers: {
        title: 'Servers',
        createServer: 'Create Server',
        status: {
          online: 'Online',
          offline: 'Offline',
          starting: 'Starting',
          stopping: 'Stopping',
          error: 'Error'
        },
        actions: {
          viewDetails: 'View Details',
          console: 'Console',
          start: 'Start',
          stop: 'Stop',
          restart: 'Restart',
          edit: 'Edit',
          delete: 'Delete'
        }
      },
      common: {
        search: 'Search'
      }
    }
  }
})

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/servers', component: ServersIndex },
    { path: '/servers/:id', component: {} }
  ]
})

describe('Servers Index Page', () => {
  let wrapper: any

  beforeEach(() => {
    vi.useFakeTimers()
    wrapper = mount(ServersIndex, {
      global: {
        plugins: [i18n, router],
        stubs: {
          UCard: true,
          UButton: true,
          UIcon: true,
          UTable: true,
          UInput: true,
          USelect: true,
          UBadge: true,
          UDropdownMenu: true
        }
      }
    })
  })

  afterEach(() => {
    vi.useRealTimers()
    wrapper.unmount()
  })

  describe('Component Mounting', () => {
    it('should mount successfully', () => {
      expect(wrapper.exists()).toBe(true)
    })

    it('should have servers data', () => {
      const vm = wrapper.vm as any
      expect(vm.servers).toBeDefined()
      expect(Array.isArray(vm.servers)).toBe(true)
      expect(vm.servers.length).toBeGreaterThan(0)
    })

    it('should have initial filter states', () => {
      const vm = wrapper.vm as any
      expect(vm.searchQuery).toBe('')
      expect(vm.selectedStatus).toBe('all')
      expect(vm.selectedGame).toBe('all')
    })
  })

  describe('Computed Properties', () => {
    it('should calculate server statistics correctly', () => {
      const vm = wrapper.vm as any
      const stats = vm.serverStats
      
      expect(stats).toBeDefined()
      expect(Array.isArray(stats)).toBe(true)
      expect(stats.length).toBe(4) // total, online, offline, error
      
      const totalStat = stats.find((s: any) => s.key === 'total')
      const onlineStat = stats.find((s: any) => s.key === 'online')
      
      expect(totalStat).toBeDefined()
      expect(onlineStat).toBeDefined()
      expect(parseInt(totalStat.value)).toBe(vm.servers.length)
    })

    it('should filter servers by search query', async () => {
      const vm = wrapper.vm as any
      
      vm.searchQuery = 'minecraft'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredServers
      expect(filtered.length).toBeGreaterThan(0)
      expect(filtered[0].name.toLowerCase()).toContain('minecraft')
    })

    it('should filter servers by status', async () => {
      const vm = wrapper.vm as any
      
      vm.selectedStatus = 'online'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredServers
      expect(filtered.every((s: any) => s.status === 'online')).toBe(true)
    })

    it('should filter servers by game', async () => {
      const vm = wrapper.vm as any
      
      vm.selectedGame = 'Minecraft'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredServers
      expect(filtered.every((s: any) => s.game === 'Minecraft')).toBe(true)
    })

    it('should have game options in filters configuration', () => {
      const vm = wrapper.vm as any
      const filters = vm.filters
      
      const gameFilter = filters.find((f: any) => f.key === 'game')
      expect(gameFilter).toBeDefined()
      expect(Array.isArray(gameFilter.options)).toBe(true)
      expect(gameFilter.options.length).toBeGreaterThan(1) // at least 'all' + unique games
      expect(gameFilter.options[0].value).toBe('all')
    })
  })

  describe('Helper Functions', () => {
    it('should generate correct filter configurations', () => {
      const vm = wrapper.vm as any
      const filters = vm.filters
      
      expect(filters).toBeDefined()
      expect(Array.isArray(filters)).toBe(true)
      expect(filters.length).toBe(2) // status and game filters
    })

    it('should return correct performance colors', () => {
      const vm = wrapper.vm as any
      
      expect(vm.getPerformanceColor(50, 50)).toBe('success') // normal usage
      expect(vm.getPerformanceColor(70, 60)).toBe('warning') // high usage
      expect(vm.getPerformanceColor(90, 85)).toBe('error') // very high usage
    })

    it('should generate action items correctly', () => {
      const vm = wrapper.vm as any
      const mockServer = { id: 1, name: 'Test Server', status: 'online' }
      const actions = vm.getActionItems(mockServer)
      
      expect(Array.isArray(actions)).toBe(true)
      expect(actions.length).toBeGreaterThan(0)
      
      // Check structure of action groups
      actions.forEach((group: any) => {
        expect(Array.isArray(group)).toBe(true)
        group.forEach((action: any) => {
          expect(action).toHaveProperty('label')
          expect(action).toHaveProperty('icon')
          expect(action).toHaveProperty('click')
        })
      })
    })
  })

  describe('Server Actions', () => {
    it('should toggle server from online to offline', async () => {
      const vm = wrapper.vm as any
      const onlineServer = vm.servers.find((s: any) => s.status === 'online')
      
      if (onlineServer) {
        vm.toggleServer(onlineServer)
        expect(onlineServer.status).toBe('stopping')
        
        // Fast-forward through the timeout
        vi.advanceTimersByTime(2000)
        await wrapper.vm.$nextTick()
        
        expect(onlineServer.status).toBe('offline')
        expect(onlineServer.cpu).toBe(0)
        expect(onlineServer.memory).toBe(0)
      }
    })

    it('should toggle server from offline to online', async () => {
      const vm = wrapper.vm as any
      const offlineServer = vm.servers.find((s: any) => s.status === 'offline')
      
      if (offlineServer) {
        vm.toggleServer(offlineServer)
        expect(offlineServer.status).toBe('starting')
        
        // Fast-forward through the timeout
        vi.advanceTimersByTime(3000)
        await wrapper.vm.$nextTick()
        
        expect(offlineServer.status).toBe('online')
        expect(offlineServer.cpu).toBeGreaterThan(0)
        expect(offlineServer.memory).toBeGreaterThan(0)
      }
    })

    it('should restart server through complete cycle', async () => {
      const vm = wrapper.vm as any
      const onlineServer = vm.servers.find((s: any) => s.status === 'online')
      
      if (onlineServer) {
        vm.restartServer(onlineServer)
        expect(onlineServer.status).toBe('stopping')
        
        // Fast-forward to stopping -> starting transition
        vi.advanceTimersByTime(2000)
        await wrapper.vm.$nextTick()
        expect(onlineServer.status).toBe('starting')
        
        // Fast-forward to starting -> online transition
        vi.advanceTimersByTime(3000)
        await wrapper.vm.$nextTick()
        expect(onlineServer.status).toBe('online')
      }
    })
  })

  describe('Reactive Updates', () => {
    it('should update stats when server status changes', async () => {
      const vm = wrapper.vm as any
      const initialStats = vm.serverStats
      const initialErrorCount = parseInt(initialStats.find((s: any) => s.key === 'error').value)
      const server = vm.servers.find((s: any) => s.status !== 'error')
      
      if (server) {
        const originalStatus = server.status
        server.status = 'error'
        await wrapper.vm.$nextTick()
        
        const newStats = vm.serverStats
        const newErrorCount = parseInt(newStats.find((s: any) => s.key === 'error').value)
        expect(newErrorCount).toBe(initialErrorCount + 1)
        
        // Restore original status
        server.status = originalStatus
      }
    })

    it('should update filtered results when search changes', async () => {
      const vm = wrapper.vm as any
      const initialCount = vm.filteredServers.length
      
      vm.searchQuery = 'nonexistent'
      await wrapper.vm.$nextTick()
      expect(vm.filteredServers.length).toBe(0)
      
      vm.searchQuery = ''
      await wrapper.vm.$nextTick()
      expect(vm.filteredServers.length).toBe(initialCount)
    })

    it('should update filtered results when filters change', async () => {
      const vm = wrapper.vm as any
      const initialCount = vm.filteredServers.length
      
      vm.selectedStatus = 'online'
      vm.selectedGame = 'Minecraft'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredServers
      expect(filtered.length).toBeLessThanOrEqual(initialCount)
      
      // Reset filters
      vm.selectedStatus = 'all'
      vm.selectedGame = 'all'
      await wrapper.vm.$nextTick()
      expect(vm.filteredServers.length).toBe(initialCount)
    })
  })
}) 