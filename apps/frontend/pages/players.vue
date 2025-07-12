<script setup lang="ts">
import { h, resolveComponent } from 'vue'

definePageMeta({
  layout: 'default'
})

const { t } = useI18n()
const router = useRouter()

// Define Player interface
interface Player {
  id: number
  name: string
  server: string
  status: string
  playtime: string
  lastSeen: string
  avatar?: string | null
}

// Mock players data
const players = ref<Player[]>([
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

// Resolve components for use in cell renderers
const UAvatar = resolveComponent('UAvatar')
const UBadge = resolveComponent('UBadge')

const getStatusColor = (status: string) => {
  switch (status) {
    case 'online': return 'success'
    case 'offline': return 'neutral'
    default: return 'neutral'
  }
}

const columns: any[] = [
  {
    accessorKey: 'name',
    header: 'Player',
    cell: ({ row }: any) => {
      const player = row.original
      return h('div', { class: 'flex items-center gap-3' }, [
        h(UAvatar, {
          src: player.avatar || undefined,
          alt: player.name,
          size: 'sm'
        }, () => h('span', { class: 'text-xs font-medium text-primary-600 dark:text-primary-400' }, 
          player.name.split('').slice(0, 2).join('').toUpperCase()
        )),
        h('div', [
          h('span', { class: 'font-medium text-gray-900 dark:text-gray-100' }, player.name)
        ])
      ])
    }
  },
  {
    accessorKey: 'server',
    header: 'Server',
    cell: ({ row }: any) => {
      const player = row.original
      return h('span', { class: 'text-gray-600 dark:text-gray-400' }, player.server)
    }
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }: any) => {
      const player = row.original
      return h(UBadge, {
        color: getStatusColor(player.status),
        variant: 'subtle',
        class: [
          'capitalize',
          player.status === 'online' ? 'text-green-700 dark:text-green-300' : '',
          player.status === 'offline' ? 'text-gray-700 dark:text-gray-300' : ''
        ]
      }, () => player.status)
    }
  },
  {
    accessorKey: 'playtime',
    header: 'Playtime',
    cell: ({ row }: any) => {
      const player = row.original
      return h('span', { class: 'text-gray-900 dark:text-gray-100' }, player.playtime)
    }
  },
  {
    accessorKey: 'lastSeen',
    header: 'Last Seen',
    cell: ({ row }: any) => {
      const player = row.original
      return h('span', { class: 'text-sm text-gray-500 dark:text-gray-400' }, player.lastSeen)
    }
  }
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
      <UTable :data="filteredPlayers" :columns="columns" />

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