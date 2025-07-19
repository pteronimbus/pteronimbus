<template>
  <UBadge 
    :color="getStatusColor(status)"
    variant="subtle"
    :class="[
      'capitalize',
      getTextColorClass(status)
    ]"
  >
    {{ displayLabel }}
  </UBadge>
</template>

<script setup lang="ts">
interface Props {
  status: string
  type?: 'server' | 'user' | 'alert' | 'player' | 'controller' | 'custom'
  customColors?: Record<string, string>
  label?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'custom'
})

const { t } = useI18n()

const displayLabel = computed(() => {
  if (props.label) {
    return props.label
  }
  
  // Use i18n keys based on type
  if (props.type === 'server') {
    return t(`servers.status.${props.status}`)
  } else if (props.type === 'user') {
    return t(`users.status.${props.status}`)
  } else if (props.type === 'alert') {
    return t(`alerts.severity.${props.status}`)
  } else if (props.type === 'controller') {
    return t(`admin.controllers.statuses.${props.status}`)
  }
  
  return props.status
})

const getStatusColor = (status: string) => {
  // Check for custom colors first
  if (props.customColors && props.customColors[status]) {
    return props.customColors[status] as any
  }

  // Default color mappings
  const colorMap: Record<string, string> = {
    // Server statuses
    online: 'success',
    offline: 'error',
    starting: 'warning',
    stopping: 'warning',
    error: 'error',
    degraded: 'orange',
    
    // Controller statuses
    pending_approval: 'warning',
    active: 'success',
    degraded: 'orange',
    rejected: 'error',
    
    // User statuses
    banned: 'error',
    suspended: 'warning',
    active: 'success',
    
    // Alert severities
    critical: 'error',
    warning: 'warning',
    info: 'primary',
    
    // General statuses
    success: 'success',
    pending: 'warning',
    failed: 'error'
  }
  
  return (colorMap[status] || 'neutral') as any
}

const getTextColorClass = (status: string) => {
  const color = getStatusColor(status)
  const textColorMap: Record<string, string> = {
    success: 'text-green-700 dark:text-green-300',
    error: 'text-red-700 dark:text-red-300',
    warning: 'text-yellow-700 dark:text-yellow-300',
    primary: 'text-blue-700 dark:text-blue-300',
    neutral: 'text-gray-700 dark:text-gray-300'
  }
  
  return textColorMap[color] || 'text-gray-700 dark:text-gray-300'
}
</script> 