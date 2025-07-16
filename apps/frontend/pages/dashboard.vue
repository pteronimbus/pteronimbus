<script setup lang="ts">
definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const { t } = useI18n()
const { user, initializeAuth } = useAuth()
const { currentTenant, initializeTenant } = useTenant()

// Initialize auth and tenant state
initializeAuth()
initializeTenant()

const router = useRouter()

// Redirect to tenant selection if no tenant is selected
onMounted(() => {
  if (!currentTenant.value) {
    router.push('/tenants')
  } else {
    // Redirect to tenant-specific dashboard
    router.push(`/tenant/${currentTenant.value.id}/dashboard`)
  }
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center">
      <UIcon name="heroicons:arrow-path" class="w-8 h-8 animate-spin mx-auto mb-4 text-primary-500" />
      <p class="text-gray-600">Redirecting to dashboard...</p>
    </div>
  </div>
</template> 