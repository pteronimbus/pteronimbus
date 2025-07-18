<template>
  <div class="text-center py-12">
    <UIcon name="heroicons:server" class="w-16 h-16 text-gray-400 dark:text-gray-500 mx-auto mb-6" />
    <h3 class="text-xl font-medium text-gray-900 dark:text-gray-100 mb-4">
      No Discord Servers Found
    </h3>
    <p class="text-gray-600 dark:text-gray-400 mb-8 max-w-md mx-auto">
      You haven't added any Discord servers to Pteronimbus yet. 
      Add your first server to start managing game servers for your community.
    </p>
    <AddTenantModal :available-guilds="availableGuilds" @refresh="$emit('refresh')">
      <UButton
        size="lg"
        color="primary"
        variant="solid"
        class="shadow-lg bg-primary-500 hover:bg-primary-600 text-white font-medium"
      >
        <UIcon name="heroicons:plus" class="w-5 h-5 mr-2" />
        Add Your First Server
      </UButton>
    </AddTenantModal>
  </div>
</template>

<script setup lang="ts">
interface Props {
  availableGuilds?: any[]
}

interface Emits {
  (e: 'refresh'): void
}

const props = withDefaults(defineProps<Props>(), {
  availableGuilds: () => []
})

defineEmits<Emits>()
</script>