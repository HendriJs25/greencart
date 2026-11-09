export default defineNuxtRouteMiddleware(async (to) => {
    if (!to.meta.requiresAuth) {
        return
    }

    const authStore = useAuthStore()

    await authStore.checkAuth()

    if (!authStore.isAuthenticated) {
        return navigateTo('/auth/signin')
    }
})