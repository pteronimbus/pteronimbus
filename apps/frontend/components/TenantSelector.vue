<template>
  <div class="tenant-selector">
    <UDropdownMenu :items="dropdownItems" :ui="{
      content: 'max-h-64 overflow-y-auto',
    }" :content="{ class: 'bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700' }">
      <UButton variant="ghost" size="sm"
        class="flex items-center gap-2 text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-100 dark:hover:bg-gray-800">
        <UAvatar v-if="currentTenant" :src="getTenantIcon(currentTenant)" :alt="currentTenant.name" size="xs" />
        <UIcon v-else name="heroicons:server" class="w-4 h-4 text-gray-500 dark:text-gray-400" />
        <span class="text-sm font-medium">
          {{ currentTenant?.name || 'Select Server' }}
        </span>
        <UIcon name="lucide:chevron-down" class="w-3 h-3 text-gray-400 dark:text-gray-500" />
      </UButton>

      <template #add-new-item>
        <AddTenantModal :available-guilds="availableGuildsForTenant" @refresh="loadAllData">
          <div @click.stop
            class="w-full flex items-center gap-2 p-1.5 text-sm text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer">
            <UIcon name="i-heroicons-plus" class="w-4 h-4" />
            <span>Add New Server</span>
          </div>
        </AddTenantModal>
      </template>
    </UDropdownMenu>
  </div>
</template>

<script setup lang="ts">
const {
  tenants,
  currentTenant,
  availableGuilds,
  fetchUserTenants,
  fetchAvailableGuilds,
  switchTenant
} = useTenant()

// No local state needed - AddTenantModal manages its own state

// Load tenants on mount
onMounted(async () => {
  await loadTenants()
})

// Computed properties
const availableGuildsForTenant = computed(() => {
  if (!availableGuilds.value || !tenants.value) {
    return []
  }
  const tenantGuildIds = new Set(tenants.value.map(tenant => tenant.discord_server_id))
  return availableGuilds.value.filter(guild => !tenantGuildIds.has(guild.id))
})

const dropdownItems = computed(() => {
  const items = []

  // Add tenant items
  if (tenants.value && tenants.value.length > 0) {
    const tenantItems = tenants.value.map(tenant => ({
      label: tenant.name,
      avatar: { src: getTenantIcon(tenant) },
      onSelect: () => selectTenant(tenant),
      checked: currentTenant.value?.id === tenant.id,
      type: 'checkbox' as const
    }))
    items.push(tenantItems)

    // Add separator and "Add New" item
    items.push([
      {
        label: 'Add New Server',
        icon: 'i-heroicons-plus',
        slot: 'add-new-item' as const
      }
    ])
  } else {
    // If no tenants, just show the "Add New" option
    items.push([
      {
        label: 'Add New Server',
        icon: 'i-heroicons-plus',
        slot: 'add-new-item' as const
      }
    ])
  }

  return items
})

// Methods
const loadTenants = async () => {
  try {
    await fetchUserTenants()
  } catch (error) {
    console.error('Failed to load tenants:', error)
  }
}

const loadAllData = async () => {
  try {
    await Promise.all([
      fetchUserTenants(),
      fetchAvailableGuilds()
    ])
  } catch (error) {
    console.error('Failed to load data:', error)
  }
}

const selectTenant = async (tenant: any) => {
  try {
    await switchTenant(tenant)
  } catch (error) {
    console.error('Failed to switch tenant:', error)
  }
}



const getTenantIcon = (tenant: any) => {
  if (tenant.icon) {
    return `https://cdn.discordapp.com/icons/${tenant.discord_server_id}/${tenant.icon}.png`
  }
  return undefined
}

// No watch needed - AddTenantModal handles its own state
</script>

<style scoped>
.tenant-selector {
  position: relative;
}
</style>