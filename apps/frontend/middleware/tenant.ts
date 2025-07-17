export default defineNuxtRouteMiddleware(async (to) => {
  const { currentTenant, initializeTenant, fetchTenant, storeCurrentTenant } = useTenant()
  
  // Initialize tenant state from localStorage
  initializeTenant()
  
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
        storeCurrentTenant(tenant)
        console.log('Successfully switched to tenant:', tenant.name)
      } catch (error) {
        console.error('Failed to switch tenant:', error)
        // If we can't fetch the tenant, redirect to tenant selection
        return navigateTo('/tenants')
      }
    }
  }
})