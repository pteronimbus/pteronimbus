import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { createRouter, createWebHistory } from 'vue-router'
import Players from '~/pages/players.vue'

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
    en: {}
  }
})

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/players', component: Players }
  ]
})

describe('Players Page', () => {
  let wrapper: any

  beforeEach(() => {
    wrapper = mount(Players, {
      global: {
        plugins: [i18n, router],
        stubs: {
          UCard: true,
          UIcon: true,
          UTable: true,
          UInput: true,
          USelect: true,
          UBadge: true,
          UAvatar: true
        }
      }
    })
  })

  afterEach(() => {
    wrapper.unmount()
  })

  describe('Component Mounting', () => {
    it('should mount successfully', () => {
      expect(wrapper.exists()).toBe(true)
    })

    it('should have players data', () => {
      const vm = wrapper.vm as any
      expect(vm.players).toBeDefined()
      expect(Array.isArray(vm.players)).toBe(true)
      expect(vm.players.length).toBeGreaterThan(0)
    })

    it('should have initial filter states', () => {
      const vm = wrapper.vm as any
      expect(vm.searchQuery).toBe('')
      expect(vm.selectedStatus).toBe('all')
    })
  })

  describe('Player Data Structure', () => {
    it('should have expected player properties', () => {
      const vm = wrapper.vm as any
      const player = vm.players[0]
      
      expect(player).toHaveProperty('id')
      expect(player).toHaveProperty('name')
      expect(player).toHaveProperty('server')
      expect(player).toHaveProperty('status')
      expect(player).toHaveProperty('playtime')
      expect(player).toHaveProperty('lastSeen')
      expect(player).toHaveProperty('avatar')
    })

    it('should have valid player statuses', () => {
      const vm = wrapper.vm as any
      const validStatuses = ['online', 'offline']
      
      vm.players.forEach((player: any) => {
        expect(validStatuses).toContain(player.status)
      })
    })

    it('should have players from different servers', () => {
      const vm = wrapper.vm as any
      const servers = vm.players.map((p: any) => p.server)
      const uniqueServers = [...new Set(servers)]
      
      expect(uniqueServers.length).toBeGreaterThan(1)
    })
  })

  describe('Component Structure', () => {
    it('should use StatusBadge component for status display', () => {
      const vm = wrapper.vm as any
      const columns = vm.columns
      
      const statusColumn = columns.find((col: any) => col.accessorKey === 'status')
      expect(statusColumn).toBeDefined()
      expect(statusColumn.cell).toBeDefined()
      expect(typeof statusColumn.cell).toBe('function')
    })
  })

  describe('Player Statistics', () => {
    it('should calculate correct player statistics', () => {
      const vm = wrapper.vm as any
      const stats = vm.playerStats
      
      expect(stats).toBeDefined()
      expect(Array.isArray(stats)).toBe(true)
      expect(stats.length).toBe(3) // total, online, offline
      
      const totalStat = stats.find((s: any) => s.key === 'total')
      const onlineStat = stats.find((s: any) => s.key === 'online')
      const offlineStat = stats.find((s: any) => s.key === 'offline')
      
      expect(totalStat).toBeDefined()
      expect(onlineStat).toBeDefined()
      expect(offlineStat).toBeDefined()
      expect(parseInt(totalStat.value)).toBe(vm.players.length)
    })

    it('should have total equal to online plus offline', () => {
      const vm = wrapper.vm as any
      const stats = vm.playerStats
      
      const totalStat = stats.find((s: any) => s.key === 'total')
      const onlineStat = stats.find((s: any) => s.key === 'online')
      const offlineStat = stats.find((s: any) => s.key === 'offline')
      
      expect(parseInt(totalStat.value)).toBe(parseInt(onlineStat.value) + parseInt(offlineStat.value))
    })

    it('should update statistics when player status changes', async () => {
      const vm = wrapper.vm as any
      const initialStats = vm.playerStats
      const initialOnlineCount = parseInt(initialStats.find((s: any) => s.key === 'online').value)
      const initialOfflineCount = parseInt(initialStats.find((s: any) => s.key === 'offline').value)
      
      // Change a player status
      const onlinePlayer = vm.players.find((p: any) => p.status === 'online')
      if (onlinePlayer) {
        onlinePlayer.status = 'offline'
        await wrapper.vm.$nextTick()
        
        const newStats = vm.playerStats
        const newOnlineCount = parseInt(newStats.find((s: any) => s.key === 'online').value)
        const newOfflineCount = parseInt(newStats.find((s: any) => s.key === 'offline').value)
        const newTotalCount = parseInt(newStats.find((s: any) => s.key === 'total').value)
        
        expect(newOnlineCount).toBe(initialOnlineCount - 1)
        expect(newOfflineCount).toBe(initialOfflineCount + 1)
        expect(newTotalCount).toBe(vm.players.length) // total unchanged
      }
    })
  })

  describe('Player Filtering', () => {
    it('should filter players by search query - name', async () => {
      const vm = wrapper.vm as any
      
      vm.searchQuery = 'PlayerOne'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredPlayers
      expect(filtered.length).toBe(1)
      expect(filtered[0].name).toBe('PlayerOne')
    })

    it('should filter players by search query - server', async () => {
      const vm = wrapper.vm as any
      
      vm.searchQuery = 'Minecraft'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredPlayers
      expect(filtered.length).toBeGreaterThan(0)
      expect(filtered.every((p: any) => p.server.toLowerCase().includes('minecraft'))).toBe(true)
    })

    it('should filter players by status', async () => {
      const vm = wrapper.vm as any
      
      vm.selectedStatus = 'online'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredPlayers
      expect(filtered.every((p: any) => p.status === 'online')).toBe(true)
    })

    it('should apply multiple filters simultaneously', async () => {
      const vm = wrapper.vm as any
      
      vm.searchQuery = 'Minecraft'
      vm.selectedStatus = 'online'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredPlayers
      filtered.forEach((player: any) => {
        expect(player.server.toLowerCase()).toContain('minecraft')
        expect(player.status).toBe('online')
      })
    })

    it('should handle case-insensitive search', async () => {
      const vm = wrapper.vm as any
      
      vm.searchQuery = 'PLAYERONE'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredPlayers
      expect(filtered.length).toBe(1)
      expect(filtered[0].name).toBe('PlayerOne')
    })

    it('should return empty array when no matches', async () => {
      const vm = wrapper.vm as any
      
      vm.searchQuery = 'nonexistent'
      await wrapper.vm.$nextTick()
      
      expect(vm.filteredPlayers.length).toBe(0)
    })
  })

  describe('Filter Configuration', () => {
    it('should have filters with status options', () => {
      const vm = wrapper.vm as any
      const filters = vm.filters
      
      expect(filters).toBeDefined()
      expect(Array.isArray(filters)).toBe(true)
      expect(filters.length).toBe(1) // only status filter
      
      const statusFilter = filters.find((f: any) => f.key === 'status')
      expect(statusFilter).toBeDefined()
      expect(Array.isArray(statusFilter.options)).toBe(true)
      expect(statusFilter.options.length).toBe(3) // all, online, offline
      expect(statusFilter.options[0].value).toBe('all')
    })

    it('should include all valid status values', () => {
      const vm = wrapper.vm as any
      const filters = vm.filters
      const statusFilter = filters.find((f: any) => f.key === 'status')
      const statusValues = statusFilter.options.map((opt: any) => opt.value)
      
      expect(statusValues).toContain('all')
      expect(statusValues).toContain('online')
      expect(statusValues).toContain('offline')
    })
  })

  describe('Table Columns Configuration', () => {
    it('should have correct column structure', () => {
      const vm = wrapper.vm as any
      
      expect(vm.columns).toBeDefined()
      expect(Array.isArray(vm.columns)).toBe(true)
      expect(vm.columns.length).toBe(5)
    })

    it('should have all required columns', () => {
      const vm = wrapper.vm as any
      const columnKeys = vm.columns.map((col: any) => col.accessorKey)
      
      expect(columnKeys).toContain('name')
      expect(columnKeys).toContain('server')
      expect(columnKeys).toContain('status')
      expect(columnKeys).toContain('playtime')
      expect(columnKeys).toContain('lastSeen')
    })

    it('should have cell renderers for each column', () => {
      const vm = wrapper.vm as any
      
      vm.columns.forEach((column: any) => {
        expect(column).toHaveProperty('cell')
        expect(typeof column.cell).toBe('function')
      })
    })
  })

  describe('Reactive Updates', () => {
    it('should update filtered results when search changes', async () => {
      const vm = wrapper.vm as any
      const initialCount = vm.filteredPlayers.length
      
      vm.searchQuery = 'nonexistent'
      await wrapper.vm.$nextTick()
      expect(vm.filteredPlayers.length).toBe(0)
      
      vm.searchQuery = ''
      await wrapper.vm.$nextTick()
      expect(vm.filteredPlayers.length).toBe(initialCount)
    })

    it('should update filtered results when status filter changes', async () => {
      const vm = wrapper.vm as any
      const initialCount = vm.filteredPlayers.length
      
      vm.selectedStatus = 'online'
      await wrapper.vm.$nextTick()
      const onlineCount = vm.filteredPlayers.length
      expect(onlineCount).toBeLessThanOrEqual(initialCount)
      
      vm.selectedStatus = 'all'
      await wrapper.vm.$nextTick()
      expect(vm.filteredPlayers.length).toBe(initialCount)
    })
  })
}) 