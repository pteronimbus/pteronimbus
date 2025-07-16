import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useTenant } from '~/composables/useTenant'

// Mock the useAuth composable
const mockApiRequest = vi.fn()
vi.mock('~/composables/useAuth', () => ({
  useAuth: () => ({
    apiRequest: mockApiRequest,
    initializeAuth: vi.fn()
  })
}))

// Mock router
const mockPush = vi.fn()

// Mock Nuxt runtime config
vi.mock('#app', () => ({
  useRuntimeConfig: () => ({
    public: {
      backendUrl: 'http://localhost:8080'
    }
  }),
  useRouter: () => ({
    push: mockPush
  })
}))

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
})

describe('useTenant', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Reset import.meta.client to true for tests
    Object.defineProperty(import.meta, 'client', {
      value: true,
      writable: true
    })
    // Reset tenant state between tests
    const { resetTenantState } = useTenant()
    resetTenantState()
  })

  describe('fetchUserTenants', () => {
    it('should fetch user tenants successfully', async () => {
      const mockTenants = [
        {
          id: 'tenant-1',
          discord_server_id: 'guild-1',
          name: 'Test Guild 1',
          icon: 'icon-1',
          owner_id: 'user-1',
          config: {},
          created_at: '2023-01-01T00:00:00Z',
          updated_at: '2023-01-01T00:00:00Z'
        },
        {
          id: 'tenant-2',
          discord_server_id: 'guild-2',
          name: 'Test Guild 2',
          icon: 'icon-2',
          owner_id: 'user-1',
          config: {},
          created_at: '2023-01-01T00:00:00Z',
          updated_at: '2023-01-01T00:00:00Z'
        }
      ]

      mockApiRequest.mockResolvedValue({ tenants: mockTenants })

      const { fetchUserTenants, tenants } = useTenant()
      
      const result = await fetchUserTenants()

      expect(mockApiRequest).toHaveBeenCalledWith('http://localhost:8080/api/tenants')
      expect(result).toEqual(mockTenants)
      expect(tenants.value).toEqual(mockTenants)
    })

    it('should handle fetch error', async () => {
      const errorMessage = 'Failed to fetch tenants'
      mockApiRequest.mockRejectedValue({
        data: { message: errorMessage }
      })

      const { fetchUserTenants, error } = useTenant()

      await expect(fetchUserTenants()).rejects.toThrow(errorMessage)
      expect(error.value).toBe(errorMessage)
    })
  })

  describe('fetchAvailableGuilds', () => {
    it('should fetch available guilds successfully', async () => {
      const mockGuilds = [
        {
          id: 'guild-1',
          name: 'Available Guild 1',
          icon: 'icon-1',
          owner: true,
          permissions: '2147483647',
          features: []
        },
        {
          id: 'guild-2',
          name: 'Available Guild 2',
          icon: 'icon-2',
          owner: false,
          permissions: '32',
          features: []
        }
      ]

      mockApiRequest.mockResolvedValue({ guilds: mockGuilds })

      const { fetchAvailableGuilds, availableGuilds } = useTenant()
      
      const result = await fetchAvailableGuilds()

      expect(mockApiRequest).toHaveBeenCalledWith('http://localhost:8080/api/tenants/available-guilds')
      expect(result).toEqual(mockGuilds)
      expect(availableGuilds.value).toEqual(mockGuilds)
    })
  })

  describe('createTenant', () => {
    it('should create tenant successfully', async () => {
      const guildId = 'guild-123'
      const mockTenant = {
        id: 'tenant-123',
        discord_server_id: guildId,
        name: 'New Guild',
        icon: 'icon-123',
        owner_id: 'user-1',
        config: {},
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T00:00:00Z'
      }

      mockApiRequest.mockResolvedValue({ tenant: mockTenant })

      const { createTenant, tenants, currentTenant } = useTenant()
      
      const result = await createTenant(guildId)

      expect(mockApiRequest).toHaveBeenCalledWith('http://localhost:8080/api/tenants', {
        method: 'POST',
        body: { guild_id: guildId }
      })
      expect(result).toEqual(mockTenant)
      expect(tenants.value).toContainEqual(mockTenant)
      expect(currentTenant.value).toEqual(mockTenant)
    })

    it('should handle create tenant error', async () => {
      const guildId = 'guild-123'
      const errorMessage = 'Failed to create tenant'
      
      mockApiRequest.mockRejectedValue({
        data: { message: errorMessage }
      })

      const { createTenant, error } = useTenant()

      await expect(createTenant(guildId)).rejects.toThrow(errorMessage)
      expect(error.value).toBe(errorMessage)
    })
  })

  describe('updateTenantConfig', () => {
    it('should update tenant config successfully', async () => {
      const tenantId = 'tenant-123'
      const config = {
        resource_limits: {
          max_game_servers: 10,
          max_cpu: '4',
          max_memory: '8Gi',
          max_storage: '20Gi'
        },
        settings: {
          notification_channel: 'general'
        }
      }

      // Set up initial tenant state by mocking fetchUserTenants
      const mockTenant = {
        id: tenantId,
        discord_server_id: 'guild-123',
        name: 'Test Guild',
        config: {}
      }

      // First mock the fetch to set up initial state
      mockApiRequest.mockResolvedValueOnce({ tenants: [mockTenant] })
      
      const { updateTenantConfig, fetchUserTenants, tenants, currentTenant } = useTenant()
      
      // Set up initial state
      await fetchUserTenants()
      
      // Mock the update request
      mockApiRequest.mockResolvedValueOnce({})

      await updateTenantConfig(tenantId, config)

      expect(mockApiRequest).toHaveBeenCalledWith(`http://localhost:8080/api/tenants/${tenantId}/config`, {
        method: 'PUT',
        body: config
      })
    })
  })

  describe('deleteTenant', () => {
    it('should delete tenant successfully', async () => {
      const tenantId = 'tenant-123'
      const mockTenant = {
        id: tenantId,
        discord_server_id: 'guild-123',
        name: 'Test Guild'
      }

      // First set up initial state by fetching tenants
      mockApiRequest.mockResolvedValueOnce({ tenants: [mockTenant] })
      
      const { deleteTenant, fetchUserTenants, tenants, currentTenant } = useTenant()
      
      // Set up initial state
      await fetchUserTenants()
      
      // Mock the delete request
      mockApiRequest.mockResolvedValueOnce({})

      await deleteTenant(tenantId)

      expect(mockApiRequest).toHaveBeenCalledWith(`http://localhost:8080/api/tenants/${tenantId}`, {
        method: 'DELETE'
      })
      expect(tenants.value).not.toContain(mockTenant)
    })
  })

  describe('switchTenant', () => {
    it('should switch to tenant and store in localStorage', async () => {
      const tenant = {
        id: 'tenant-123',
        discord_server_id: 'guild-123',
        name: 'Test Guild'
      }

      const { switchTenant, currentTenant } = useTenant()
      
      try {
        await switchTenant(tenant)
      } catch (error) {
        // Router navigation might fail in test environment, but we still want to test the core functionality
      }

      expect(currentTenant.value).toEqual(tenant)
      expect(localStorageMock.setItem).toHaveBeenCalledWith('current_tenant', JSON.stringify(tenant))
    })
  })

  describe('initializeTenant', () => {
    it('should initialize tenant from localStorage', () => {
      const mockTenant = {
        id: 'tenant-123',
        discord_server_id: 'guild-123',
        name: 'Test Guild'
      }

      localStorageMock.getItem.mockReturnValue(JSON.stringify(mockTenant))

      const { initializeTenant, currentTenant } = useTenant()
      
      initializeTenant()

      expect(localStorageMock.getItem).toHaveBeenCalledWith('current_tenant')
      expect(currentTenant.value).toEqual(mockTenant)
    })

    it('should handle invalid localStorage data', () => {
      localStorageMock.getItem.mockReturnValue('invalid-json')

      const { initializeTenant, currentTenant } = useTenant()
      
      initializeTenant()

      expect(localStorageMock.removeItem).toHaveBeenCalledWith('current_tenant')
      expect(currentTenant.value).toBeNull()
    })
  })

  describe('getTenantHeaders', () => {
    it('should return tenant headers when tenant is selected', async () => {
      const tenant = {
        id: 'tenant-123',
        discord_server_id: 'guild-123',
        name: 'Test Guild'
      }

      // Mock the create tenant to set current tenant
      mockApiRequest.mockResolvedValue({ tenant })

      const { getTenantHeaders, createTenant } = useTenant()
      
      // Set current tenant by creating one
      await createTenant('guild-123')

      const headers = getTenantHeaders()

      expect(headers).toEqual({
        'X-Tenant-ID': tenant.id
      })
    })

    it('should return empty headers when no tenant is selected', () => {
      const { getTenantHeaders, resetTenantState } = useTenant()
      
      // Ensure no tenant is selected
      resetTenantState()

      const headers = getTenantHeaders()

      expect(headers).toEqual({})
    })
  })

  describe('tenantApiRequest', () => {
    it('should make API request with tenant headers', async () => {
      const tenant = {
        id: 'tenant-123',
        discord_server_id: 'guild-123',
        name: 'Test Guild'
      }

      const mockResponse = { data: 'test' }
      
      // First mock create tenant to set current tenant
      mockApiRequest.mockResolvedValueOnce({ tenant })
      
      const { tenantApiRequest, createTenant } = useTenant()
      
      // Set current tenant by creating one
      await createTenant('guild-123')
      
      // Mock the API request
      mockApiRequest.mockResolvedValueOnce(mockResponse)

      const url = '/api/test'
      const options = { method: 'GET' }

      const result = await tenantApiRequest(url, options)

      expect(mockApiRequest).toHaveBeenCalledWith(url, {
        ...options,
        headers: {
          'X-Tenant-ID': tenant.id
        }
      })
      expect(result).toEqual(mockResponse)
    })

    it('should make API request without tenant headers when no tenant selected', async () => {
      const mockResponse = { data: 'test' }
      mockApiRequest.mockResolvedValue(mockResponse)

      const { tenantApiRequest } = useTenant()

      const url = '/api/test'
      const options = { method: 'GET' }

      const result = await tenantApiRequest(url, options)

      expect(mockApiRequest).toHaveBeenCalledWith(url, {
        ...options,
        headers: {}
      })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('clearCurrentTenant', () => {
    it('should clear current tenant', async () => {
      const tenant = {
        id: 'tenant-123',
        discord_server_id: 'guild-123',
        name: 'Test Guild'
      }

      // Mock create tenant to set current tenant
      mockApiRequest.mockResolvedValue({ tenant })

      const { clearCurrentTenant, currentTenant, createTenant } = useTenant()
      
      // Set initial tenant by creating one
      await createTenant('guild-123')

      clearCurrentTenant()

      expect(currentTenant.value).toBeNull()
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('current_tenant')
    })
  })
})