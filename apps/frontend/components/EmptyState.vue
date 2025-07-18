<template>
  <div class="text-center py-12">
    <div class="flex justify-center mb-6">
      <div class="w-16 h-16 bg-gray-100 dark:bg-gray-800 rounded-2xl flex items-center justify-center shadow-sm">
        <UIcon 
          :name="icon" 
          class="w-8 h-8 text-gray-400 dark:text-gray-500" 
        />
      </div>
    </div>
    <div class="max-w-sm mx-auto">
      <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">
        {{ title }}
      </h3>
      <p class="text-gray-600 dark:text-gray-400 text-sm leading-relaxed" :class="{ 'mb-6': hasAction }">
        {{ description }}
      </p>
    </div>
    <div v-if="hasAction" class="mt-6">
      <slot name="actions">
        <UButton 
          v-if="actionLabel"
          :icon="actionIcon"
          :class="['shadow-sm hover:shadow-md transition-all duration-200', actionClass]"
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
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md'
})

const emit = defineEmits<{
  action: []
}>()

const hasAction = computed(() => {
  return !!props.actionLabel || !!useSlots().actions
})
</script> 