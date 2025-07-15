export default defineNuxtRouteMiddleware((to) => {
  const { isAuthenticated, initializeAuth } = useAuth()

  // Initialize auth state from localStorage
  initializeAuth()

  // If not authenticated, redirect to login
  if (!isAuthenticated.value) {
    return navigateTo('/login')
  }
})