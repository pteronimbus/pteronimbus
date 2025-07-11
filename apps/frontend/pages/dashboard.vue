<script setup lang="ts">
definePageMeta({
  layout: 'default'
})

const { t } = useI18n()
const user = useUser()
const router = useRouter()

// Enhanced stats with more comprehensive data
const stats = computed(() => [
  { 
    key: 'activeServers',
    label: t('dashboard.stats.activeServers'),
    value: '12',
    total: '15',
    color: 'emerald',
    icon: 'i-heroicons-server-20-solid',
    route: '/servers',
    trend: '+2',
    trendColor: 'text-green-500'
  },
  { 
    key: 'totalPlayers',
    label: t('dashboard.stats.totalPlayers'),
    value: '150',
    total: '200',
    color: 'blue',
    icon: 'i-heroicons-users-20-solid',
    route: '/players',
    trend: '+12',
    trendColor: 'text-green-500'
  },
  { 
    key: 'totalUsers',
    label: t('dashboard.stats.totalUsers'),
    value: '45',
    total: '50',
    color: 'purple',
    icon: 'i-heroicons-user-group-20-solid',
    route: '/users',
    trend: '+3',
    trendColor: 'text-green-500'
  },
  { 
    key: 'onlineUsers',
    label: t('dashboard.stats.onlineUsers'),
    value: '28',
    total: '45',
    color: 'green',
    icon: 'i-heroicons-signal-20-solid',
    route: '/users?status=online',
    trend: '+5',
    trendColor: 'text-green-500'
  },
  { 
    key: 'cpuUsage',
    label: t('dashboard.stats.cpuUsage'),
    value: '75%',
    color: 'yellow',
    icon: 'i-heroicons-cpu-chip-20-solid',
    route: '/monitoring',
    trend: '+5%',
    trendColor: 'text-yellow-500'
  },
  { 
    key: 'memoryUsage',
    label: t('dashboard.stats.memoryUsage'),
    value: '60%',
    color: 'cyan',
    icon: 'i-heroicons-circle-stack-20-solid',
    route: '/monitoring',
    trend: '-2%',
    trendColor: 'text-green-500'
  },
  { 
    key: 'diskUsage',
    label: t('dashboard.stats.diskUsage'),
    value: '45%',
    color: 'orange',
    icon: 'i-heroicons-archive-box-20-solid',
    route: '/monitoring',
    trend: '+1%',
    trendColor: 'text-yellow-500'
  },
  { 
    key: 'alertsActive',
    label: t('dashboard.stats.alertsActive'),
    value: '3',
    color: 'red',
    icon: 'i-heroicons-exclamation-triangle-20-solid',
    route: '/alerts',
    trend: '+1',
    trendColor: 'text-red-500'
  }
])

// Recent activity data
const recentActivity = computed(() => [
  {
    id: 1,
    type: 'server_started',
    message: t('dashboard.activity.serverStarted', { name: 'Minecraft Survival' }),
    timestamp: '2 minutes ago',
    icon: 'i-heroicons-play-circle-20-solid',
    color: 'green'
  },
  {
    id: 2,
    type: 'user_joined',
    message: t('dashboard.activity.userJoined', { name: 'PlayerOne', server: 'Valheim Dedicated' }),
    timestamp: '5 minutes ago',
    icon: 'i-heroicons-user-plus-20-solid',
    color: 'blue'
  },
  {
    id: 3,
    type: 'server_stopped',
    message: t('dashboard.activity.serverStopped', { name: 'CS:GO Competitive' }),
    timestamp: '10 minutes ago',
    icon: 'i-heroicons-stop-circle-20-solid',
    color: 'red'
  },
  {
    id: 4,
    type: 'user_banned',
    message: t('dashboard.activity.userBanned', { name: 'BadPlayer', server: 'Minecraft Survival' }),
    timestamp: '15 minutes ago',
    icon: 'i-heroicons-no-symbol-20-solid',
    color: 'red'
  },
  {
    id: 5,
    type: 'server_created',
    message: t('dashboard.activity.serverCreated', { name: 'Terraria Adventure' }),
    timestamp: '1 hour ago',
    icon: 'i-heroicons-plus-circle-20-solid',
    color: 'purple'
  }
])

// Active alerts
const activeAlerts = computed(() => [
  {
    id: 1,
    type: 'high_cpu',
    message: t('dashboard.alerts.highCpu'),
    severity: 'warning',
    timestamp: '5 minutes ago',
    icon: 'i-heroicons-cpu-chip-20-solid'
  },
  {
    id: 2,
    type: 'server_down',
    message: t('dashboard.alerts.serverDown', { name: 'CS:GO Competitive' }),
    severity: 'error',
    timestamp: '10 minutes ago',
    icon: 'i-heroicons-x-circle-20-solid'
  },
  {
    id: 3,
    type: 'disk_space',
    message: t('dashboard.alerts.diskSpace'),
    severity: 'warning',
    timestamp: '30 minutes ago',
    icon: 'i-heroicons-archive-box-20-solid'
  }
])

// Handle stat card clicks
const handleStatClick = (stat) => {
  if (stat.route) {
    router.push(stat.route)
  }
}

// Resource usage for chart
const resourceData = computed(() => ({
  labels: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00'],
  datasets: [
    {
      label: 'CPU Usage',
      data: [45, 52, 68, 75, 82, 75],
      borderColor: 'rgb(59, 130, 246)',
      backgroundColor: 'rgba(59, 130, 246, 0.1)',
      fill: true
    },
    {
      label: 'Memory Usage',
      data: [35, 45, 55, 60, 65, 60],
      borderColor: 'rgb(16, 185, 129)',
      backgroundColor: 'rgba(16, 185, 129, 0.1)',
      fill: true
    }
  ]
}))
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">{{ t('dashboard.welcome', { name: user?.name || 'User' }) }}</h1>
      <p class="mt-1 text-gray-500 dark:text-gray-400">{{ t('dashboard.overview') }}</p>
    </div>

    <!-- Stats Grid -->
    <div class="mb-8 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <UCard 
        v-for="stat in stats" 
        :key="stat.key"
        :class="[
          'cursor-pointer transition-all duration-200 hover:shadow-lg hover:scale-105',
          stat.route ? 'hover:bg-gray-50 dark:hover:bg-gray-800' : ''
        ]"
        @click="handleStatClick(stat)"
      >
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ stat.label }}</p>
            <div class="flex items-baseline mt-1">
              <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ stat.value }}</p>
              <p v-if="stat.total" class="ml-2 text-sm text-gray-500 dark:text-gray-400">/ {{ stat.total }}</p>
            </div>
            <div v-if="stat.trend" class="flex items-center mt-2">
              <span :class="[stat.trendColor, 'text-xs font-medium']">{{ stat.trend }}</span>
              <span class="ml-1 text-xs text-gray-500 dark:text-gray-400">from last hour</span>
            </div>
          </div>
          <div class="flex-shrink-0">
            <div :class="[
              'p-3 rounded-full',
              stat.color === 'emerald' ? 'bg-emerald-100 dark:bg-emerald-900' : '',
              stat.color === 'blue' ? 'bg-blue-100 dark:bg-blue-900' : '',
              stat.color === 'purple' ? 'bg-purple-100 dark:bg-purple-900' : '',
              stat.color === 'green' ? 'bg-green-100 dark:bg-green-900' : '',
              stat.color === 'yellow' ? 'bg-yellow-100 dark:bg-yellow-900' : '',
              stat.color === 'cyan' ? 'bg-cyan-100 dark:bg-cyan-900' : '',
              stat.color === 'orange' ? 'bg-orange-100 dark:bg-orange-900' : '',
              stat.color === 'red' ? 'bg-red-100 dark:bg-red-900' : ''
            ]">
              <UIcon 
                :name="stat.icon" 
                :class="[
                  'w-6 h-6',
                  stat.color === 'emerald' ? 'text-emerald-600 dark:text-emerald-400' : '',
                  stat.color === 'blue' ? 'text-blue-600 dark:text-blue-400' : '',
                  stat.color === 'purple' ? 'text-purple-600 dark:text-purple-400' : '',
                  stat.color === 'green' ? 'text-green-600 dark:text-green-400' : '',
                  stat.color === 'yellow' ? 'text-yellow-600 dark:text-yellow-400' : '',
                  stat.color === 'cyan' ? 'text-cyan-600 dark:text-cyan-400' : '',
                  stat.color === 'orange' ? 'text-orange-600 dark:text-orange-400' : '',
                  stat.color === 'red' ? 'text-red-600 dark:text-red-400' : ''
                ]"
              />
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Resource Monitoring Chart -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('dashboard.stats.resourceMonitoring') }}</h3>
            <UButton color="gray" variant="ghost" icon="i-heroicons-arrow-path-20-solid" size="sm" />
          </div>
        </template>
        <div class="h-64 bg-gray-100 dark:bg-gray-800 rounded-md flex items-center justify-center">
          <div class="text-center">
            <UIcon name="i-heroicons-chart-bar-20-solid" class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-2" />
            <p class="text-gray-500 dark:text-gray-400">Resource usage chart</p>
            <p class="text-sm text-gray-400 dark:text-gray-500">Chart component would be displayed here</p>
          </div>
        </div>
      </UCard>

      <!-- Recent Activity -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('dashboard.activity.title') }}</h3>
            <UButton color="gray" variant="ghost" size="sm">{{ t('common.viewAll') }}</UButton>
          </div>
        </template>
        <div class="space-y-4">
          <div 
            v-for="activity in recentActivity.slice(0, 5)" 
            :key="activity.id"
            class="flex items-start space-x-3"
          >
            <div :class="[
              'flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center',
              activity.color === 'green' ? 'bg-green-100 dark:bg-green-900' : '',
              activity.color === 'blue' ? 'bg-blue-100 dark:bg-blue-900' : '',
              activity.color === 'red' ? 'bg-red-100 dark:bg-red-900' : '',
              activity.color === 'purple' ? 'bg-purple-100 dark:bg-purple-900' : ''
            ]">
              <UIcon 
                :name="activity.icon" 
                :class="[
                  'w-4 h-4',
                  activity.color === 'green' ? 'text-green-600 dark:text-green-400' : '',
                  activity.color === 'blue' ? 'text-blue-600 dark:text-blue-400' : '',
                  activity.color === 'red' ? 'text-red-600 dark:text-red-400' : '',
                  activity.color === 'purple' ? 'text-purple-600 dark:text-purple-400' : ''
                ]"
              />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm text-gray-800 dark:text-gray-100">{{ activity.message }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ activity.timestamp }}</p>
            </div>
          </div>
          <div v-if="recentActivity.length === 0" class="text-center py-8">
            <UIcon name="i-heroicons-clock-20-solid" class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-2" />
            <p class="text-gray-500 dark:text-gray-400">{{ t('dashboard.activity.noActivity') }}</p>
          </div>
        </div>
      </UCard>

      <!-- System Alerts -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('dashboard.alerts.title') }}</h3>
            <UBadge v-if="activeAlerts.length > 0" :color="activeAlerts.some(a => a.severity === 'error') ? 'red' : 'yellow'" variant="subtle">
              {{ activeAlerts.length }}
            </UBadge>
          </div>
        </template>
        <div class="space-y-4">
          <div 
            v-for="alert in activeAlerts" 
            :key="alert.id"
            class="flex items-start space-x-3 p-3 rounded-lg"
            :class="[
              alert.severity === 'error' ? 'bg-red-50 dark:bg-red-900/20' : 'bg-yellow-50 dark:bg-yellow-900/20'
            ]"
          >
            <UIcon 
              :name="alert.icon" 
              :class="[
                'w-5 h-5 flex-shrink-0',
                alert.severity === 'error' ? 'text-red-600 dark:text-red-400' : 'text-yellow-600 dark:text-yellow-400'
              ]"
            />
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ alert.message }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ alert.timestamp }}</p>
            </div>
            <UButton 
              :color="alert.severity === 'error' ? 'red' : 'yellow'" 
              variant="ghost" 
              size="xs"
              icon="i-heroicons-x-mark-20-solid"
            />
          </div>
          <div v-if="activeAlerts.length === 0" class="text-center py-8">
            <UIcon name="i-heroicons-check-circle-20-solid" class="w-12 h-12 text-green-400 dark:text-green-500 mx-auto mb-2" />
            <p class="text-gray-500 dark:text-gray-400">{{ t('dashboard.alerts.noAlerts') }}</p>
          </div>
        </div>
      </UCard>

      <!-- Quick Actions -->
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">Quick Actions</h3>
        </template>
        <div class="grid grid-cols-2 gap-4">
          <UButton 
            color="blue" 
            variant="soft" 
            size="lg" 
            icon="i-heroicons-plus-circle-20-solid"
            class="justify-start"
            @click="router.push('/servers/create')"
          >
            {{ t('servers.createServer') }}
          </UButton>
          <UButton 
            color="purple" 
            variant="soft" 
            size="lg" 
            icon="i-heroicons-user-plus-20-solid"
            class="justify-start"
            @click="router.push('/users/create')"
          >
            {{ t('users.createUser') }}
          </UButton>
          <UButton 
            color="green" 
            variant="soft" 
            size="lg" 
            icon="i-heroicons-arrow-path-20-solid"
            class="justify-start"
          >
            {{ t('common.refresh') }}
          </UButton>
          <UButton 
            color="orange" 
            variant="soft" 
            size="lg" 
            icon="i-heroicons-cog-6-tooth-20-solid"
            class="justify-start"
            @click="router.push('/settings')"
          >
            Settings
          </UButton>
        </div>
      </UCard>
    </div>
  </div>
</template> 