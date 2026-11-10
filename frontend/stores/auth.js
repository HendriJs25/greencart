
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null,
        csrfToken: null,
        loading: false,
        error: null,
        validationErrors: []
    }),

    getters: {
        isAuthenticated: (state) => Boolean(state.user)
    },

    actions: {

        async signin(email, password) {
            this.loading = true
            this.error = null
            this.validationErrors = []

            try {
                const result = await $fetch('/api/user/auth/signin', {
                    method: 'POST',
                    body: {
                        email,
                        password
                    }
                })

                if (!result?.data?.user) {
                    throw new Error('ログインの応答が正しくありません。')
                }

                this.user = result.data.user
                this.csrfToken = result.data.csrf_token ?? null

                if (import.meta.client && this.csrfToken) {
                    sessionStorage.setItem(
                        'greencart_csrf',
                        this.csrfToken
                    )
                }

                const authenticated = await this.checkAuth()

                if (!authenticated) {
                    throw new Error('ログインセッションを確認できませんでした。')
                }

                return result

            } catch (err) {
                this.user = null

                // $fetch / ofetch menempatkan respons backend di .data
                const payload = err?.data ?? err?.response?._data ?? {}
                const status = err?.response?.status ?? err?.statusCode

                // Mendukung format:
                // { message: [{ field, message }] }
                // { errors: [{ field, message }] }
                // { errors: { email: ["..."] } }
                const rawErrors = Array.isArray(payload.message)
                    ? payload.message
                    : payload.errors ?? []

                if (Array.isArray(rawErrors)) {
                    this.validationErrors = rawErrors
                        .filter(item =>
                            item &&
                            typeof item === 'object' &&
                            item.field &&
                            typeof item.message === 'string'
                        )
                        .map(item => ({
                            field: String(item.field).toLowerCase().split('.').pop(),
                            message: item.message
                        }))
                } else if (rawErrors && typeof rawErrors === 'object') {
                    this.validationErrors = Object.entries(rawErrors)
                        .flatMap(([field, messages]) =>
                            (Array.isArray(messages) ? messages : [messages])
                                .filter(message => typeof message === 'string')
                                .map(message => ({
                                    field: field.toLowerCase().split('.').pop(),
                                    message
                                }))
                        )
                }

                if (this.validationErrors.length === 0) {
                    const backendMessage = [
                        payload.message,
                        payload.error,
                        payload.detail
                    ].find(value => typeof value === 'string' && value.trim())

                    this.error = status >= 500
                        ? 'サーバーエラーが発生しました。しばらくしてからもう一度お試しください。'
                        : backendMessage ||
                        (status ? 'ログインに失敗しました。' : err.message || '接続エラーが発生しました。')
                }

                throw err
            } finally {
                this.loading = false
            }
        },

        async checkAuth() {
            const fetcher = import.meta.server
                ? useRequestFetch()
                : $fetch

            try {
                const result = await fetcher('/api/user/me')

                if (!result?.data) {
                    throw new Error('Respons sesi tidak valid')
                }

                this.user = result.data
                return true
            } catch (err) {
                if (
                    err.statusCode === 401 ||
                    err.status === 401
                ) {
                    this.user = null
                    return false
                }

                throw err
            }
        },

        async getProfile() {
            const fetcher = import.meta.server
                ? useRequestFetch()
                : $fetch

            const result = await fetcher('/api/user/profile')
            return result.data
        },

        async updateProfile(payload) {
            const csrfToken = this.csrfToken ||
                (import.meta.client
                    ? sessionStorage.getItem('greencart_csrf')
                    : null)

            if (!csrfToken) {
                throw new Error(
                    'CSRF token tidak tersedia. Silakan login kembali.'
                )
            }

            const result = await $fetch('/api/user/profile', {
                method: 'PATCH',
                headers: {
                    'X-CSRF-Token': csrfToken
                },
                body: payload
            })

            await this.checkAuth()
            return result
        }
    }
})
