<script setup lang="ts">
import { h, resolveComponent } from 'vue'

definePageMeta({
  layout: 'default'
})

const { t } = useI18n()
const router = useRouter()

// Define server interface for better type safety
interface Server {
  id: number
  name: string
  game: string
  status: string
  players: string
  ip: string
  port: number
  version: string
  uptime: string
  cpu: number
  memory: number
  createdAt: string
}

// Mock server data - in real app this would come from API
const servers = ref<Server[]>([
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

// Resolve components for use in cell renderers
const UIcon = resolveComponent('UIcon')
const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')
const UDropdownMenu = resolveComponent('UDropdownMenu')

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

// Action items for dropdown
const getActionItems = (server: Server) => [
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

// Action functions
const toggleServer = (server: Server) => {
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

const restartServer = (server: Server) => {
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

const createBackup = (server: Server) => {
  console.log('Creating backup for:', server.name)
  // Implementation for creating backup
}

const deleteServer = (server: Server) => {
  console.log('Deleting server:', server.name)
  // Implementation for server deletion
}

// Navigation functions
const viewServer = (server: Server) => {
  router.push(`/servers/${server.id}`)
}

const createServer = () => {
  router.push('/servers/create')
}

const columns: any[] = [
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }: any) => {
      const server = row.original
      return h('div', { class: 'flex items-center gap-3' }, [
        h('div', { class: 'flex-shrink-0 w-10 h-10 bg-blue-100 dark:bg-blue-900 rounded-lg flex items-center justify-center' }, [
          h(UIcon, { name: 'i-heroicons-server-20-solid', class: 'w-5 h-5 text-blue-600 dark:text-blue-400' })
        ]),
        h('div', [
          h('button', {
            onClick: () => viewServer(server),
            class: 'font-medium text-gray-900 dark:text-gray-100 hover:text-primary-600 dark:hover:text-primary-400 transition-colors'
          }, server.name),
          h('p', { class: 'text-sm text-gray-500 dark:text-gray-400' }, `${server.ip}:${server.port}`)
        ])
      ])
    }
  },
  {
    accessorKey: 'game',
    header: 'Game',
    cell: ({ row }: any) => {
      const server = row.original
      return h('div', [
        h('p', { class: 'font-medium text-gray-900 dark:text-gray-100' }, server.game),
        h('p', { class: 'text-sm text-gray-500 dark:text-gray-400' }, `v${server.version}`)
      ])
    }
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }: any) => {
      const server = row.original
      return h('div', { class: 'space-y-1' }, [
        h(UBadge, { 
          color: getStatusColor(server.status), 
          variant: 'subtle', 
          class: 'capitalize' 
        }, () => t(`servers.status.${server.status}`)),
        server.status === 'online' ? h('div', { class: 'flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400' }, [
          h(UIcon, { name: 'i-heroicons-clock-20-solid', class: 'w-3 h-3' }),
          server.uptime
        ]) : null
      ])
    }
  },
  {
    accessorKey: 'players',
    header: 'Players',
    cell: ({ row }: any) => {
      const server = row.original
      return h('div', { class: 'space-y-1' }, [
        h('p', { class: 'font-medium text-gray-900 dark:text-gray-100' }, server.players),
        server.status === 'online' ? h('div', { class: 'flex items-center gap-2 text-xs' }, [
          h('div', { class: 'flex items-center gap-1' }, [
            h('div', { class: 'w-2 h-2 rounded-full bg-blue-500' }),
            h('span', { class: 'text-gray-500 dark:text-gray-400' }, `CPU: ${server.cpu}%`)
          ]),
          h('div', { class: 'flex items-center gap-1' }, [
            h('div', { class: 'w-2 h-2 rounded-full bg-green-500' }),
            h('span', { class: 'text-gray-500 dark:text-gray-400' }, `MEM: ${server.memory}%`)
          ])
        ]) : null
      ])
    }
  },
  {
    id: 'actions',
    header: 'Actions',
    cell: ({ row }: any) => {
      const server = row.original
      return h(UDropdownMenu, {
        items: getActionItems(server)
      }, () => h(UButton, {
        variant: 'ghost',
        color: 'neutral',
        size: 'sm',
        icon: 'i-heroicons-ellipsis-horizontal-20-solid'
      }))
    }
  }
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
      <UTable :data="filteredServers" :columns="columns" />

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