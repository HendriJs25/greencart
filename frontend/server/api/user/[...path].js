const allowedRoutes = new Set([
    'POST auth/signin',
    'GET me',
    'GET profile',
    'PATCH profile'
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

    const target = `${baseURL.replace(/\/+$/, '')}/${path}`

    return proxyRequest(event, target)
})