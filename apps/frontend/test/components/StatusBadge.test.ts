import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import StatusBadge from '~/components/StatusBadge.vue'
import { defineComponent } from 'vue'
import type { App } from 'vue'

// Create a proper i18n instance for testing
const i18n = createI18n({
  legacy: false, // Use Composition API
  locale: 'en',
  messages: {
    en: {
      servers: {
        status: {
          online: 'Online',
          offline: 'Offline',
          starting: 'Starting',
          stopping: 'Stopping',
          error: 'Error'
        }
      },
      users: {
        status: {
          online: 'Online',
          offline: 'Offline',
          banned: 'Banned',
          suspended: 'Suspended',
          active: 'Active'
        }
      },
      alerts: {
        severity: {
          critical: 'Critical',
          error: 'Error',
          warning: 'Warning',
          info: 'Info'
        }
      }
    }
  }
})

// Mock Nuxt auto-imports
vi.mock('#imports', () => ({
  useI18n: () => i18n.global
}))

// UBadge stub plugin
const UBadgeStubPlugin = {
  install(app: App) {
    app.component('UBadge', defineComponent({
      props: ['color', 'variant', 'class'],
      template: '<span data-testid="badge" :color="color" :variant="variant" :class="class"><slot /></span>'
    }))
  }
}

describe('StatusBadge Component', () => {
  const createWrapper = (props = {}) => {
    return mount(StatusBadge, {
      props: {
        status: 'online',
        type: 'server',
        ...props
      },
      global: {
        plugins: [i18n],
        stubs: {
          UBadge: true
        }
      }
    })
  }

  beforeEach(() => {
    // Reset any mocks between tests
    vi.clearAllMocks()
  })

  it('displays translated text for server status', () => {
    const wrapper = createWrapper({
      status: 'online',
      type: 'server'
    })
    
    expect(wrapper.text()).toContain('Online')
  })

  it('displays custom label when provided', () => {
    const wrapper = createWrapper({
      status: 'online',
      type: 'server',
      label: 'Custom Status Label'
    })
    
    expect(wrapper.text()).toContain('Custom Status Label')
  })

  it('handles different server statuses', () => {
    const wrapper = createWrapper({
      status: 'offline',
      type: 'server'
    })
    
    expect(wrapper.text()).toContain('Offline')
  })

  it('handles user statuses with translations', () => {
    const wrapper = createWrapper({
      status: 'banned',
      type: 'user'
    })
    
    expect(wrapper.text()).toContain('Banned')
  })

  it('handles alert severities with translations', () => {
    const wrapper = createWrapper({
      status: 'critical',
      type: 'alert'
    })
    
    expect(wrapper.text()).toContain('Critical')
  })

  it('falls back to status value for unknown types', () => {
    const wrapper = createWrapper({
      status: 'custom-status',
      type: 'custom'
    })
    
    expect(wrapper.text()).toContain('custom-status')
  })

  it('renders without errors for basic props', () => {
    expect(() => createWrapper()).not.toThrow()
  })

  it('handles custom colors prop', () => {
    expect(() => createWrapper({
      status: 'special',
      customColors: { special: 'purple' }
    })).not.toThrow()
  })

  it('renders consistently with different prop combinations', () => {
    const combos = [
      { status: 'online', type: 'server' },
      { status: 'active', type: 'user' },
      { status: 'warning', type: 'alert' },
      { status: 'custom', type: 'custom', label: 'Custom Label' }
    ]
    
    combos.forEach(props => {
      expect(() => createWrapper(props)).not.toThrow()
    })
  })

  it('preserves i18n functionality', () => {
    const wrapper = createWrapper({
      status: 'online',
      type: 'server'
    })
    
    // Verify translated content appears
    expect(wrapper.text()).toContain('Online')
  })

  it('handles edge cases gracefully', () => {
    const edgeCases = [
      { status: '', type: 'server' },
      { status: 'unknown', type: 'server' },
      { status: 'test', type: undefined }
    ]
    
    edgeCases.forEach(props => {
      expect(() => createWrapper(props)).not.toThrow()
    })
  })

  it('handles multiple status types correctly', () => {
    const testCases = [
      { status: 'online', type: 'server', expectedText: 'Online' },
      { status: 'banned', type: 'user', expectedText: 'Banned' },
      { status: 'critical', type: 'alert', expectedText: 'Critical' }
    ]
    
    testCases.forEach(({ status, type, expectedText }) => {
      const wrapper = createWrapper({ status, type })
      expect(wrapper.text()).toContain(expectedText)
    })
  })
}) 