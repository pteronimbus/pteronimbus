import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import StatusBadge from '~/components/StatusBadge.vue'

// Create a complete mock for the component's setup
const mockUseI18n = vi.fn(() => ({
  t: (key: string) => {
    // Simple mock that returns a predictable translation
    const translations: Record<string, string> = {
      'servers.status.online': 'Online',
      'servers.status.offline': 'Offline',
      'users.status.banned': 'Banned',
      'users.status.active': 'Active',
      'alerts.severity.critical': 'Critical',
      'alerts.severity.warning': 'Warning'
    }
    return translations[key] || key
  }
}))

// Mock the auto-imports
vi.mock('#imports', () => ({
  useI18n: mockUseI18n
}))

describe('StatusBadge Component', () => {
  const createWrapper = (props = {}) => {
    return mount(StatusBadge, {
      props: {
        status: 'online',
        type: 'server',
        ...props
      },
      global: {
        stubs: {
          UBadge: {
            template: '<span class="test-badge" v-bind="$attrs"><slot /></span>',
            props: ['color', 'variant', 'class']
          }
        }
      }
    })
  }

  it('renders the component successfully', () => {
    const wrapper = createWrapper()
    expect(wrapper.find('.test-badge').exists()).toBe(true)
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
    
    // Verify that useI18n was called (ensuring i18n integration works)
    expect(mockUseI18n).toHaveBeenCalled()
    
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

  it('maintains proper component structure', () => {
    const wrapper = createWrapper({
      status: 'online',
      type: 'server'
    })
    
    expect(wrapper.vm).toBeTruthy()
    expect(wrapper.html()).toBeTruthy()
    expect(wrapper.find('.test-badge').exists()).toBe(true)
  })
}) 