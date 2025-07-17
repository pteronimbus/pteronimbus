export default defineNuxtRouteMiddleware((to) => {
  const { isAuthenticated } = useAuth()

  console.log('Auth middleware check:', { 
    route: to.path, 
    isAuthenticated: isAuthenticated.value,
    hasTokens: !!(process.client && localStorage.getItem('access_token'))
  })

  // If not authenticated, redirect to login
  if (!isAuthenticated.value) {
    console.log('Redirecting to login from:', to.path)
    return navigateTo('/login')
  }
})