export default defineNuxtRouteMiddleware((to, from) => {
  const user = useUser()

  if (to.meta.auth === false) {
    return
  }

  if (!user.value && to.path !== '/login') {
    return navigateTo('/login')
  }

  if (user.value && to.path === '/login') {
    return navigateTo('/dashboard')
  }
}) 