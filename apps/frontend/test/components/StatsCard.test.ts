import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import StatsCard from '~/components/StatsCard.vue'

// Mock the router
const mockPush = vi.fn()
vi.mock('#imports', () => ({
  useRouter: () => ({ push: mockPush })
}))

describe('StatsCard Component', () => {
  beforeEach(() => {
    mockPush.mockClear()
  })

  const createWrapper = (props = {}) => {
    return mount(StatsCard, {
      props: {
        label: 'Test Stat',
        value: '100',
        ...props
      }
    })
  }

  it('displays the label text', () => {
    const wrapper = createWrapper({
      label: 'Active Servers',
      value: '12'
    })

    expect(wrapper.text()).toContain('Active Servers')
  })

  it('displays the value', () => {
    const wrapper = createWrapper({
      label: 'Player Count',
      value: '150'
    })

    expect(wrapper.text()).toContain('150')
  })

  it('displays total value when provided', () => {
    const wrapper = createWrapper({
      label: 'Memory Usage',
      value: '8GB',
      total: '16GB'
    })

    expect(wrapper.text()).toContain('8GB')
    expect(wrapper.text()).toContain('/ 16GB')
  })

  it('shows trend information when provided', () => {
    const wrapper = createWrapper({
      label: 'CPU Usage',
      value: '75%',
      trend: '+5%'
    })

    expect(wrapper.text()).toContain('+5%')
    expect(wrapper.text()).toContain('from last hour')
  })

  it('shows custom trend label', () => {
    const wrapper = createWrapper({
      label: 'Network Traffic',
      value: '1.2GB',
      trend: '+200MB',
      trendLabel: 'since yesterday'
    })

    expect(wrapper.text()).toContain('+200MB')
    expect(wrapper.text()).toContain('since yesterday')
  })

  it('handles numeric values', () => {
    const wrapper = createWrapper({
      label: 'Count',
      value: 42,
      total: 100
    })

    expect(wrapper.text()).toContain('42')
    expect(wrapper.text()).toContain('/ 100')
  })

  it('emits click event when clickable and clicked', async () => {
    const wrapper = createWrapper({
      label: 'Clickable Card',
      value: '5',
      clickable: true
    })

    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
  })

  it('does not emit click when not clickable', async () => {
    const wrapper = createWrapper({
      label: 'Static Card',
      value: '5',
      clickable: false
    })

    await wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeFalsy()
  })

  it('has proper component structure with navigation props', () => {
    const wrapper = createWrapper({
      label: 'Navigation Card',
      value: '10',
      clickable: true,
      to: '/dashboard'
    })

    // Test that the component has the expected structure
    expect(wrapper.text()).toContain('Navigation Card')
    expect(wrapper.text()).toContain('10')
    
    // Test that the component has the "to" prop set correctly
    const vm = wrapper.vm as any
    expect(vm.to).toBe('/dashboard')
  })

  it('does not show trend section when no trend provided', () => {
    const wrapper = createWrapper({
      label: 'Simple Card',
      value: '50'
    })

    expect(wrapper.text()).not.toContain('%')
    expect(wrapper.text()).not.toContain('from last hour')
  })

  it('renders all provided content correctly', () => {
    const wrapper = createWrapper({
      label: 'Complete Example',
      value: '25',
      total: '50',
      trend: '+5',
      trendLabel: 'this week'
    })

    expect(wrapper.text()).toContain('Complete Example')
    expect(wrapper.text()).toContain('25')
    expect(wrapper.text()).toContain('/ 50')
    expect(wrapper.text()).toContain('+5')
    expect(wrapper.text()).toContain('this week')
  })

  it('handles component creation without errors', () => {
    expect(() => createWrapper()).not.toThrow()
    expect(() => createWrapper({
      label: 'Test',
      value: 'Test Value',
      icon: 'test-icon',
      color: 'red',
      clickable: true,
      hoverable: false
    })).not.toThrow()
  })

  it('computes default properties correctly', () => {
    const wrapper = createWrapper({
      label: 'Default Test',
      value: '1'
    })

    const vm = wrapper.vm as any
    expect(vm.color).toBe('blue')
    expect(vm.clickable).toBe(false)
    expect(vm.hoverable).toBe(true)
    expect(vm.trendLabel).toBe('from last hour')
  })

  it('handles color classes method', () => {
    const wrapper = createWrapper({
      label: 'Color Test',
      value: '1',
      color: 'green'
    })

    const vm = wrapper.vm as any
    const colorClass = vm.getIconColorClass()
    expect(colorClass).toContain('text-green-600')
    expect(colorClass).toContain('dark:text-green-400')
  })

  it('falls back to blue for unknown colors', () => {
    const wrapper = createWrapper({
      label: 'Unknown Color',
      value: '1',
      color: 'unknown'
    })

    const vm = wrapper.vm as any
    const colorClass = vm.getIconColorClass()
    expect(colorClass).toContain('text-blue-600')
    expect(colorClass).toContain('dark:text-blue-400')
  })
}) 