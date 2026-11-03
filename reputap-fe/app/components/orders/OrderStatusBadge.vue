<script setup lang="ts">
const props = defineProps<{
  status: string
  type?: 'order' | 'payment'
}>()

const orderStatusColor: Record<string, string> = {
  pending: 'warning',
  confirmed: 'info',
  processing: 'primary',
  production: 'primary',
  qc: 'secondary',
  ready_to_ship: 'info',
  shipped: 'primary',
  completed: 'success',
  cancelled: 'error',
  on_hold: 'neutral',
}

const paymentStatusColor: Record<string, string> = {
  unpaid: 'error',
  partial: 'warning',
  paid: 'success',
  refunded: 'neutral',
}

const color = computed(() => {
  const map = props.type === 'payment' ? paymentStatusColor : orderStatusColor
  return map[props.status] || 'neutral'
})

const label = computed(() => props.status.replace(/_/g, ' '))
</script>

<template>
  <UBadge variant="subtle" :color="color" class="capitalize">
    {{ label }}
  </UBadge>
</template>
