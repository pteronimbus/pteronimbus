export default defineNuxtRouteMiddleware((to) => {
  const { currentTenant, initializeTenant } = useTenant()
  
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
      // This could happen if user bookmarked a URL or shared a link
      // We should either switch to the tenant or redirect to tenant selection
      console.warn('Tenant ID mismatch:', {
        current: currentTenant.value.id,
        route: tenantIdFromRoute
      })
      
      // For now, redirect to tenant selection
      return navigateTo('/tenants')
    }
  }
})