<template>
  <div>
    <header class="flex justify-between items-center p-4 mb-2 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900">
      <div class="flex items-center gap-4">
        <button 
          @click="router.push('/')"
          class="text-2xl font-bold text-gray-900 dark:text-white hover:text-primary-600 dark:hover:text-primary-400 transition-colors cursor-pointer"
        >
          Pteronimbus
        </button>
        <TenantSelector />
      </div>
      <div class="flex items-center gap-4">
        <UDropdownMenu :items="userMenuItems" class=" dark:text-white" :ui="{ content: 'dark:bg-gray-800 dark:border-gray-700 dark:text-white' }">
          <UButton
            color="neutral"
            variant="ghost"
            icon="i-heroicons-user-circle"
            size="sm"
            class="text-xl"
          />
        </UDropdownMenu>
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
import ThemeSwitcher from '~/components/ThemeSwitcher.vue'
import TenantSelector from '~/components/TenantSelector.vue'

const router = useRouter()
const { user, signOut, initializeAuth } = useAuth()

// Initialize auth state
initializeAuth()

const userMenuItems = computed(() => [
  [{
    label: user.value?.email || user.value?.username || 'User',
    disabled: true
  }],
  [{
    label: 'Sign out',
    icon: 'i-heroicons-arrow-left-on-rectangle',
    click: async () => {
      await signOut()
    }
  }]
])
</script> 