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

// Mock Vue's onMounted to execute immediately
vi.mock('vue', async () => {
  const actual = await vi.importActual('vue')
  return {
    ...actual,
    onMounted: (callback: () => void) => {
      // Execute the callback immediately when onMounted is called
      setTimeout(callback, 0)
    }
  }
})

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
    // Wait for the component to fully mount and onMounted to execute
    await nextTick()
    await new Promise(resolve => setTimeout(resolve, 10))
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
      expect(mockNavigateTo).toHaveBeenCalledWith('/login')
    })

    it('should call navigateTo only once', () => {
      expect(mockNavigateTo).toHaveBeenCalledTimes(1)
    })
  })
}) 