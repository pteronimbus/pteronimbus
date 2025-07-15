<script setup lang="ts">
const { signIn, isLoading } = useAuth()

definePageMeta({
  layout: 'login',
  middleware: 'guest'
})

const loginWithDiscord = async () => {
  try {
    await signIn('discord', { callbackUrl: '/dashboard' })
  } catch (error) {
    console.error('Login failed:', error)
    // You could add a toast notification here
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
        <UButton @click="loginWithDiscord" block color="primary" :loading="isLoading">
          Login with Discord
        </UButton>
      </div>
    </UCard>
  </div>
</template> 