export default defineNuxtRouteMiddleware(async (to) => {
  const { currentTenant, initializeTenant, fetchTenant } = useTenant()
  const { initializeAuth } = useAuth()
  
  // Initialize auth state first to ensure it's loaded
  initializeAuth()
  
  // Initialize tenant state from localStorage
  initializeTenant()
  
  // Wait a tick to ensure the states are properly set
  await nextTick()
  
  // Skip tenant middleware for admin routes
  if (to.path.startsWith('/admin/')) {
    return
  }
  
  // Check if route requires tenant context
  const requiresTenant = to.path.startsWith('/tenant/') || 
                        to.meta.requiresTenant === true
  
  if (requiresTenant && !currentTenant.value) {
    // Redirect to tenant selection if no tenant is selected
    return navigateTo('/tenants')
  }
  
  // If accessing tenant-specific route, validate tenant ID
  if (to.path.startsWith('/tenant/')) {
    const tenantIdFromRoute = to.params.tenantId as string
    
    if (currentTenant.value && currentTenant.value.id !== tenantIdFromRoute) {
      // Tenant ID in URL doesn't match current tenant
      console.log('Tenant ID mismatch, attempting to switch tenant:', {
        current: currentTenant.value.id,
        route: tenantIdFromRoute
      })
      
      try {
        // Try to fetch and switch to the tenant from the route
        const tenant = await fetchTenant(tenantIdFromRoute)
        console.log('Successfully fetched tenant:', tenant.name)
        // Note: We don't need to store it here since the route will handle the switch
      } catch (error) {
        console.error('Failed to switch tenant:', error)
        // If we can't fetch the tenant, redirect to tenant selection
        return navigateTo('/tenants')
      }
    }
  }
})