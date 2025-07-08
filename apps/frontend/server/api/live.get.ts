interface HealthResponse {
  status: string
  timestamp: string
  service: string
}

export default defineEventHandler(async (event): Promise<HealthResponse> => {
  return {
    status: 'alive',
    timestamp: new Date().toISOString(),
    service: 'pteronimbus-frontend'
  }
}) 