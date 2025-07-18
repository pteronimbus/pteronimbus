<template>
  <div class="h-screen flex flex-col bg-gray-50 dark:bg-gray-900">
    <!-- Fixed Header -->
    <header class="sticky top-0 z-50 w-full border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 shadow-sm">
      <div class="w-full max-w-none">
        <div class="flex justify-between items-center h-16 px-4 sm:px-6 lg:px-8">
          <div class="flex items-center gap-4 min-w-0 flex-1">
            <button 
              @click="router.push('/')"
              class="text-2xl font-bold text-gray-900 dark:text-white hover:text-primary-600 dark:hover:text-primary-400 transition-colors cursor-pointer flex-shrink-0"
            >
              Pteronimbus
            </button>
            <div class="flex-shrink-0">
              <TenantSelector />
            </div>
          </div>
          <div class="flex items-center gap-4 flex-shrink-0">
            <UDropdownMenu :items="userMenuItems" :ui="{ content: 'dark:bg-gray-800 dark:border-gray-700 dark:text-white' }">
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
        </div>
      </div>
    </header>

    <!-- Main Content Area -->
    <div class="flex-1 w-full overflow-y-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full">
        <main class="py-6 w-full">
          <slot />
        </main>
      </div>
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