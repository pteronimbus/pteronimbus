import { ref, computed, nextTick } from 'vue'

interface Tenant {
  id: string
  discord_server_id: string
  name: string
  icon: string
  owner_id: string
  config: TenantConfig
  created_at: string
  updated_at: string
}

interface TenantConfig {
  default_game_template?: string
  resource_limits?: ResourceLimits
  notification_channels?: string[]
  settings?: Record<string, string>
}

interface ResourceLimits {
  max_game_servers: number
  max_cpu: string
  max_memory: string
  max_storage: string
}

interface DiscordGuild {
  id: string
  name: string
  icon: string
  owner: boolean
  permissions: string
  features: string[]
}

interface TenantState {
  tenants: Tenant[]
  currentTenant: Tenant | null
  availableGuilds: DiscordGuild[]
  isLoading: boolean
  error: string | null
}

// Global tenant state
const tenantState = ref<TenantState>({
  tenants: [],
  currentTenant: null,
  availableGuilds: [],
  isLoading: false,
  error: null
})

export const useTenant = () => {
  const { apiRequest } = useAuth()
  const config = useRuntimeConfig()
  const router = useRouter()

  // Computed properties
  const tenants = computed(() => tenantState.value.tenants)
  const currentTenant = computed(() => tenantState.value.currentTenant)
  const availableGuilds = computed(() => tenantState.value.availableGuilds)
  const isLoading = computed(() => tenantState.value.isLoading)
  const error = computed(() => tenantState.value.error)

  // Set loading state
  const setLoading = (loading: boolean) => {
    tenantState.value.isLoading = loading
  }

  // Set error state
  const setError = (errorMessage: string) => {
    tenantState.value.error = errorMessage
    tenantState.value.isLoading = false
  }

  // Clear error state
  const clearError = () => {
    tenantState.value.error = null
  }

  // Initialize tenant state from localStorage
  const initializeTenant = () => {
    if (import.meta.client) {
      const currentTenantData = localStorage.getItem('current_tenant')
      if (currentTenantData) {
        try {
          tenantState.value.currentTenant = JSON.parse(currentTenantData)
        } catch (error) {
          console.error('Failed to parse stored tenant data:', error)
          localStorage.removeItem('current_tenant')
        }
      }
    }
  }

  // Store current tenant
  const storeCurrentTenant = (tenant: Tenant | null) => {
    tenantState.value.currentTenant = tenant
    if (import.meta.client) {
      if (tenant) {
        localStorage.setItem('current_tenant', JSON.stringify(tenant))
      } else {
        localStorage.removeItem('current_tenant')
      }
    }
  }

  // Get user's tenants
  const fetchUserTenants = async () => {
    setLoading(true)
    clearError()

    try {
      const response = await apiRequest<{ tenants: Tenant[] }>(`${config.public.backendUrl}/api/tenants`)
      tenantState.value.tenants = response.tenants
      return response.tenants
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to fetch tenants'
      console.error('Failed to fetch tenants:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Get available Discord guilds for tenant creation
  const fetchAvailableGuilds = async () => {
    setLoading(true)
    clearError()

    try {
      const response = await apiRequest<{ guilds: DiscordGuild[] }>(`${config.public.backendUrl}/api/tenants/available-guilds`)
      tenantState.value.availableGuilds = response.guilds
      return response.guilds
    } catch (error: any) {
      // Handle specific Discord token missing error
      if (error?.data?.code === 'DISCORD_TOKEN_MISSING') {
        const errorMsg = 'Please log in again to refresh your Discord connection'
        console.error('Discord token missing, user needs to re-authenticate:', error)
        setError(errorMsg)

        // Redirect to login after a short delay
        setTimeout(() => {
          window.location.href = '/login?reason=discord_token_expired'
        }, 2000)

        throw new Error(errorMsg)
      }

      const errorMsg = error?.data?.message || 'Failed to fetch available guilds'
      console.error('Failed to fetch available guilds:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Create a new tenant
  const createTenant = async (guildId: string) => {
    setLoading(true)
    clearError()

    try {
      const response = await apiRequest<{ tenant: Tenant }>(`${config.public.backendUrl}/api/tenants`, {
        method: 'POST',
        body: {
          guild_id: guildId
        }
      })

      // Add the new tenant to the list
      tenantState.value.tenants.push(response.tenant)

      // Set as current tenant
      storeCurrentTenant(response.tenant)

      return response.tenant
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to create tenant'
      console.error('Failed to create tenant:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Get a specific tenant
  const fetchTenant = async (tenantId: string) => {
    setLoading(true)
    clearError()

    try {
      const response = await apiRequest<{ tenant: Tenant }>(`${config.public.backendUrl}/api/tenants/${tenantId}`)
      return response.tenant
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to fetch tenant'
      console.error('Failed to fetch tenant:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Update tenant configuration
  const updateTenantConfig = async (tenantId: string, tenantConfig: TenantConfig) => {
    setLoading(true)
    clearError()

    try {
      await apiRequest(`${config.public.backendUrl}/api/tenants/${tenantId}/config`, {
        method: 'PUT',
        body: tenantConfig
      })

      // Update the tenant in the list if it exists
      const tenantIndex = tenantState.value.tenants.findIndex(t => t.id === tenantId)
      if (tenantIndex !== -1) {
        tenantState.value.tenants[tenantIndex].config = tenantConfig
      }

      // Update current tenant if it's the same
      if (tenantState.value.currentTenant?.id === tenantId) {
        tenantState.value.currentTenant.config = tenantConfig
        storeCurrentTenant(tenantState.value.currentTenant)
      }
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to update tenant config'
      console.error('Failed to update tenant config:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Sync tenant data with Discord
  const syncTenantData = async (tenantId: string) => {
    setLoading(true)
    clearError()

    try {
      await apiRequest(`${config.public.backendUrl}/api/tenants/${tenantId}/sync`, {
        method: 'POST'
      })
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to sync tenant data'
      console.error('Failed to sync tenant data:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Delete a tenant
  const deleteTenant = async (tenantId: string) => {
    setLoading(true)
    clearError()

    try {
      await apiRequest(`${config.public.backendUrl}/api/tenants/${tenantId}`, {
        method: 'DELETE'
      })

      // Remove from tenants list
      tenantState.value.tenants = tenantState.value.tenants.filter(t => t.id !== tenantId)

      // Clear current tenant if it was deleted
      if (tenantState.value.currentTenant?.id === tenantId) {
        storeCurrentTenant(null)
      }
    } catch (error: any) {
      const errorMsg = error?.data?.message || 'Failed to delete tenant'
      console.error('Failed to delete tenant:', error)
      setError(errorMsg)
      throw new Error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  // Switch to a different tenant
  const switchTenant = async (tenant: Tenant) => {
    // Store the tenant first
    storeCurrentTenant(tenant)

    // Wait a tick to ensure the state is updated
    await nextTick()

    // Navigate to tenant dashboard
    await router.push(`/tenant/${tenant.id}/dashboard`)
  }

  // Clear current tenant
  const clearCurrentTenant = () => {
    storeCurrentTenant(null)
  }

  // Reset all tenant state (for testing)
  const resetTenantState = () => {
    tenantState.value = {
      tenants: [],
      currentTenant: null,
      availableGuilds: [],
      isLoading: false,
      error: null
    }
  }

  // Get tenant context for API requests
  const getTenantHeaders = () => {
    if (tenantState.value.currentTenant) {
      return {
        'X-Tenant-ID': tenantState.value.currentTenant.id
      }
    }
    return {}
  }

  // API request with tenant context
  const tenantApiRequest = async <T>(url: string, options: any = {}) => {
    const tenantHeaders = getTenantHeaders()
    const requestOptions = {
      ...options,
      headers: {
        ...options.headers,
        ...tenantHeaders
      }
    }

    // Ensure the URL includes the backend URL if it's a relative path
    const fullUrl = url.startsWith('http') ? url : `${config.public.backendUrl}${url}`

    return await apiRequest<T>(fullUrl, requestOptions)
  }

  return {
    // State
    tenants: readonly(tenants),
    currentTenant: readonly(currentTenant),
    availableGuilds: readonly(availableGuilds),
    isLoading: readonly(isLoading),
    error: readonly(error),

    // Methods
    fetchUserTenants,
    fetchAvailableGuilds,
    createTenant,
    fetchTenant,
    updateTenantConfig,
    syncTenantData,
    deleteTenant,
    switchTenant,
    clearCurrentTenant,
    initializeTenant,
    clearError,
    tenantApiRequest,
    getTenantHeaders,
    resetTenantState
  }
}