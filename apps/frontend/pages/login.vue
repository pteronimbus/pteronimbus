<script setup lang="ts">
const { signIn } = useAuth()

definePageMeta({
  layout: 'login',
  auth: {
    unauthenticatedOnly: true,
    navigateAuthenticatedTo: '/dashboard'
  }
})

const loginWithDiscord = async () => {
  await signIn('discord', { callbackUrl: '/dashboard' })
}
</script>

<template>
  <div class="flex items-center justify-center">
    <UCard class="w-96 max-w-md">
      <template #header>
        <h2 class="text-xl font-bold text-center text-gray-900 dark:text-gray-100">Login</h2>
      </template>

      <div class="space-y-4">
        <UButton @click="loginWithDiscord" block color="primary">
          Login with Discord
        </UButton>
      </div>
    </UCard>
  </div>
</template> 