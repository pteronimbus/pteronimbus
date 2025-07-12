<script setup lang="ts">
definePageMeta({
  layout: 'default'
})

const { t } = useI18n()

// Mock alerts data
const alerts = ref([
  {
    id: 1,
    type: 'high_cpu',
    title: 'High CPU Usage Alert',
    message: 'CPU usage has exceeded 85% for more than 5 minutes',
    severity: 'warning',
    timestamp: '5 minutes ago',
    icon: 'i-heroicons-cpu-chip-20-solid',
    server: 'Minecraft Survival',
    acknowledged: false
  },
  {
    id: 2,
    type: 'server_down',
    title: 'Server Offline',
    message: 'CS:GO Competitive server is not responding',
    severity: 'error',
    timestamp: '10 minutes ago',
    icon: 'i-heroicons-x-circle-20-solid',
    server: 'CS:GO Competitive',
    acknowledged: false
  },
  {
    id: 3,
    type: 'disk_space',
    title: 'Low Disk Space',
    message: 'Available disk space is below 15%',
    severity: 'warning',
    timestamp: '30 minutes ago',
    icon: 'i-heroicons-archive-box-20-solid',
    server: 'System',
    acknowledged: true
  },
  {
    id: 4,
    type: 'memory_usage',
    title: 'High Memory Usage',
    message: 'Memory usage has exceeded 90% on Valheim server',
    severity: 'error',
    timestamp: '1 hour ago',
    icon: 'i-heroicons-circle-stack-20-solid',
    server: 'Valheim Dedicated',
    acknowledged: true
  },
  {
    id: 5,
    type: 'network_latency',
    title: 'Network Latency',
    message: 'High network latency detected (>200ms)',
    severity: 'warning',
    timestamp: '2 hours ago',
    icon: 'i-heroicons-signal-20-solid',
    server: 'System',
    acknowledged: false
  }
])

const searchQuery = ref('')
const selectedSeverity = ref('all')
const selectedStatus = ref('all')

const severityOptions = [
  { label: 'All Severities', value: 'all' },
  { label: 'Critical', value: 'critical' },
  { label: 'Error', value: 'error' },
  { label: 'Warning', value: 'warning' },
  { label: 'Info', value: 'info' }
]

const statusOptions = [
  { label: 'All Status', value: 'all' },
  { label: 'Active', value: 'active' },
  { label: 'Acknowledged', value: 'acknowledged' }
]

interface Alert {
  id: number
  type: string
  title: string
  message: string
  severity: string
  timestamp: string
  icon: string
  server: string
  acknowledged: boolean
}

const columns = [
  { key: 'alert', label: 'Alert', id: 'alert' },
  { key: 'server', label: 'Server', id: 'server' },
  { key: 'severity', label: 'Severity', id: 'severity' },
  { key: 'timestamp', label: 'Time', id: 'timestamp' },
  { key: 'status', label: 'Status', id: 'status' },
  { key: 'actions', label: 'Actions', id: 'actions' }
]

const filteredAlerts = computed(() => {
  let filtered = alerts.value

  if (searchQuery.value) {
    filtered = filtered.filter(alert =>
      alert.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      alert.message.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      alert.server.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  if (selectedSeverity.value !== 'all') {
    filtered = filtered.filter(alert => alert.severity === selectedSeverity.value)
  }

  if (selectedStatus.value !== 'all') {
    if (selectedStatus.value === 'active') {
      filtered = filtered.filter(alert => !alert.acknowledged)
    } else if (selectedStatus.value === 'acknowledged') {
      filtered = filtered.filter(alert => alert.acknowledged)
    }
  }

  return filtered
})

const getSeverityColor = (severity: string) => {
  switch (severity) {
    case 'critical': return 'error'
    case 'error': return 'error'
    case 'warning': return 'warning'
    case 'info': return 'primary'
    default: return 'neutral'
  }
}

const acknowledgeAlert = (alertId: number) => {
  const alert = alerts.value.find(a => a.id === alertId)
  if (alert) {
    alert.acknowledged = true
  }
}

const dismissAlert = (alertId: number) => {
  const index = alerts.value.findIndex(a => a.id === alertId)
  if (index !== -1) {
    alerts.value.splice(index, 1)
  }
}

const alertStats = computed(() => ({
  total: alerts.value.length,
  active: alerts.value.filter(a => !a.acknowledged).length,
  critical: alerts.value.filter(a => a.severity === 'critical').length,
  warnings: alerts.value.filter(a => a.severity === 'warning').length
}))
</script>

<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
      <div>
        <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">System Alerts</h1>
        <p class="mt-1 text-gray-500 dark:text-gray-400">
          Monitor and manage system alerts and notifications
        </p>
      </div>
      <div class="mt-4 sm:mt-0 flex items-center gap-2">
        <UButton 
          color="primary" 
          icon="i-heroicons-check-20-solid" 
          size="sm"
          class="text-primary-700 dark:text-primary-300"
        >
          Acknowledge All
        </UButton>
        <UButton 
          color="neutral" 
          variant="ghost" 
          icon="i-heroicons-arrow-path-20-solid" 
          size="sm"
          class="text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200"
        >
          Refresh
        </UButton>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Total Alerts</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ alertStats.total }}</p>
          </div>
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-full">
            <UIcon name="i-heroicons-bell-20-solid" class="w-6 h-6 text-blue-600 dark:text-blue-400" />
          </div>
        </div>
      </UCard>
      
      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Active Alerts</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ alertStats.active }}</p>
          </div>
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-full">
            <UIcon name="i-heroicons-exclamation-triangle-20-solid" class="w-6 h-6 text-red-600 dark:text-red-400" />
          </div>
        </div>
      </UCard>
      
      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Critical</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ alertStats.critical }}</p>
          </div>
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-full">
            <UIcon name="i-heroicons-fire-20-solid" class="w-6 h-6 text-red-600 dark:text-red-400" />
          </div>
        </div>
      </UCard>
      
      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Warnings</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ alertStats.warnings }}</p>
          </div>
          <div class="p-3 bg-yellow-100 dark:bg-yellow-900 rounded-full">
            <UIcon name="i-heroicons-exclamation-circle-20-solid" class="w-6 h-6 text-yellow-600 dark:text-yellow-400" />
          </div>
        </div>
      </UCard>
    </div>

    <!-- Filters -->
    <div class="mb-6 flex flex-col sm:flex-row gap-4">
      <div class="flex-1">
        <UInput
          v-model="searchQuery"
          placeholder="Search alerts..."
          icon="i-heroicons-magnifying-glass-20-solid"
          size="md"
        />
      </div>
      <div class="flex gap-2">
        <USelect
          v-model="selectedSeverity"
          :options="severityOptions"
          size="md"
          class="w-40"
        />
        <USelect
          v-model="selectedStatus"
          :options="statusOptions"
          size="md"
          class="w-40"
        />
      </div>
    </div>

    <!-- Alerts Table -->
    <UCard>
      <UTable :rows="filteredAlerts" :columns="columns">
        <!-- Alert column with icon and details -->
        <template #alert-data="{ row }">
          <div class="flex items-start gap-3">
            <div :class="[
              'flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center',
              (row as Alert).severity === 'error' || (row as Alert).severity === 'critical' ? 'bg-red-100 dark:bg-red-900' : 
              (row as Alert).severity === 'warning' ? 'bg-yellow-100 dark:bg-yellow-900' : 
              'bg-blue-100 dark:bg-blue-900'
            ]">
              <UIcon 
                :name="(row as Alert).icon" 
                :class="[
                  'w-5 h-5',
                  (row as Alert).severity === 'error' || (row as Alert).severity === 'critical' ? 'text-red-600 dark:text-red-400' : 
                  (row as Alert).severity === 'warning' ? 'text-yellow-600 dark:text-yellow-400' : 
                  'text-blue-600 dark:text-blue-400'
                ]"
              />
            </div>
            <div>
              <p class="font-medium text-gray-900 dark:text-gray-100">{{ (row as Alert).title }}</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ (row as Alert).message }}</p>
            </div>
          </div>
        </template>

        <!-- Server column -->
        <template #server-data="{ row }">
          <span class="text-gray-600 dark:text-gray-400">{{ (row as Alert).server }}</span>
        </template>

        <!-- Severity column -->
        <template #severity-data="{ row }">
          <UBadge 
            :color="getSeverityColor((row as Alert).severity)" 
            variant="subtle"
            class="capitalize"
            :class="[
              (row as Alert).severity === 'critical' ? 'text-red-700 dark:text-red-300' : '',
              (row as Alert).severity === 'error' ? 'text-red-700 dark:text-red-300' : '',
              (row as Alert).severity === 'warning' ? 'text-yellow-700 dark:text-yellow-300' : '',
              (row as Alert).severity === 'info' ? 'text-blue-700 dark:text-blue-300' : ''
            ]"
          >
            {{ (row as Alert).severity }}
          </UBadge>
        </template>

        <!-- Timestamp column -->
        <template #timestamp-data="{ row }">
          <span class="text-sm text-gray-500 dark:text-gray-400">{{ (row as Alert).timestamp }}</span>
        </template>

        <!-- Status column -->
        <template #status-data="{ row }">
          <UBadge 
            :color="(row as Alert).acknowledged ? 'success' : 'warning'" 
            variant="subtle"
            :class="[
              (row as Alert).acknowledged ? 'text-green-700 dark:text-green-300' : 'text-yellow-700 dark:text-yellow-300'
            ]"
          >
            {{ (row as Alert).acknowledged ? 'Acknowledged' : 'Active' }}
          </UBadge>
        </template>

        <!-- Actions column -->
        <template #actions-data="{ row }">
          <div class="flex items-center gap-2">
            <UButton 
              v-if="!(row as Alert).acknowledged"
              color="primary" 
              variant="ghost" 
              size="sm"
              icon="i-heroicons-check-20-solid"
              class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-200"
              @click="acknowledgeAlert((row as Alert).id)"
            >
              Acknowledge
            </UButton>
            <UButton 
              color="error" 
              variant="ghost" 
              size="sm"
              icon="i-heroicons-x-mark-20-solid"
              class="text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-200"
              @click="dismissAlert((row as Alert).id)"
            >
              Dismiss
            </UButton>
          </div>
        </template>
      </UTable>

      <!-- Empty state -->
      <div v-if="filteredAlerts.length === 0" class="text-center py-12">
        <UIcon name="i-heroicons-bell-slash-20-solid" class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" />
        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
          {{ searchQuery || selectedSeverity !== 'all' || selectedStatus !== 'all' ? 'No alerts found' : 'No active alerts' }}
        </h3>
        <p class="text-gray-500 dark:text-gray-400">
          {{ searchQuery || selectedSeverity !== 'all' || selectedStatus !== 'all' 
            ? 'Try adjusting your search or filters' 
            : 'All systems are running smoothly' }}
        </p>
      </div>
    </UCard>
  </div>
</template> 