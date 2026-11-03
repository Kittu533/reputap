<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { Product } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()
const toast = useToast()

const UButton = resolveComponent('UButton')
const UBadge = resolveComponent('UBadge')

const search = ref('')
const statusFilter = ref<string | undefined>(undefined)
const page = ref(1)
const perPage = ref(20)
const data = ref<Product[]>([])
const total = ref(0)
const loading = ref(false)

const statusOptions = [
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
  { label: 'Discontinued', value: 'discontinued' },
]

const statusColor: Record<string, string> = {
  active: 'success',
  inactive: 'neutral',
  discontinued: 'error',
}

function formatRupiah(amount: number | null): string {
  if (amount == null) return '—'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount)
}

function formatDimensions(product: Product): string {
  const { width_mm, height_mm, thickness_mm } = product
  if (width_mm == null && height_mm == null && thickness_mm == null) return '—'
  return `${width_mm ?? '—'} × ${height_mm ?? '—'} × ${thickness_mm ?? '—'} mm`
}

const columns: TableColumn<Product>[] = [
  { accessorKey: 'sku', header: 'SKU' },
  { accessorKey: 'name', header: 'Name' },
  {
    accessorKey: 'material',
    header: 'Material',
    cell: ({ row }) => row.original.material || '—',
  },
  {
    id: 'dimensions',
    header: 'Dimensions',
    cell: ({ row }) => formatDimensions(row.original),
  },
  {
    accessorKey: 'selling_price',
    header: 'Selling Price',
    cell: ({ row }) => formatRupiah(row.original.selling_price),
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => {
      const status = row.getValue('status') as string
      return h(UBadge, { variant: 'subtle', color: statusColor[status] || 'neutral', class: 'capitalize' }, () => status)
    },
  },
  {
    id: 'actions',
    header: '',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => {
      return h('div', { class: 'flex justify-end gap-1' }, [
        h(UButton, {
          icon: 'i-lucide-pencil',
          color: 'neutral',
          variant: 'ghost',
          size: 'xs',
          onClick: (e: Event) => { e.stopPropagation(); openEdit(row.original) },
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

  const response = await apiFetch<Product[]>('/api/v1/products', { query })
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
watch(page, () => fetchData())
onMounted(() => fetchData())

const slideoverOpen = ref(false)
const editingProduct = ref<Product | null>(null)

function openCreate() {
  editingProduct.value = null
  slideoverOpen.value = true
}

function openEdit(product: Product) {
  editingProduct.value = product
  slideoverOpen.value = true
}

function onSaved() {
  slideoverOpen.value = false
  fetchData()
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">Products</h1>
      <UButton icon="i-lucide-plus" label="Add Product" @click="openCreate" />
    </div>

    <div class="flex items-center gap-4 mb-4">
      <UInput
        v-model="search"
        placeholder="Search products..."
        icon="i-lucide-search"
        class="max-w-sm"
      />
      <USelect
        v-model="statusFilter"
        :items="statusOptions"
        placeholder="All Status"
        class="w-40"
      />
    </div>

    <UTable :data="data" :columns="columns" :loading="loading" class="w-full" />

    <div v-if="total > perPage" class="flex justify-end mt-4">
      <UPagination v-model:page="page" :items-per-page="perPage" :total="total" />
    </div>

    <USlideover v-model:open="slideoverOpen" :title="editingProduct ? 'Edit Product' : 'Add Product'">
      <template #body>
        <ProductsProductForm
          :product="editingProduct"
          @saved="onSaved"
          @cancel="slideoverOpen = false"
        />
      </template>
    </USlideover>
  </div>
</template>
