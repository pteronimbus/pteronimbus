<script setup lang="ts">
definePageMeta({
  layout: 'default'
})

const { t } = useI18n()
const router = useRouter()

// Mock server data - in real app this would come from API
const servers = ref([
  { 
    id: 1, 
    name: 'Minecraft Survival', 
    game: 'Minecraft', 
    status: 'online', 
    players: '10/20', 
    ip: '192.168.1.100',
    port: 25565,
    version: '1.20.4',
    uptime: '2d 14h 30m',
    cpu: 45,
    memory: 60,
    createdAt: '2024-01-15'
  },
  { 
    id: 2, 
    name: 'Valheim Dedicated', 
    game: 'Valheim', 
    status: 'online', 
    players: '5/10', 
    ip: '192.168.1.101',
    port: 2456,
    version: '0.217.46',
    uptime: '5d 8h 15m',
    cpu: 32,
    memory: 45,
    createdAt: '2024-01-12'
  },
  { 
    id: 3, 
    name: 'CS:GO Competitive', 
    game: 'CS:GO', 
    status: 'offline', 
    players: '0/10', 
    ip: '192.168.1.102',
    port: 27015,
    version: '1.38.0.1',
    uptime: '0m',
    cpu: 0,
    memory: 0,
    createdAt: '2024-01-10'
  },
  { 
    id: 4, 
    name: 'Terraria Adventure', 
    game: 'Terraria', 
    status: 'starting', 
    players: '0/16', 
    ip: '192.168.1.103',
    port: 7777,
    version: '1.4.4.9',
    uptime: '0m',
    cpu: 15,
    memory: 20,
    createdAt: '2024-01-20'
  },
  { 
    id: 5, 
    name: 'Rust Survival', 
    game: 'Rust', 
    status: 'error', 
    players: '0/100', 
    ip: '192.168.1.104',
    port: 28015,
    version: '2023.12.7',
    uptime: '0m',
    cpu: 0,
    memory: 0,
    createdAt: '2024-01-18'
  }
])

const searchQuery = ref('')
const selectedStatus = ref('all')
const selectedGame = ref('all')

const statusOptions = [
  { value: 'all', label: 'All Status' },
  { value: 'online', label: t('servers.status.online') },
  { value: 'offline', label: t('servers.status.offline') },
  { value: 'starting', label: t('servers.status.starting') },
  { value: 'stopping', label: t('servers.status.stopping') },
  { value: 'error', label: t('servers.status.error') }
]

const gameOptions = computed(() => {
  const games = ['all', ...new Set(servers.value.map(s => s.game))]
  return games.map(game => ({
    value: game,
    label: game === 'all' ? 'All Games' : game
  }))
})

const columns = [
  { key: 'server', label: t('servers.columns.name'), id: 'server' },
  { key: 'game', label: t('servers.columns.game'), id: 'game' },
  { key: 'status', label: t('servers.columns.status'), id: 'status' },
  { key: 'players', label: t('servers.columns.players'), id: 'players' },
  { key: 'performance', label: 'Performance', id: 'performance' },
  { key: 'uptime', label: t('servers.columns.uptime'), id: 'uptime' },
  { key: 'actions', label: t('servers.columns.actions'), id: 'actions' }
]

// Filtered servers based on search and filters
const filteredServers = computed(() => {
  return servers.value.filter(server => {
    const matchesSearch = searchQuery.value === '' || 
      server.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      server.game.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      server.ip.includes(searchQuery.value)
    
    const matchesStatus = selectedStatus.value === 'all' || server.status === selectedStatus.value
    const matchesGame = selectedGame.value === 'all' || server.game === selectedGame.value
    
    return matchesSearch && matchesStatus && matchesGame
  })
})

// Action items for dropdown
const getActionItems = (server) => [
  [{
    label: t('servers.actions.viewDetails'),
    icon: 'i-heroicons-eye-20-solid',
    click: () => router.push(`/servers/${server.id}`)
  }, {
    label: t('servers.actions.console'),
    icon: 'i-heroicons-command-line-20-solid',
    click: () => router.push(`/servers/${server.id}?tab=console`)
  }],
  [{
    label: server.status === 'online' ? t('servers.actions.stop') : t('servers.actions.start'),
    icon: server.status === 'online' ? 'i-heroicons-stop-20-solid' : 'i-heroicons-play-20-solid',
    click: () => toggleServer(server)
  }, {
    label: t('servers.actions.restart'),
    icon: 'i-heroicons-arrow-path-20-solid',
    click: () => restartServer(server),
    disabled: server.status !== 'online'
  }],
  [{
    label: t('servers.actions.edit'),
    icon: 'i-heroicons-pencil-square-20-solid',
    click: () => router.push(`/servers/${server.id}/edit`)
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

// Helper functions
const getStatusColor = (status: string) => {
  switch (status) {
    case 'online': return 'success'
    case 'offline': return 'error'
    case 'starting': return 'warning'
    case 'stopping': return 'warning'
    case 'error': return 'error'
    default: return 'neutral'
  }
}

const getPerformanceColor = (cpu: number, memory: number) => {
  const maxUsage = Math.max(cpu, memory)
  if (maxUsage > 80) return 'error'
  if (maxUsage > 60) return 'warning'
  return 'success'
}

// Action functions
const toggleServer = (server: any) => {
  const index = servers.value.findIndex(s => s.id === server.id)
  if (index !== -1) {
    if (server.status === 'online') {
      servers.value[index].status = 'stopping'
      setTimeout(() => {
        servers.value[index].status = 'offline'
        servers.value[index].cpu = 0
        servers.value[index].memory = 0
        servers.value[index].uptime = '0m'
      }, 2000)
    } else {
      servers.value[index].status = 'starting'
      setTimeout(() => {
        servers.value[index].status = 'online'
        servers.value[index].cpu = Math.floor(Math.random() * 60) + 20
        servers.value[index].memory = Math.floor(Math.random() * 50) + 30
      }, 3000)
    }
  }
}

const restartServer = (server: any) => {
  const index = servers.value.findIndex(s => s.id === server.id)
  if (index !== -1) {
    servers.value[index].status = 'stopping'
    setTimeout(() => {
      servers.value[index].status = 'starting'
      setTimeout(() => {
        servers.value[index].status = 'online'
      }, 3000)
    }, 2000)
  }
}

const createBackup = (server: any) => {
  console.log('Creating backup for:', server.name)
  // Implementation for creating backup
}

const deleteServer = (server: any) => {
  console.log('Deleting server:', server.name)
  // Implementation for server deletion
}

// Navigation functions
const viewServer = (server: any) => {
  router.push(`/servers/${server.id}`)
}

const createServer = () => {
  router.push('/servers/create')
}

// Stats
const serverStats = computed(() => ({
  total: servers.value.length,
  online: servers.value.filter(s => s.status === 'online').length,
  offline: servers.value.filter(s => s.status === 'offline').length,
  error: servers.value.filter(s => s.status === 'error').length
}))
</script>

<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
      <div>
        <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">{{ t('servers.title') }}</h1>
        <p class="mt-1 text-gray-500 dark:text-gray-400">
          Manage and monitor your game servers
        </p>
      </div>
      <div class="mt-4 sm:mt-0">
        <UButton 
          icon="i-heroicons-plus-circle" 
          size="lg"
          @click="createServer"
        >
          {{ t('servers.createServer') }}
        </UButton>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <UCard class="cursor-pointer hover:shadow-lg transition-shadow" @click="selectedStatus = 'all'">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Total Servers</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ serverStats.total }}</p>
          </div>
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-full">
            <UIcon name="i-heroicons-server-20-solid" class="w-6 h-6 text-blue-600 dark:text-blue-400" />
          </div>
        </div>
      </UCard>
      
      <UCard class="cursor-pointer hover:shadow-lg transition-shadow" @click="selectedStatus = 'online'">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Online</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ serverStats.online }}</p>
          </div>
          <div class="p-3 bg-green-100 dark:bg-green-900 rounded-full">
            <UIcon name="i-heroicons-check-circle-20-solid" class="w-6 h-6 text-green-600 dark:text-green-400" />
          </div>
        </div>
      </UCard>
      
      <UCard class="cursor-pointer hover:shadow-lg transition-shadow" @click="selectedStatus = 'offline'">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Offline</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ serverStats.offline }}</p>
          </div>
          <div class="p-3 bg-gray-100 dark:bg-gray-900 rounded-full">
            <UIcon name="i-heroicons-x-circle-20-solid" class="w-6 h-6 text-gray-600 dark:text-gray-400" />
          </div>
        </div>
      </UCard>
      
      <UCard class="cursor-pointer hover:shadow-lg transition-shadow" @click="selectedStatus = 'error'">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Errors</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ serverStats.error }}</p>
          </div>
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-full">
            <UIcon name="i-heroicons-exclamation-triangle-20-solid" class="w-6 h-6 text-red-600 dark:text-red-400" />
          </div>
        </div>
      </UCard>
    </div>

    <!-- Filters -->
    <div class="mb-6 flex flex-col sm:flex-row gap-4">
      <div class="flex-1">
        <UInput
          v-model="searchQuery"
          :placeholder="t('common.search') + ' servers...'"
          icon="i-heroicons-magnifying-glass-20-solid"
          size="md"
        />
      </div>
      <div class="flex gap-2">
        <USelect
          v-model="selectedStatus"
          :options="statusOptions"
          size="md"
          class="w-40"
        />
        <USelect
          v-model="selectedGame"
          :options="gameOptions"
          size="md"
          class="w-40"
        />
      </div>
    </div>

    <!-- Servers Table -->
    <UCard>
      <UTable :rows="filteredServers" :columns="columns">
        <!-- Server column with icon and name -->
        <template #server-data="{ row }">
          <div class="flex items-center gap-3">
            <div class="flex-shrink-0 w-10 h-10 bg-blue-100 dark:bg-blue-900 rounded-lg flex items-center justify-center">
              <UIcon name="i-heroicons-server-20-solid" class="w-5 h-5 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <button
                @click="viewServer(row)"
                class="font-medium text-gray-900 dark:text-gray-100 hover:text-primary-600 dark:hover:text-primary-400 transition-colors"
              >
                {{ row.name }}
              </button>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ row.ip }}:{{ row.port }}</p>
            </div>
          </div>
        </template>

        <!-- Game column -->
        <template #game-data="{ row }">
          <div>
            <span class="font-medium text-gray-900 dark:text-gray-100">{{ row.game }}</span>
            <p class="text-sm text-gray-500 dark:text-gray-400">{{ row.version }}</p>
          </div>
        </template>

        <!-- Status column -->
        <template #status-data="{ row }">
          <UBadge 
            :color="getStatusColor(row.status)" 
            variant="subtle"
            class="capitalize"
            :class="[
              row.status === 'online' ? 'text-green-700 dark:text-green-300' : '',
              row.status === 'offline' ? 'text-red-700 dark:text-red-300' : '',
              row.status === 'starting' ? 'text-yellow-700 dark:text-yellow-300' : '',
              row.status === 'stopping' ? 'text-orange-700 dark:text-orange-300' : '',
              row.status === 'error' ? 'text-red-700 dark:text-red-300' : ''
            ]"
          >
            {{ t(`servers.status.${row.status}`) }}
          </UBadge>
        </template>

        <!-- Players column -->
        <template #players-data="{ row }">
          <div class="flex items-center gap-2">
            <UIcon name="i-heroicons-users-20-solid" class="w-4 h-4 text-gray-400" />
            <span class="font-medium text-gray-900 dark:text-gray-100">{{ row.players }}</span>
          </div>
        </template>

        <!-- Performance column -->
        <template #performance-data="{ row }">
          <div class="flex items-center gap-2">
            <div class="flex flex-col gap-1">
              <div class="flex items-center gap-2">
                <span class="text-xs text-gray-500 dark:text-gray-400">CPU:</span>
                <div class="w-12 h-2 bg-gray-200 dark:bg-gray-700 rounded-full">
                  <div 
                    :class="[
                      'h-2 rounded-full transition-all duration-300',
                      getPerformanceColor(row.cpu, row.memory) === 'success' ? 'bg-green-500' : '',
                      getPerformanceColor(row.cpu, row.memory) === 'warning' ? 'bg-yellow-500' : '',
                      getPerformanceColor(row.cpu, row.memory) === 'error' ? 'bg-red-500' : ''
                    ]"
                    :style="{ width: `${row.cpu}%` }"
                  />
                </div>
                <span class="text-xs font-medium text-gray-900 dark:text-gray-100">{{ row.cpu }}%</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-xs text-gray-500 dark:text-gray-400">RAM:</span>
                <div class="w-12 h-2 bg-gray-200 dark:bg-gray-700 rounded-full">
                  <div 
                    :class="[
                      'h-2 rounded-full transition-all duration-300',
                      getPerformanceColor(row.cpu, row.memory) === 'success' ? 'bg-green-500' : '',
                      getPerformanceColor(row.cpu, row.memory) === 'warning' ? 'bg-yellow-500' : '',
                      getPerformanceColor(row.cpu, row.memory) === 'error' ? 'bg-red-500' : ''
                    ]"
                    :style="{ width: `${row.memory}%` }"
                  />
                </div>
                <span class="text-xs font-medium text-gray-900 dark:text-gray-100">{{ row.memory }}%</span>
              </div>
            </div>
          </div>
        </template>

        <!-- Uptime column -->
        <template #uptime-data="{ row }">
          <span class="text-sm text-gray-600 dark:text-gray-400">{{ row.uptime }}</span>
        </template>

        <!-- Actions column -->
        <template #actions-data="{ row }">
          <div class="flex items-center gap-2">
            <UButton 
              color="primary" 
              variant="ghost" 
              icon="i-heroicons-eye-20-solid"
              size="sm"
              class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-200"
              @click="viewServer(row)"
            />
            <UButton 
              v-if="row.status === 'online'"
              color="error" 
              variant="ghost" 
              icon="i-heroicons-stop-20-solid"
              size="sm"
              class="text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-200"
              @click="toggleServer(row)"
            />
            <UButton 
              v-else-if="row.status === 'offline'"
              color="success" 
              variant="ghost" 
              icon="i-heroicons-play-20-solid"
              size="sm"
              class="text-green-600 dark:text-green-400 hover:text-green-800 dark:hover:text-green-200"
              @click="toggleServer(row)"
            />
            <UDropdown :items="getActionItems(row)">
              <UButton 
                color="neutral" 
                variant="ghost" 
                icon="i-heroicons-ellipsis-horizontal-20-solid"
                size="sm"
                class="text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200"
              />
            </UDropdown>
          </div>
        </template>
      </UTable>

      <!-- Empty state -->
      <div v-if="filteredServers.length === 0" class="text-center py-12">
        <UIcon name="i-heroicons-server-20-solid" class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" />
        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
          {{ searchQuery || selectedStatus !== 'all' || selectedGame !== 'all' ? 'No servers found' : t('servers.noServers') }}
        </h3>
        <p class="text-gray-500 dark:text-gray-400 mb-6">
          {{ searchQuery || selectedStatus !== 'all' || selectedGame !== 'all' 
            ? 'Try adjusting your search or filters' 
            : 'Get started by creating your first server' }}
        </p>
        <UButton 
          v-if="!searchQuery && selectedStatus === 'all' && selectedGame === 'all'"
          @click="createServer"
          icon="i-heroicons-plus-circle"
          class="text-blue-700 dark:text-blue-300"
        >
          {{ t('servers.createServer') }}
        </UButton>
      </div>
    </UCard>
  </div>
</template> 