import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import PageHeader from '~/components/PageHeader.vue'

describe('PageHeader Component', () => {
  const createWrapper = (props = {}) => {
    return mount(PageHeader, {
      props: {
        title: 'Test Title',
        ...props
      }
    })
  }

  it('displays the title', () => {
    const wrapper = createWrapper({
      title: 'Dashboard'
    })
    
    expect(wrapper.text()).toContain('Dashboard')
  })

  it('displays the description when provided', () => {
    const wrapper = createWrapper({
      title: 'Dashboard',
      description: 'Monitor and manage your game servers'
    })
    
    expect(wrapper.text()).toContain('Monitor and manage your game servers')
  })

  it('does not show description when not provided', () => {
    const wrapper = createWrapper({
      title: 'Simple Header'
    })
    
    // Should only contain the title, not any description text
    expect(wrapper.text()).toContain('Simple Header')
    expect(wrapper.text()).not.toContain('Monitor and manage')
  })

  it('renders action buttons when provided', () => {
    const wrapper = createWrapper({
      title: 'Actions Test',
      actions: [
        {
          label: 'Create Server',
          color: 'primary',
          onClick: vi.fn()
        }
      ]
    })
    
    expect(wrapper.text()).toContain('Create Server')
  })

  it('renders badge actions when provided', () => {
    const wrapper = createWrapper({
      title: 'Badge Test',
      actions: [
        {
          label: 'System Healthy',
          type: 'badge',
          color: 'success'
        }
      ]
    })
    
    expect(wrapper.text()).toContain('System Healthy')
  })

  it('handles mixed action types', () => {
    const wrapper = createWrapper({
      title: 'Mixed Actions',
      actions: [
        {
          label: 'Create Button',
          color: 'primary',
          onClick: vi.fn()
        },
        {
          label: 'Status Badge',
          type: 'badge',
          color: 'success'
        }
      ]
    })
    
    expect(wrapper.text()).toContain('Create Button')
    expect(wrapper.text()).toContain('Status Badge')
  })

  it('handles component creation without errors', () => {
    expect(() => createWrapper()).not.toThrow()
    expect(() => createWrapper({
      title: 'Complete Test',
      description: 'Full featured header',
      actions: [
        {
          label: 'Action',
          color: 'primary',
          variant: 'solid',
          size: 'sm',
          onClick: vi.fn()
        }
      ]
    })).not.toThrow()
  })

  it('renders multiple actions correctly', () => {
    const wrapper = createWrapper({
      title: 'Multiple Actions',
      actions: [
        { label: 'Action 1', onClick: vi.fn() },
        { label: 'Action 2', onClick: vi.fn() },
        { label: 'Action 3', type: 'badge', color: 'info' }
      ]
    })
    
    expect(wrapper.text()).toContain('Action 1')
    expect(wrapper.text()).toContain('Action 2')
    expect(wrapper.text()).toContain('Action 3')
  })

  it('handles empty actions array', () => {
    const wrapper = createWrapper({
      title: 'No Actions',
      actions: []
    })
    
    expect(wrapper.text()).toContain('No Actions')
  })
}) 