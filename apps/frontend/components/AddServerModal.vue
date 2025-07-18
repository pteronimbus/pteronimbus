<template>
  <UModal v-model="isOpen" prevent-close>
    <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
            {{ t('servers.modals.add.title') }}
          </h3>
          <UButton
            color="gray"
            variant="ghost"
            icon="i-heroicons-x-mark-20-solid"
            class="-my-1"
            @click="closeModal"
          />
        </div>
      </template>

      <div class="space-y-6">
        <!-- Server Name -->
        <UFormGroup :label="t('servers.modals.add.fields.name')" required>
          <UInput
            v-model="form.name"
            :placeholder="t('servers.modals.add.placeholders.name')"
            :error="errors.name"
          />
        </UFormGroup>

        <!-- Game Selection -->
        <UFormGroup :label="t('servers.modals.add.fields.game')" required>
          <USelectMenu
            v-model="form.game"
            :options="gameOptions"
            :placeholder="t('servers.modals.add.placeholders.game')"
            :error="errors.game"
          />
        </UFormGroup>

        <!-- Server Configuration -->
        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="t('servers.modals.add.fields.maxPlayers')">
            <UInput
              v-model="form.maxPlayers"
              type="number"
              min="1"
              max="100"
              :placeholder="t('servers.modals.add.placeholders.maxPlayers')"
            />
          </UFormGroup>

          <UFormGroup :label="t('servers.modals.add.fields.port')">
            <UInput
              v-model="form.port"
              type="number"
              min="1024"
              max="65535"
              :placeholder="t('servers.modals.add.placeholders.port')"
            />
          </UFormGroup>
        </div>

        <!-- Description -->
        <UFormGroup :label="t('servers.modals.add.fields.description')">
          <UTextarea
            v-model="form.description"
            :placeholder="t('servers.modals.add.placeholders.description')"
            :rows="3"
          />
        </UFormGroup>

        <!-- Auto Start -->
        <UFormGroup>
          <UCheckbox
            v-model="form.autoStart"
            :label="t('servers.modals.add.fields.autoStart')"
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
            :loading="isCreating"
            @click="createServer"
          >
            {{ t('servers.modals.add.create') }}
          </UButton>
        </div>
      </template>
    </UCard>
  </UModal>
</template>

<script setup lang="ts">
interface ServerForm {
  name: string
  game: string
  maxPlayers: number
  port: number
  description: string
  autoStart: boolean
}

interface Props {
  modelValue: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'created', server: ServerForm): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const { t } = useI18n()

const isOpen = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const form = ref<ServerForm>({
  name: '',
  game: '',
  maxPlayers: 20,
  port: 25565,
  description: '',
  autoStart: false
})

const errors = ref<Record<string, string>>({})
const error = ref<string>('')
const isCreating = ref(false)

const gameOptions = [
  { label: 'Minecraft', value: 'minecraft' },
  { label: 'Valheim', value: 'valheim' },
  { label: 'CS:GO', value: 'csgo' },
  { label: 'Terraria', value: 'terraria' },
  { label: 'Rust', value: 'rust' },
  { label: 'ARK: Survival Evolved', value: 'ark' },
  { label: 'Palworld', value: 'palworld' },
  { label: 'Satisfactory', value: 'satisfactory' }
]

const validateForm = (): boolean => {
  errors.value = {}
  
  if (!form.value.name.trim()) {
    errors.value.name = t('servers.modals.add.validation.nameRequired')
    return false
  }
  
  if (!form.value.game) {
    errors.value.game = t('servers.modals.add.validation.gameRequired')
    return false
  }
  
  return true
}

const createServer = async () => {
  if (!validateForm()) return
  
  isCreating.value = true
  error.value = ''
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    // Emit the created server data
    emit('created', { ...form.value })
    
    // Reset form and close modal
    resetForm()
    closeModal()
  } catch (err: any) {
    error.value = err.message || t('servers.modals.add.errors.createFailed')
  } finally {
    isCreating.value = false
  }
}

const resetForm = () => {
  form.value = {
    name: '',
    game: '',
    maxPlayers: 20,
    port: 25565,
    description: '',
    autoStart: false
  }
  errors.value = {}
  error.value = ''
}

const closeModal = () => {
  emit('update:modelValue', false)
}

const clearError = () => {
  error.value = ''
}

// Reset form when modal opens
watch(isOpen, (newValue) => {
  if (newValue) {
    resetForm()
  }
})
</script>