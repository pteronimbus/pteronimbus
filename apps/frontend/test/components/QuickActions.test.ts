import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import QuickActions from '~/components/QuickActions.vue'

describe('QuickActions Component', () => {
  const createWrapper = (props = {}) => {
    return mount(QuickActions, {
      props: {
        title: 'Quick Actions',
        actions: [],
        ...props
      }
    })
  }

  it('displays the title', () => {
    const wrapper = createWrapper({
      title: 'Dashboard Actions'
    })
    
    expect(wrapper.text()).toContain('Dashboard Actions')
  })

  it('uses default title when not provided', () => {
    const wrapper = createWrapper()
    
    expect(wrapper.text()).toContain('Quick Actions')
  })

  it('renders action buttons when actions provided', () => {
    const wrapper = createWrapper({
      actions: [
        {
          label: 'Create Server',
          icon: 'i-heroicons-plus-circle-20-solid',
          onClick: vi.fn()
        }
      ]
    })
    
    expect(wrapper.text()).toContain('Create Server')
  })

  it('handles button clicks', async () => {
    const mockClick = vi.fn()
    const wrapper = createWrapper({
      actions: [
        {
          label: 'Test Action',
          onClick: mockClick
        }
      ]
    })
    
    const button = wrapper.find('button')
    await button.trigger('click')
    
    expect(mockClick).toHaveBeenCalled()
  })

  it('renders multiple actions', () => {
    const wrapper = createWrapper({
      actions: [
        { label: 'Action 1', onClick: vi.fn() },
        { label: 'Action 2', onClick: vi.fn() },
        { label: 'Action 3', onClick: vi.fn() }
      ]
    })
    
    expect(wrapper.text()).toContain('Action 1')
    expect(wrapper.text()).toContain('Action 2')
    expect(wrapper.text()).toContain('Action 3')
    
    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(3)
  })

  it('handles actions with different colors', () => {
    const wrapper = createWrapper({
      actions: [
        {
          label: 'Primary Action',
          color: 'primary',
          onClick: vi.fn()
        },
        {
          label: 'Warning Action', 
          color: 'warning',
          onClick: vi.fn()
        }
      ]
    })
    
    expect(wrapper.text()).toContain('Primary Action')
    expect(wrapper.text()).toContain('Warning Action')
  })

  it('handles different grid column configurations', () => {
    const wrapper2Col = createWrapper({
      gridCols: 2,
      actions: [
        { label: 'Action 1', onClick: vi.fn() },
        { label: 'Action 2', onClick: vi.fn() }
      ]
    })
    
    const wrapper4Col = createWrapper({
      gridCols: 4,
      actions: [
        { label: 'Action 1', onClick: vi.fn() },
        { label: 'Action 2', onClick: vi.fn() },
        { label: 'Action 3', onClick: vi.fn() },
        { label: 'Action 4', onClick: vi.fn() }
      ]
    })
    
    expect(wrapper2Col.text()).toContain('Action 1')
    expect(wrapper4Col.text()).toContain('Action 4')
  })

  it('handles empty actions array', () => {
    const wrapper = createWrapper({
      actions: []
    })
    
    expect(wrapper.text()).toContain('Quick Actions')
    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(0)
  })

  it('handles component creation without errors', () => {
    expect(() => createWrapper()).not.toThrow()
    expect(() => createWrapper({
      title: 'Custom Actions',
      gridCols: 3,
      actions: [
        {
          label: 'Complete Action',
          icon: 'i-heroicons-cog-6-tooth-20-solid',
          color: 'secondary',
          variant: 'outline',
          size: 'md',
          class: 'custom-class',
          onClick: vi.fn()
        }
      ]
    })).not.toThrow()
  })

  it('renders actions with icons', () => {
    const wrapper = createWrapper({
      actions: [
        {
          label: 'Settings',
          icon: 'i-heroicons-cog-6-tooth-20-solid',
          onClick: vi.fn()
        }
      ]
    })
    
    expect(wrapper.text()).toContain('Settings')
  })

  it('handles default properties correctly', () => {
    const wrapper = createWrapper()
    
    const vm = wrapper.vm as any
    expect(vm.title).toBe('Quick Actions')
    expect(vm.gridCols).toBe(2) // default grid columns
    expect(Array.isArray(vm.actions)).toBe(true)
  })

  it('renders dashboard-style actions', () => {
    const wrapper = createWrapper({
      title: 'Dashboard Actions',
      actions: [
        {
          label: 'Create Server',
          icon: 'i-heroicons-plus-circle-20-solid',
          color: 'primary',
          onClick: vi.fn()
        },
        {
          label: 'Create User',
          icon: 'i-heroicons-user-plus-20-solid',
          color: 'secondary',
          onClick: vi.fn()
        },
        {
          label: 'Refresh',
          icon: 'i-heroicons-arrow-path-20-solid',
          color: 'success',
          onClick: vi.fn()
        },
        {
          label: 'Settings',
          icon: 'i-heroicons-cog-6-tooth-20-solid',
          color: 'warning',
          onClick: vi.fn()
        }
      ]
    })
    
    expect(wrapper.text()).toContain('Create Server')
    expect(wrapper.text()).toContain('Create User')
    expect(wrapper.text()).toContain('Refresh')
    expect(wrapper.text()).toContain('Settings')
  })

  it('handles monitoring-style actions', () => {
    const wrapper = createWrapper({
      title: 'System Actions',
      gridCols: 2,
      actions: [
        {
          label: 'Restart Services',
          icon: 'i-heroicons-arrow-path-20-solid',
          color: 'primary',
          onClick: vi.fn()
        },
        {
          label: 'Clear Cache',
          icon: 'i-heroicons-trash-20-solid',
          color: 'warning',
          onClick: vi.fn()
        }
      ]
    })
    
    expect(wrapper.text()).toContain('Restart Services')
    expect(wrapper.text()).toContain('Clear Cache')
  })

  it('passes through button properties correctly', () => {
    const wrapper = createWrapper({
      actions: [
        {
          label: 'Styled Action',
          variant: 'outline',
          size: 'lg',
          class: 'custom-button-class',
          onClick: vi.fn()
        }
      ]
    })
    
    expect(wrapper.text()).toContain('Styled Action')
  })
}) 