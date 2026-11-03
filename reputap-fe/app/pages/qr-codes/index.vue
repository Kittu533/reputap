<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { QRCodeDetail } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()
const config = useRuntimeConfig()
const { accessToken } = useAuth()

const UButton = resolveComponent('UButton')
const UBadge = resolveComponent('UBadge')

const search = ref('')
const statusFilter = ref<string | undefined>(undefined)
const page = ref(1)
const perPage = ref(20)
const data = ref<QRCodeDetail[]>([])
const total = ref(0)
const loading = ref(false)

const statusOptions = [
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
  { label: 'Suspended', value: 'suspended' },
  { label: 'Retired', value: 'retired' },
]

const statusColor: Record<string, string> = {
  active: 'success',
  inactive: 'neutral',
  suspended: 'warning',
  retired: 'error',
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

function truncateUrl(url: string, maxLen = 40): string {
  if (!url) return '—'
  return url.length > maxLen ? `${url.substring(0, maxLen)}…` : url
}

const columns: TableColumn<QRCodeDetail>[] = [
  { accessorKey: 'code', header: 'Code' },
  {
    accessorKey: 'store_id',
    header: 'Store ID',
    cell: ({ row }) => `#${row.original.store_id}`,
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => h(UBadge, {
      label: row.original.status,
      color: statusColor[row.original.status] || 'neutral',
      variant: 'subtle',
      size: 'xs',
      class: 'capitalize',
    }),
  },
  {
    accessorKey: 'redirect_url',
    header: 'Redirect URL',
    cell: ({ row }) => {
      const url = row.original.redirect_url
      if (!url) return '—'
      return h('a', {
        href: url,
        target: '_blank',
        rel: 'noopener noreferrer',
        class: 'text-primary hover:underline text-xs',
      }, truncateUrl(url))
    },
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
          to: `/qr-codes/${row.original.id}`,
        }),
        h(UButton, {
          icon: 'i-lucide-download',
          color: 'neutral',
          variant: 'ghost',
          size: 'xs',
          onClick: () => downloadQR(row.original),
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

  const response = await apiFetch<QRCodeDetail[]>('/api/v1/qr-codes', { query })
  if (response.success) {
    data.value = response.data || []
    total.value = response.meta?.total || 0
  }
  loading.value = false
}

async function downloadQR(qr: QRCodeDetail) {
  const apiBase = config.public.apiBase as string
  const response = await fetch(`${apiBase}/api/v1/qr-codes/${qr.id}/download?size=512`, {
    headers: { Authorization: `Bearer ${accessToken.value}` },
  })
  const blob = await response.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${qr.code}.png`
  a.click()
  URL.revokeObjectURL(url)
}

let searchTimeout: ReturnType<typeof setTimeout>
watch(search, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => { page.value = 1; fetchData() }, 300)
})
watch(statusFilter, () => { page.value = 1; fetchData() })
watch(page, () => fetchData())
onMounted(() => fetchData())
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">QR Codes</h1>
    </div>

    <div class="flex items-center gap-4 mb-4">
      <UInput
        v-model="search"
        placeholder="Search by code..."
        icon="i-lucide-search"
        class="max-w-sm"
      />
      <USelect
        v-model="statusFilter"
        :items="statusOptions"
        placeholder="All Status"
        class="w-44"
      />
    </div>

    <UTable :data="data" :columns="columns" :loading="loading" class="w-full" />

    <div v-if="total > perPage" class="flex justify-end mt-4">
      <UPagination v-model:page="page" :items-per-page="perPage" :total="total" />
    </div>
  </div>
</template>
