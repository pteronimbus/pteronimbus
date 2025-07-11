export default defineAppConfig({
  ui: {
    colors: {
      primary: 'blue',
      neutral: 'cool'
    },
    formField: {
      slots: {
        label: 'block font-medium text-gray-900 dark:text-gray-100'
      }
    },
    card: {
      slots: {
        root: 'bg-white dark:bg-gray-800 ring-1 ring-gray-200 dark:ring-gray-700 divide-gray-200 dark:divide-gray-700 rounded-lg shadow-sm'
      }
    },
    input: {
      slots: {
        leadingIcon: 'text-gray-600 dark:text-gray-300',
        trailingIcon: 'text-gray-600 dark:text-gray-300'
      }
    }
  }
}) 