<template>
  <div>
    <PageHeader
      :title="t('logs.title')"
      :description="t('logs.description')"
    />

    <div class="space-y-6">
      <UCard>
        <div class="text-center py-12">
          <UIcon name="heroicons:document-text" class="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
            {{ t('logs.comingSoon') }}
          </h3>
          <p class="text-gray-600 dark:text-gray-400 mb-6">
            {{ t('logs.comingSoonDesc') }}
          </p>
          <UButton
            variant="outline"
            @click="router.push(`/tenant/${currentTenant?.id}/dashboard`)"
          >
            {{ t('common.back') }}
          </UButton>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: ['auth', 'tenant']
})

const { t } = useI18n()
const router = useRouter()
const { currentTenant } = useTenant()
</script>