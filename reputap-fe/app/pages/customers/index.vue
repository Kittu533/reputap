<script setup lang="ts">
import { h, resolveComponent } from 'vue'
import type { TableColumn } from '@nuxt/ui'
import type { Customer } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()
const toast = useToast()

const UButton = resolveComponent('UButton')

const search = ref('')
const page = ref(1)
const perPage = ref(20)
const data = ref<Customer[]>([])
const total = ref(0)
const loading = ref(false)

const columns: TableColumn<Customer>[] = [
  { accessorKey: 'customer_code', header: 'Code' },
  { accessorKey: 'name', header: 'Name' },
  { accessorKey: 'business_name', header: 'Business' },
  { accessorKey: 'email', header: 'Email' },
  { accessorKey: 'phone', header: 'Phone' },
  { accessorKey: 'city', header: 'City' },
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
          onClick: (e: Event) => {
            e.stopPropagation()
            openEdit(row.original)
          },
        }),
        h(UButton, {
          icon: 'i-lucide-trash-2',
          color: 'error',
          variant: 'ghost',
          size: 'xs',
          onClick: (e: Event) => {
            e.stopPropagation()
            deleteCustomer(row.original)
          },
        }),
      ])
    },
  },
]

async function fetchData() {
  loading.value = true
  const response = await apiFetch<Customer[]>('/api/v1/customers', {
    query: {
      page: page.value,
      per_page: perPage.value,
      search: search.value || undefined,
    },
  })
  if (response.success) {
    data.value = response.data || []
    total.value = response.meta?.total || 0
  }
  loading.value = false
}

async function deleteCustomer(customer: Customer) {
  if (!confirm(`Delete customer "${customer.name}"?`)) return

  const response = await apiFetch(`/api/v1/customers/${customer.id}`, { method: 'DELETE' })
  if (response.success) {
    toast.add({ title: 'Deleted', description: `Customer ${customer.customer_code} deleted`, color: 'success' })
    fetchData()
  }
  else {
    toast.add({ title: 'Error', description: response.error?.message || 'Delete failed', color: 'error' })
  }
}

let searchTimeout: ReturnType<typeof setTimeout>
watch(search, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    fetchData()
  }, 300)
})

watch(page, () => fetchData())

onMounted(() => fetchData())

const slideoverOpen = ref(false)
const editingCustomer = ref<Customer | null>(null)

function openCreate() {
  editingCustomer.value = null
  slideoverOpen.value = true
}

function openEdit(customer: Customer) {
  editingCustomer.value = customer
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
      <h1 class="text-2xl font-bold">Customers</h1>
      <UButton icon="i-lucide-plus" label="Add Customer" @click="openCreate" />
    </div>

    <div class="flex items-center gap-4 mb-4">
      <UInput
        v-model="search"
        placeholder="Search customers..."
        icon="i-lucide-search"
        class="max-w-sm"
      />
    </div>

    <UTable
      :data="data"
      :columns="columns"
      :loading="loading"
      class="w-full"
    />

    <div v-if="total > perPage" class="flex justify-end mt-4">
      <UPagination
        v-model:page="page"
        :items-per-page="perPage"
        :total="total"
      />
    </div>

    <USlideover v-model:open="slideoverOpen" :title="editingCustomer ? 'Edit Customer' : 'Add Customer'">
      <template #body>
        <CustomersCustomerForm
          :customer="editingCustomer"
          @saved="onSaved"
          @cancel="slideoverOpen = false"
        />
      </template>
    </USlideover>
  </div>
</template>
