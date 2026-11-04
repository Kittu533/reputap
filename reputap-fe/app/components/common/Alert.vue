<script setup lang="ts">
withDefaults(
  defineProps<{
    type?: 'info' | 'success' | 'warning' | 'error'
    title?: string
    description?: string
    dismissible?: boolean
  }>(),
  {
    type: 'info',
    dismissible: false,
  }
)
</script>

<template>
  <CommonAlertBanner
    :type="type"
    :title="title"
    :description="description"
    :dismissible="dismissible"
  >
    <slot />
  </CommonAlertBanner>
</template>
