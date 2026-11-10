export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',

  modules: [
    '@pinia/nuxt',
    '@nuxt/icon'
  ],

  css: [
    '@tabler/icons-webfont/dist/tabler-icons.min.css',
    '~/assets/css/theme.min.css'
  ],

  icon: {
    serverBundle: {
      collections: ['tabler']
    },
    fallbackToApi: false
  },

  runtimeConfig: {
    apiUserServerBaseUrl: '',
    public: {
      apiProductBaseUrl: '/api/v1'
    }
  },

  app: {
    head: {
      title: 'GreenCart',
      htmlAttrs: {
        lang: 'ja'
      }
    }
  }
})