<script setup lang="ts">
import { h, resolveComponent } from 'vue'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const { t } = useI18n()

// Define Alert interface
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

// Mock alerts data
const alerts = ref<Alert[]>([
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

// Filter configurations for SearchAndFilters component
const filters = computed(() => [
  {
    key: 'severity',
    value: selectedSeverity.value,
    options: [
      { label: 'All Severities', value: 'all' },
      { label: 'Critical', value: 'critical' },
      { label: 'Error', value: 'error' },
      { label: 'Warning', value: 'warning' },
      { label: 'Info', value: 'info' }
    ],
    class: 'w-40'
  },
  {
    key: 'status',
    value: selectedStatus.value,
    options: [
      { label: 'All Status', value: 'all' },
      { label: 'Active', value: 'active' },
      { label: 'Acknowledged', value: 'acknowledged' }
    ],
    class: 'w-40'
  }
])

// Page header actions
const headerActions = computed(() => [
  {
    label: 'Acknowledge All',
    icon: 'i-heroicons-check-20-solid',
    color: 'primary' as const,
    onClick: () => acknowledgeAll()
  },
  {
    label: 'Refresh',
    icon: 'i-heroicons-arrow-path-20-solid',
    color: 'neutral' as const,
    variant: 'ghost' as const,
    onClick: () => window.location.reload()
  }
])

// Alert stats for StatsCard components
const alertStats = computed(() => [
  {
    key: 'total',
    label: 'Total Alerts',
    value: alerts.value.length.toString(),
    icon: 'i-heroicons-bell-20-solid',
    color: 'blue'
  },
  {
    key: 'active',
    label: 'Active Alerts',
    value: alerts.value.filter(a => !a.acknowledged).length.toString(),
    icon: 'i-heroicons-exclamation-triangle-20-solid',
    color: 'red'
  },
  {
    key: 'critical',
    label: 'Critical',
    value: alerts.value.filter(a => a.severity === 'critical').length.toString(),
    icon: 'i-heroicons-fire-20-solid',
    color: 'red'
  },
  {
    key: 'warnings',
    label: 'Warnings',
    value: alerts.value.filter(a => a.severity === 'warning').length.toString(),
    icon: 'i-heroicons-exclamation-circle-20-solid',
    color: 'yellow'
  }
])

// Resolve components for use in cell renderers
const UIcon = resolveComponent('UIcon')
const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')
const StatusBadge = resolveComponent('StatusBadge')

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

const acknowledgeAll = () => {
  alerts.value.forEach(alert => {
    alert.acknowledged = true
  })
}

const columns: any[] = [
  {
    accessorKey: 'title',
    header: 'Alert',
    cell: ({ row }: any) => {
      const alert = row.original
      return h('div', { class: 'flex items-start gap-3' }, [
        h('div', {
          class: [
            'flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center',
            alert.severity === 'error' || alert.severity === 'critical' ? 'bg-red-100 dark:bg-red-900' : 
            alert.severity === 'warning' ? 'bg-yellow-100 dark:bg-yellow-900' : 
            'bg-blue-100 dark:bg-blue-900'
          ]
        }, [
          h(UIcon, {
            name: alert.icon,
            class: [
              'w-5 h-5',
              alert.severity === 'error' || alert.severity === 'critical' ? 'text-red-600 dark:text-red-400' : 
              alert.severity === 'warning' ? 'text-yellow-600 dark:text-yellow-400' : 
              'text-blue-600 dark:text-blue-400'
            ]
          })
        ]),
        h('div', [
          h('p', { class: 'font-medium text-gray-900 dark:text-gray-100' }, alert.title),
          h('p', { class: 'text-sm text-gray-500 dark:text-gray-400' }, alert.message)
        ])
      ])
    }
  },
  {
    accessorKey: 'server',
    header: 'Server',
    cell: ({ row }: any) => {
      const alert = row.original
      return h('span', { class: 'text-gray-600 dark:text-gray-400' }, alert.server)
    }
  },
  {
    accessorKey: 'severity',
    header: 'Severity',
    cell: ({ row }: any) => {
      const alert = row.original
      return h(UBadge, {
        color: getSeverityColor(alert.severity),
        variant: 'subtle',
        class: [
          'capitalize',
          alert.severity === 'critical' ? 'text-red-700 dark:text-red-300' : '',
          alert.severity === 'error' ? 'text-red-700 dark:text-red-300' : '',
          alert.severity === 'warning' ? 'text-yellow-700 dark:text-yellow-300' : '',
          alert.severity === 'info' ? 'text-blue-700 dark:text-blue-300' : ''
        ]
      }, () => alert.severity)
    }
  },
  {
    accessorKey: 'timestamp',
    header: 'Time',
    cell: ({ row }: any) => {
      const alert = row.original
      return h('span', { class: 'text-sm text-gray-500 dark:text-gray-400' }, alert.timestamp)
    }
  },
  {
    accessorKey: 'acknowledged',
    header: 'Status',
    cell: ({ row }: any) => {
      const alert = row.original
      return h(StatusBadge, {
        status: alert.acknowledged ? 'acknowledged' : 'active'
      })
    }
  },
  {
    id: 'actions',
    header: 'Actions',
    cell: ({ row }: any) => {
      const alert = row.original
      return h('div', { class: 'flex items-center gap-2' }, [
        !alert.acknowledged ? h(UButton, {
          color: 'primary',
          variant: 'ghost',
          size: 'sm',
          icon: 'i-heroicons-check-20-solid',
          class: 'text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-200',
          onClick: () => acknowledgeAlert(alert.id)
        }, () => 'Acknowledge') : null,
        h(UButton, {
          color: 'error',
          variant: 'ghost',
          size: 'sm',
          icon: 'i-heroicons-x-mark-20-solid',
          class: 'text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-200',
          onClick: () => dismissAlert(alert.id)
        }, () => 'Dismiss')
      ])
    }
  }
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

const handleFilterUpdate = (newFilters: any) => {
  selectedSeverity.value = newFilters.severity
  selectedStatus.value = newFilters.status
}

const hasActiveFilters = computed(() => {
  return searchQuery.value || selectedSeverity.value !== 'all' || selectedStatus.value !== 'all'
})
</script>

<template>
  <div>
    <!-- Page Header -->
    <PageHeader 
      title="System Alerts"
      description="Monitor and manage system alerts and notifications"
      :actions="headerActions"
    />

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <StatsCard
        v-for="stat in alertStats"
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
      search-placeholder="Search alerts..."
      @update:filter="handleFilterUpdate"
      class="mb-6"
    />

    <!-- Alerts Table -->
    <UCard>
      <UTable :data="filteredAlerts" :columns="columns" />

      <!-- Empty state -->
      <EmptyState
        v-if="filteredAlerts.length === 0"
        icon="i-heroicons-bell-slash-20-solid"
        :title="hasActiveFilters ? 'No alerts found' : 'No active alerts'"
        :description="hasActiveFilters 
          ? 'Try adjusting your search or filters' 
          : 'All systems are running smoothly'"
      />
    </UCard>
  </div>
</template> 