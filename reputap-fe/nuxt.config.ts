// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  app: {
    head: {
      htmlAttrs: { lang: 'id' },
      link: [
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32.png?v=2' },
        { rel: 'icon', type: 'image/png', sizes: '48x48', href: '/favicon-48.png?v=2' },
        { rel: 'icon', type: 'image/png', sizes: '192x192', href: '/icon-192.png?v=2' },
        { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png?v=2' },
      ],
    },
  },
  modules: ['@nuxt/eslint', '@nuxt/ui'],
  css: ['~/assets/css/main.css'],
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  colorMode: {
    preference: 'light',
  },
  runtimeConfig: {
    public: {
      siteUrl: process.env.NUXT_PUBLIC_SITE_URL || 'https://reputap-alpha.vercel.app',
      whatsappNumber: process.env.NUXT_PUBLIC_WHATSAPP_NUMBER || '6285834653796',
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'https://api.kitangoding.my.id',
      redirectBase: process.env.NUXT_PUBLIC_REDIRECT_BASE || 'https://api.kitangoding.my.id/r',
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
