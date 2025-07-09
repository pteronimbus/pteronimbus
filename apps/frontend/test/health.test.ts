import { describe, it, expect } from 'vitest'
import { setup, $fetch } from '@nuxt/test-utils'

interface HealthResponse {
  status: string
  timestamp: string
  service: string
  version?: string
}

describe('Health Check Endpoints', async () => {
  await setup({
    // Test context is automatically handled by @nuxt/test-utils
  })

  it('should respond to /api/health with healthy status', async () => {
    const data: HealthResponse = await $fetch('/api/health')
    
    expect(data.status).toBe('healthy')
    expect(data.service).toBe('pteronimbus-frontend')
    expect(data.version).toBe('0.1.0')
    expect(data.timestamp).toBeDefined()
    
    // Verify timestamp is recent (within last 10 seconds)
    const timestamp = new Date(data.timestamp)
    const now = new Date()
    const timeDiff = now.getTime() - timestamp.getTime()
    expect(timeDiff).toBeLessThan(10000)
  })

  it('should respond to /api/healthz with healthy status', async () => {
    const data: HealthResponse = await $fetch('/api/healthz')
    
    expect(data.status).toBe('healthy')
    expect(data.service).toBe('pteronimbus-frontend')
    expect(data.version).toBe('0.1.0')
    expect(data.timestamp).toBeDefined()
  })

  it('should respond to /api/ready with ready status', async () => {
    const data: HealthResponse = await $fetch('/api/ready')
    
    expect(data.status).toBe('ready')
    expect(data.service).toBe('pteronimbus-frontend')
    expect(data.timestamp).toBeDefined()
    expect(data.version).toBeUndefined()
  })

  it('should respond to /api/live with alive status', async () => {
    const data: HealthResponse = await $fetch('/api/live')
    
    expect(data.status).toBe('alive')
    expect(data.service).toBe('pteronimbus-frontend')
    expect(data.timestamp).toBeDefined()
    expect(data.version).toBeUndefined()
  })

  it('should have consistent timestamp format across all endpoints', async () => {
    const endpoints = ['/api/health', '/api/healthz', '/api/ready', '/api/live']
    
    for (const endpoint of endpoints) {
      const data: HealthResponse = await $fetch(endpoint)
      expect(data.timestamp).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/)
    }
  })
}) 