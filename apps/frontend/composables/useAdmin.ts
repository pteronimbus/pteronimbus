import type { Controller } from '~/types/controller'

interface AdminState {
  controllers: Controller[]
  isLoading: boolean
  error: string | null
}

// Global admin state
const adminState = ref<AdminState>({
  controllers: [],
  isLoading: false,
  error: null
})

export const useAdmin = () => {
  const { apiRequest } = useAuth()
  const config = useRuntimeConfig()

  // Computed properties
  const controllers = computed(() => adminState.value.controllers)
  const isLoading = computed(() => adminState.value.isLoading)
  const error = computed(() => adminState.value.error)

  // Set loading state
  const setLoading = (loading: boolean) => {
    adminState.value.isLoading = loading
  }

  // Set error state
  const setError = (errorMessage: string) => {
    adminState.value.error = errorMessage
    adminState.value.isLoading = false
  }

  // Clear error state
  const clearError = () => {
    adminState.value.error = null
  }

  // Fetch all controllers
  const fetchControllers = async () => {
    setLoading(true)
    clearError()

    try {
      const response = await apiRequest<{ controllers: Controller[] }>(`${config.public.backendUrl}/api/controllers`)
      // Ensure we always have an array, even if the backend returns null
      adminState.value.controllers = response.controllers || []
      return adminState.value.controllers
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to fetch controllers'
      console.error('Failed to fetch controllers:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Refresh controllers
  const refreshControllers = async () => {
    return await fetchControllers()
  }

  // Cleanup inactive controllers
  const cleanupInactiveControllers = async () => {
    setLoading(true)
    clearError()

    try {
      await apiRequest(`${config.public.backendUrl}/api/admin/cleanup-controllers`, {
        method: 'POST'
      })

      // Refresh controllers after cleanup
      await fetchControllers()
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to cleanup inactive controllers'
      console.error('Failed to cleanup inactive controllers:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Get controller status
  const getControllerStatus = async (controllerId: string) => {
    try {
      const response = await apiRequest<{ controller: Controller }>(`${config.public.backendUrl}/api/controllers/${controllerId}`)
      return response.controller
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to get controller status'
      console.error('Failed to get controller status:', error)
      throw new Error(errorMsg)
    }
  }

  // Restart controller (placeholder for future implementation)
  const restartController = async (controllerId: string) => {
    // TODO: Implement controller restart functionality
    console.log('Restarting controller:', controllerId)
    throw new Error('Controller restart not yet implemented')
  }

  // Remove controller (placeholder for future implementation)
  const removeController = async (controllerId: string) => {
    // TODO: Implement controller removal functionality
    console.log('Removing controller:', controllerId)
    throw new Error('Controller removal not yet implemented')
  }

  // Approve controller
  const approveController = async (controllerId: string) => {
    setLoading(true)
    clearError()

    try {
      await apiRequest(`${config.public.backendUrl}/api/controllers/${controllerId}/approve`, {
        method: 'POST'
      })

      // Refresh controllers after approval
      await fetchControllers()
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to approve controller'
      console.error('Failed to approve controller:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Reject controller
  const rejectController = async (controllerId: string, reason?: string) => {
    setLoading(true)
    clearError()

    try {
      await apiRequest(`${config.public.backendUrl}/api/controllers/${controllerId}/reject`, {
        method: 'POST',
        body: {
          action: 'reject',
          reason: reason || ''
        }
      })

      // Refresh controllers after rejection
      await fetchControllers()
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to reject controller'
      console.error('Failed to reject controller:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Reset admin state (for testing)
  const resetAdminState = () => {
    adminState.value = {
      controllers: [],
      isLoading: false,
      error: null
    }
  }

  return {
    // State
    controllers,
    isLoading: readonly(isLoading),
    error: readonly(error),

    // Methods
    fetchControllers,
    refreshControllers,
    cleanupInactiveControllers,
    getControllerStatus,
    restartController,
    removeController,
    approveController,
    rejectController,
    clearError,
    resetAdminState
  }
} 