export default defineNuxtRouteMiddleware(async (to) => {
  // Skip auth middleware for login and public routes
  if (to.path === '/login' || to.path === '/auth/callback') {
    return
  }
  
  // Only run auth checks on client side
  if (!process.client) {
    return
  }
  
  const { isAuthenticated, initializeAuth } = useAuth()
  
  // Initialize auth state to ensure it's loaded
  initializeAuth()
  
  // Wait a tick to ensure the auth state is properly set
  await nextTick()
  
  // Also check localStorage directly as a fallback
  const hasTokens = localStorage.getItem('access_token') && localStorage.getItem('refresh_token')
  const userData = localStorage.getItem('user_data')

  console.log('Auth middleware check:', { 
    route: to.path, 
    isAuthenticated: isAuthenticated.value,
    hasTokens: !!hasTokens,
    userData: !!userData,
    localStorageTokens: !!localStorage.getItem('access_token')
  })

  // If not authenticated and no tokens in localStorage, redirect to login
  if (!isAuthenticated.value && !hasTokens) {
    console.log('Redirecting to login from:', to.path)
    return navigateTo('/login')
  }
})