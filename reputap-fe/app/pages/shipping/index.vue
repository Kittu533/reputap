<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { Shipment } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()
const toast = useToast()

const UButton = resolveComponent('UButton')
const UBadge = resolveComponent('UBadge')

const statusFilter = ref<string | undefined>(undefined)
const page = ref(1)
const perPage = ref(20)
const data = ref<Shipment[]>([])
const total = ref(0)
const loading = ref(false)

const statusOptions = [
  { label: 'Waiting', value: 'waiting' },
  { label: 'Packed', value: 'packed' },
  { label: 'Shipped', value: 'shipped' },
  { label: 'Delivered', value: 'delivered' },
  { label: 'Failed', value: 'failed' },
]

const statusColor: Record<string, string> = {
  waiting: 'neutral',
  packed: 'info',
  shipped: 'warning',
  delivered: 'success',
  failed: 'error',
}

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function formatStatus(status: string): string {
  return status.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

const columns: TableColumn<Shipment>[] = [
  { accessorKey: 'id', header: 'ID' },
  {
    accessorKey: 'order_id',
    header: 'Order ID',
    cell: ({ row }) => `#${row.original.order_id}`,
  },
  { accessorKey: 'recipient_name', header: 'Recipient' },
  {
    accessorKey: 'courier',
    header: 'Courier',
    cell: ({ row }) => row.original.courier || '—',
  },
  {
    accessorKey: 'tracking_number',
    header: 'Tracking Number',
    cell: ({ row }) => row.original.tracking_number || '—',
  },
  {
    accessorKey: 'shipping_status',
    header: 'Status',
    cell: ({ row }) => {
      const status = row.original.shipping_status
      return h(UBadge, { variant: 'subtle', color: statusColor[status] || 'neutral', class: 'capitalize' }, () => formatStatus(status))
    },
  },
  {
    accessorKey: 'shipped_at',
    header: 'Shipped At',
    cell: ({ row }) => formatDate(row.original.shipped_at),
  },
  {
    accessorKey: 'delivered_at',
    header: 'Delivered At',
    cell: ({ row }) => formatDate(row.original.delivered_at),
  },
  {
    id: 'actions',
    header: '',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => {
      return h('div', { class: 'flex justify-end gap-1' }, [
        h(UButton, {
          icon: 'i-lucide-refresh-cw',
          color: 'neutral',
          variant: 'ghost',
          size: 'xs',
          onClick: (e: Event) => { e.stopPropagation(); openUpdateStatus(row.original) },
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
  if (statusFilter.value) query.status = statusFilter.value

  const response = await apiFetch<Shipment[]>('/api/v1/shipments', { query })
  if (response.success) {
    data.value = response.data || []
    total.value = response.meta?.total || 0
  }
  loading.value = false
}

watch(statusFilter, () => { page.value = 1; fetchData() })
watch(page, () => fetchData())
onMounted(() => fetchData())

// Create slideover
const createOpen = ref(false)
const submitting = ref(false)
const form = reactive({
  order_id: null as number | null,
  recipient_name: '',
  recipient_address: '',
  recipient_phone: '',
  courier: '',
  tracking_number: '',
})

function resetForm() {
  form.order_id = null
  form.recipient_name = ''
  form.recipient_address = ''
  form.recipient_phone = ''
  form.courier = ''
  form.tracking_number = ''
}

function openCreate() {
  resetForm()
  createOpen.value = true
}

async function submitCreate() {
  if (!form.order_id || !form.recipient_name || !form.recipient_address) {
    toast.add({ title: 'Please fill required fields', color: 'warning' })
    return
  }

  submitting.value = true
  const response = await apiFetch('/api/v1/shipments', {
    method: 'POST',
    body: {
      order_id: form.order_id,
      recipient_name: form.recipient_name,
      recipient_address: form.recipient_address,
      recipient_phone: form.recipient_phone || undefined,
      courier: form.courier || undefined,
      tracking_number: form.tracking_number || undefined,
    },
  })

  if (response.success) {
    toast.add({ title: 'Shipment created', color: 'success' })
    createOpen.value = false
    fetchData()
  }
  else {
    toast.add({ title: response.error?.message || 'Failed to create shipment', color: 'error' })
  }
  submitting.value = false
}

// Update status slideover
const updateStatusOpen = ref(false)
const updatingShipment = ref<Shipment | null>(null)
const statusSubmitting = ref(false)
const statusForm = reactive({
  status: undefined as string | undefined,
  tracking_number: '',
  courier: '',
})

function openUpdateStatus(shipment: Shipment) {
  updatingShipment.value = shipment
  statusForm.status = shipment.shipping_status
  statusForm.tracking_number = shipment.tracking_number || ''
  statusForm.courier = shipment.courier || ''
  updateStatusOpen.value = true
}

async function submitUpdateStatus() {
  if (!updatingShipment.value || !statusForm.status) {
    toast.add({ title: 'Please select a status', color: 'warning' })
    return
  }

  statusSubmitting.value = true
  const response = await apiFetch(`/api/v1/shipments/${updatingShipment.value.id}/status`, {
    method: 'PATCH',
    body: {
      status: statusForm.status,
      tracking_number: statusForm.tracking_number || undefined,
      courier: statusForm.courier || undefined,
    },
  })

  if (response.success) {
    toast.add({ title: 'Status updated', color: 'success' })
    updateStatusOpen.value = false
    fetchData()
  }
  else {
    toast.add({ title: response.error?.message || 'Failed to update status', color: 'error' })
  }
  statusSubmitting.value = false
}
</script>

<template>
  <div class="space-y-6">
    <CommonPageHeader
      title="Shipping"
      description="Manage parcel dispatches, courier tracking, and order fulfillment delivery."
      :breadcrumbs="[{ label: 'Shipping' }]"
    >
      <template #actions>
        <CommonButton icon="i-lucide-plus" @click="openCreate">
          Create Shipment
        </CommonButton>
      </template>
    </CommonPageHeader>

    <CommonDataTable
      v-model:page="page"
      :data="data"
      :columns="columns"
      :total="total"
      :per-page="perPage"
      :loading="loading"
      :show-search="false"
      empty-title="No shipments found"
      empty-description="There are currently no shipment dispatches matching your criteria."
      empty-action-label="Create Shipment"
      @empty-action="openCreate"
    >
      <template #filters>
        <USelect
          v-model="statusFilter"
          :items="statusOptions"
          placeholder="All Status"
          class="w-44"
        />
      </template>
    </CommonDataTable>

    <!-- Create Shipment Modal -->
    <CommonModal
      v-model:open="createOpen"
      title="Create Shipment"
      description="Register dispatch details and assign tracking to an order."
    >
      <div class="space-y-4">
        <CommonFormField label="Order ID" required>
          <CommonInput
            v-model.number="form.order_id"
            type="number"
            placeholder="Enter order ID"
          />
        </CommonFormField>

        <CommonFormField label="Recipient Name" required>
          <CommonInput
            v-model="form.recipient_name"
            placeholder="Enter recipient name"
          />
        </CommonFormField>

        <CommonFormField label="Recipient Address" required>
          <textarea
            v-model="form.recipient_address"
            rows="3"
            placeholder="Enter complete shipping address"
            class="w-full p-3 text-sm bg-white border border-slate-200 rounded-xl text-slate-900 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500"
          />
        </CommonFormField>

        <CommonFormField label="Recipient Phone">
          <CommonInput
            v-model="form.recipient_phone"
            placeholder="Enter phone number"
          />
        </CommonFormField>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <CommonFormField label="Courier">
            <CommonInput
              v-model="form.courier"
              placeholder="e.g. JNE, J&T, SiCepat"
            />
          </CommonFormField>

          <CommonFormField label="Tracking Number">
            <CommonInput
              v-model="form.tracking_number"
              placeholder="Enter tracking number"
            />
          </CommonFormField>
        </div>
      </div>

      <template #footer>
        <div class="flex items-center justify-end gap-2 w-full">
          <CommonButton variant="outline" @click="createOpen = false">
            Cancel
          </CommonButton>
          <CommonButton :loading="submitting" @click="submitCreate">
            Create Shipment
          </CommonButton>
        </div>
      </template>
    </CommonModal>

    <!-- Update Status Modal -->
    <CommonModal
      v-model:open="updateStatusOpen"
      title="Update Shipment Status"
      :description="updatingShipment ? `Shipment #${updatingShipment.id} — Order #${updatingShipment.order_id}` : ''"
    >
      <div v-if="updatingShipment" class="space-y-4">
        <div class="grid grid-cols-2 gap-3 bg-slate-50 p-4 rounded-xl border border-slate-200/60">
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Recipient</p>
            <p class="text-sm font-semibold text-slate-900">{{ updatingShipment.recipient_name }}</p>
          </div>
          <div>
            <p class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Current Status</p>
            <UBadge :color="statusColor[updatingShipment.shipping_status] || 'neutral'" variant="subtle" class="capitalize mt-0.5">
              {{ formatStatus(updatingShipment.shipping_status) }}
            </UBadge>
          </div>
        </div>

        <CommonFormField label="New Status" required>
          <USelect
            v-model="statusForm.status"
            :items="statusOptions"
            placeholder="Select status"
            class="w-full"
          />
        </CommonFormField>

        <CommonFormField label="Tracking Number">
          <CommonInput
            v-model="statusForm.tracking_number"
            placeholder="Enter tracking number"
          />
        </CommonFormField>

        <CommonFormField label="Courier">
          <CommonInput
            v-model="statusForm.courier"
            placeholder="Enter courier name"
          />
        </CommonFormField>
      </div>

      <template #footer>
        <div class="flex items-center justify-end gap-2 w-full">
          <CommonButton variant="outline" @click="updateStatusOpen = false">
            Cancel
          </CommonButton>
          <CommonButton :loading="statusSubmitting" @click="submitUpdateStatus">
            Update Status
          </CommonButton>
        </div>
      </template>
    </CommonModal>
  </div>
</template>
