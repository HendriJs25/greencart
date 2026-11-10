
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null,
        csrfToken: null,
        loading: false,
        error: null
    }),

    getters: {
        isAuthenticated: (state) => Boolean(state.user)
    },

    actions: {
        async signin(email, password) {
            this.loading = true
            this.error = null

            try {
                const result = await $fetch('/api/user/auth/signin', {
                    method: 'POST',
                    body: { email, password }
                })

                if (!result?.data?.user) {
                    throw new Error('Respons login tidak valid')
                }

                this.user = result.data.user
                this.csrfToken = result.data.csrf_token ?? null

                if (import.meta.client && this.csrfToken) {
                    sessionStorage.setItem(
                        'greencart_csrf',
                        this.csrfToken
                    )
                }

                // Pastikan sesi benar-benar diterima backend
                const authenticated = await this.checkAuth()

                if (!authenticated) {
                    throw new Error(
                        'Login berhasil tetapi sesi tidak tersimpan'
                    )
                }

                return true
            } catch (err) {
                this.user = null
                this.error =
                    err.data?.message ||
                    err.message ||
                    'Login gagal'
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
