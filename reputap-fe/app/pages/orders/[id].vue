<script setup lang="ts">
import type { OrderDetail, StatusHistory, ProductInstance } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const route = useRoute()
const { apiFetch } = useApi()
const toast = useToast()

const orderId = route.params.id as string
const order = ref<OrderDetail | null>(null)
const history = ref<StatusHistory[]>([])
const loading = ref(true)
const statusLoading = ref(false)
const paymentLoading = ref(false)

const allowedTransitions: Record<string, string[]> = {
  pending: ['confirmed', 'cancelled', 'on_hold'],
  confirmed: ['processing', 'cancelled', 'on_hold'],
  processing: ['production', 'cancelled', 'on_hold'],
  production: ['qc', 'on_hold'],
  qc: ['ready_to_ship', 'production', 'on_hold'],
  ready_to_ship: ['shipped', 'on_hold'],
  shipped: ['completed'],
  on_hold: ['pending', 'confirmed', 'processing', 'production', 'qc', 'ready_to_ship'],
}

const transitionColors: Record<string, string> = {
  confirmed: 'info',
  processing: 'primary',
  production: 'primary',
  qc: 'secondary',
  ready_to_ship: 'info',
  shipped: 'primary',
  completed: 'success',
  cancelled: 'error',
  on_hold: 'neutral',
  pending: 'warning',
}

const paymentStatusOptions = [
  { label: 'Unpaid', value: 'unpaid' },
  { label: 'Partial', value: 'partial' },
  { label: 'Paid', value: 'paid' },
  { label: 'Refunded', value: 'refunded' },
]

const paymentStatusModel = ref<string | undefined>(undefined)

function formatRupiah(amount: number | null): string {
  if (amount == null) return '—'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount)
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function formatLabel(status: string): string {
  return status.replace(/_/g, ' ')
}

async function fetchOrder() {
  loading.value = true
  const response = await apiFetch<OrderDetail>(`/api/v1/orders/${orderId}`)
  if (response.success) {
    order.value = response.data
    paymentStatusModel.value = response.data.payment_status
  } else {
    toast.add({ title: 'Error', description: response.error?.message || 'Failed to load order', color: 'error' })
  }
  loading.value = false
}

async function fetchHistory() {
  const response = await apiFetch<StatusHistory[]>(`/api/v1/orders/${orderId}/history`)
  if (response.success) {
    history.value = response.data || []
  }
}

async function updateStatus(newStatus: string) {
  const notes = prompt(`Notes for status change to "${formatLabel(newStatus)}" (optional):`)
  if (notes === null) return

  statusLoading.value = true
  const response = await apiFetch(`/api/v1/orders/${orderId}/status`, {
    method: 'PATCH',
    body: { status: newStatus, notes: notes || undefined },
  })
  statusLoading.value = false

  if (response.success) {
    toast.add({ title: 'Updated', description: `Status changed to ${formatLabel(newStatus)}`, color: 'success' })
    fetchOrder()
    fetchHistory()
  } else {
    toast.add({ title: 'Error', description: response.error?.message || 'Update failed', color: 'error' })
  }
}

async function updatePaymentStatus() {
  if (!order.value || paymentStatusModel.value === order.value.payment_status) return

  paymentLoading.value = true
  const response = await apiFetch(`/api/v1/orders/${orderId}/payment-status`, {
    method: 'PATCH',
    body: { payment_status: paymentStatusModel.value },
  })
  paymentLoading.value = false

  if (response.success) {
    toast.add({ title: 'Updated', description: `Payment status changed to ${formatLabel(paymentStatusModel.value)}`, color: 'success' })
    fetchOrder()
  } else {
    toast.add({ title: 'Error', description: response.error?.message || 'Update failed', color: 'error' })
    paymentStatusModel.value = order.value.payment_status
  }
}

const availableTransitions = computed(() => {
  if (!order.value) return []
  return allowedTransitions[order.value.order_status] || []
})

const instances = ref<ProductInstance[]>([])
const instancesLoading = ref(false)
const generateLoading = ref(false)

const nfcStatusColor: Record<string, string> = {
  not_configured: 'neutral',
  configured: 'info',
  tested: 'success',
  failed: 'error',
  replaced: 'warning',
}

const prodStatusColor: Record<string, string> = {
  pending: 'neutral',
  in_progress: 'info',
  completed: 'success',
  failed: 'error',
}

const qcStatusColor: Record<string, string> = {
  pending: 'neutral',
  passed: 'success',
  failed: 'error',
  rework: 'warning',
}

async function fetchInstances() {
  instancesLoading.value = true
  const response = await apiFetch<ProductInstance[]>(`/api/v1/orders/${orderId}/instances`)
  if (response.success) {
    instances.value = response.data || []
  }
  instancesLoading.value = false
}

async function generateInstances() {
  generateLoading.value = true
  const response = await apiFetch<ProductInstance[]>(`/api/v1/orders/${orderId}/instances/generate`, {
    method: 'POST',
  })
  generateLoading.value = false

  if (response.success) {
    instances.value = response.data || []
    toast.add({ title: 'Generated', description: `${instances.value.length} product instances created`, color: 'success' })
  } else {
    if (response.error?.code === 'ALREADY_GENERATED') {
      toast.add({ title: 'Info', description: 'Instances already generated for this order', color: 'info' })
      fetchInstances()
    } else {
      toast.add({ title: 'Error', description: response.error?.message || 'Generation failed', color: 'error' })
    }
  }
}

onMounted(() => {
  fetchOrder()
  fetchHistory()
  fetchInstances()
})
</script>

<template>
  <div>
    <div class="flex items-center gap-2 mb-6">
      <UButton icon="i-lucide-arrow-left" color="neutral" variant="ghost" to="/orders" />
      <h1 class="text-2xl font-bold">
        {{ order?.order_number || 'Order Detail' }}
      </h1>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <UIcon name="i-lucide-loader-2" class="animate-spin text-2xl" />
    </div>

    <div v-else-if="order" class="space-y-6">
      <!-- Order Header -->
      <div class="rounded-lg border border-(--ui-border) p-4">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <p class="text-sm text-(--ui-text-muted)">Order Number</p>
            <p class="font-semibold">{{ order.order_number }}</p>
          </div>
          <div>
            <p class="text-sm text-(--ui-text-muted)">Status</p>
            <OrdersOrderStatusBadge :status="order.order_status" />
          </div>
          <div>
            <p class="text-sm text-(--ui-text-muted)">Payment</p>
            <OrdersOrderStatusBadge :status="order.payment_status" type="payment" />
          </div>
          <div>
            <p class="text-sm text-(--ui-text-muted)">Created</p>
            <p class="text-sm">{{ formatDate(order.created_at) }}</p>
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4 pt-4 border-t border-(--ui-border)">
          <div>
            <p class="text-sm text-(--ui-text-muted)">Customer ID</p>
            <p class="font-mono text-sm">#{{ order.customer_id }}</p>
          </div>
          <div>
            <p class="text-sm text-(--ui-text-muted)">Store ID</p>
            <p class="font-mono text-sm">#{{ order.store_id }}</p>
          </div>
          <div>
            <p class="text-sm text-(--ui-text-muted)">Created By</p>
            <p class="font-mono text-sm">#{{ order.created_by }}</p>
          </div>
          <div>
            <p class="text-sm text-(--ui-text-muted)">Updated</p>
            <p class="text-sm">{{ formatDate(order.updated_at) }}</p>
          </div>
        </div>
      </div>

      <!-- Items Table -->
      <div class="rounded-lg border border-(--ui-border) p-4">
        <h2 class="text-lg font-semibold mb-3">Items</h2>
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-(--ui-border) text-left text-(--ui-text-muted)">
              <th class="pb-2">Product ID</th>
              <th class="pb-2 text-right">Qty</th>
              <th class="pb-2 text-right">Unit Price</th>
              <th class="pb-2 text-right">Subtotal</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="item in order.items"
              :key="item.id"
              class="border-b border-(--ui-border) last:border-b-0"
            >
              <td class="py-2 font-mono">#{{ item.product_id }}</td>
              <td class="py-2 text-right">{{ item.quantity }}</td>
              <td class="py-2 text-right">{{ formatRupiah(item.unit_price) }}</td>
              <td class="py-2 text-right">{{ formatRupiah(item.subtotal) }}</td>
            </tr>
          </tbody>
        </table>

        <div class="mt-3 pt-3 border-t border-(--ui-border) space-y-1 text-sm">
          <div class="flex justify-between">
            <span class="text-(--ui-text-muted)">Subtotal</span>
            <span>{{ formatRupiah(order.subtotal) }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-(--ui-text-muted)">Discount</span>
            <span class="text-(--ui-text-muted)">- {{ formatRupiah(order.discount) }}</span>
          </div>
          <div class="flex justify-between font-semibold text-base pt-1 border-t border-(--ui-border)">
            <span>Total</span>
            <span>{{ formatRupiah(order.total) }}</span>
          </div>
        </div>
      </div>

      <!-- Product Instances -->
      <div class="rounded-lg border border-(--ui-border) p-4">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2">
            <h2 class="text-lg font-semibold">Product Instances</h2>
            <UBadge v-if="instances.length > 0" :label="String(instances.length)" color="primary" variant="subtle" size="sm" />
          </div>
          <UButton
            v-if="instances.length === 0"
            label="Generate Instances"
            icon="i-lucide-wand-2"
            size="sm"
            :loading="generateLoading"
            @click="generateInstances"
          />
        </div>

        <div v-if="instancesLoading" class="flex justify-center py-4">
          <UIcon name="i-lucide-loader-2" class="animate-spin" />
        </div>

        <div v-else-if="instances.length === 0" class="text-sm text-(--ui-text-muted)">
          No instances generated yet. Click "Generate Instances" to create them.
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-(--ui-border) text-left text-(--ui-text-muted)">
                <th class="pb-2">Serial Number</th>
                <th class="pb-2">NFC Status</th>
                <th class="pb-2">Production</th>
                <th class="pb-2">QC</th>
                <th class="pb-2">QR Code</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="inst in instances" :key="inst.id" class="border-b border-(--ui-border) last:border-b-0">
                <td class="py-2 font-mono text-xs">{{ inst.serial_number }}</td>
                <td class="py-2">
                  <UBadge :label="inst.nfc_status.replace(/_/g, ' ')" :color="nfcStatusColor[inst.nfc_status] || 'neutral'" variant="subtle" size="xs" class="capitalize" />
                </td>
                <td class="py-2">
                  <UBadge :label="inst.production_status.replace(/_/g, ' ')" :color="prodStatusColor[inst.production_status] || 'neutral'" variant="subtle" size="xs" class="capitalize" />
                </td>
                <td class="py-2">
                  <UBadge :label="inst.qc_status" :color="qcStatusColor[inst.qc_status] || 'neutral'" variant="subtle" size="xs" class="capitalize" />
                </td>
                <td class="py-2">
                  <NuxtLink :to="`/qr-codes`" class="text-primary hover:underline text-xs font-mono">
                    {{ inst.qr_code }}
                  </NuxtLink>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Notes -->
      <div v-if="order.notes" class="rounded-lg border border-(--ui-border) p-4">
        <h2 class="text-lg font-semibold mb-2">Notes</h2>
        <p class="text-sm text-(--ui-text-muted) whitespace-pre-wrap">{{ order.notes }}</p>
      </div>

      <!-- Status Actions -->
      <div class="rounded-lg border border-(--ui-border) p-4">
        <h2 class="text-lg font-semibold mb-3">Actions</h2>

        <div class="space-y-4">
          <div v-if="availableTransitions.length > 0">
            <p class="text-sm text-(--ui-text-muted) mb-2">Transition order status to:</p>
            <div class="flex flex-wrap gap-2">
              <UButton
                v-for="status in availableTransitions"
                :key="status"
                :label="formatLabel(status)"
                :color="transitionColors[status] || 'neutral'"
                variant="soft"
                size="sm"
                class="capitalize"
                :loading="statusLoading"
                @click="updateStatus(status)"
              />
            </div>
          </div>
          <div v-else>
            <p class="text-sm text-(--ui-text-muted)">No status transitions available.</p>
          </div>

          <div class="pt-3 border-t border-(--ui-border)">
            <p class="text-sm text-(--ui-text-muted) mb-2">Update payment status:</p>
            <div class="flex items-center gap-2">
              <USelect
                v-model="paymentStatusModel"
                :items="paymentStatusOptions"
                class="w-40"
              />
              <UButton
                label="Update"
                size="sm"
                :loading="paymentLoading"
                :disabled="paymentStatusModel === order.payment_status"
                @click="updatePaymentStatus"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Status History Timeline -->
      <div class="rounded-lg border border-(--ui-border) p-4">
        <h2 class="text-lg font-semibold mb-3">Status History</h2>

        <div v-if="history.length === 0" class="text-sm text-(--ui-text-muted)">
          No status history available.
        </div>

        <div v-else class="relative space-y-0">
          <div
            v-for="(entry, index) in history"
            :key="entry.id"
            class="flex gap-3 pb-4 last:pb-0"
          >
            <div class="flex flex-col items-center">
              <div class="w-2.5 h-2.5 rounded-full bg-(--ui-primary) mt-1.5 shrink-0" />
              <div v-if="index < history.length - 1" class="w-px flex-1 bg-(--ui-border)" />
            </div>
            <div class="pb-1">
              <div class="flex items-center gap-2">
                <OrdersOrderStatusBadge v-if="entry.from_status" :status="entry.from_status" />
                <UIcon v-if="entry.from_status" name="i-lucide-arrow-right" class="text-(--ui-text-muted) text-xs" />
                <OrdersOrderStatusBadge :status="entry.to_status" />
              </div>
              <p v-if="entry.notes" class="text-sm text-(--ui-text-muted) mt-1">{{ entry.notes }}</p>
              <p class="text-xs text-(--ui-text-dimmed) mt-0.5">
                {{ formatDate(entry.created_at) }}
                · User #{{ entry.changed_by }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
