<template>
  <UCard class="overflow-hidden">
    <template #header>
      <div class="flex items-center space-x-2">
        <div class="w-2 h-2 bg-primary-500 rounded-full"></div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ title }}</h3>
      </div>
    </template>
    <div :class="[
      'grid gap-4',
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
        :class="[
          'justify-start h-auto py-4 px-4 shadow-sm hover:shadow-md transition-all duration-200 hover:-translate-y-0.5 group',
          action.class
        ]"
        @click="action.onClick"
      >
        <div class="flex items-center space-x-3 w-full">
          <div v-if="action.icon" class="flex-shrink-0">
            <UIcon 
              :name="action.icon" 
              class="w-5 h-5 group-hover:scale-110 transition-transform duration-200" 
            />
          </div>
          <div class="flex-1 text-left">
            <div class="font-medium">{{ action.label }}</div>
            <div v-if="action.description" class="text-xs opacity-75 mt-1">{{ action.description }}</div>
          </div>
        </div>
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
  description?: string
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