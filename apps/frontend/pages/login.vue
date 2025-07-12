<script setup lang="ts">
import { reactive } from 'vue'

definePageMeta({
  layout: 'login',
  auth: false
})

const router = useRouter()
const { saveUser } = useUser()

const state = reactive({
  email: '',
  password: ''
})

function login() {
  // Mock login with persistence
  saveUser({ email: state.email, name: 'Test User' })
  router.push('/dashboard')
}
</script>

<template>
  <div class="flex items-center justify-center">
    <UCard class="w-96 max-w-md">
      <template #header>
        <h2 class="text-xl font-bold text-center text-gray-900 dark:text-gray-100">Login</h2>
      </template>

      <UForm :state="state" @submit="login">
        <UFormField label="Email" name="email" class="mb-4">
          <UInput v-model="state.email" type="email" placeholder="you@example.com" icon="i-heroicons-envelope" class="w-full" />
        </UFormField>

        <UFormField label="Password" name="password" class="mb-4">
          <UInput v-model="state.password" type="password" placeholder="********" icon="i-heroicons-lock-closed" class="w-full" />
        </UFormField>

        <UButton type="submit" block>Login</UButton>
      </UForm>
    </UCard>
  </div>
</template> 