<script setup lang="ts">
definePageMeta({
  layout: 'default'
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

const getUsageColor = (usage: number) => {
  if (usage > 80) return 'error'
  if (usage > 60) return 'warning'
  return 'success'
}
</script>

<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
      <div>
        <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">System Monitoring</h1>
        <p class="mt-1 text-gray-500 dark:text-gray-400">
          Monitor system resources and performance metrics
        </p>
      </div>
      <div class="mt-4 sm:mt-0 flex items-center gap-2">
        <UBadge color="success" variant="subtle">
          <UIcon name="i-heroicons-check-circle-20-solid" class="w-4 h-4 mr-1" />
          System Healthy
        </UBadge>
        <UButton color="neutral" variant="ghost" icon="i-heroicons-arrow-path-20-solid" size="sm">
          Refresh
        </UButton>
      </div>
    </div>

    <!-- System Overview Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <!-- CPU Usage -->
      <UCard>
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">CPU Usage</p>
            <div class="flex items-baseline mt-1">
              <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ systemStats.cpu.usage }}%</p>
            </div>
            <div class="mt-2">
              <UProgress :value="systemStats.cpu.usage" :color="getUsageColor(systemStats.cpu.usage)" size="sm" />
            </div>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ systemStats.cpu.cores }} cores • {{ systemStats.cpu.temperature }}°C</p>
          </div>
          <div class="flex-shrink-0">
            <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-full">
              <UIcon name="i-heroicons-cpu-chip-20-solid" class="w-6 h-6 text-blue-600 dark:text-blue-400" />
            </div>
          </div>
        </div>
      </UCard>

      <!-- Memory Usage -->
      <UCard>
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">Memory Usage</p>
            <div class="flex items-baseline mt-1">
              <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ systemStats.memory.usage }}%</p>
            </div>
            <div class="mt-2">
              <UProgress :value="systemStats.memory.usage" :color="getUsageColor(systemStats.memory.usage)" size="sm" />
            </div>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ systemStats.memory.used }}GB / {{ systemStats.memory.total }}GB</p>
          </div>
          <div class="flex-shrink-0">
            <div class="p-3 bg-green-100 dark:bg-green-900 rounded-full">
              <UIcon name="i-heroicons-circle-stack-20-solid" class="w-6 h-6 text-green-600 dark:text-green-400" />
            </div>
          </div>
        </div>
      </UCard>

      <!-- Disk Usage -->
      <UCard>
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">Disk Usage</p>
            <div class="flex items-baseline mt-1">
              <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ systemStats.disk.usage }}%</p>
            </div>
            <div class="mt-2">
              <UProgress :value="systemStats.disk.usage" :color="getUsageColor(systemStats.disk.usage)" size="sm" />
            </div>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ systemStats.disk.used }}GB / {{ systemStats.disk.total }}GB</p>
          </div>
          <div class="flex-shrink-0">
            <div class="p-3 bg-orange-100 dark:bg-orange-900 rounded-full">
              <UIcon name="i-heroicons-archive-box-20-solid" class="w-6 h-6 text-orange-600 dark:text-orange-400" />
            </div>
          </div>
        </div>
      </UCard>

      <!-- Network I/O -->
      <UCard>
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400">Network I/O</p>
            <div class="flex items-baseline mt-1">
              <p class="text-lg font-bold text-gray-800 dark:text-gray-100">↑{{ systemStats.network.upload }}MB/s</p>
              <p class="text-lg font-bold text-gray-800 dark:text-gray-100 ml-2">↓{{ systemStats.network.download }}MB/s</p>
            </div>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">Total: ↑{{ systemStats.network.totalUpload }}GB ↓{{ systemStats.network.totalDownload }}GB</p>
          </div>
          <div class="flex-shrink-0">
            <div class="p-3 bg-purple-100 dark:bg-purple-900 rounded-full">
              <UIcon name="i-heroicons-signal-20-solid" class="w-6 h-6 text-purple-600 dark:text-purple-400" />
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Charts and Details -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Resource Usage Charts -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">Resource Usage Over Time</h3>
            <UButton color="neutral" variant="ghost" icon="i-heroicons-arrow-path-20-solid" size="sm" />
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
            <UButton color="neutral" variant="ghost" size="sm">View All</UButton>
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
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">Quick Actions</h3>
        </template>
        <div class="grid grid-cols-1 gap-3">
          <UButton color="primary" variant="soft" icon="i-heroicons-arrow-path-20-solid" class="justify-start">
            Restart All Services
          </UButton>
          <UButton color="warning" variant="soft" icon="i-heroicons-trash-20-solid" class="justify-start">
            Clear System Cache
          </UButton>
          <UButton color="success" variant="soft" icon="i-heroicons-document-arrow-down-20-solid" class="justify-start">
            Download System Report
          </UButton>
          <UButton color="primary" variant="soft" icon="i-heroicons-cog-6-tooth-20-solid" class="justify-start">
            System Settings
          </UButton>
        </div>
      </UCard>
    </div>
  </div>
</template> 