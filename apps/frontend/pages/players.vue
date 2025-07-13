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

// Filter configurations for SearchAndFilters component
const filters = computed(() => [
  {
    key: 'status',
    value: selectedStatus.value,
    options: [
      { label: 'All Status', value: 'all' },
      { label: 'Online', value: 'online' },
      { label: 'Offline', value: 'offline' }
    ],
    class: 'w-40'
  }
])

// Player stats for StatsCard components
const playerStats = computed(() => [
  {
    key: 'total',
    label: 'Total Players',
    value: players.value.length.toString(),
    icon: 'i-heroicons-users-20-solid',
    color: 'blue'
  },
  {
    key: 'online',
    label: 'Online Now',
    value: players.value.filter(p => p.status === 'online').length.toString(),
    icon: 'i-heroicons-check-circle-20-solid',
    color: 'green'
  },
  {
    key: 'offline',
    label: 'Offline',
    value: players.value.filter(p => p.status === 'offline').length.toString(),
    icon: 'i-heroicons-x-circle-20-solid',
    color: 'gray'
  }
])

// Resolve components for use in cell renderers
const UAvatar = resolveComponent('UAvatar')
const StatusBadge = resolveComponent('StatusBadge')

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
      return h(StatusBadge, {
        status: player.status,
        type: 'player'
      })
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

// Handle filter updates
const handleFilterUpdate = (key: string, value: string) => {
  if (key === 'status') {
    selectedStatus.value = value
  }
}

// Check if filters are active
const hasActiveFilters = computed(() => {
  return searchQuery.value !== '' || selectedStatus.value !== 'all'
})
</script>

<template>
  <div>
    <!-- Page Header -->
    <PageHeader 
      title="Players"
      description="Monitor player activity across all servers"
    />

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
      <StatsCard
        v-for="stat in playerStats"
        :key="stat.key"
        :label="stat.label"
        :value="stat.value"
        :icon="stat.icon"
        :color="stat.color"
      />
    </div>

    <!-- Search and Filters -->
    <SearchAndFilters
      v-model:search-query="searchQuery"
      :filters="filters"
      search-placeholder="Search players..."
      @update:filter="handleFilterUpdate"
      class="mb-6"
    />

    <!-- Players Table -->
    <UCard>
      <UTable :data="filteredPlayers" :columns="columns" />

      <!-- Empty state -->
      <EmptyState
        v-if="filteredPlayers.length === 0"
        icon="i-heroicons-users-20-solid"
        :title="hasActiveFilters ? 'No players found' : 'No players online'"
        :description="hasActiveFilters 
          ? 'Try adjusting your search or filters' 
          : 'Players will appear here when they join servers'"
      />
    </UCard>
  </div>
</template> 