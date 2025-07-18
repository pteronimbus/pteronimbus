<template>
  <UModal :title="t('tenants.modals.add.title')" :ui="{
    overlay: 'fixed inset-0 bg-gray-200/75 dark:bg-gray-900/75 backdrop-blur-sm'
  }">
    <slot />

    <template #body>
      <div v-if="step === 'select-guild'" class="space-y-6">
        <div class="text-sm text-gray-600 dark:text-gray-400">
          <p class="mb-4">
            {{ t('tenants.modals.add.description') }}
          </p>
          <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
            <div class="flex items-start">
              <UIcon name="i-heroicons-information-circle" class="w-5 h-5 text-blue-500 mr-2 mt-0.5" />
              <div class="text-blue-700 dark:text-blue-300 text-sm">
                <p class="font-medium mb-1">{{ t('tenants.modals.add.infoTitle') }}</p>
                <ul class="list-disc list-inside space-y-1 text-xs">
                  <li v-for="item in infoItems" :key="item">{{ item }}</li>
                </ul>
              </div>
            </div>
          </div>
        </div>

        <div v-if="isLoading" class="text-center py-8">
          <UIcon name="i-heroicons-arrow-path" class="w-6 h-6 animate-spin mx-auto mb-2 text-primary-500" />
          <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('tenants.modals.add.loadingServers') }}</p>
        </div>

        <div v-else-if="availableGuilds.length > 0" class="space-y-3">
          <h4 class="font-medium text-gray-900 dark:text-gray-100 mb-3">{{ t('tenants.modals.add.availableServers') }}
          </h4>
          <div v-for="guild in availableGuilds" :key="guild.id"
            class="guild-item border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:border-gray-300 dark:hover:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors cursor-pointer"
            @click="createTenantFromGuild(guild)">
            <div class="flex items-center space-x-3">
              <UAvatar :src="getGuildIcon(guild)" :alt="guild.name" size="md" />
              <div class="flex-1 min-w-0">
                <h4 class="font-medium text-gray-900 dark:text-gray-100 truncate">{{ guild.name }}</h4>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ guild.owner ? t('tenants.modals.add.owner') : t('tenants.modals.add.manager') }}
                </p>
              </div>
              <UButton size="sm" :loading="creatingTenant === guild.id" @click.stop="createTenantFromGuild(guild)">
                {{ t('tenants.modals.add.addButton') }}
              </UButton>
            </div>
          </div>
        </div>

        
        <div v-else class="text-center py-8">
          <UIcon name="i-heroicons-shield-exclamation"
            class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" />
          <h4 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">{{
            t('tenants.modals.add.noAvailableServers')
          }}</h4>
          <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">
            {{ t('tenants.modals.add.noServersDescription') }}
          </p>
          <UButton variant="outline" @click="loadAvailableGuilds">
            <UIcon name="i-heroicons-arrow-path" class="w-4 h-4 mr-2" />
            {{ t('tenants.modals.add.refreshList') }}
          </UButton>
        </div>

        
        <UAlert v-if="error" icon="i-heroicons-exclamation-triangle" color="error" variant="soft"
          :title="t('common.error')" :description="error"
          :close-button="{ icon: 'i-heroicons-x-mark-20-solid', color: 'neutral', variant: 'link', padded: false }"
          @close="clearError" />
      </div>
      <div v-if="step === 'install-bot'" class="text-center py-8">
        <UIcon name="i-heroicons-puzzle-piece" class="w-12 h-12 text-primary-500 mx-auto mb-4" />
        <h4 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
          {{ t('tenants.modals.add.installBotTitle') }}
        </h4>
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-6">
          {{ t('tenants.modals.add.installBotDescription') }}
        </p>
        <UButton size="lg" @click="installBot">
          <UIcon name="i-heroicons-arrow-top-right-on-square" class="w-5 h-5 mr-2" />
          {{ t('tenants.modals.add.installBotButton') }}
        </UButton>
        <div class="mt-8">
          <UButton variant="ghost" @click="finish">
            {{ t('tenants.modals.add.finishButton') }}
          </UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
interface Props {
  availableGuilds?: any[]
}

interface Emits {
  (e: 'refresh'): void
  (e: 'finish', tenant: any): void
}

const props = withDefaults(defineProps<Props>(), {
  availableGuilds: () => []
})
const emit = defineEmits<Emits>()

const { t } = useI18n()
const toast = useToast()
const {
  fetchAvailableGuilds,
  createTenant,
  switchTenant
} = useTenant()

// Internal state
const isLoading = ref(false)
const creatingTenant = ref<string | null>(null)
const error = ref<string | null>(null)
const step = ref('select-guild') // 'select-guild' or 'install-bot'
const newTenant = ref<any>(null)

// Use the guilds passed from parent
const availableGuilds = computed(() => {
  console.log('Debug: AddTenantModal availableGuilds', {
    propsLength: props.availableGuilds?.length,
    guilds: props.availableGuilds?.map(g => ({ id: g.id, name: g.name }))
  })
  return props.availableGuilds
})

// Watch for changes in the props to ensure reactivity
watch(() => props.availableGuilds, (newGuilds) => {
  console.log('Debug: AddTenantModal props changed', {
    length: newGuilds?.length,
    guilds: newGuilds?.map(g => ({ id: g.id, name: g.name }))
  })
}, { immediate: true, deep: true })

const loadAvailableGuilds = async () => {
  isLoading.value = true
  error.value = null
  try {
    // Emit to parent to refresh the guild data
    emit('refresh')
    await fetchAvailableGuilds()
  } catch (err: any) {
    error.value = err.message || 'Failed to load available guilds'
  } finally {
    isLoading.value = false
  }
}

const createTenantFromGuild = async (guild: any) => {
  creatingTenant.value = guild.id
  error.value = null
  try {
    const createdTenant = await createTenant(guild.id)
    newTenant.value = createdTenant
    step.value = 'install-bot'

    // Show success notification for tenant creation
    toast.add({
      title: t('tenants.modals.add.successTitle'),
      description: t('tenants.modals.add.successDescription', { serverName: guild.name }),
      color: 'success'
    })
  } catch (err: any) {
    error.value = err.message || 'Failed to add Discord server'
  } finally {
    creatingTenant.value = null
  }
}

const installBot = () => {
  const config = useRuntimeConfig()
  const clientId = config.public.discordClientId
  if (!clientId || !newTenant.value) {
    error.value = 'Configuration error: Missing Discord Client ID.'
    return
  }
  const permissions = '8' // Administrator
  const url = `https://discord.com/api/oauth2/authorize?client_id=${clientId}&guild_id=${newTenant.value.discord_server_id}&permissions=${permissions}&scope=bot%20applications.commands`
  window.open(url, '_blank')
}

const finish = async () => {
  if (newTenant.value) {
    emit('finish', newTenant.value)
  }
}

const clearError = () => {
  error.value = null
}



const getGuildIcon = (guild: any) => {
  if (guild.icon) {
    return `https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png`
  }
  return undefined
}

// Computed property to properly handle the i18n array
const infoItems = computed(() => {
  // Get the raw translation data to access the array directly
  const { tm } = useI18n()
  const items = tm('tenants.modals.add.infoItems')

  // If tm returns an array, use it; otherwise use fallback
  if (Array.isArray(items)) {
    return items
  }

  // Fallback to hardcoded values if i18n array doesn't work
  return [
    t('tenants.modals.add.infoItems.0') || 'Your Discord roles will be synced for permission management',
    t('tenants.modals.add.infoItems.1') || 'You can manage game servers through Discord commands',
    t('tenants.modals.add.infoItems.2') || 'Server notifications will be sent to Discord channels'
  ]
})

// No need to load guilds on mount - they're provided by parent
</script>

<style scoped>
.guild-item {
  transition: all 0.2s;
}
</style>