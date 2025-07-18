<template>
  <UModal prevent-close>
    <slot />
    
    <template #header>
      <div class="flex items-center gap-3">
        <UIcon name="i-heroicons-exclamation-triangle" class="w-6 h-6 text-red-500" />
        <h3 class="text-lg font-semibold text-red-600 dark:text-red-400">
          {{ t('tenants.modals.delete.title') }}
        </h3>
      </div>
    </template>

    <template #body>
      <div class="space-y-4">
        <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
          <p class="text-gray-900 dark:text-gray-100 font-medium mb-2">
            {{ t('tenants.modals.delete.confirmMessage', { name: tenant?.name || '' }) }}
          </p>
          <p class="text-sm text-gray-600 dark:text-gray-400">
            {{ t('tenants.modals.delete.warningMessage') }}
          </p>
        </div>

        <!-- Tenant Details -->
        <div v-if="tenant" class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
          <div class="flex items-center gap-3">
            <UAvatar
              :src="getTenantIcon(tenant)"
              :alt="tenant.name"
              size="md"
            />
            <div>
              <h4 class="font-medium text-gray-900 dark:text-gray-100">{{ tenant.name }}</h4>
              <p class="text-sm text-gray-500 dark:text-gray-400">Discord Server</p>
            </div>
          </div>
        </div>

        <!-- Error Display -->
        <UAlert
          v-if="error"
          icon="i-heroicons-exclamation-triangle"
          color="error"
          variant="soft"
          :title="t('common.error')"
          :description="error"
          :close-button="{ icon: 'i-heroicons-x-mark-20-solid', color: 'neutral', variant: 'link', padded: false }"
          @close="clearError"
        />
      </div>
    </template>

    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton
          color="neutral"
          variant="ghost"
          @click="closeModal"
        >
          {{ t('common.cancel') }}
        </UButton>
        <UButton
          color="error"
          :loading="isDeleting"
          @click="confirmDelete"
        >
          {{ t('tenants.modals.delete.confirmButton') }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
interface Props {
  tenant: any | null
}

const props = defineProps<Props>()

const { t } = useI18n()
const toast = useToast()
const { deleteTenant } = useTenant()

// Internal state
const isDeleting = ref(false)
const error = ref<string | null>(null)

const closeModal = () => {
  // The modal will close automatically when the user clicks outside or presses escape
  // since we're not using prevent-close for this action
}

const confirmDelete = async () => {
  if (!props.tenant) return

  isDeleting.value = true
  error.value = null
  try {
    await deleteTenant(props.tenant.id)
    
    toast.add({
      title: 'Server Removed',
      description: 'The Discord server has been removed from Pteronimbus',
      color: 'success'
    })
    
    // The modal will close automatically after successful deletion
    // since the tenant list will be updated and the modal trigger will be removed
  } catch (err: any) {
    console.error('Failed to delete tenant:', err)
    error.value = err.message || t('tenants.modals.delete.errors.deleteFailed')
  } finally {
    isDeleting.value = false
  }
}

const clearError = () => {
  error.value = null
}

const getTenantIcon = (tenant: any) => {
  if (tenant.icon) {
    return `https://cdn.discordapp.com/icons/${tenant.discord_server_id}/${tenant.icon}.png`
  }
  return undefined
}
</script>