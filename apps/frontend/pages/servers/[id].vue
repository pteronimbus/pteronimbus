<script setup lang="ts">
definePageMeta({
  layout: 'default'
})

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const serverId = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id

// Define types for better type safety
interface Player {
  id: number
  name: string
  joinedAt: string
  playtime: string
  ping: number
}

interface Backup {
  id: number
  name: string
  size: string
  createdAt: string
  type: string
}

// Mock server data - in real app this would come from API
const server = ref({
  id: parseInt(serverId),
  name: 'Minecraft Survival',
  game: 'Minecraft',
  status: 'online',
  ip: '192.168.1.100',
  port: 25565,
  version: '1.20.4',
  players: {
    current: 10,
    max: 20
  },
  uptime: '2d 14h 30m',
  createdAt: '2024-01-15',
  lastRestart: '2024-01-20 08:00',
  performance: {
    cpu: 45,
    memory: 60,
    disk: 25,
    network: {
      in: '1.2 MB/s',
      out: '0.8 MB/s'
    }
  },
  settings: {
    autoRestart: true,
    backupEnabled: true,
    backupInterval: '6h',
    difficulty: 'normal',
    gameMode: 'survival',
    maxPlayers: 20,
    pvp: true,
    onlineMode: true
  },
  currentPlayers: [
    { id: 1, name: 'PlayerOne', joinedAt: '2024-01-20 15:30', playtime: '2h 30m', ping: 45 },
    { id: 2, name: 'GamerTwo', joinedAt: '2024-01-20 14:15', playtime: '3h 45m', ping: 32 },
    { id: 3, name: 'BuilderThree', joinedAt: '2024-01-20 16:00', playtime: '1h 00m', ping: 28 }
  ],
  recentLogs: [
    { id: 1, timestamp: '2024-01-20 16:30:15', level: 'INFO', message: 'PlayerOne joined the game' },
    { id: 2, timestamp: '2024-01-20 16:29:45', level: 'INFO', message: 'Server tick took 52ms (should be max 50ms)' },
    { id: 3, timestamp: '2024-01-20 16:28:30', level: 'WARN', message: 'Memory usage at 85%' },
    { id: 4, timestamp: '2024-01-20 16:27:10', level: 'INFO', message: 'BuilderThree left the game' },
    { id: 5, timestamp: '2024-01-20 16:25:00', level: 'INFO', message: 'Auto-save completed' }
  ],
  backups: [
    { id: 1, name: 'auto-backup-2024-01-20-16-00', size: '1.2 GB', createdAt: '2024-01-20 16:00:00', type: 'automatic' },
    { id: 2, name: 'manual-backup-2024-01-20-12-00', size: '1.1 GB', createdAt: '2024-01-20 12:00:00', type: 'manual' },
    { id: 3, name: 'auto-backup-2024-01-20-10-00', size: '1.0 GB', createdAt: '2024-01-20 10:00:00', type: 'automatic' }
  ]
})

const activeTab = ref('overview')
const consoleInput = ref('')
const consoleOutput = ref([
  '[16:30:15] [Server thread/INFO]: PlayerOne joined the game',
  '[16:29:45] [Server thread/INFO]: Server tick took 52ms (should be max 50ms)',
  '[16:28:30] [Server thread/WARN]: Memory usage at 85%',
  '[16:27:10] [Server thread/INFO]: BuilderThree left the game',
  '[16:25:00] [Server thread/INFO]: Auto-save completed'
])

const tabs = [
  { key: 'overview', label: 'Overview', icon: 'i-heroicons-home-20-solid' },
  { key: 'console', label: t('servers.details.console'), icon: 'i-heroicons-command-line-20-solid' },
  { key: 'players', label: t('servers.details.players'), icon: 'i-heroicons-users-20-solid' },
  { key: 'performance', label: t('servers.details.performance'), icon: 'i-heroicons-chart-bar-20-solid' },
  { key: 'settings', label: t('servers.details.settings'), icon: 'i-heroicons-cog-6-tooth-20-solid' },
  { key: 'backups', label: t('servers.details.backups'), icon: 'i-heroicons-archive-box-20-solid' },
  { key: 'logs', label: t('servers.details.logs'), icon: 'i-heroicons-document-text-20-solid' }
]

// Stats configuration for StatsCard components
const serverStats = computed(() => [
  {
    key: 'players',
    label: t('servers.columns.players'),
    value: `${server.value.players.current}/${server.value.players.max}`,
    icon: 'i-heroicons-users-20-solid',
    color: 'blue'
  },
  {
    key: 'uptime',
    label: t('servers.columns.uptime'),
    value: server.value.uptime,
    icon: 'i-heroicons-clock-20-solid',
    color: 'green'
  },
  {
    key: 'cpu',
    label: 'CPU Usage',
    value: `${server.value.performance.cpu}%`,
    icon: 'i-heroicons-cpu-chip-20-solid',
    color: server.value.performance.cpu > 80 ? 'red' : server.value.performance.cpu > 60 ? 'yellow' : 'green'
  },
  {
    key: 'memory',
    label: 'Memory Usage',
    value: `${server.value.performance.memory}%`,
    icon: 'i-heroicons-circle-stack-20-solid',
    color: server.value.performance.memory > 80 ? 'red' : server.value.performance.memory > 60 ? 'yellow' : 'green'
  }
])

// Page header actions
const headerActions = computed(() => {
  const actions = []
  
  if (server.value.status === 'offline') {
    actions.push({
      label: t('servers.actions.start'),
      icon: 'i-heroicons-play-20-solid',
      color: 'success' as const,
      onClick: startServer
    })
  }
  
  if (server.value.status === 'online') {
    actions.push({
      label: t('servers.actions.stop'),
      icon: 'i-heroicons-stop-20-solid',
      color: 'error' as const,
      onClick: stopServer
    })
    actions.push({
      label: t('servers.actions.restart'),
      icon: 'i-heroicons-arrow-path-20-solid',
      color: 'warning' as const,
      onClick: restartServer
    })
  }
  
  return actions
})

// Helper functions
const getLogLevelColor = (level: string) => {
  switch (level) {
    case 'INFO': return 'info'
    case 'WARN': return 'warning'
    case 'ERROR': return 'error'
    case 'DEBUG': return 'neutral'
    default: return 'neutral'
  }
}

const startServer = () => {
  server.value.status = 'starting'
  // Simulate server start
  setTimeout(() => {
    server.value.status = 'online'
  }, 3000)
}

const stopServer = () => {
  server.value.status = 'stopping'
  // Simulate server stop
  setTimeout(() => {
    server.value.status = 'offline'
  }, 2000)
}

const restartServer = () => {
  server.value.status = 'stopping'
  // Simulate server restart
  setTimeout(() => {
    server.value.status = 'starting'
    setTimeout(() => {
      server.value.status = 'online'
    }, 3000)
  }, 2000)
}

const sendConsoleCommand = () => {
  if (consoleInput.value.trim()) {
    consoleOutput.value.push(`> ${consoleInput.value}`)
    // Simulate command response
    setTimeout(() => {
      consoleOutput.value.push(`[Server thread/INFO]: Command executed: ${consoleInput.value}`)
    }, 500)
    consoleInput.value = ''
  }
}

const kickPlayer = (player: Player) => {
  // Remove player from current players list
  const index = server.value.currentPlayers.findIndex(p => p.id === player.id)
  if (index !== -1) {
    server.value.currentPlayers.splice(index, 1)
    server.value.players.current--
  }
}

const banPlayer = (player: Player) => {
  // Ban player and remove from current players list
  kickPlayer(player)
  console.log('Banned player:', player.name)
}

const createBackup = () => {
  const newBackup = {
    id: Date.now(),
    name: `manual-backup-${new Date().toISOString().replace(/[:.]/g, '-')}`,
    size: '1.2 GB',
    createdAt: new Date().toISOString(),
    type: 'manual'
  }
  server.value.backups.unshift(newBackup)
}

const deleteBackup = (backup: Backup) => {
  const index = server.value.backups.findIndex(b => b.id === backup.id)
  if (index !== -1) {
    server.value.backups.splice(index, 1)
  }
}

const restoreBackup = (backup: Backup) => {
  console.log('Restore backup:', backup.name)
}

const downloadBackup = (backup: Backup) => {
  console.log('Download backup:', backup.name)
}

const goBack = () => {
  router.push('/servers')
}

// Auto-scroll console output
const consoleContainer = ref<HTMLElement | null>(null)
watch(consoleOutput, () => {
  nextTick(() => {
    if (consoleContainer.value) {
      consoleContainer.value.scrollTop = consoleContainer.value.scrollHeight
    }
  })
})
</script>

<template>
  <div>
    <!-- Page Header -->
    <PageHeader 
      :title="t('servers.details.title')"
      :description="`${server.name} - ${server.game} ${server.version}`"
      :actions="headerActions"
      show-back-button
      @back="goBack"
    />

    <!-- Server Info Card -->
    <UCard class="mb-6">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-6">
          <div class="flex-shrink-0 w-16 h-16 bg-blue-100 dark:bg-blue-900 rounded-lg flex items-center justify-center">
            <UIcon name="i-heroicons-server-20-solid" class="w-8 h-8 text-blue-600 dark:text-blue-400" />
          </div>
          <div>
            <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ server.name }}</h2>
            <p class="text-gray-600 dark:text-gray-400">{{ server.game }} {{ server.version }}</p>
            <div class="flex items-center gap-3 mt-2">
              <StatusBadge :status="server.status" type="server" />
              <span class="text-sm text-gray-500 dark:text-gray-400">{{ server.ip }}:{{ server.port }}</span>
            </div>
          </div>
        </div>
        
        <!-- Additional action dropdown -->
        <div class="flex items-center gap-2">
          <UDropdownMenu :items="[
            [{
              label: t('servers.actions.edit'),
              icon: 'i-heroicons-pencil-square-20-solid',
              click: () => router.push(`/servers/${serverId}/edit`)
            }],
            [{
              label: 'Create Backup',
              icon: 'i-heroicons-archive-box-20-solid',
              click: createBackup
            }],
            [{
              label: t('servers.actions.delete'),
              icon: 'i-heroicons-trash-20-solid',
              click: () => console.log('Delete server')
            }]
          ]">
            <UButton 
              color="neutral" 
              variant="ghost" 
              icon="i-heroicons-ellipsis-horizontal-20-solid"
              class="text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200"
            />
          </UDropdownMenu>
        </div>
      </div>
    </UCard>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <StatsCard
        v-for="stat in serverStats"
        :key="stat.key"
        :label="stat.label"
        :value="stat.value"
        :icon="stat.icon"
        :color="stat.color"
      />
    </div>

    <!-- Tabs -->
    <div class="mb-6">
      <div class="border-b border-gray-200 dark:border-gray-700">
        <nav class="flex space-x-8 overflow-x-auto">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            @click="activeTab = tab.key"
            :class="[
              'flex items-center gap-2 py-4 px-1 border-b-2 font-medium text-sm transition-colors whitespace-nowrap',
              activeTab === tab.key
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
            ]"
          >
            <UIcon :name="tab.icon" class="w-4 h-4" />
            {{ tab.label }}
          </button>
        </nav>
      </div>
    </div>

    <!-- Tab Content -->
    <div>
      <!-- Overview Tab -->
      <div v-if="activeTab === 'overview'">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Server Information -->
          <UCard>
            <template #header>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Server Information</h3>
            </template>
            <div class="space-y-4">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Game:</span>
                <span class="font-medium text-gray-900 dark:text-gray-100">{{ server.game }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Version:</span>
                <span class="font-medium text-gray-900 dark:text-gray-100">{{ server.version }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Status:</span>
                <StatusBadge :status="server.status" type="server" />
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Address:</span>
                <span class="font-medium text-gray-900 dark:text-gray-100">{{ server.ip }}:{{ server.port }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Created:</span>
                <span class="font-medium text-gray-900 dark:text-gray-100">{{ server.createdAt }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Last Restart:</span>
                <span class="font-medium text-gray-900 dark:text-gray-100">{{ server.lastRestart }}</span>
              </div>
            </div>
          </UCard>

          <!-- Performance Overview -->
          <UCard>
            <template #header>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Performance</h3>
            </template>
            <div class="space-y-4">
              <div>
                <div class="flex justify-between mb-2">
                  <span class="text-sm text-gray-600 dark:text-gray-400">CPU Usage</span>
                  <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ server.performance.cpu }}%</span>
                </div>
                <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div 
                    class="bg-blue-500 h-2 rounded-full transition-all duration-300"
                    :style="{ width: `${server.performance.cpu}%` }"
                  />
                </div>
              </div>
              <div>
                <div class="flex justify-between mb-2">
                  <span class="text-sm text-gray-600 dark:text-gray-400">Memory Usage</span>
                  <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ server.performance.memory }}%</span>
                </div>
                <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div 
                    class="bg-green-500 h-2 rounded-full transition-all duration-300"
                    :style="{ width: `${server.performance.memory}%` }"
                  />
                </div>
              </div>
              <div>
                <div class="flex justify-between mb-2">
                  <span class="text-sm text-gray-600 dark:text-gray-400">Disk Usage</span>
                  <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ server.performance.disk }}%</span>
                </div>
                <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div 
                    class="bg-yellow-500 h-2 rounded-full transition-all duration-300"
                    :style="{ width: `${server.performance.disk}%` }"
                  />
                </div>
              </div>
              <div class="flex justify-between">
                <span class="text-sm text-gray-600 dark:text-gray-400">Network In:</span>
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ server.performance.network.in }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-sm text-gray-600 dark:text-gray-400">Network Out:</span>
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ server.performance.network.out }}</span>
              </div>
            </div>
          </UCard>
        </div>
      </div>

      <!-- Console Tab -->
      <div v-if="activeTab === 'console'">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('servers.details.console') }}</h3>
          </template>
          <div class="space-y-4">
            <!-- Console Output -->
            <div 
              ref="consoleContainer"
              class="h-96 bg-gray-900 text-green-400 p-4 rounded-lg overflow-y-auto font-mono text-sm"
            >
              <div v-for="(line, index) in consoleOutput" :key="index" class="mb-1">
                {{ line }}
              </div>
            </div>
            
            <!-- Console Input -->
            <div class="flex gap-2">
              <UInput
                v-model="consoleInput"
                placeholder="Enter command..."
                class="flex-1"
                @keyup.enter="sendConsoleCommand"
              />
              <UButton 
                color="success" 
                icon="i-heroicons-paper-airplane-20-solid"
                class="text-green-700 dark:text-green-300"
                @click="sendConsoleCommand"
              >
                Send
              </UButton>
            </div>
          </div>
        </UCard>
      </div>

      <!-- Players Tab -->
      <div v-if="activeTab === 'players'">
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('servers.details.players') }}</h3>
              <UBadge 
                color="primary" 
                variant="subtle"
                class="text-blue-700 dark:text-blue-300"
              >
                {{ server.players.current }}/{{ server.players.max }}
              </UBadge>
            </div>
          </template>
          <div class="space-y-4">
            <div 
              v-for="player in server.currentPlayers" 
              :key="player.id"
              class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800 rounded-lg"
            >
              <div class="flex items-center gap-4">
                <UAvatar :alt="player.name" size="sm">
                  <span class="text-xs font-medium text-primary-600 dark:text-primary-400">{{ player.name.charAt(0) }}</span>
                </UAvatar>
                <div>
                  <p class="font-medium text-gray-900 dark:text-gray-100">{{ player.name }}</p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">
                    Joined {{ player.joinedAt }} • Playtime: {{ player.playtime }} • Ping: {{ player.ping }}ms
                  </p>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <UButton 
                  color="warning" 
                  variant="ghost" 
                  icon="i-heroicons-user-minus-20-solid"
                  size="sm"
                  class="text-yellow-600 dark:text-yellow-400 hover:text-yellow-800 dark:hover:text-yellow-200"
                  @click="kickPlayer(player)"
                />
                <UButton 
                  color="error" 
                  variant="ghost" 
                  icon="i-heroicons-no-symbol-20-solid"
                  size="sm"
                  class="text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-200"
                  @click="banPlayer(player)"
                />
              </div>
            </div>
            <div v-if="server.currentPlayers.length === 0" class="text-center py-8">
              <UIcon name="i-heroicons-users-20-solid" class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-2" />
              <p class="text-gray-500 dark:text-gray-400">No players currently online</p>
            </div>
          </div>
        </UCard>
      </div>

      <!-- Performance Tab -->
      <div v-if="activeTab === 'performance'">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('servers.details.performance') }}</h3>
          </template>
          <div class="h-64 bg-gray-100 dark:bg-gray-800 rounded-md flex items-center justify-center">
            <div class="text-center">
              <UIcon name="i-heroicons-chart-bar-20-solid" class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-2" />
              <p class="text-gray-500 dark:text-gray-400">Performance charts</p>
              <p class="text-sm text-gray-400 dark:text-gray-500">Real-time performance monitoring would be displayed here</p>
            </div>
          </div>
        </UCard>
      </div>

      <!-- Settings Tab -->
      <div v-if="activeTab === 'settings'">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('servers.details.settings') }}</h3>
          </template>
          <div class="space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Max Players</label>
                <UInput v-model="server.settings.maxPlayers" type="number" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Difficulty</label>
                <USelect 
                  v-model="server.settings.difficulty"
                  :options="[
                    { value: 'peaceful', label: 'Peaceful' },
                    { value: 'easy', label: 'Easy' },
                    { value: 'normal', label: 'Normal' },
                    { value: 'hard', label: 'Hard' }
                  ]"
                />
              </div>
            </div>
            
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <div>
                  <label class="text-sm font-medium text-gray-700 dark:text-gray-300">Auto Restart</label>
                  <p class="text-xs text-gray-500 dark:text-gray-400">Automatically restart server on crash</p>
                </div>
                <UToggle v-model="server.settings.autoRestart" />
              </div>
              
              <div class="flex items-center justify-between">
                <div>
                  <label class="text-sm font-medium text-gray-700 dark:text-gray-300">Backup Enabled</label>
                  <p class="text-xs text-gray-500 dark:text-gray-400">Enable automatic backups</p>
                </div>
                <UToggle v-model="server.settings.backupEnabled" />
              </div>
              
              <div class="flex items-center justify-between">
                <div>
                  <label class="text-sm font-medium text-gray-700 dark:text-gray-300">PvP</label>
                  <p class="text-xs text-gray-500 dark:text-gray-400">Enable player vs player combat</p>
                </div>
                <UToggle v-model="server.settings.pvp" />
              </div>
              
              <div class="flex items-center justify-between">
                <div>
                  <label class="text-sm font-medium text-gray-700 dark:text-gray-300">Online Mode</label>
                  <p class="text-xs text-gray-500 dark:text-gray-400">Verify player authentication</p>
                </div>
                <UToggle v-model="server.settings.onlineMode" />
              </div>
            </div>
            
            <div class="flex justify-end">
              <UButton 
                color="primary"
                class="text-blue-700 dark:text-blue-300"
              >
                Save Settings
              </UButton>
            </div>
          </div>
        </UCard>
      </div>

      <!-- Backups Tab -->
      <div v-if="activeTab === 'backups'">
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('servers.details.backups') }}</h3>
              <UButton 
                color="primary" 
                icon="i-heroicons-plus-circle-20-solid" 
                class="text-blue-700 dark:text-blue-300"
                @click="createBackup"
              >
                Create Backup
              </UButton>
            </div>
          </template>
          <div class="space-y-4">
            <div 
              v-for="backup in server.backups" 
              :key="backup.id"
              class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-800 rounded-lg"
            >
              <div class="flex items-center gap-4">
                <div class="flex-shrink-0 w-10 h-10 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center">
                  <UIcon name="i-heroicons-archive-box-20-solid" class="w-5 h-5 text-blue-600 dark:text-blue-400" />
                </div>
                <div>
                  <p class="font-medium text-gray-900 dark:text-gray-100">{{ backup.name }}</p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">
                    <span class="text-gray-700 dark:text-gray-300">{{ backup.size }}</span> • 
                    <span class="text-gray-700 dark:text-gray-300">{{ backup.createdAt }}</span> • 
                    <span class="text-gray-700 dark:text-gray-300 capitalize">{{ backup.type }}</span>
                  </p>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <UButton 
                  color="primary" 
                  variant="ghost" 
                  icon="i-heroicons-arrow-down-tray-20-solid"
                  size="sm"
                  class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-200"
                  @click="downloadBackup(backup)"
                />
                <UButton 
                  color="success" 
                  variant="ghost" 
                  icon="i-heroicons-arrow-uturn-left-20-solid"
                  size="sm"
                  class="text-green-600 dark:text-green-400 hover:text-green-800 dark:hover:text-green-200"
                  @click="restoreBackup(backup)"
                />
                <UButton 
                  color="error" 
                  variant="ghost" 
                  icon="i-heroicons-trash-20-solid"
                  size="sm"
                  class="text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-200"
                  @click="deleteBackup(backup)"
                />
              </div>
            </div>
          </div>
        </UCard>
      </div>

      <!-- Logs Tab -->
      <div v-if="activeTab === 'logs'">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('servers.details.logs') }}</h3>
          </template>
          <div class="space-y-2">
            <div 
              v-for="log in server.recentLogs" 
              :key="log.id"
              class="flex items-start gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
            >
              <UBadge 
                :color="getLogLevelColor(log.level)" 
                variant="subtle"
                size="xs"
                :class="[
                  log.level === 'INFO' ? 'text-blue-700 dark:text-blue-300' : '',
                  log.level === 'WARN' ? 'text-yellow-700 dark:text-yellow-300' : '',
                  log.level === 'ERROR' ? 'text-red-700 dark:text-red-300' : '',
                  log.level === 'DEBUG' ? 'text-gray-700 dark:text-gray-300' : ''
                ]"
              >
                {{ log.level }}
              </UBadge>
              <div class="flex-1">
                <p class="text-sm text-gray-900 dark:text-gray-100 font-mono">{{ log.message }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ log.timestamp }}</p>
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template> 