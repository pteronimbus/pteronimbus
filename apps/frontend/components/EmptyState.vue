<template>
  <div class="text-center py-12">
    <UIcon 
      :name="icon" 
      class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" 
    />
    <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
      {{ title }}
    </h3>
    <p class="text-gray-500 dark:text-gray-400" :class="{ 'mb-6': hasAction }">
      {{ description }}
    </p>
    <div v-if="hasAction">
      <slot name="action">
        <UButton 
          v-if="actionLabel"
          :icon="actionIcon"
          :class="actionClass"
          @click="$emit('action')"
        >
          {{ actionLabel }}
        </UButton>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  icon: string
  title: string
  description: string
  actionLabel?: string
  actionIcon?: string
  actionClass?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  action: []
}>()

const hasAction = computed(() => {
  return !!props.actionLabel || !!useSlots().action
})
</script> 