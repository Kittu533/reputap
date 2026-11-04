<script setup lang="ts">
const props = withDefaults(
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

const visible = ref(true)

const config = computed(() => {
  switch (props.type) {
    case 'success':
      return {
        bg: 'bg-emerald-50/90',
        border: 'border-emerald-200',
        text: 'text-emerald-900',
        desc: 'text-emerald-700',
        icon: 'i-lucide-check-circle-2',
        iconColor: 'text-emerald-600',
      }
    case 'warning':
      return {
        bg: 'bg-amber-50/90',
        border: 'border-amber-200',
        text: 'text-amber-900',
        desc: 'text-amber-700',
        icon: 'i-lucide-alert-triangle',
        iconColor: 'text-amber-600',
      }
    case 'error':
      return {
        bg: 'bg-rose-50/90',
        border: 'border-rose-200',
        text: 'text-rose-900',
        desc: 'text-rose-700',
        icon: 'i-lucide-alert-circle',
        iconColor: 'text-rose-600',
      }
    default:
      return {
        bg: 'bg-blue-50/90',
        border: 'border-blue-200',
        text: 'text-blue-900',
        desc: 'text-blue-700',
        icon: 'i-lucide-info',
        iconColor: 'text-blue-600',
      }
  }
})
</script>

<template>
  <div
    v-if="visible"
    class="rounded-xl border p-4 shadow-xs transition-all duration-200"
    :class="[config.bg, config.border]"
  >
    <div class="flex items-start gap-3">
      <UIcon :name="config.icon" class="size-5 shrink-0 mt-0.5" :class="config.iconColor" />
      <div class="flex-1 space-y-0.5">
        <h4 v-if="title" class="text-sm font-semibold" :class="config.text">
          {{ title }}
        </h4>
        <p v-if="description" class="text-xs sm:text-sm leading-relaxed" :class="config.desc">
          {{ description }}
        </p>
        <slot />
      </div>
      <button
        v-if="dismissible"
        class="text-slate-400 hover:text-slate-600 p-1 -mr-1 rounded-md transition-colors"
        @click="visible = false"
      >
        <UIcon name="i-lucide-x" class="size-4" />
      </button>
    </div>
  </div>
</template>
