// https://nuxt.com/docs/api/configuration/nuxt-config

export default defineNuxtConfig({
  compatibilityDate: '2025-05-15',
  devtools: { enabled: true },

  css: ['~/assets/css/main.css'],

  modules: [
    '@nuxt/icon',
    '@nuxt/fonts',
    '@nuxt/test-utils/module',
    '@nuxt/ui',
    '@nuxtjs/i18n',
    '@sidebase/nuxt-auth'
  ],

  auth: {
    origin: 'http://localhost:3000',
    baseUrl: '/api/auth',
    enableGlobalAppMiddleware: true
  },

  i18n: {
    // vueI18n: './i18n.config.ts' // No longer needed, will be auto-detected
    bundle: {
      optimizeTranslationDirective: false
    }
  }
})