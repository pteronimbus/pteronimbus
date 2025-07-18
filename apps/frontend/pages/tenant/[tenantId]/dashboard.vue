<template>
  <div>
    <!-- Tenant Header -->
    <div class="mb-8">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <UAvatar
            :src="getTenantIcon(currentTenant) || undefined"
            :alt="currentTenant?.name"
            size="lg"
          />
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {{ currentTenant?.name || 'Loading...' }}
            </h1>
          </div>
        </div>
        <div class="flex items-center space-x-3">
          <UButton
            variant="outline"
            size="sm"
            icon="heroicons:arrow-path"
            @click="refreshData"
            :loading="isRefreshing"
          >
            {{ t('common.refresh') }}
          </UButton>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center py-12">
      <UIcon name="heroicons:arrow-path" class="w-8 h-8 animate-spin mx-auto mb-4 text-primary-500" />
      <p class="text-gray-600 dark:text-gray-400">Loading dashboard data...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-center py-12">
      <UIcon name="heroicons:exclamation-triangle" class="w-12 h-12 text-red-500 mx-auto mb-4" />
      <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">Failed to Load Dashboard</h3>
      <p class="text-gray-600 dark:text-gray-400 mb-6">{{ error }}</p>
      <UButton @click="refreshData" variant="outline">
        <UIcon name="heroicons:arrow-path" class="w-4 h-4 mr-2" />
        Try Again
      </UButton>
    </div>

    <!-- Dashboard Content -->
    <div v-else>
      <!-- Stats Grid -->
      <div class="mb-8 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatsCard
          v-for="stat in tenantStats"
          :key="stat.key"
          :label="stat.label"
          :value="stat.value"
          :total="stat.total"
          :icon="stat.icon"
          :color="stat.color"
          :trend="stat.trend"
          :trend-color="stat.trendColor"
          :to="stat.route"
          :clickable="!!stat.route"
          class="cursor-pointer"
        />
      </div>

      <!-- Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Quick Actions -->
        <QuickActions 
          :title="t('dashboard.quickActions')"
          :actions="tenantQuickActions"
          :grid-cols="2"
        />

        <!-- Game Servers Status -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-800">
                {{ t('dashboard.gameServers') }}
              </h3>
              <UButton 
                color="primary" 
                variant="ghost" 
                size="sm"
                :to="`/tenant/${currentTenant?.id}/servers`"
              >
                {{ t('common.viewAll') }}
              </UButton>
            </div>
          </template>
          <div class="space-y-4">
            <div 
              v-for="server in gameServers.slice(0, 5)" 
              :key="server.id"
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div class="flex items-center space-x-3">
                <div class="flex-shrink-0">
                  <UIcon 
                    name="heroicons:server" 
                    class="w-5 h-5 text-gray-600"
                  />
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-900">{{ server.name }}</p>
                  <p class="text-xs text-gray-500">{{ server.game_type }}</p>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <StatusBadge :status="server.status.phase" />
                <span class="text-xs text-gray-500">
                  {{ server.status.player_count || 0 }} players
                </span>
              </div>
            </div>
            <EmptyState
              v-if="gameServers.length === 0"
              icon="heroicons:server"
              :title="t('dashboard.noGameServers')"
              :description="t('dashboard.noGameServersDesc')"
              size="sm"
            >
              <template #actions>
                <UButton
                  size="sm"
                  :to="`/tenant/${currentTenant?.id}/servers/create`"
                >
                  {{ t('servers.createServer') }}
                </UButton>
              </template>
            </EmptyState>
          </div>
        </UCard>

        <!-- Recent Activity -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-800">
                {{ t('dashboard.activity.title') }}
              </h3>
              <UButton 
                color="neutral" 
                variant="ghost" 
                size="sm"
                :to="`/tenant/${currentTenant?.id}/activity`"
              >
                {{ t('common.viewAll') }}
              </UButton>
            </div>
          </template>
          <div class="space-y-4">
            <div 
              v-for="activity in recentActivity.slice(0, 5)" 
              :key="activity.id"
              class="flex items-start space-x-3"
            >
              <div :class="[
                'flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center',
                getActivityColorClass(activity.type)
              ]">
                <UIcon 
                  :name="getActivityIcon(activity.type)" 
                  :class="[
                    'w-4 h-4',
                    getActivityIconColorClass(activity.type)
                  ]"
                />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm text-gray-800">{{ activity.message }}</p>
                <p class="text-xs text-gray-500">{{ formatTimestamp(activity.timestamp) }}</p>
              </div>
            </div>
            <EmptyState
              v-if="recentActivity.length === 0"
              icon="heroicons:clock"
              :title="t('dashboard.activity.noActivity')"
              :description="t('dashboard.activity.noActivityDesc')"
              size="sm"
            />
          </div>
        </UCard>

        <!-- Discord Integration Status -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-800">
                {{ t('dashboard.discordIntegration') }}
              </h3>
              <UButton 
                color="neutral" 
                variant="ghost" 
                size="sm"
                @click="syncDiscordData"
                :loading="isSyncing"
              >
                {{ t('dashboard.sync') }}
              </UButton>
            </div>
          </template>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <UIcon name="heroicons:users" class="w-5 h-5 text-blue-600" />
                <span class="text-sm text-gray-700">Discord Members</span>
              </div>
              <span class="text-sm font-medium">{{ discordStats.memberCount }}</span>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <UIcon name="heroicons:shield-check" class="w-5 h-5 text-green-600" />
                <span class="text-sm text-gray-700">Roles Synced</span>
              </div>
              <span class="text-sm font-medium">{{ discordStats.roleCount }}</span>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <UIcon name="heroicons:clock" class="w-5 h-5 text-gray-600" />
                <span class="text-sm text-gray-700">Last Sync</span>
              </div>
              <span class="text-sm text-gray-500">{{ formatTimestamp(discordStats.lastSync) }}</span>
            </div>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth', 'tenant']
})

const { t } = useI18n()
const router = useRouter()
const { 
  currentTenant, 
  tenantApiRequest,
  syncTenantData
} = useTenant()

// Local state
const isLoading = ref(true)
const isRefreshing = ref(false)
const isSyncing = ref(false)
const error = ref<string | null>(null)
const gameServers = ref<any[]>([])
const recentActivity = ref<any[]>([])
const discordStats = ref({
  memberCount: 0,
  roleCount: 0,
  lastSync: new Date().toISOString()
})

// Computed properties
const tenantStats = computed(() => [
  { 
    key: 'gameServers',
    label: t('dashboard.stats.gameServers'),
    value: gameServers.value.length.toString(),
    total: currentTenant.value?.config?.resource_limits?.max_game_servers?.toString() || '5',
    color: 'emerald',
    icon: 'heroicons:server',
    route: `/tenant/${currentTenant.value?.id}/servers`,
    trend: '+0',
    trendColor: 'text-gray-500'
  },
  { 
    key: 'activeServers',
    label: t('dashboard.stats.activeServers'),
    value: gameServers.value.filter(s => s.status.phase === 'Running').length.toString(),
    total: gameServers.value.length.toString(),
    color: 'blue',
    icon: 'heroicons:play-circle',
    route: `/tenant/${currentTenant.value?.id}/servers?status=running`,
    trend: '+0',
    trendColor: 'text-gray-500'
  },
  { 
    key: 'totalPlayers',
    label: t('dashboard.stats.totalPlayers'),
    value: gameServers.value.reduce((sum, s) => sum + (s.status.player_count || 0), 0).toString(),
    color: 'purple',
    icon: 'heroicons:users',
    route: `/tenant/${currentTenant.value?.id}/players`,
    trend: '+0',
    trendColor: 'text-gray-500'
  },
  { 
    key: 'discordMembers',
    label: t('dashboard.stats.discordMembers'),
    value: discordStats.value.memberCount.toString(),
    color: 'indigo',
    icon: 'heroicons:user-group',
    route: `/tenant/${currentTenant.value?.id}/members`,
    trend: '+0',
    trendColor: 'text-gray-500'
  }
])

const tenantQuickActions = computed(() => [
  {
    label: t('servers.createServer'),
    icon: 'heroicons:plus-circle',
    color: 'primary' as const,
    onClick: () => router.push(`/tenant/${currentTenant.value?.id}/servers/create`)
  },
  {
    label: t('dashboard.manageRoles'),
    icon: 'heroicons:shield-check',
    color: 'secondary' as const,
    onClick: () => router.push(`/tenant/${currentTenant.value?.id}/roles`)
  },
  {
    label: t('dashboard.viewLogs'),
    icon: 'heroicons:document-text',
    color: 'success' as const,
    onClick: () => router.push(`/tenant/${currentTenant.value?.id}/logs`)
  },
  {
    label: t('dashboard.settings'),
    icon: 'heroicons:cog-6-tooth',
    color: 'warning' as const,
    onClick: () => router.push(`/tenant/${currentTenant.value?.id}/settings`)
  }
])

// Methods
const loadDashboardData = async () => {
  if (!currentTenant.value) return

  try {
    error.value = null
    
    // Load game servers
    const serversResponse = await tenantApiRequest<{ servers: any[] }>(
      `/api/tenant/servers`
    )
    gameServers.value = serversResponse.servers || []

    // Load recent activity
    const activityResponse = await tenantApiRequest<{ activities: any[] }>(
      `/api/tenant/activity?limit=10`
    )
    recentActivity.value = activityResponse.activities || []

    // Load Discord stats
    const discordResponse = await tenantApiRequest<{ stats: any }>(
      `/api/tenant/discord/stats`
    )
    discordStats.value = {
      ...discordStats.value,
      ...discordResponse.stats
    }

  } catch (err: any) {
    console.error('Failed to load dashboard data:', err)
    error.value = err?.data?.message || 'Failed to load dashboard data'
  }
}

const refreshData = async () => {
  isRefreshing.value = true
  try {
    await loadDashboardData()
  } finally {
    isRefreshing.value = false
  }
}

const syncDiscordData = async () => {
  if (!currentTenant.value) return

  isSyncing.value = true
  try {
    await syncTenantData(currentTenant.value.id)
    await loadDashboardData()
    
    const toast = useToast()
    toast.add({
      title: 'Discord Data Synced',
      description: 'Discord roles and members have been synchronized',
      color: 'success'
    })
  } catch (err: any) {
    console.error('Failed to sync Discord data:', err)
    const toast = useToast()
    toast.add({
      title: 'Sync Failed',
      description: 'Failed to synchronize Discord data',
      color: 'error'
    })
  } finally {
    isSyncing.value = false
  }
}

// Helper functions
const getTenantIcon = (tenant: any) => {
  if (tenant?.icon) {
    return `https://cdn.discordapp.com/icons/${tenant.discord_server_id}/${tenant.icon}.png`
  }
  return null
}

const getActivityIcon = (type: string) => {
  const icons: Record<string, string> = {
    server_started: 'heroicons:play-circle',
    server_stopped: 'heroicons:stop-circle',
    server_created: 'heroicons:plus-circle',
    user_joined: 'heroicons:user-plus',
    user_left: 'heroicons:user-minus',
    role_updated: 'heroicons:shield-check'
  }
  return icons[type] || 'heroicons:information-circle'
}

const getActivityColorClass = (type: string) => {
  const colors: Record<string, string> = {
    server_started: 'bg-green-100',
    server_stopped: 'bg-red-100',
    server_created: 'bg-blue-100',
    user_joined: 'bg-purple-100',
    user_left: 'bg-gray-100',
    role_updated: 'bg-yellow-100'
  }
  return colors[type] || 'bg-gray-100'
}

const getActivityIconColorClass = (type: string) => {
  const colors: Record<string, string> = {
    server_started: 'text-green-600',
    server_stopped: 'text-red-600',
    server_created: 'text-blue-600',
    user_joined: 'text-purple-600',
    user_left: 'text-gray-600',
    role_updated: 'text-yellow-600'
  }
  return colors[type] || 'text-gray-600'
}

const formatTimestamp = (timestamp: string) => {
  return new Date(timestamp).toLocaleString()
}

// Load data on mount
onMounted(async () => {
  isLoading.value = true
  try {
    await loadDashboardData()
  } finally {
    isLoading.value = false
  }
})

// Watch for tenant changes
watch(currentTenant, async (newTenant) => {
  if (newTenant) {
    await refreshData()
  }
})
</script>