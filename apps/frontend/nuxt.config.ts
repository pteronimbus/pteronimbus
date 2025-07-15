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
    '@nuxtjs/i18n'
  ],

  runtimeConfig: {
    public: {
      backendUrl: process.env.BACKEND_URL || 'http://localhost:8080'
    }
  },

  i18n: {
    // vueI18n: './i18n.config.ts' // No longer needed, will be auto-detected
    bundle: {
      optimizeTranslationDirective: false
    }
  }
})