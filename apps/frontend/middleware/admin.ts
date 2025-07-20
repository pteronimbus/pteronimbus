export default defineNuxtRouteMiddleware(async (to) => {
  // Only run on client side
  if (!process.client) {
    return
  }
  
  const { user, isAuthenticated, initializeAuth, apiRequest, isSuperAdmin } = useAuth()
  const config = useRuntimeConfig()
  
  // Initialize auth state to ensure it's loaded
  initializeAuth()
  
  // Wait a tick to ensure the auth state is properly set
  await nextTick()
  
  // Also check localStorage directly as a fallback
  const hasTokens = localStorage.getItem('access_token') && localStorage.getItem('refresh_token')
  const userData = localStorage.getItem('user_data')
  
  console.log('Admin middleware check:', { 
    route: to.path, 
    user: user.value,
    isAuthenticated: isAuthenticated.value,
    isSuperAdmin: isSuperAdmin.value,
    hasTokens: !!hasTokens,
    userData: !!userData,
    localStorageTokens: !!localStorage.getItem('access_token')
  })
  
  // Check if user is authenticated
  if (!isAuthenticated.value && !hasTokens) {
    console.log('User not authenticated, redirecting to login')
    return navigateTo('/login')
  }

  // Quick check: if user is super admin from JWT, allow access immediately
  if (isSuperAdmin.value) {
    console.log('User is super admin, allowing access')
    return
  }

  // Check admin permissions via backend API for non-super admins
  try {
    const response = await apiRequest<{ hasAdminAccess: boolean }>(`${config.public.backendUrl}/api/admin/check-access`)
    
    console.log('Admin access check:', { 
      username: user.value?.username, 
      discord_user_id: user.value?.discord_user_id,
      hasAdminAccess: response.hasAdminAccess 
    })
    
    if (!response.hasAdminAccess) {
      throw createError({
        statusCode: 403,
        statusMessage: 'Access Denied',
        message: 'You do not have permission to access this page'
      })
    }
  } catch (error: any) {
    console.error('Admin permission check failed:', error)
    
    // If it's a 403 error, show access denied
    if (error.statusCode === 403) {
      throw createError({
        statusCode: 403,
        statusMessage: 'Access Denied',
        message: 'You do not have permission to access this page'
      })
    }
    
    // For other errors (like network issues), redirect to login
    console.log('Admin check failed, redirecting to login')
    return navigateTo('/login')
  }
}) 