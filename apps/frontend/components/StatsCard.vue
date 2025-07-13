<template>
  <UCard 
    :class="[
      'transition-all duration-200',
      clickable ? 'cursor-pointer hover:shadow-lg hover:scale-105' : '',
      clickable && hoverable ? 'hover:bg-gray-50 dark:hover:bg-gray-800' : ''
    ]"
    @click="handleClick"
  >
    <div class="flex items-center justify-between">
      <div class="flex-1">
        <p class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ label }}</p>
        <div class="flex items-baseline mt-1">
          <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ value }}</p>
          <p v-if="total" class="ml-2 text-sm text-gray-500 dark:text-gray-400">/ {{ total }}</p>
        </div>
        <div v-if="trend" class="flex items-center mt-2">
          <span :class="[trendColor, 'text-xs font-medium']">{{ trend }}</span>
          <span v-if="trendLabel" class="ml-1 text-xs text-gray-500 dark:text-gray-400">{{ trendLabel }}</span>
        </div>
      </div>
      <div v-if="icon" class="flex-shrink-0">
        <div :class="[
          'p-3 rounded-full',
          getIconBackgroundClass()
        ]">
          <UIcon 
            :name="icon" 
            :class="[
              'w-6 h-6',
              getIconColorClass()
            ]"
          />
        </div>
      </div>
    </div>
  </UCard>
</template>

<script setup lang="ts">
interface Props {
  label: string
  value: string | number
  total?: string | number
  icon?: string
  color?: string
  trend?: string
  trendColor?: string
  trendLabel?: string
  clickable?: boolean
  hoverable?: boolean
  to?: string
}

const props = withDefaults(defineProps<Props>(), {
  color: 'blue',
  clickable: false,
  hoverable: true,
  trendLabel: 'from last hour'
})

const emit = defineEmits<{
  click: []
}>()

const router = useRouter()

const handleClick = () => {
  if (props.clickable) {
    if (props.to) {
      router.push(props.to)
    }
    emit('click')
  }
}

const getIconBackgroundClass = () => {
  const colorMap: Record<string, string> = {
    emerald: 'bg-emerald-100 dark:bg-emerald-900',
    blue: 'bg-blue-100 dark:bg-blue-900',
    purple: 'bg-purple-100 dark:bg-purple-900',
    green: 'bg-green-100 dark:bg-green-900',
    yellow: 'bg-yellow-100 dark:bg-yellow-900',
    cyan: 'bg-cyan-100 dark:bg-cyan-900',
    orange: 'bg-orange-100 dark:bg-orange-900',
    red: 'bg-red-100 dark:bg-red-900',
    gray: 'bg-gray-100 dark:bg-gray-900'
  }
  return colorMap[props.color] || colorMap.blue
}

const getIconColorClass = () => {
  const colorMap: Record<string, string> = {
    emerald: 'text-emerald-600 dark:text-emerald-400',
    blue: 'text-blue-600 dark:text-blue-400',
    purple: 'text-purple-600 dark:text-purple-400',
    green: 'text-green-600 dark:text-green-400',
    yellow: 'text-yellow-600 dark:text-yellow-400',
    cyan: 'text-cyan-600 dark:text-cyan-400',
    orange: 'text-orange-600 dark:text-orange-400',
    red: 'text-red-600 dark:text-red-400',
    gray: 'text-gray-600 dark:text-gray-400'
  }
  return colorMap[props.color] || colorMap.blue
}
</script> 