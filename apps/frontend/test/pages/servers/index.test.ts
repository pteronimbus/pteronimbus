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
      expect(stats.total).toBe(vm.servers.length)
      expect(stats.online).toBe(vm.servers.filter((s: any) => s.status === 'online').length)
      expect(stats.offline).toBe(vm.servers.filter((s: any) => s.status === 'offline').length)
      expect(stats.error).toBe(vm.servers.filter((s: any) => s.status === 'error').length)
    })

    it('should filter servers by search query', async () => {
      const vm = wrapper.vm as any
      
      // Modify searchQuery directly through the component's ref
      vm.searchQuery = 'minecraft'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredServers
      expect(filtered.length).toBe(1)
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

    it('should generate game options from servers', () => {
      const vm = wrapper.vm as any
      const gameOptions = vm.gameOptions
      
      expect(gameOptions).toBeDefined()
      expect(gameOptions.length).toBeGreaterThan(1) // at least 'all' + unique games
      expect(gameOptions[0].value).toBe('all')
    })
  })

  describe('Helper Functions', () => {
    it('should return correct status colors', () => {
      const vm = wrapper.vm as any
      
      expect(vm.getStatusColor('online')).toBe('success')
      expect(vm.getStatusColor('offline')).toBe('error')
      expect(vm.getStatusColor('starting')).toBe('warning')
      expect(vm.getStatusColor('stopping')).toBe('warning')
      expect(vm.getStatusColor('error')).toBe('error')
      expect(vm.getStatusColor('unknown')).toBe('neutral')
    })

    it('should return correct performance colors', () => {
      const vm = wrapper.vm as any
      
      expect(vm.getPerformanceColor(30, 40)).toBe('success')
      expect(vm.getPerformanceColor(70, 50)).toBe('warning')
      expect(vm.getPerformanceColor(85, 70)).toBe('error')
      expect(vm.getPerformanceColor(45, 85)).toBe('error') // memory dominates
    })

    it('should generate action items correctly', () => {
      const vm = wrapper.vm as any
      const server = vm.servers[0] // online server
      const actions = vm.getActionItems(server)
      
      expect(actions).toBeDefined()
      expect(Array.isArray(actions)).toBe(true)
      expect(actions.length).toBeGreaterThan(0)
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
      const initialStats = { ...vm.serverStats }
      
      // Change a server status
      const server = vm.servers[0]
      const originalStatus = server.status
      server.status = 'error'
      
      await wrapper.vm.$nextTick()
      
      const newStats = vm.serverStats
      expect(newStats.total).toBe(initialStats.total) // total unchanged
      expect(newStats.error).toBe(initialStats.error + 1) // error count increased
      
      // Restore original status
      server.status = originalStatus
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
      await wrapper.vm.$nextTick()
      const onlineCount = vm.filteredServers.length
      expect(onlineCount).toBeLessThanOrEqual(initialCount)
      
      vm.selectedStatus = 'all'
      await wrapper.vm.$nextTick()
      expect(vm.filteredServers.length).toBe(initialCount)
    })
  })
}) 