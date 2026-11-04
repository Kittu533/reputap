<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    type?: string
    placeholder?: string
    icon?: string
    trailingIcon?: string
    disabled?: boolean
    readonly?: boolean
    error?: string
  }>(),
  {
    type: 'text',
    placeholder: '',
    disabled: false,
    readonly: false,
  }
)

const model = defineModel<string | number | null>({ default: '' })
</script>

<template>
  <div class="relative w-full">
    <UIcon
      v-if="icon"
      :name="icon"
      class="size-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
    />
    <input
      v-model="model"
      :type="type"
      :placeholder="placeholder"
      :disabled="disabled"
      :readonly="readonly"
      class="w-full py-2 text-sm bg-white border rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 transition-all disabled:bg-slate-50 disabled:text-slate-500"
      :class="[
        icon ? 'pl-9.5' : 'pl-3.5',
        trailingIcon ? 'pr-9.5' : 'pr-3.5',
        error
          ? 'border-rose-300 focus:ring-rose-500/20 focus:border-rose-500'
          : 'border-slate-200 focus:ring-blue-500/20 focus:border-blue-500'
      ]"
    >
    <UIcon
      v-if="trailingIcon"
      :name="trailingIcon"
      class="size-4 text-slate-400 absolute right-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
    />
  </div>
</template>
