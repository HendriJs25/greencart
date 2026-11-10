export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',

  modules: [
    '@pinia/nuxt'
  ],

  runtimeConfig: {
    apiUserServerBaseUrl: ''
  },

  app: {
    head: {
      title: 'GreenCart',
      meta: [
        {
          name: 'viewport',
          content: 'width=device-width, initial-scale=1'
        }
      ]
    }
  }
})