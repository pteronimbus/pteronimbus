<template>
  <div>
    <header class="flex justify-between items-center p-4 mb-2 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Pteronimbus</h1>
      <ThemeSwitcher />
    </header>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <NuxtPage />
    </div>
  </div>
</template>

<script setup>
import ThemeSwitcher from '~/components/ThemeSwitcher.vue'
</script>
