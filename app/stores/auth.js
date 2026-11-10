import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        csrfToken: null,
        loading: false,
        error: null,
        validationErrors: [],
        user: null,
        profile: null
    }),


    actions: {
        async signin(email, password) {
            try {
                this.loading = true
                this.error = null
                this.validationErrors = []

                const response = await fetch(
                    `${useRuntimeConfig().public.apiUserBaseUrl}/auth/signin`,
                    {
                        method: 'POST',
                        credentials: 'same-origin',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({email, password})
                    }
                )

                const data = await response.json()

                if (!response.ok) {
                    if (Array.isArray(data.message)) {
                        this.validationErrors = data.message
                        throw new Error('Validation failed')
                    }

                    throw new Error(data.message || 'Something went wrong!')
                }

                this.user = data.data.user
                this.csrfToken = data.data.csrf_token

                if (import.meta.client && this.csrfToken) {
                    sessionStorage.setItem(
                        'greencart_csrf',
                        this.csrfToken
                    )
                }

                return data
            } catch (error) {
                if (this.validationErrors.length === 0) {
                    this.error = error.message
                }

                throw error
            } finally {
                this.loading = false
            }
        },

        isSuperAdmin() {
            return this.user?.role?.includes('admin') ?? false
        },

        logout() {
            this.user = null
            this.token = null
            const tokenCookie = useCookie('token')
            tokenCookie.value = null
            const userCookie = useCookie('user')
            userCookie.value = null
            // localStorage.removeItem('token')
        },

        async checkAuth() {
            const config = useRuntimeConfig()

            const baseURL = (
                import.meta.server
                    ? config.apiUserServerBaseUrl
                    : config.public.apiUserBaseUrl
            ).replace(/\/$/, '')

            const headers = import.meta.server
                ? useRequestHeaders(['cookie'])
                : {}

            const response = await fetch(`${baseURL}/me`, {
                method: 'GET',
                credentials: 'same-origin',
                headers: {
                    ...headers,
                    Accept: 'application/json'
                }
            })

            if (response.status === 401) {
                this.user = null
                this.profile = null
                this.csrfToken = null

                if (import.meta.client) {
                    sessionStorage.removeItem('greencart_csrf')
                }

                return false
            }

            const result = await response.json()

            if (!response.ok) {
                throw new Error(
                    typeof result.message === 'string'
                        ? result.message
                        : 'Gagal memeriksa sesi login'
                )
            }

            this.user = result.data

            if (import.meta.client && !this.csrfToken) {
                this.csrfToken =
                    sessionStorage.getItem('greencart_csrf')
            }

            return true
        },

        async signup(userData) {
            try {
                this.loading = true
                this.error = null

                const response = await fetch(`${useRuntimeConfig().public.apiUserBaseUrl}/signup`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(userData)
                })

                const data = await response.json()

                if (!response.ok) {
                    console.log(userData)
                    throw new Error(data.message || 'Terjadi kesalahan saat registrasi')
                }

                return data
            } catch (error) {
                this.error = error.message
                throw error
            } finally {
                this.loading = false
            }
        },

        async forgotPassword(email) {
            try {
                this.loading = true
                this.error = null

                const response = await fetch(`${useRuntimeConfig().public.apiUserBaseUrl}/forgot-password`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({email})
                })

                const data = await response.json()
                if (!response.ok) {
                    throw new Error(data.message || 'Something went wrong!');
                }

                return data
            } catch (error) {
                this.error = error.message
                throw error
            } finally {
                this.loading = false
            }
        },

        async verifyAccount(token) {
            try {
                this.loading = true
                this.error = null

                const response = await fetch(`${useRuntimeConfig().public.apiUserBaseUrl}/verify-account?token=${token}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                })

                const data = await response.json()
                if (!response.ok) {
                    throw new Error(data.message || 'Something went wrong!');
                }

                return data
            } catch (error) {
                this.error = error.message
                throw error
            } finally {
                this.loading = false
            }
        },

        async updatePasswordNoAuth(password_new, password_confirmation, token) {
            try {
                this.loading = true
                this.error = null

                const response = await fetch(`${useRuntimeConfig().public.apiUserBaseUrl}/update-password?token=${token}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({password_new, password_confirmation})
                })

                const data = await response.json()
                if (!response.ok) {
                    throw new Error(data.message || 'Something went wrong!');
                }

                return data
            } catch (error) {
                this.error = error.message
                throw error
            } finally {
                this.loading = false
            }
        },

        async getProfile() {
            this.loading = true
            this.error = null

            try {
                const baseURL = useRuntimeConfig()
                    .public.apiUserBaseUrl.replace(/\/$/, '')

                const response = await fetch(`${baseURL}/profile`, {
                    method: 'GET',
                    credentials: 'same-origin',
                    headers: {
                        Accept: 'application/json'
                    }
                })

                const result = await response.json()

                if (!response.ok) {
                    if (response.status === 401) {
                        this.user = null
                        this.profile = null
                    }

                    throw new Error(
                        typeof result.message === 'string'
                            ? result.message
                            : 'Gagal mengambil profil customer'
                    )
                }

                this.profile = result.data

                return result

            } catch (error) {
                this.error = error.message
                throw error
            } finally {
                this.loading = false
            }
        },

        async updateProfile(userData) {
            this.loading = true
            this.error = null

            try {
                const csrfToken = this.csrfToken ||
                    (import.meta.client
                        ? sessionStorage.getItem('greencart_csrf')
                        : null)

                if (!csrfToken) {
                    throw new Error(
                        'CSRF token tidak tersedia. Silakan login kembali.'
                    )
                }

                const baseURL = useRuntimeConfig()
                    .public.apiUserBaseUrl.replace(/\/$/, '')

                const response = await fetch(`${baseURL}/profile`, {
                    method: 'PATCH',
                    credentials: 'same-origin',
                    headers: {
                        'Content-Type': 'application/json',
                        'Accept': 'application/json',
                        'X-CSRF-Token': csrfToken
                    },
                    body: JSON.stringify(userData)
                })

                const result = await response.json()

                if (!response.ok) {
                    throw new Error(
                        typeof result.message === 'string'
                            ? result.message
                            : 'Gagal memperbarui profil customer'
                    )
                }

                // Backend tidak mengembalikan profil terbaru
                await this.getProfile()

                // Sinkronkan nama pada identitas user
                if (this.user && this.profile) {
                    this.user = {
                        ...this.user,
                        name: this.profile.name
                    }
                }

                return result

            } catch (error) {
                this.error = error.message
                throw error
            } finally {
                this.loading = false
            }
        },

        getters: {
            isAuthenticated: (state) => !!state.user,
            getUser: (state) => state.user,
            userRole: (state) => state.user?.role
        }
    }
})