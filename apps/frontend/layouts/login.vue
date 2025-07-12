<script setup lang="ts">
import ThemeSwitcher from '~/components/ThemeSwitcher.vue'
</script>

<template>
  <div class="relative h-screen flex flex-col items-center justify-center text-center bg-gray-50 dark:bg-gray-900">
    <!-- Theme Switcher in top-right corner -->
    <div class="absolute top-4 right-4 z-10">
      <ThemeSwitcher />
    </div>
    
    <div class="mb-10">
      <h1 class="text-5xl font-bold text-gray-900 dark:text-white">
        Pteronimbus
      </h1>
      <p class="text-gray-600 dark:text-gray-300">
        The cloud-native game server manager
      </p>
    </div>
    <slot />
  </div>
</template> 