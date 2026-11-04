<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { Order } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()
const toast = useToast()

const UButton = resolveComponent('UButton')

const search = ref('')
const statusFilter = ref<string | undefined>(undefined)
const paymentFilter = ref<string | undefined>(undefined)
const page = ref(1)
const perPage = ref(20)
const data = ref<Order[]>([])
const total = ref(0)
const loading = ref(false)

const orderStatusOptions = [
  { label: 'Pending', value: 'pending' },
  { label: 'Confirmed', value: 'confirmed' },
  { label: 'Processing', value: 'processing' },
  { label: 'Production', value: 'production' },
  { label: 'QC', value: 'qc' },
  { label: 'Ready to Ship', value: 'ready_to_ship' },
  { label: 'Shipped', value: 'shipped' },
  { label: 'Completed', value: 'completed' },
  { label: 'Cancelled', value: 'cancelled' },
  { label: 'On Hold', value: 'on_hold' },
]

const paymentStatusOptions = [
  { label: 'Unpaid', value: 'unpaid' },
  { label: 'Partial', value: 'partial' },
  { label: 'Paid', value: 'paid' },
  { label: 'Refunded', value: 'refunded' },
]

const OrderStatusBadge = resolveComponent('OrdersOrderStatusBadge')

function formatRupiah(amount: number | null): string {
  if (amount == null) return '—'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount)
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

const columns: TableColumn<Order>[] = [
  { accessorKey: 'order_number', header: 'Order #' },
  {
    accessorKey: 'customer_id',
    header: 'Customer',
    cell: ({ row }) => `#${row.original.customer_id}`,
  },
  {
    accessorKey: 'order_status',
    header: 'Status',
    cell: ({ row }) => h(OrderStatusBadge, { status: row.original.order_status }),
  },
  {
    accessorKey: 'payment_status',
    header: 'Payment',
    cell: ({ row }) => h(OrderStatusBadge, { status: row.original.payment_status, type: 'payment' }),
  },
  {
    accessorKey: 'total',
    header: 'Total',
    cell: ({ row }) => formatRupiah(row.original.total),
  },
  {
    accessorKey: 'created_at',
    header: 'Created',
    cell: ({ row }) => formatDate(row.original.created_at),
  },
  {
    id: 'actions',
    header: '',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => {
      return h('div', { class: 'flex justify-end gap-1' }, [
        h(UButton, {
          icon: 'i-lucide-eye',
          color: 'neutral',
          variant: 'ghost',
          size: 'xs',
          to: `/orders/${row.original.id}`,
        }),
      ])
    },
  },
]

async function fetchData() {
  loading.value = true
  const query: Record<string, any> = {
    page: page.value,
    per_page: perPage.value,
  }
  if (search.value) query.search = search.value
  if (statusFilter.value) query.status = statusFilter.value
  if (paymentFilter.value) query.payment_status = paymentFilter.value

  const response = await apiFetch<Order[]>('/api/v1/orders', { query })
  if (response.success) {
    data.value = response.data || []
    total.value = response.meta?.total || 0
  }
  loading.value = false
}

let searchTimeout: ReturnType<typeof setTimeout>
watch(search, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => { page.value = 1; fetchData() }, 300)
})
watch(statusFilter, () => { page.value = 1; fetchData() })
watch(paymentFilter, () => { page.value = 1; fetchData() })
watch(page, () => fetchData())
onMounted(() => fetchData())
</script>

<template>
  <div class="space-y-6">
    <CommonPageHeader
      title="Orders"
      description="Track NFC review stand orders, fulfillment lifecycle, and payment status."
      :breadcrumbs="[{ label: 'Orders' }]"
    >
      <template #actions>
        <CommonButton to="/orders/new" icon="i-lucide-plus">
          Create Order
        </CommonButton>
      </template>
    </CommonPageHeader>

    <CommonDataTable
      v-model:search="search"
      v-model:page="page"
      :data="data"
      :columns="columns"
      :total="total"
      :per-page="perPage"
      :loading="loading"
      search-placeholder="Search order number..."
      empty-title="No orders found"
      empty-description="There are no orders matching your criteria."
      empty-action-label="Create Order"
      @empty-action="navigateTo('/orders/new')"
    >
      <template #filters>
        <USelect
          v-model="statusFilter"
          :items="orderStatusOptions"
          placeholder="All Status"
          class="w-44"
        />
        <USelect
          v-model="paymentFilter"
          :items="paymentStatusOptions"
          placeholder="All Payment"
          class="w-40"
        />
      </template>
    </CommonDataTable>
  </div>
</template>
