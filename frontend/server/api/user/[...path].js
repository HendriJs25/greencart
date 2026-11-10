const allowedRoutes = new Set([
    'POST auth/signin',
    'POST auth/signup',
    'POST auth/verify-account',
    'POST auth/forgot-password',
    'PATCH auth/reset-password',
    'GET me',
    'GET profile',
    'PATCH profile',
    'PATCH profile/image-upload',
    'GET admin/customers'
])

export default defineEventHandler((event) => {
    const path = (getRouterParam(event, 'path') || '')
        .replace(/^\/+|\/+$/g, '')

    const method = event.node.req.method?.toUpperCase() || 'GET'

    const isCustomerDetail = method === 'GET' &&
        /^admin\/customers\/[1-9]\d{0,18}$/.test(path)

    if (!allowedRoutes.has(`${method} ${path}`) && !isCustomerDetail) {
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

    if (path === 'admin/customers' || isCustomerDetail) {
        setResponseHeader(event, 'Cache-Control', 'private, no-store')
    }

    const query = ['auth/verify-account', 'auth/reset-password', 'admin/customers'].includes(path)
        ? getRequestURL(event).search
        : ''

    const target = `${baseURL.replace(/\/+$/, '')}/${path}${query}`

    return proxyRequest(event, target)
})