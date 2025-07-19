export interface Controller {
  id: string
  cluster_id: string
  cluster_name: string
  version: string
  status: string
  last_heartbeat: string
  is_online: boolean
  uptime: string
  approved_at?: string
  approved_by?: string
  created_at: string
} 