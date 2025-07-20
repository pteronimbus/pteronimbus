<script setup lang="ts">
import { useToast } from '#imports'
const route = useRoute()
const router = useRouter()

onMounted(() => {
  // Show a toast if error_description is present in the query
  const errorDescription = route.query.error_description as string
  if (errorDescription) {
    const toast = useToast()
    toast.add({
      title: 'Login Error',
      description: errorDescription,
      color: 'error',
    })
    // Remove error_description from the URL after showing the toast
    router.replace({ query: { ...route.query, error_description: undefined } })
  }
})

const { signIn, isLoading, error, clearError } = useAuth()

definePageMeta({
  layout: 'login',
  middleware: 'guest'
})

const loginWithDiscord = async () => {
  try {
    clearError()
    await signIn('discord', { callbackUrl: '/dashboard' })
  } catch (error) {
    console.error('Login failed:', error)
    // Error is already handled in the composable
  }
}
</script>

<template>
  <div class="flex items-center justify-center">
    <UCard class="w-96 max-w-md">
      <template #header>
        <h2 class="text-xl font-bold text-center text-gray-900 dark:text-gray-100">Login</h2>
      </template>

      <div class="space-y-4">
        <div v-if="error" class="p-3 text-sm text-red-600 bg-red-50 border border-red-200 rounded-md dark:bg-red-900/20 dark:text-red-400 dark:border-red-800">
          {{ error }}
        </div>
        
        <UButton @click="loginWithDiscord" block color="primary" :loading="isLoading">
          Login with Discord
        </UButton>
      </div>
    </UCard>
  </div>
</template> 