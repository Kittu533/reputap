<script setup lang="ts">
import type { Customer } from '~/utils/api'

definePageMeta({
  middleware: 'auth',
})

const { apiFetch } = useApi()
const toast = useToast()

const search = ref('')
const page = ref(1)
const perPage = ref(15)
const data = ref<Customer[]>([])
const total = ref(0)
const loading = ref(false)

// Delete Modal State
const deleteModalOpen = ref(false)
const customerToDelete = ref<Customer | null>(null)
const deleteLoading = ref(false)

async function fetchData() {
  loading.value = true
  try {
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
  }
  finally {
    loading.value = false
  }
}

function promptDelete(customer: Customer) {
  customerToDelete.value = customer
  deleteModalOpen.value = true
}

async function confirmDelete() {
  if (!customerToDelete.value) return
  deleteLoading.value = true
  try {
    const response = await apiFetch(`/api/v1/customers/${customerToDelete.value.id}`, {
      method: 'DELETE',
    })
    if (response.success) {
      toast.add({
        title: 'Customer Deleted',
        description: `Customer "${customerToDelete.value.name}" (${customerToDelete.value.customer_code}) has been deleted.`,
        color: 'success',
      })
      deleteModalOpen.value = false
      fetchData()
    }
    else {
      toast.add({
        title: 'Failed to Delete',
        description: response.error?.message || 'Customer could not be deleted',
        color: 'error',
      })
    }
  }
  finally {
    deleteLoading.value = false
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
</script>

<template>
  <div class="space-y-6">
    <!-- Header with Breadcrumbs & Action -->
    <CommonPageHeader
      title="Customers"
      description="Manage business clients, review locations, and stand ordering contacts."
      :breadcrumbs="[{ label: 'Customers' }]"
    >
      <template #actions>
        <CommonButton to="/customers/new" icon="i-lucide-plus">
          Add Customer
        </CommonButton>
      </template>
    </CommonPageHeader>

    <!-- Reusable Data Table -->
    <CommonDataTable
      v-model:search="search"
      v-model:page="page"
      :data="data"
      :total="total"
      :per-page="perPage"
      :loading="loading"
      search-placeholder="Search by name, code, email, city..."
      empty-title="No customers found"
      empty-description="There are no registered customers matching your criteria. Click below to add a new client."
      empty-action-label="Add Customer"
      @empty-action="navigateTo('/customers/new')"
    >
      <template #table>
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-100 bg-slate-50/70 text-[11px] font-bold uppercase tracking-wider text-slate-400">
              <th class="py-3.5 px-6">Code</th>
              <th class="py-3.5 px-6">Name</th>
              <th class="py-3.5 px-6">Business</th>
              <th class="py-3.5 px-6">Email</th>
              <th class="py-3.5 px-6">Phone</th>
              <th class="py-3.5 px-6">City</th>
              <th class="py-3.5 px-6 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-sm">
            <tr
              v-for="customer in data"
              :key="customer.id"
              class="hover:bg-slate-50/70 transition-colors"
            >
              <td class="py-4 px-6 font-mono text-xs font-bold text-blue-600">
                {{ customer.customer_code }}
              </td>
              <td class="py-4 px-6 font-semibold text-slate-900">
                {{ customer.name }}
              </td>
              <td class="py-4 px-6 text-slate-600">
                {{ customer.business_name || '-' }}
              </td>
              <td class="py-4 px-6 text-xs text-slate-500">
                {{ customer.email || '-' }}
              </td>
              <td class="py-4 px-6 text-xs text-slate-500">
                {{ customer.phone || '-' }}
              </td>
              <td class="py-4 px-6 text-xs text-slate-600">
                {{ customer.city || '-' }}
              </td>
              <td class="py-4 px-6 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <NuxtLink
                    :to="`/customers/${customer.id}/edit`"
                    class="size-8 rounded-lg text-slate-500 hover:text-blue-600 hover:bg-blue-50 flex items-center justify-center transition-colors"
                    title="Edit Customer"
                  >
                    <UIcon name="i-lucide-pencil" class="size-4" />
                  </NuxtLink>
                  <button
                    class="size-8 rounded-lg text-slate-500 hover:text-rose-600 hover:bg-rose-50 flex items-center justify-center transition-colors"
                    title="Delete Customer"
                    @click="promptDelete(customer)"
                  >
                    <UIcon name="i-lucide-trash-2" class="size-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </template>
    </CommonDataTable>

    <!-- Confirmation Delete Modal -->
    <CommonConfirmModal
      v-model:open="deleteModalOpen"
      title="Delete Customer"
      :description="`Are you sure you want to delete customer '${customerToDelete?.name}' (${customerToDelete?.customer_code})? All associated store and order history might be affected.`"
      confirm-text="Yes, Delete Customer"
      confirm-color="error"
      icon="i-lucide-alert-triangle"
      :loading="deleteLoading"
      @confirm="confirmDelete"
    />
  </div>
</template>
