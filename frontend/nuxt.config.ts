
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
    provider: 'none',

    clientBundle: {
      scan: true,
      icons: [
        'tabler:search',
        'tabler:user',
        'tabler:shopping-bag',
        'tabler:shopping-cart',
        'tabler:menu-2',
        'tabler:layout-grid',
        'tabler:arrow-right',
        'tabler:truck-delivery',
        'tabler:leaf',
        'tabler:wallet',
        'tabler:x',
        'tabler:trash',
        'tabler:brand-facebook',
        'tabler:brand-x',
        'tabler:brand-instagram'
      ]
    }
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
