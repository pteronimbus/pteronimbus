export const useAppButtons = () => {
  const getButtonConfig = (type: 'primary' | 'secondary' | 'success' | 'error' | 'warning' | 'info' | 'neutral') => {
    const configs = {
      primary: {
        color: 'primary' as const,
        variant: 'outline' as const,
        class: 'bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors'
      },
      secondary: {
        color: 'neutral' as const,
        variant: 'ghost' as const,
        class: 'bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors'
      },
      success: {
        color: 'success' as const,
        variant: 'outline' as const,
        class: 'bg-white dark:bg-gray-700 hover:bg-green-50 dark:hover:bg-green-900/20 border-green-200 dark:border-green-700 hover:border-green-300 dark:hover:border-green-600 transition-colors'
      },
      error: {
        color: 'error' as const,
        variant: 'outline' as const,
        class: 'bg-white dark:bg-gray-700 hover:bg-red-50 dark:hover:bg-red-900/20 border-red-200 dark:border-red-700 hover:border-red-300 dark:hover:border-red-600 transition-colors'
      },
      warning: {
        color: 'warning' as const,
        variant: 'outline' as const,
        class: 'bg-white dark:bg-gray-700 hover:bg-yellow-50 dark:hover:bg-yellow-900/20 border-yellow-200 dark:border-yellow-700 hover:border-yellow-300 dark:hover:border-yellow-600 transition-colors'
      },
      info: {
        color: 'info' as const,
        variant: 'outline' as const,
        class: 'bg-white dark:bg-gray-700 hover:bg-blue-50 dark:hover:bg-blue-900/20 border-blue-200 dark:border-blue-700 hover:border-blue-300 dark:hover:border-blue-600 transition-colors'
      },
      neutral: {
        color: 'neutral' as const,
        variant: 'outline' as const,
        class: 'bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors'
      }
    }
    
    return configs[type]
  }

  const createActionButton = (label: string, onClick: () => void, options: {
    type?: 'primary' | 'secondary' | 'success' | 'error' | 'warning' | 'info' | 'neutral'
    icon?: string
    size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
    loading?: boolean
  } = {}) => {
    const { type = 'primary', icon, size = 'sm', loading = false } = options
    const config = getButtonConfig(type)
    
    return {
      label,
      onClick,
      icon,
      size,
      loading,
      ...config
    }
  }

  return {
    getButtonConfig,
    createActionButton
  }
} 