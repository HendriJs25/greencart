export default defineNuxtRouteMiddleware(async () => {
  const authStore = useAuthStore()

  await authStore.checkAuth()

  if (!authStore.isAuthenticated) {
    return navigateTo('/auth/signin')
  }

  if (!authStore.isSuperAdmin()) {
    return navigateTo('/')
  }
})