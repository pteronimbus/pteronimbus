<template>
  <UCard>
    <template #header>
      <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ title }}</h3>
    </template>
    <div :class="[
      'grid gap-3',
      gridCols === 1 ? 'grid-cols-1' : '',
      gridCols === 2 ? 'grid-cols-2' : '',
      gridCols === 3 ? 'grid-cols-3' : '',
      gridCols === 4 ? 'grid-cols-4' : ''
    ]">
      <UButton
        v-for="action in actions"
        :key="action.label"
        :color="action.color || 'primary'"
        :variant="action.variant || 'soft'"
        :size="action.size || 'lg'"
        :icon="action.icon"
        :class="['justify-start', action.class]"
        @click="action.onClick"
      >
        {{ action.label }}
      </UButton>
    </div>
  </UCard>
</template>

<script setup lang="ts">
type ButtonColor = 'success' | 'neutral' | 'error' | 'warning' | 'primary' | 'info' | 'secondary'
type ButtonVariant = 'link' | 'solid' | 'outline' | 'soft' | 'subtle' | 'ghost'
type ButtonSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl'

interface QuickAction {
  label: string
  icon?: string
  color?: ButtonColor
  variant?: ButtonVariant
  size?: ButtonSize
  class?: string
  onClick: () => void
}

interface Props {
  title?: string
  actions: QuickAction[]
  gridCols?: number
}

withDefaults(defineProps<Props>(), {
  title: 'Quick Actions',
  gridCols: 2
})
</script> 