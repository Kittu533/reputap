<script setup lang="ts">
defineProps<{
  label?: string
  name?: string
  required?: boolean
  hint?: string
  error?: string
}>()
</script>

<template>
  <div class="space-y-1.5 w-full">
    <div v-if="label" class="flex items-center justify-between">
      <label :for="name" class="block text-xs font-bold uppercase tracking-wider text-slate-700">
        {{ label }}
        <span v-if="required" class="text-rose-500 ml-0.5">*</span>
      </label>
      <span v-if="hint" class="text-[11px] text-slate-400">{{ hint }}</span>
    </div>

    <slot />

    <p v-if="error" class="text-xs text-rose-500 font-medium mt-1">
      {{ error }}
    </p>
  </div>
</template>
