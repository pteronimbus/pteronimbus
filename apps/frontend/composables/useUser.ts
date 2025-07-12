export const useUser = () => {
  // Initialize state as null, will be restored by plugin
  const user = useState('user', () => null)

  // Save user data to localStorage
  const saveUser = (userData: any) => {
    user.value = userData
    if (process.client) {
      try {
        localStorage.setItem('pteronimbus-user', JSON.stringify(userData))
      } catch (error) {
        console.warn('Failed to save user to localStorage:', error)
      }
    }
  }

  // Clear user data from both state and localStorage
  const clearUser = () => {
    user.value = null
    if (process.client) {
      try {
        localStorage.removeItem('pteronimbus-user')
      } catch (error) {
        console.warn('Failed to clear user from localStorage:', error)
      }
    }
  }

  return {
    user,
    saveUser,
    clearUser
  }
} 