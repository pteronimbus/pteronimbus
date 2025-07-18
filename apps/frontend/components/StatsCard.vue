<template>
  <UCard 
    :class="[
      'transition-all duration-200 overflow-hidden border border-gray-200 dark:border-gray-700',
      clickable ? 'cursor-pointer hover:shadow-lg hover:scale-105 hover:border-primary-300 dark:hover:border-primary-600' : '',
      clickable && hoverable ? 'hover:bg-gray-50 dark:hover:bg-gray-800' : ''
    ]"
    @click="handleClick"
  >
    <div class="p-6">
      <div class="flex items-center justify-between">
        <div class="flex-1">
          <div class="flex items-center justify-between mb-2">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">{{ label }}</p>
            <div v-if="icon" class="flex-shrink-0">
              <div :class="[
                'w-12 h-12 rounded-xl flex items-center justify-center shadow-sm',
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
          
          <div class="flex items-baseline space-x-2 mb-3">
            <p class="text-3xl font-bold text-gray-900 dark:text-gray-100">{{ value }}</p>
            <p v-if="total" class="text-lg text-gray-500 dark:text-gray-400">/ {{ total }}</p>
          </div>
          
          <div v-if="trend" class="flex items-center">
            <div class="flex items-center space-x-1">
              <UIcon 
                :name="trend.startsWith('+') ? 'heroicons:arrow-trending-up' : trend.startsWith('-') ? 'heroicons:arrow-trending-down' : 'heroicons:minus'"
                :class="[
                  'w-4 h-4',
                  trend.startsWith('+') ? 'text-green-500' : trend.startsWith('-') ? 'text-red-500' : 'text-gray-500'
                ]"
              />
              <span :class="[
                'text-sm font-medium',
                trend.startsWith('+') ? 'text-green-600 dark:text-green-400' : 
                trend.startsWith('-') ? 'text-red-600 dark:text-red-400' : 
                'text-gray-600 dark:text-gray-400'
              ]">{{ trend }}</span>
            </div>
            <span v-if="trendLabel" class="ml-2 text-sm text-gray-500 dark:text-gray-400">{{ trendLabel }}</span>
          </div>
        </div>
      </div>
      
      <!-- Progress bar for total/value comparison -->
      <div v-if="total && value !== total" class="mt-4">
        <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
          <div 
            :class="[
              'h-2 rounded-full transition-all duration-300',
              getProgressBarColor()
            ]"
            :style="`width: ${Math.min((Number(value) / Number(total)) * 100, 100)}%`"
          ></div>
        </div>
        <div class="flex justify-between text-xs text-gray-500 dark:text-gray-400 mt-1">
          <span>{{ Math.round((Number(value) / Number(total)) * 100) }}% used</span>
          <span>{{ Number(total) - Number(value) }} remaining</span>
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
    emerald: 'bg-emerald-100 dark:bg-emerald-900/30',
    blue: 'bg-blue-100 dark:bg-blue-900/30',
    purple: 'bg-purple-100 dark:bg-purple-900/30',
    green: 'bg-green-100 dark:bg-green-900/30',
    yellow: 'bg-yellow-100 dark:bg-yellow-900/30',
    cyan: 'bg-cyan-100 dark:bg-cyan-900/30',
    orange: 'bg-orange-100 dark:bg-orange-900/30',
    red: 'bg-red-100 dark:bg-red-900/30',
    gray: 'bg-gray-100 dark:bg-gray-900/30',
    indigo: 'bg-indigo-100 dark:bg-indigo-900/30'
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
    gray: 'text-gray-600 dark:text-gray-400',
    indigo: 'text-indigo-600 dark:text-indigo-400'
  }
  return colorMap[props.color] || colorMap.blue
}

const getProgressBarColor = () => {
  const colorMap: Record<string, string> = {
    emerald: 'bg-emerald-500',
    blue: 'bg-blue-500',
    purple: 'bg-purple-500',
    green: 'bg-green-500',
    yellow: 'bg-yellow-500',
    cyan: 'bg-cyan-500',
    orange: 'bg-orange-500',
    red: 'bg-red-500',
    gray: 'bg-gray-500',
    indigo: 'bg-indigo-500'
  }
  return colorMap[props.color] || colorMap.blue
}
</script> 