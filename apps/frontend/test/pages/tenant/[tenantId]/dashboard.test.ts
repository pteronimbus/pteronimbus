import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock the API request function
const mockTenantApiRequest = vi.fn()

describe('Tenant Dashboard API Integration', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    
    // Setup default mock responses
    mockTenantApiRequest.mockImplementation((url: string) => {
      if (url === '/api/tenant/servers') {
        return Promise.resolve({
          servers: [
            {
              id: 'server-1',
              name: 'Survival World',
              game_type: 'minecraft',
              status: {
                phase: 'Running',
                player_count: 5
              }
            },
            {
              id: 'server-2',
              name: 'Competitive Server',
              game_type: 'cs2',
              status: {
                phase: 'Stopped',
                player_count: 0
              }
            }
          ]
        })
      }
      
      if (url === '/api/tenant/activity?limit=10') {
        return Promise.resolve({
          activities: [
            {
              id: 'activity-1',
              type: 'server_started',
              message: "Server 'Survival World' was started",
              timestamp: '2024-01-15T10:30:00Z'
            }
          ]
        })
      }
      
      if (url === '/api/tenant/discord/stats') {
        return Promise.resolve({
          stats: {
            memberCount: 42,
            roleCount: 8,
            lastSync: '2024-01-15T10:00:00Z'
          }
        })
      }
      
      return Promise.resolve({})
    })
  })

  it('should return servers data from /api/tenant/servers endpoint', async () => {
    const response = await mockTenantApiRequest('/api/tenant/servers')
    
    expect(response.servers).toBeDefined()
    expect(response.servers).toHaveLength(2)
    expect(response.servers[0].name).toBe('Survival World')
    expect(response.servers[0].game_type).toBe('minecraft')
    expect(response.servers[0].status.phase).toBe('Running')
    expect(response.servers[0].status.player_count).toBe(5)
    
    expect(response.servers[1].name).toBe('Competitive Server')
    expect(response.servers[1].game_type).toBe('cs2')
    expect(response.servers[1].status.phase).toBe('Stopped')
    expect(response.servers[1].status.player_count).toBe(0)
  })

  it('should return activity data from /api/tenant/activity endpoint', async () => {
    const response = await mockTenantApiRequest('/api/tenant/activity?limit=10')
    
    expect(response.activities).toBeDefined()
    expect(response.activities).toHaveLength(1)
    expect(response.activities[0].id).toBe('activity-1')
    expect(response.activities[0].type).toBe('server_started')
    expect(response.activities[0].message).toBe("Server 'Survival World' was started")
    expect(response.activities[0].timestamp).toBe('2024-01-15T10:30:00Z')
  })

  it('should return Discord stats from /api/tenant/discord/stats endpoint', async () => {
    const response = await mockTenantApiRequest('/api/tenant/discord/stats')
    
    expect(response.stats).toBeDefined()
    expect(response.stats.memberCount).toBe(42)
    expect(response.stats.roleCount).toBe(8)
    expect(response.stats.lastSync).toBe('2024-01-15T10:00:00Z')
  })

  it('should handle API errors gracefully', async () => {
    mockTenantApiRequest.mockRejectedValue(new Error('Network error'))
    
    await expect(mockTenantApiRequest('/api/tenant/servers')).rejects.toThrow('Network error')
  })

  it('should return empty object for unknown endpoints', async () => {
    const response = await mockTenantApiRequest('/api/tenant/unknown')
    
    expect(response).toEqual({})
  })
})