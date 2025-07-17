<script setup lang="ts">
const { handleTokensFromUrl, error: authError, clearError } = useAuth()
const route = useRoute()
const router = useRouter()

definePageMeta({
  layout: 'login'
})

const localError = ref<string | null>(null)
const isProcessing = ref(true)

// Computed error that combines local and auth errors
const error = computed(() => localError.value || authError.value)

onMounted(async () => {
  try {
    clearError()
    
    // Check for OAuth error from Discord
    const errorParam = route.query.error as string
    if (errorParam) {
      throw new Error(`OAuth error: ${errorParam}`)
    }

    // Get tokens from query parameters (sent by backend)
    const accessToken = route.query.access_token as string
    const refreshToken = route.query.refresh_token as string
    const expiresIn = route.query.expires_in as string

    if (!accessToken || !refreshToken) {
      throw new Error('Missing authentication tokens')
    }

    // Handle tokens and get user info
    await handleTokensFromUrl(accessToken, refreshToken, parseInt(expiresIn) || 3600)

    // Redirect to callback URL or dashboard
    const callbackUrl = import.meta.client ? localStorage.getItem('auth_callback_url') : null
    if (import.meta.client) {
      localStorage.removeItem('auth_callback_url')
    }

    await router.push(callbackUrl || '/dashboard')
  } catch (err: any) {
    console.error('OAuth callback error:', err)
    localError.value = err.message || 'Authentication failed'
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