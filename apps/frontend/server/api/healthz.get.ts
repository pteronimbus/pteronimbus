interface HealthResponse {
  status: string
  timestamp: string
  service: string
  version?: string
}

export default defineEventHandler(async (event): Promise<HealthResponse> => {
  return {
    status: 'healthy',
    timestamp: new Date().toISOString(),
    service: 'pteronimbus-frontend',
    version: '0.1.0'
  }
}) 