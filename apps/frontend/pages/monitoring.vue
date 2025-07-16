<script setup lang="ts">
definePageMeta({
  layout: 'default',
  middleware: ['auth', 'tenant']
})

const { t } = useI18n()

// Mock monitoring data
const systemStats = ref({
  cpu: {
    usage: 75,
    cores: 8,
    temperature: 65,
    processes: 156
  },
  memory: {
    usage: 60,
    total: 16,
    used: 9.6,
    available: 6.4
  },
  disk: {
    usage: 45,
    total: 500,
    used: 225,
    available: 275
  },
  network: {
    upload: 2.5,
    download: 8.3,
    totalUpload: 125.6,
    totalDownload: 543.2
  }
})

const uptime = ref('7 days, 14 hours, 23 minutes')

// Mock historical data for charts
const chartData = ref({
  cpu: [45, 52, 68, 75, 82, 75, 68, 72, 78, 75],
  memory: [35, 45, 55, 60, 65, 60, 58, 62, 64, 60],
  disk: [40, 41, 42, 43, 44, 45, 45, 45, 45, 45],
  network: [12, 15, 18, 22, 25, 20, 18, 16, 14, 12]
})

const activeProcesses = ref([
  { name: 'minecraft-server', cpu: 25.4, memory: 2048, pid: 1234 },
  { name: 'valheim-server', cpu: 18.2, memory: 1536, pid: 1235 },
  { name: 'node', cpu: 8.5, memory: 512, pid: 1236 },
  { name: 'postgres', cpu: 5.2, memory: 256, pid: 1237 },
  { name: 'nginx', cpu: 2.1, memory: 128, pid: 1238 }
])

// Page header actions
const headerActions = computed(() => [
  {
    label: 'Refresh',
    icon: 'i-heroicons-arrow-path-20-solid',
    color: 'neutral' as const,
    variant: 'ghost' as const,
    onClick: () => window.location.reload()
  }
])

// System overview stats for StatsCard components
const monitoringStats = computed(() => [
  {
    key: 'cpu',
    label: 'CPU Usage',
    value: `${systemStats.value.cpu.usage}%`,
    total: '100%',
    icon: 'i-heroicons-cpu-chip-20-solid',
    color: getUsageColor(systemStats.value.cpu.usage),
    description: `${systemStats.value.cpu.cores} cores • ${systemStats.value.cpu.temperature}°C`
  },
  {
    key: 'memory',
    label: 'Memory Usage',
    value: `${systemStats.value.memory.usage}%`,
    total: '100%',
    icon: 'i-heroicons-circle-stack-20-solid',
    color: getUsageColor(systemStats.value.memory.usage),
    description: `${systemStats.value.memory.used}GB / ${systemStats.value.memory.total}GB`
  },
  {
    key: 'disk',
    label: 'Disk Usage',
    value: `${systemStats.value.disk.usage}%`,
    total: '100%',
    icon: 'i-heroicons-archive-box-20-solid',
    color: getUsageColor(systemStats.value.disk.usage),
    description: `${systemStats.value.disk.used}GB / ${systemStats.value.disk.total}GB`
  },
  {
    key: 'network',
    label: 'Network I/O',
    value: `↑${systemStats.value.network.upload}MB/s`,
    icon: 'i-heroicons-signal-20-solid',
    color: 'purple',
    description: `↓${systemStats.value.network.download}MB/s`
  }
])

// Quick actions configuration
const quickActions = computed(() => [
  {
    label: 'Restart All Services',
    icon: 'i-heroicons-arrow-path-20-solid',
    color: 'primary' as const,
    onClick: () => console.log('Restart all services')
  },
  {
    label: 'Clear System Cache',
    icon: 'i-heroicons-trash-20-solid',
    color: 'warning' as const,
    onClick: () => console.log('Clear system cache')
  },
  {
    label: 'Download System Report',
    icon: 'i-heroicons-document-arrow-down-20-solid',
    color: 'success' as const,
    onClick: () => console.log('Download system report')
  },
  {
    label: 'System Settings',
    icon: 'i-heroicons-cog-6-tooth-20-solid',
    color: 'secondary' as const,
    onClick: () => console.log('Open system settings')
  }
])

const getUsageColor = (usage: number) => {
  if (usage > 80) return 'red'
  if (usage > 60) return 'yellow'
  return 'green'
}
</script>

<template>
  <div>
    <!-- Page Header -->
    <PageHeader 
      title="System Monitoring"
      description="Monitor system resources and performance metrics"
      :actions="headerActions"
    >
      <template #extra>
        <UBadge color="success" variant="subtle">
          <UIcon name="i-heroicons-check-circle-20-solid" class="w-4 h-4 mr-1" />
          System Healthy
        </UBadge>
      </template>
    </PageHeader>

    <!-- System Overview Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <StatsCard
        v-for="stat in monitoringStats"
        :key="stat.key"
        :label="stat.label"
        :value="stat.value"
        :total="stat.total"
        :icon="stat.icon"
        :color="stat.color"
        :description="stat.description"
      />
    </div>

    <!-- Charts and Details -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Resource Usage Charts -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">Resource Usage Over Time</h3>
            <UButton 
              color="neutral" 
              variant="ghost" 
              icon="i-heroicons-arrow-path-20-solid" 
              size="sm"
              class="text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200"
            />
          </div>
        </template>
        <div class="h-64 bg-gray-100 dark:bg-gray-800 rounded-md flex items-center justify-center">
          <div class="text-center">
            <UIcon name="i-heroicons-chart-line-20-solid" class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-2" />
            <p class="text-gray-500 dark:text-gray-400">Resource usage chart</p>
            <p class="text-sm text-gray-400 dark:text-gray-500">Chart component would be displayed here</p>
          </div>
        </div>
      </UCard>

      <!-- Top Processes -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">Top Processes</h3>
            <UButton 
              color="neutral" 
              variant="ghost" 
              size="sm"
              class="text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200"
            >
              View All
            </UButton>
          </div>
        </template>
        <div class="space-y-4">
          <div 
            v-for="process in activeProcesses" 
            :key="process.pid"
            class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 bg-blue-100 dark:bg-blue-900 rounded-full flex items-center justify-center">
                <UIcon name="i-heroicons-command-line-20-solid" class="w-4 h-4 text-blue-600 dark:text-blue-400" />
              </div>
              <div>
                <p class="font-medium text-gray-900 dark:text-gray-100">{{ process.name }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400">PID: {{ process.pid }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ process.cpu }}%</p>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ process.memory }}MB</p>
            </div>
          </div>
        </div>
        
        <!-- Empty state for processes -->
        <EmptyState
          v-if="activeProcesses.length === 0"
          icon="i-heroicons-command-line-20-solid"
          title="No active processes"
          description="No processes are currently running on the system."
          size="sm"
        />
      </UCard>

      <!-- System Information -->
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">System Information</h3>
        </template>
        <div class="space-y-4">
          <div class="flex justify-between">
            <span class="text-gray-600 dark:text-gray-400">Uptime:</span>
            <span class="font-medium text-gray-900 dark:text-gray-100">{{ uptime }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600 dark:text-gray-400">Load Average:</span>
            <span class="font-medium text-gray-900 dark:text-gray-100">0.85, 0.72, 0.68</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600 dark:text-gray-400">Active Processes:</span>
            <span class="font-medium text-gray-900 dark:text-gray-100">{{ systemStats.cpu.processes }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600 dark:text-gray-400">Kernel Version:</span>
            <span class="font-medium text-gray-900 dark:text-gray-100">5.4.0-74-generic</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600 dark:text-gray-400">Architecture:</span>
            <span class="font-medium text-gray-900 dark:text-gray-100">x86_64</span>
          </div>
        </div>
      </UCard>

      <!-- Quick Actions -->
      <QuickActions 
        title="Quick Actions"
        :actions="quickActions"
        :grid-cols="1"
      />
    </div>
  </div>
</template> 