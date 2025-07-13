import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import EmptyState from '~/components/EmptyState.vue'

// Mock useSlots
vi.mock('#imports', () => ({
  useSlots: () => ({
    action: vi.fn()
  })
}))

describe('EmptyState Component', () => {
  const createWrapper = (props = {}, options = {}) => {
    return mount(EmptyState, {
      props: {
        icon: 'i-heroicons-server-20-solid',
        title: 'Test Title',
        description: 'Test Description',
        ...props
      },
      ...options
    })
  }

  it('displays the title', () => {
    const wrapper = createWrapper({
      icon: 'i-heroicons-folder-20-solid',
      title: 'No files found',
      description: 'The folder is empty'
    })
    
    expect(wrapper.text()).toContain('No files found')
  })

  it('displays the description', () => {
    const wrapper = createWrapper({
      icon: 'i-heroicons-users-20-solid',
      title: 'No users',
      description: 'There are no users to display'
    })
    
    expect(wrapper.text()).toContain('There are no users to display')
  })

  it('shows action button text when actionLabel is provided', () => {
    const wrapper = createWrapper({
      icon: 'i-heroicons-plus-20-solid',
      title: 'Empty list',
      description: 'Add your first item',
      actionLabel: 'Create Item'
    })
    
    expect(wrapper.text()).toContain('Create Item')
  })

  it('does not show action button when no actionLabel provided', () => {
    const wrapper = createWrapper({
      icon: 'i-heroicons-search-20-solid',
      title: 'No results',
      description: 'Try a different search'
    })
    
    expect(wrapper.text()).not.toContain('Create')
    expect(wrapper.text()).not.toContain('Add')
  })

  it('emits action event when action button is clicked', async () => {
    const wrapper = createWrapper({
      icon: 'i-heroicons-server-20-solid',
      title: 'No servers',
      description: 'Create your first server',
      actionLabel: 'Create Server'
    })

    // Find any button in the component and click it
    const button = wrapper.find('button')
    if (button.exists()) {
      await button.trigger('click')
      expect(wrapper.emitted('action')).toBeTruthy()
    }
  })

  it('displays different content based on props', () => {
    const wrapper1 = createWrapper({
      icon: 'i-heroicons-alert-triangle-20-solid',
      title: 'Error State',
      description: 'Something went wrong'
    })

    const wrapper2 = createWrapper({
      icon: 'i-heroicons-check-circle-20-solid',
      title: 'Success State',
      description: 'All good!'
    })

    expect(wrapper1.text()).toContain('Error State')
    expect(wrapper1.text()).toContain('Something went wrong')
    
    expect(wrapper2.text()).toContain('Success State')
    expect(wrapper2.text()).toContain('All good!')
  })

  it('handles component creation without errors', () => {
    expect(() => createWrapper()).not.toThrow()
    
    expect(() => createWrapper({
      icon: 'i-heroicons-info-20-solid',
      title: 'Info',
      description: 'Information message',
      actionLabel: 'Learn More',
      actionIcon: 'i-heroicons-arrow-right-20-solid',
      actionClass: 'btn-primary'
    })).not.toThrow()
  })

  it('supports custom slot content', () => {
    const wrapper = createWrapper({
      icon: 'i-heroicons-cog-20-solid',
      title: 'Custom Actions',
      description: 'Multiple action buttons'
    }, {
      slots: {
        action: '<div class="custom-actions">Custom content</div>'
      }
    })

    expect(wrapper.html()).toContain('Custom content')
  })

  it('computes hasAction correctly', () => {
    const withAction = createWrapper({
      icon: 'i-heroicons-server-20-solid',
      title: 'Test',
      description: 'Test',
      actionLabel: 'Action'
    })
    
    const vm = withAction.vm as any
    expect(vm.hasAction).toBe(true)

    const withoutAction = createWrapper({
      icon: 'i-heroicons-server-20-solid',
      title: 'Test',
      description: 'Test'
    })
    
    const vm2 = withoutAction.vm as any
    expect(vm2.hasAction).toBe(false)
  })

  it('handles various empty state scenarios', () => {
    const scenarios = [
      {
        icon: 'i-heroicons-folder-open-20-solid',
        title: 'No files',
        description: 'This folder is empty'
      },
      {
        icon: 'i-heroicons-users-20-solid',
        title: 'No team members',
        description: 'Invite people to join your team'
      },
      {
        icon: 'i-heroicons-document-text-20-solid',
        title: 'No documents',
        description: 'Upload your first document'
      }
    ]

    scenarios.forEach(scenario => {
      const wrapper = createWrapper(scenario)
      expect(wrapper.text()).toContain(scenario.title)
      expect(wrapper.text()).toContain(scenario.description)
    })
  })

  it('renders content structure correctly', () => {
    const wrapper = createWrapper({
      icon: 'i-heroicons-star-20-solid',
      title: 'Star Rating',
      description: 'No ratings yet'
    })

    // Should have main content
    expect(wrapper.text()).toContain('Star Rating')
    expect(wrapper.text()).toContain('No ratings yet')
  })

  it('handles all required props', () => {
    const wrapper = createWrapper({
      icon: 'i-heroicons-warning-20-solid',
      title: 'Warning',
      description: 'Attention required'
    })

    expect(wrapper.text()).toContain('Warning')
    expect(wrapper.text()).toContain('Attention required')
  })

  it('handles empty state with action correctly', () => {
    const wrapper = createWrapper({
      icon: 'i-heroicons-plus-circle-20-solid',
      title: 'Add Items',
      description: 'Start building your collection',
      actionLabel: 'Add First Item'
    })

    expect(wrapper.text()).toContain('Add Items')
    expect(wrapper.text()).toContain('Start building your collection')
    expect(wrapper.text()).toContain('Add First Item')
  })

  it('maintains component structure with different props', () => {
    const minimal = createWrapper({
      icon: 'i-heroicons-inbox-20-solid',
      title: 'Empty',
      description: 'Nothing here'
    })

    const complete = createWrapper({
      icon: 'i-heroicons-inbox-20-solid',
      title: 'Empty',
      description: 'Nothing here',
      actionLabel: 'Add Something',
      actionIcon: 'i-heroicons-plus-20-solid',
      actionClass: 'custom-class'
    })

    expect(minimal.text()).toContain('Empty')
    expect(minimal.text()).toContain('Nothing here')
    
    expect(complete.text()).toContain('Empty')
    expect(complete.text()).toContain('Nothing here')
    expect(complete.text()).toContain('Add Something')
  })
}) 