<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'

const { user, clearUser } = useUser()
const router = useRouter()

const items: DropdownMenuItem[][] = [
  [{
    label: (user.value as any)?.email || 'User',
    slot: 'account',
    disabled: true
  }],
  [{
    label: 'Settings',
    icon: 'i-heroicons-cog-8-tooth'
    // Add click handler to navigate to settings page when created
  }, {
    label: 'Sign out',
    icon: 'i-heroicons-arrow-left-on-rectangle',
    onSelect: () => {
      clearUser()
      router.push('/login')
    }
  }]
]
</script>

<template>
  <UDropdownMenu :items="items" :ui="{ content: 'w-48' }">
    <UAvatar src="https://avatars.githubusercontent.com/u/739984?v=4" />

    <template #account="{ item }">
      <div class="text-left">
        <p>
          Signed in as
        </p>
        <p class="truncate font-medium text-gray-900 dark:text-white">
          {{ (item as any).label }}
        </p>
      </div>
    </template>

    <template #item="{ item }">
      <span class="truncate">{{ (item as any).label }}</span>
      <UIcon v-if="(item as any).icon" :name="(item as any).icon" class="flex-shrink-0 h-4 w-4 text-gray-400 dark:text-gray-500 ms-auto" />
    </template>
  </UDropdownMenu>
</template> 