import { describe, it, expect, vi } from 'vitest'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import LoginPage from '~/pages/login.vue'

// Mock the router
const mockRouter = {
  push: vi.fn()
}
vi.mock('vue-router', () => ({
  useRouter: () => mockRouter
}))

// Mock useUser composable
const mockUser = { value: null }
vi.mock('~/composables/useUser', () => ({
  useUser: () => mockUser
}))

describe('Login Page', () => {
  it('should allow a user to log in and redirect to the dashboard', async () => {
    const wrapper = await mountSuspended(LoginPage)

    // Find input fields and button
    const emailInput = wrapper.find('input[type="email"]')
    const passwordInput = wrapper.find('input[type="password"]')
    const form = wrapper.find('form')

    // Set values
    await emailInput.setValue('test@example.com')
    await passwordInput.setValue('password')

    // Submit the form
    await form.trigger('submit.prevent')
    
    // Assert user state is set
    expect(mockUser.value).toEqual({ email: 'test@example.com', name: 'Test User' })

    // Assert redirection
    expect(mockRouter.push).toHaveBeenCalledWith('/dashboard')
  })
}) 