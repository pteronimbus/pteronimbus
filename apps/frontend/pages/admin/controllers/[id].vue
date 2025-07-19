<script setup lang="ts">
definePageMeta({
  layout: 'default',
  middleware: ['auth', 'admin']
})

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const controllerId = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id

// Use admin composable for real API data
const { 
  controllers, 
  isLoading, 
  error,
  fetchControllers, 
  approveController,
  rejectController,
  clearError
} = useAdmin()

// Get the specific controller
const controller = computed(() => {
  return controllers.value.find(c => c.id === controllerId)
})

// Load controllers on mount
onMounted(async () => {
  console.log('Controller details page mounted, controllerId:', controllerId)
  console.log('Current route:', route.path)
  console.log('Current controllers:', controllers.value)
  
  if (controllers.value.length === 0) {
    console.log('No controllers loaded, fetching...')
    await fetchControllers()
  }
  
  console.log('Controller found:', controller.value)
})

const activeTab = ref('overview')

const tabs = [
  { key: 'overview', label: 'Overview', icon: 'i-heroicons-home-20-solid' },
  { key: 'metrics', label: 'Metrics', icon: 'i-heroicons-chart-bar-20-solid' },
  { key: 'heartbeats', label: 'Heartbeats', icon: 'i-heroicons-heart-20-solid' },
  { key: 'logs', label: 'Logs', icon: 'i-heroicons-document-text-20-solid' }
]

// Format functions
const formatUptime = (uptime: string) => {
  if (!uptime) return 'Offline'
  
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

const formatCreatedAt = (timestamp: string) => {
  return new Date(timestamp).toLocaleString()
}

// Controller stats for StatsCard components
const controllerStats = computed(() => {
  if (!controller.value) return []
  
  return [
    {
      key: 'status',
      label: 'Status',
      value: t(`admin.controllers.statuses.${controller.value.status}`),
      icon: controller.value.status === 'pending_approval' ? 'i-heroicons-clock-20-solid' :
            controller.value.is_online ? 'i-heroicons-check-circle-20-solid' : 'i-heroicons-x-circle-20-solid',
      color: controller.value.status === 'pending_approval' ? 'yellow' :
             controller.value.status === 'degraded' ? 'orange' :
             controller.value.is_online ? 'green' : 'red'
    },
    {
      key: 'uptime',
      label: 'Uptime',
      value: formatUptime(controller.value.uptime),
      icon: 'i-heroicons-clock-20-solid',
      color: 'blue'
    },
    {
      key: 'version',
      label: 'Version',
      value: controller.value.version,
      icon: 'i-heroicons-tag-20-solid',
      color: 'gray'
    },
    {
      key: 'last_heartbeat',
      label: 'Last Heartbeat',
      value: formatLastHeartbeat(controller.value.last_heartbeat),
      icon: 'i-heroicons-heart-20-solid',
      color: 'purple'
    }
  ]
})

// Page header actions
const headerActions = computed(() => {
  const actions = [
    {
      label: 'Back to Controllers',
      icon: 'i-heroicons-arrow-left-20-solid',
      color: 'neutral' as const,
      variant: 'ghost' as const,
      onClick: () => router.push('/admin/controllers')
    },
    {
      label: 'Refresh',
      icon: 'i-heroicons-arrow-path-20-solid',
      color: 'primary' as const,
      onClick: () => fetchControllers(),
      loading: isLoading.value
    }
  ]

  // Add approval actions if controller is pending
  if (controller.value?.status === 'pending_approval') {
    actions.push(
      {
        label: 'Approve',
        icon: 'i-heroicons-check-circle-20-solid',
        color: 'success' as const,
        onClick: () => handleApproveController(),
        loading: isLoading.value
      },
      {
        label: 'Reject',
        icon: 'i-heroicons-x-circle-20-solid',
        color: 'error' as const,
        onClick: () => handleRejectController(),
        loading: isLoading.value
      }
    )
  }

  return actions
})

// Handle not found
watch(controller, (newController) => {
  if (!isLoading.value && !newController) {
    router.push('/admin/controllers')
  }
})

// Handler functions
const handleRefresh = async () => {
  await fetchControllers()
}

const handleApproveController = async () => {
  if (!controller.value) return
  
  const confirmed = confirm(`Are you sure you want to approve ${controller.value.cluster_name}? This action cannot be undone.`)
  if (!confirmed) return

  try {
    await approveController(controller.value.id)
    const toast = useToast()
    toast.add({
      title: 'Controller Approved',
      description: `${controller.value.cluster_name} has been approved`,
      color: 'success'
    })
    await fetchControllers()
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Approve Failed',
      description: err?.data?.message || 'Failed to approve controller',
      color: 'error'
    })
  }
}

const handleRejectController = async () => {
  if (!controller.value) return
  
  const confirmed = confirm(`Are you sure you want to reject ${controller.value.cluster_name}? This action cannot be undone.`)
  if (!confirmed) return

  try {
    await rejectController(controller.value.id)
    const toast = useToast()
    toast.add({
      title: 'Controller Rejected',
      description: `${controller.value.cluster_name} has been rejected`,
      color: 'error'
    })
    await fetchControllers()
  } catch (err: any) {
    const toast = useToast()
    toast.add({
      title: 'Reject Failed',
      description: err?.data?.message || 'Failed to reject controller',
      color: 'error'
    })
  }
}
</script>

<template>
  <div>
    <!-- Loading state -->
    <div v-if="isLoading" class="flex items-center justify-center py-12">
      <UIcon name="i-heroicons-arrow-path-20-solid" class="w-8 h-8 animate-spin text-primary-500" />
      <span class="ml-2 text-gray-600 dark:text-gray-400">Loading controller details...</span>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="text-center py-12">
      <UIcon name="i-heroicons-exclamation-triangle-20-solid" class="w-12 h-12 text-red-500 mx-auto mb-4" />
      <h2 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Error Loading Controller</h2>
      <p class="text-gray-600 dark:text-gray-400 mb-4">{{ error }}</p>
      <UButton @click="handleRefresh" color="primary">Try Again</UButton>
    </div>

    <!-- Controller details -->
    <div v-else-if="controller">
      <!-- Page Header -->
      <PageHeader 
        :title="controller.cluster_name"
        :description="`Controller ID: ${controller.cluster_id}`"
        :actions="headerActions"
      >
        <template #extra>
          <UBadge 
            :color="controller.status === 'pending_approval' ? 'warning' : 
                    controller.status === 'degraded' ? 'orange' :
                    controller.is_online ? 'success' : 'error'" 
            variant="subtle"
            class="ml-2"
          >
            <UIcon 
              :name="controller.status === 'pending_approval' ? 'i-heroicons-clock-20-solid' :
                      controller.is_online ? 'i-heroicons-check-circle-20-solid' : 'i-heroicons-x-circle-20-solid'" 
              class="w-4 h-4 mr-1" 
            />
            {{ t(`admin.controllers.statuses.${controller.status}`) }}
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

      <!-- Tabs -->
      <UCard>
        <UTabs v-model="activeTab" :items="tabs" />

        <!-- Overview Tab -->
        <div v-if="activeTab === 'overview'" class="mt-6 space-y-6">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Controller Information -->
            <div class="space-y-4">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Controller Information</h3>
              <div class="space-y-3">
                <div class="flex justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Controller ID:</span>
                  <span class="font-mono text-sm">{{ controller.id }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Cluster ID:</span>
                  <span class="font-mono text-sm">{{ controller.cluster_id }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Version:</span>
                  <span>{{ controller.version }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Status:</span>
                  <UBadge :color="controller.status === 'pending_approval' ? 'warning' : 
                                   controller.status === 'degraded' ? 'orange' :
                                   controller.is_online ? 'success' : 'error'" variant="subtle">
                    {{ t(`admin.controllers.statuses.${controller.status}`) }}
                  </UBadge>
                </div>
                <div v-if="controller.approved_at" class="flex justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Approved:</span>
                  <span>{{ formatCreatedAt(controller.approved_at) }}</span>
                </div>
                <div v-if="controller.approved_by" class="flex justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Approved By:</span>
                  <span class="font-mono text-sm">{{ controller.approved_by }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Created:</span>
                  <span>{{ formatCreatedAt(controller.created_at) }}</span>
                </div>
              </div>
            </div>

            <!-- Runtime Information -->
            <div class="space-y-4">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Runtime Information</h3>
              <div class="space-y-3">
                <div class="flex justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Uptime:</span>
                  <span>{{ formatUptime(controller.uptime) }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Last Heartbeat:</span>
                  <span>{{ formatLastHeartbeat(controller.last_heartbeat) }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Status:</span>
                  <span>{{ controller.status }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Other tabs placeholder -->
        <div v-else class="mt-6 text-center py-12">
          <UIcon name="i-heroicons-wrench-screwdriver-20-solid" class="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">{{ tabs.find(t => t.key === activeTab)?.label }} Coming Soon</h3>
          <p class="text-gray-600 dark:text-gray-400">This feature is under development.</p>
        </div>
      </UCard>
    </div>

    <!-- Not found state -->
    <div v-else class="text-center py-12">
      <UIcon name="i-heroicons-magnifying-glass-20-solid" class="w-12 h-12 text-gray-400 mx-auto mb-4" />
      <h2 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Controller Not Found</h2>
      <p class="text-gray-600 dark:text-gray-400 mb-4">The controller you're looking for doesn't exist or has been removed.</p>
      <UButton @click="router.push('/admin/controllers')" color="primary">Back to Controllers</UButton>
    </div>
  </div>
</template> 