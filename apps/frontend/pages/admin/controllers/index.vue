<script setup lang="ts">
import { h, resolveComponent } from 'vue'

definePageMeta({
  layout: 'default',
  middleware: ['auth', 'admin']
})

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

// Use admin composable for real API data
const { 
  controllers, 
  isLoading, 
  error,
  fetchControllers, 
  refreshControllers, 
  cleanupInactiveControllers,
  restartController,
  removeController,
  clearError
} = useAdmin()

// Local state
const searchQuery = ref('')
const selectedStatus = ref('all')
const selectedCluster = ref('all')

// Load controllers on mount
onMounted(async () => {
  console.log('Controllers list page mounted, route:', route.path)
  await loadControllers()
})

// Load controllers data
const loadControllers = async () => {
  try {
    await fetchControllers()
  } catch (err: any) {
    console.error('Failed to load controllers:', err)
    // Error is already handled by the composable
  }
}

// Filter configurations for SearchAndFilters component
const filters = computed(() => [
  {
    key: 'status',
    value: selectedStatus.value,
    options: [
      { label: 'All Status', value: 'all' },
      { label: 'Online', value: 'online' },
      { label: 'Offline', value: 'offline' },
      { label: 'Error', value: 'error' }
    ],
    class: 'w-40'
  },
  {
    key: 'cluster',
    value: selectedCluster.value,
    options: [
      { label: 'All Clusters', value: 'all' },
      { label: 'Production', value: 'prod' },
      { label: 'Staging', value: 'staging' },
      { label: 'Development', value: 'dev' }
    ],
    class: 'w-40'
  }
])

// Page header actions
const headerActions = computed(() => [
  {
    label: t('admin.controllers.refreshControllers'),
    icon: 'i-heroicons-arrow-path-20-solid',
    color: 'primary' as const,
    onClick: () => handleRefreshControllers(),
    loading: isLoading.value
  },
])

// Controller stats for StatsCard components
const controllerStats = computed(() => [
  {
    key: 'total',
    label: t('admin.controllers.totalControllers'),
    value: controllers.value.length.toString(),
    icon: 'i-heroicons-cpu-chip-20-solid',
    color: 'blue'
  },
  {
    key: 'online',
    label: t('admin.controllers.online'),
    value: controllers.value.filter(c => c.is_online).length.toString(),
    icon: 'i-heroicons-check-circle-20-solid',
    color: 'green'
  },
  {
    key: 'offline',
    label: t('admin.controllers.offline'),
    value: controllers.value.filter(c => !c.is_online).length.toString(),
    icon: 'i-heroicons-x-circle-20-solid',
    color: 'red'
  },
  {
    key: 'errors',
    label: t('admin.controllers.errors'),
    value: controllers.value.filter(c => c.status === 'error').length.toString(),
    icon: 'i-heroicons-exclamation-triangle-20-solid',
    color: 'yellow'
  }
])

// Resolve components for use in cell renderers
const UIcon = resolveComponent('UIcon')
const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')
const UDropdownMenu = resolveComponent('UDropdownMenu')
const StatusBadge = resolveComponent('StatusBadge')

// Helper functions
const getStatusColor = (status: string, isOnline: boolean) => {
  if (!isOnline) return 'error'
  if (status === 'error') return 'error'
  if (status === 'active') return 'success'
  return 'warning'
}

const formatUptime = (uptime: string) => {
  if (!uptime) return 'Offline'
  
  // Parse the uptime string (e.g., "7m22.705186214s")
  const match = uptime.match(/(\d+)m(\d+\.\d+)s/)
  if (!match) return uptime
  
  const minutes = parseInt(match[1])
  const seconds = parseFloat(match[2])
  
  if (minutes < 60) {
    return `${minutes}m ${Math.floor(seconds)}s`
  }
  
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  
  if (hours < 24) {
    return `${hours}h ${remainingMinutes}m`
  }
  
  const days = Math.floor(hours / 24)
  const remainingHours = hours % 24
  
  return `${days}d ${remainingHours}h`
}

const formatLastHeartbeat = (timestamp: string) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / (1000 * 60))
  
  if (minutes < 1) return 'Just now'
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
}

// Action items for dropdown
const getControllerActions = (controller: any) => [
  {
    label: t('admin.controllers.viewDetails'),
    icon: 'i-heroicons-eye-20-solid',
    onClick: () => viewControllerDetails(controller)
  },
  {
    label: t('admin.controllers.viewLogs'),
    icon: 'i-heroicons-document-text-20-solid',
    onClick: () => viewControllerLogs(controller)
  },
  {
    label: t('admin.controllers.restartController'),
    icon: 'i-heroicons-arrow-path-20-solid',
    onClick: () => handleRestartController(controller),
    disabled: !controller.is_online
  },
  {
    label: t('admin.controllers.removeController'),
    icon: 'i-heroicons-trash-20-solid',
    onClick: () => handleRemoveController(controller),
    color: 'red' as const
  }
]

// Table columns configuration
const columns = computed(() => [
  {
    accessorKey: 'cluster_name',
    header: t('admin.controllers.clusterName'),
    cell: ({ row }: any) => {
      const controller = row.original
      return h('div', { class: 'flex items-center space-x-3' }, [
        h(UIcon, { name: 'i-heroicons-server-20-solid', class: 'w-4 h-4 text-gray-400' }),
        h('div', [
          h('button', {
            onClick: () => viewControllerDetails(controller),
            class: 'font-medium text-gray-900 dark:text-gray-100 hover:text-primary-600 dark:hover:text-primary-400 transition-colors text-left'
          }, controller.cluster_name),
          h('div', { class: 'text-xs text-gray-500 dark:text-gray-400' }, controller.cluster_id)
        ])
      ])
    }
  },
  {
    accessorKey: 'status',
    header: t('admin.controllers.status'),
    cell: ({ row }: any) => {
      const controller = row.original
      return h(UBadge, {
        color: controller.is_online ? 'green' : 'red',
        variant: 'subtle',
        size: 'sm'
      }, controller.is_online ? 'Online' : 'Offline')
    }
  },
  {
    accessorKey: 'version',
    header: t('admin.controllers.version'),
    cell: ({ row }: any) => {
      const controller = row.original
      return h(UBadge, {
        color: 'gray',
        variant: 'subtle',
        size: 'sm'
      }, controller.version)
    }
  },
  {
    accessorKey: 'uptime',
    header: t('admin.controllers.uptime'),
    cell: ({ row }: any) => {
      const controller = row.original
      return h('span', { class: 'text-sm text-gray-600 dark:text-gray-400' }, formatUptime(controller.uptime))
    }
  },
  {
    id: 'actions',
    header: t('admin.controllers.actions'),
    cell: ({ row }: any) => {
      const controller = row.original
      return h(UDropdownMenu, {
        items: getControllerActions(controller)
      }, {
        trigger: () => h(UButton, {
          color: 'gray',
          variant: 'ghost',
          icon: 'i-heroicons-ellipsis-vertical-20-solid',
          size: 'xs'
        })
      })
    }
  }
] as any)

// Filtered controllers
const filteredControllers = computed(() => {
  let filtered = controllers.value

  // Apply search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(controller =>
      controller.cluster_name.toLowerCase().includes(query) ||
      controller.cluster_id.toLowerCase().includes(query) ||
      controller.version.toLowerCase().includes(query)
    )
  }

  // Apply status filter
  if (selectedStatus.value !== 'all') {
    if (selectedStatus.value === 'online') {
      filtered = filtered.filter(controller => controller.is_online)
    } else if (selectedStatus.value === 'offline') {
      filtered = filtered.filter(controller => !controller.is_online)
    } else if (selectedStatus.value === 'error') {
      filtered = filtered.filter(controller => controller.status === 'error')
    }
  }

  // Apply cluster filter
  if (selectedCluster.value !== 'all') {
    filtered = filtered.filter(controller => 
      controller.cluster_id.toLowerCase().includes(selectedCluster.value)
    )
  }

  return filtered
})

// Check if filters are active
const hasActiveFilters = computed(() => {
  return searchQuery.value !== '' || selectedStatus.value !== 'all' || selectedCluster.value !== 'all'
})

// Action handlers
const handleRefreshControllers = async () => {
  try {
    await refreshControllers()
    const toast = useToast()
    toast.add({
      title: 'Controllers Refreshed',
      description: 'Controller data has been updated',
      color: 'success'
    })
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Refresh Failed',
      description: err?.data?.message || 'Failed to refresh controllers',
      color: 'error'
    })
  }
}

const handleCleanupInactiveControllers = async () => {
  const confirmed = confirm('Are you sure you want to cleanup inactive controllers? This action cannot be undone.')
  if (!confirmed) return

  try {
    await cleanupInactiveControllers()
    const toast = useToast()
    toast.add({
      title: 'Cleanup Complete',
      description: 'Inactive controllers have been cleaned up',
      color: 'success'
    })
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Cleanup Failed',
      description: err?.data?.message || 'Failed to cleanup inactive controllers',
      color: 'error'
    })
  }
}

const viewControllerDetails = (controller: any) => {
  router.push(`/admin/controllers/${controller.id}`)
}

const viewControllerLogs = (controller: any) => {
  router.push(`/admin/controllers/${controller.id}/logs`)
}

const handleRestartController = async (controller: any) => {
  try {
    await restartController(controller.id)
    const toast = useToast()
    toast.add({
      title: 'Controller Restarted',
      description: `${controller.cluster_name} has been restarted`,
      color: 'success'
    })
    await loadControllers()
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Restart Failed',
      description: err?.data?.message || 'Failed to restart controller',
      color: 'error'
    })
  }
}

const handleRemoveController = async (controller: any) => {
  const confirmed = confirm(`Are you sure you want to remove ${controller.cluster_name}? This action cannot be undone.`)
  if (!confirmed) return

  try {
    await removeController(controller.id)
    const toast = useToast()
    toast.add({
      title: 'Controller Removed',
      description: `${controller.cluster_name} has been removed`,
      color: 'success'
    })
    await loadControllers()
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Remove Failed',
      description: err?.data?.message || 'Failed to remove controller',
      color: 'error'
    })
  }
}

const handleFilterUpdate = (key: string, value: string) => {
  if (key === 'status') {
    selectedStatus.value = value
  } else if (key === 'cluster') {
    selectedCluster.value = value
  }
}
</script>

<template>
  <div>
    <!-- Page Header -->
    <PageHeader 
      :title="t('admin.controllers.title')"
      :description="t('admin.controllers.description')"
      :actions="headerActions"
    >
      <template #extra>
        <UBadge color="success" variant="subtle">
          <UIcon name="i-heroicons-shield-check-20-solid" class="w-4 h-4 mr-1" />
          Admin Access
        </UBadge>
      </template>
    </PageHeader>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <StatsCard
        v-for="stat in controllerStats"
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
      search-placeholder="Search controllers..."
      @update:filter="handleFilterUpdate"
      class="mb-6"
    />

    <!-- Controllers Table -->
    <UCard>
      <UTable :data="filteredControllers" :columns="columns" />

      <!-- Empty state -->
      <EmptyState
        v-if="filteredControllers.length === 0"
        icon="i-heroicons-cpu-chip-20-solid"
        :title="hasActiveFilters ? 'No controllers found' : t('admin.controllers.noControllers')"
        :description="hasActiveFilters 
          ? 'Try adjusting your search or filters' 
          : t('admin.controllers.noControllersDesc')"
      />
    </UCard>
  </div>
</template> 