export default defineNuxtRouteMiddleware((to) => {
  const { isAuthenticated, initializeAuth } = useAuth()

  // Initialize auth state from localStorage
  initializeAuth()

  // If authenticated, redirect to dashboard
  if (isAuthenticated.value) {
    return navigateTo('/dashboard')
  }
})