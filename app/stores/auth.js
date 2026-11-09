import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        csrfToken: null,
        loading: false,
        error: null,
        validationErrors: [],
        user: null
    }),

    getters: {
        isAuthenticated: (state) => !!state.user,
        getUser: (state) => state.user,
        userRole: (state) => state.user?.role
    },

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
                        body: JSON.stringify({ email, password })
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

            const baseURL = import.meta.server
                ? config.apiUserServerBaseUrl
                : config.public.apiUserBaseUrl

            const headers = import.meta.server
                ? useRequestHeaders(['cookie'])
                : undefined

            const response = await fetch(`${baseURL}/me`, {
                method: 'GET',
                credentials: 'same-origin',
                headers
            })

            if (response.status === 401) {
                this.user = null
                this.csrfToken = null
                return false
            }

            if (!response.ok) {
                throw new Error('Failed to check authentication')
            }

            const data = await response.json()

            this.user = data.data.user ?? data.data

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

        async forgotPassword (email) {
            try {
                this.loading = true
                this.error = null

                const response = await fetch(`${useRuntimeConfig().public.apiUserBaseUrl}/forgot-password`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ email })
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

        async verifyAccount ( token) {
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
        
        async updatePasswordNoAuth (password_new, password_confirmation, token) {
            try {
                this.loading = true
                this.error = null

                const response = await fetch(`${useRuntimeConfig().public.apiUserBaseUrl}/update-password?token=${token}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ password_new, password_confirmation  })
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
        async getProfile () {
            try {
                this.loading = true
        
                const response = await fetch(`${useRuntimeConfig().public.apiUserBaseUrl}/auth/profile`, {
                    headers: {
                        'Authorization': `Bearer ${useCookie('token').value}`,
                        'Content-Type': 'application/json'
                    }
                })
        
                const result = await response.json()
                
                this.user = result.data

                return result
                
            } catch (err) {
                this.error = err.message
                throw error
            } finally {
                this.loading = false
            }
        },

        async updateProfile (userData) {
            try {
                this.loading = true
                this.error = null

                const response = await fetch(`${useRuntimeConfig().public.apiUserBaseUrl}/auth/profile`, {
                    method: 'PUT',
                    headers: {
                        'Authorization': `Bearer ${useCookie('token').value}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(userData)
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
        }
    },

    getters: {
        isAuthenticated: (state) => !!state.user,
        getUser: (state) => state.user,
        userRole: (state) => state.user?.role
    }
})