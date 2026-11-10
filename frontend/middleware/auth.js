export default defineNuxtRouteMiddleware(async (to) => {
    const authStore = useAuthStore()

    try {
        const authenticated = await authStore.checkAuth()

        if (!authenticated) {
            return navigateTo('/auth/signin')
        }

        if (to.meta.adminOnly && !authStore.isAdmin) {
            return navigateTo('/')
        }
    } catch {
        return abortNavigation(
            createError({
                statusCode: 502,
                statusMessage: 'Gagal memeriksa sesi autentikasi'
            })
        )
    }
})