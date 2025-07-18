<template>
  <div>
    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <UButton
            variant="ghost"
            icon="heroicons:arrow-left"
            @click="goBack"
          >
            {{ t('common.back') }}
          </UButton>
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {{ server?.name || 'Loading...' }}
            </h1>
            <p class="text-gray-600 dark:text-gray-400">
              {{ server?.game_type }} • {{ server?.version || 'Latest' }}
            </p>
          </div>
        </div>
        <div class="flex items-center space-x-3">
          <StatusBadge v-if="server" :status="server.status.phase" />
          <UDropdown v-if="server" :items="serverActions">
            <UButton
              color="gray"
              variant="outline"
              icon="i-heroicons-ellipsis-vertical-20-solid"
            >
              {{ t('common.actions') }}
            </UButton>
          </UDropdown>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center py-12">
      <UIcon name="heroicons:arrow-path" class="w-8 h-8 animate-spin mx-auto mb-4 text-primary-500" />
      <p class="text-gray-600 dark:text-gray-400">Loading server details...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-center py-12">
      <UIcon name="heroicons:exclamation-triangle" class="w-12 h-12 text-red-500 mx-auto mb-4" />
      <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">Failed to Load Server</h3>
      <p class="text-gray-600 dark:text-gray-400 mb-6">{{ error }}</p>
      <UButton @click="loadServer" variant="outline">
        <UIcon name="heroicons:arrow-path" class="w-4 h-4 mr-2" />
        Try Again
      </UButton>
    </div>

    <!-- Server Details -->
    <div v-else-if="server" class="space-y-6">
      <!-- Stats Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatsCard
          :label="t('servers.stats.status')"
          :value="server.status.phase"
          icon="heroicons:signal"
          :color="getStatusColor(server.status.phase)"
        />
        <StatsCard
          :label="t('servers.stats.players')"
          :value="`${server.status.player_count || 0}/${server.max_players || 20}`"
          icon="heroicons:users"
          color="blue"
        />
        <StatsCard
          :label="t('servers.stats.uptime')"
          :value="formatUptime(server.status.uptime)"
          icon="heroicons:clock"
          color="green"
        />
        <StatsCard
          :label="t('servers.stats.memory')"
          :value="formatMemory(server.status.memory_usage)"
          icon="heroicons:cpu-chip"
          color="purple"
        />
      </div>

      <!-- Tabs -->
      <UTabs v-model="activeTab" :items="tabs">
        <!-- Overview Tab -->
        <template #overview>
          <div class="space-y-6">
            <!-- Server Information -->
            <UCard>
              <template #header>
                <h3 class="text-lg font-semibold">{{ t('servers.details.information') }}</h3>
              </template>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    {{ t('servers.fields.name') }}
                  </label>
                  <p class="text-gray-900 dark:text-gray-100">{{ server.name }}</p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    {{ t('servers.fields.gameType') }}
                  </label>
                  <p class="text-gray-900 dark:text-gray-100">{{ server.game_type }}</p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    {{ t('servers.fields.version') }}
                  </label>
                  <p class="text-gray-900 dark:text-gray-100">{{ server.version || 'Latest' }}</p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    {{ t('servers.fields.maxPlayers') }}
                  </label>
                  <p class="text-gray-900 dark:text-gray-100">{{ server.max_players || 20 }}</p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    {{ t('servers.fields.created') }}
                  </label>
                  <p class="text-gray-900 dark:text-gray-100">{{ formatDate(server.created_at) }}</p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    {{ t('servers.fields.lastUpdated') }}
                  </label>
                  <p class="text-gray-900 dark:text-gray-100">{{ formatDate(server.updated_at) }}</p>
                </div>
              </div>
            </UCard>

            <!-- Resource Usage -->
            <UCard>
              <template #header>
                <h3 class="text-lg font-semibold">{{ t('servers.details.resources') }}</h3>
              </template>
              <div class="space-y-4">
                <div>
                  <div class="flex justify-between text-sm mb-1">
                    <span>{{ t('servers.stats.cpu') }}</span>
                    <span>{{ formatCpu(server.status.cpu_usage) }}</span>
                  </div>
                  <UProgress :value="server.status.cpu_usage || 0" color="blue" />
                </div>
                <div>
                  <div class="flex justify-between text-sm mb-1">
                    <span>{{ t('servers.stats.memory') }}</span>
                    <span>{{ formatMemory(server.status.memory_usage) }}</span>
                  </div>
                  <UProgress :value="(server.status.memory_usage || 0) / 1024 / 1024 / 10" color="green" />
                </div>
              </div>
            </UCard>
          </div>
        </template>

        <!-- Console Tab -->
        <template #console>
          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold">{{ t('servers.console.title') }}</h3>
                <div class="flex items-center space-x-2">
                  <UButton
                    variant="ghost"
                    size="sm"
                    icon="heroicons:arrow-path"
                    @click="refreshLogs"
                    :loading="isRefreshingLogs"
                  >
                    {{ t('common.refresh') }}
                  </UButton>
                  <UButton
                    variant="ghost"
                    size="sm"
                    icon="heroicons:arrow-down-tray"
                    @click="downloadLogs"
                  >
                    {{ t('servers.console.download') }}
                  </UButton>
                </div>
              </div>
            </template>
            <div class="space-y-4">
              <!-- Console Output -->
              <div class="bg-gray-900 text-green-400 p-4 rounded-lg font-mono text-sm h-96 overflow-y-auto">
                <div v-for="(line, index) in consoleLogs" :key="index" class="whitespace-pre-wrap">
                  {{ line }}
                </div>
                <div v-if="consoleLogs.length === 0" class="text-gray-500">
                  No logs available
                </div>
              </div>
              
              <!-- Command Input -->
              <div class="flex space-x-2">
                <UInput
                  v-model="consoleCommand"
                  placeholder="Enter command..."
                  class="flex-1"
                  @keyup.enter="sendCommand"
                />
                <UButton
                  @click="sendCommand"
                  :disabled="!consoleCommand.trim()"
                >
                  {{ t('servers.console.send') }}
                </UButton>
              </div>
            </div>
          </UCard>
        </template>

        <!-- Files Tab -->
        <template #files>
          <UCard>
            <template #header>
              <h3 class="text-lg font-semibold">{{ t('servers.files.title') }}</h3>
            </template>
            <div class="text-center py-12">
              <UIcon name="heroicons:folder" class="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <p class="text-gray-600 dark:text-gray-400">File management coming soon</p>
            </div>
          </UCard>
        </template>

        <!-- Settings Tab -->
        <template #settings>
          <UCard>
            <template #header>
              <h3 class="text-lg font-semibold">{{ t('servers.settings.title') }}</h3>
            </template>
            <div class="text-center py-12">
              <UIcon name="heroicons:cog-6-tooth" class="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <p class="text-gray-600 dark:text-gray-400">Server settings coming soon</p>
            </div>
          </UCard>
        </template>
      </UTabs>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth', 'tenant']
})

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const { currentTenant, tenantApiRequest } = useTenant()

// Get server ID from route
const serverId = computed(() => route.params.id as string)

// Types
interface Server {
  id: string
  name: string
  game_type: string
  version?: string
  max_players?: number
  status: {
    phase: string
    player_count?: number
    uptime?: string
    memory_usage?: number
    cpu_usage?: number
  }
  created_at: string
  updated_at: string
}

// State
const isLoading = ref(true)
const error = ref<string | null>(null)
const server = ref<Server | null>(null)
const activeTab = ref(0)
const consoleLogs = ref<string[]>([])
const consoleCommand = ref('')
const isRefreshingLogs = ref(false)

// Tabs configuration
const tabs = [
  { key: 'overview', label: t('servers.tabs.overview') },
  { key: 'console', label: t('servers.tabs.console') },
  { key: 'files', label: t('servers.tabs.files') },
  { key: 'settings', label: t('servers.tabs.settings') }
]

// Server actions dropdown
const serverActions = computed(() => {
  if (!server.value) return []
  
  return [
    [{
      label: server.value.status.phase === 'Running' ? t('servers.actions.stop') : t('servers.actions.start'),
      icon: server.value.status.phase === 'Running' ? 'i-heroicons-stop-circle-20-solid' : 'i-heroicons-play-circle-20-solid',
      click: () => toggleServer()
    }, {
      label: t('servers.actions.restart'),
      icon: 'i-heroicons-arrow-path-20-solid',
      click: () => restartServer()
    }],
    [{
      label: t('servers.actions.edit'),
      icon: 'i-heroicons-pencil-square-20-solid',
      click: () => router.push(`/tenant/${currentTenant.value?.id}/servers/${serverId.value}/edit`)
    }],
    [{
      label: t('servers.actions.delete'),
      icon: 'i-heroicons-trash-20-solid',
      click: () => deleteServer()
    }]
  ]
})

// Methods
const loadServer = async () => {
  if (!currentTenant.value || !serverId.value) return

  try {
    isLoading.value = true
    error.value = null

    const response = await tenantApiRequest<{ server: Server }>(`/api/tenant/servers/${serverId.value}`)
    server.value = response.server
  } catch (err: any) {
    console.error('Failed to load server:', err)
    error.value = err?.data?.message || 'Failed to load server'
  } finally {
    isLoading.value = false
  }
}

const goBack = () => {
  router.push(`/tenant/${currentTenant.value?.id}/servers`)
}

const toggleServer = async () => {
  if (!server.value) return

  try {
    const action = server.value.status.phase === 'Running' ? 'stop' : 'start'
    await tenantApiRequest(`/api/tenant/servers/${serverId.value}/${action}`, {
      method: 'POST'
    })
    
    const toast = useToast()
    toast.add({
      title: `Server ${action === 'start' ? 'Started' : 'Stopped'}`,
      description: `${server.value.name} has been ${action === 'start' ? 'started' : 'stopped'}`,
      color: 'success'
    })
    
    await loadServer()
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Action Failed',
      description: err?.data?.message || 'Failed to perform server action',
      color: 'error'
    })
  }
}

const restartServer = async () => {
  if (!server.value) return

  try {
    await tenantApiRequest(`/api/tenant/servers/${serverId.value}/restart`, {
      method: 'POST'
    })
    
    const toast = useToast()
    toast.add({
      title: 'Server Restarted',
      description: `${server.value.name} has been restarted`,
      color: 'success'
    })
    
    await loadServer()
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Restart Failed',
      description: err?.data?.message || 'Failed to restart server',
      color: 'error'
    })
  }
}

const deleteServer = async () => {
  if (!server.value) return

  const confirmed = confirm(`Are you sure you want to delete ${server.value.name}? This action cannot be undone.`)
  if (!confirmed) return

  try {
    await tenantApiRequest(`/api/tenant/servers/${serverId.value}`, {
      method: 'DELETE'
    })
    
    const toast = useToast()
    toast.add({
      title: 'Server Deleted',
      description: `${server.value.name} has been deleted`,
      color: 'success'
    })
    
    router.push(`/tenant/${currentTenant.value?.id}/servers`)
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Delete Failed',
      description: err?.data?.message || 'Failed to delete server',
      color: 'error'
    })
  }
}

const refreshLogs = async () => {
  if (!serverId.value) return

  try {
    isRefreshingLogs.value = true
    const response = await tenantApiRequest<{ logs: string[] }>(`/api/tenant/servers/${serverId.value}/logs`)
    consoleLogs.value = response.logs || []
  } catch (err: any) {
    console.error('Failed to refresh logs:', err)
    const toast = useToast()
    toast.add({
      title: 'Failed to Load Logs',
      description: err?.data?.message || 'Failed to load server logs',
      color: 'error'
    })
  } finally {
    isRefreshingLogs.value = false
  }
}

const sendCommand = async () => {
  if (!consoleCommand.value.trim() || !serverId.value) return

  try {
    await tenantApiRequest(`/api/tenant/servers/${serverId.value}/command`, {
      method: 'POST',
      body: { command: consoleCommand.value }
    })
    
    // Add command to logs
    consoleLogs.value.push(`> ${consoleCommand.value}`)
    consoleCommand.value = ''
    
    // Refresh logs after a short delay
    setTimeout(() => {
      refreshLogs()
    }, 1000)
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Command Failed',
      description: err?.data?.message || 'Failed to send command',
      color: 'error'
    })
  }
}

const downloadLogs = () => {
  if (consoleLogs.value.length === 0) return

  const logContent = consoleLogs.value.join('\n')
  const blob = new Blob([logContent], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${server.value?.name || 'server'}-logs.txt`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

// Utility functions
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    'Running': 'green',
    'Stopped': 'red',
    'Pending': 'yellow',
    'Failed': 'red'
  }
  return colors[status] || 'gray'
}

const formatUptime = (uptime?: string) => {
  if (!uptime) return 'N/A'
  return uptime
}

const formatMemory = (memory?: number) => {
  if (!memory) return 'N/A'
  return `${Math.round(memory / 1024 / 1024)}MB`
}

const formatCpu = (cpu?: number) => {
  if (!cpu) return 'N/A'
  return `${Math.round(cpu)}%`
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString()
}

// Load server on mount
onMounted(() => {
  loadServer()
  
  // Load logs if console tab is active
  if (activeTab.value === 1) {
    refreshLogs()
  }
})

// Watch for tab changes
watch(activeTab, (newTab) => {
  if (newTab === 1 && consoleLogs.value.length === 0) {
    refreshLogs()
  }
})

// Watch for tenant changes
watch(currentTenant, () => {
  if (currentTenant.value) {
    loadServer()
  }
})
</script>