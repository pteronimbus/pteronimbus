<script setup lang="ts">
definePageMeta({
  middleware: ['auth', 'tenant']
})

const { t } = useI18n()
const router = useRouter()
const { 
  currentTenant, 
  tenantApiRequest,
  syncTenantData,
  clearCurrentTenant
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
const botStatus = ref<{ present: boolean, missingPermissions: string[] } | null>(null)
const botModalOpen = ref(false)
const pollingInterval = ref<NodeJS.Timeout | null>(null)

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
    description: 'Launch a new game server',
    icon: 'heroicons:plus-circle',
    color: 'primary' as const,
    onClick: () => router.push(`/tenant/${currentTenant.value?.id}/servers/create`)
  },
  {
    label: t('dashboard.manageRoles'),
    description: 'Configure user permissions',
    icon: 'heroicons:shield-check',
    color: 'secondary' as const,
    onClick: () => router.push(`/tenant/${currentTenant.value?.id}/roles`)
  },
  {
    label: t('dashboard.viewLogs'),
    description: 'Monitor server activity',
    icon: 'heroicons:document-text',
    color: 'success' as const,
    onClick: () => router.push(`/tenant/${currentTenant.value?.id}/logs`)
  },
  {
    label: t('dashboard.settings'),
    description: 'Server configuration',
    icon: 'heroicons:cog-6-tooth',
    color: 'warning' as const,
    onClick: () => router.push(`/tenant/${currentTenant.value?.id}/settings`)
  }
])

const runtimeConfig = useRuntimeConfig()
const inviteUrl = computed(() => {
  if (!currentTenant.value) return '#'
  const clientId = runtimeConfig.public.discordClientId
  const guildId = currentTenant.value.discord_server_id
  const perms = 277293902870 // MANAGE_ROLES + SEND_MESSAGES + VIEW_CHANNEL
  return `https://discord.com/oauth2/authorize?client_id=${clientId}&scope=bot+applications.commands&permissions=${perms}&guild_id=${guildId}`
})

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
    // Enhanced error handling for forbidden/not found
    const code = err?.data?.code || err?.response?.data?.code
    if (code === 'FORBIDDEN' || code === 'NOT_FOUND') {
      clearCurrentTenant()
      await router.replace('/tenants')
      return
    }
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

const checkBotStatus = async () => {
  if (!currentTenant.value) return
  try {
    const res = await tenantApiRequest<{ present: boolean, missingPermissions: string[] }>(`/api/tenants/${currentTenant.value.id}/bot-status`)
    botStatus.value = res
    botModalOpen.value = !res.present || (res.missingPermissions && res.missingPermissions.length > 0)
  } catch (err) {
    botStatus.value = { present: false, missingPermissions: ["unknown"] }
    botModalOpen.value = true
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
    await checkBotStatus()
  } finally {
    isLoading.value = false
  }
})

// Watch for tenant changes
watch(currentTenant, async (newTenant) => {
  if (newTenant) {
    await refreshData()
    await checkBotStatus()
  }
})

const startBotStatusPolling = () => {
  if (pollingInterval.value) return // already polling
  pollingInterval.value = setInterval(async () => {
    await checkBotStatus()
    // If bot is present, stop polling
    if (botStatus.value && botStatus.value.present && (!botStatus.value.missingPermissions || botStatus.value.missingPermissions.length === 0)) {
      stopBotStatusPolling()
    }
  }, 2000)
}

const stopBotStatusPolling = () => {
  if (pollingInterval.value) {
    clearInterval(pollingInterval.value)
    pollingInterval.value = null
  }
}

watch(
  () => botModalOpen.value,
  (open) => {
    if (open) {
      startBotStatusPolling()
    } else {
      stopBotStatusPolling()
    }
  }
)

onUnmounted(() => {
  stopBotStatusPolling()
})
</script>

<template>
  <div class="space-y-8">
    <UModal v-model:open="botModalOpen"
      :dismissible="false"
      :close="false"
      :title="t('dashboard.botRequiredTitle', 'Bot Required')"
      :ui="{
        overlay: 'fixed inset-0 bg-gray-200/75 dark:bg-gray-900/75 backdrop-blur-sm'
      }"
    >
      <template #body>
        <div class="text-center py-8">
          <UIcon name="i-heroicons-wrench-screwdriver" class="w-12 h-12 text-primary-500 mx-auto mb-4" />
          <p class="text-sm text-gray-500 dark:text-gray-400 mb-6">
            {{ t('dashboard.botRequiredDescription', 'To use all features, please invite the Pteronimbus bot to your Discord server.') }}
          </p>
          <UButton size="lg" color="primary" :to="inviteUrl" target="_blank">
            {{ t('dashboard.inviteBot', 'Invite Bot') }}
          </UButton>
          <div v-if="botStatus && botStatus.missingPermissions && botStatus.missingPermissions.length > 0" class="mt-4">
            <p class="text-sm text-red-500">Missing permissions: {{ botStatus.missingPermissions.join(', ') }}</p>
          </div>
        </div>
      </template>
    </UModal>
    <!-- Tenant Header - Enhanced -->
    <div class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6 shadow-sm">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <div class="relative">
            <UAvatar
              :src="getTenantIcon(currentTenant) || undefined"
              :alt="currentTenant?.name"
              size="xl"
              class="ring-4 ring-white dark:ring-gray-800 shadow-lg"
            />
            <div class="absolute -bottom-1 -right-1 w-6 h-6 bg-green-500 border-2 border-white dark:border-gray-800 rounded-full flex items-center justify-center">
              <UIcon name="heroicons:check" class="w-3 h-3 text-white" />
            </div>
          </div>
          <div>
            <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-1">
              {{ currentTenant?.name || 'Loading...' }}
            </h1>
            <p class="text-gray-600 dark:text-gray-400 flex items-center">
              <UIcon name="heroicons:server" class="w-4 h-4 mr-1" />
              Discord Server • Active
            </p>
          </div>
        </div>
        <div class="flex items-center space-x-3">
          <UButton
            variant="outline"
            size="sm"
            icon="heroicons:arrow-path"
            @click="refreshData"
            :loading="isRefreshing"
            class="bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors"
          >
            {{ t('common.refresh') }}
          </UButton>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center py-16">
      <div class="inline-flex items-center justify-center w-16 h-16 bg-primary-100 dark:bg-primary-900 rounded-full mb-4">
        <UIcon name="heroicons:arrow-path" class="w-8 h-8 animate-spin text-primary-600 dark:text-primary-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">Loading Dashboard</h3>
      <p class="text-gray-600 dark:text-gray-400">Gathering your server data...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-center py-16">
      <div class="inline-flex items-center justify-center w-16 h-16 bg-red-100 dark:bg-red-900 rounded-full mb-4">
        <UIcon name="heroicons:exclamation-triangle" class="w-8 h-8 text-red-600 dark:text-red-400" />
      </div>
      <h3 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Failed to Load Dashboard</h3>
      <p class="text-gray-600 dark:text-gray-400 mb-6 max-w-md mx-auto">{{ error }}</p>
      <UButton @click="refreshData" variant="outline" size="lg">
        <UIcon name="heroicons:arrow-path" class="w-4 h-4 mr-2" />
        Try Again
      </UButton>
    </div>

    <!-- Dashboard Content -->
    <div v-else class="space-y-8">
      <!-- Stats Grid - Enhanced -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <div 
          v-for="stat in tenantStats" 
          :key="stat.key"
          class="group"
        >
          <StatsCard
            :label="stat.label"
            :value="stat.value"
            :total="stat.total"
            :icon="stat.icon"
            :color="stat.color"
            :trend="stat.trend"
            :trend-color="stat.trendColor"
            :to="stat.route"
            :clickable="!!stat.route"
            class="h-full group-hover:shadow-lg group-hover:-translate-y-0.5 transition-all duration-200"
          />
        </div>
      </div>

      <!-- Content Grid - Enhanced -->
      <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
        <!-- Quick Actions - Enhanced -->
        <div class="space-y-6">
          <QuickActions 
            :title="t('dashboard.quickActions')"
            :actions="tenantQuickActions"
            :grid-cols="2"
            class="shadow-sm hover:shadow-md transition-shadow duration-200"
          />
        </div>

        <!-- Game Servers Status - Enhanced -->
        <UCard class="shadow-sm hover:shadow-md transition-shadow duration-200">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <div class="w-2 h-2 bg-emerald-500 rounded-full"></div>
                <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  {{ t('dashboard.gameServers') }}
                </h3>
                                 <UBadge color="success" variant="soft" size="sm">
                   {{ gameServers.length }}
                 </UBadge>
               </div>
               <UButton 
                 color="neutral" 
                variant="ghost" 
                size="sm"
                :to="`/tenant/${currentTenant?.id}/servers`"
                class="hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
              >
                {{ t('common.viewAll') }}
                <UIcon name="heroicons:arrow-right" class="w-4 h-4 ml-1" />
              </UButton>
            </div>
          </template>
          <div class="space-y-3">
            <div 
              v-for="server in gameServers.slice(0, 5)" 
              :key="server.id"
              class="group flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-transparent hover:border-gray-200 dark:hover:border-gray-700 hover:bg-white dark:hover:bg-gray-800 transition-all duration-200 cursor-pointer"
            >
              <div class="flex items-center space-x-3 flex-1 min-w-0">
                <div class="flex-shrink-0 w-10 h-10 bg-white dark:bg-gray-700 rounded-lg flex items-center justify-center shadow-sm group-hover:shadow-md transition-shadow">
                  <UIcon 
                    name="heroicons:server" 
                    class="w-5 h-5 text-gray-600 dark:text-gray-400"
                  />
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{{ server.name }}</p>
                  <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wide">{{ server.game_type }}</p>
                </div>
              </div>
              <div class="flex items-center space-x-3 flex-shrink-0">
                <StatusBadge :status="server.status.phase" />
                <div class="text-right">
                  <p class="text-xs font-medium text-gray-900 dark:text-gray-100">
                    {{ server.status.player_count || 0 }}
                  </p>
                  <p class="text-xs text-gray-500 dark:text-gray-400">players</p>
                </div>
              </div>
            </div>
            
            <EmptyState
              v-if="gameServers.length === 0"
              icon="heroicons:server"
              :title="t('dashboard.noGameServers')"
              :description="t('dashboard.noGameServersDesc')"
              size="sm"
              class="py-8"
            >
              <template #actions>
                <UButton
                  size="sm"
                  :to="`/tenant/${currentTenant?.id}/servers/create`"
                  class="shadow-sm"
                >
                  <UIcon name="heroicons:plus" class="w-4 h-4 mr-1" />
                  {{ t('servers.createServer') }}
                </UButton>
              </template>
            </EmptyState>
          </div>
        </UCard>

        <!-- Recent Activity - Enhanced -->
        <UCard class="shadow-sm hover:shadow-md transition-shadow duration-200">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <div class="w-2 h-2 bg-blue-500 rounded-full"></div>
                <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  {{ t('dashboard.activity.title') }}
                </h3>
              </div>
                             <UButton 
                 color="neutral" 
                 variant="ghost" 
                 size="sm"
                 :to="`/tenant/${currentTenant?.id}/activity`"
                 class="hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
               >
                {{ t('common.viewAll') }}
                <UIcon name="heroicons:arrow-right" class="w-4 h-4 ml-1" />
              </UButton>
            </div>
          </template>
          <div class="space-y-4">
            <div 
              v-for="activity in recentActivity.slice(0, 5)" 
              :key="activity.id"
              class="flex items-start space-x-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors duration-150"
            >
              <div :class="[
                'flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center shadow-sm',
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
                <p class="text-sm text-gray-900 dark:text-gray-100 leading-relaxed">{{ activity.message }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ formatTimestamp(activity.timestamp) }}</p>
              </div>
            </div>
            
            <EmptyState
              v-if="recentActivity.length === 0"
              icon="heroicons:clock"
              :title="t('dashboard.activity.noActivity')"
              :description="t('dashboard.activity.noActivityDesc')"
              size="sm"
              class="py-6"
            />
          </div>
        </UCard>

        <!-- Discord Integration Status - Enhanced -->
        <UCard class="shadow-sm hover:shadow-md transition-shadow duration-200">
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <div class="w-2 h-2 bg-indigo-500 rounded-full"></div>
                <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  {{ t('dashboard.discordIntegration') }}
                </h3>
                                 <UBadge color="info" variant="soft" size="sm">
                   Connected
                 </UBadge>
               </div>
               <UButton 
                 color="info"
                variant="ghost" 
                size="sm"
                @click="syncDiscordData"
                :loading="isSyncing"
                class="hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors px-3"
              >
                <UIcon name="heroicons:arrow-path" class="w-4 h-4 mr-1" />
                {{ t('dashboard.sync') }}
              </UButton>
            </div>
          </template>
          <div class="space-y-4">
            <div class="flex items-center justify-between p-3 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
              <div class="flex items-center space-x-3">
                <div class="w-8 h-8 bg-blue-100 dark:bg-blue-900 rounded-lg flex items-center justify-center">
                  <UIcon name="heroicons:users" class="w-4 h-4 text-blue-600 dark:text-blue-400" />
                </div>
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100">Discord Members</span>
              </div>
              <span class="text-lg font-bold text-blue-600 dark:text-blue-400">{{ discordStats.memberCount }}</span>
            </div>
            
            <div class="flex items-center justify-between p-3 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800">
              <div class="flex items-center space-x-3">
                <div class="w-8 h-8 bg-green-100 dark:bg-green-900 rounded-lg flex items-center justify-center">
                  <UIcon name="heroicons:shield-check" class="w-4 h-4 text-green-600 dark:text-green-400" />
                </div>
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100">Roles Synced</span>
              </div>
              <span class="text-lg font-bold text-green-600 dark:text-green-400">{{ discordStats.roleCount }}</span>
            </div>
            
            <div class="flex items-center justify-between p-3 rounded-lg bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
              <div class="flex items-center space-x-3">
                <div class="w-8 h-8 bg-gray-100 dark:bg-gray-700 rounded-lg flex items-center justify-center">
                  <UIcon name="heroicons:clock" class="w-4 h-4 text-gray-600 dark:text-gray-400" />
                </div>
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100">Last Sync</span>
              </div>
              <span class="text-sm text-gray-600 dark:text-gray-400">{{ formatTimestamp(discordStats.lastSync) }}</span>
            </div>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>