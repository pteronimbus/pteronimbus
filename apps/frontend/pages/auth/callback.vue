<script setup lang="ts">
const { handleCallback } = useAuth()
const route = useRoute()

definePageMeta({
  layout: 'login'
})

const error = ref<string | null>(null)
const isProcessing = ref(true)

onMounted(async () => {
  try {
    const code = route.query.code as string
    const state = route.query.state as string
    const errorParam = route.query.error as string

    if (errorParam) {
      throw new Error(`OAuth error: ${errorParam}`)
    }

    if (!code || !state) {
      throw new Error('Missing authorization code or state parameter')
    }

    await handleCallback(code, state)
  } catch (err: any) {
    console.error('OAuth callback error:', err)
    error.value = err.message || 'Authentication failed'
    isProcessing.value = false
  }
})
</script>

<template>
  <div class="flex items-center justify-center">
    <UCard class="w-96 max-w-md">
      <template #header>
        <h2 class="text-xl font-bold text-center text-gray-900 dark:text-gray-100">
          {{ error ? 'Authentication Failed' : 'Authenticating...' }}
        </h2>
      </template>

      <div class="space-y-4 text-center">
        <div v-if="isProcessing && !error" class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
          <p class="text-gray-600 dark:text-gray-400">
            Processing your Discord authentication...
          </p>
        </div>

        <div v-else-if="error" class="space-y-4">
          <div class="text-red-500">
            <Icon name="heroicons:exclamation-triangle" class="h-8 w-8 mx-auto mb-2" />
            <p>{{ error }}</p>
          </div>
          <UButton @click="$router.push('/login')" block color="primary">
            Try Again
          </UButton>
        </div>
      </div>
    </UCard>
  </div>
</template>