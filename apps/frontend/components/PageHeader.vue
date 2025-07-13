<template>
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
    <div>
      <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">{{ title }}</h1>
      <p v-if="description" class="mt-1 text-gray-500 dark:text-gray-400">{{ description }}</p>
    </div>
    <div v-if="$slots.actions || actions.length > 0" class="mt-4 sm:mt-0 flex items-center gap-2">
      <slot name="actions">
        <template v-for="action in actions" :key="action.label">
          <UBadge 
            v-if="action.type === 'badge'"
            :color="(action.color as any) || 'success'" 
            variant="subtle"
          >
            <UIcon v-if="action.icon" :name="action.icon" class="w-4 h-4 mr-1" />
            {{ action.label }}
          </UBadge>
          <UButton 
            v-else
            :color="(action.color as any) || 'primary'"
            :variant="(action.variant as any) || 'solid'"
            :icon="action.icon"
            :size="(action.size as any) || 'sm'"
            :class="action.class"
            @click="action.onClick"
          >
            {{ action.label }}
          </UButton>
        </template>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
type ButtonColor = 'success' | 'neutral' | 'error' | 'warning' | 'primary' | 'info' | 'secondary'
type ButtonVariant = 'link' | 'solid' | 'outline' | 'soft' | 'subtle' | 'ghost'
type ButtonSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl'

interface Action {
  label: string
  type?: 'button' | 'badge'
  color?: ButtonColor
  variant?: ButtonVariant
  icon?: string
  size?: ButtonSize
  class?: string
  onClick?: () => void
}

interface Props {
  title: string
  description?: string
  actions?: Action[]
}

withDefaults(defineProps<Props>(), {
  actions: () => []
})
</script> 