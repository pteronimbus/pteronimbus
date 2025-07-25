<template>
  <UCard 
    class="h-full transition-all duration-200 cursor-pointer overflow-hidden border border-gray-200 dark:border-gray-700 hover:shadow-lg hover:scale-105 hover:border-primary-300 dark:hover:border-primary-600" 
    @click="$emit('select', tenant)"
  >
    <div class="flex items-center space-x-4 mb-4">
      <UAvatar
        :src="getTenantIcon(tenant) || undefined"
        :alt="tenant.name"
        size="lg"
      />
      <div class="flex-1 min-w-0">
        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 truncate">
          {{ tenant.name }}
        </h3>
        <p class="text-sm text-gray-500 dark:text-gray-400">
          {{ isOwner ? 'Owner' : 'Member' }}
        </p>
      </div>
    </div>
    
    <div class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
      <div class="flex items-center">
        <UIcon name="heroicons:calendar" class="w-4 h-4 mr-2" />
        Added {{ formatDate(tenant.created_at) }}
      </div>
      <div class="flex items-center">
        <UIcon name="heroicons:server" class="w-4 h-4 mr-2" />
        {{ tenant.config?.resource_limits?.max_game_servers || 5 }} server limit
      </div>
    </div>
  </UCard>
</template>

<script setup lang="ts">
interface Props {
  tenant: any
  isOwner: boolean
}

interface Emits {
  (e: 'select', tenant: any): void
  (e: 'delete', tenant: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const getTenantIcon = (tenant: any) => {
  if (tenant.icon) {
    return `https://cdn.discordapp.com/icons/${tenant.discord_server_id}/${tenant.icon}.png`
  }
  return null
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const tenantActions = computed(() => [
  [{
    label: 'Manage Settings',
    icon: 'heroicons:cog-6-tooth',
    click: () => navigateTo(`/tenant/${props.tenant.id}/settings`)
  }],
  [{
    label: 'Remove Server',
    icon: 'heroicons:trash',
    click: () => emit('delete', props.tenant)
  }]
])
</script>

<style scoped>
.tenant-card {
  @apply transition-transform hover:scale-105;
}
</style>