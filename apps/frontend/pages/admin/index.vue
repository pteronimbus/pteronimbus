<script setup lang="ts">
import { h, resolveComponent } from 'vue'

definePageMeta({
  layout: 'default',
  middleware: ['auth', 'admin']
})

const { t } = useI18n()
const router = useRouter()

// Admin navigation items
const adminNavItems = computed(() => [
  {
    title: t('admin.dashboard.controllers.title'),
    description: t('admin.dashboard.controllers.description'),
    icon: 'i-heroicons-cpu-chip-20-solid',
    color: 'blue',
    href: '/admin/controllers',
    stats: {
      total: 0,
      online: 0,
      pending: 0
    }
  },
  {
    title: t('admin.dashboard.users.title'),
    description: t('admin.dashboard.users.description'),
    icon: 'i-heroicons-users-20-solid',
    color: 'green',
    href: '/admin/users',
    stats: {
      total: 0,
      active: 0,
      new: 0
    }
  },
  {
    title: t('admin.dashboard.tenants.title'),
    description: t('admin.dashboard.tenants.description'),
    icon: 'i-heroicons-building-office-20-solid',
    color: 'purple',
    href: '/admin/tenants',
    stats: {
      total: 0,
      active: 0,
      pending: 0
    }
  },
  {
    title: t('admin.dashboard.moderation.title'),
    description: t('admin.dashboard.moderation.description'),
    icon: 'i-heroicons-shield-exclamation-20-solid',
    color: 'orange',
    href: '/admin/moderation',
    stats: {
      reports: 0,
      bans: 0,
      warnings: 0
    }
  }
])

// Audit log data (mock for now)
const auditLogs = ref([
  {
    id: 1,
    action: 'controller_registered',
    user: 'thatguy164.',
    target: 'prod-cluster-01',
    details: 'New controller registered from production cluster',
    timestamp: new Date(Date.now() - 5 * 60 * 1000), // 5 minutes ago
    severity: 'info'
  },
  {
    id: 2,
    action: 'user_signup',
    user: 'newuser123',
    target: 'Discord User',
    details: 'New user signed up via Discord OAuth2',
    timestamp: new Date(Date.now() - 15 * 60 * 1000), // 15 minutes ago
    severity: 'info'
  },
  {
    id: 3,
    action: 'tenant_created',
    user: 'adminuser',
    target: 'Gaming Community Server',
    details: 'New tenant created for Discord server',
    timestamp: new Date(Date.now() - 30 * 60 * 1000), // 30 minutes ago
    severity: 'info'
  },
  {
    id: 4,
    action: 'controller_approved',
    user: 'thatguy164.',
    target: 'staging-cluster-02',
    details: 'Controller approved and activated',
    timestamp: new Date(Date.now() - 45 * 60 * 1000), // 45 minutes ago
    severity: 'success'
  },
  {
    id: 5,
    action: 'user_banned',
    user: 'moderator1',
    target: 'troublemaker',
    details: 'User banned for violating community guidelines',
    timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000), // 2 hours ago
    severity: 'warning'
  },
  {
    id: 6,
    action: 'tenant_removed',
    user: 'adminuser',
    target: 'Inactive Server',
    details: 'Tenant removed due to inactivity',
    timestamp: new Date(Date.now() - 4 * 60 * 60 * 1000), // 4 hours ago
    severity: 'error'
  }
])

// Filter state
const selectedAction = ref('all')
const selectedSeverity = ref('all')
const searchQuery = ref('')

// Filter options
const actionOptions = computed(() => [
  { label: 'All Actions', value: 'all' },
  { label: 'Controller Events', value: 'controller' },
  { label: 'User Events', value: 'user' },
  { label: 'Tenant Events', value: 'tenant' },
  { label: 'Moderation Events', value: 'moderation' }
])

const severityOptions = computed(() => [
  { label: 'All Severities', value: 'all' },
  { label: 'Info', value: 'info' },
  { label: 'Success', value: 'success' },
  { label: 'Warning', value: 'warning' },
  { label: 'Error', value: 'error' }
])

// Filtered audit logs
const filteredAuditLogs = computed(() => {
  let filtered = auditLogs.value

  // Apply search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(log =>
      log.user.toLowerCase().includes(query) ||
      log.target.toLowerCase().includes(query) ||
      log.details.toLowerCase().includes(query)
    )
  }

  // Apply action filter
  if (selectedAction.value !== 'all') {
    filtered = filtered.filter(log => log.action.startsWith(selectedAction.value))
  }

  // Apply severity filter
  if (selectedSeverity.value !== 'all') {
    filtered = filtered.filter(log => log.severity === selectedSeverity.value)
  }

  return filtered
})

// Helper functions
const getSeverityColor = (severity: string) => {
  switch (severity) {
    case 'success': return 'green'
    case 'warning': return 'orange'
    case 'error': return 'red'
    default: return 'blue'
  }
}

const getActionIcon = (action: string) => {
  if (action.startsWith('controller')) return 'i-heroicons-cpu-chip-20-solid'
  if (action.startsWith('user')) return 'i-heroicons-user-20-solid'
  if (action.startsWith('tenant')) return 'i-heroicons-building-office-20-solid'
  if (action.startsWith('moderation') || action.includes('ban')) return 'i-heroicons-shield-exclamation-20-solid'
  return 'i-heroicons-information-circle-20-solid'
}

const formatTimestamp = (timestamp: Date) => {
  const now = new Date()
  const diff = now.getTime() - timestamp.getTime()
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (minutes < 1) return 'Just now'
  if (minutes < 60) return `${minutes}m ago`
  if (hours < 24) return `${hours}h ago`
  return `${days}d ago`
}

// Resolve components
const UIcon = resolveComponent('UIcon')
const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')

// Stats for the dashboard
const dashboardStats = computed(() => [
  {
    key: 'total_controllers',
    label: 'Total Controllers',
    value: '12',
    icon: 'i-heroicons-cpu-chip-20-solid',
    color: 'blue'
  },
  {
    key: 'active_users',
    label: 'Active Users',
    value: '1,247',
    icon: 'i-heroicons-users-20-solid',
    color: 'green'
  },
  {
    key: 'total_tenants',
    label: 'Total Tenants',
    value: '89',
    icon: 'i-heroicons-building-office-20-solid',
    color: 'purple'
  },
  {
    key: 'pending_actions',
    label: 'Pending Actions',
    value: '3',
    icon: 'i-heroicons-clock-20-solid',
    color: 'orange'
  }
])

// Handle filter updates
const handleFilterUpdate = (key: string, value: string) => {
  if (key === 'action') {
    selectedAction.value = value
  } else if (key === 'severity') {
    selectedSeverity.value = value
  }
}
</script>

<template>
  <div>
    <!-- Page Header -->
    <PageHeader 
      title="Admin Dashboard"
      description="Manage system-wide operations, monitor activity, and oversee platform administration"
    >
      <template #extra>
        <UBadge color="success" variant="subtle">
          <UIcon name="i-heroicons-shield-check-20-solid" class="w-4 h-4 mr-1" />
          Super Admin Access
        </UBadge>
      </template>
    </PageHeader>

    <!-- Dashboard Stats -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <StatsCard
        v-for="stat in dashboardStats"
        :key="stat.key"
        :label="stat.label"
        :value="stat.value"
        :icon="stat.icon"
        :color="stat.color"
      />
    </div>

    <!-- Admin Navigation Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <UCard
        v-for="item in adminNavItems"
        :key="item.href"
        class="cursor-pointer hover:shadow-lg transition-shadow"
        @click="router.push(item.href)"
      >
        <div class="flex items-start space-x-3">
          <div :class="`p-2 rounded-lg bg-${item.color}-100 dark:bg-${item.color}-900/20`">
            <UIcon :name="item.icon" :class="`w-6 h-6 text-${item.color}-600 dark:text-${item.color}-400`" />
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="text-sm font-medium text-gray-900 dark:text-white">
              {{ item.title }}
            </h3>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
              {{ item.description }}
            </p>
          </div>
          <UIcon name="i-heroicons-arrow-right-20-solid" class="w-4 h-4 text-gray-400" />
        </div>
      </UCard>
    </div>

    <!-- Audit Log Section -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">
              System Audit Log
            </h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              Recent system-wide activities and administrative actions
            </p>
          </div>
          <UButton
            color="primary"
            variant="ghost"
            icon="i-heroicons-arrow-path-20-solid"
            size="sm"
          >
            Refresh
          </UButton>
        </div>
      </template>

      <!-- Search and Filters -->
      <div class="flex flex-col sm:flex-row gap-4 mb-4">
        <UInput
          v-model="searchQuery"
          placeholder="Search audit logs..."
          icon="i-heroicons-magnifying-glass-20-solid"
          class="flex-1"
        />
        <USelect
          v-model="selectedAction"
          :options="actionOptions"
          placeholder="Filter by action"
          class="w-40"
        />
        <USelect
          v-model="selectedSeverity"
          :options="severityOptions"
          placeholder="Filter by severity"
          class="w-40"
        />
      </div>

      <!-- Audit Log Table -->
      <div class="space-y-3">
        <div
          v-for="log in filteredAuditLogs"
          :key="log.id"
          class="flex items-start space-x-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
        >
          <div :class="`p-2 rounded-lg bg-${getSeverityColor(log.severity)}-100 dark:bg-${getSeverityColor(log.severity)}-900/20`">
            <UIcon :name="getActionIcon(log.action)" :class="`w-4 h-4 text-${getSeverityColor(log.severity)}-600 dark:text-${getSeverityColor(log.severity)}-400`" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center space-x-2 mb-1">
              <span class="text-sm font-medium text-gray-900 dark:text-white">
                {{ log.user }}
              </span>
              <UBadge :color="getSeverityColor(log.severity)" variant="subtle" size="xs">
                {{ log.action.replace('_', ' ').toUpperCase() }}
              </UBadge>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-300">
              {{ log.details }}
            </p>
            <div class="flex items-center space-x-4 mt-2 text-xs text-gray-500 dark:text-gray-400">
              <span>Target: {{ log.target }}</span>
              <span>{{ formatTimestamp(log.timestamp) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <EmptyState
        v-if="filteredAuditLogs.length === 0"
        icon="i-heroicons-document-text-20-solid"
        title="No audit logs found"
        description="No audit logs match your current filters"
      />
    </UCard>
  </div>
</template> 