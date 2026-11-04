<script setup lang="ts">
defineProps<{
  title?: string
  description?: string
  noPadding?: boolean
}>()
</script>

<template>
  <div class="bg-white rounded-2xl border border-slate-200/80 shadow-xs overflow-hidden">
    <div
      v-if="title || $slots.header || $slots.actions"
      class="px-6 py-4 border-b border-slate-100 flex items-center justify-between gap-4"
    >
      <div>
        <h3 v-if="title" class="font-bold text-slate-900 text-base leading-tight">{{ title }}</h3>
        <p v-if="description" class="text-xs text-slate-500 mt-1">{{ description }}</p>
        <slot name="header" />
      </div>
      <div v-if="$slots.actions" class="flex items-center gap-2">
        <slot name="actions" />
      </div>
    </div>

    <div :class="noPadding ? '' : 'p-6'">
      <slot />
    </div>

    <div v-if="$slots.footer" class="px-6 py-4 bg-slate-50/70 border-t border-slate-100 flex items-center justify-between">
      <slot name="footer" />
    </div>
  </div>
</template>
