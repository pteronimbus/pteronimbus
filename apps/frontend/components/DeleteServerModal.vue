<template>
  <UModal 
    v-model="isOpen" 
    :dismissible="false"
    :ui="{ 
      overlay: 'fixed inset-0 bg-gray-200/75 dark:bg-gray-900/75 backdrop-blur-sm' 
    }"
  >
    <slot />
    
    <template #header>
      <div class="flex items-center gap-3">
        <UIcon name="i-heroicons-exclamation-triangle" class="w-6 h-6 text-red-500" />
        <h3 class="text-lg font-semibold text-red-600 dark:text-red-400">
          {{ t('servers.modals.delete.title') }}
        </h3>
      </div>
    </template>

    <template #body>
      <div class="space-y-4">
        <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
          <p class="text-gray-900 dark:text-gray-100 font-medium mb-2">
            {{ t('servers.modals.delete.confirmMessage', { name: server?.name || '' }) }}
          </p>
          <p class="text-sm text-gray-600 dark:text-gray-400">
            {{ t('servers.modals.delete.warningMessage') }}
          </p>
        </div>

        <!-- Server Details -->
        <div v-if="server" class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
          <div class="flex items-center gap-3 mb-3">
            <div class="w-10 h-10 bg-blue-100 dark:bg-blue-900 rounded-lg flex items-center justify-center">
              <UIcon name="i-heroicons-server-20-solid" class="w-5 h-5 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <h4 class="font-medium text-gray-900 dark:text-gray-100">{{ server.name }}</h4>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ server.game }}</p>
            </div>
          </div>
          
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span class="text-gray-500 dark:text-gray-400">{{ t('servers.columns.status') }}:</span>
              <StatusBadge :status="server.status" type="server" class="ml-2" />
            </div>
            <div>
              <span class="text-gray-500 dark:text-gray-400">{{ t('servers.columns.players') }}:</span>
              <span class="ml-2 text-gray-900 dark:text-gray-100">{{ server.players }}</span>
            </div>
            <div>
              <span class="text-gray-500 dark:text-gray-400">{{ t('servers.columns.ip') }}:</span>
              <span class="ml-2 text-gray-900 dark:text-gray-100">{{ server.ip }}:{{ server.port }}</span>
            </div>
            <div>
              <span class="text-gray-500 dark:text-gray-400">{{ t('servers.columns.uptime') }}:</span>
              <span class="ml-2 text-gray-900 dark:text-gray-100">{{ server.uptime }}</span>
            </div>
          </div>
        </div>

        <!-- Confirmation Input -->
        <UFormGroup :label="t('servers.modals.delete.confirmationLabel', { name: server?.name || '' })">
          <UInput
            v-model="confirmationText"
            :placeholder="server?.name || ''"
            :error="confirmationError"
          />
        </UFormGroup>

        <!-- Error Display -->
        <UAlert
          v-if="error"
          icon="i-heroicons-exclamation-triangle"
          color="red"
          variant="soft"
          :title="t('common.error')"
          :description="error"
          :close-button="{ icon: 'i-heroicons-x-mark-20-solid', color: 'gray', variant: 'link', padded: false }"
          @close="clearError"
        />
      </div>
    </template>

    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton
          color="gray"
          variant="ghost"
          @click="closeModal"
        >
          {{ t('common.cancel') }}
        </UButton>
        <UButton
          color="red"
          :loading="isDeleting"
          :disabled="!canDelete"
          @click="deleteServer"
        >
          {{ t('servers.modals.delete.confirm') }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
interface Server {
  id: number
  name: string
  game: string
  status: string
  players: string
  ip: string
  port: number
  uptime: string
}

interface Props {
  modelValue?: boolean
  server?: Server | null
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'deleted', serverId: number): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  server: null
})
const emit = defineEmits<Emits>()

const { t } = useI18n()

// Internal state for when used with slot
const internalOpen = ref(false)

const isOpen = computed({
  get: () => props.modelValue !== undefined ? props.modelValue : internalOpen.value,
  set: (value) => {
    if (props.modelValue !== undefined) {
      emit('update:modelValue', value)
    } else {
      internalOpen.value = value
    }
  }
})

const confirmationText = ref('')
const confirmationError = ref('')
const error = ref('')
const isDeleting = ref(false)

const canDelete = computed(() => {
  return confirmationText.value === props.server?.name
})

const validateConfirmation = (): boolean => {
  confirmationError.value = ''
  
  if (!confirmationText.value) {
    confirmationError.value = t('servers.modals.delete.validation.confirmationRequired')
    return false
  }
  
  if (confirmationText.value !== props.server?.name) {
    confirmationError.value = t('servers.modals.delete.validation.confirmationMismatch')
    return false
  }
  
  return true
}

const deleteServer = async () => {
  if (!validateConfirmation() || !props.server) return
  
  isDeleting.value = true
  error.value = ''
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // Emit the deleted server ID
    emit('deleted', props.server.id)
    
    // Reset and close modal
    resetForm()
    closeModal()
  } catch (err: any) {
    error.value = err.message || t('servers.modals.delete.errors.deleteFailed')
  } finally {
    isDeleting.value = false
  }
}

const resetForm = () => {
  confirmationText.value = ''
  confirmationError.value = ''
  error.value = ''
}

const closeModal = () => {
  isOpen.value = false
}

const clearError = () => {
  error.value = ''
}

// Reset form when modal opens or server changes
watch([isOpen, () => props.server], ([newIsOpen, newServer]) => {
  if (newIsOpen && newServer) {
    resetForm()
  }
})
</script>