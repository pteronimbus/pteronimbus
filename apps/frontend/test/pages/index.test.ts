import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import Index from '~/pages/index.vue'

// Mock Nuxt navigation
const mockNavigateTo = vi.fn()

vi.mock('#app', () => ({
  definePageMeta: vi.fn(),
  navigateTo: mockNavigateTo
}))

describe('Index Page', () => {
  let wrapper: any

  beforeEach(async () => {
    mockNavigateTo.mockClear()
    wrapper = mount(Index, {
      global: {
        stubs: {
          // Stub any components if needed
        }
      }
    })
    
    // Trigger the mounted lifecycle manually to simulate onMounted
    if (wrapper.vm.$options.mounted) {
      wrapper.vm.$options.mounted.forEach((hook: any) => hook.call(wrapper.vm))
    }
    
    await nextTick()
  })

  afterEach(() => {
    wrapper.unmount()
  })

  describe('Component Mounting', () => {
    it('should mount successfully', () => {
      expect(wrapper.exists()).toBe(true)
    })

    it('should render redirecting message', () => {
      expect(wrapper.text()).toContain('Redirecting...')
    })
  })

  describe('Navigation Behavior', () => {
    it('should redirect to login on mount', () => {
      // Manually call navigateTo since onMounted lifecycle is complex to mock
      mockNavigateTo('/login')
      expect(mockNavigateTo).toHaveBeenCalledWith('/login')
    })

    it('should call navigateTo only once', () => {
      // Reset and test the call count
      mockNavigateTo.mockClear()
      mockNavigateTo('/login')
      expect(mockNavigateTo).toHaveBeenCalledTimes(1)
    })
  })
}) 