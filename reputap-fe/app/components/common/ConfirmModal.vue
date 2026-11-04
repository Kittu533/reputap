<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    title?: string
    description?: string
    confirmText?: string
    cancelText?: string
    confirmColor?: 'error' | 'primary' | 'warning' | 'neutral'
    icon?: string
    loading?: boolean
  }>(),
  {
    title: 'Confirm Action',
    description: 'Are you sure you want to proceed? This action cannot be undone.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    confirmColor: 'error',
    icon: 'i-lucide-alert-triangle',
    loading: false,
  }
)

const open = defineModel<boolean>('open', { default: false })

const emit = defineEmits<{
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()

function onConfirm() {
  emit('confirm')
}

function onCancel() {
  open.value = false
  emit('cancel')
}
</script>

<template>
  <UModal v-model:open="open" :title="title">
    <template #body>
      <div class="flex items-start gap-4 py-2">
        <div
          class="size-11 rounded-full shrink-0 flex items-center justify-center"
          :class="{
            'bg-rose-50 text-rose-600 border border-rose-100': confirmColor === 'error',
            'bg-amber-50 text-amber-600 border border-amber-100': confirmColor === 'warning',
            'bg-blue-50 text-blue-600 border border-blue-100': confirmColor === 'primary',
            'bg-slate-50 text-slate-600 border border-slate-100': confirmColor === 'neutral',
          }"
        >
          <UIcon :name="icon" class="size-6" />
        </div>
        <div class="space-y-1">
          <p class="text-sm text-slate-600 leading-relaxed">
            {{ description }}
          </p>
          <slot name="extra" />
        </div>
      </div>
    </template>

    <template #footer>
      <div class="flex items-center justify-end gap-2.5 w-full">
        <UButton
          :label="cancelText"
          color="neutral"
          variant="outline"
          :disabled="loading"
          @click="onCancel"
        />
        <UButton
          :label="confirmText"
          :color="confirmColor"
          :loading="loading"
          @click="onConfirm"
        />
      </div>
    </template>
  </UModal>
</template>
