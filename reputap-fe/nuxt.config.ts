// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ['@nuxt/eslint', '@nuxt/ui'],
  css: ['~/assets/css/main.css'],
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  colorMode: {
    preference: 'light',
  },
  runtimeConfig: {
    public: {
      whatsappNumber: process.env.NUXT_PUBLIC_WHATSAPP_NUMBER || '6285834653796',
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8081',
    },
  },
  vite: {
    optimizeDeps: {
      include: [
        '@vue/devtools-core',
      ],
    },
  },
})
