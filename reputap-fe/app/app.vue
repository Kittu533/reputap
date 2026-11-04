<script setup lang="ts">
const colorMode = useColorMode()

onMounted(() => {
  // Sync and ensure default is light mode
  const stored = localStorage.getItem('nuxt-color-mode')
  if (!stored || stored === 'dark' || stored === 'system') {
    colorMode.preference = 'light'
    colorMode.value = 'light'
    localStorage.setItem('nuxt-color-mode', 'light')
    document.documentElement.classList.remove('dark')
    document.documentElement.classList.add('light')
  }
})
</script>

<template>
  <UApp>
    <NuxtRouteAnnouncer />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <UToaster />
  </UApp>
</template>
