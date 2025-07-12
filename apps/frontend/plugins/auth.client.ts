export default defineNuxtPlugin(() => {
  // This plugin only runs on client-side after hydration
  const { user, saveUser } = useUser()

  // Restore user from localStorage if available
  if (process.client && !user.value) {
    try {
      const stored = localStorage.getItem('pteronimbus-user')
      if (stored) {
        const userData = JSON.parse(stored)
        user.value = userData
      }
    } catch (error) {
      console.warn('Failed to restore user from localStorage:', error)
      // Clear invalid data
      localStorage.removeItem('pteronimbus-user')
    }
  }
}) 