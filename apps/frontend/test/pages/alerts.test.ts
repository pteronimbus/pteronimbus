import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import AlertsPage from '~/pages/alerts.vue'

// Mock Nuxt components and composables
vi.mock('#app', () => ({
  definePageMeta: vi.fn(),
  useI18n: () => ({
    t: (key: string) => key
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

describe('Alerts Page', () => {
  let wrapper: any

  beforeEach(() => {
    wrapper = mount(AlertsPage, {
      global: {
        plugins: [i18n],
        stubs: {
          UCard: true,
          UButton: true,
          UIcon: true,
          UTable: true,
          UInput: true,
          USelect: true,
          UBadge: true
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

    it('should have alerts data', () => {
      const vm = wrapper.vm as any
      expect(vm.alerts).toBeDefined()
      expect(Array.isArray(vm.alerts)).toBe(true)
      expect(vm.alerts.length).toBeGreaterThan(0)
    })

    it('should have initial filter states', () => {
      const vm = wrapper.vm as any
      expect(vm.searchQuery).toBe('')
      expect(vm.selectedSeverity).toBe('all')
      expect(vm.selectedStatus).toBe('all')
    })
  })

  describe('Computed Properties', () => {
    it('should calculate alert statistics correctly', () => {
      const vm = wrapper.vm as any
      const stats = vm.alertStats
      
      expect(stats).toBeDefined()
      expect(stats.total).toBe(vm.alerts.length)
      expect(stats.active).toBe(vm.alerts.filter((a: any) => !a.acknowledged).length)
      expect(stats.critical).toBe(vm.alerts.filter((a: any) => a.severity === 'critical').length)
      expect(stats.warnings).toBe(vm.alerts.filter((a: any) => a.severity === 'warning').length)
    })

    it('should filter alerts by search query - title', async () => {
      const vm = wrapper.vm as any
      
      vm.searchQuery = 'cpu'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredAlerts
      expect(filtered.length).toBe(1)
      expect(filtered[0].title.toLowerCase()).toContain('cpu')
    })

    it('should filter alerts by search query - message', async () => {
      const vm = wrapper.vm as any
      
      vm.searchQuery = 'exceeded'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredAlerts
      expect(filtered.length).toBeGreaterThan(0)
      expect(filtered[0].message.toLowerCase()).toContain('exceeded')
    })

    it('should filter alerts by search query - server', async () => {
      const vm = wrapper.vm as any
      
      vm.searchQuery = 'minecraft'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredAlerts
      expect(filtered.length).toBeGreaterThan(0)
      expect(filtered[0].server.toLowerCase()).toContain('minecraft')
    })

    it('should filter alerts by severity', async () => {
      const vm = wrapper.vm as any
      
      vm.selectedSeverity = 'warning'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredAlerts
      expect(filtered.every((a: any) => a.severity === 'warning')).toBe(true)
    })

    it('should filter alerts by status - active', async () => {
      const vm = wrapper.vm as any
      
      vm.selectedStatus = 'active'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredAlerts
      expect(filtered.every((a: any) => !a.acknowledged)).toBe(true)
    })

    it('should filter alerts by status - acknowledged', async () => {
      const vm = wrapper.vm as any
      
      vm.selectedStatus = 'acknowledged'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredAlerts
      expect(filtered.every((a: any) => a.acknowledged)).toBe(true)
    })

    it('should apply multiple filters simultaneously', async () => {
      const vm = wrapper.vm as any
      
      vm.searchQuery = 'minecraft'
      vm.selectedSeverity = 'warning'
      vm.selectedStatus = 'active'
      await wrapper.vm.$nextTick()
      
      const filtered = vm.filteredAlerts
      if (filtered.length > 0) {
        expect(filtered[0].server.toLowerCase()).toContain('minecraft')
        expect(filtered[0].severity).toBe('warning')
        expect(filtered[0].acknowledged).toBe(false)
      }
    })
  })

  describe('Helper Functions', () => {
    it('should return correct severity colors', () => {
      const vm = wrapper.vm as any
      
      expect(vm.getSeverityColor('critical')).toBe('error')
      expect(vm.getSeverityColor('error')).toBe('error')
      expect(vm.getSeverityColor('warning')).toBe('warning')
      expect(vm.getSeverityColor('info')).toBe('primary')
      expect(vm.getSeverityColor('unknown')).toBe('neutral')
    })
  })

  describe('Alert Actions', () => {
    it('should acknowledge an alert', async () => {
      const vm = wrapper.vm as any
      const activeAlert = vm.alerts.find((a: any) => !a.acknowledged)
      
      if (activeAlert) {
        const originalAcknowledged = activeAlert.acknowledged
        vm.acknowledgeAlert(activeAlert.id)
        expect(activeAlert.acknowledged).toBe(true)
        expect(activeAlert.acknowledged).not.toBe(originalAcknowledged)
      }
    })

    it('should not acknowledge an already acknowledged alert', async () => {
      const vm = wrapper.vm as any
      const acknowledgedAlert = vm.alerts.find((a: any) => a.acknowledged)
      
      if (acknowledgedAlert) {
        const beforeAcknowledged = acknowledgedAlert.acknowledged
        vm.acknowledgeAlert(acknowledgedAlert.id)
        expect(acknowledgedAlert.acknowledged).toBe(beforeAcknowledged) // should remain the same
      }
    })

    it('should dismiss an alert', async () => {
      const vm = wrapper.vm as any
      const initialCount = vm.alerts.length
      const alertToDismiss = vm.alerts[0]
      
      vm.dismissAlert(alertToDismiss.id)
      expect(vm.alerts.length).toBe(initialCount - 1)
      expect(vm.alerts.find((a: any) => a.id === alertToDismiss.id)).toBeUndefined()
    })

    it('should handle non-existent alert gracefully', async () => {
      const vm = wrapper.vm as any
      const initialCount = vm.alerts.length
      
      vm.acknowledgeAlert(999) // non-existent ID
      vm.dismissAlert(999) // non-existent ID
      
      expect(vm.alerts.length).toBe(initialCount) // should remain unchanged
    })
  })

  describe('Alert Data Structure', () => {
    it('should have expected alert properties', () => {
      const vm = wrapper.vm as any
      const alert = vm.alerts[0]
      
      expect(alert).toHaveProperty('id')
      expect(alert).toHaveProperty('type')
      expect(alert).toHaveProperty('title')
      expect(alert).toHaveProperty('message')
      expect(alert).toHaveProperty('severity')
      expect(alert).toHaveProperty('timestamp')
      expect(alert).toHaveProperty('icon')
      expect(alert).toHaveProperty('server')
      expect(alert).toHaveProperty('acknowledged')
    })

    it('should have valid alert severities', () => {
      const vm = wrapper.vm as any
      const validSeverities = ['critical', 'error', 'warning', 'info']
      
      vm.alerts.forEach((alert: any) => {
        expect(validSeverities).toContain(alert.severity)
      })
    })

    it('should have boolean acknowledged field', () => {
      const vm = wrapper.vm as any
      
      vm.alerts.forEach((alert: any) => {
        expect(typeof alert.acknowledged).toBe('boolean')
      })
    })
  })

  describe('Reactive Updates', () => {
    it('should update stats when alert is acknowledged', async () => {
      const vm = wrapper.vm as any
      const initialStats = { ...vm.alertStats }
      const activeAlert = vm.alerts.find((a: any) => !a.acknowledged)
      
      if (activeAlert) {
        vm.acknowledgeAlert(activeAlert.id)
        await wrapper.vm.$nextTick()
        
        const newStats = vm.alertStats
        expect(newStats.active).toBe(initialStats.active - 1)
        expect(newStats.total).toBe(initialStats.total) // total unchanged
      }
    })

    it('should update stats when alert is dismissed', async () => {
      const vm = wrapper.vm as any
      const initialStats = { ...vm.alertStats }
      const alertToDismiss = vm.alerts[0]
      
      vm.dismissAlert(alertToDismiss.id)
      await wrapper.vm.$nextTick()
      
      const newStats = vm.alertStats
      expect(newStats.total).toBe(initialStats.total - 1)
    })

    it('should update filtered results when search changes', async () => {
      const vm = wrapper.vm as any
      const initialCount = vm.filteredAlerts.length
      
      vm.searchQuery = 'nonexistent'
      await wrapper.vm.$nextTick()
      expect(vm.filteredAlerts.length).toBe(0)
      
      vm.searchQuery = ''
      await wrapper.vm.$nextTick()
      expect(vm.filteredAlerts.length).toBe(initialCount)
    })

    it('should update filtered results when severity filter changes', async () => {
      const vm = wrapper.vm as any
      const initialCount = vm.filteredAlerts.length
      
      vm.selectedSeverity = 'critical'
      await wrapper.vm.$nextTick()
      const criticalCount = vm.filteredAlerts.length
      expect(criticalCount).toBeLessThanOrEqual(initialCount)
      
      vm.selectedSeverity = 'all'
      await wrapper.vm.$nextTick()
      expect(vm.filteredAlerts.length).toBe(initialCount)
    })

    it('should update filtered results when status filter changes', async () => {
      const vm = wrapper.vm as any
      const initialCount = vm.filteredAlerts.length
      
      vm.selectedStatus = 'active'
      await wrapper.vm.$nextTick()
      const activeCount = vm.filteredAlerts.length
      expect(activeCount).toBeLessThanOrEqual(initialCount)
      
      vm.selectedStatus = 'all'
      await wrapper.vm.$nextTick()
      expect(vm.filteredAlerts.length).toBe(initialCount)
    })
  })

  describe('Options Arrays', () => {
    it('should have correct severity options', () => {
      const vm = wrapper.vm as any
      
      expect(vm.severityOptions).toBeDefined()
      expect(Array.isArray(vm.severityOptions)).toBe(true)
      expect(vm.severityOptions.length).toBeGreaterThan(0)
      expect(vm.severityOptions[0].value).toBe('all')
    })

    it('should have correct status options', () => {
      const vm = wrapper.vm as any
      
      expect(vm.statusOptions).toBeDefined()
      expect(Array.isArray(vm.statusOptions)).toBe(true)
      expect(vm.statusOptions.length).toBeGreaterThan(0)
      expect(vm.statusOptions[0].value).toBe('all')
    })
  })
}) 