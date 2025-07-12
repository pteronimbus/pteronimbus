<template>
  <div>
    <header class="flex justify-between items-center p-4 mb-2 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900">
      <button 
        @click="router.push('/')"
        class="text-2xl font-bold text-gray-900 dark:text-white hover:text-primary-600 dark:hover:text-primary-400 transition-colors cursor-pointer"
      >
        Pteronimbus
      </button>
      <div class="flex items-center gap-4">
        <UserMenu />
        <ThemeSwitcher />
      </div>
    </header>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <main class="py-6">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import UserMenu from '~/components/UserMenu.vue'
import ThemeSwitcher from '~/components/ThemeSwitcher.vue'

const router = useRouter()
</script> 