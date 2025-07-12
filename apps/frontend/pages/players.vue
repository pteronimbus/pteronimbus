<script setup lang="ts">
definePageMeta({
  layout: 'default'
})

const { t } = useI18n()
const router = useRouter()

// Mock players data
const players = ref([
  {
    id: 1,
    name: 'PlayerOne',
    server: 'Minecraft Survival',
    status: 'online',
    playtime: '4h 23m',
    lastSeen: '2 minutes ago',
    avatar: null
  },
  {
    id: 2,
    name: 'GamerGirl',
    server: 'Valheim Dedicated',
    status: 'online',
    playtime: '12h 45m',
    lastSeen: '5 minutes ago',
    avatar: null
  },
  {
    id: 3,
    name: 'ProPlayer',
    server: 'CS:GO Competitive',
    status: 'offline',
    playtime: '156h 12m',
    lastSeen: '1 hour ago',
    avatar: null
  },
  {
    id: 4,
    name: 'CasualGamer',
    server: 'Terraria Adventure',
    status: 'online',
    playtime: '8h 34m',
    lastSeen: 'Just now',
    avatar: null
  }
])

const searchQuery = ref('')
const selectedStatus = ref('all')

const statusOptions = [
  { label: 'All Status', value: 'all' },
  { label: 'Online', value: 'online' },
  { label: 'Offline', value: 'offline' }
]

const columns = [
  { key: 'player', label: 'Player', id: 'player' },
  { key: 'server', label: 'Server', id: 'server' },
  { key: 'status', label: 'Status', id: 'status' },
  { key: 'playtime', label: 'Playtime', id: 'playtime' },
  { key: 'lastSeen', label: 'Last Seen', id: 'lastSeen' }
]

const filteredPlayers = computed(() => {
  let filtered = players.value

  if (searchQuery.value) {
    filtered = filtered.filter(player =>
      player.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      player.server.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  if (selectedStatus.value !== 'all') {
    filtered = filtered.filter(player => player.status === selectedStatus.value)
  }

  return filtered
})

const getStatusColor = (status: string) => {
  switch (status) {
    case 'online': return 'success'
    case 'offline': return 'neutral'
    default: return 'neutral'
  }
}

const playerStats = computed(() => ({
  total: players.value.length,
  online: players.value.filter(p => p.status === 'online').length,
  offline: players.value.filter(p => p.status === 'offline').length
}))
</script>

<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
      <div>
        <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">Players</h1>
        <p class="mt-1 text-gray-500 dark:text-gray-400">
          Monitor player activity across all servers
        </p>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Total Players</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ playerStats.total }}</p>
          </div>
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-full">
            <UIcon name="i-heroicons-users-20-solid" class="w-6 h-6 text-blue-600 dark:text-blue-400" />
          </div>
        </div>
      </UCard>
      
      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Online Now</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ playerStats.online }}</p>
          </div>
          <div class="p-3 bg-green-100 dark:bg-green-900 rounded-full">
            <UIcon name="i-heroicons-check-circle-20-solid" class="w-6 h-6 text-green-600 dark:text-green-400" />
          </div>
        </div>
      </UCard>
      
      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Offline</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ playerStats.offline }}</p>
          </div>
          <div class="p-3 bg-gray-100 dark:bg-gray-900 rounded-full">
            <UIcon name="i-heroicons-x-circle-20-solid" class="w-6 h-6 text-gray-600 dark:text-gray-400" />
          </div>
        </div>
      </UCard>
    </div>

    <!-- Filters -->
    <div class="mb-6 flex flex-col sm:flex-row gap-4">
      <div class="flex-1">
        <UInput
          v-model="searchQuery"
          placeholder="Search players..."
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
      </div>
    </div>

    <!-- Players Table -->
    <UCard>
      <UTable :rows="filteredPlayers" :columns="columns">
        <!-- Player column with avatar and name -->
        <template #player-data="{ row }">
          <div class="flex items-center gap-3">
            <UAvatar
              :src="row.avatar"
              :alt="row.name"
              size="sm"
              :ui="{ background: 'bg-primary-100 dark:bg-primary-900' }"
            >
              <span class="text-xs font-medium text-primary-600 dark:text-primary-400">
                {{ row.name.split('').slice(0, 2).join('').toUpperCase() }}
              </span>
            </UAvatar>
            <div>
              <span class="font-medium text-gray-900 dark:text-gray-100">{{ row.name }}</span>
            </div>
          </div>
        </template>

        <!-- Server column -->
        <template #server-data="{ row }">
          <span class="text-gray-600 dark:text-gray-400">{{ row.server }}</span>
        </template>

        <!-- Status column -->
        <template #status-data="{ row }">
          <UBadge 
            :color="getStatusColor(row.status)" 
            variant="subtle"
            class="capitalize"
            :class="[
              row.status === 'online' ? 'text-green-700 dark:text-green-300' : '',
              row.status === 'offline' ? 'text-gray-700 dark:text-gray-300' : ''
            ]"
          >
            {{ row.status }}
          </UBadge>
        </template>

        <!-- Playtime column -->
        <template #playtime-data="{ row }">
          <span class="text-gray-900 dark:text-gray-100">{{ row.playtime }}</span>
        </template>

        <!-- Last Seen column -->
        <template #lastSeen-data="{ row }">
          <span class="text-sm text-gray-500 dark:text-gray-400">{{ row.lastSeen }}</span>
        </template>
      </UTable>

      <!-- Empty state -->
      <div v-if="filteredPlayers.length === 0" class="text-center py-12">
        <UIcon name="i-heroicons-users-20-solid" class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" />
        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
          {{ searchQuery || selectedStatus !== 'all' ? 'No players found' : 'No players online' }}
        </h3>
        <p class="text-gray-500 dark:text-gray-400 mb-6">
          {{ searchQuery || selectedStatus !== 'all' 
            ? 'Try adjusting your search or filters' 
            : 'Players will appear here when they join servers' }}
        </p>
      </div>
    </UCard>
  </div>
</template> 