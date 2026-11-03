<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { Store, ApiResponse } from '~/utils/api'

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
const data = ref<Store[]>([])
const total = ref(0)
const loading = ref(false)

const statusOptions = [
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
  { label: 'Suspended', value: 'suspended' },
]

const statusColor: Record<string, string> = {
  active: 'success',
  inactive: 'neutral',
  suspended: 'error',
}

const columns: TableColumn<Store>[] = [
  { accessorKey: 'store_code', header: 'Code' },
  { accessorKey: 'name', header: 'Name' },
  {
    accessorKey: 'category',
    header: 'Category',
    cell: ({ row }) => h(UBadge, { variant: 'subtle', color: 'neutral', class: 'capitalize' }, () => row.getValue('category')?.replace('_', ' ')),
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
    accessorKey: 'review_url',
    header: 'Review URL',
    cell: ({ row }) => {
      const url = row.original.review_url
      const verified = row.original.review_url_verified
      if (!url) return h('span', { class: 'text-muted' }, '—')
      return h(UBadge, { variant: 'subtle', color: verified ? 'success' : 'warning' }, () => verified ? 'Verified' : 'Not Verified')
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

  const response = await apiFetch<Store[]>('/api/v1/stores', { query })
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
const editingStore = ref<Store | null>(null)

function openCreate() {
  editingStore.value = null
  slideoverOpen.value = true
}

function openEdit(store: Store) {
  editingStore.value = store
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
      <h1 class="text-2xl font-bold">Stores</h1>
      <UButton icon="i-lucide-plus" label="Add Store" @click="openCreate" />
    </div>

    <div class="flex items-center gap-4 mb-4">
      <UInput
        v-model="search"
        placeholder="Search stores..."
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

    <USlideover v-model:open="slideoverOpen" :title="editingStore ? 'Edit Store' : 'Add Store'">
      <template #body>
        <StoresStoreForm
          :store="editingStore"
          @saved="onSaved"
          @cancel="slideoverOpen = false"
        />
      </template>
    </USlideover>
  </div>
</template>
