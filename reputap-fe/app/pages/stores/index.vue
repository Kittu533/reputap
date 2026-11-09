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

const deleteModalOpen = ref(false)
const deleteLoading = ref(false)
const storeToDelete = ref<Store | null>(null)

function promptDelete(store: Store) {
  storeToDelete.value = store
  deleteModalOpen.value = true
}

async function confirmDelete() {
  if (!storeToDelete.value) return
  deleteLoading.value = true
  try {
    const res: any = await apiFetch(`/api/v1/stores/${storeToDelete.value.id}`, {
      method: 'DELETE',
    })
    if (res.success) {
      toast.add({
        title: 'Store Deleted',
        description: `Store "${storeToDelete.value.name}" has been deleted.`,
        color: 'success',
      })
      deleteModalOpen.value = false
      fetchData()
    } else {
      toast.add({
        title: 'Delete Failed',
        description: res.error?.message || 'Could not delete store.',
        color: 'error',
      })
    }
  } catch (err: any) {
    toast.add({
      title: 'Delete Failed',
      description: err?.message || 'Error deleting store',
      color: 'error',
    })
  } finally {
    deleteLoading.value = false
  }
}

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
          to: `/stores/${row.original.id}/edit`,
        }),
        h(UButton, {
          icon: 'i-lucide-trash-2',
          color: 'error',
          variant: 'ghost',
          size: 'xs',
          onClick: () => promptDelete(row.original),
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
</script>

<template>
  <div class="space-y-6">
    <CommonPageHeader
      title="Stores"
      description="Manage customer branches, review target links, and physical store locations."
      :breadcrumbs="[{ label: 'Stores' }]"
    >
      <template #actions>
        <CommonButton to="/stores/new" icon="i-lucide-plus">
          Add Store
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
      search-placeholder="Search stores by name, code, category..."
      empty-title="No stores found"
      empty-description="There are no stores matching your criteria."
      empty-action-label="Add Store"
      @empty-action="navigateTo('/stores/new')"
    >
      <template #filters>
        <USelect
          v-model="statusFilter"
          :items="statusOptions"
          placeholder="All Status"
          class="w-40"
        />
      </template>
    </CommonDataTable>

    <!-- Confirmation Delete Modal -->
    <CommonConfirmModal
      v-model:open="deleteModalOpen"
      title="Delete Store"
      :description="`Are you sure you want to delete store '${storeToDelete?.name}' (${storeToDelete?.store_code})?`"
      confirm-text="Yes, Delete Store"
      confirm-color="error"
      icon="i-lucide-alert-triangle"
      :loading="deleteLoading"
      @confirm="confirmDelete"
    />
  </div>
</template>
