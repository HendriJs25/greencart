export default defineNuxtRouteMiddleware(async () => {
    const authStore = useAuthStore()

    try {
        const authenticated = await authStore.checkAuth()

        if (!authenticated) {
            return navigateTo('/auth/signin')
        }
    } catch (err) {
        return abortNavigation(
            createError({
                statusCode: 502,
                statusMessage: 'Gagal memeriksa sesi autentikasi'
            })
        )
    }
})