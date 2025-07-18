<template>
  <div>
    <PageHeader
      :title="t('servers.title')"
      :description="t('servers.description')"
    />

    <div class="space-y-6">
      <!-- Search and Filters -->
      <SearchAndFilters
        v-model:search="searchQuery"
        v-model:filters="activeFilters"
        :filter-options="filterOptions"
        :placeholder="t('servers.search.placeholder')"
        @clear-filters="clearFilters"
      />

      <!-- Servers Grid/List -->
      <UCard>
        <div v-if="isLoading" class="text-center py-12">
          <UIcon name="heroicons:arrow-path" class="w-8 h-8 animate-spin mx-auto mb-4 text-primary-500" />
          <p class="text-gray-600 dark:text-gray-400">{{ t('servers.loading') }}</p>
        </div>

        <div v-else-if="error" class="text-center py-12">
          <UIcon name="heroicons:exclamation-triangle" class="w-12 h-12 text-red-500 mx-auto mb-4" />
          <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">{{ t('common.error') }}</h3>
          <p class="text-gray-600 dark:text-gray-400 mb-6">{{ error }}</p>
          <UButton @click="loadServers" variant="outline">
            <UIcon name="heroicons:arrow-path" class="w-4 h-4 mr-2" />
            {{ t('common.tryAgain') }}
          </UButton>
        </div>

        <div v-else-if="filteredServers.length === 0">
          <EmptyState
            v-if="!hasActiveFilters"
            icon="heroicons:server"
            :title="t('servers.empty.title')"
            :description="t('servers.empty.description')"
            :action-label="t('servers.createServer')"
            action-icon="i-heroicons-plus-circle"
            @action="router.push(`/tenant/${currentTenant?.id}/servers/create`)"
          />
          <EmptyState
            v-else
            icon="heroicons:magnifying-glass"
            :title="t('servers.noResults.title')"
            :description="t('servers.noResults.description')"
            :action-label="!hasActiveFilters ? t('servers.createServer') : undefined"
            action-icon="i-heroicons-plus-circle"
            @action="router.push(`/tenant/${currentTenant?.id}/servers/create`)"
          />
        </div>

        <div v-else class="space-y-4">
          <!-- Server Cards -->
          <div
            v-for="server in filteredServers"
            :key="server.id"
            class="border border-gray-200 dark:border-gray-700 rounded-lg p-6 hover:border-gray-300 dark:hover:border-gray-600 transition-colors cursor-pointer"
            @click="viewServer(server)"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <div class="flex-shrink-0">
                  <UIcon 
                    name="heroicons:server" 
                    class="w-8 h-8 text-gray-600 dark:text-gray-400"
                  />
                </div>
                <div>
                  <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100">
                    {{ server.name }}
                  </h3>
                  <p class="text-sm text-gray-500 dark:text-gray-400">
                    {{ server.game_type }} • {{ server.version || 'Latest' }}
                  </p>
                </div>
              </div>
              <div class="flex items-center space-x-4">
                <div class="text-right">
                  <StatusBadge :status="server.status.phase" />
                  <p class="text-xs text-gray-500 mt-1">
                    {{ server.status.player_count || 0 }}/{{ server.max_players || 20 }} players
                  </p>
                </div>
                <UDropdown :items="getServerActions(server)">
                  <UButton
                    color="gray"
                    variant="ghost"
                    icon="i-heroicons-ellipsis-vertical-20-solid"
                    @click.stop
                  />
                </UDropdown>
              </div>
            </div>

            <!-- Server Stats -->
            <div class="mt-4 grid grid-cols-3 gap-4 text-sm">
              <div>
                <p class="text-gray-500 dark:text-gray-400">{{ t('servers.stats.uptime') }}</p>
                <p class="font-medium text-gray-900 dark:text-gray-100">
                  {{ formatUptime(server.status.uptime) }}
                </p>
              </div>
              <div>
                <p class="text-gray-500 dark:text-gray-400">{{ t('servers.stats.memory') }}</p>
                <p class="font-medium text-gray-900 dark:text-gray-100">
                  {{ formatMemory(server.status.memory_usage) }}
                </p>
              </div>
              <div>
                <p class="text-gray-500 dark:text-gray-400">{{ t('servers.stats.cpu') }}</p>
                <p class="font-medium text-gray-900 dark:text-gray-100">
                  {{ formatCpu(server.status.cpu_usage) }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth', 'tenant']
})

const { t } = useI18n()
const router = useRouter()
const { currentTenant, tenantApiRequest } = useTenant()

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
const servers = ref<Server[]>([])
const searchQuery = ref('')
const activeFilters = ref<Record<string, any>>({})

// Computed
const filteredServers = computed(() => {
  let filtered = servers.value

  // Apply search
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(server => 
      server.name.toLowerCase().includes(query) ||
      server.game_type.toLowerCase().includes(query)
    )
  }

  // Apply filters
  if (activeFilters.value.status) {
    filtered = filtered.filter(server => 
      server.status.phase.toLowerCase() === activeFilters.value.status.toLowerCase()
    )
  }

  if (activeFilters.value.game_type) {
    filtered = filtered.filter(server => 
      server.game_type === activeFilters.value.game_type
    )
  }

  return filtered
})

const hasActiveFilters = computed(() => {
  return Object.keys(activeFilters.value).length > 0 || searchQuery.value.length > 0
})

const filterOptions = computed(() => [
  {
    key: 'status',
    label: t('servers.filters.status'),
    options: [
      { label: t('servers.status.running'), value: 'running' },
      { label: t('servers.status.stopped'), value: 'stopped' },
      { label: t('servers.status.pending'), value: 'pending' },
      { label: t('servers.status.failed'), value: 'failed' }
    ]
  },
  {
    key: 'game_type',
    label: t('servers.filters.gameType'),
    options: [
      { label: 'Minecraft', value: 'minecraft' },
      { label: 'CS2', value: 'cs2' },
      { label: 'Valheim', value: 'valheim' },
      { label: 'Terraria', value: 'terraria' }
    ]
  }
])

// Server actions dropdown
const getServerActions = (server: Server) => [
  [{
    label: t('servers.actions.viewDetails'),
    icon: 'i-heroicons-eye-20-solid',
    click: () => router.push(`/tenant/${currentTenant.value?.id}/servers/${server.id}`)
  }, {
    label: t('servers.actions.console'),
    icon: 'i-heroicons-command-line-20-solid',
    click: () => router.push(`/tenant/${currentTenant.value?.id}/servers/${server.id}?tab=console`)
  }],
  [{
    label: server.status.phase === 'Running' ? t('servers.actions.stop') : t('servers.actions.start'),
    icon: server.status.phase === 'Running' ? 'i-heroicons-stop-circle-20-solid' : 'i-heroicons-play-circle-20-solid',
    click: () => toggleServer(server)
  }, {
    label: t('servers.actions.restart'),
    icon: 'i-heroicons-arrow-path-20-solid',
    click: () => restartServer(server)
  }],
  [{
    label: t('servers.actions.edit'),
    icon: 'i-heroicons-pencil-square-20-solid',
    click: () => router.push(`/tenant/${currentTenant.value?.id}/servers/${server.id}/edit`)
  }, {
    label: 'Create Backup',
    icon: 'i-heroicons-archive-box-20-solid',
    click: () => createBackup(server)
  }],
  [{
    label: t('servers.actions.delete'),
    icon: 'i-heroicons-trash-20-solid',
    click: () => deleteServer(server)
  }]
]

// Methods
const loadServers = async () => {
  if (!currentTenant.value) return

  try {
    isLoading.value = true
    error.value = null

    const response = await tenantApiRequest<{ servers: Server[] }>('/api/tenant/servers')
    servers.value = response.servers || []
  } catch (err: any) {
    console.error('Failed to load servers:', err)
    error.value = err?.data?.message || 'Failed to load servers'
  } finally {
    isLoading.value = false
  }
}

// Navigation functions
const viewServer = (server: Server) => {
  router.push(`/tenant/${currentTenant.value?.id}/servers/${server.id}`)
}

const clearFilters = () => {
  activeFilters.value = {}
  searchQuery.value = ''
}

// Server actions
const toggleServer = async (server: Server) => {
  try {
    const action = server.status.phase === 'Running' ? 'stop' : 'start'
    await tenantApiRequest(`/api/tenant/servers/${server.id}/${action}`, {
      method: 'POST'
    })
    
    const toast = useToast()
    toast.add({
      title: `Server ${action === 'start' ? 'Started' : 'Stopped'}`,
      description: `${server.name} has been ${action === 'start' ? 'started' : 'stopped'}`,
      color: 'success'
    })
    
    await loadServers()
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Action Failed',
      description: err?.data?.message || 'Failed to perform server action',
      color: 'error'
    })
  }
}

const restartServer = async (server: Server) => {
  try {
    await tenantApiRequest(`/api/tenant/servers/${server.id}/restart`, {
      method: 'POST'
    })
    
    const toast = useToast()
    toast.add({
      title: 'Server Restarted',
      description: `${server.name} has been restarted`,
      color: 'success'
    })
    
    await loadServers()
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Restart Failed',
      description: err?.data?.message || 'Failed to restart server',
      color: 'error'
    })
  }
}

const createBackup = async (server: Server) => {
  try {
    await tenantApiRequest(`/api/tenant/servers/${server.id}/backup`, {
      method: 'POST'
    })
    
    const toast = useToast()
    toast.add({
      title: 'Backup Created',
      description: `Backup for ${server.name} has been created`,
      color: 'success'
    })
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Backup Failed',
      description: err?.data?.message || 'Failed to create backup',
      color: 'error'
    })
  }
}

const deleteServer = async (server: Server) => {
  const confirmed = confirm(`Are you sure you want to delete ${server.name}? This action cannot be undone.`)
  if (!confirmed) return

  try {
    await tenantApiRequest(`/api/tenant/servers/${server.id}`, {
      method: 'DELETE'
    })
    
    const toast = useToast()
    toast.add({
      title: 'Server Deleted',
      description: `${server.name} has been deleted`,
      color: 'success'
    })
    
    await loadServers()
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Delete Failed',
      description: err?.data?.message || 'Failed to delete server',
      color: 'error'
    })
  }
}

// Utility functions
const formatUptime = (uptime?: string) => {
  if (!uptime) return 'N/A'
  // Parse uptime and format nicely
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

// Load servers on mount
onMounted(() => {
  loadServers()
})

// Watch for tenant changes
watch(currentTenant, () => {
  if (currentTenant.value) {
    loadServers()
  }
})
</script>