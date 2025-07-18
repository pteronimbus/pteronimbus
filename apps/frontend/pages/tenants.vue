<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="max-w-6xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          {{ t('tenants.title') }}
        </h1>
        <p class="text-lg text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
          Manage game servers for your Discord communities. Select a server to get started.
        </p>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="text-center py-12">
        <UIcon name="i-heroicons-arrow-path" class="w-8 h-8 animate-spin mx-auto mb-4 text-primary-500" />
        <p class="text-gray-600 dark:text-gray-400">{{ t('tenants.loading') }}</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-12">
        <UAlert
          icon="i-heroicons-exclamation-triangle"
          color="error"
          variant="soft"
          :title="t('common.error')"
          :description="error"
          :close-button="{ icon: 'i-heroicons-x-mark-20-solid', color: 'neutral', variant: 'link', padded: false }"
          @close="clearError"
        />
        <div class="mt-6">
          <UButton @click="loadTenants" variant="outline">
            <UIcon name="i-heroicons-arrow-path" class="w-4 h-4 mr-2" />
            {{ t('common.refresh') }}
          </UButton>
        </div>
      </div>

      <!-- Tenant Grid -->
      <div v-else-if="tenants.length > 0" class="space-y-8">
        <!-- Action Bar -->
        <div class="flex justify-between items-center">
          <div class="flex items-center space-x-4">
            <h2 class="text-xl font-semibold text-gray-900 dark:text-gray-100">
              Your Servers ({{ tenants.length }})
            </h2>
            <p class="text-sm text-gray-500">Debug: availableGuildsForTenant.length = {{ availableGuildsForTenant.length }}</p>
          </div>
          <AddTenantModal v-if="!isLoading" :available-guilds="availableGuildsForTenant" @refresh="loadAllData" :key="`modal-${availableGuildsForTenant.length}`">
            <UButton>
              <UIcon name="i-heroicons-plus" class="w-4 h-4 mr-2" />
              {{ t('tenants.addServer') }}
            </UButton>
          </AddTenantModal>
        </div>

        <!-- Tenant Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <TenantCard
            v-for="tenant in tenants"
            :key="tenant.id"
            :tenant="tenant"
            :is-owner="isOwner(tenant)"
            @select="selectTenant"
          >
            <template #delete-button>
              <DeleteTenantModal :tenant="tenant">
                <UButton
                  color="error"
                  variant="ghost"
                  size="sm"
                  icon="i-heroicons-trash"
                />
              </DeleteTenantModal>
            </template>
          </TenantCard>
        </div>
      </div>

      <!-- Empty State -->
      <EmptyTenantState v-else>
        <template #add-server-button>
          <AddTenantModal v-if="!isLoading" :available-guilds="availableGuildsForTenant" @refresh="loadAllData" :key="`modal-${availableGuildsForTenant.length}`">
            <UButton>
              <UIcon name="i-heroicons-plus" class="w-4 h-4 mr-2" />
              {{ t('tenants.addServer') }}
            </UButton>
          </AddTenantModal>
        </template>
      </EmptyTenantState>
    </div>


  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const { t } = useI18n()
const { user } = useAuth()
const { 
  tenants, 
  availableGuilds,
  isLoading, 
  error,
  fetchUserTenants,
  fetchAvailableGuilds,
  switchTenant,
  clearError
} = useTenant()

// Computed property to get guilds that are not already tenants
const availableGuildsForTenant = computed(() => {
  console.log('Debug: availableGuildsForTenant computed running', {
    availableGuildsExists: !!availableGuilds.value,
    availableGuildsLength: availableGuilds.value?.length,
    tenantsExists: !!tenants.value,
    tenantsLength: tenants.value?.length
  })
  
  if (!availableGuilds.value || !tenants.value) {
    console.log('Debug: Missing data', { 
      availableGuilds: availableGuilds.value?.length, 
      tenants: tenants.value?.length 
    })
    return []
  }
  
  const tenantGuildIds = new Set(tenants.value.map(tenant => tenant.discord_server_id))
  const filtered = availableGuilds.value.filter(guild => !tenantGuildIds.has(guild.id))
  
  console.log('Debug: Filtering guilds', {
    totalGuilds: availableGuilds.value.length,
    existingTenants: tenants.value.length,
    tenantGuildIds: Array.from(tenantGuildIds),
    availableForTenant: filtered.length,
    filteredGuilds: filtered.map(g => ({ id: g.id, name: g.name }))
  })
  
  return filtered
})

// Load data on mount
onMounted(async () => {
  await loadAllData()
})

// Watch the computed property to ensure it's reactive
watch(availableGuildsForTenant, (newValue) => {
  console.log('Debug: availableGuildsForTenant changed', {
    length: newValue?.length,
    guilds: newValue?.map(g => ({ id: g.id, name: g.name }))
  })
}, { immediate: true })

// Methods
const loadAllData = async () => {
  try {
    clearError()
    console.log('Debug: Loading all data...')
    // Load both tenants and available guilds in parallel to avoid rate limiting
    await Promise.all([
      fetchUserTenants(),
      fetchAvailableGuilds()
    ])
    console.log('Debug: Data loaded', {
      tenantsCount: tenants.value?.length,
      guildsCount: availableGuilds.value?.length
    })
  } catch (error) {
    console.error('Failed to load data:', error)
  }
}

const loadTenants = async () => {
  try {
    clearError()
    await fetchUserTenants()
  } catch (error) {
    console.error('Failed to load tenants:', error)
  }
}

const selectTenant = async (tenant: any) => {
  try {
    await switchTenant(tenant)
  } catch (error) {
    console.error('Failed to switch tenant:', error)
    const toast = useToast()
    toast.add({
      title: 'Failed to Switch Server',
      description: 'There was an error switching to the selected server',
      color: 'error'
    })
  }
}





// Helper functions
const isOwner = (tenant: any) => {
  return tenant.owner_id === user.value?.id
}
</script>

