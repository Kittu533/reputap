export default defineNuxtPlugin(() => {
  // Ensure default is light mode and clear stale 'dark' / 'system' settings
  if (import.meta.client) {
    const stored = localStorage.getItem('nuxt-color-mode')
    if (!stored || stored === 'dark' || stored === 'system') {
      localStorage.setItem('nuxt-color-mode', 'light')
      document.documentElement.classList.remove('dark')
      document.documentElement.classList.add('light')
    }
  }
})
