<script setup lang="ts">
withDefaults(
  defineProps<{
    title?: string
    description?: string
  }>(),
  {
    title: '',
    description: '',
  }
)

const open = defineModel<boolean>('open', { default: false })
</script>

<template>
  <UModal v-model:open="open" :title="title" :description="description">
    <template #body>
      <div class="py-2">
        <slot />
      </div>
    </template>
    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </UModal>
</template>
