const allowedRoutes = new Set([
    'POST auth/signin',
    'POST auth/signup',
    'POST auth/verify-account',
    'POST auth/forgot-password',
    'PATCH auth/reset-password',
    'GET me',
    'GET profile',
    'PATCH profile',
    'PATCH profile/image-upload'
])

export default defineEventHandler((event) => {
    const path = (getRouterParam(event, 'path') || '')
        .replace(/^\/+|\/+$/g, '')

    const method = event.node.req.method?.toUpperCase() || 'GET'

    if (!allowedRoutes.has(`${method} ${path}`)) {
        throw createError({
            statusCode: 404,
            statusMessage: 'Route not found'
        })
    }

    const config = useRuntimeConfig(event)
    const baseURL = config.apiUserServerBaseUrl

    if (!baseURL) {
        throw createError({
            statusCode: 503,
            statusMessage: 'User API belum dikonfigurasi'
        })
    }

    const query = path === 'auth/verify-account' || path === 'auth/reset-password'
        ? getRequestURL(event).search
        : ''
    const target = `${baseURL.replace(/\/+$/, '')}/${path}${query}`

    return proxyRequest(event, target)
})