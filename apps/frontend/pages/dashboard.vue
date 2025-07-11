<script setup lang="ts">
definePageMeta({
  layout: 'default'
})

const user = useUser()

const stats = [
  { label: 'Active Servers', value: '12' },
  { label: 'Total Players', value: '150' },
  { label: 'CPU Usage', value: '75%' },
  { label: 'Memory Usage', value: '60%' }
]
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">Welcome, {{ user?.name || 'User' }}!</h1>
    <p class="mt-1 text-gray-500 dark:text-gray-400">Here's a quick overview of your game server environment.</p>

    <!-- Stats Grid -->
    <div class="mt-8 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <UCard v-for="stat in stats" :key="stat.label">
        <template #header>
          <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ stat.label }}</h3>
        </template>
        <p class="text-3xl font-bold text-gray-800 dark:text-gray-100">{{ stat.value }}</p>
      </UCard>
    </div>

    <!-- Chart Placeholder -->
    <div class="mt-8">
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">Player Activity</h3>
        </template>
        <div class="h-64 bg-gray-100 dark:bg-gray-800 rounded-md flex items-center justify-center">
          <p class="text-gray-500 dark:text-gray-400">[Chart will be displayed here]</p>
        </div>
      </UCard>
    </div>
  </div>
</template> 