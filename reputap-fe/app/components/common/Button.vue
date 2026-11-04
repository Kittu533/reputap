<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger'
    size?: 'xs' | 'sm' | 'md' | 'lg'
    icon?: string
    trailingIcon?: string
    loading?: boolean
    disabled?: boolean
    to?: string
    type?: 'button' | 'submit' | 'reset'
  }>(),
  {
    variant: 'primary',
    size: 'md',
    type: 'button',
    loading: false,
    disabled: false,
  }
)
</script>

<template>
  <NuxtLink
    v-if="to"
    :to="to"
    class="inline-flex items-center justify-center font-semibold rounded-xl transition-all cursor-pointer select-none"
    :class="[
      size === 'xs' ? 'px-2.5 py-1 text-xs gap-1' :
      size === 'sm' ? 'px-3 py-1.5 text-xs gap-1.5' :
      size === 'lg' ? 'px-5 py-2.5 text-base gap-2' :
      'px-4 py-2 text-sm gap-2',
      variant === 'primary' ? 'bg-blue-600 hover:bg-blue-700 text-white shadow-xs' :
      variant === 'secondary' ? 'bg-slate-100 hover:bg-slate-200 text-slate-800' :
      variant === 'outline' ? 'border border-slate-200 bg-white hover:bg-slate-50 text-slate-700' :
      variant === 'ghost' ? 'text-slate-600 hover:bg-slate-100 hover:text-slate-900' :
      variant === 'danger' ? 'bg-rose-600 hover:bg-rose-700 text-white shadow-xs' : '',
      disabled || loading ? 'opacity-60 pointer-events-none' : ''
    ]"
  >
    <UIcon v-if="loading" name="i-lucide-loader-2" class="animate-spin shrink-0" :class="size === 'xs' ? 'size-3' : 'size-4'" />
    <UIcon v-else-if="icon" :name="icon" class="shrink-0" :class="size === 'xs' ? 'size-3.5' : 'size-4'" />
    <slot />
    <UIcon v-if="trailingIcon && !loading" :name="trailingIcon" class="shrink-0" :class="size === 'xs' ? 'size-3.5' : 'size-4'" />
  </NuxtLink>

  <button
    v-else
    :type="type"
    :disabled="disabled || loading"
    class="inline-flex items-center justify-center font-semibold rounded-xl transition-all cursor-pointer select-none"
    :class="[
      size === 'xs' ? 'px-2.5 py-1 text-xs gap-1' :
      size === 'sm' ? 'px-3 py-1.5 text-xs gap-1.5' :
      size === 'lg' ? 'px-5 py-2.5 text-base gap-2' :
      'px-4 py-2 text-sm gap-2',
      variant === 'primary' ? 'bg-blue-600 hover:bg-blue-700 text-white shadow-xs' :
      variant === 'secondary' ? 'bg-slate-100 hover:bg-slate-200 text-slate-800' :
      variant === 'outline' ? 'border border-slate-200 bg-white hover:bg-slate-50 text-slate-700' :
      variant === 'ghost' ? 'text-slate-600 hover:bg-slate-100 hover:text-slate-900' :
      variant === 'danger' ? 'bg-rose-600 hover:bg-rose-700 text-white shadow-xs' : '',
      disabled || loading ? 'opacity-60 cursor-not-allowed' : ''
    ]"
  >
    <UIcon v-if="loading" name="i-lucide-loader-2" class="animate-spin shrink-0" :class="size === 'xs' ? 'size-3' : 'size-4'" />
    <UIcon v-else-if="icon" :name="icon" class="shrink-0" :class="size === 'xs' ? 'size-3.5' : 'size-4'" />
    <slot />
    <UIcon v-if="trailingIcon && !loading" :name="trailingIcon" class="shrink-0" :class="size === 'xs' ? 'size-3.5' : 'size-4'" />
  </button>
</template>
