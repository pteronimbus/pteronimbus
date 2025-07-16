<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-4xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="text-center mb-12">
        <h1 class="text-3xl font-bold text-gray-900 mb-4">
          Select Discord Server
        </h1>
        <p class="text-lg text-gray-600 max-w-2xl mx-auto">
          Choose a Discord server to manage game servers for your community.
          You can switch between servers anytime.
        </p>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="text-center py-12">
        <UIcon name="heroicons:arrow-path" class="w-8 h-8 animate-spin mx-auto mb-4 text-primary-500" />
        <p class="text-gray-600">Loading your Discord servers...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-12">
        <UIcon name="heroicons:exclamation-triangle" class="w-12 h-12 text-red-500 mx-auto mb-4" />
        <h3 class="text-lg font-medium text-gray-900 mb-2">Failed to Load Servers</h3>
        <p class="text-gray-600 mb-6">{{ error }}</p>
        <UButton @click="loadTenants" variant="outline">
          <UIcon name="heroicons:arrow-path" class="w-4 h-4 mr-2" />
          Try Again
        </UButton>
      </div>

      <!-- Tenant Grid -->
      <div v-else-if="tenants.length > 0" class="space-y-8">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="tenant in tenants"
            :key="tenant.id"
            class="tenant-card"
          >
            <UCard class="h-full hover:shadow-lg transition-shadow cursor-pointer" @click="selectTenant(tenant)">
              <div class="flex items-center space-x-4 mb-4">
                <UAvatar
                  :src="getTenantIcon(tenant)"
                  :alt="tenant.name"
                  size="lg"
                />
                <div class="flex-1 min-w-0">
                  <h3 class="text-lg font-medium text-gray-900 truncate">
                    {{ tenant.name }}
                  </h3>
                  <p class="text-sm text-gray-500">
                    {{ isOwner(tenant) ? 'Owner' : 'Member' }}
                  </p>
                </div>
              </div>
              
              <div class="space-y-2 text-sm text-gray-600">
                <div class="flex items-center">
                  <UIcon name="heroicons:calendar" class="w-4 h-4 mr-2" />
                  Added {{ formatDate(tenant.created_at) }}
                </div>
                <div class="flex items-center">
                  <UIcon name="heroicons:server" class="w-4 h-4 mr-2" />
                  {{ tenant.config?.resource_limits?.max_game_servers || 5 }} server limit
                </div>
              </div>

              <div class="mt-6 flex space-x-2">
                <UButton
                  class="flex-1"
                  @click.stop="selectTenant(tenant)"
                >
                  Select Server
                </UButton>
                <UDropdown
                  v-if="isOwner(tenant)"
                  :items="getTenantActions(tenant)"
                  :popper="{ placement: 'bottom-end' }"
                >
                  <UButton
                    variant="ghost"
                    size="sm"
                    icon="heroicons:ellipsis-vertical"
                    @click.stop
                  />
                </UDropdown>
              </div>
            </UCard>
          </div>
        </div>

        <!-- Add New Server Button -->
        <div class="text-center">
          <UButton
            size="lg"
            variant="outline"
            @click="showCreateModal = true"
          >
            <UIcon name="heroicons:plus" class="w-5 h-5 mr-2" />
            Add New Discord Server
          </UButton>
        </div>
      </div>

      <!-- No Tenants State -->
      <div v-else class="text-center py-12">
        <UIcon name="heroicons:server" class="w-16 h-16 text-gray-400 mx-auto mb-6" />
        <h3 class="text-xl font-medium text-gray-900 mb-4">
          No Discord Servers Found
        </h3>
        <p class="text-gray-600 mb-8 max-w-md mx-auto">
          You haven't added any Discord servers to Pteronimbus yet. 
          Add your first server to start managing game servers for your community.
        </p>
        <UButton
          size="lg"
          @click="showCreateModal = true"
        >
          <UIcon name="heroicons:plus" class="w-5 h-5 mr-2" />
          Add Your First Server
        </UButton>
      </div>
    </div>

    <!-- Create Tenant Modal -->
    <UModal v-model="showCreateModal" :ui="{ width: 'sm:max-w-lg' }">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Add Discord Server</h3>
            <UButton
              variant="ghost"
              size="sm"
              icon="heroicons:x-mark"
              @click="showCreateModal = false"
            />
          </div>
        </template>

        <div class="space-y-6">
          <div class="text-sm text-gray-600">
            <p class="mb-4">
              Select a Discord server where you have "Manage Server" permissions. 
              This will allow Pteronimbus to integrate with your Discord server for game server management.
            </p>
            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <div class="flex items-start">
                <UIcon name="heroicons:information-circle" class="w-5 h-5 text-blue-500 mr-2 mt-0.5" />
                <div class="text-blue-700 text-sm">
                  <p class="font-medium mb-1">What happens when you add a server?</p>
                  <ul class="list-disc list-inside space-y-1 text-xs">
                    <li>Your Discord roles will be synced for permission management</li>
                    <li>You can manage game servers through Discord commands</li>
                    <li>Server notifications will be sent to Discord channels</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>

          <!-- Loading Available Guilds -->
          <div v-if="loadingGuilds" class="text-center py-8">
            <UIcon name="heroicons:arrow-path" class="w-6 h-6 animate-spin mx-auto mb-2" />
            <p class="text-sm text-gray-500">Loading your Discord servers...</p>
          </div>

          <!-- Available Guilds -->
          <div v-else-if="availableGuilds.length > 0" class="space-y-3">
            <div
              v-for="guild in availableGuilds"
              :key="guild.id"
              class="guild-item border border-gray-200 rounded-lg p-4 hover:border-gray-300 hover:bg-gray-50 transition-colors cursor-pointer"
              @click="createTenantFromGuild(guild)"
            >
              <div class="flex items-center space-x-3">
                <UAvatar
                  :src="getGuildIcon(guild)"
                  :alt="guild.name"
                  size="md"
                />
                <div class="flex-1 min-w-0">
                  <h4 class="font-medium text-gray-900 truncate">{{ guild.name }}</h4>
                  <p class="text-sm text-gray-500">
                    {{ guild.owner ? 'Owner' : 'Manager' }}
                  </p>
                </div>
                <UButton
                  size="sm"
                  :loading="creatingTenant === guild.id"
                  @click.stop="createTenantFromGuild(guild)"
                >
                  Add Server
                </UButton>
              </div>
            </div>
          </div>

          <!-- No Available Guilds -->
          <div v-else class="text-center py-8">
            <UIcon name="heroicons:shield-exclamation" class="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h4 class="text-lg font-medium mb-2">No Available Servers</h4>
            <p class="text-sm text-gray-500 mb-4">
              You need "Manage Server" permissions to add a Discord server to Pteronimbus.
              Make sure you're an administrator or have the required permissions.
            </p>
            <UButton variant="outline" @click="loadAvailableGuilds">
              <UIcon name="heroicons:arrow-path" class="w-4 h-4 mr-2" />
              Refresh List
            </UButton>
          </div>
        </div>
      </UCard>
    </UModal>

    <!-- Delete Confirmation Modal -->
    <UModal v-model="showDeleteModal" :ui="{ width: 'sm:max-w-md' }">
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold text-red-600">Delete Server</h3>
        </template>

        <div class="space-y-4">
          <div class="flex items-start space-x-3">
            <UIcon name="heroicons:exclamation-triangle" class="w-6 h-6 text-red-500 mt-1" />
            <div>
              <p class="text-gray-900 font-medium mb-2">
                Are you sure you want to remove "{{ tenantToDelete?.name }}"?
              </p>
              <p class="text-sm text-gray-600">
                This will permanently delete all game servers, configurations, and data 
                associated with this Discord server. This action cannot be undone.
              </p>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="flex justify-end space-x-3">
            <UButton
              variant="ghost"
              @click="showDeleteModal = false"
            >
              Cancel
            </UButton>
            <UButton
              color="red"
              :loading="deletingTenant"
              @click="confirmDeleteTenant"
            >
              Delete Server
            </UButton>
          </div>
        </template>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const { user } = useAuth()
const { 
  tenants, 
  availableGuilds, 
  isLoading, 
  error,
  fetchUserTenants,
  fetchAvailableGuilds,
  createTenant,
  deleteTenant,
  switchTenant,
  clearError
} = useTenant()

// Local state
const showCreateModal = ref(false)
const showDeleteModal = ref(false)
const loadingGuilds = ref(false)
const creatingTenant = ref<string | null>(null)
const deletingTenant = ref(false)
const tenantToDelete = ref<any>(null)

// Load data on mount
onMounted(async () => {
  await loadTenants()
})

// Methods
const loadTenants = async () => {
  try {
    clearError()
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
  } catch (error) {
    console.error('Failed to switch tenant:', error)
    const toast = useToast()
    toast.add({
      title: 'Failed to Switch Server',
      description: 'There was an error switching to the selected server',
      color: 'red'
    })
  }
}

const createTenantFromGuild = async (guild: any) => {
  creatingTenant.value = guild.id
  try {
    const newTenant = await createTenant(guild.id)
    showCreateModal.value = false
    
    // Show success notification
    const toast = useToast()
    toast.add({
      title: 'Server Added Successfully',
      description: `${guild.name} has been added to Pteronimbus`,
      color: 'green'
    })

    // Automatically switch to the new tenant
    await switchTenant(newTenant)
  } catch (error) {
    console.error('Failed to create tenant:', error)
    const toast = useToast()
    toast.add({
      title: 'Failed to Add Server',
      description: 'There was an error adding the Discord server',
      color: 'red'
    })
  } finally {
    creatingTenant.value = null
  }
}

const confirmDeleteTenant = async () => {
  if (!tenantToDelete.value) return

  deletingTenant.value = true
  try {
    await deleteTenant(tenantToDelete.value.id)
    showDeleteModal.value = false
    tenantToDelete.value = null

    const toast = useToast()
    toast.add({
      title: 'Server Removed',
      description: 'The Discord server has been removed from Pteronimbus',
      color: 'green'
    })
  } catch (error) {
    console.error('Failed to delete tenant:', error)
    const toast = useToast()
    toast.add({
      title: 'Failed to Remove Server',
      description: 'There was an error removing the Discord server',
      color: 'red'
    })
  } finally {
    deletingTenant.value = false
  }
}

// Helper functions
const isOwner = (tenant: any) => {
  return tenant.owner_id === user.value?.id
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

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const getTenantActions = (tenant: any) => {
  return [
    [{
      label: 'Manage Settings',
      icon: 'heroicons:cog-6-tooth',
      click: () => navigateTo(`/tenant/${tenant.id}/settings`)
    }],
    [{
      label: 'Remove Server',
      icon: 'heroicons:trash',
      click: () => {
        tenantToDelete.value = tenant
        showDeleteModal.value = true
      }
    }]
  ]
}

// Watch for create modal opening
watch(showCreateModal, (newValue) => {
  if (newValue && availableGuilds.value.length === 0) {
    loadAvailableGuilds()
  }
})
</script>

<style scoped>
.tenant-card {
  @apply transition-transform hover:scale-105;
}

.guild-item {
  @apply transition-all duration-200;
}


</style>