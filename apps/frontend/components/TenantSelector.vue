<template>
  <div class="tenant-selector">
    <!-- Current Tenant Display -->
    <div v-if="currentTenant" class="current-tenant">
      <UButton variant="ghost" size="lg" class="w-full justify-between" @click="showSelector = !showSelector">
        <div class="flex items-center space-x-3">
          <UAvatar :src="getTenantIcon(currentTenant)" :alt="currentTenant.name" size="sm" />
          <div class="text-left">
            <div class="font-medium">{{ currentTenant.name }}</div>
            <div class="text-xs text-gray-500">Current Server</div>
          </div>
        </div>
        <UIcon name="lucide:chevron-down" class="w-4 h-4" />
      </UButton>
    </div>

    <!-- No Tenant Selected -->
    <div v-else class="no-tenant">
      <UButton variant="outline" size="lg" class="w-full" @click="showSelector = true">
        <UIcon name="heroicons:server" class="w-4 h-4 mr-2" />
        Select Server
      </UButton>
    </div>

    <!-- Tenant Selection Modal -->
    <UModal v-model="showSelector" :ui="{ width: 'sm:max-w-md' }">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Select Server</h3>
            <UButton variant="ghost" size="sm" icon="heroicons:x-mark" @click="showSelector = false" />
          </div>
        </template>

        <div class="space-y-4">
          <!-- Loading State -->
          <div v-if="isLoading" class="text-center py-8">
            <UIcon name="heroicons:arrow-path" class="w-6 h-6 animate-spin mx-auto mb-2" />
            <p class="text-sm text-gray-500">Loading servers...</p>
          </div>

          <!-- Error State -->
          <div v-else-if="error" class="text-center py-8">
            <UIcon name="heroicons:exclamation-triangle" class="w-6 h-6 text-red-500 mx-auto mb-2" />
            <p class="text-sm text-red-500 mb-4">{{ error }}</p>
            <UButton variant="outline" size="sm" @click="loadTenants">
              Try Again
            </UButton>
          </div>

          <!-- Tenant List -->
          <div v-else-if="tenants.length > 0" class="space-y-2">
            <div v-for="tenant in tenants" :key="tenant.id" class="tenant-item">
              <UButton variant="ghost" size="lg" class="w-full justify-start" :class="{
                'bg-primary-50 border-primary-200': currentTenant?.id === tenant.id
              }" @click="selectTenant(tenant)">
                <div class="flex items-center space-x-3">
                  <UAvatar :src="getTenantIcon(tenant)" :alt="tenant.name" size="sm" />
                  <div class="text-left">
                    <div class="font-medium">{{ tenant.name }}</div>
                    <div class="text-xs text-gray-500">
                      {{ tenant.id === tenant.owner_id ? 'Owner' : 'Member' }}
                    </div>
                  </div>
                </div>
                <UIcon v-if="currentTenant?.id === tenant.id" name="heroicons:check-circle"
                  class="w-4 h-4 text-primary-500 ml-auto" />
              </UButton>
            </div>
          </div>

          <!-- No Tenants -->
          <div v-else class="text-center py-8">
            <UIcon name="heroicons:server" class="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h4 class="text-lg font-medium mb-2">No Servers Found</h4>
            <p class="text-sm text-gray-500 mb-4">
              You don't have access to any servers yet. Create one to get started.
            </p>
            <UButton @click="showCreateTenant = true">
              <UIcon name="heroicons:plus" class="w-4 h-4 mr-2" />
              Add Server
            </UButton>
          </div>

          <!-- Create Tenant Button -->
          <div v-if="tenants.length > 0" class="border-t pt-4">
            <UButton variant="outline" class="w-full" @click="showCreateTenant = true">
              <UIcon name="heroicons:plus" class="w-4 h-4 mr-2" />
              Add New Server
            </UButton>
          </div>
        </div>
      </UCard>
    </UModal>

    <!-- Create Tenant Modal -->
    <UModal v-model="showCreateTenant" :ui="{ width: 'sm:max-w-md' }">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Add Discord Server</h3>
            <UButton variant="ghost" size="sm" icon="heroicons:x-mark" @click="showCreateTenant = false" />
          </div>
        </template>

        <div class="space-y-4">
          <!-- Loading Available Guilds -->
          <div v-if="loadingGuilds" class="text-center py-8">
            <UIcon name="heroicons:arrow-path" class="w-6 h-6 animate-spin mx-auto mb-2" />
            <p class="text-sm text-gray-500">Loading Discord servers...</p>
          </div>

          <!-- Available Guilds -->
          <div v-else-if="availableGuilds.length > 0" class="space-y-2">
            <p class="text-sm text-gray-600 mb-4">
              Select a Discord server where you have manage permissions:
            </p>
            <div v-for="guild in availableGuilds" :key="guild.id" class="guild-item">
              <UButton variant="ghost" size="lg" class="w-full justify-start" @click="createTenantFromGuild(guild)"
                :loading="creatingTenant">
                <div class="flex items-center space-x-3">
                  <UAvatar :src="getGuildIcon(guild)" :alt="guild.name" size="sm" />
                  <div class="text-left">
                    <div class="font-medium">{{ guild.name }}</div>
                    <div class="text-xs text-gray-500">
                      {{ guild.owner ? 'Owner' : 'Manager' }}
                    </div>
                  </div>
                </div>
              </UButton>
            </div>
          </div>

          <!-- No Available Guilds -->
          <div v-else class="text-center py-8">
            <UIcon name="heroicons:shield-exclamation" class="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h4 class="text-lg font-medium mb-2">No Available Servers</h4>
            <p class="text-sm text-gray-500 mb-4">
              You need "Manage Server" permissions to add a Discord server to Pteronimbus.
            </p>
            <UButton variant="outline" @click="loadAvailableGuilds">
              Refresh
            </UButton>
          </div>
        </div>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
interface Props {
  showLabel?: boolean
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  showLabel: true,
  size: 'md'
})

const {
  tenants,
  currentTenant,
  availableGuilds,
  isLoading,
  error,
  fetchUserTenants,
  fetchAvailableGuilds,
  createTenant,
  switchTenant,
  clearError
} = useTenant()

// Local state
const showSelector = ref(false)
const showCreateTenant = ref(false)
const loadingGuilds = ref(false)
const creatingTenant = ref(false)

// Load tenants on mount
onMounted(async () => {
  await loadTenants()
})

// Methods
const loadTenants = async () => {
  try {
    await fetchUserTenants()
  } catch (error) {
    console.error('Failed to load tenants:', error)
  }
}

const loadAvailableGuilds = async () => {
  loadingGuilds.value = true
  try {
    await fetchAvailableGuilds()
  } catch (error) {
    console.error('Failed to load available guilds:', error)
  } finally {
    loadingGuilds.value = false
  }
}

const selectTenant = async (tenant: any) => {
  try {
    await switchTenant(tenant)
    showSelector.value = false
  } catch (error) {
    console.error('Failed to switch tenant:', error)
  }
}

const createTenantFromGuild = async (guild: any) => {
  creatingTenant.value = true
  try {
    const newTenant = await createTenant(guild.id)
    showCreateTenant.value = false
    showSelector.value = false

    // Show success notification
    const toast = useToast()
    toast.add({
      title: 'Server Added',
      description: `${guild.name} has been added to Pteronimbus`,
      color: 'green'
    })
  } catch (error) {
    console.error('Failed to create tenant:', error)
    const toast = useToast()
    toast.add({
      title: 'Failed to Add Server',
      description: 'There was an error adding the Discord server',
      color: 'red'
    })
  } finally {
    creatingTenant.value = false
  }
}

const getTenantIcon = (tenant: any) => {
  if (tenant.icon) {
    return `https://cdn.discordapp.com/icons/${tenant.discord_server_id}/${tenant.icon}.png`
  }
  return null
}

const getGuildIcon = (guild: any) => {
  if (guild.icon) {
    return `https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`
  }
  return null
}

// Watch for create tenant modal opening
watch(showCreateTenant, (newValue) => {
  if (newValue && availableGuilds.value.length === 0) {
    loadAvailableGuilds()
  }
})
</script>

<style scoped>
@reference "~/assets/css/main.css";

.tenant-selector {
  @apply relative;
}

.current-tenant,
.no-tenant {
  @apply w-full;
}

.tenant-item,
.guild-item {
  @apply rounded-lg border border-gray-200 hover:border-gray-300 transition-colors;
}

.tenant-item:hover,
.guild-item:hover {
  @apply bg-gray-50;
}
</style>