export default defineNuxtPlugin(() => {
  const { initializeAuth } = useAuth()
  
  // Initialize auth state from localStorage on client side
  initializeAuth()
})