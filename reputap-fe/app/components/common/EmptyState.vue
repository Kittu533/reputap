<script setup lang="ts">
withDefaults(
  defineProps<{
    title?: string
    description?: string
    icon?: string
    actionLabel?: string
    actionIcon?: string
  }>(),
  {
    title: 'No items found',
    description: 'Get started by creating a new record or adjust your search filters.',
    icon: 'i-lucide-inbox',
  }
)

const emit = defineEmits<{
  (e: 'action'): void
}>()
</script>

<template>
  <div class="py-12 px-4 text-center flex flex-col items-center justify-center">
    <div class="size-12 rounded-xl bg-slate-100/80 text-slate-400 flex items-center justify-center mb-3.5 border border-slate-200/60 shadow-xs">
      <UIcon :name="icon" class="size-6 text-slate-500" />
    </div>
    <h3 class="text-base font-semibold text-slate-900 mb-1">
      {{ title }}
    </h3>
    <p class="text-sm text-slate-500 max-w-sm mb-5 leading-relaxed">
      {{ description }}
    </p>
    <div v-if="$slots.action || actionLabel" class="flex items-center gap-2">
      <slot name="action">
        <UButton
          v-if="actionLabel"
          :label="actionLabel"
          :icon="actionIcon || 'i-lucide-plus'"
          color="primary"
          @click="emit('action')"
        />
      </slot>
    </div>
  </div>
</template>
